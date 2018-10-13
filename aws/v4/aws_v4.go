// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package v4

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

// NewV4SigningClient returns an *http.Client that will sign all requests with AWS V4 Signing.
func NewV4SigningClient(creds *credentials.Credentials, region string) *http.Client {
	return NewV4SigningClientWithHTTPClient(creds, region, http.DefaultClient)
}

// NewV4SigningClientWithHTTPClient returns an *http.Client that will sign all requests with AWS V4 Signing.
func NewV4SigningClientWithHTTPClient(creds *credentials.Credentials, region string, httpClient *http.Client) *http.Client {
	return &http.Client{
		Transport: Transport{
			client: httpClient,
			creds:  creds,
			signer: v4.NewSigner(creds),
			region: region,
		},
	}
}

// Transport is a RoundTripper that will sign requests with AWS V4 Signing
type Transport struct {
	client *http.Client
	creds  *credentials.Credentials
	signer *v4.Signer
	region string
}

// RoundTrip uses the underlying RoundTripper transport, but signs request first with AWS V4 Signing
func (st Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// AWS signer needs an io.ReadSeeker; however, req.Body is an io.ReadCloser.
	// TODO Maybe there's a more efficient way to get an io.ReadSeeker than to read the whole thing.
	var body io.ReadSeeker
	if req.Body != nil {
		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(d)
	}
	_, err := st.signer.Sign(req, body, "es", st.region, time.Now())
	if err != nil {
		return nil, err
	}
	return st.client.Do(req)
}
