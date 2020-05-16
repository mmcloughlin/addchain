package cli

import (
	"io"
	"io/ioutil"
	"os"
)

// OpenInput is a convenience for possibly opening an input file, or otherwise returning standard in.
func OpenInput(filename string) (string, io.ReadCloser, error) {
	if filename == "" {
		return "<stdin>", ioutil.NopCloser(os.Stdin), nil
	}
	f, err := os.Open(filename)
	return filename, f, err
}
