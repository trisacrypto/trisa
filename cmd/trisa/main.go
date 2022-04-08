/*
 The TRISA CLI client allows you to create and execute TRISA requests from the command
 line for development or testing purposes. For more information on how to use the CLI,
 run `trisa --help` or see the documenation at https://trisa.dev.
*/
package main

import (
	"os"

	"github.com/trisacrypto/trisa/pkg"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "trisa"
	app.Usage = "create and execute TRISA requests"
	app.Version = pkg.Version()
	app.Before = initClient
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c"},
			Usage:   "specify the path to a configuration file to use",
			EnvVars: []string{"TRISA_CONF"},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "transfer",
			Usage:  "execute a TRISA transfer with a TRISA peer",
			Action: transfer,
			Flags:  []cli.Flag{},
		},
		{
			Name:    "exchange",
			Aliases: []string{"key-exchange", "keys"},
			Usage:   "manage profiles",
			Action:  exchange,
			Flags:   []cli.Flag{},
		},
	}
	app.Run(os.Args)
}

//====================================================================================
// RPC Commands
//====================================================================================

func transfer(c *cli.Context) (err error) {
	return nil
}

func exchange(c *cli.Context) (err error) {
	return nil
}

//====================================================================================
// Helper Commands
//====================================================================================

func initClient(c *cli.Context) (err error) {
	return nil
}
