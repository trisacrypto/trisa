package ivms101

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
