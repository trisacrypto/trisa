package cmd

import (
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/trisacrypto/trisa/pkg/ca"
	"github.com/trisacrypto/trisa/pkg/trisa/config"
)

var (
	csrPath = "/etc/trisa"
	keyFile = "server.key"
	csrFile = "server.csr"
	crtFile = "server.crt"

	initHostname    string
	initOrg         string
	trustedRootCA   string
	listenAddr      string
	listenAddrAdmin string
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config management",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			stat, err := os.Stat(csrPath)
			if err != nil {
				return err
			}
			if !stat.IsDir() {
				return fmt.Errorf("%s is not a directory", csrPath)
			}
			return nil
		},
	}

	cmd.AddCommand(
		NewConfigCSRCmd(),
		NewConfigGenerateCmd(),
	)

	cmd.PersistentFlags().StringVar(&csrPath, "path", csrPath, "Path where private key and CSR is created")

	return cmd
}

func NewConfigCSRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "csr",
		Short: "Create Certificate Signing Request",
		RunE:  runConfigCSRCmd,
	}

	cmd.Flags().StringVar(&initHostname, "hostname", "", "Hostname where the server will be running")
	cmd.Flags().StringVar(&initOrg, "org", "", "Organization")

	cmd.MarkFlagRequired("hostname")
	cmd.MarkFlagRequired("org")

	return cmd
}

func runConfigCSRCmd(cmd *cobra.Command, args []string) error {
	key, err := ca.GenerateRSAPrivateKey(4096)
	if err != nil {
		return err
	}

	csr, err := ca.CreateCSR(pkix.Name{
		CommonName:   initHostname,
		Organization: []string{initOrg},
	}, key)

	if err != nil {
		return err
	}

	keyPEM, err := ca.PEMEncodePrivateKey(key)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(csrPath, keyFile), keyPEM, os.ModePerm); err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(csrPath, csrFile), csr, os.ModePerm)
}

func NewConfigGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate configuration",
		RunE:  runConfigGenerateCmd,
	}

	cmd.Flags().StringVar(&listenAddr, "listen", listenAddr, "Listen address")
	cmd.Flags().StringVar(&listenAddrAdmin, "listen-admin", listenAddrAdmin, "Listen address admin")
	cmd.Flags().StringVar(&trustedRootCA, "trust", "", "Trusted root CA URL")

	cmd.MarkFlagRequired("trust")

	return cmd
}

func runConfigGenerateCmd(cmd *cobra.Command, args []string) error {

	// Decode certificate so we can deduct the hostname from CN.
	crtPEM, err := ioutil.ReadFile(filepath.Join(csrPath, crtFile))
	if err != nil {
		return err
	}

	crt, err := ca.PEMDecodeCertificate(crtPEM)
	if err != nil {
		return err
	}

	c := &config.Config{
		TLS: &config.TLS{
			CertificateFile: filepath.Join(csrPath, crtFile),
			PrivateKeyFile:  filepath.Join(csrPath, keyFile),
			TrustedRootCAs:  []string{trustedRootCA},
		},
		Server: &config.Server{
			ListenAddress:      listenAddr,
			ListenAddressAdmin: listenAddrAdmin,
			Hostname:           crt.Subject.CommonName,
		},
	}

	if err := c.Save(configFile); err != nil {
		return err
	}

	fmt.Printf("config stored in %s\n", configFile)
	return nil
}
