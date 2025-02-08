package client

import (
	"errors"
	"net/http"
	"strings"

	"github.com/trisacrypto/trisa/pkg/openvasp"
	"github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

type Response struct {
	http.Response
	info *trp.Info
	err  error
}

func (r *Response) Info() *trp.Info {
	if r.info == nil {
		r.info = &trp.Info{
			Address:           r.Request.URL.String(),
			APIVersion:        r.Header.Get(openvasp.APIVersionHeader),
			RequestIdentifier: r.Header.Get(openvasp.RequestIdentifierHeader),
		}

		if extensions := r.Header.Get(openvasp.APIExtensionsHeader); extensions != "" {
			parts := strings.Split(extensions, ",")
			r.info.APIExtensions = make([]string, 0, len(parts))
			for _, part := range parts {
				r.info.APIExtensions = append(r.info.APIExtensions, strings.TrimSpace(part))
			}
		}
	}
	return r.info
}

func (r *Response) Err() error {
	return r.err
}

func (r *Response) StatusError() (*trp.StatusError, bool) {
	if r.err != nil {
		var status *trp.StatusError
		if errors.As(r.err, &status) {
			return status, true
		}
		return nil, false
	}
	return nil, true
}
