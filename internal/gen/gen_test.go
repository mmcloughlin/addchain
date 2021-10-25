package gen

import "testing"

func TestLoadBuiltinTemplates(t *testing.T) {
	names := BuiltinTemplateNames()
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			_, err := BuiltinTemplate(name)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
