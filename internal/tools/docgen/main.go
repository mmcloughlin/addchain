package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/mmcloughlin/addchain/internal/results"
)

//go:generate assets -d templates/ -o ztemplates.go -map templates

func main() {
	log.SetPrefix("docgen: ")
	log.SetFlags(0)
	if err := mainerr(); err != nil {
		log.Fatal(err)
	}
}

var (
	typ    = flag.String("type", "", "documentation type")
	tmpl   = flag.String("tmpl", "", "explicit template file (overrides -type)")
	output = flag.String("output", "", "path to output file (default stdout)")
)

func mainerr() (err error) {
	flag.Parse()

	// Initialize template.
	t := template.New("doc")

	t.Funcs(template.FuncMap{
		"include": include,
		"anchor":  anchor,
	})

	// Load template.
	s, err := load()
	if err != nil {
		return err
	}

	if _, err := t.Parse(s); err != nil {
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
		defer func() {
			if errc := f.Close(); errc != nil && err == nil {
				err = errc
			}
		}()

		w = f
	}

	if err := t.Execute(w, data); err != nil {
		return err
	}

	return nil
}

// load template.
func load() (string, error) {
	// Prefer explicit filename, if provided.
	if *tmpl != "" {
		b, err := ioutil.ReadFile(*tmpl)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	// Otherwise expect a named type.
	if *typ == "" {
		return "", errors.New("missing documentation type")
	}
	key := fmt.Sprintf("/%s.tmpl", *typ)
	s, ok := templates[key]
	if !ok {
		return "", fmt.Errorf("unknown documentation type %q", *typ)
	}

	return s, nil
}

// include template function.
func include(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// anchor returns the anchor for a heading in Github.
func anchor(heading string) string {
	r := strings.NewReplacer(" ", "-", "(", "", ")", "")
	return r.Replace(strings.ToLower((heading)))
}
