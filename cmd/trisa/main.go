/*
 The TRISA CLI client allows you to create and execute TRISA requests from the command
 line for development or testing purposes. For more information on how to use the CLI,
 run `trisa --help` or see the documenation at https://trisa.dev.
*/
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/trisa/pkg"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/mtls"
	"github.com/trisacrypto/trisa/pkg/trust"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Clients to connect to various services and make RPC requests.
var (
	peer      api.TRISANetworkClient
	ping      api.TRISAHealthClient
	directory gds.TRISADirectoryClient
)

// Directory service endpoints
const (
	testnet = "api.trisatest.net:443"
	mainnet = "api.vaspdirectory.net:443"
)

// Aliases that map directory names to endpoints
var directoryAliases = map[string]string{
	"testnet":           testnet,
	"trisatest":         testnet,
	"trisatest.net":     testnet,
	"mainnet":           mainnet,
	"vaspdirectory":     mainnet,
	"vaspdirectory.net": mainnet,
}

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
			EnvVars: []string{"TRISA_DIRECTORY", "TRISA_DIRECTORY_URL", "GDS_DIRECTORY_URL"},
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
			Action:  envelope,
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
			Before:  initHealthCheckClient,
			Action:  health,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "insecure",
					Aliases: []string{"S"},
					Usage:   "do not connect to the peer health check endpoint with mTLS",
				},
			},
		},
		{
			Name:      "lookup",
			Usage:     "lookup a TRISA record on the directory service",
			UsageText: "trisa lookup [-dir value] -[in]\nlookup a VASP by ID or common name\nspecify registered directory for alternative network issuers",
			Before:    initDirectoryClient,
			Action:    lookup,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "vasp-id",
					Aliases: []string{"id", "i"},
					Usage:   "the UUID of the VASP to lookup",
				},
				&cli.StringFlag{
					Name:    "registered-directory",
					Aliases: []string{"dir", "d"},
					Usage:   "the registered directory of the VASP record",
				},
				&cli.StringFlag{
					Name:    "common-name",
					Aliases: []string{"cn", "n"},
					Usage:   "the common name of the VASP to lookup",
				},
			},
		},
		{
			Name:   "search",
			Usage:  "search for a VASP on the directory service",
			Before: initDirectoryClient,
			Action: search,
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "one or more names of VASPs to search for",
				},
				&cli.StringSliceFlag{
					Name:    "website",
					Aliases: []string{"w"},
					Usage:   "one or more urls of VASP websites to search for",
				},
				&cli.StringSliceFlag{
					Name:    "country",
					Aliases: []string{"c"},
					Usage:   "one or more countries to filter requests on",
				},
				&cli.StringSliceFlag{
					Name:    "category",
					Aliases: []string{"C"},
					Usage:   "one or more categories to filter requests on",
				},
			},
		},
	}
	app.Run(os.Args)
}

//====================================================================================
// RPC Commands
//====================================================================================

func envelope(c *cli.Context) (err error) {
	return nil
}

func transfer(c *cli.Context) (err error) {
	req := &api.SecureEnvelope{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var rep *api.SecureEnvelope
	if rep, err = peer.Transfer(ctx, req); err != nil {
		return rpcerr(err)
	}

	return printJSON(rep)
}

func exchange(c *cli.Context) (err error) {
	req := &api.SigningKey{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var rep *api.SigningKey
	if rep, err = peer.KeyExchange(ctx, req); err != nil {
		return rpcerr(err)
	}

	return printJSON(rep)
}

func confirm(c *cli.Context) (err error) {
	return cli.Exit("unimplemented: the address confirmation protocol has not been fully specified by the TRISA working group", 9)
}

func health(c *cli.Context) (err error) {
	// Performs a status check with an empty health check request.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var rep *api.ServiceState
	if rep, err = ping.Status(ctx, &api.HealthCheck{}); err != nil {
		return rpcerr(err)
	}
	return printJSON(rep)
}

func lookup(c *cli.Context) (err error) {
	req := &gds.LookupRequest{
		Id:                  c.String("vasp-id"),
		RegisteredDirectory: c.String("registered-directory"),
		CommonName:          c.String("common-name"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var rep *gds.LookupReply
	if rep, err = directory.Lookup(ctx, req); err != nil {
		return rpcerr(err)
	}

	return printJSON(rep)
}

func search(c *cli.Context) (err error) {
	req := &gds.SearchRequest{
		Name:             c.StringSlice("name"),
		Website:          c.StringSlice("website"),
		Country:          c.StringSlice("country"),
		BusinessCategory: make([]models.BusinessCategory, 0, len(c.StringSlice("category"))),
		VaspCategory:     make([]string, 0, len(c.StringSlice("category"))),
	}

	for _, cat := range c.StringSlice("category") {
		if enum, ok := models.BusinessCategory_value[cat]; ok {
			req.BusinessCategory = append(req.BusinessCategory, models.BusinessCategory(enum))
		} else {
			req.VaspCategory = append(req.VaspCategory, cat)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var rep *gds.SearchReply
	if rep, err = directory.Search(ctx, req); err != nil {
		return rpcerr(err)
	}

	return printJSON(rep)
}

//====================================================================================
// Helper Commands
//====================================================================================

func initClient(c *cli.Context) (err error) {
	var endpoint string
	if endpoint = c.String("endpoint"); endpoint == "" {
		return cli.Exit("specify endpoint of TRISA peer to connect to", 1)
	}

	var creds grpc.DialOption
	if creds, err = loadCreds(endpoint, c); err != nil {
		return err
	}

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(endpoint, creds); err != nil {
		return cli.Exit(err, 1)
	}

	peer = api.NewTRISANetworkClient(cc)
	return nil
}

func initHealthCheckClient(c *cli.Context) (err error) {
	var endpoint string
	if endpoint = c.String("endpoint"); endpoint == "" {
		return cli.Exit("specify endpoint of TRISA peer to connect to", 1)
	}

	var opts []grpc.DialOption
	if c.Bool("insecure") {
		fmt.Println("warning: connecting in insecure mode is not supported by all TRISA peers")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		var creds grpc.DialOption
		if creds, err = loadCreds(endpoint, c); err != nil {
			return err
		}
		opts = append(opts, creds)
	}

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(endpoint, opts...); err != nil {
		return cli.Exit(err, 1)
	}

	ping = api.NewTRISAHealthClient(cc)
	return nil
}

func initDirectoryClient(c *cli.Context) (err error) {
	endpoint := c.String("directory")
	if _, ok := directoryAliases[endpoint]; ok {
		endpoint = directoryAliases[endpoint]
	}

	if endpoint == "" {
		return cli.Exit("specify endpoint or name of directory service to connect to", 1)
	}

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))); err != nil {
		return cli.Exit(err, 1)
	}

	directory = gds.NewTRISADirectoryClient(cc)
	return nil
}

func loadCreds(endpoint string, c *cli.Context) (creds grpc.DialOption, err error) {
	var (
		certs *trust.Provider
		pool  trust.ProviderPool
	)

	if certs, pool, err = loadCerts(c); err != nil {
		return nil, cli.Exit(err, 1)
	}

	if creds, err = mtls.ClientCreds(endpoint, certs, pool); err != nil {
		return nil, cli.Exit(err, 1)
	}
	return creds, nil
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
			return nil, nil, cli.Exit(err, 1)
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

func printJSON(msg interface{}) (err error) {
	var data []byte
	switch m := msg.(type) {
	case proto.Message:
		opts := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			AllowPartial:    true,
			UseProtoNames:   true,
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
		}

		if data, err = opts.Marshal(m); err != nil {
			return cli.Exit(err, 1)
		}
	default:
		if data, err = json.MarshalIndent(msg, "", "  "); err != nil {
			return cli.Exit(err, 1)
		}
	}

	fmt.Println(string(data))
	return nil
}

func rpcerr(err error) error {
	if serr, ok := status.FromError(err); ok {
		return cli.Exit(fmt.Errorf("%s: %s", serr.Code().String(), serr.Message()), int(serr.Code()))
	}
	return cli.Exit(err, 2)
}
