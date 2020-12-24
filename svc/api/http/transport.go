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
	var request UppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	if request.String == ""{
		return nil, errors.New("invalid request body")
	}
	return request, nil
}

func decodeVersionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request VersionRequest
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

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
func ErrorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

// This is used to set the http status, see an example here :
// https://github.com/go-kit/kit/blob/master/examples/addsvc/pkg/addtransport/http.go#L133
func err2code(err error) int {
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}

func MakeHandler(service svc.Service) http.Handler {

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	upperCaseHandler := kithttp.NewServer(
		MakeUppercaseEndpoint(service),
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

	versionHandler := kithttp.NewServer(
		MakeVersionEndpoint(service),
		decodeVersionRequest,
		encodeResponse,
		opts ...)

	r.Post("/svc/uppercase", upperCaseHandler)

	r.Post("/svc/count", countHandler)

	r.Get("/svc/version", versionHandler)

	return r
}
