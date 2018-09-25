// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package v4

import (
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
		Transport: V4Transport{
			client: httpClient,
			creds:  creds,
			signer: v4.NewSigner(creds),
			region: region,
		},
	}
}

// V4Transport is a RoundTripper that will sign requests with AWS V4 Signing
type V4Transport struct {
	client *http.Client
	creds  *credentials.Credentials
	signer *v4.Signer
	region string
}

// RoundTrip uses the underlying RoundTripper transport, but signs request first with AWS V4 Signing
func (st V4Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	_, err := st.signer.Sign(req, nil, "es", st.region, time.Unix(0, 0))
	if err != nil {
		return nil, err
	}
	return st.client.Do(req)
}
