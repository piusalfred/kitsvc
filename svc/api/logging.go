package api

import (
	"github.com/go-kit/kit/log"
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

func (l loggingMiddleware) UpperCase(s string) (output string, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = l.svc.UpperCase(s)
	return
}

func (l loggingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = l.svc.Count(s)
	return
}
