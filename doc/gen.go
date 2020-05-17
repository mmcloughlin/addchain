// +build ignore

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/mmcloughlin/addchain/internal/results"
)

func main() {
	log.SetPrefix("docgen: ")
	log.SetFlags(0)
	if err := mainerr(); err != nil {
		log.Fatal(err)
	}
}

var (
	tmpl   = flag.String("tmpl", "", "template")
	output = flag.String("output", "", "path to output file (default stdout)")
)

func mainerr() error {
	flag.Parse()

	// Initialize template.
	t := template.New("doc")

	t.Funcs(template.FuncMap{
		"include": include,
	})

	// Load template.
	b, err := ioutil.ReadFile(*tmpl)
	if err != nil {
		return err
	}

	if _, err := t.Parse(string(b)); err != nil {
		return err
	}

	// Prepare template data.
	data := map[string]interface{}{
		"Results": results.Results,
	}

	// Execute.
	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return err
		}
		defer f.Close()

		w = f
	}

	if err := t.Execute(w, data); err != nil {
		return err
	}

	return nil
}

func include(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
