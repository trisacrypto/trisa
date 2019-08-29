package discovery

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/trisacrypto/trisa/pkg/ca"

	"github.com/trisacrypto/trisa/proto/tvca/discovery"
)

const (
	WELL_KNOWN = "/.well-known/trisa"
	USER_AGENT = "Trisa Golang Client"
)

func New(ca string) (*Trisa, error) {
	u, err := url.Parse(ca)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, WELL_KNOWN)

	return &Trisa{
		wkURL: u.String(),
		hc:    &http.Client{},
	}, nil
}

type Trisa struct {
	wkURL     string
	hc        *http.Client
	disco     *discovery.Trisa
	RootCAs   []*x509.Certificate
	IssuerCAs []*x509.Certificate
}

func (t *Trisa) AddRootCA(crt *x509.Certificate) {
	t.RootCAs = append(t.RootCAs, crt)
}

func (t *Trisa) AddIssuerCA(crt *x509.Certificate) {
	t.IssuerCAs = append(t.IssuerCAs, crt)
}

func (t *Trisa) Init(ctx context.Context) error {
	req, err := http.NewRequest("GET", t.wkURL, nil)

	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", USER_AGENT)

	res, err := t.hc.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)

	t.disco = &discovery.Trisa{}
	if err != nil {
		return err
	}

	return json.Unmarshal(body, t.disco)
}

func (t *Trisa) LoadAll(ctx context.Context) error {

	rootStore, err := t.loadX509Store(ctx, t.disco.X509RootStore)
	if err != nil {
		return err
	}
	if err := t.loadStore(rootStore, t.AddRootCA); err != nil {
		return err
	}

	issuerStore, err := t.loadX509Store(ctx, t.disco.X509IssuerStore)
	if err != nil {
		return err
	}
	if err := t.loadStore(issuerStore, t.AddIssuerCA); err != nil {
		return err
	}

	return nil
}

func (t *Trisa) loadStore(store *discovery.X509Store, add func(*x509.Certificate)) error {
	for _, entry := range store.Store {
		crt, err := ca.PEMDecodeCertificate([]byte(entry.Pem))
		if err != nil {
			return err
		}
		add(crt)
	}
	return nil
}

func (t *Trisa) loadX509Store(ctx context.Context, url string) (*discovery.X509Store, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", USER_AGENT)

	res, err := t.hc.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	store := &discovery.X509Store{}
	if err := json.Unmarshal(body, store); err != nil {
		return nil, err
	}
	return store, nil
}
