package openvasp

type InquiryHandler interface {
	OnInquiry(*TRP) (*InquiryResolution, error)
}

type ConfirmationHandler interface {
	OnConfirmation(*Confirmation) error
}
