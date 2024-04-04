package ivms101_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	"google.golang.org/protobuf/proto"
)

func TestIVMS101Database(t *testing.T) {

	t.Run("IdentityPayload", func(t *testing.T) {
		ip := &ivms101.IdentityPayload{}
		err := loadFixture("testdata/identity_payload.json", ip)
		require.NoError(t, err, "could not load identity payload fixture")

		value, err := ip.Value()
		require.NoError(t, err, "could not fetch value for identity payload")

		cp := &ivms101.IdentityPayload{}
		err = cp.Scan(value)
		require.NoError(t, err, "could not scan value into identity payload")
		require.True(t, proto.Equal(ip, cp), "loaded value not equal to copied value")
	})

	t.Run("Person", func(t *testing.T) {
		ip := &ivms101.Person{}
		err := loadFixture("testdata/person_legal_person.json", ip)
		require.NoError(t, err, "could not load legal person fixture")

		value, err := ip.Value()
		require.NoError(t, err, "could not fetch value for person")

		cp := &ivms101.Person{}
		err = cp.Scan(value)
		require.NoError(t, err, "could not scan value into person")
		require.True(t, proto.Equal(ip, cp), "loaded value not equal to copied value")
	})

	t.Run("LegalPerson", func(t *testing.T) {
		ip := &ivms101.LegalPerson{}
		err := loadFixture("testdata/legal_person.json", ip)
		require.NoError(t, err, "could not load legal person fixture")

		value, err := ip.Value()
		require.NoError(t, err, "could not fetch value for legal person")

		cp := &ivms101.LegalPerson{}
		err = cp.Scan(value)
		require.NoError(t, err, "could not scan value into legal person")
		require.True(t, proto.Equal(ip, cp), "loaded value not equal to copied value")
	})

	t.Run("NaturalPerson", func(t *testing.T) {
		ip := &ivms101.NaturalPerson{}
		err := loadFixture("testdata/natural_person.json", ip)
		require.NoError(t, err, "could not load natural person fixture")

		value, err := ip.Value()
		require.NoError(t, err, "could not fetch value for natural person")

		cp := &ivms101.NaturalPerson{}
		err = cp.Scan(value)
		require.NoError(t, err, "could not scan value into natural person")

		require.True(t, proto.Equal(ip, cp), "loaded value not equal to copied value")
	})

	t.Run("Address", func(t *testing.T) {
		ip := &ivms101.Address{}
		err := loadFixture("testdata/address.json", ip)
		require.NoError(t, err, "could not load address fixture")

		value, err := ip.Value()
		require.NoError(t, err, "could not fetch value for address")

		cp := &ivms101.Address{}
		err = cp.Scan(value)
		require.NoError(t, err, "could not scan value into address")
		require.True(t, proto.Equal(ip, cp), "loaded value not equal to copied value")
	})
}

func loadFixture(path string, obj interface{}) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(obj)
}
