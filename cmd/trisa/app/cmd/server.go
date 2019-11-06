package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/trisacrypto/trisa/pkg/trisa/config"
	"github.com/trisacrypto/trisa/pkg/trisa/handler"
	"github.com/trisacrypto/trisa/pkg/trisa/server"
	"github.com/trisacrypto/trisa/pkg/trisa/trust"
	bitcoin "github.com/trisacrypto/trisa/proto/trisa/data/bitcoin/v1alpha1"
	us "github.com/trisacrypto/trisa/proto/trisa/identity/us/v1alpha1"
	pb "github.com/trisacrypto/trisa/proto/trisa/protocol/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start VASP TRISA server",
		Run:   runServerCmd,
	}

	return cmd
}

func runServerCmd(cmd *cobra.Command, args []string) {

	c, err := config.FromFile(configFile)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	chain, err := ioutil.ReadFile(c.TLS.TrustChainFile)
	if err != nil {
		log.Fatalf("load trust chain: %v", err)
	}
	tp := trust.NewProvider(chain)

	crt, err := tls.LoadX509KeyPair(c.TLS.CertificateFile, c.TLS.PrivateKeyFile)
	if err != nil {
		log.Fatalf("load x509 key pair: %v", err)
	}

	baseTLSCfg := &tls.Config{
		Certificates: []tls.Certificate{crt},
		MinVersion:   tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	handler := handler.NewDemoHandler()
	pServer := server.New(handler, crt, tp.GetCertPool())

	errs := make(chan error, 2)

	go func() {

		r := mux.NewRouter()

		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			out, _ := json.Marshal(struct {
				Hello string
			}{
				Hello: "World",
			})
			w.Write(out)
		})

		r.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
			out, _ := json.Marshal(mTLSConnectionTest(
				ctx,
				r.URL.Query(),
				crt,
				tp.GetCertPool(),
			))
			w.Write(out)
		})

		r.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {

			identity, _ := ptypes.MarshalAny(&us.Identity{
				FirstName:     "John",
				LastName:      "Doe",
				Ssn:           "001-0434-4983",
				DriverLicense: "FA-387463",
				State:         "CA",
			})

			data, _ := ptypes.MarshalAny(&bitcoin.Data{
				Source:      uuid.New().String(),
				Destination: uuid.New().String(),
			})

			tData := &pb.TransactionData{
				Identity: identity,
				Data:     data,
			}

			if err := pServer.SendRequest(r.Context(), r.URL.Query().Get("target"), uuid.New().String(), tData); err != nil {
				fmt.Fprintf(w, "error: %v", err)
				return
			}

			fmt.Fprint(w, ".")
		})

		srv := &http.Server{
			Addr:      c.Server.ListenAddressAdmin,
			Handler:   r,
			TLSConfig: baseTLSCfg,
		}

		log.WithFields(log.Fields{
			"component": "admin",
			"tls":       "listening",
			"port":      c.Server.ListenAddressAdmin,
		}).Info("starting TRISA admin server")

		errs <- srv.ListenAndServeTLS(c.TLS.CertificateFile, c.TLS.PrivateKeyFile)
	}()

	/*go func() {

		r := mux.NewRouter()

		r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

			w.Header().Add("Content-Type", "application/json")
			params := r.URL.Query()
			response := &PingResponse{
				Message:  params.Get("msg"),
				ClientCN: r.TLS.PeerCertificates[0].Subject.CommonName,
				ServerCN: r.TLS.ServerName,
			}
			out, _ := json.Marshal(response)
			w.Write(out)
		})

		mTLSCfg := &tls.Config{}
		copier.Copy(&mTLSCfg, &baseTLSCfg)

		mTLSCfg.ClientAuth = tls.RequireAndVerifyClientCert
		mTLSCfg.ClientCAs = tp.GetCertPool()

		srv := &http.Server{
			Addr:      c.Server.ListenAddress,
			Handler:   r,
			TLSConfig: mTLSCfg,
		}

		log.WithFields(log.Fields{
			"component": "service",
			"tls":       "listening",
			"port":      c.Server.ListenAddress,
		}).Info("starting TRISA server")

		errs <- srv.ListenAndServeTLS(c.TLS.CertificateFile, c.TLS.PrivateKeyFile)

	}()*/

	go func() {
		lis, err := net.Listen("tcp", c.Server.ListenAddress)
		if err != nil {
			errs <- err
		}

		mTLSCfg := &tls.Config{}
		copier.Copy(&mTLSCfg, &baseTLSCfg)

		mTLSCfg.ClientAuth = tls.RequireAndVerifyClientCert
		mTLSCfg.ClientCAs = tp.GetCertPool()

		tc := credentials.NewTLS(mTLSCfg)
		s := grpc.NewServer(grpc.Creds(tc))
		pb.RegisterTrisaPeer2PeerServer(s, pServer)

		log.WithFields(log.Fields{
			"component": "grpc",
			"tls":       "listening",
			"port":      c.Server.ListenAddress,
		}).Info("starting TRISA server")

		errs <- s.Serve(lis)
	}()

	log.Fatalf("terminated %v", <-errs)
}

type PingResponse struct {
	Status   string `json:"status,omitempty"`
	Message  string `json:"message,omitempty"`
	ServerCN string `json:"server_cn,omitempty"`
	ClientCN string `json:"client_cn,omitempty"`
}

func mTLSConnectionTest(ctx context.Context, params url.Values, crt tls.Certificate, certPool *x509.CertPool) *PingResponse {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{crt},
				RootCAs:      certPool,
			},
		},
	}

	url := fmt.Sprintf("https://%s/ping?msg=%s", params.Get("target"), params.Get("msg"))

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return responseFailure(err)
	}

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return responseFailure(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return responseFailure(err)
	}

	out := &PingResponse{}
	json.Unmarshal(body, out)

	return out
}

func responseFailure(err error) *PingResponse {
	return &PingResponse{
		Message: "something went wrong",
		Status:  err.Error(),
	}
}
