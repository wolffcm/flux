package dependenciestest

import (
	"io/ioutil"
	"net/http"

	"github.com/influxdata/flux/dependencies"
	"github.com/influxdata/flux/mock"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

var StatusOK int = 200

func defaultTestFunction(req *http.Request) *http.Response {
	body := (*req).Body
	// Test request parameters
	return &http.Response{
		StatusCode: StatusOK,
		Status:     "Body generated by test client",

		// Send response to be tested
		Body: ioutil.NopCloser(body),

		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}
}

type defaultDependencies struct {
	httpclient    *http.Client
	secretservice dependencies.SecretService
}

func (d defaultDependencies) HTTPClient() (*http.Client, error) {
	return d.httpclient, nil
}

func (d defaultDependencies) SecretService() (dependencies.SecretService, error) {
	return d.secretservice, nil
}

func NewTestDependenciesInterface() dependencies.Interface {
	return defaultDependencies{
		httpclient: &http.Client{
			Transport: RoundTripFunc(defaultTestFunction),
		},
		secretservice: mock.SecretService{},
	}
}