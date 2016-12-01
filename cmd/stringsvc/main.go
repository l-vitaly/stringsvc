package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"

	"github.com/l-vitaly/consul"
	"github.com/l-vitaly/stringsvc"
	pb "github.com/l-vitaly/stringsvc/stringsvcpb"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"

	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"github.com/l-vitaly/stringsvc/stringsvcgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	githash    = "dev"
	buildstamp = time.Now().Format(time.RFC822)
	debugAddr  = flag.String("debug.addr", "localhost:62101", "Debug addr via host:port")
	zipkinAddr = flag.String("zipkin.addr", "localhost:9411", "Zipkin tracing via host:port")
	consulAddr = flag.String("consul.addr", "localhost:8400", "Consul addr via host:port")
	svcAddr    = flag.String("addr", "localhost:62001", "Service addr via host:port")
)

const (
	SERVICE_NAME = "stringsvc"
	DEBUG        = true
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
}

func main() {
	var logger log.Logger

	logger = log.NewJSONLogger(os.Stdout)
	logger = log.NewContext(logger).With("@timestamp", log.DefaultTimestampUTC)
	logger = log.NewContext(logger).With("@message", "info")
	logger = log.NewContext(logger).With("caller", log.DefaultCaller)

	logger.Log("version", githash)
	logger.Log("builddate", buildstamp)
	logger.Log("msg", "hello")
	defer logger.Log("msg", "goodbye")

	if *svcAddr == "" {
		logger.Log("err", "service addr is required")
		os.Exit(1)
	}

	if *consulAddr == "" {
		logger.Log("err", "consul addr is required")
		os.Exit(1)
	}

	if *zipkinAddr == "" {
		logger.Log("err", "zipkin addr is required")
		os.Exit(1)
	}

	consul, err := consul.NewConsulClient(*consulAddr)
	if err != nil {
		logger.Log("err", err)
	}

	err = consul.RegisterService(SERVICE_NAME, *svcAddr)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	// Metrics.
	duration := prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: SERVICE_NAME,
		Name:      "request_duration_ns",
		Help:      "Request duration in nanoseconds.",
	}, []string{"method", "success"})

	// Tracing.
	collector, err := zipkin.NewHTTPCollector(fmt.Sprintf("http://%s/api/v1/spans", *zipkinAddr))
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	recorder := zipkin.NewRecorder(collector, DEBUG, *svcAddr, SERVICE_NAME)
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	svc := stringsvc.NewService()
	svc = stringsvc.ServiceLoggingMiddleware(logger)(svc)
	svc = stringsvc.ServiceInstrumentingMiddleware()(svc)

	uppercaseDuration := duration.With("method", "Uppercase")
	uppercaseLogger := log.NewContext(logger).With("method", "Uppercase")

	uppercaseEndpoint := stringsvc.MakeUppercaseEndpoint(svc)
	uppercaseEndpoint = opentracing.TraceServer(tracer, "Uppercase")(uppercaseEndpoint)
	uppercaseEndpoint = stringsvc.EndpointInstrumentingMiddleware(uppercaseDuration)(uppercaseEndpoint)
	uppercaseEndpoint = stringsvc.EndpointLoggingMiddleware(uppercaseLogger)(uppercaseEndpoint)

	endpoints := stringsvc.Endpoints{
		UppercaseEndpoint: uppercaseEndpoint,
	}

	ctx := context.Background()
	errc := make(chan error)

	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Debug handler.
	go func() {
		logger := log.NewContext(logger).With("transport", "debug")

		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		m.Handle("/metrics", stdprometheus.Handler())

		logger.Log("addr", *debugAddr)
		errc <- http.ListenAndServe(*debugAddr, m)
	}()

	go func() {
		logger := log.NewContext(logger).With("transport", "gRPC")

		listener, err := net.Listen("tcp", *svcAddr)
		if err != nil {
			errc <- err
			return
		}

		srv := stringsvcgrpc.NewServer(ctx, endpoints, tracer, logger)
		s := grpc.NewServer()
		pb.RegisterStringServer(s, srv)

		logger.Log("addr", *svcAddr)
		errc <- s.Serve(listener)
	}()

	logger.Log("exit", <-errc)
}
