package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/addchain/internal/cli"
	"github.com/mmcloughlin/addchain/internal/zenodo"
)

// upload subcommand.
type upload struct {
	cli.Command

	httpclient HTTPClient
	zenodo     Zenodo
	varsfile   VarsFile
	metadata   string
	publish    bool
}

func (*upload) Name() string     { return "upload" }
func (*upload) Synopsis() string { return "upload version" }
func (*upload) Usage() string {
	return `Usage: upload <version>

Bump version and update related files.

`
}

func (cmd *upload) SetFlags(f *flag.FlagSet) {
	cmd.httpclient.SetFlags(f)
	cmd.zenodo.SetFlags(f)
	cmd.varsfile.SetFlags(f)
	f.StringVar(&cmd.metadata, "metadata", RepoPath(".zenodo.json"), "path to zenodo metadata")
	f.BoolVar(&cmd.publish, "publish", false, "publish the zenodo record")
}

func (cmd *upload) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Fetch version.
	version, err := cmd.varsfile.Get("releaseversion")
	if err != nil {
		return cmd.Error(err)
	}

	// HTTP Client.
	httpclient, err := cmd.httpclient.Client()
	if err != nil {
		return cmd.Error(err)
	}

	// Download archive from github.
	archive := bytes.NewBuffer(nil)
	url := "https://api.github.com/repos/mmcloughlin/addchain/zipball/v" + version
	if err := Download(ctx, archive, httpclient, url); err != nil {
		return cmd.Error(err)
	}

	cmd.Log.Printf("downloaded zip archive bytes %d", archive.Len())

	// Zenodo client.
	client, err := cmd.zenodo.Client(httpclient)
	if err != nil {
		return cmd.Error(err)
	}

	// Load zenodo metadata.
	metadata, err := LoadZenodoMetadata(cmd.metadata)
	if err != nil {
		return cmd.Error(err)
	}

	cmd.Log.Printf("loaded zenodo metadata from %q", cmd.metadata)

	// Set metadata.
	id, err := cmd.varsfile.Get("zenodoid")
	if err != nil {
		return cmd.Error(err)
	}

	d, err := client.DepositionUpdate(ctx, id, metadata)
	if err != nil {
		return cmd.Error(err)
	}

	cmd.Log.Printf("updated zenodo metadata")

	// Clear Zenodo files.
	if err := client.DepositionFilesDeleteAll(ctx, id); err != nil {
		return cmd.Error(err)
	}

	cmd.Log.Printf("cleared all files from zenodo deposit")

	// Upload new Zenodo file.
	archivefilename := fmt.Sprintf("addchain_%s.zip", version)
	file, err := client.DepositionFilesCreate(ctx, id, archivefilename, "application/zip", archive)
	if err != nil {
		return cmd.Error(err)
	}

	cmd.Log.Printf("uploaded file %q", file.Filename)

	// Optionally publish zenodo deposit.
	if cmd.publish {
		if _, err := client.DepositionPublish(ctx, id); err != nil {
			return cmd.Error(err)
		}
		cmd.Log.Printf("published zenodo record")
	} else {
		cmd.Log.Printf("publish skipped: review and publish at %q", d.Links["html"])
	}

	return subcommands.ExitSuccess
}

// Download a file over HTTP.
func Download(ctx context.Context, w io.Writer, c *http.Client, url string) (err error) {
	// Issue GET request.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if errc := res.Body.Close(); errc != nil && err == nil {
			err = errc
		}
	}()

	// Copy to writer.
	_, err = io.Copy(w, res.Body)
	return
}

// LoadZenodoMetadata reads and parses a Zenodo metadata file.
func LoadZenodoMetadata(filename string) (*zenodo.DepositionMetadata, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	m := new(zenodo.DepositionMetadata)
	if err := json.Unmarshal(b, m); err != nil {
		return nil, err
	}

	return m, nil
}
