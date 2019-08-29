package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the Trisa Server configuration.
type Config struct {
	TLS     *TLS               `yaml:"tls,omitempty"`
	Server  *Server            `yaml:"server,omitempty"`
	Wallets map[string]*Wallet `yaml:"wallets,omitempty"`
}

type TLS struct {
	PrivateKeyFile  string   `yaml:"privateKeyFile,omitempty"`
	CertificateFile string   `yaml:"certificateFile,omitempty"`
	TrustedRootCAs  []string `yaml:"trustedRootCAs,omitempty"`
}

type Server struct {
	ListenAddress      string `yaml:"listenAddress,omitempty"`
	ListenAddressAdmin string `yaml:"listenAddressAdmin,omitempty"`
	Hostname           string `yaml:"hostname,omitempty"`
}

type Wallet struct {
	FirstName string
	LastName  string
}

func FromFile(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var c Config

	err = yaml.UnmarshalStrict(data, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) Save(file string) error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, out, os.ModePerm)
}
