package openvasp

import (
	"mime"
	"net/http"
)

// APIChecks is middleware that asserts that the headers in the TRP request are correct
// and valid, ensuring that the core protocol is implemented correctly.
func APIChecks(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Enforce Application Version
		apiVersion := r.Header.Get(APIVersionHeader)
		if apiVersion != APIVersion {
			http.Error(w, "must specify api version header "+APIVersion, http.StatusBadRequest)
			return
		}

		// Set the APIVersion header in the outgoing response
		w.Header().Add(APIVersionHeader, APIVersion)

		// Must specify a request identifier
		var requestIdentifier string
		if requestIdentifier = r.Header.Get(RequestIdentifierHeader); requestIdentifier == "" {
			http.Error(w, "must specify request identifier header", http.StatusBadRequest)
			return
		}

		// Echo back the request identifier in the outgoing response
		w.Header().Add(RequestIdentifierHeader, requestIdentifier)

		// Enforce JSON content type; if no content-type is specified assume JSON
		contentType := r.Header.Get(ContentTypeHeader)
		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "malformed content-type header", http.StatusBadRequest)
				return
			}

			if mt != ContentMediaType {
				http.Error(w, "content-type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
