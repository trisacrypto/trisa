package openvasp

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Client is used to make HTTPS requests with mTLS configured to TRP servers.
type Client struct {
	http.Client
}

func NewClient() *Client {
	cookies, _ := cookiejar.New(nil)
	return &Client{
		Client: http.Client{
			// TODO: update transport with mTLS credentials
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           cookies,
			Timeout:       1 * time.Minute,
		},
	}
}

// Post a request to the specified Address, which can be a Travel Address, LNURL, or
// plain-text URL. The content type is automatically set to application/json and the
// body should read JSON data. Required OpenVASP headers are set by the TRPInfo struct.
// If the APIVersion is omitted, the default APIVersion is used, similarly if no request
// identifier is present, one is generated for the reuqest.
func (c *Client) Post(info *TRPInfo, body io.Reader) (_ *http.Response, err error) {
	var req *http.Request
	if req, err = NewRequest(info, body); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewRequest generates a new TRP POST request from the specified info, using the
// Address to determine the endpoint and setting headers accordingly. If the APIVersion
// is omitted, the default APIVersion is used. If there is no request identifier, one
// is generated and added to the info struct.
func NewRequest(info *TRPInfo, body io.Reader) (req *http.Request, err error) {
	if info.APIVersion == "" {
		info.APIVersion = APIVersion
	}

	if info.RequestIdentifier == "" {
		info.RequestIdentifier = uuid.NewString()
	}

	var endpoint string
	if endpoint, err = info.GetURL(); err != nil {
		return nil, err
	}

	if req, err = http.NewRequest(http.MethodPost, endpoint, body); err != nil {
		return nil, err
	}

	req.Header.Set(ContentTypeHeader, ContentTypeValue)
	req.Header.Set(APIVersionHeader, info.APIVersion)
	req.Header.Set(RequestIdentifierHeader, info.RequestIdentifier)
	if len(info.APIExtensions) > 0 {
		req.Header.Set(APIExtensionsHeader, strings.Join(info.APIExtensions, ","))
	}
	return req, nil
}
