// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package v4

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
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
	if h, ok := req.Header["Authorization"]; ok && len(h) > 0 && strings.HasPrefix(h[0], "AWS4") {
		// Received a signed request, just pass it on.
		return st.client.Do(req)
	}

	if strings.Contains(req.URL.RawPath, "%2C") {
		// Escaping path
		req.URL.RawPath = url.PathEscape(req.URL.RawPath)
	}

	now := time.Now().UTC()
	req.Header.Set("Date", now.Format(time.RFC3339))

	var err error
	switch req.Body {
	case nil:
		_, err = st.signer.Sign(req, nil, "es", st.region, now)
	default:
		switch body := req.Body.(type) {
		case io.ReadSeeker:
			_, err = st.signer.Sign(req, body, "es", st.region, now)
		default:
			buf, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			req.Body = ioutil.NopCloser(bytes.NewReader(buf))
			_, err = st.signer.Sign(req, bytes.NewReader(buf), "es", st.region, time.Now().UTC())
		}
	}
	if err != nil {
		return nil, err
	}
	return st.client.Do(req)
}
