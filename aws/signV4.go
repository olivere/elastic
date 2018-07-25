package aws

import (
	"net/http"

	"github.com/smartystreets/go-aws-auth"
)

// NewV4SigningClient returns an *http.Client that will sign all requests with AWS V4 Signing.
func NewV4SigningClient(credentials awsauth.Credentials) *http.Client {
	return &http.Client{
		Transport: V4Transport{
			HTTPClient:  http.DefaultClient,
			Credentials: credentials,
		},
	}
}

// V4Transport is a RoundTripper that will sign requests with AWS V4 Signing
type V4Transport struct {
	HTTPClient  *http.Client
	Credentials awsauth.Credentials
}

// RoundTrip uses the underlying RoundTripper transport, but signs request first with AWS V4 Signing
func (st V4Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Instead of directly modifying the request then calling http.DefaultTransport,
	// instead restart the request with the HTTPClient.Do function,
	// because the HTTPClient includes safeguards around not forwarding the
	// signed Authorization header to untrusted domains.
	return st.HTTPClient.Do(awsauth.Sign4(req, st.Credentials))
}
