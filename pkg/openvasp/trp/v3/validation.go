package trp

func (i Identity) Validate() error {
	// All fields are optional.
	return nil
}

func (i Inquiry) Validate() (err error) {
	// TODO: return field based validation errors.
	if i.Asset != nil {
		if err = i.Asset.Validate(); err != nil {
			return err
		}
	}

	if i.Amount == 0 {
		return ErrNoAmount
	}

	if i.Callback == "" {
		return ErrEmptyCallback
	}

	if i.IVMS101 == nil {
		return ErrMissingIVMS101
	}

	// TODO: Validate IVMS101 payload
	return nil
}

func (a Asset) Validate() error {
	if a.DTI == "" && a.SLIP044 == 0 {
		return ErrEmptyAsset
	}

	// TODO: validate DTI
	return nil
}

func (r Resolution) Validate() error {
	switch {
	case r.Version != "":
		if r.Approved != nil || r.Rejected != "" {
			return ErrAmbiguousResolution
		}

		// TODO: check to ensure version is semver?
		return nil

	case r.Approved != nil:
		if r.Version != "" || r.Rejected != "" {
			return ErrAmbiguousResolution
		}

		return r.Approved.Validate()

	case r.Rejected != "":
		if r.Version != "" || r.Approved != nil {
			return ErrAmbiguousResolution
		}
		return nil

	default:
		return ErrEmptyResolution
	}
}

func (a Approval) Validate() error {
	if a.Address == "" {
		return ErrEmptyAddress
	}

	if a.Callback == "" {
		return ErrEmptyCallback
	}

	return nil
}

func (c Confirmation) Validate() error {
	if c.TXID == "" && c.Canceled == "" {
		return ErrEmptyConfirmation
	}

	if c.TXID != "" && c.Canceled != "" {
		return ErrAmbiguousConfirmation
	}

	return nil
}
