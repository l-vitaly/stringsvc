package transportgrpc

// Encode request is a transport/grpc.EncodeRequestFunc that converts a
// user-domain request to a gRPC request. Primarily useful in a client.

import (
	"github.com/l-vitaly/eutils"
	"github.com/l-vitaly/stringsvc"
	"github.com/l-vitaly/stringsvc/pb"
	"golang.org/x/net/context"
)

func encodeUppercaseRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(stringsvc.UppercaseRequest)
	return &pb.UppercaseRequest{S: req.S}, nil
}

func encodeUppercaseResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(stringsvc.UppercaseResponse)
	return &pb.UppercaseResponse{V: resp.V, Err: eutils.Err2Str(resp.Err)}, nil
}
