package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/piusalfred/kitsvc/svc"
	svchttp "github.com/piusalfred/kitsvc/svc/api/http"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// New returns an AddService backed by an HTTP server living at the remote
// instance. We expect instance to come from a service discovery system, so
// likely of the form "host:port".
func NewClient(instance string, options map[string][]kithttp.ClientOption) (svc.Service, error) {

	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	var countEndpoint endpoint.Endpoint
	{
		countEndpoint = kithttp.NewClient(
			"POST", copyURL(u, "/count"),
			encodeHTTPGenericRequest,
			decodeCountResponse,
			options["Count"]...).Endpoint()
	}

	var uppercaseEndpoint endpoint.Endpoint
	{
		uppercaseEndpoint = kithttp.NewClient("POST", copyURL(u, "/uppercase"), encodeHTTPGenericRequest, decodeUppercaseResponse, options["Uppercase"]...).Endpoint()
	}

	var versionEndpoint endpoint.Endpoint
	{
		versionEndpoint = kithttp.NewClient("GET", copyURL(u, "/version"), encodeHTTPGenericRequest, decodeVersionResponse, options["Version"]...).Endpoint()
	}

	return svchttp.Endpoints{
		CountEndpoint:     countEndpoint,
		UppercaseEndpoint: uppercaseEndpoint,
		VersionEndpoint:   versionEndpoint,
	}, nil
}

// EncodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// SON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeCountResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeCountResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, svchttp.ErrorDecoder(r)
	}
	var resp svchttp.CountResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeUppercaseResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeUppercaseResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, svchttp.ErrorDecoder(r)
	}
	var resp svchttp.UppercaseResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeVersionResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeVersionResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, svchttp.ErrorDecoder(r)
	}
	var resp svchttp.VersionResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
func copyURL(base *url.URL, path string) (next *url.URL) {
	n := *base
	n.Path = path
	next = &n
	return
}