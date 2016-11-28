package stringsvc_grpc

import (
	pb "github.com/l-vitaly/stringsvc/stringsvc_pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/l-vitaly/stringsvc"
	"google.golang.org/grpc"
)

func NewClient(conn *grpc.ClientConn) stringsvc.StringService {
	uppercaseEndpoint := grpctransport.NewClient(
		conn, "String", "Uppercase", encodeUppercaseRequest, decodeUppercaseResponse, pb.UppercaseResponse{},
	).Endpoint()

	return stringsvc.Endpoints{
		UppercaseEndpoint: uppercaseEndpoint,
	}
}
