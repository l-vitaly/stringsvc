package stringsvc

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints.

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// Endpoints collects for String service.
type Endpoints struct {
	UppercaseEndpoint endpoint.Endpoint
}

// Uppercase implement Service.
func (e Endpoints) Uppercase(ctx context.Context, s string) (string, error) {
	req := UppercaseRequest{S: s}
	resp, err := e.UppercaseEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.(UppercaseResponse).V, resp.(UppercaseResponse).Err
}

// MakeUppercaseEndpoint returns an endpoint that invokes Uppercase on the service.
func MakeUppercaseEndpoint(svc StringSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UppercaseRequest)
		v, err := svc.Uppercase(ctx, req.S)
		return UppercaseResponse{V: v, Err: err}, nil
	}
}
