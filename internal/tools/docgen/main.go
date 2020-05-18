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
	"regexp"
	"strings"
	"text/template"
	"unicode"

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
	var lines []string
	output := false
	s := bufio.NewScanner(strings.NewReader(data))
	for s.Scan() {
		line := s.Text()
		if startx.MatchString(line) {
			output = true
		}
		if output {
			lines = append(lines, line)
		}
		if endx.MatchString(line) {
			output = false
		}
	}

	// Strip common prefix.
	var buf bytes.Buffer
	for _, line := range stripCommonWhitespacePrefix(lines) {
		fmt.Fprintln(&buf, line)
	}

	return buf.String(), nil
}

// anchor returns the anchor for a heading in Github.
func anchor(heading string) string {
	r := strings.NewReplacer(" ", "-", "(", "", ")", "", "/", "")
	return r.Replace(strings.ToLower((heading)))
}

// stripCommonWhitespacePrefix returns a list of strings with any common
// whitespace prefix removed. Empty strings are preserved but ignored for the
// purposes of common prefix identification.
func stripCommonWhitespacePrefix(strs []string) []string {
	// Find common prefix, ignoring empty strings.
	nonempty := []string{}
	for _, str := range strs {
		if str != "" {
			nonempty = append(nonempty, str)
		}
	}
	n := commonWhitespacePrefixLen(nonempty)

	// Strip common prefix.
	stripped := []string{}
	for _, str := range strs {
		if str != "" {
			stripped = append(stripped, str[n:])
		} else {
			stripped = append(stripped, "")
		}
	}

	return stripped
}

// commonWhitespacePrefixLen returns the length of the common whitespace prefix
// in strs.
func commonWhitespacePrefixLen(strs []string) int {
	l := 0
	for ; ; l++ {
		for _, str := range strs {
			if !(l < len(str) && isspace(str[l]) && str[l] == strs[0][l]) {
				return l
			}
		}
	}
}

// isspace reports whether ch is a whitespace character.
func isspace(ch byte) bool {
	return unicode.IsSpace(rune(ch))
}
