package middleware

import (
	"time"

	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
    "github.com/l-vitaly/stringsvc"
)

// ServiceLogging returns a service middleware that logs the
// parameters and result of each method invocation.
func ServiceLogging(logger log.Logger) Middleware {
	return func(next stringsvc.StringSvc) stringsvc.StringSvc {
		return serviceLogging{
			logger: logger,
			next:   next,
		}
	}
}

type serviceLogging struct {
	logger log.Logger
	next   stringsvc.StringSvc
}

func (mw serviceLogging) Uppercase(ctx context.Context, s string) (v string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Uppercase", "s", s, v, "error", err, "took", time.Since(begin))
	}(time.Now())
	return mw.next.Uppercase(ctx, s)
}
