package test

import (
	"flag"
	"testing"

	"golang.org/x/net/context"

	"github.com/l-vitaly/stringsvc/transportgrpc"
	"google.golang.org/grpc"
)

var testAddr = flag.String("test.addr", "", "")

func TestGRPCServer(t *testing.T) {
	flag.Parse()

	conn, err := grpc.Dial(*testAddr, grpc.WithInsecure())
	if err != nil {
		t.Errorf("Fail dial \"%s\" gRPC service: %s", *testAddr, err)
		return
	}

	svc := transportgrpc.NewClient(conn)
	for _, tt := range uppercaseProvider {
		actual, err := svc.Uppercase(context.Background(), tt.value)
		if err != nil {
			t.Errorf("Uppercase: %s", err)
			return
		}
		if actual != tt.expected {
			t.Errorf("Uppercase(%s): expected %s, actual %s", tt.value, tt.expected, actual)
			return
		}
	}
}
