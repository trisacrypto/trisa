package dti

type DTI struct {
	Identifier string   // the identifier of the digital token
	LongName   string   // the description or long form name of the digital token
	ShortNames []string // any short names or coin symbols for the digital token
	Related    []string // auxiliary distributed ledger or underlying asset external identifier
}
