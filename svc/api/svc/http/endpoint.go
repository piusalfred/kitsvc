package http

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/piusalfred/kitsvc/svc"
)


func MakeUpperCaseEndpoint(service svc.Service)endpoint.Endpoint {

	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(upperCaseReq)
		v,err := service.UpperCase(req.S)

		if err != nil {
			return upperCaseResp{
				V:   v,
				Err: err.Error(),
			}, nil
		}

		return upperCaseResp{
			V:   v,
			Err: "",
		}, nil
	}
}


func MakeCountEndpoint(service svc.Service) endpoint.Endpoint{
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(countReq)
		v := service.Count(req.S)

		return countResp{V: v},nil
	}
}
