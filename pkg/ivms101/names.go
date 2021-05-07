package ivms101

import "strings"

// Names returns a list of strings that collects all of the names from the NaturalPerson
// including name identifiers, local name identifiers, and phonetic name identifiers,
// in that order. The secondary and primary names are joined with a space such that the
// secondary name is first. No identifier codes or sorting are returned.
func (p *NaturalPerson) Names() []string {
	if p.Name == nil {
		return nil
	}

	// Create the name array
	nameCount := len(p.Name.NameIdentifiers) + len(p.Name.LocalNameIdentifiers) + len(p.Name.PhoneticNameIdentifiers)
	names := make([]string, 0, nameCount)

	for _, nameID := range p.Name.NameIdentifiers {
		name := strings.TrimSpace(nameID.SecondaryIdentifier + " " + nameID.PrimaryIdentifier)
		names = append(names, name)
	}

	for _, nameID := range p.Name.LocalNameIdentifiers {
		name := strings.TrimSpace(nameID.SecondaryIdentifier + " " + nameID.PrimaryIdentifier)
		names = append(names, name)
	}

	for _, nameID := range p.Name.PhoneticNameIdentifiers {
		name := strings.TrimSpace(nameID.SecondaryIdentifier + " " + nameID.PrimaryIdentifier)
		names = append(names, name)
	}

	return names
}

// Names returns a list of strings that collects all of the names from the LegalPerson
// including name identifiers, local name identifiers, and phonetic name identifiers,
// in that order. No identifier codes or sorting are returned.
func (p *LegalPerson) Names() []string {
	if p.Name == nil {
		return nil
	}

	nameCount := len(p.Name.NameIdentifiers) + len(p.Name.LocalNameIdentifiers) + len(p.Name.PhoneticNameIdentifiers)
	names := make([]string, 0, nameCount)

	for _, nameID := range p.Name.NameIdentifiers {
		names = append(names, nameID.LegalPersonName)
	}

	for _, nameID := range p.Name.LocalNameIdentifiers {
		names = append(names, nameID.LegalPersonName)
	}

	for _, nameID := range p.Name.PhoneticNameIdentifiers {
		names = append(names, nameID.LegalPersonName)
	}

	return names
}
