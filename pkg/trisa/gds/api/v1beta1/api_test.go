package api_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestRegistrationForm(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/alice-trixo.json")
	require.NoError(t, err)

	form := &api.RegisterRequest{}
	err = protojson.Unmarshal(data, form)
	require.NoError(t, err)

	// Compare JSON
	opts := protojson.MarshalOptions{Multiline: true, Indent: "  ", UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: false}
	out, err := opts.Marshal(form)
	require.NoError(t, err)

	ioutil.WriteFile("testdata/alice-trixo-parsed.json", out, 0644)
}
