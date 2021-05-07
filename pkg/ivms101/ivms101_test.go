package ivms101_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
)

func TestLegalPerson(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/legalperson.json")
	require.NoError(t, err)

	// Should be able to load a valid legal person from JSON data
	var person *ivms101.LegalPerson
	require.NoError(t, json.Unmarshal(data, &person))
	require.NoError(t, person.Validate())

	// Should be able to convert a legal person into a generic Persion
	gp := person.Person()
	require.Nil(t, gp.GetNaturalPerson())
	require.Equal(t, person, gp.GetLegalPerson())
}

func TestNaturalPerson(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/naturalperson.json")
	require.NoError(t, err)

	// Should be able to load a valid legal person from JSON data
	var person *ivms101.NaturalPerson
	require.NoError(t, json.Unmarshal(data, &person))
	require.NoError(t, person.Validate())

	// Should be able to convert a legal person into a generic Persion
	gp := person.Person()
	require.Nil(t, gp.GetLegalPerson())
	require.Equal(t, person, gp.GetNaturalPerson())
}
