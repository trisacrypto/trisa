package trust_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/trust"
	"software.sslmate.com/src/go-pkcs12"
)

func TestNewSerializer(t *testing.T) {
	stream, err := trust.NewSerializer(false)
	require.NoError(t, err)
	require.False(t, stream.Private)
	require.Equal(t, stream.Format, trust.CompressionAuto)
	require.Empty(t, stream.Password)

	stream, err = trust.NewSerializer(true)
	require.NoError(t, err)
	require.True(t, stream.Private)
	require.Equal(t, pkcs12.DefaultPassword, stream.Password)
	require.Equal(t, stream.Format, trust.CompressionAuto)

	stream, err = trust.NewSerializer(true, "supersecret")
	require.NoError(t, err)
	require.True(t, stream.Private)
	require.Equal(t, "supersecret", stream.Password)
	require.Equal(t, stream.Format, trust.CompressionAuto)

	stream, err = trust.NewSerializer(false, "", trust.CompressionZIP)
	require.NoError(t, err)
	require.False(t, stream.Private)
	require.Empty(t, stream.Password)
	require.Equal(t, stream.Format, trust.CompressionZIP)

	stream, err = trust.NewSerializer(true, "supersecret", trust.CompressionNone)
	require.NoError(t, err)
	require.True(t, stream.Private)
	require.Equal(t, "supersecret", stream.Password)
	require.Equal(t, stream.Format, trust.CompressionNone)

	stream, err = trust.NewSerializer(true, "", "foo")
	require.Nil(t, stream)
	require.Error(t, err)
}

func TestSerializer(t *testing.T) {
	// Test Serializer on Sectigo data, e.g. a client can read a private Provider from
	// a PCKS12 encrypted file and that the Directory service can extract the Public
	// keys and write them in a gzip format that can then be decrompressed and loaded by
	// other clients. This should exercise the primary serialization code.

	serializer, err := trust.NewSerializer(true, "qDhAwnfMjgDEzzUC")
	require.NoError(t, err)

	// Decode from Sectigo encoding
	provider, err := serializer.ReadFile("testdata/128106.zip")
	require.NoError(t, err)
	require.True(t, provider.IsPrivate())

	// Decode as a pool to test pool file reader
	providerPool, err := serializer.ReadPoolFile("testdata/128106.zip")
	require.NoError(t, err)
	for _, v := range providerPool {
		require.False(t, v.IsPrivate())
	}

	// Create temporary file for reading and writing
	f, err := ioutil.TempFile("", "128106*.gz")
	require.NoError(t, err)
	path := f.Name()
	f.Close()
	defer os.Remove(path)
	t.Logf("created temporary file at %s", path)

	// Compress public provider
	serializer, err = trust.NewSerializer(false)
	require.NoError(t, err)
	provData, err := serializer.Compress(provider.Public())
	require.NoError(t, err)
	require.Equal(t, len(provData), 4402)

	// Write public provider to gzip file
	err = serializer.WriteFile(provider.Public(), path)
	require.NoError(t, err)

	// Read public provider from gzip file
	serializer, err = trust.NewSerializer(false, "", trust.CompressionGZIP)
	require.NoError(t, err)
	o, err := serializer.ReadFile(path)
	require.NoError(t, err)
	require.False(t, o.IsPrivate())

	pb, err := provider.Public().Encode()
	require.NoError(t, err)

	ob, err := o.Encode()
	require.NoError(t, err)

	require.True(t, bytes.Equal(pb, ob))

	// Compress provider pool
	pool := trust.NewPool(provider.Public())
	poolData, err := serializer.CompressPool(pool)
	require.NoError(t, err)
	require.Equal(t, len(poolData), 4402)

	// Write provider pool to gzip file
	err = serializer.WritePoolFile(pool, path)
	require.NoError(t, err)
}
