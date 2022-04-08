/*
 The TRISA CLI client allows you to create and execute TRISA requests from the command
 line for development or testing purposes. For more information on how to use the CLI,
 run `trisa --help` or see the documenation at https://trisa.dev.
*/
package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/trisa/pkg"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/mtls"
	"github.com/trisacrypto/trisa/pkg/trust"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var (
	client    api.TRISANetworkClient
	directory gds.TRISADirectoryClient
)

func main() {
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "trisa"
	app.Usage = "create and execute TRISA requests"
	app.Version = pkg.Version()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "endpoint",
			Aliases: []string{"e"},
			Usage:   "the endpoint of the TRISA peer to connect to",
			EnvVars: []string{"TRISA_ENDPOINT"},
		},
		&cli.StringFlag{
			Name:    "directory",
			Aliases: []string{"d"},
			Usage:   "the endpoint or name of the directory service to use",
			EnvVars: []string{"TRISA_DIRECTORY"},
			Value:   "testnet",
		},
		&cli.StringFlag{
			Name:    "certs",
			Aliases: []string{"c"},
			Usage:   "specify the path to your mTLS TRISA identity certificates",
			EnvVars: []string{"TRISA_CERTS"},
		},
		&cli.StringFlag{
			Name:    "chain",
			Aliases: []string{"t"},
			Usage:   "the path to the trust chain without private keys if separate from the certs",
			EnvVars: []string{"TRISA_TRUST_CHAIN"},
		},
		&cli.StringFlag{
			Name:    "pkcs12password",
			Aliases: []string{"P"},
			Usage:   "the pkcs12 password of the certs if they are encrypted",
			EnvVars: []string{"TRISA_CERTS_PASSWORD"},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "make",
			Aliases: []string{"envelope"},
			Usage:   "create a secure envelope template from a payload",
			Action:  make,
			Flags:   []cli.Flag{},
		},
		{
			Name:   "transfer",
			Usage:  "execute a TRISA transfer with a TRISA peer",
			Before: initClient,
			Action: transfer,
			Flags:  []cli.Flag{},
		},
		{
			Name:    "exchange",
			Aliases: []string{"key-exchange"},
			Usage:   "exchange public sealing keys with a TRISA peer",
			Before:  initClient,
			Action:  exchange,
			Flags:   []cli.Flag{},
		},
		{
			Name:    "confirm",
			Aliases: []string{"confirm-address"},
			Usage:   "execute an address confirmation request with a TRISA peer",
			Before:  initClient,
			Action:  confirm,
			Flags:   []cli.Flag{},
		},
		{
			Name:    "status",
			Aliases: []string{"health-check"},
			Usage:   "execute a health check against a TRISA peer and directory service",
			Before:  initDirectoryClient,
			Action:  health,
			Flags:   []cli.Flag{},
		},
		{
			Name:   "lookup",
			Usage:  "lookup a TRISA record on the directory service",
			Before: initDirectoryClient,
			Action: lookup,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "search",
			Usage:  "search for a VASP on the directory service",
			Before: initDirectoryClient,
			Action: search,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "list",
			Usage:  "list the available VASPs on the directory service",
			Before: initDirectoryClient,
			Action: list,
			Flags:  []cli.Flag{},
		},
	}
	app.Run(os.Args)
}

//====================================================================================
// RPC Commands
//====================================================================================

func make(c *cli.Context) (err error) {
	return nil
}

func transfer(c *cli.Context) (err error) {
	return nil
}

func exchange(c *cli.Context) (err error) {
	return nil
}

func confirm(c *cli.Context) (err error) {
	return nil
}

func health(c *cli.Context) (err error) {
	return nil
}

func lookup(c *cli.Context) (err error) {
	return nil
}

func search(c *cli.Context) (err error) {
	return nil
}

func list(c *cli.Context) (err error) {
	return nil
}

//====================================================================================
// Helper Commands
//====================================================================================

func initClient(c *cli.Context) (err error) {
	endpoint := c.String("endpoint")
	if endpoint == "" {
		return cli.Exit("specify endpoint of TRISA peer to connect to", 1)
	}

	certs, pool, err := loadCerts(c)
	if err != nil {
		return err
	}

	creds, err := mtls.ClientCreds(endpoint, certs, pool)
	if err != nil {
		return cli.Exit(err, 1)
	}

	cc, err := grpc.Dial(endpoint, creds)
	if err != nil {
		return cli.Exit(err, 1)
	}

	client = api.NewTRISANetworkClient(cc)
	return nil
}

func initDirectoryClient(c *cli.Context) (err error) {
	return nil
}

func loadCerts(c *cli.Context) (certs *trust.Provider, pool trust.ProviderPool, err error) {
	// Get configuration for certificates
	certPath := c.String("certs")
	if certPath == "" {
		return nil, nil, cli.Exit("path to identity certificates required for this command", 1)
	}

	var sz *trust.Serializer
	if passwd := c.String("pkcs12password"); passwd != "" {
		if sz, err = trust.NewSerializer(true, passwd); err != nil {
			return nil, nil, cli.Exit(err, 1)
		}
	} else {
		if sz, err = trust.NewSerializer(false); err != nil {
			return nil, nil, err
		}
	}

	if certs, err = sz.ReadFile(certPath); err != nil {
		return nil, nil, cli.Exit(err, 1)
	}

	if chainPath := c.String("chain"); chainPath != "" {
		if pool, err = sz.ReadPoolFile(chainPath); err != nil {
			return nil, nil, cli.Exit(err, 1)
		}
	} else {
		if pool, err = sz.ReadPoolFile(certPath); err != nil {
			return nil, nil, cli.Exit(err, 1)
		}
	}

	return certs, pool, nil
}
