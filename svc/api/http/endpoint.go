package http

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/piusalfred/kitsvc"
	"github.com/piusalfred/kitsvc/svc"
)

// CountRequest collects the request parameters for the Count method.
type CountRequest struct {
	String string `json:"string"`
}

// CountResponse collects the response parameters for the Count method.
type CountResponse struct {
	Count int `json:"count"`
	Err   error `json:"err,omitempty"`
}

// MakeCountEndpoint returns an endpoint that invokes Count on the service.
func MakeCountEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CountRequest)
		i0, e1 := s.Count(ctx, req.String)
		return CountResponse{
			Err:   e1,
			Count: i0,
		}, nil
	}
}

// Failed implements Failer.
func (r CountResponse) Failed() error {
	return r.Err
}

// UppercaseRequest collects the request parameters for the Uppercase method.
type UppercaseRequest struct {
	String string `json:"string"`
}

// UppercaseResponse collects the response parameters for the Uppercase method.
type UppercaseResponse struct {
	String string `json:"string,omitempty"`
	Err    error  `json:"err,omitempty"`
}

// MakeUppercaseEndpoint returns an endpoint that invokes Uppercase on the service.
func MakeUppercaseEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UppercaseRequest)
		s0, e1 := s.Uppercase(ctx, req.String)
		return UppercaseResponse{
			Err:    e1,
			String: s0,
		}, nil
	}
}

// Failed implements Failer.
func (r UppercaseResponse) Failed() error {
	return r.Err
}

// VersionRequest collects the request parameters for the Version method.
type VersionRequest struct{}

// VersionResponse collects the response parameters for the Version method.
type VersionResponse struct {
	Version kitsvc.Version `json:"version"`
	Err     error          `json:"err"`
}

// MakeVersionEndpoint returns an endpoint that invokes Version on the service.
func MakeVersionEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		k0, e1 := s.Version(ctx)
		return VersionResponse{
			Err:     e1,
			Version: k0,
		}, nil
	}
}

// Failed implements Failer.
func (r VersionResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	CountEndpoint     endpoint.Endpoint
	UppercaseEndpoint endpoint.Endpoint
	VersionEndpoint   endpoint.Endpoint
}

func MakeEndpoints(s svc.Service) Endpoints {
	return Endpoints{
		CountEndpoint:     MakeCountEndpoint(s),
		UppercaseEndpoint: MakeUppercaseEndpoint(s),
		VersionEndpoint:   MakeVersionEndpoint(s),
	}
}

// Count implements Service. Primarily useful in a client.
func (e Endpoints) Count(ctx context.Context, string string) (i0 int, e1 error) {
	request := CountRequest{String: string}
	response, err := e.CountEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(CountResponse).Count, response.(CountResponse).Err
}

// Uppercase implements Service. Primarily useful in a client.
func (e Endpoints) Uppercase(ctx context.Context, string2 string) (s0 string, e1 error) {
	request := UppercaseRequest{String: string2}
	response, err := e.UppercaseEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(UppercaseResponse).String, response.(UppercaseResponse).Err
}

// Version implements Service. Primarily useful in a client.
func (e Endpoints) Version(ctx context.Context) (k0 kitsvc.Version, e1 error) {
	request := VersionRequest{}
	response, err := e.VersionEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(VersionResponse).Version, response.(VersionResponse).Err
}
