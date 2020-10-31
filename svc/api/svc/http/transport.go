package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-zoo/bone"
	"github.com/piusalfred/kitsvc/svc"
	"io"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)


const (
	contentType = "application/json"
)

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	errInvalidQueryParams     = errors.New("invalid query params")
)

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request upperCaseReq
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countReq
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case svc.ErrStringEmpty:
		w.WriteHeader(http.StatusBadRequest)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case errInvalidQueryParams:
		w.WriteHeader(http.StatusBadRequest)
	case io.ErrUnexpectedEOF:
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		w.WriteHeader(http.StatusBadRequest)
	default:
		switch err.(type) {
		case *json.SyntaxError:
			w.WriteHeader(http.StatusBadRequest)
		case *json.UnmarshalTypeError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}


func MakeHandler(service svc.Service) http.Handler{

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	upperCaseHandler := kithttp.NewServer(
		MakeUpperCaseEndpoint(service),
		decodeUppercaseRequest,
		encodeResponse,
		opts...,
	)

	countHandler := kithttp.NewServer(
		MakeCountEndpoint(service),
		decodeCountRequest,
		encodeResponse,
		opts...,
	)

	r.Post("/svc/uppercase",upperCaseHandler)

	r.Post("/svc/count",countHandler)


	return r
}
