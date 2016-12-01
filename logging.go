package stringsvc

import (
	"time"

	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

func ServiceLoggingMiddleware(logger log.Logger) Middleware {
	return func(next StringSvc) StringSvc {
		return serviceLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   StringSvc
}

func (mw serviceLoggingMiddleware) Uppercase(ctx context.Context, s string) (v string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Uppercase",
			"s", s, v, "error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Uppercase(ctx, s)
}
