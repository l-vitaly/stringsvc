package stringsvc

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints.

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/l-vitaly/eutils"
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

// EndpointInstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func EndpointInstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// EndpointLoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log("error", eutils.Err2Str(err), "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// UppercaseRequest struct for Uppercase request
type UppercaseRequest struct {
	S string `json:"s"`
}

// UppercaseResponse struct for Uppercase response
type UppercaseResponse struct {
	V   string `json:"v"`
	Err error  `json:"err"`
}
