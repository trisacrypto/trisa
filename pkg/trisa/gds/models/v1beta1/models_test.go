package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"
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
