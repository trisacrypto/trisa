package ivms101

import "encoding/json"

//
// IdentityPayload JSON
//

type serialIdentityPayload struct {
	Originator      *Originator      `json:"originator,omitempty"`
	Beneficiary     *Beneficiary     `json:"beneficiary,omitempty"`
	OriginatingVASP *OriginatingVasp `json:"originatingVASP,omitempty"`
	BeneficiaryVASP *BeneficiaryVasp `json:"beneficiaryVASP,omitempty"`
	TransferPath    *TransferPath    `json:"transferPath,omitempty"`
	PayloadMetadata *PayloadMetadata `json:"payloadMetadata,omitempty"`
}

func (i *IdentityPayload) MarshalJSON() ([]byte, error) {
	middle := &serialIdentityPayload{
		Originator:      i.Originator,
		Beneficiary:     i.Beneficiary,
		OriginatingVASP: i.OriginatingVasp,
		BeneficiaryVASP: i.BeneficiaryVasp,
		TransferPath:    i.TransferPath,
		PayloadMetadata: i.PayloadMetadata,
	}
	return json.Marshal(middle)
}

func (i *IdentityPayload) UnmarshalJSON(data []byte) error {
	middle := serialIdentityPayload{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	p := IdentityPayload{
		Originator:      middle.Originator,
		Beneficiary:     middle.Beneficiary,
		OriginatingVasp: middle.OriginatingVASP,
		BeneficiaryVasp: middle.BeneficiaryVASP,
		TransferPath:    middle.TransferPath,
		PayloadMetadata: middle.PayloadMetadata,
	}

	// TODO warning: assignment copies lock value to *n
	*i = p
	return nil
}

//
// Identity Natural Persons JSON
//

type serialOriginator struct {
	Originator     []*Person `json:"originatorPersons,omitempty"`
	AccountNumbers []string  `json:"accountNumber,omitempty"`
}

func (o *Originator) MarshalJSON() ([]byte, error) {
	middle := serialOriginator{
		Originator:     o.OriginatorPersons,
		AccountNumbers: o.AccountNumbers,
	}
	return json.Marshal(middle)
}

func (o *Originator) UnmarshalJSON(data []byte) error {
	middle := serialOriginator{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := Originator{
		OriginatorPersons: middle.Originator,
		AccountNumbers:    middle.AccountNumbers,
	}

	// TODO warning: assignment copies lock value to *n
	*o = i
	return nil
}

type serialBeneficiary struct {
	Beneficiary    []*Person `json:"beneficiaryPersons,omitempty"`
	AccountNumbers []string  `json:"accountNumber,omitempty"`
}

func (b *Beneficiary) MarshalJSON() ([]byte, error) {
	middle := serialBeneficiary{
		Beneficiary:    b.BeneficiaryPersons,
		AccountNumbers: b.AccountNumbers,
	}
	return json.Marshal(middle)
}

func (b *Beneficiary) UnmarshalJSON(data []byte) error {
	middle := serialBeneficiary{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := Beneficiary{
		BeneficiaryPersons: middle.Beneficiary,
		AccountNumbers:     middle.AccountNumbers,
	}

	// TODO warning: assignment copies lock value to *n
	*b = i
	return nil
}

//
// Identity Legal Persons JSON
//

type serialOriginatorVASP struct {
	Originator *Person `json:"originatingVASP,omitempty"`
}

func (o *OriginatingVasp) MarshalJSON() ([]byte, error) {
	middle := serialOriginatorVASP{
		Originator: o.OriginatingVasp,
	}
	return json.Marshal(middle)
}

func (o *OriginatingVasp) UnmarshalJSON(data []byte) error {
	middle := serialOriginatorVASP{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := OriginatingVasp{
		OriginatingVasp: middle.Originator,
	}

	// TODO warning: assignment copies lock value to *n
	*o = i
	return nil
}

type serialBeneficiaryVASP struct {
	Beneficiary *Person `json:"beneficiaryVASP,omitempty"`
}

func (b *BeneficiaryVasp) MarshalJSON() ([]byte, error) {
	middle := serialBeneficiaryVASP{
		Beneficiary: b.BeneficiaryVasp,
	}
	return json.Marshal(middle)
}

func (b *BeneficiaryVasp) UnmarshalJSON(data []byte) error {
	middle := serialBeneficiaryVASP{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := BeneficiaryVasp{
		BeneficiaryVasp: middle.Beneficiary,
	}

	// TODO warning: assignment copies lock value to *n
	*b = i
	return nil
}

//
// Transfer Path JSON
//

type serialIntermediaryVASP struct {
	Intermediary *Person `json:"intermediaryVASP,omitempty"`
	Sequence     uint64  `json:"sequence,omitempty"`
}

func (v *IntermediaryVasp) MarshalJSON() ([]byte, error) {
	middle := serialIntermediaryVASP{
		Intermediary: v.IntermediaryVasp,
		Sequence:     v.Sequence,
	}
	return json.Marshal(middle)
}

func (v *IntermediaryVasp) UnmarshalJSON(data []byte) error {
	middle := serialIntermediaryVASP{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := IntermediaryVasp{
		IntermediaryVasp: middle.Intermediary,
		Sequence:         middle.Sequence,
	}

	// TODO warning: assignment copies lock value to *n
	*v = i
	return nil
}

type serialTransferPath struct {
	TransferPath []*IntermediaryVasp `json:"transferPath,omitempty"`
}

func (p *TransferPath) MarshalJSON() ([]byte, error) {
	middle := serialTransferPath{
		TransferPath: p.TransferPath,
	}
	return json.Marshal(middle)
}

func (p *TransferPath) UnmarshalJSON(data []byte) error {
	middle := serialTransferPath{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := TransferPath{
		TransferPath: middle.TransferPath,
	}

	// TODO warning: assignment copies lock value to *n
	*p = i
	return nil
}

//
// Payload Metadata JSON
//

type serialPayloadMetadata struct {
	TransliterationMethod []TransliterationMethodCode `json:"transliterationMethod,omitempty"`
}

func (p *PayloadMetadata) MarshalJSON() ([]byte, error) {
	middle := serialPayloadMetadata{
		TransliterationMethod: p.TransliterationMethod,
	}
	return json.Marshal(middle)
}

func (p *PayloadMetadata) UnmarshalJSON(data []byte) error {
	middle := serialPayloadMetadata{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := PayloadMetadata{
		TransliterationMethod: middle.TransliterationMethod,
	}

	// TODO warning: assignment copies lock value to *n
	*p = i
	return nil
}
