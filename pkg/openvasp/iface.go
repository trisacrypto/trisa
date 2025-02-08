package openvasp

import "github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"

type InquiryHandler interface {
	OnInquiry(*trp.Inquiry) (*trp.Resolution, error)
}

type ConfirmationHandler interface {
	OnConfirmation(*trp.Confirmation) error
}
