package middleware

import (
    "golang.org/x/net/context"
    "github.com/l-vitaly/stringsvc"
)

// ServiceInstrumenting returns a service middleware that instruments
// the collect other metrics.
func ServiceInstrumenting() Middleware {
	return func(next stringsvc.StringSvc) stringsvc.StringSvc {
		return serviceInstrumenting{
			next: next,
		}
	}
}

type serviceInstrumenting struct {
	next stringsvc.StringSvc
}

func (mw serviceInstrumenting) Uppercase(ctx context.Context, s string) (string, error) {
	v, err := mw.next.Uppercase(ctx, s)
	return v, err
}
