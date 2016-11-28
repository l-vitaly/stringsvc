package main

import (
	"flag"
	"net"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"log"

	"github.com/l-vitaly/stringsvc"
	"github.com/l-vitaly/stringsvc/stringsvc_grpc"
	pb "github.com/l-vitaly/stringsvc/stringsvc_pb"
	"github.com/spf13/viper"
)

const APP_NAME = "stringsvc"

var (
	githash    = "dev"
	buildstamp = time.Now().Format(time.RFC822)
	configName = flag.String("config-name", APP_NAME+"-config", "Set config file name")
	host       = flag.String("host", "", "")
	port       = flag.String("port", "8082", "")
)

func init() {
	flag.Parse()

	viper.SetConfigName(*configName)
	viper.AddConfigPath("/etc/" + APP_NAME + "/")
	viper.AddConfigPath("$HOME/." + APP_NAME)
	viper.AddConfigPath(".")
}

func main() {
	log.Println("Version:", githash, "Build:", buildstamp)

	svc := stringsvc.NewService()
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

	log.Println("Server listen:", addr)

	errc <- s.Serve(ln)
}
