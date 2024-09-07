package ivms101_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	reflect "reflect"
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

func TestIVMS101DatabaseScan(t *testing.T) {
	type Object struct {
		Identity      *ivms101.IdentityPayload `json:"identity"`
		Person        *ivms101.Person          `json:"person"`
		NaturalPerson *ivms101.NaturalPerson   `json:"naturalPerson"`
		LegalPerson   *ivms101.LegalPerson     `json:"legalPerson"`
		Address       *ivms101.Address         `json:"address"`
	}

	model := &Object{}
	mock := &MockRow{}
	err := mock.Open("testdata/identity_payload.json", "testdata/person_legal_person.json", "testdata/natural_person.json", "testdata/legal_person.json", "testdata/address.json")
	require.NoError(t, err, "could not load mock fixtures")
	require.Len(t, mock.raw, 5)

	err = mock.Scan(&model.Identity, &model.Person, &model.NaturalPerson, &model.LegalPerson, &model.Address)
	require.NoError(t, err)

	require.NotEmpty(t, model.Identity)
	require.NotEmpty(t, model.Person)
	require.NotEmpty(t, model.NaturalPerson)
	require.NotEmpty(t, model.LegalPerson)
	require.NotEmpty(t, model.Address)
}

func TestIVMS101DatabaseScanNil(t *testing.T) {
	type Object struct {
		Identity      *ivms101.IdentityPayload `json:"identity"`
		Person        *ivms101.Person          `json:"person"`
		NaturalPerson *ivms101.NaturalPerson   `json:"naturalPerson"`
		LegalPerson   *ivms101.LegalPerson     `json:"legalPerson"`
		Address       *ivms101.Address         `json:"address"`
	}

	model := &Object{}
	mock := &MockRow{}
	err := mock.Open("", "", "", "", "")
	require.NoError(t, err, "could not load mock fixtures")
	require.Len(t, mock.raw, 5)

	err = mock.Scan(&model.Identity, &model.Person, &model.NaturalPerson, &model.LegalPerson, &model.Address)
	require.NoError(t, err)

	require.Empty(t, model.Identity)
	require.Empty(t, model.Person)
	require.Empty(t, model.NaturalPerson)
	require.Empty(t, model.LegalPerson)
	require.Empty(t, model.Address)
}

func TestIVMS101DatabaseScanEmptyBytes(t *testing.T) {
	type Object struct {
		Identity      *ivms101.IdentityPayload `json:"identity"`
		Person        *ivms101.Person          `json:"person"`
		NaturalPerson *ivms101.NaturalPerson   `json:"naturalPerson"`
		LegalPerson   *ivms101.LegalPerson     `json:"legalPerson"`
		Address       *ivms101.Address         `json:"address"`
	}

	model := &Object{}
	mock := &MockRow{
		raw: [][]byte{{}, {}, {}, {}, {}},
	}

	err := mock.Scan(&model.Identity, &model.Person, &model.NaturalPerson, &model.LegalPerson, &model.Address)
	require.EqualError(t, err, "unexpected end of JSON input")
}

func TestIVMS101DatabaseScanNullJSON(t *testing.T) {
	type Object struct {
		Identity      *ivms101.IdentityPayload `json:"identity"`
		Person        *ivms101.Person          `json:"person"`
		NaturalPerson *ivms101.NaturalPerson   `json:"naturalPerson"`
		LegalPerson   *ivms101.LegalPerson     `json:"legalPerson"`
		Address       *ivms101.Address         `json:"address"`
	}

	model := &Object{}
	mock := &MockRow{
		raw: [][]byte{{110, 117, 108, 108}, {110, 117, 108, 108}, {110, 117, 108, 108}, {110, 117, 108, 108}, {110, 117, 108, 108}},
	}

	err := mock.Scan(&model.Identity, &model.Person, &model.NaturalPerson, &model.LegalPerson, &model.Address)
	require.NoError(t, err)
	require.Empty(t, model.Identity)
	require.Empty(t, model.Person)
	require.Empty(t, model.NaturalPerson)
	require.Empty(t, model.LegalPerson)
	require.Empty(t, model.Address)
}

type MockRow struct {
	raw [][]byte
}

var errNilPtr = errors.New("destination pointer is nil") // embedded in descriptive error

func (r *MockRow) Open(paths ...string) error {
	r.raw = make([][]byte, 0, len(paths))
	for _, path := range paths {
		if err := r.open(path); err != nil {
			return err
		}
	}
	return nil
}

func (r *MockRow) open(path string) (err error) {
	if path == "" {
		r.raw = append(r.raw, nil)
		return nil
	}

	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return err
	}
	r.raw = append(r.raw, data)
	return nil
}

func (r *MockRow) Scan(dest ...any) error {
	if len(dest) != len(r.raw) {
		return fmt.Errorf("sql: expected %d destination arguments in Scan, not %d", len(r.raw), len(dest))
	}

	for i, raw := range r.raw {
		dst := dest[i]

		if raw == nil {
			dst = nil
			continue
		}

		if scanner, ok := dst.(sql.Scanner); ok {
			if err := scanner.Scan(raw); err != nil {
				return err
			}
		}

		dpv := reflect.ValueOf(dst)
		if dpv.Kind() != reflect.Pointer {
			return errors.New("destination not a pointer")
		}
		if dpv.IsNil() {
			return errNilPtr
		}

		dv := reflect.Indirect(dpv)
		switch dv.Kind() {
		case reflect.Pointer:
			if dst == nil {
				dv.SetZero()
			}
			dv.Set(reflect.New(dv.Type().Elem()))
			dvi := dv.Interface()
			if scanner, ok := dvi.(sql.Scanner); ok {
				if err := scanner.Scan(raw); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func loadFixture(path string, obj interface{}) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(obj)
}
