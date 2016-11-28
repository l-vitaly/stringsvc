package stringsvc_test

import (
	"testing"

	"google.golang.org/grpc"

	"github.com/l-vitaly/stringsvc/stringsvc_grpc"
	"golang.org/x/net/context"
)

const DEFAULT_SERVER_ADDR = ":8082"

func TestGRPCServer(t *testing.T) {
	conn, err := grpc.Dial(DEFAULT_SERVER_ADDR, grpc.WithInsecure())
	if err != nil {
		t.Errorf("Fail dial \"%s\" gRPC service: %s", DEFAULT_SERVER_ADDR, err)
		return
	}

	svc := stringsvc_grpc.NewClient(conn)
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
