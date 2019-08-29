package app

import (
	"io"

	"github.com/trisacrypto/trisa/cmd/trisa/app/cmd"
)

func Run(out, err io.Writer) error {
	c := cmd.NewTRISACommand(out, err)
	return c.Execute()
}
