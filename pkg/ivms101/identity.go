package ivms101

import "encoding/json"

//===========================================================================
// IdentityPayload Methods
//===========================================================================

type serialIdentityPayload struct {
	Originator      *Originator      `json:"originator,omitempty"`
	Beneficiary     *Beneficiary     `json:"beneficiary,omitempty"`
	OriginatingVASP *OriginatingVasp `json:"originatingVASP,omitempty"`
	BeneficiaryVASP *BeneficiaryVasp `json:"beneficiaryVASP,omitempty"`
	TransferPath    *TransferPath    `json:"transferPath,omitempty"`
	PayloadMetadata *PayloadMetadata `json:"payloadMetadata,omitempty"`
}

var serialIdentityPayloadFields = map[string]string{
	"originator":       "originator",
	"originators":      "originator",
	"beneficiary":      "beneficiary",
	"beneficiaries":    "beneficiary",
	"originatingVASP":  "originatingVASP",
	"originating_vasp": "originatingVASP",
	"beneficiaryVASP":  "beneficiaryVASP",
	"beneficiary_vasp": "beneficiaryVASP",
	"transferPath":     "transferPath",
	"transfer_path":    "transferPath",
	"payloadMetadata":  "payloadMetadata",
	"payload_metadata": "payloadMetadata",
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

func (i *IdentityPayload) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialIdentityPayloadFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialIdentityPayload{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate the identity payload value
	i.Originator = middle.Originator
	i.Beneficiary = middle.Beneficiary
	i.OriginatingVasp = middle.OriginatingVASP
	i.BeneficiaryVasp = middle.BeneficiaryVASP
	i.TransferPath = middle.TransferPath
	i.PayloadMetadata = middle.PayloadMetadata

	return nil
}

//===========================================================================
// Originator Methods
//===========================================================================

type serialOriginator struct {
	Originator     []*Person `json:"originatorPersons,omitempty"`
	AccountNumbers []string  `json:"accountNumber,omitempty"`
}

var serialOriginatorFields = map[string]string{
	"originatorPersons":  "originatorPersons",
	"originatorPerson":   "originatorPersons",
	"originator_persons": "originatorPersons",
	"originator_person":  "originatorPersons",
	"originator":         "originatorPersons",
	"originators":        "originatorPersons",
	"accountNumber":      "accountNumber",
	"accountNumbers":     "accountNumber",
	"account_number":     "accountNumber",
	"account_numbers":    "accountNumber",
}

func (o *Originator) MarshalJSON() ([]byte, error) {
	middle := serialOriginator{
		Originator:     o.OriginatorPersons,
		AccountNumbers: o.AccountNumbers,
	}
	return json.Marshal(middle)
}

func (o *Originator) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialOriginatorFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialOriginator{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate originator values
	o.OriginatorPersons = middle.Originator
	o.AccountNumbers = middle.AccountNumbers

	return nil
}

//===========================================================================
// Beneficiary Methods
//===========================================================================

type serialBeneficiary struct {
	Beneficiary    []*Person `json:"beneficiaryPersons,omitempty"`
	AccountNumbers []string  `json:"accountNumber,omitempty"`
}

var serialBeneficiaryFields = map[string]string{
	"beneficiaryPersons":  "beneficiaryPersons",
	"beneficiaryPerson":   "beneficiaryPersons",
	"beneficiary_persons": "beneficiaryPersons",
	"beneficiary_person":  "beneficiaryPersons",
	"beneficiary":         "beneficiaryPersons",
	"beneficiaries":       "beneficiaryPersons",
	"accountNumber":       "accountNumber",
	"accountNumbers":      "accountNumber",
	"account_number":      "accountNumber",
	"account_numbers":     "accountNumber",
}

func (b *Beneficiary) MarshalJSON() ([]byte, error) {
	middle := serialBeneficiary{
		Beneficiary:    b.BeneficiaryPersons,
		AccountNumbers: b.AccountNumbers,
	}
	return json.Marshal(middle)
}

func (b *Beneficiary) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialBeneficiaryFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialBeneficiary{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate originator values
	b.BeneficiaryPersons = middle.Beneficiary
	b.AccountNumbers = middle.AccountNumbers

	return nil
}

//===========================================================================
// OriginatorVASP Methods
//===========================================================================

type serialOriginatorVASP struct {
	Originator *Person `json:"originatingVASP,omitempty"`
}

var serialOriginatorVASPFields = map[string]string{
	"originatingVASP":  "originatingVASP",
	"originating_vasp": "originatingVASP",
}

func (o *OriginatingVasp) MarshalJSON() ([]byte, error) {
	middle := serialOriginatorVASP{
		Originator: o.OriginatingVasp,
	}
	return json.Marshal(middle)
}

func (o *OriginatingVasp) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialOriginatorVASPFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialOriginatorVASP{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate originator vasp values.
	o.OriginatingVasp = middle.Originator

	return nil
}

//===========================================================================
// BeneficiaryVASP Methods
//===========================================================================

type serialBeneficiaryVASP struct {
	Beneficiary *Person `json:"beneficiaryVASP,omitempty"`
}

var serialBeneficiaryVASPFields = map[string]string{
	"beneficiaryVASP":  "beneficiaryVASP",
	"beneficiary_vasp": "beneficiaryVASP",
}

func (b *BeneficiaryVasp) MarshalJSON() ([]byte, error) {
	middle := serialBeneficiaryVASP{
		Beneficiary: b.BeneficiaryVasp,
	}
	return json.Marshal(middle)
}

func (b *BeneficiaryVasp) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialBeneficiaryVASPFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialBeneficiaryVASP{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate beneficiary vasp values
	b.BeneficiaryVasp = middle.Beneficiary

	return nil
}

//===========================================================================
// IntermediaryVASP Methods
//===========================================================================

type serialIntermediaryVASP struct {
	Intermediary *Person `json:"intermediaryVASP,omitempty"`
	Sequence     uint64  `json:"sequence,omitempty"`
}

var serialIntermediaryVASPFields = map[string]string{
	"intermediaryVASP":   "intermediaryVASP",
	"intermediaryVASPs":  "intermediaryVASP",
	"intermediary_vasp":  "intermediaryVASP",
	"intermediary_vasps": "intermediaryVASP",
	"sequence":           "sequence",
}

func (v *IntermediaryVasp) MarshalJSON() ([]byte, error) {
	middle := serialIntermediaryVASP{
		Intermediary: v.IntermediaryVasp,
		Sequence:     v.Sequence,
	}
	return json.Marshal(middle)
}

func (v *IntermediaryVasp) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialIntermediaryVASPFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialIntermediaryVASP{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate intermediary vasp values
	v.IntermediaryVasp = middle.Intermediary
	v.Sequence = middle.Sequence

	return nil
}

//===========================================================================
// TransferPath Methods
//===========================================================================

type serialTransferPath struct {
	TransferPath []*IntermediaryVasp `json:"transferPath,omitempty"`
}

var serialTransferPathFields = map[string]string{
	"transferPath":  "transferPath",
	"transfer_path": "transferPath",
}

func (p *TransferPath) MarshalJSON() ([]byte, error) {
	middle := serialTransferPath{
		TransferPath: p.TransferPath,
	}
	return json.Marshal(middle)
}

func (p *TransferPath) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialTransferPathFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialTransferPath{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate transfer path values
	p.TransferPath = middle.TransferPath

	return nil
}

//===========================================================================
// PayloadMetadata Methods
//===========================================================================

type serialPayloadMetadata struct {
	TransliterationMethod []TransliterationMethodCode `json:"transliterationMethod,omitempty"`
}

var serialPayloadMetadataFields = map[string]string{
	"transliterationMethod":   "transliterationMethod",
	"transliteration_method":  "transliterationMethod",
	"transliterationMethods":  "transliterationMethod",
	"transliteration_methods": "transliterationMethod",
	"methods":                 "transliterationMethod",
}

func (p *PayloadMetadata) MarshalJSON() ([]byte, error) {
	middle := serialPayloadMetadata{
		TransliterationMethod: p.TransliterationMethod,
	}
	return json.Marshal(middle)
}

func (p *PayloadMetadata) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialPayloadMetadataFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialPayloadMetadata{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate payload metadata values
	p.TransliterationMethod = middle.TransliterationMethod

	return nil
}
