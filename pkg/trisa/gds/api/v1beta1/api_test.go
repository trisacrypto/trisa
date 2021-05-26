package api_test

import (
	"io/ioutil"
	"os"
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

	// Compare JSON manually if required
	path := os.Getenv("TRISA_TEST_COMPARE_REGISTRATION_FORM")
	if path != "" {
		opts := protojson.MarshalOptions{Multiline: true, Indent: "  ", UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: false}
		out, err := opts.Marshal(form)
		require.NoError(t, err)

		ioutil.WriteFile(path, out, 0644)
	}
}
