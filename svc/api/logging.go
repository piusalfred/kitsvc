package api

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/piusalfred/kitsvc"
	"github.com/piusalfred/kitsvc/svc"
	"time"
)

var _ svc.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    svc.Service
}

func LoggingMiddleware(svc svc.Service, logger log.Logger) svc.Service {
	return &loggingMiddleware{logger, svc}
}

func (l loggingMiddleware) Count(ctx context.Context, string string) (count int, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "count",
			"input", string,
			"output", count,
			"error",err,
			"took", time.Since(begin),
		)
	}(time.Now())

	count,err = l.svc.Count(ctx,string)
	return
}

func (l loggingMiddleware) Uppercase(ctx context.Context, string2 string) (str string, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "uppercase",
			"input", string2,
			"output", str,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	str, err = l.svc.Uppercase(ctx,string2)
	return
}

func (l loggingMiddleware) Version(ctx context.Context) (version kitsvc.Version, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "version",
			"output", version,
			"took", time.Since(begin),
		)
	}(time.Now())

	version,err = l.svc.Version(ctx)
	return
}
