package main

import (
	"flag"
	"net"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/l-vitaly/stringsvc"
	"github.com/l-vitaly/stringsvc/stringsvc_grpc"
	pb "github.com/l-vitaly/stringsvc/stringsvc_pb"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	githash    = "dev"
	buildstamp = time.Now().Format(time.RFC822)
	zipkinAddr = flag.String("zipkin.addr", "", "Zipkin tracing via host:port")
	consulAddr = flag.String("consulAddr", "localhost:8301", "Consul addr via host:port")
	host       = flag.String("host", "", "")
	port       = flag.String("port", "8082", "")
)

func init() {
	flag.Parse()
}

func main() {
	var logger log.Logger

	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
	logger = log.NewContext(logger).With("caller", log.DefaultCaller)

	logger.Log("version", githash)
	logger.Log("builddate", buildstamp)
	logger.Log("msg", "hello")
	defer logger.Log("msg", "goodbye")

	if *consulAddr == "" {
		logger.Log("err", "consul addr is required")
		os.Exit(1)
	}

	viper.AddRemoteProvider("consul", *consulAddr, "/config/config.json")
	viper.SetConfigType("json")

	err := viper.ReadRemoteConfig()
	if err != nil {
		logger.Log("err", "consul error: "+err.Error())
		os.Exit(1)
	}

	if *zipkinAddr == "" {
		logger.Log("err", "zipkin addr is required")
		os.Exit(1)
	}

	// Metrics.
	duration := prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "stringsvc",
		Name:      "request_duration_ns",
		Help:      "Request duration in nanoseconds.",
	}, []string{"method", "success"})

	// Tracing.
	collector, err := zipkin.NewKafkaCollector(
		strings.Split(*zipkinAddr, ","),
		zipkin.KafkaLogger(
			log.NewContext(logger).With("tracer", "Zipkin"),
		),
	)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "localhost:80", "stringsvc"),
	)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	svc := stringsvc.NewService()

	uppercaseDuration := duration.With("method", "Uppercase")
	uppercaseLogger := log.NewContext(logger).With("method", "Uppercase")

	uppercaseEndpoint := stringsvc.MakeUppercaseEndpoint(svc)
	uppercaseEndpoint = opentracing.TraceServer(tracer, "Uppercase")(uppercaseEndpoint)
	uppercaseEndpoint = stringsvc.EndpointInstrumentingMiddleware(uppercaseDuration)(uppercaseEndpoint)
	uppercaseEndpoint = stringsvc.EndpointLoggingMiddleware(uppercaseLogger)(uppercaseEndpoint)

	endpoints := stringsvc.Endpoints{
		UppercaseEndpoint: stringsvc.MakeUppercaseEndpoint(svc),
	}

	ctx := context.Background()
	srv := stringsvc_grpc.NewServer(ctx, endpoints)

	errc := make(chan error)

	addr := *host + ":" + *port

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		errc <- err
		return
	}

	s := grpc.NewServer()
	pb.RegisterStringServer(s, srv)

	errc <- s.Serve(ln)
}
