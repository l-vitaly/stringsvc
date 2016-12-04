package transportgrpc

import (
	"github.com/l-vitaly/stringsvc/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/l-vitaly/stringsvc"
	"google.golang.org/grpc"
)

func NewClient(conn *grpc.ClientConn) stringsvc.StringSvc {
	uppercaseEndpoint := grpctransport.NewClient(
		conn, "String", "Uppercase", encodeUppercaseRequest, decodeUppercaseResponse, pb.UppercaseResponse{},
	).Endpoint()

	return stringsvc.Endpoints{
		UppercaseEndpoint: uppercaseEndpoint,
	}
}
