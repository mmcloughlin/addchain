package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
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
		"snippet": snippet,
		"anchor":  anchor,
		"pkg":     pkg,
		"sym":     sym,
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

// snippet of a file between start and end regular expressions.
func snippet(filename, start, end string) (string, error) {
	// Parse regular expressions.
	startx, err := regexp.Compile(start)
	if err != nil {
		return "", err
	}

	endx, err := regexp.Compile(end)
	if err != nil {
		return "", err
	}

	// Read the full file.
	data, err := include(filename)
	if err != nil {
		return "", err
	}

	// Collect matched lines.
	var buf bytes.Buffer
	output := false
	s := bufio.NewScanner(strings.NewReader(data))
	for s.Scan() {
		line := s.Text()
		if startx.MatchString(line) {
			output = true
		}
		if output {
			fmt.Fprintln(&buf, line)
		}
		if endx.MatchString(line) {
			output = false
		}
	}

	return buf.String(), nil
}

// anchor returns the anchor for a heading in Github.
func anchor(heading string) string {
	r := strings.NewReplacer(" ", "-", "(", "", ")", "", "/", "")
	return r.Replace(strings.ToLower((heading)))
}

// pkg returns markdown for a package with a link to documentation.
func pkg(name string) string {
	return fmt.Sprintf("[`%s`](%s)", name, pkgurl(name))
}

// sym returns markdown for a symbol with a
func sym(pkg, name string) string {
	return fmt.Sprintf("[`%s.%s`](%s#%s)", path.Base(pkg), name, pkgurl(pkg), name)
}

// pkgurl returns url to go.dev documentation on the given sub-package.
func pkgurl(name string) string {
	return "https://pkg.go.dev/" + path.Join("github.com/mmcloughlin/addchain", name)
}
