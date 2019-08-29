package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

var (
	configFile string
)

func NewTRISACommand(out, err io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:  "trisa",
		Long: "VASP Travel Rule Information Sharing Architecture",
	}

	root.PersistentFlags().StringVarP(&configFile, "config", "c", "/etc/trisa/config.yaml", "Configuration file")

	root.AddCommand(
		NewServerCmd(),
		NewConfigCmd(),
	)

	return root
}
