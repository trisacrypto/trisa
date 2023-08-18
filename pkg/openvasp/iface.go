package openvasp

type InquiryHandler interface {
	OnInquiry(*Inquiry) (*InquiryResolution, error)
}

type ConfirmationHandler interface {
	OnConfirmation(*Confirmation) error
}
