package transportgrpc

// Decode request is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.

import (
	"github.com/l-vitaly/eutils"
	"github.com/l-vitaly/stringsvc"
	"github.com/l-vitaly/stringsvc/pb"
	"golang.org/x/net/context"
)

func decodeUppercaseRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UppercaseRequest)
	return stringsvc.UppercaseRequest{S: req.S}, nil

}

func decodeUppercaseResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*pb.UppercaseResponse)
	return stringsvc.UppercaseResponse{V: resp.V, Err: eutils.Str2Err(resp.Err)}, nil
}
