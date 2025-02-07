package openvasp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// Discoverability Extension Endpoints
	DiscoverabilityVersion    = "/version"
	DiscoverabilityUptime     = "/uptime"
	DiscoverabilityExtensions = "/extensions"
	DiscoverabilityIdentity   = "/identity"
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

// Inquiry posts a Travel Rule Inquiry to the counterparty as specified by the info in
// the inquiry request and returns the response from the client.
func (c *Client) Inquiry(inquiry *Inquiry) (_ *TravelRuleResponse, err error) {
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(inquiry); err != nil {
		return nil, err
	}
	return c.Post(inquiry.TRP, body)
}

// Confirm posts a Travel Rule Confirmation to the counterparty as specified by the info
// in the confirmation request and returns the response from the client.
func (c *Client) Confirmation(confirm *Confirmation) (_ *TravelRuleResponse, err error) {
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(confirm); err != nil {
		return nil, err
	}
	return c.Post(confirm.TRP, body)
}

// Gets the identity of the TRP server specified by the travel address.
func (c *Client) Identity(info *TRPInfo) (out *IdentityInfo, err error) {
	// Use the identity endpoint
	var req *TRPInfo
	if req, err = discoverabilityEndpoint(info, DiscoverabilityIdentity); err != nil {
		return nil, err
	}

	var rep *TravelRuleResponse
	if rep, err = c.Get(req); err != nil {
		return nil, err
	}

	out = &IdentityInfo{}
	if err = rep.Decode(out); err != nil {
		return nil, err
	}

	return out, nil
}

// Gets the version and vendor of the TRP server specified by the travel address.
func (c *Client) Version(info *TRPInfo) (out *VersionInfo, err error) {
	// Use the version endpoint
	var req *TRPInfo
	if req, err = discoverabilityEndpoint(info, DiscoverabilityVersion); err != nil {
		return nil, err
	}

	var rep *TravelRuleResponse
	if rep, err = c.Get(req); err != nil {
		return nil, err
	}

	out = &VersionInfo{}
	if err = rep.Decode(out); err != nil {
		return nil, err
	}

	return out, nil
}

// Gets the supported extensions of the TRP server specified by the travel address.
func (c *Client) Extensions(info *TRPInfo) (out *ExtensionsInfo, err error) {
	// Use the extensions endpoint
	var req *TRPInfo
	if req, err = discoverabilityEndpoint(info, DiscoverabilityExtensions); err != nil {
		return nil, err
	}

	var rep *TravelRuleResponse
	if rep, err = c.Get(req); err != nil {
		return nil, err
	}

	out = &ExtensionsInfo{}
	if err = rep.Decode(out); err != nil {
		return nil, err
	}

	return out, nil
}

// Gets the uptime of the TRP server specified by the travel address.
func (c *Client) Uptime(info *TRPInfo) (out time.Duration, err error) {
	// Use the uptime endpoint
	var req *TRPInfo
	if req, err = discoverabilityEndpoint(info, DiscoverabilityUptime); err != nil {
		return 0, err
	}

	var rep *TravelRuleResponse
	if rep, err = c.Get(req); err != nil {
		return 0, err
	}

	// Check the content type to ensure JSON data was returned
	contentType := rep.Header.Get(ContentTypeHeader)
	if contentType != "" {
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return 0, err
		}

		if mediaType != ContentPlainText {
			return 0, fmt.Errorf("could not decode response content type %s", contentType)
		}
	}
	defer rep.Body.Close()

	var data []byte
	if data, err = io.ReadAll(rep.Body); err != nil {
		return 0, fmt.Errorf("could not read response body: %w", err)
	}

	// Read the integer from the body
	var uptime int64
	if uptime, err = strconv.ParseInt(string(data), 10, 64); err != nil {
		return 0, fmt.Errorf("could not parse uptime response: %w", err)
	}

	return time.Duration(uptime) * time.Second, nil
}

// Create the discoverability endpoint info without overwriting the original.
func discoverabilityEndpoint(info *TRPInfo, endpoint string) (_ *TRPInfo, err error) {
	var uri string
	if uri, err = info.GetURL(); err != nil {
		return nil, err
	}

	u, _ := url.Parse(uri)
	u.Path = endpoint
	u.Fragment = ""
	u.RawQuery = ""

	return &TRPInfo{
		Address:           u.String(),
		APIVersion:        info.APIVersion,
		RequestIdentifier: info.RequestIdentifier,
		APIExtensions:     info.APIExtensions,
	}, nil
}

// Get a request to the specified Address, which can be a Travel Address, LNURL, or
// plain-text URL. The content type is automatically set to application/json and the
// body should read JSON data. Required OpenVASP headers are set by the TRPInfo struct.
func (c *Client) Get(info *TRPInfo) (_ *TravelRuleResponse, err error) {
	var req *http.Request
	if req, err = NewRequest(http.MethodGet, info, nil); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = c.Client.Do(req); err != nil {
		return nil, err
	}

	trr := &TravelRuleResponse{Response: *rep}
	if serr := trr.StatusError(); serr != nil {
		return nil, serr
	}
	return trr, nil
}

// Post a request to the specified Address, which can be a Travel Address, LNURL, or
// plain-text URL. The content type is automatically set to application/json and the
// body should read JSON data. Required OpenVASP headers are set by the TRPInfo struct.
// If the APIVersion is omitted, the default APIVersion is used, similarly if no request
// identifier is present, one is generated for the reuqest.
func (c *Client) Post(info *TRPInfo, body io.Reader) (_ *TravelRuleResponse, err error) {
	var req *http.Request
	if req, err = NewRequest(http.MethodPost, info, body); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = c.Client.Do(req); err != nil {
		return nil, err
	}

	trr := &TravelRuleResponse{Response: *rep}
	if serr := trr.StatusError(); serr != nil {
		return nil, serr
	}
	return trr, nil
}

// NewRequest generates a new TRP POST request from the specified info, using the
// Address to determine the endpoint and setting headers accordingly. If the APIVersion
// is omitted, the default APIVersion is used. If there is no request identifier, one
// is generated and added to the info struct.
func NewRequest(method string, info *TRPInfo, body io.Reader) (req *http.Request, err error) {
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

	if req, err = http.NewRequest(method, endpoint, body); err != nil {
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

// TravelRuleResponse embeds an http Response and makes provisions for parsing the
// response data back from the counterparty.
type TravelRuleResponse struct {
	http.Response
	info        *TRPInfo
	statusError *StatusError
}

// Info returns the travel rule info from the headers of the response
func (t *TravelRuleResponse) Info() *TRPInfo {
	if t.info == nil {
		t.info = &TRPInfo{
			Address:           t.Request.URL.String(),
			APIVersion:        t.Header.Get(APIVersionHeader),
			RequestIdentifier: t.Header.Get(RequestIdentifierHeader),
		}

		if extensions := t.Header.Get(APIExtensionsHeader); extensions != "" {
			t.info.APIExtensions = strings.Split(extensions, ",")
		}
	}
	return t.info
}

// If a 400 or 500 status code was received, the status error is returned.
func (t *TravelRuleResponse) StatusError() *StatusError {
	if t.statusError == nil {
		if t.StatusCode < 200 || t.StatusCode >= 300 {
			t.statusError = &StatusError{
				Code: t.StatusCode,
			}

			contentType := t.Header.Get(ContentTypeHeader)
			if contentType != "" {
				if mediaType, _, err := mime.ParseMediaType(contentType); err == nil {
					if mediaType == "text/plain" {
						msg, _ := io.ReadAll(t.Body)
						t.statusError.Message = strings.TrimSpace(string(msg))
						t.Body.Close()
					}
				}
			}

		}
	}
	return t.statusError
}

// InquiryResolution attempts to parse a travel rule response from the response body,
// closing the response body when it is complete. If the body has already been closed
// then an EOF error will be returned.
func (t *TravelRuleResponse) InquiryResolution() (*InquiryResolution, error) {
	// Check the status and the header before reading the body
	// There is a possibility that nil, nil is returned if the server returned a 204.
	if t.StatusCode != http.StatusOK {
		return nil, t.StatusError()
	}

	resolution := &InquiryResolution{}
	if err := t.Decode(resolution); err != nil {
		return nil, err
	}

	return resolution, nil
}

// Decodes JSON from the response body into the specified interface.
func (t *TravelRuleResponse) Decode(v interface{}) error {
	// Check the content type to ensure JSON data was returned
	contentType := t.Header.Get(ContentTypeHeader)
	if contentType != "" {
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return err
		}

		if mediaType != ContentMediaType {
			return fmt.Errorf("could not decode response content type %s", contentType)
		}
	}

	defer t.Body.Close()

	// Create the JSON decoder
	decoder := json.NewDecoder(t.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
