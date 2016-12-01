package stringsvc

import "golang.org/x/net/context"

func ServiceInstrumentingMiddleware() Middleware {
	return func(next StringSvc) StringSvc {
		return serviceInstrumentingMiddleware{
			next: next,
		}
	}
}

type serviceInstrumentingMiddleware struct {
	next StringSvc
}

func (mw serviceInstrumentingMiddleware) Uppercase(ctx context.Context, s string) (string, error) {
	v, err := mw.next.Uppercase(ctx, s)
	return v, err
}
