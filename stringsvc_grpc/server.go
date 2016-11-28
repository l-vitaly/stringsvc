package stringsvc_grpc

import (
	"golang.org/x/net/context"

	"github.com/go-kit/kit/transport/grpc"
	"github.com/l-vitaly/eutils"
	"github.com/l-vitaly/stringsvc"
	pb "github.com/l-vitaly/stringsvc/stringsvc_pb"
)

type gRPCServer struct {
	uppercase grpc.Handler
}

// NewServer makes a set of endpoints.
func NewServer(ctx context.Context, endpoints stringsvc.Endpoints) *gRPCServer {
	return &gRPCServer{
		uppercase: grpc.NewServer(
			ctx,
			endpoints.UppercaseEndpoint,
			decodeUppercaseRequest,
			encodeUppercaseResponse,
		),
	}
}

func (s *gRPCServer) Uppercase(ctx context.Context, req *pb.UppercaseRequest) (*pb.UppercaseResponse, error) {
	_, resp, err := s.uppercase.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UppercaseResponse), nil
}

// decodeUppercaseRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeUppercaseRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UppercaseRequest)
	return stringsvc.UppercaseRequest{S: req.S}, nil
}

// encodeUppercaseRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain request to a gRPC request. Primarily useful in a client.
func encodeUppercaseRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(stringsvc.UppercaseRequest)
	return &pb.UppercaseRequest{S: req.S}, nil
}

// decodeUppercaseResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC response to a user-domain response. Primarily useful in a client.
func decodeUppercaseResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*pb.UppercaseResponse)
	return stringsvc.UppercaseResponse{V: resp.V, Err: eutils.Str2Err(resp.Err)}, nil
}

// encodeUppercaseResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC response. Primarily useful in a server.
func encodeUppercaseResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(stringsvc.UppercaseResponse)
	return &pb.UppercaseResponse{V: resp.V, Err: eutils.Err2Str(resp.Err)}, nil
}
