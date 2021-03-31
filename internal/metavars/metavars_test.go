package metavars

import (
	"bytes"
	"reflect"
	"testing"
)

func TestRoundtrip(t *testing.T) {
	cases := []*File{
		// TODO: mixed doc comments
		// TODO: quoted strings
		{
			Package: "standard",
			Properties: []Property{
				{Name: "A", Doc: "A is the first variable.", Value: "one"},
				{Name: "B", Doc: "B is the second variable.", Value: "two"},
				{Name: "C", Doc: "C is the third variable.", Value: "three"},
			},
		},
		{
			Package: "single",
			Properties: []Property{
				{Name: "A", Doc: "A is the first variable.", Value: "one"},
			},
		},
		{
			Package: "empty",
		},
		{
			Package: "nodoc",
			Properties: []Property{
				{Name: "A", Value: "one"},
				{Name: "B", Value: "two"},
			},
		},
	}
	for _, f := range cases {
		t.Run(f.Package, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			if err := Write(buf, f); err != nil {
				t.Fatal(err)
			}

			got, err := Read(buf)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, f) {
				t.Logf("got    = %#v", got)
				t.Logf("expect = %#v", f)
				t.Fatal("roundtrip fail")
			}
		})
	}
}
