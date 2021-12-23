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

	// VASP must contain at least one contact with an email address
	checkValidateError(t, vasp)
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
