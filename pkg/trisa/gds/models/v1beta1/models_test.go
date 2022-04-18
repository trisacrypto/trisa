package models_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func TestParseEnums(t *testing.T) {
	bcat, err := pb.ParseBusinessCategory("unknown entity")
	require.NoError(t, err)
	require.Equal(t, pb.BusinessCategoryUnknown, bcat)

	bcat, err = pb.ParseBusinessCategory("PRIVATE_ORGANIZATION")
	require.NoError(t, err)
	require.Equal(t, pb.BusinessCategoryPrivate, bcat)

	bcat, err = pb.ParseBusinessCategory("Government Entity")
	require.NoError(t, err)
	require.Equal(t, pb.BusinessCategoryGovernment, bcat)

	bcat, err = pb.ParseBusinessCategory("Business_entity")
	require.NoError(t, err)
	require.Equal(t, pb.BusinessCategoryBusiness, bcat)

	bcat, err = pb.ParseBusinessCategory("non commercial entity")
	require.NoError(t, err)
	require.Equal(t, pb.BusinessCategoryNonCommercial, bcat)

	_, err = pb.ParseBusinessCategory("foo bar")
	require.Error(t, err)
}

func TestName(t *testing.T) {
	// VASP with no entity
	vasp := &pb.VASP{
		Id: "vasp_id",
	}
	_, err := vasp.Name()
	require.Error(t, err)
	vasp.Entity = &ivms101.LegalPerson{}
	_, err = vasp.Name()
	require.Error(t, err)

	// No name identifiers
	vasp.Entity = &ivms101.LegalPerson{
		Name: &ivms101.LegalPersonName{
			NameIdentifiers: []*ivms101.LegalPersonNameId{},
		},
	}
	_, err = vasp.Name()
	require.Error(t, err)

	// No supported name identifiers
	vasp.Entity.Name.NameIdentifiers = []*ivms101.LegalPersonNameId{
		{
			LegalPersonNameIdentifierType: -1,
		},
	}
	_, err = vasp.Name()
	require.Error(t, err)

	// Trading name takes precedence
	vasp.Entity.Name.NameIdentifiers = []*ivms101.LegalPersonNameId{
		{
			LegalPersonNameIdentifierType: -1,
			LegalPersonName:               "Invalid Name",
		},
		{
			LegalPersonNameIdentifierType: ivms101.LegalPersonShort,
			LegalPersonName:               "Short Name",
		},
		{
			LegalPersonNameIdentifierType: ivms101.LegalPersonLegal,
			LegalPersonName:               "Legal Name",
		},
		{
			LegalPersonNameIdentifierType: ivms101.LegalPersonTrading,
			LegalPersonName:               "Trading Name",
		},
	}
	name, err := vasp.Name()
	require.NoError(t, err)
	require.Equal(t, "Trading Name", name)

	// Short name is selected if no trading name
	vasp.Entity.Name.NameIdentifiers = []*ivms101.LegalPersonNameId{
		{
			LegalPersonNameIdentifierType: -1,
			LegalPersonName:               "Invalid Name",
		},
		{
			LegalPersonNameIdentifierType: ivms101.LegalPersonLegal,
			LegalPersonName:               "Legal Name",
		},
		{
			LegalPersonNameIdentifierType: ivms101.LegalPersonShort,
			LegalPersonName:               "Short Name",
		},
	}
	name, err = vasp.Name()
	require.NoError(t, err)
	require.Equal(t, "Short Name", name)

	// Legal name is selected if no short name or trading name
	vasp.Entity.Name.NameIdentifiers = []*ivms101.LegalPersonNameId{
		{
			LegalPersonNameIdentifierType: -1,
			LegalPersonName:               "Invalid Name",
		},
		{
			LegalPersonNameIdentifierType: ivms101.LegalPersonLegal,
			LegalPersonName:               "Legal Name",
		},
	}
	name, err = vasp.Name()
	require.NoError(t, err)
	require.Equal(t, "Legal Name", name)
}

func TestValidate(t *testing.T) {
	// VASP must contain an ID if partial is false
	vasp := &pb.VASP{}
	require.Error(t, vasp.Validate(false))
	vasp.Id = "vasp_id"

	// VASP must contain a registered directory if it has an ID
	checkValidateError(t, vasp)
	vasp.RegisteredDirectory = "trisatest.net"

	// VASP must contain an entity
	checkValidateError(t, vasp)
	vasp.Entity = &ivms101.LegalPerson{}

	// VASP entity must be ivms101 valid
	checkValidateError(t, vasp)
	data, err := ioutil.ReadFile(filepath.Join("..", "..", "..", "..", "ivms101", "testdata", "legalperson.json"))
	require.NoError(t, err)
	entity := &ivms101.LegalPerson{}
	require.NoError(t, json.Unmarshal(data, entity))
	require.NoError(t, entity.Validate())
	vasp.Entity = entity

	// VASP must contain at least one contact
	checkValidateError(t, vasp)
	vasp.Contacts = &pb.Contacts{
		Technical: &pb.Contact{},
	}

	// VASP must contain at least one contact with a name and email address
	checkValidateError(t, vasp)
	vasp.Contacts.Technical.Name = "Jason Bourne"
	vasp.Contacts.Technical.Email = "technical@example.com"

	// VASP must contain a TRISA endpoint and common name
	checkValidateError(t, vasp)
	vasp.CommonName = "CommonName"
	vasp.TrisaEndpoint = "trisa.example.com"

	// TRISA endpoint must be in a valid host:port format
	checkValidateError(t, vasp)
	vasp.TrisaEndpoint = ":12345"
	checkValidateError(t, vasp)
	vasp.TrisaEndpoint = "trisa.example.com:12345"

	// VASP cannot be verified without the verified timestamp
	vasp.VerificationStatus = pb.VerificationState_VERIFIED
	checkValidateError(t, vasp)
	vasp.VerifiedOn = "2020-01-02T00:00:00Z"

	// VASP must contain FirstListed and LastUpdated timestamp for full validation
	require.Error(t, vasp.Validate(false))
	vasp.FirstListed = "2020-01-01T00:00:00Z"
	vasp.LastUpdated = "2020-01-02T00:00:00Z"

	// VASP must contain a search signature for full validation
	require.Error(t, vasp.Validate(false))
	vasp.Signature = []byte("foo")
	checkValidateNoError(t, vasp)
}

// checkValidateError verifies that the VASP fails with both partial and full
// validation
func checkValidateError(t *testing.T, vasp *pb.VASP) {
	require.Error(t, vasp.Validate(true))
	require.Error(t, vasp.Validate(false))
}

// checkValidateNoError verifies that the VASP succeeds with both partial and full
// validation
func checkValidateNoError(t *testing.T, vasp *pb.VASP) {
	require.NoError(t, vasp.Validate(true))
	require.NoError(t, vasp.Validate(false))
}

func TestContactsValidation(t *testing.T) {
	contacts := &pb.Contacts{}
	require.EqualError(t, contacts.Validate(), "no contact specified on the VASP entity", "at least one non-nil contact is required")

	contacts = &pb.Contacts{
		Administrative: &pb.Contact{},
		Technical:      &pb.Contact{},
		Legal:          &pb.Contact{},
		Billing:        &pb.Contact{},
	}
	require.EqualError(t, contacts.Validate(), "no contact specified on the VASP entity", "at least one non-zero-valued contact is required")

	contacts.Administrative.Name = "k"
	require.EqualError(t, contacts.Validate(), "administrative contact invalid: contact name must be longer than one character", "non-zero administrative contact must be valid")

	contacts.Administrative.Name = "Kreg Balin"
	contacts.Administrative.Email = "kreg@example.com"
	contacts.Technical.Name = "Visual Nygard"
	require.EqualError(t, contacts.Validate(), "technical contact invalid: contact email is required", "non-zero technical contact must be valid")

	contacts.Technical.Email = "nygard@example.com"
	contacts.Legal.Email = "roger@example.com"
	require.EqualError(t, contacts.Validate(), "legal contact invalid: contact name is required", "non-zero legal contact must be valid")

	contacts.Legal.Name = "Roger Rabbit"
	contacts.Billing.Name = "Jessica Jones"
	contacts.Billing.Email = "notanemail"
	require.EqualError(t, contacts.Validate(), "billing contact invalid: could not parse email address", "non-zero billing contact must be valid")

	contacts.Billing.Email = "jessica@example.com"
	require.NoError(t, contacts.Validate(), "expected 4 valid contacts")

	require.NoError(t, (&pb.Contacts{Administrative: contacts.Administrative}).Validate(), "contacts should be valid with one valid administrative contact")
	require.NoError(t, (&pb.Contacts{Technical: contacts.Technical}).Validate(), "contacts should be valid with one valid technical contact")
	require.NoError(t, (&pb.Contacts{Legal: contacts.Legal}).Validate(), "contacts should be valid with one valid legal contact")
	require.NoError(t, (&pb.Contacts{Billing: contacts.Billing}).Validate(), "contacts should be valid with one valid billing contact")
}

func TestContactValidation(t *testing.T) {
	contact := &pb.Contact{}

	contact.Name = ""
	contact.Email = "foo@example.com"
	require.EqualError(t, contact.Validate(), "contact name is required")

	contact.Name = "d"
	require.EqualError(t, contact.Validate(), "contact name must be longer than one character")

	contact.Name = "Crazy Joe"
	contact.Email = ""
	require.EqualError(t, contact.Validate(), "contact email is required")

	contact.Email = "verybad%notanemail.com"
	require.EqualError(t, contact.Validate(), "could not parse email address")

	contact.Name = "Darlene Frederick"
	contact.Email = "darlene@example.com"
	contact.Phone = ""
	require.NoError(t, contact.Validate(), "only name and email are required")

	contact.Phone = "+18882123921"
	require.NoError(t, contact.Validate(), "only name and email are validated")
}
