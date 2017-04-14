package elastic

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/signer/v4"
)

const SERVICE_NAME = "es"

type AWSSigningRoundTripper struct {
	region      string
	credentials *credentials.Credentials
	inner       http.RoundTripper
}

func NewAWSSigningRoundTripper(inner http.RoundTripper, region string, credentials *credentials.Credentials) *AWSSigningRoundTripper {
	if inner == nil {
		inner = http.DefaultTransport
	}
	p := &AWSSigningRoundTripper{inner: inner, region: region, credentials: credentials}
	return p
}

func (tr *AWSSigningRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	var err error

	if req.Body != nil {
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	if req.Method == "GET" || req.Method == "HEAD" {
		delete(req.Header, "Content-Length")
	}

	oldPath := req.URL.Path
	if oldPath != "" {
		// Escape the path before signing so that the path in the signature and
		// the path in the request match.
		req.URL.Path = req.URL.EscapedPath()
	}

	awsReq := &request.Request{}
	awsReq.Config.Credentials = tr.credentials
	awsReq.Config.Region = aws.String(tr.region)
	awsReq.ClientInfo.ServiceName = SERVICE_NAME
	awsReq.HTTPRequest = req
	awsReq.Time = time.Now()
	awsReq.ExpireTime = 0

	if body != nil {
		awsReq.Body = bytes.NewReader(body)
	}

	v4.Sign(awsReq)

	if awsReq.Error != nil {
		return nil, awsReq.Error
	}

	req.URL.Path = oldPath
	if body != nil {
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	res, err := tr.inner.RoundTrip(req)

	if err != nil {
		return nil, err
	} else {
		return res, err
	}
}
