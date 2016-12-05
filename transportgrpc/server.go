package transportgrpc

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/l-vitaly/stringsvc"
	"github.com/l-vitaly/stringsvc/pb"
	stdopentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

type server struct {
	uppercase grpctransport.Handler
}

// NewServer makes a set of endpoints.
func NewServer(ctx context.Context, endpoints stringsvc.Endpoints,
	tracer stdopentracing.Tracer, logger log.Logger) pb.StringServer {
	defaultOpt := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	uppercaseServerBefore := grpctransport.ServerBefore(
		opentracing.FromGRPCRequest(tracer, "Uppercase", logger),
	)
	return &server{
		uppercase: grpctransport.NewServer(
			ctx,
			endpoints.UppercaseEndpoint,
			decodeUppercaseRequest,
			encodeUppercaseResponse,
			append(defaultOpt, uppercaseServerBefore)...,
		),
	}
}

func (s *server) Uppercase(ctx context.Context, req *pb.UppercaseRequest) (*pb.UppercaseResponse, error) {
	_, resp, err := s.uppercase.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UppercaseResponse), nil
}
