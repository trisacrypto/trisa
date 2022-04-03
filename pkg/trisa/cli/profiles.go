package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/shibukawa/configdir"
	"gopkg.in/yaml.v3"
)

const (
	ProfileYAML     = "config.yaml"
	ProfileVersion  = "v1"
	ProfileDefault  = "default"
	ProfileDirEnv   = "TRISA_CONF_DIR"
	ProfilePathEnv  = "TRISA_CONF"
	vendorName      = "trisa"
	applicationName = "trisacli"
)

// Profiles allows multiple independent configurations to specify different connections
// to different TRISA nodes. The currently active profile and profile version is
// specified at the top of the YAML file to assist with parsing and validation. Every
// profile should have a profile that is the "default" profile; if a configuration is
// missing a value, the value in the "default" profile is used.
//
// The profiles object is the primary access point for all TRISA CLI operations, hence
// it is the object that is returned from the top-level functions in the package.
type Profiles struct {
	Version     string              `yaml:"version"`
	Active      string              `yaml:"active"`
	Profiles    map[string]*Profile `yaml:"profiles"`
	LastUpdated time.Time           `yaml:"lasted_updated,omitempty"`
	folder      *configdir.Config
	conf        string
}

// Profile represents a configuration used to manage the connection to a single TRISA
// node and its associated directory service.
type Profile struct {
	Endpoint    string    `yaml:"endpoint,omitempty"`
	Certs       string    `yaml:"certs,omitempty"`
	Chain       string    `yaml:"chain,omitempty"`
	SealingKey  string    `yaml:"sealing_key,omitempty"`
	LastUpdated time.Time `yaml:"last_updated,omitempty"`
	Info        struct {
		ID                  string    `yaml:"vasp_id,omitempty"`
		Name                string    `yaml:"name,omitempty"`
		CommonName          string    `yaml:"common_name,omitempty"`
		RegisteredDirectory string    `yaml:"registered_directory,omitempty"`
		ResolvedAt          time.Time `yaml:"resolved_at,omitempty"`
	} `yaml:"directory_info,omitempty"`
	folder *configdir.Config
}

// New returns a new configuration with an empty default profile. The returned Profile
// is not valid but can be updated interactively and validated before saving to disk.
// For most cases, use Load() instead of New() to fetch the profile from disk.
func New() *Profiles {
	return &Profiles{
		Version: ProfileVersion,
		Active:  ProfileDefault,
		Profiles: map[string]*Profile{
			ProfileDefault: {},
		},
	}
}

//====================================================================================
// Load and Save Profiles
//====================================================================================

// Load the profiles from disk searching the configuration directories. The profiles are
// central to all CLI applications, so loading it is the first step to running a command.
func Load() (_ *Profiles, err error) {
	// If the $TRISA_CONF environment variable is set, load the profile from the path
	if path := os.Getenv(ProfilePathEnv); path != "" {
		return LoadPath(path)
	}

	// Create the configdir search which searches the local path first, then the user
	// configuration folder, then the system configuration folder, returning an error if
	// one of the configurations is not found.
	cfgd := configdir.New(vendorName, applicationName)

	// Set the local path to the $TRISA_CONF_DIR environment variable or the local path
	// if that environment variable is not set; this will set the search priority.
	if cfgd.LocalPath = os.Getenv(ProfileDirEnv); cfgd.LocalPath == "" {
		cfgd.LocalPath, _ = filepath.Abs(".")
	}

	var folder *configdir.Config
	if folder = cfgd.QueryFolderContainsFile(ProfileYAML); folder == nil {
		return nil, ErrProfileNotFound
	}
	return load(folder, ProfileYAML)
}

// LoadPath loads a profile using the specified path.
func LoadPath(path string) (_ *Profiles, err error) {
	// Ensure the file exists and is a file, otherwise configdir parent directory search will fail
	var info fs.FileInfo
	if info, err = os.Stat(path); os.IsNotExist(err) {
		return nil, ErrInvalidProfilePath
	}

	if info.IsDir() {
		return nil, ErrInvalidProfilePath
	}

	// Get the basename and the directory name to create the config folder
	parent := filepath.Dir(path)
	name := filepath.Base(path)

	cfgd := configdir.New(vendorName, applicationName)
	cfgd.LocalPath = parent

	var folder *configdir.Config
	if folder = cfgd.QueryFolderContainsFile(name); folder == nil {
		return nil, ErrProfileNotFound
	}
	return load(folder, name)
}

func load(folder *configdir.Config, conf string) (profiles *Profiles, err error) {
	if conf == "" {
		conf = ProfileYAML
	}

	// Load and unmarshal the data from the configuration
	var data []byte
	if data, err = folder.ReadFile(conf); err != nil {
		return nil, fmt.Errorf("could not read configuration: %s", err)
	}

	profiles = &Profiles{folder: folder, conf: conf}
	if err = yaml.Unmarshal(data, profiles); err != nil {
		return nil, fmt.Errorf("could not unmarshal configuration: %s", err)
	}

	if err = profiles.Validate(); err != nil {
		return nil, err
	}
	return profiles, nil
}

// Save the profiles to disk in the specified configuration folder.
func (p *Profiles) Save() (err error) {
	if err = p.resolveFolder(); err != nil {
		return err
	}

	p.LastUpdated = time.Now()

	var data []byte
	if data, err = yaml.Marshal(p); err != nil {
		return fmt.Errorf("could not marshal profiles: %s", err)
	}

	if err = p.folder.WriteFile(p.conf, data); err != nil {
		return fmt.Errorf("could not write profiles to disk: %s", err)
	}
	return nil
}

//====================================================================================
// Validation
//====================================================================================

// Validate the profiles are configured correctly
func (p *Profiles) Validate() (err error) {
	if p.Version != ProfileVersion {
		return ErrIncorrectVersion
	}

	if p.Active == "" {
		return ErrNoActiveProfile
	}

	if _, ok := p.Profiles[p.Active]; !ok {
		return ErrNoActiveProfile
	}

	for name, profile := range p.Profiles {
		if err = profile.Validate(); err != nil {
			return fmt.Errorf("invalid profile %q: %s", name, err)
		}
	}
	return nil
}

// Validate a single profile
func (p *Profile) Validate() (err error) {
	if p.Endpoint == "" {
		return ErrInvalidEndpoint
	}
	return nil
}

//====================================================================================
// CLI Handlers
//====================================================================================

// GetActive returns the currently active profile and gives it access to the folder.
func (p *Profiles) GetActive() *Profile {
	profile := p.Profiles[p.Active]
	profile.folder = p.folder
	return profile
}

// SetActive activates the profile with the specified name if it exists and saves the
// profiles back to disk to ensure the profile remains active on the next load.
func (p *Profiles) SetActive(name string) (err error) {
	if _, ok := p.Profiles[name]; !ok {
		return fmt.Errorf("no profile named %q found", name)
	}

	p.Active = name
	if err = p.Save(); err != nil {
		return err
	}
	return nil
}

func (p *Profiles) Path() (path string, err error) {
	if err = p.resolveFolder(); err != nil {
		return "", err
	}
	return filepath.Join(p.folder.Path, p.conf), nil
}

// SetPath to a configuration file. This should only be used on install; if the profile
// already has a folder and a path, then an error will be returned.
func (p *Profiles) SetPath(path string) error {
	if p.folder != nil || p.conf != "" {
		return ErrCannotSetPath
	}

	p.folder = &configdir.Config{
		Path: filepath.Dir(path),
		Type: configdir.Existing,
	}
	p.conf = filepath.Base(path)
	return nil
}

//====================================================================================
// Helper Functions
//====================================================================================

// Update the profiles folder to the correct configuration folder.
func (p *Profiles) resolveFolder() (err error) {
	if p.folder == nil {
		if p.folder, err = findFolder(); err != nil {
			return err
		}
	}

	if p.conf == "" {
		p.conf = ProfileYAML
	}
	return nil
}

// search for best configuration location to save to including the $TRISA_CONF_DIR
func findFolder() (folder *configdir.Config, err error) {
	cfgd := configdir.New(vendorName, applicationName)
	cfgd.LocalPath = os.Getenv(ProfileDirEnv)

	if folder = cfgd.QueryFolderContainsFile(ProfileYAML); folder != nil {
		return folder, nil
	}

	// Use the local path from the environment if it is supplied
	if cfgd.LocalPath != "" {
		return &configdir.Config{Path: cfgd.LocalPath, Type: configdir.Local}, nil
	}

	// If there is no current folder with the configuration, attempt to
	// create the configuration in the user directory.
	if folders := cfgd.QueryFolders(configdir.Global); len(folders) > 0 {
		return folders[0], nil
	}
	return nil, ErrNoProfileDirectory
}
