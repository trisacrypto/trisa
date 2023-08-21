package openvasp_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/trisa/pkg/openvasp"
	"github.com/trisacrypto/trisa/pkg/slip0044"
)

const (
	requestIdentifier = "f716c465-1b37-4a7b-aa32-f49346500721"
	originatorURL     = "https://originator.com/inquiryResolution?q=458839457"
	beneficiaryURL    = "https://beneficiary.com/transferConfirmation?q=3454366424"
)

func TestTransferInquiry(t *testing.T) {
	// Load the payload fixture
	payload, err := loadPayload("testdata/payload.json")
	require.NoError(t, err, "could not open fixture at testdata/payload.json")
	require.NotEmpty(t, payload, "payload does not contain any data")

	// Create requests to execute against the inquiry handler
	makeRequest := func(t *testing.T) (*httptest.ResponseRecorder, *http.Request) {
		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(payload)
		require.NoError(t, err, "could not encode json payload")

		req := httptest.NewRequest(http.MethodPost, originatorURL, &body)
		req.Header.Set(APIVersionHeader, APIVersion)
		req.Header.Set(RequestIdentifierHeader, requestIdentifier)
		req.Header.Set(ContentTypeHeader, ContentTypeValue)

		return httptest.NewRecorder(), req
	}

	// Create the mock handler for testing
	mock := &MockHandler{}
	handler := TransferInquiry(mock)

	t.Run("Standard", func(t *testing.T) {
		defer mock.Reset()
		mock.CallInquiry = func(i *Inquiry) (*InquiryResolution, error) {
			// Make some assertions about the inquiry
			if i.TRP == nil || i.TRP.RequestIdentifier != requestIdentifier || i.TRP.APIVersion != APIVersion {
				return nil, errors.New("invalid TRP info")
			}

			if i.Amount == 0 || i.Asset == nil || i.Asset.SLIP044 != slip0044.CoinType_BITCOIN || i.Callback != beneficiaryURL || i.IVMS101 == nil {
				return nil, errors.New("invalid payload")
			}

			// Return nil and nil to have the API version returned
			return nil, nil
		}

		w, r := makeRequest(t)
		handler.ServeHTTP(w, r)
		require.Equal(t, 1, mock.Calls(CallInquiry))
		require.Equal(t, 0, mock.Calls(CallConfirmation))

		rep := w.Result()
		require.Equal(t, http.StatusOK, rep.StatusCode)
		require.Equal(t, requestIdentifier, rep.Header.Get(RequestIdentifierHeader))
		require.Equal(t, APIVersion, rep.Header.Get(APIVersionHeader))
		require.Equal(t, ContentTypeValue, rep.Header.Get(ContentTypeHeader))

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{"version": "3.1.0"}`, string(data))
	})

	t.Run("Approval", func(t *testing.T) {
		defer mock.Reset()
		mock.CallInquiry = func(i *Inquiry) (*InquiryResolution, error) {
			return &InquiryResolution{
				Approved: &Approval{
					Address:  "some payment address",
					Callback: beneficiaryURL,
				},
			}, nil
		}

		w, r := makeRequest(t)
		handler.ServeHTTP(w, r)
		require.Equal(t, 1, mock.Calls(CallInquiry))
		require.Equal(t, 0, mock.Calls(CallConfirmation))

		rep := w.Result()
		require.Equal(t, http.StatusOK, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		expected := fmt.Sprintf(`{"approved": {"address": "some payment address", "callback": "%s"}}`, beneficiaryURL)
		require.JSONEq(t, expected, string(data))
	})

	t.Run("Rejection", func(t *testing.T) {
		defer mock.Reset()
		mock.CallInquiry = func(i *Inquiry) (*InquiryResolution, error) {
			return &InquiryResolution{
				Rejected: "human readable comment",
			}, nil
		}

		w, r := makeRequest(t)
		handler.ServeHTTP(w, r)
		require.Equal(t, 1, mock.Calls(CallInquiry))
		require.Equal(t, 0, mock.Calls(CallConfirmation))

		rep := w.Result()
		require.Equal(t, http.StatusOK, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{"rejected": "human readable comment"}`, string(data))
	})

	t.Run("Error", func(t *testing.T) {
		defer mock.Reset()
		mock.UseError(CallInquiry, errors.New("whoopsie"))

		w, r := makeRequest(t)
		handler.ServeHTTP(w, r)
		require.Equal(t, 1, mock.Calls(CallInquiry))
		require.Equal(t, 0, mock.Calls(CallConfirmation))

		rep := w.Result()
		require.Equal(t, http.StatusInternalServerError, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "whoopsie\n", string(data))
	})

	t.Run("StatusError", func(t *testing.T) {
		defer mock.Reset()
		mock.UseError(CallInquiry, &StatusError{Code: http.StatusConflict})

		w, r := makeRequest(t)
		handler.ServeHTTP(w, r)
		require.Equal(t, 1, mock.Calls(CallInquiry))
		require.Equal(t, 0, mock.Calls(CallConfirmation))

		rep := w.Result()
		require.Equal(t, http.StatusConflict, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "Conflict\n", string(data))
	})
}

func TestTransferConfirmation(t *testing.T) {
	// Create requests to execute against the inquiry handler
	makeRequest := func(t *testing.T, payload *Confirmation) (*httptest.ResponseRecorder, *http.Request) {
		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(payload)
		require.NoError(t, err, "could not encode json payload")

		req := httptest.NewRequest(http.MethodPost, beneficiaryURL, &body)
		req.Header.Set(APIVersionHeader, APIVersion)
		req.Header.Set(RequestIdentifierHeader, requestIdentifier)
		req.Header.Set(ContentTypeHeader, ContentTypeValue)

		return httptest.NewRecorder(), req
	}

	// Create the mock handler for testing
	mock := &MockHandler{}
	handler := TransferConfirmation(mock)

	t.Run("TXID", func(t *testing.T) {
		defer mock.Reset()
		mock.CallConfirmation = func(c *Confirmation) error {
			if c.TRP == nil || c.TRP.RequestIdentifier != requestIdentifier || c.TRP.APIVersion != APIVersion {
				return errors.New("invalid TRP info")
			}

			if err := c.Validate(); err != nil {
				return err
			}

			if c.TXID != "foo" {
				return errors.New("unexpected txid")
			}
			return nil
		}

		w, r := makeRequest(t, &Confirmation{TXID: "foo"})
		handler.ServeHTTP(w, r)
		require.Equal(t, 1, mock.Calls(CallConfirmation))
		require.Equal(t, 0, mock.Calls(CallInquiry))

		rep := w.Result()
		require.Equal(t, http.StatusNoContent, rep.StatusCode)
		require.Equal(t, requestIdentifier, rep.Header.Get(RequestIdentifierHeader))
		require.Equal(t, APIVersion, rep.Header.Get(APIVersionHeader))
		require.Empty(t, rep.Header.Get(ContentTypeHeader))
	})

	t.Run("Canceled", func(t *testing.T) {
		defer mock.Reset()
		mock.CallConfirmation = func(c *Confirmation) error {
			if c.TRP == nil || c.TRP.RequestIdentifier != requestIdentifier || c.TRP.APIVersion != APIVersion {
				return errors.New("invalid TRP info")
			}

			if err := c.Validate(); err != nil {
				return err
			}

			if c.Canceled != "foo" {
				return errors.New("unexpected canceled")
			}
			return nil
		}

		w, r := makeRequest(t, &Confirmation{Canceled: "foo"})
		handler.ServeHTTP(w, r)
		require.Equal(t, 1, mock.Calls(CallConfirmation))
		require.Equal(t, 0, mock.Calls(CallInquiry))

		rep := w.Result()
		require.Equal(t, http.StatusNoContent, rep.StatusCode)
		require.Equal(t, requestIdentifier, rep.Header.Get(RequestIdentifierHeader))
		require.Equal(t, APIVersion, rep.Header.Get(APIVersionHeader))
		require.Empty(t, rep.Header.Get(ContentTypeHeader))
	})

	t.Run("Invalid", func(t *testing.T) {
		defer mock.Reset()
		mock.UseError(CallConfirmation, errors.New("whoopsie"))

		w, r := makeRequest(t, &Confirmation{})
		handler.ServeHTTP(w, r)
		require.Equal(t, 0, mock.Calls(CallInquiry))
		require.Equal(t, 0, mock.Calls(CallConfirmation))

		rep := w.Result()
		require.Equal(t, http.StatusBadRequest, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "must specify either txid or canceled in confirmation\n", string(data))
	})

	t.Run("Error", func(t *testing.T) {
		defer mock.Reset()
		mock.UseError(CallConfirmation, errors.New("whoopsie"))

		w, r := makeRequest(t, &Confirmation{TXID: "foo"})
		handler.ServeHTTP(w, r)
		require.Equal(t, 0, mock.Calls(CallInquiry))
		require.Equal(t, 1, mock.Calls(CallConfirmation))

		rep := w.Result()
		require.Equal(t, http.StatusInternalServerError, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "whoopsie\n", string(data))
	})

	t.Run("StatusError", func(t *testing.T) {
		defer mock.Reset()
		mock.UseError(CallConfirmation, &StatusError{Code: http.StatusExpectationFailed})

		w, r := makeRequest(t, &Confirmation{TXID: "foo"})
		handler.ServeHTTP(w, r)
		require.Equal(t, 0, mock.Calls(CallInquiry))
		require.Equal(t, 1, mock.Calls(CallConfirmation))

		rep := w.Result()
		require.Equal(t, http.StatusExpectationFailed, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "Expectation Failed\n", string(data))
	})
}

func TestAPIChecks(t *testing.T) {
	handler := APIChecks(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	t.Run("MethodNotAllowed", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, originatorURL, nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, r)
		rep := w.Result()

		require.Equal(t, http.StatusMethodNotAllowed, rep.StatusCode, "expected method not allowed on GET")
	})

	t.Run("APIVersion", func(t *testing.T) {
		testCases := []string{"", "9.9.999"}
		for i, tc := range testCases {
			r := httptest.NewRequest(http.MethodPost, originatorURL, nil)
			w := httptest.NewRecorder()

			if tc != "" {
				r.Header.Set(APIVersionHeader, tc)
			}

			handler.ServeHTTP(w, r)
			rep := w.Result()

			require.Equal(t, http.StatusBadRequest, rep.StatusCode, "test case %d failed", i)
			data, _ := io.ReadAll(rep.Body)
			require.Equal(t, "must specify api version header 3.1.0\n", string(data), "test case %d failed", i)
		}
	})

	t.Run("Request Identifier", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, originatorURL, nil)
		w := httptest.NewRecorder()

		r.Header.Set(APIVersionHeader, APIVersion)

		handler.ServeHTTP(w, r)
		rep := w.Result()

		require.Equal(t, http.StatusBadRequest, rep.StatusCode)
		data, _ := io.ReadAll(rep.Body)
		require.Equal(t, "must specify request identifier\n", string(data))
	})

	t.Run("Malformed Content Type", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, originatorURL, nil)
		w := httptest.NewRecorder()

		r.Header.Set(APIVersionHeader, APIVersion)
		r.Header.Set(RequestIdentifierHeader, "foo")
		r.Header.Set(ContentTypeHeader, "baz/;;;--")

		handler.ServeHTTP(w, r)
		rep := w.Result()

		require.Equal(t, http.StatusUnsupportedMediaType, rep.StatusCode)
		data, _ := io.ReadAll(rep.Body)
		require.Equal(t, "malformed content-type header\n", string(data))
	})

	t.Run("Incorrect Content Type", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, originatorURL, nil)
		w := httptest.NewRecorder()

		r.Header.Set(APIVersionHeader, APIVersion)
		r.Header.Set(RequestIdentifierHeader, "foo")
		r.Header.Set(ContentTypeHeader, "application/xml")

		handler.ServeHTTP(w, r)
		rep := w.Result()

		require.Equal(t, http.StatusUnsupportedMediaType, rep.StatusCode)
		data, _ := io.ReadAll(rep.Body)
		require.Equal(t, "content-type header must be application/json\n", string(data))
	})

	t.Run("Happy", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, originatorURL, nil)
		w := httptest.NewRecorder()

		r.Header.Set(APIVersionHeader, APIVersion)
		r.Header.Set(RequestIdentifierHeader, "foo")
		r.Header.Set(ContentTypeHeader, ContentTypeValue)

		handler.ServeHTTP(w, r)
		rep := w.Result()

		require.Equal(t, http.StatusNoContent, rep.StatusCode)
		require.Equal(t, "3.1.0", rep.Header.Get(APIVersionHeader))
		require.Equal(t, "foo", rep.Header.Get(RequestIdentifierHeader))
	})
}

func TestParseTRPInfo(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, originatorURL, nil)
		info := ParseTRPInfo(r)
		require.NotNil(t, info, "expected info to not be nil")
		require.Equal(t, originatorURL, info.LNURL)
		require.Zero(t, info.APIVersion)
		require.Zero(t, info.RequestIdentifier)
		require.Zero(t, info.APIExtensions)
	})

	t.Run("Populated", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, originatorURL, nil)
		r.Header.Set(APIVersionHeader, APIVersion)
		r.Header.Set(RequestIdentifierHeader, requestIdentifier)
		r.Header.Set(APIExtensionsHeader, "request-signing,beneficiary-details-response")

		info := ParseTRPInfo(r)
		require.NotNil(t, info, "expected info to not be nil")
		require.Equal(t, originatorURL, info.LNURL)
		require.Equal(t, APIVersion, info.APIVersion)
		require.Equal(t, requestIdentifier, info.RequestIdentifier)
		require.Equal(t, []string{"request-signing", "beneficiary-details-response"}, info.APIExtensions)
	})
}

func TestDecodeJSON(t *testing.T) {
	mock := &MockHandler{}
	inquiry := TransferInquiry(mock)
	confirm := TransferConfirmation(mock)

	makeRequest := func(t *testing.T, data string) (*httptest.ResponseRecorder, *http.Request) {
		var body bytes.Buffer
		_, err := body.Write([]byte(data))
		require.NoError(t, err, "could not encode json payload")

		req := httptest.NewRequest(http.MethodPost, originatorURL, &body)
		req.Header.Set(APIVersionHeader, APIVersion)
		req.Header.Set(RequestIdentifierHeader, requestIdentifier)
		req.Header.Set(ContentTypeHeader, ContentTypeValue)

		return httptest.NewRecorder(), req
	}

	t.Run("SyntaxError", func(t *testing.T) {
		w, r := makeRequest(t, `{"data: "foo"}`)
		inquiry.ServeHTTP(w, r)

		rep := w.Result()
		require.Equal(t, http.StatusBadRequest, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "request body contains badly-formed JSON (at position 10)\n", string(data))
	})

	t.Run("UnexpectedEOF", func(t *testing.T) {
		w, r := makeRequest(t, `{"data": "foo"`)
		confirm.ServeHTTP(w, r)

		rep := w.Result()
		require.Equal(t, http.StatusBadRequest, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "request body contains badly-formed JSON\n", string(data))
	})

	t.Run("EOF", func(t *testing.T) {
		w, r := makeRequest(t, "")
		inquiry.ServeHTTP(w, r)

		rep := w.Result()
		require.Equal(t, http.StatusBadRequest, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "Bad Request\n", string(data))
	})

	t.Run("TypeError", func(t *testing.T) {
		w, r := makeRequest(t, `{"txid": [1, 2, 3]}`)
		confirm.ServeHTTP(w, r)

		rep := w.Result()
		require.Equal(t, http.StatusBadRequest, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "request body contains an invalid value for the \"txid\" field (at 10)\n", string(data))
	})

	t.Run("Unknown Field", func(t *testing.T) {
		w, r := makeRequest(t, `{"data": "foo"}`)
		inquiry.ServeHTTP(w, r)

		rep := w.Result()
		require.Equal(t, http.StatusBadRequest, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "request body contains unknown field \"data\"\n", string(data))
	})

	t.Run("TooLarge", func(t *testing.T) {
		bigval := strings.Repeat(`money money money ... must be funny`, 700000)
		w, r := makeRequest(t, fmt.Sprintf(`{"txid": "%s"}`, bigval))
		confirm.ServeHTTP(w, r)

		rep := w.Result()
		require.Equal(t, http.StatusRequestEntityTooLarge, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "Request Entity Too Large\n", string(data))
	})

	t.Run("DoubleData", func(t *testing.T) {
		w, r := makeRequest(t, `{"txid": "foo"}{"canceled": "reason"}`)
		confirm.ServeHTTP(w, r)

		rep := w.Result()
		require.Equal(t, http.StatusBadRequest, rep.StatusCode)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "request body must only contain a single JSON object\n", string(data))
	})
}

const (
	CallInquiry      = "OnInquiry"
	CallConfirmation = "OnConfirmation"
)

// MockHandler implements the InquiryHandler and Confirmation Handler interface
type MockHandler struct {
	sync.RWMutex
	calls map[string]int

	CallInquiry      func(*Inquiry) (*InquiryResolution, error)
	CallConfirmation func(*Confirmation) error
}

func (m *MockHandler) UseError(call string, err error) {
	switch call {
	case CallInquiry:
		m.CallInquiry = func(*Inquiry) (*InquiryResolution, error) { return nil, err }
	case CallConfirmation:
		m.CallConfirmation = func(*Confirmation) error { return err }
	default:
		panic(fmt.Errorf("unknown call %q", call))
	}
}

func (m *MockHandler) OnInquiry(in *Inquiry) (*InquiryResolution, error) {
	m.incr(CallInquiry)
	if m.CallInquiry != nil {
		return m.CallInquiry(in)
	}
	return nil, errors.New("no mock on inquiry handler defined")
}

func (m *MockHandler) OnConfirmation(in *Confirmation) error {
	m.incr(CallConfirmation)
	if m.CallConfirmation != nil {
		return m.CallConfirmation(in)
	}
	return errors.New("no mock on inquiry handler defined")
}

func (m *MockHandler) Reset() {
	m.Lock()
	defer m.Unlock()
	for call := range m.calls {
		m.calls[call] = 0
	}

	m.CallInquiry = nil
	m.CallConfirmation = nil
}

func (m *MockHandler) Calls(call string) int {
	m.RLock()
	defer m.RUnlock()
	if m.calls == nil {
		return 0
	}
	return m.calls[call]
}

func (m *MockHandler) incr(call string) {
	m.Lock()
	defer m.Unlock()
	if m.calls == nil {
		m.calls = make(map[string]int)
	}
	m.calls[call]++
}

func loadPayload(path string) (inquiry *Inquiry, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&inquiry); err != nil {
		return nil, err
	}

	return inquiry, nil
}
