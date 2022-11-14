/*
The TRISA CLI client allows you to create and execute TRISA requests from the command
line for development or testing purposes. For more information on how to use the CLI,
run `trisa --help` or see the documenation at https://trisa.dev.
*/
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/trisacrypto/trisa/pkg"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	apierr "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1/errors"
	generic "github.com/trisacrypto/trisa/pkg/trisa/data/generic/v1beta1"
	env "github.com/trisacrypto/trisa/pkg/trisa/envelope"
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
	"google.golang.org/protobuf/types/known/anypb"
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
			Usage:   "create a secure envelope or payload template from a payload",
			Action:  envelope,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "identity",
					Aliases: []string{"i"},
					Usage:   "path to identity payload JSON to load",
				},
				&cli.StringFlag{
					Name:    "transaction",
					Aliases: []string{"t"},
					Usage:   "path to transaction payload JSON to load",
				},
				&cli.StringFlag{
					Name:        "out",
					Aliases:     []string{"o"},
					Usage:       "path to save sealed secure envelope to disk",
					DefaultText: "stdout",
				},
				&cli.StringFlag{
					Name:    "sealing-key",
					Aliases: []string{"s", "seal"},
					Usage:   "path to recipient's public key to seal outgoing envelope with (optional)",
				},
				&cli.StringFlag{
					Name:    "envelope-id",
					Aliases: []string{"id", "I"},
					Usage:   "specify the envelope ID for the outgoing secure envelope (optional)",
				},
				&cli.StringFlag{
					Name:        "sent-at",
					Aliases:     []string{"sent", "S"},
					Usage:       "specify a sent at timestamp for the payload in RFC3339 format",
					DefaultText: "now",
				},
				&cli.StringFlag{
					Name:    "received-at",
					Aliases: []string{"received", "R"},
					Usage:   "specify a received at timestamp for the payload in RFC3339 format or the keyword \"now\" (optional)",
				},
				&cli.StringFlag{
					Name:    "error-code",
					Aliases: []string{"C"},
					Usage:   "add an error with the specified code to the outgoing envelope",
					Value:   "REJECTED",
				},
				&cli.StringFlag{
					Name:    "error-message",
					Aliases: []string{"error", "E"},
					Usage:   "add an error message to the outgoing envelope",
				},
			},
		},
		{
			Name:   "seal",
			Usage:  "seal a secure envelope from a payload template",
			Action: seal,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "in",
					Aliases: []string{"i"},
					Usage:   "path to secure envelope or payload template to load",
				},
				&cli.StringFlag{
					Name:        "out",
					Aliases:     []string{"o"},
					Usage:       "path to save sealed secure envelope to disk",
					DefaultText: "stdout",
				},
				&cli.StringFlag{
					Name:    "sealing-key",
					Aliases: []string{"s", "seal"},
					Usage:   "path to recipient's public key to seal outgoing envelope with",
				},
				&cli.StringFlag{
					Name:        "envelope-id",
					Aliases:     []string{"id", "I"},
					Usage:       "specify the envelope ID for the outgoing secure envelope",
					DefaultText: "original or random UUID",
				},
				&cli.StringFlag{
					Name:        "sent-at",
					Aliases:     []string{"sent", "S"},
					Usage:       "specify a sent at timestamp for the payload in RFC3339 format",
					DefaultText: "original or now",
				},
				&cli.StringFlag{
					Name:        "received-at",
					Aliases:     []string{"received", "R"},
					Usage:       "specify a received at timestamp for the payload in RFC3339 format or the keyword \"now\"",
					DefaultText: "original or empty",
				},
				&cli.StringFlag{
					Name:    "error-code",
					Aliases: []string{"C"},
					Usage:   "add an error with the specified code to the outgoing envelope",
					Value:   "REJECTED",
				},
				&cli.StringFlag{
					Name:    "error-message",
					Aliases: []string{"error", "E"},
					Usage:   "add an error message to the outgoing envelope",
				},
			},
		},
		{
			Name:   "open",
			Usage:  "unseal a secure envelope from an envelope saved to disk",
			Action: open,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "in",
					Aliases:  []string{"i"},
					Usage:    "path to secure envelope to open and unseal or parse",
					Required: true,
				},
				&cli.StringFlag{
					Name:        "out",
					Aliases:     []string{"o"},
					Usage:       "path to save unsealed envelope to disk",
					DefaultText: "stdout",
				},
				&cli.BoolFlag{
					Name:    "payload",
					Aliases: []string{"p"},
					Usage:   "extract payload when saving to disk (if -out is specified)",
				},
				&cli.BoolFlag{
					Name:    "error",
					Aliases: []string{"E"},
					Usage:   "extract error from secure envelope",
				},
				&cli.StringFlag{
					Name:    "unsealing-key",
					Aliases: []string{"key", "k"},
					Usage:   "path to private key to unseal the secure envelope",
				},
			},
		},
		{
			Name:      "transfer",
			Usage:     "execute a TRISA transfer with a TRISA peer",
			UsageText: "trisa transfer -i sealed_envelope.json\ntrisa transfer -i payload.json -s public.pem\ntrisa transfer -I [envelope-id] -E [error message]",
			Before:    initClient,
			Action:    transfer,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "in",
					Aliases: []string{"i"},
					Usage:   "path to secure envelope or payload template to load",
				},
				&cli.StringFlag{
					Name:        "out",
					Aliases:     []string{"o"},
					Usage:       "path to save response envelope to disk",
					DefaultText: "stdout",
				},
				&cli.StringFlag{
					Name:    "unsealing-key",
					Aliases: []string{"key", "k"},
					Usage:   "path to private key to unseal the incoming secure envelope",
				},
				&cli.BoolFlag{
					Name:    "payload",
					Aliases: []string{"p"},
					Usage:   "extract payload when saving to disk (if -out and -key are specified)",
				},
				&cli.StringFlag{
					Name:    "sealing-key",
					Aliases: []string{"s", "seal"},
					Usage:   "path to recipient's public key to seal outgoing envelope with",
				},
				&cli.StringFlag{
					Name:        "envelope-id",
					Aliases:     []string{"id", "I"},
					Usage:       "specify the envelope ID for the outgoing secure envelope",
					DefaultText: "original or random UUID",
				},
				&cli.StringFlag{
					Name:        "sent-at",
					Aliases:     []string{"sent", "S"},
					Usage:       "specify a sent at timestamp for the payload in RFC3339 format",
					DefaultText: "original or now",
				},
				&cli.StringFlag{
					Name:        "received-at",
					Aliases:     []string{"recv", "R"},
					Usage:       "specify a received at timestamp for the payload in RFC3339 format or the keyword \"now\"",
					DefaultText: "original",
				},
				&cli.StringFlag{
					Name:    "error-code",
					Aliases: []string{"C"},
					Usage:   "add an error with the specified code to the outgoing envelope",
					Value:   "REJECTED",
				},
				&cli.StringFlag{
					Name:    "error-message",
					Aliases: []string{"error", "E"},
					Usage:   "add an error message to the outgoing envelope",
				},
			},
		},
		{
			Name:    "exchange",
			Aliases: []string{"key-exchange"},
			Usage:   "exchange public sealing keys with a TRISA peer",
			Before:  initClient,
			Action:  exchange,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "in",
					Aliases:     []string{"i"},
					Usage:       "path to PEM encoded public key to send to remote",
					DefaultText: "TRISA certs",
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "path to write PEM encoded keys sent from remote",
				},
			},
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
// Envelope Handling Commands
//====================================================================================

func envelope(c *cli.Context) (err error) {
	// Is this an error only envelope?
	if c.String("error-message") != "" {
		if c.String("envelope-id") == "" {
			return cli.Exit("an envelope id is required to create an error only secure envelope", 1)
		}

		if c.String("identity") != "" || c.String("transaction") != "" {
			return cli.Exit("this command does not create secure envelopes with both an error and a payload, specify either -error or -transaction and -identity", 1)
		}

		var msg *api.SecureEnvelope
		if msg, err = errorMessage(c.String("envelope-id"), c.String("error-message"), c.String("error-code")); err != nil {
			return cli.Exit(err, 1)
		}

		if out := c.String("out"); out != "" {
			if err = dumpProto(msg, out); err != nil {
				return cli.Exit(err, 1)
			}
			return nil
		}
		return printJSON(msg)
	}

	// Creating an envelope with a payload, identity and transaction are required
	if c.String("identity") == "" || c.String("transaction") == "" {
		return cli.Exit("a path to both the identity and transaction JSON is required", 1)
	}

	var sealingKey interface{}
	if path := c.String("sealing-key"); path != "" {
		if sealingKey, err = loadSealingKey(path); err != nil {
			return cli.Exit(err, 1)
		}
	}

	// Create the payload
	payload := &api.Payload{}
	ts := c.String("sent-at")
	if strings.ToLower(ts) == "now" || ts == "" {
		ts = time.Now().Format(time.RFC3339)
	}
	payload.SentAt = ts

	if ts := c.String("received-at"); ts != "" {
		if strings.ToLower(ts) == "now" {
			ts = time.Now().Format(time.RFC3339)
		}
		payload.ReceivedAt = ts
	}

	// Load the payloads from JSON
	if payload.Identity, err = loadIdentity(c.String("identity")); err != nil {
		return cli.Exit(err, 1)
	}

	if payload.Transaction, err = loadTransaction(c.String("transaction")); err != nil {
		return cli.Exit(err, 1)
	}

	envelopeID := c.String("envelope-id")
	if envelopeID == "" {
		envelopeID = uuid.NewString()
	}

	// Create the envelope
	var handler *env.Envelope
	if handler, err = env.New(payload, env.WithEnvelopeID(envelopeID)); err != nil {
		return cli.Exit(err, 1)
	}

	// Create the unsealed envelope
	if handler, _, err = handler.Encrypt(); err != nil {
		return cli.Exit(err, 1)
	}

	// Create the unsealed envelope if necessary
	if sealingKey != nil {
		if handler, _, err = handler.Seal(env.WithSealingKey(sealingKey)); err != nil {
			return cli.Exit(err, 1)
		}
	}

	// Save to disk
	if out := c.String("out"); out != "" {
		if err = dumpProto(handler.Proto(), out); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}
	return printJSON(handler.Proto())
}

func seal(c *cli.Context) (err error) {
	var msg *api.SecureEnvelope
	if in := c.String("in"); in != "" {
		// Load the envelope or payload template
		if msg, err = loadEnvelope(in); err != nil {
			return cli.Exit(err, 1)
		}

		// If the sealing key is provided, use it to seal the envelope
		var sealKey interface{}
		if path := c.String("sealing-key"); path != "" {
			if sealKey, err = loadSealingKey(path); err != nil {
				return cli.Exit(err, 1)
			}
		} else {
			return cli.Exit("a path to the public sealing keys is required to seal envelope", 1)
		}

		if msg, err = updateEnvelope(msg, c); err != nil {
			return cli.Exit(err, 1)
		}

		if msg, err = sealEnvelope(msg, sealKey); err != nil {
			return cli.Exit(err, 1)
		}
	} else {
		// Attempt to create an error-only envelope
		if c.String("error-message") == "" || c.String("envelope-id") == "" {
			return cli.Exit("specify an envelope to load or an error message and id", 1)
		}

		if msg, err = errorMessage(c.String("envelope-id"), c.String("error-message"), c.String("error-code")); err != nil {
			return cli.Exit(err, 1)
		}
	}

	// Always use the current timestamp for ordering purposes
	msg.Timestamp = time.Now().Format(time.RFC3339Nano)

	// Did we manage to load a valid secure envelope?
	if err = env.Validate(msg); err != nil {
		return cli.Exit(fmt.Errorf("could not load envelope to send: %s", err), 1)
	}

	if out := c.String("out"); out != "" {
		if err = dumpProto(msg, out); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}
	return printJSON(msg)
}

func open(c *cli.Context) (err error) {
	var msg *api.SecureEnvelope
	if msg, err = loadEnvelope(c.String("in")); err != nil {
		return cli.Exit(err, 1)
	}

	// Extract the error if requested - no cryptography required
	if c.Bool("error") {
		if msg.Error == nil || apierr.IsZero(msg.Error) {
			return cli.Exit("there is no error on the secure envelope", 1)
		}

		if out := c.String("out"); out != "" {
			if err = dumpProto(msg.Error, out); err != nil {
				return cli.Exit(err, 1)
			}
		}
		return printJSON(msg.Error)
	}

	var handler *env.Envelope
	if handler, err = env.Wrap(msg); err != nil {
		return cli.Exit(err, 1)
	}

	// Has the unsealing key been provided?
	var unsealingKey interface{}
	if path := c.String("unsealing-key"); path != "" {
		if unsealingKey, err = loadPrivateKey(path); err != nil {
			return cli.Exit(err, 1)
		}
	}

	var (
		payload          *api.Payload
		unsealedEnvelope *api.SecureEnvelope
	)

	// Figure out if we can unseal the envelope
	switch handler.State() {
	case env.Sealed, env.SealedError:
		if unsealingKey == nil {
			return cli.Exit("must specify unsealing key to open sealed envelope", 1)
		}

		if handler, _, err = handler.Unseal(env.WithUnsealingKey(unsealingKey)); err != nil {
			return cli.Exit(err, 1)
		}

		unsealedEnvelope = handler.Proto()
		if handler, _, err = handler.Decrypt(); err != nil {
			return cli.Exit(err, 1)
		}

		payload, _ = handler.Payload()

	case env.Unsealed, env.UnsealedError:
		unsealedEnvelope = handler.Proto()
		if handler, _, err = handler.Decrypt(); err != nil {
			return cli.Exit(err, 1)
		}
		payload, _ = handler.Payload()

	case env.Clear, env.ClearError:
		payload, _ = handler.Payload()
	default:
		return cli.Exit(fmt.Errorf("envelope in unhandled state %s", handler.State()), 1)
	}

	if out := c.String("out"); out != "" {
		if c.Bool("payload") {
			if err = dumpProto(payload, out); err != nil {
				return cli.Exit(err, 1)
			}
			return nil
		}

		if err = dumpProto(unsealedEnvelope, out); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}

	if c.Bool("payload") {
		return printJSON(payload)
	}
	return printJSON(unsealedEnvelope)
}

//====================================================================================
// TRISA RPC Commands
//====================================================================================

func transfer(c *cli.Context) (err error) {
	// There are three cases for loading a secure envelope to send:
	// 1. a sealed secure envelope is unmarshaled from disk
	// 2. an error and envelope ID is specified on the command line
	// 3. a payload template is unmarshaled from disk with a sealing key
	var req *api.SecureEnvelope
	if in := c.String("in"); in != "" {
		// Load the envelope or payload template
		if req, err = loadEnvelope(in); err != nil {
			return cli.Exit(err, 1)
		}

		// If the sealing key is provided, use it to seal the envelope
		if seal := c.String("sealing-key"); seal != "" {
			var sealKey interface{}
			if sealKey, err = loadSealingKey(seal); err != nil {
				return cli.Exit(err, 1)
			}

			if req, err = updateEnvelope(req, c); err != nil {
				return cli.Exit(err, 1)
			}

			if req, err = sealEnvelope(req, sealKey); err != nil {
				return cli.Exit(err, 1)
			}
		}
	} else {
		// Attempt to create an error-only envelope
		if c.String("error-message") == "" || c.String("envelope-id") == "" {
			return cli.Exit("specify an envelope to load or an error message and id", 1)
		}

		if req, err = errorMessage(c.String("envelope-id"), c.String("error-message"), c.String("error-code")); err != nil {
			return cli.Exit(err, 1)
		}
	}

	// Always use the current timestamp for ordering purposes
	req.Timestamp = time.Now().Format(time.RFC3339Nano)

	// Did we manage to load a valid secure envelope?
	if err = env.Validate(req); err != nil {
		return cli.Exit(fmt.Errorf("could not load envelope to send: %s", err), 1)
	}

	// Is the envelope sealed?
	if !req.Sealed && (req.Error == nil || apierr.IsZero(req.Error)) {
		return cli.Exit("envelope has not been sealed", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var rep *api.SecureEnvelope
	if rep, err = peer.Transfer(ctx, req); err != nil {
		return rpcerr(err)
	}

	// If a private key is provided, unseal the envelope
	if unseal := c.String("unsealing-key"); unseal != "" {
		var unsealKey interface{}
		if unsealKey, err = loadPrivateKey(unseal); err != nil {
			return cli.Exit(err, 1)
		}

		var handler *env.Envelope
		if handler, err = env.Wrap(rep, env.WithUnsealingKey(unsealKey)); err != nil {
			return cli.Exit(err, 1)
		}

		if handler, _, err = handler.Unseal(); err != nil {
			return cli.Exit(err, 1)
		}

		// Are we saving to disk or printing JSON?
		if out := c.String("out"); out != "" {
			if extractPayload := c.Bool("payload"); extractPayload {
				if handler, _, err = handler.Decrypt(); err != nil {
					return cli.Exit(err, 1)
				}
				payload, _ := handler.Payload()
				if err = dumpProto(payload, out); err != nil {
					return cli.Exit(err, 1)
				}
				return nil
			}

			if err = dumpProto(handler.Proto(), out); err != nil {
				return cli.Exit(err, 1)
			}
			return nil
		}

		// Print the decrypted envelope
		if handler, _, err = handler.Decrypt(); err != nil {
			return cli.Exit(err, 1)
		}
		payload, _ := handler.Payload()
		return printJSON(payload)
	}

	if out := c.String("out"); out != "" {
		if err = dumpProto(rep, out); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}
	return printJSON(rep)
}

func exchange(c *cli.Context) (err error) {
	var req *api.SigningKey
	if in := c.String("in"); in != "" {
		if req, err = loadPublicKeys(in); err != nil {
			return cli.Exit(err, 1)
		}
	} else {
		// By default use the TRISA identity certificates in the key exchange
		var provider *trust.Provider
		if provider, _, err = loadCerts(c); err != nil {
			return err
		}

		var certs *x509.Certificate
		if certs, err = provider.GetLeafCertificate(); err != nil {
			return cli.Exit(err, 1)
		}

		req = &api.SigningKey{}
		req.Version = int64(certs.Version)
		req.Signature = certs.Signature
		req.SignatureAlgorithm = certs.SignatureAlgorithm.String()
		req.PublicKeyAlgorithm = certs.PublicKeyAlgorithm.String()
		req.NotBefore = certs.NotBefore.Format(time.RFC3339)
		req.NotAfter = certs.NotAfter.Format(time.RFC3339)

		if req.Data, err = x509.MarshalPKIXPublicKey(certs.PublicKey); err != nil {
			return cli.Exit("could not create public sealing key from certs", 1)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var rep *api.SigningKey
	if rep, err = peer.KeyExchange(ctx, req); err != nil {
		return rpcerr(err)
	}

	if out := c.String("out"); out != "" {
		if err = dumpKeys(rep, out); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
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

//====================================================================================
// Directory RPC Commands
//====================================================================================

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
// Helper Commands - Clients
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

//====================================================================================
// Helper Commands - Serialization and Deserialization
//====================================================================================

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

func loadPublicKeys(path string) (key *api.SigningKey, err error) {
	key = new(api.SigningKey)

	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, fmt.Errorf("could not read key file: %s", err)
	}

	switch filepath.Ext(path) {
	case ".json":
		opts := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}
		if err = opts.Unmarshal(data, key); err != nil {
			return nil, fmt.Errorf("could not unmarshal json keys: %s", err)
		}
	case ".pb":
		if err = proto.Unmarshal(data, key); err != nil {
			return nil, fmt.Errorf("could not unmarshal pb keys: %s", err)
		}
	case ".pem", ".crt":
		var block *pem.Block
	pemblocks:
		for {
			block, data = pem.Decode(data)
			if block == nil {
				return nil, fmt.Errorf("could not find public key in key file %s", path)
			}

			switch block.Type {
			case trust.BlockPublicKey, trust.BlockRSAPublicKey:
				key.Data = block.Bytes
				break pemblocks
			case trust.BlockCertificate:
				// assumes that the first certificate in a trust chain is the "leaf"
				var certs *x509.Certificate
				if certs, err = x509.ParseCertificate(block.Bytes); err != nil {
					return nil, fmt.Errorf("could not parse certificate: %s", err)
				}
				key.Version = int64(certs.Version)
				key.Signature = certs.Signature
				key.SignatureAlgorithm = certs.SignatureAlgorithm.String()
				key.PublicKeyAlgorithm = certs.PublicKeyAlgorithm.String()
				key.NotBefore = certs.NotBefore.Format(time.RFC3339)
				key.NotAfter = certs.NotAfter.Format(time.RFC3339)

				if key.Data, err = x509.MarshalPKIXPublicKey(certs.PublicKey); err != nil {
					return nil, fmt.Errorf("could not parse certificate: %s", err)
				}
				break pemblocks
			}
		}
	default:
		return nil, fmt.Errorf("unhandled extension %q use .json or .pem", filepath.Ext(path))
	}

	return key, nil
}

func loadSealingKey(path string) (key interface{}, err error) {
	var sealingKey *api.SigningKey
	if sealingKey, err = loadPublicKeys(path); err != nil {
		return nil, err
	}

	if key, err = x509.ParsePKIXPublicKey(sealingKey.Data); err != nil {
		return nil, err
	}
	return key, nil
}

func loadPrivateKey(path string) (key interface{}, err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, fmt.Errorf("could not read key file: %s", err)
	}

	var block *pem.Block
	for {
		block, data = pem.Decode(data)
		if block == nil {
			break
		}

		switch block.Type {
		case trust.BlockPrivateKey, trust.BlockRSAPrivateKey, trust.BlockECPrivateKey:
			return trust.ParsePrivateKey(block)
		}
	}
	return nil, fmt.Errorf("could not find private key in key file %s", path)
}

func dumpKeys(key *api.SigningKey, path string) (err error) {
	var data []byte
	switch filepath.Ext(path) {
	case ".json":
		opts := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			AllowPartial:    true,
			UseProtoNames:   true,
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
		}

		if data, err = opts.Marshal(key); err != nil {
			return cli.Exit(err, 1)
		}
	case ".pb":
		if data, err = proto.Marshal(key); err != nil {
			return cli.Exit(err, 1)
		}
	case ".pem":
		var out interface{}
		if out, err = x509.ParsePKIXPublicKey(key.Data); err != nil {
			return fmt.Errorf("invalid PKIX public key received from remote")
		}
		if data, err = trust.PEMEncodePublicKey(out); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown extension %q use .json or .pem", filepath.Ext(path))
	}

	if err = os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write keys to disk: %s", err)
	}
	fmt.Printf("saved keys to %s\n", path)
	return nil
}

func loadIdentity(path string) (_ *anypb.Any, err error) {
	opts := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: false,
	}

	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, fmt.Errorf("could not read identity from %s", err)
	}

	// Attempt to unmarshal the data as IVMS 101
	identity := &ivms101.IdentityPayload{}
	if err = opts.Unmarshal(data, identity); err == nil {
		return anypb.New(identity)
	}

	// Attempt to unmarshal the data as a serialized any
	msg := &anypb.Any{}
	if err = opts.Unmarshal(data, msg); err == nil {
		return msg, nil
	}
	return nil, fmt.Errorf("could not unmarshal identity: unknown type or format")
}

func loadTransaction(path string) (_ *anypb.Any, err error) {
	opts := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: false,
	}

	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, fmt.Errorf("could not read transaction from %s", err)
	}

	// Attempt to unmarshal the data as a generic Transaction
	transaction := &generic.Transaction{}
	if err = opts.Unmarshal(data, transaction); err == nil {
		return anypb.New(transaction)
	}

	// Attempt to unmarshal the data as a generic Pending
	pending := &generic.Pending{}
	if err = opts.Unmarshal(data, pending); err == nil {
		return anypb.New(pending)
	}

	// Attempt to unmarshal the data as a serialized any
	msg := &anypb.Any{}
	if err = opts.Unmarshal(data, msg); err == nil {
		return msg, nil
	}

	return nil, fmt.Errorf("could not unmarshal transaction: unknown type or format")
}

func loadEnvelope(path string) (msg *api.SecureEnvelope, err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, fmt.Errorf("could not read envelope file: %s", err)
	}

	msg = &api.SecureEnvelope{}
	switch filepath.Ext(path) {
	case ".json":
		opts := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}
		if err = opts.Unmarshal(data, msg); err != nil {
			return nil, fmt.Errorf("could not unmarshal json envelope: %s", err)
		}
	case ".pb":
		if err = proto.Unmarshal(data, msg); err != nil {
			return nil, fmt.Errorf("could not unmarshal pb envelope: %s", err)
		}

	default:
		return nil, fmt.Errorf("unhandled extension %q use .json or .pb", filepath.Ext(path))
	}
	return msg, nil
}

func dumpProto(msg proto.Message, path string) (err error) {
	var data []byte
	switch filepath.Ext(path) {
	case ".json":
		opts := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			AllowPartial:    true,
			UseProtoNames:   true,
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
		}

		if data, err = opts.Marshal(msg); err != nil {
			return cli.Exit(err, 1)
		}
	case ".pb":
		if data, err = proto.Marshal(msg); err != nil {
			return cli.Exit(err, 1)
		}
	default:
		return fmt.Errorf("unknown extension %q use .json or .pb", filepath.Ext(path))
	}
	if err = os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write message to disk: %s", err)
	}
	fmt.Printf("message saved to %s\n", path)
	return nil
}

//====================================================================================
// Helper Commands - Managing Secure Envelopes
//====================================================================================

func errorMessage(id, message, code string) (_ *api.SecureEnvelope, err error) {
	code = strings.ToUpper(code)
	val, ok := api.Error_Code_value[code]
	if !ok {
		return nil, fmt.Errorf("%q is not a valid error code", code)
	}

	return &api.SecureEnvelope{
		Id: id,
		Error: &api.Error{
			Code:    api.Error_Code(val),
			Message: message,
		},
	}, nil
}

func updateEnvelope(in *api.SecureEnvelope, c *cli.Context) (out *api.SecureEnvelope, err error) {
	// Validation: cannot update error message or update the payload if the envelope is sealed
	if emsg := c.String("error-message"); emsg != "" {
		return nil, errors.New("error message is ignored when -in is supplied")
	}

	if in.Sealed && (c.String("sent-at") != "" || c.String("received-at") != "") {
		return nil, errors.New("cannot update payload on a sealed envelope")
	}

	// Set the envelope id if it is supplied on the command line; otherwise if there
	// is no envelope id then set it to a new random id.
	if eid := c.String("envelope-id"); eid != "" {
		in.Id = eid
	}

	if in.Id == "" {
		in.Id = uuid.NewString()
	}

	// If the envelope is sealed, return here.
	if in.Sealed {
		return in, nil
	}

	// Otherwise continue to update the payload
	var handler *env.Envelope
	if handler, err = env.Wrap(in); err != nil {
		return nil, err
	}

	// Decrypt and parse the secure envelope
	if handler, _, err = handler.Decrypt(); err != nil {
		return nil, err
	}

	// Fetch the payload to update it
	var payload *api.Payload
	if payload, err = handler.Payload(); err != nil {
		return nil, err
	}

	if ts := c.String("sent-at"); ts != "" {
		payload.SentAt = ts
	}

	if payload.SentAt == "" {
		payload.SentAt = time.Now().Format(time.RFC3339)
	}

	if ts := c.String("received-at"); ts != "" {
		if strings.ToLower(ts) == "now" {
			ts = time.Now().Format(time.RFC3339)
		}
		payload.ReceivedAt = ts
	}

	if handler, err = handler.Update(payload); err != nil {
		return nil, err
	}

	if handler, _, err = handler.Encrypt(); err != nil {
		return nil, err
	}

	return handler.Proto(), nil
}

func sealEnvelope(in *api.SecureEnvelope, key interface{}) (out *api.SecureEnvelope, err error) {
	var handler *env.Envelope
	if handler, err = env.Wrap(in, env.WithSealingKey(key)); err != nil {
		return nil, err
	}

	if handler, _, err = handler.Seal(); err != nil {
		return nil, err
	}
	return handler.Proto(), nil
}

//====================================================================================
// Helper Commands - CLI output
//====================================================================================

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
