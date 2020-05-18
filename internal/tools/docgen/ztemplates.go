// Code generated by assets compiler. DO NOT EDIT.
	
package main

var (
	templates = map[string]string{
"/readme.tmpl": "<p align=\"center\">\n  <img src=\"logo.svg\" width=\"40%\" border=\"0\" alt=\"addchain\" />\n</p>\n\n<p align=\"center\">Cryptographic Addition Chain Generation in Go</p>\n\n`addchain` generates short addition chains for exponents of cryptographic\ninterest with [results](#results) rivaling the best hand-optimized chains.\nIntended as a building block in elliptic curve or other cryptographic code\ngenerators.\n\n* Suite of algorithms from academic research: continued fractions,\n  dictionary-based and Bos-Coster heuristics\n* Custom run-length techniques exploit structure of cryptographic exponents\n  with excellent results on Solinas primes\n* Generic optimization methods eliminate redundant operations\n* Simple domain-specific language for addition chain computations\n* Command-line interface or library\n\n## Results\n\n| Name | Length | Best | Delta |\n| ---- | -----: | ---: | ----: |\n{{ range .Results -}}\n{{- if gt .BestKnown 0 -}}\n| [{{ .Name }}](doc/results.md#{{ anchor .Name }}) | {{ .Length }} | {{ .BestKnown }} | {{ printf \"%+d\" .Delta }} |\n{{ end -}}\n{{ end }}\n\n## Usage\n\n### Command-line Interface\n\nInstall:\n\n```\ngo get -u github.com/mmcloughlin/addchain/cmd/addchain\n```\n\nSearch for a curve25519 field inversion addition chain with:\n\n```sh\n{{ include \"internal/examples/cli/cmd.sh\" -}}\n```\n\nOutput:\n\n```\n{{ include \"internal/examples/cli/output\" -}}\n```\n\n### Library\n\nInstall:\n\n```\ngo get -u github.com/mmcloughlin/addchain\n```\n\nAlgorithms all conform to the `alg.ChainAlgorithm` or `alg.SequenceAlgorithm`\ninterfaces and can be used directly. However the most user-friendly method\nuses the `alg/ensemble` package to instantiate a sensible default set of\nalgorithms and the `alg/exec` helper to execute them in parallel. The\nfollowing code uses this method to find an addition chain for curve25519\nfield inversion:\n\n```go\n{{ snippet \"alg/exec/example_test.go\" \"func Example\" \"^}\" -}}\n```\n\n## License\n\n`addchain` is available under the [BSD 3-Clause License](LICENSE).\n",
"/results.tmpl": "# Results\n\n{{ range .Results -}}\n* [{{ .Name }}](#{{ anchor .Name }})\n{{ end }}\n\n{{ range .Results -}}\n## {{ .Name }}\n\n| Property | Value |\n| --- | ----- |\n| _N_ | `{{ .N.String }}` |\n| _d_ | `{{ .D }}` |\n| _N_-_d_ | `{{ printf \"%x\" .Target }}` |\n| Length | {{ .Length }} |\n{{- if gt .BestKnown 0 }}\n| Best Known | {{ .BestKnown }} |\n| Delta | {{ printf \"%+d\" .Delta }} |\n{{ end }}\n\nAddition chain produced by `addchain`:\n\n```go\n{{ include (printf \"internal/results/testdata/%s.golden\" .Slug) }}```\n\n{{ end }}\n",
}
)