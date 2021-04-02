package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mmcloughlin/addchain/internal/metavars"
	"github.com/mmcloughlin/addchain/internal/zenodo"
)

type VarsFile struct {
	path string
}

func (v *VarsFile) SetFlags(f *flag.FlagSet) {
	f.StringVar(&v.path, "vars", v.DefaultPath(), "path to meta variables file")
}

// DefaultPath returns the path to the default variables file in the
// meta package.
func (v *VarsFile) DefaultPath() string {
	_, self, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(self), "../../meta/vars.go")
	return filepath.Clean(path)
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

type Zenodo struct {
	base  string
	token string
}

const zenodoTokenEnvVar = "ZENODO_TOKEN"

func (z *Zenodo) SetFlags(f *flag.FlagSet) {
	f.StringVar(&z.base, "url", zenodo.BaseURL, "zenodo api base url")
	f.StringVar(&z.token, "token", "", fmt.Sprintf("zenodo token (uses %q environment variable if empty)", zenodoTokenEnvVar))
}

func (z *Zenodo) Token() (string, error) {
	if z.token != "" {
		return z.token, nil
	}
	if token, ok := os.LookupEnv(zenodoTokenEnvVar); ok {
		return token, nil
	}
	return "", errors.New("zenodo token not specified")
}

func (z *Zenodo) Client() (*zenodo.Client, error) {
	// HTTP Client.
	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{}
	client.Transport = tr

	// Token.
	token, err := z.Token()
	if err != nil {
		return nil, err
	}

	// Zenodo client.
	c := zenodo.NewClient(client, z.base, token)

	return c, nil
}
