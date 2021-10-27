package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mmcloughlin/addchain/internal/metavars"
	"github.com/mmcloughlin/addchain/internal/zenodo"
)

// VarsFile represents a meta variables file.
type VarsFile struct {
	path string
}

// SetFlags registers command-line flags to configure the meta variables file.
func (v *VarsFile) SetFlags(f *flag.FlagSet) {
	f.StringVar(&v.path, "vars", v.DefaultPath(), "path to meta variables file")
}

// DefaultPath returns the path to the default variables file in the meta
// package. Returns the empty string if the path cannot be determined.
func (v *VarsFile) DefaultPath() string {
	return RepoPath("meta/vars.go")
}

// Get variable from the meta variables file.
func (v *VarsFile) Get(name string) (string, error) {
	f, err := metavars.ReadFile(v.path)
	if err != nil {
		return "", err
	}

	value, ok := f.Get(name)
	if !ok {
		return "", fmt.Errorf("unknown property %q", name)
	}

	return value, nil
}

// Set variable in the meta variables file.
func (v *VarsFile) Set(name, value string) error {
	f, err := metavars.ReadFile(v.path)
	if err != nil {
		return err
	}

	if err := f.Set(name, value); err != nil {
		return err
	}

	if err := metavars.WriteFile(v.path, f); err != nil {
		return err
	}

	return nil
}

// Zenodo configures a zenodo client.
type Zenodo struct {
	base    string
	sandbox bool
	token   string
}

const zenodoTokenEnvVar = "ZENODO_TOKEN"

// SetFlags registers command-line flags to configure the zenodo client.
func (z *Zenodo) SetFlags(f *flag.FlagSet) {
	f.StringVar(&z.base, "url", zenodo.BaseURL, "zenodo api base url")
	f.BoolVar(&z.sandbox, "sandbox", false, "use zenodo sandbox")
	f.StringVar(&z.token, "token", "", fmt.Sprintf("zenodo token (uses %q environment variable if empty)", zenodoTokenEnvVar))
}

// Token returns the configured zenodo token, either from the command-line or
// environment variable.
func (z *Zenodo) Token() (string, error) {
	if z.token != "" {
		return z.token, nil
	}
	if token, ok := os.LookupEnv(zenodoTokenEnvVar); ok {
		return token, nil
	}
	return "", errors.New("zenodo token not specified")
}

// URL to connect to.
func (z *Zenodo) URL() string {
	if z.sandbox {
		return zenodo.SandboxBaseURL
	}
	return z.base
}

// Client builds the configured client.
func (z *Zenodo) Client(c *http.Client) (*zenodo.Client, error) {
	// Token.
	token, err := z.Token()
	if err != nil {
		return nil, err
	}

	// Zenodo client.
	return zenodo.NewClient(c, z.URL(), token), nil
}

// HTTPClient configures a HTTP client.
type HTTPClient struct {
	cert string
}

// SetFlags registers command-line flags to configure the HTTP client.
func (h *HTTPClient) SetFlags(f *flag.FlagSet) {
	f.StringVar(&h.cert, "cert", "", "trust additional certificate authority certificate")
}

// Client builds a HTTP client according to configuration.
func (h *HTTPClient) Client() (*http.Client, error) {
	// CAs.
	roots, err := h.RootCAs()
	if err != nil {
		return nil, err
	}

	// Transport.
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    roots,
		},
	}

	// Client.
	return &http.Client{
		Transport: tr,
	}, nil
}

// RootCAs returns the configured certificate pool.
func (h *HTTPClient) RootCAs() (*x509.CertPool, error) {
	roots, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	if h.cert == "" {
		return roots, nil
	}

	data, err := ioutil.ReadFile(h.cert)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, err
	}

	roots.AddCert(cert)

	return roots, nil
}
