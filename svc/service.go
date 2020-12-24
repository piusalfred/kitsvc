package svc

import (
	"context"
	"github.com/piusalfred/kitsvc"
	"strings"
)




type Service interface {
	Count(ctx context.Context, string string) (int, error)
	Uppercase(ctx context.Context, string2 string) (string, error)
	Version(ctx context.Context) (kitsvc.Version, error)
}

type svcImpl struct {
	username string
}

var _ Service = (*svcImpl)(nil)

func New(username string) Service {
	return &svcImpl{username: username}
}

func (svc *svcImpl) Count(ctx context.Context, string string) (int, error) {

	return len(string), nil
}

func (svc *svcImpl) Uppercase(ctx context.Context, s string) (string, error) {
	if s == "" {
		return "", ErrStringEmpty
	}

	return strings.ToUpper(s), nil
}

func (svc *svcImpl) Version(ctx context.Context) (kitsvc.Version, error) {
	return kitsvc.Version{
		Name:    "registry",
		Number:  "v 0.1.0",
		CdeName: "RoadRanger",
	},nil
}






