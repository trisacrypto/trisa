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
		require.Equal(t, rep.StatusCode, http.StatusOK)
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
		require.Equal(t, rep.StatusCode, http.StatusOK)

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
		require.Equal(t, rep.StatusCode, http.StatusOK)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{"rejected": "human readable comment"}`, string(data))
	})

	t.Run("Error", func(t *testing.T) {
		defer mock.Reset()
		mock.CallInquiry = func(i *Inquiry) (*InquiryResolution, error) {
			return nil, errors.New("whoopsie")
		}

		w, r := makeRequest(t)
		handler.ServeHTTP(w, r)
		require.Equal(t, 1, mock.Calls(CallInquiry))
		require.Equal(t, 0, mock.Calls(CallConfirmation))

		rep := w.Result()
		require.Equal(t, rep.StatusCode, http.StatusInternalServerError)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "whoopsie\n", string(data))
	})

	t.Run("StatusError", func(t *testing.T) {
		defer mock.Reset()
		mock.CallInquiry = func(i *Inquiry) (*InquiryResolution, error) {
			return nil, &StatusError{Code: http.StatusConflict}
		}

		w, r := makeRequest(t)
		handler.ServeHTTP(w, r)
		require.Equal(t, 1, mock.Calls(CallInquiry))
		require.Equal(t, 0, mock.Calls(CallConfirmation))

		rep := w.Result()
		require.Equal(t, rep.StatusCode, http.StatusConflict)

		data, err := io.ReadAll(rep.Body)
		require.NoError(t, err)
		require.Equal(t, "Conflict\n", string(data))
	})
}

func TestTransferConfirmation(t *testing.T) {

}

func TestAPIChecks(t *testing.T) {

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
