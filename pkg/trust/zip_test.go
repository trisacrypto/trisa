package trust_test

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
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

	serializer, err := trust.NewSerializer(true, "supersecretsquirrel")
	require.NoError(t, err)

	// Decode from Sectigo encoding
	provider, err := serializer.ReadFile("testdata/110831.zip")
	require.NoError(t, err)
	require.True(t, provider.IsPrivate())

	// Decode as a pool to test pool file reader
	providerPool, err := serializer.ReadPoolFile("testdata/110831.zip")
	require.NoError(t, err)
	for _, v := range providerPool {
		require.False(t, v.IsPrivate())
	}

	// Create temporary file for reading and writing
	f, err := ioutil.TempFile("", "110831*.gz")
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
	require.Equal(t, 2833, len(provData))

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
	require.Equal(t, 2833, len(poolData))

	// Write provider pool to gzip file
	err = serializer.WritePoolFile(pool, path)
	require.NoError(t, err)
}

func TestZipWriter(t *testing.T) {
	// Test creating encrypted and unencrypted zip files using the trust package
	// Create private and open serializers
	serializerA, err := trust.NewSerializer(true, "supersecretsquirrel")
	require.NoError(t, err)

	serializerB, err := trust.NewSerializer(false)
	require.NoError(t, err)

	// Decode the fixture to use
	provider, err := serializerA.ReadFile("testdata/110831.zip")
	require.NoError(t, err)

	pool := make(trust.ProviderPool)
	pool[provider.String()] = provider.Public()

	// Create temporary paths for serialization
	tmpdir := t.TempDir()
	var (
		p12path  = filepath.Join(tmpdir, "testprivate.zip")
		openpath = filepath.Join(tmpdir, "testopen.zip")
		p12pool  = filepath.Join(tmpdir, "testprivatepool.zip")
		openpool = filepath.Join(tmpdir, "testopenpool.zip")
	)

	// Write the compressed file back to disk with p12 encryption
	err = serializerA.WriteFile(provider, p12path)
	require.NoError(t, err, "could not write file")
	require.FileExists(t, p12path, "private file was not written to disk")
	RequireInZip(t, p12path, "docs.trisa.dev.p12")

	// Write the compressed file back to disk without p12 encryption
	err = serializerB.WriteFile(provider, openpath)
	require.NoError(t, err, "could not write open file to disk")
	require.FileExists(t, openpath, "open file was not written to disk")
	RequireInZip(t, openpath, "docs.trisa.dev.pem")

	// Write compressed pool back to disk with p12 encryption
	// TODO: this use case doesn't make a lot of sense, why encrypt public keys?
	err = serializerA.WritePoolFile(pool, p12pool)
	require.NoError(t, err, "could not write pool")
	require.FileExists(t, p12pool, "private pool was not written to disk")
	RequireInZip(t, p12pool, "docs.trisa.dev.p12")

	// Write compressed pool back to disk without p12 encryption
	err = serializerB.WritePoolFile(pool, openpool)
	require.NoError(t, err, "could not write open pool to disk")
	require.FileExists(t, openpool, "open pool was not written to disk")
	RequireInZip(t, openpool, "docs.trisa.dev.pem")
}

func RequireInZip(t *testing.T, zipPath, contains string) {
	r, err := zip.OpenReader(zipPath)
	require.NoError(t, err, "could not open %s", zipPath)
	defer r.Close()

	found := false
	for _, f := range r.File {
		if f.Name == contains {
			found = true
			break
		}
	}
	require.True(t, found, "could not find %q in %s", contains, zipPath)
}
