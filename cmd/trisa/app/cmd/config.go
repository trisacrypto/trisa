package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/trisacrypto/trisa/pkg/ca"
	"github.com/trisacrypto/trisa/pkg/trisa/config"
)

var (
	csrPath   = "/etc/trisa"
	keyFile   = "server.key"
	crtFile   = "server.crt"
	trustFile = "trust.chain"

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
		NewConfigGenerateCmd(),
	)

	cmd.PersistentFlags().StringVar(&csrPath, "path", csrPath, "Path where private key and CSR is created")

	return cmd
}

func NewConfigGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate configuration",
		RunE:  runConfigGenerateCmd,
	}

	cmd.Flags().StringVar(&listenAddr, "listen", listenAddr, "Listen address")
	cmd.Flags().StringVar(&listenAddrAdmin, "listen-admin", listenAddrAdmin, "Listen address admin")

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
			CertificateFile: crtFile,
			PrivateKeyFile:  keyFile,
			TrustChainFile:  trustFile,
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
