package cli_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/trisa/cli"
)

func TestCreateDefaultProfiles(t *testing.T) {
	// Create temporary configuration directory
	confDir := setupConfigDir(t)

	// Create new profiles from scratch
	profiles := cli.New()
	profiles.Profiles[cli.ProfileDefault] = &cli.Profile{
		Endpoint: "trisa.dev:10001",
	}

	// Save the profiles to disk - should save to the temporary directory
	err := profiles.Save()
	require.NoError(t, err, "could not save profiles")
	require.FileExists(t, filepath.Join(confDir, cli.ProfileYAML), "profiles configuration not created")

	// Test loading profiles that were saved
	other, err := cli.Load()
	require.NoError(t, err, "could not load saved profiles")
	require.Equal(t, profiles, other)
}

func TestLoadProfiles(t *testing.T) {
	setConfigPath(t)
	profiles, err := cli.Load()
	require.NoError(t, err, "could not load profiles")
	require.Len(t, profiles.Profiles, 2, "incorrect profiles loaded")
}

func TestProfilesValidate(t *testing.T) {
	profiles := &cli.Profiles{}
	err := profiles.Validate()
	require.EqualError(t, err, cli.ErrIncorrectVersion.Error())

	profiles.Version = cli.ProfileVersion
	err = profiles.Validate()
	require.EqualError(t, err, cli.ErrNoActiveProfile.Error())

	profiles.Active = "foo"
	err = profiles.Validate()
	require.EqualError(t, err, cli.ErrNoActiveProfile.Error())

	profiles.Profiles = make(map[string]*cli.Profile)
	profiles.Profiles["foo"] = &cli.Profile{}
	err = profiles.Validate()
	require.EqualError(t, err, `invalid profile "foo": valid endpoint is required`)

	profiles.Profiles["foo"].Endpoint = "trisa.dev:10001"

	err = profiles.Validate()
	require.NoError(t, err, "expected profiles to be valid at this point")
}

func setupConfigDir(t *testing.T) string {
	tmpdir, err := ioutil.TempDir("", "trisacli-config-*")
	require.NoError(t, err, "could not create temporary configuration dir")

	current := os.Getenv(cli.ProfileDirEnv)
	os.Setenv(cli.ProfileDirEnv, tmpdir)

	t.Cleanup(func() {
		os.Setenv(cli.ProfileDirEnv, current)
		os.RemoveAll(tmpdir)
		t.Logf("cleaned up %s\n", tmpdir)
	})

	t.Logf("created $%s at %s", cli.ProfileDirEnv, tmpdir)
	return tmpdir
}

func setConfigPath(t *testing.T) {
	current := os.Getenv(cli.ProfilePathEnv)
	path, err := filepath.Abs("testdata/config.yaml")
	require.NoError(t, err, "could not get absolute path to testdata")
	os.Setenv(cli.ProfilePathEnv, path)

	t.Cleanup(func() {
		os.Setenv(cli.ProfilePathEnv, current)
	})
}
