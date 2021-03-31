package meta

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mmcloughlin/addchain/internal/print"
)

func (p *Properties) WriteCitation(w io.Writer) error {
	// Determine release time.
	date, err := p.ReleaseTime()
	if err != nil {
		return fmt.Errorf("release date: %w", err)
	}

	// Use tabwriter for field alignment.
	tw := print.NewTabWriter(w, 1, 4, 1, ' ', 0)

	field := func(key, value string) { tw.Linef("    %s\t=\t%s,", key, value) }
	str := func(key, value string) { field(key, "{"+value+"}") }

	tw.Linef("@misc{addchain,")
	str("title", "addchain: Cryptographic Addition Chain Generation in Go")
	str("author", "Michael B. McLoughlin")
	field("year", strconv.Itoa(date.Year()))
	field("month", strings.ToLower(date.Month().String()[:3]))
	str("version", p.ReleaseVersion)
	str("license", "BSD 3-Clause License")
	str("doi", p.DOI)
	str("url", p.DOIURL())
	tw.Linef("}")
	tw.Flush()

	return tw.Error()
}

func (p *Properties) Citation(w io.Writer) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := p.WriteCitation(buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
