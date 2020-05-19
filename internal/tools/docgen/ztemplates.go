// Code generated by assets compiler. DO NOT EDIT.
	
package main

var (
	templates = map[string]string{
"/readme.tmpl": "<p align=\"center\">\n  <img src=\"logo.svg\" width=\"40%\" border=\"0\" alt=\"addchain\" />\n</p>\n\n<p align=\"center\">Cryptographic Addition Chain Generation in Go</p>\n\n`addchain` generates short addition chains for exponents of cryptographic\ninterest with [results](#results) rivaling the best hand-optimized chains.\nIntended as a building block in elliptic curve or other cryptographic code\ngenerators.\n\n* Suite of algorithms from academic research: continued fractions,\n  dictionary-based and Bos-Coster heuristics\n* Custom run-length techniques exploit structure of cryptographic exponents\n  with excellent results on Solinas primes\n* Generic optimization methods eliminate redundant operations\n* Simple domain-specific language for addition chain computations\n* Command-line interface or library\n\n## Background\n\nAn [_addition chain_](https://en.wikipedia.org/wiki/Addition_chain) for a\ntarget integer _n_ is a sequence of numbers starting at 1 and ending at _n_\nsuch that every term is a sum of two numbers appearing earlier in the\nsequence. For example, an addition chain for 29 is\n\n```\n1, 2, 4, 8, 9, 17, 25, 29\n```\n\nAddition chains arise in the optimization of exponentiation algorithms with\nfixed exponents. For example, the addition chain above corresponds to the\nfollowing sequence of multiplications to compute `x^29`\n\n```\n x^2 =   x * x\n x^4 = x^2 * x^2\n x^8 = x^4 * x^4\n x^9 =   x * x^8\nx^17 = x^8 * x^9\nx^25 = x^8 * x^17\nx^29 = x^4 * x^25\n```\n\nAn exponentiation algorithm for a fixed exponent _n_ reduces to finding a\n_minimal length addition chain_ for _n_. This is especially relevent in\ncryptography where exponentiation by huge fixed exponents forms a\nperformance-critical component of finite-field arithmetic. In particular,\nconstant-time inversion modulo a prime _p_ is performed by computing `x^(p-2)\n(mod p)`, thanks to [Fermat's Little\nTheorem](https://en.wikipedia.org/wiki/Fermat%27_little_theorem). Square root\nalso reduces to exponentiation for some prime moduli. Finding short addition\nchains for these exponents is one important part of high-performance finite\nfield implementations required for elliptic curve cryptography or RSA.\n\nMinimal addition chain search is famously hard. No practical optimal\nalgorithm is known, especially for cryptographic exponents of size 256-bits\nand up. Given its importance for the performance of cryptographic\nimplementations, implementers devote significant effort to hand-tune addition\nchains. The `addchain` project aims to match or exceed the best\nhand-optimized addition chains using entirely automated approaches, building\non extensive academic research and applying new tweaks that exploit the\nunique nature of cryptographic exponents.\n\n## Results\n\nResults for common cryptographic exponents and delta compared to [best known\nhand-optimized addition\nchains](https://briansmith.org/ecc-inversion-addition-chains-01).\n\n| Name | Length | Best | Delta |\n| ---- | -----: | ---: | ----: |\n{{ range .Results -}}\n{{- if gt .BestKnown 0 -}}\n| [{{ .Name }}](doc/results.md#{{ anchor .Name }}) | {{ .Length }} | {{ .BestKnown }} | {{ if le .Delta 0 }}**{{ printf \"%+d\" .Delta }}**{{ else }}{{ printf \"%+d\" .Delta }}{{ end }} |\n{{ end -}}\n{{ end }}\n\nSee [full results listing](doc/results.md) for more detail and additional\nexponents.\n\n## Usage\n\n### Command-line Interface\n\nInstall:\n\n```\ngo get -u github.com/mmcloughlin/addchain/cmd/addchain\n```\n\nSearch for a curve25519 field inversion addition chain with:\n\n```sh\n{{ include \"internal/examples/cli/cmd.sh\" -}}\n```\n\nOutput:\n\n```\n{{ include \"internal/examples/cli/output\" -}}\n```\n\n### Library\n\nInstall:\n\n```\ngo get -u github.com/mmcloughlin/addchain\n```\n\nAlgorithms all conform to the {{ sym \"alg\" \"ChainAlgorithm\" }} or\n{{ sym \"alg\" \"SequenceAlgorithm\" }} interfaces and can be used directly. However the\nmost user-friendly method uses the {{ pkg \"alg/ensemble\" }} package to\ninstantiate a sensible default set of algorithms and the {{ pkg \"alg/exec\" }}\nhelper to execute them in parallel. The following code uses this method to\nfind an addition chain for curve25519 field inversion:\n\n```go\n{{ snippet \"alg/exec/example_test.go\" \"func Example\" \"^}\" -}}\n```\n\n## Algorithms\n\n### Binary\n\nThe {{ pkg \"alg/binary\" }} package implements the addition chain equivalent\nof the basic [square-and-multiply exponentiation\nmethod](https://en.wikipedia.org/wiki/Exponentiation_by_squaring). It is\nincluded for completeness, but is almost always outperformed by more advanced\nalgorithms below.\n\n### Continued Fractions\n\nThe {{ pkg \"alg/contfrac\" }} package implements the continued fractions\nmethods for addition sequence search introduced by\nBergeron-Berstel-Brlek-Duboc in 1989 and later extended. This approach\nutilizes a decomposition of an addition chain akin to continued fractions,\nnamely\n\n```\n(1,..., k,..., n) = (1,...,n mod k,..., k) {{ \"\\u2297\" }} (1,..., n/k) {{ \"\\u2295\" }} (n mod k).\n```\n\nfor certain special operators {{ \"\\u2297\" }} and {{ \"\\u2295\" }}. This\ndecomposition lends itself to a recursive algorithm for efficient addition\nsequence search, with results dependent on the _strategy_ for choosing the\nauxillary integer _k_. The {{ pkg \"alg/contfrac\" }} package provides a\nlaundry list of strategies from the literature: binary, co-binary,\ndichotomic, dyadic, fermat, square-root and total.\n\n#### References\n\n* {{ bibentry \"contfrac\" }}\n* {{ bibentry \"efficientcompaddchain\" }}\n* {{ bibentry \"gencontfrac\" }}\n* {{ bibentry \"hehcc:exp\" }}\n\n### Bos-Coster Heuristics\n\nBos and Coster described an iterative algorithm for efficient addition\nsequence generation in which at each step a heuristic proposes new numbers\nfor the sequence in such a way that the _maximum_ number always decreases.\nThe [original Bos-Coster paper]({{ biburl \"boscoster\" }}) defined four\nheuristics: Approximation, Divison, Halving and Lucas. Package\n{{ pkg \"alg/heuristic\" }} implements a variation on these heuristics:\n\n* **Approximation:** looks for two elements a, b in the current sequence with sum close to the largest element.\n* **Halving:** applies when the target is at least twice as big as the next largest, and if so it will propose adding a sequence of doublings.\n* **Delta Largest:** proposes adding the delta between the largest two entries in the current sequence.\n\nDivison and Lucas are not implemented due to disparities in the literature\nabout their precise definition and poor results from early experiments.\nFurthermore, this library does not apply weights to the heuristics as\nsuggested in the paper, rather it simply uses the first that applies. However\nboth of these remain [possible avenues for\nimprovement](https://github.com/mmcloughlin/addchain/issues/26).\n\n#### References\n\n* {{ bibentry \"boscoster\" }}\n* {{ bibentry \"github:kwantam/addchain\" }}\n* {{ bibentry \"hehcc:exp\" }}\n* {{ bibentry \"modboscoster\" }}\n* {{ bibentry \"mpnt\" }}\n* {{ bibentry \"speedsubgroup\" }}\n\n### Dictionary\n\nDictionary methods decompose the binary representation of a target integer _n_ into a set of dictionary _terms_, such that _n_\nmay be written as a sum\n\n```\nn = {{ \"\\u2211\" }} 2^{e_i} d_i\n```\n\nfor exponents _e_ and elements _d_ from a dictionary _D_. Given such a decomposition we can construct an addition chain for _n_ by\n\n1. Find a short addition _sequence_ containing every element of the dictionary _D_. Continued fractions and Bos-Coster heuristics can be used here.\n2. Build _n_ from the dictionary terms according to the sum decomposition.\n\nThe efficiency of this approach boils down to the decomposition method. The {{ pkg \"alg/dict\" }} package provides:\n\n* **Fixed Window:** binary representation of _n_ is broken into fixed _k_-bit windows\n* **Sliding Window**: break _n_ into _k_-bit windows, skipping zeros where possible\n* **Run Length**: decompose _n_ into runs of 1s up to a maximal length\n* **Hybrid**: mix of sliding window and run length methods\n\n#### References\n\n* {{ bibentry \"braueraddsubchains\" }}\n* {{ bibentry \"genshortchains\" }}\n* {{ bibentry \"hehcc:exp\" }}\n\n### Runs\n\nThe runs algorithm is a custom variant of the dictionary approach that decomposes\na target into runs of ones. It leverages the observation that building a\ndictionary consisting of runs of 1s of lengths `l_1, l_2, ..., l_k` can itself be\nreduced to:\n\n1. Find an addition sequence containing the run lengths `l_i`. As with\n   dictionary approaches we can use Bos-Coster heuristics and continued\n   fractions here. However here we have the advantage that the `l_i` are\n   typically very _small_, meaning that a wider range of algorithms can\n   be brought to bear.\n2. Use the addition sequence for the run lengths `l_i` to build an addition\n   sequence for the runs themselves `r(l_i)` where `r(e) = 2^e-1`. See\n   {{ sym \"alg/dict\" \"RunsChain\" }}.\n\nThis approach has proved highly effective against cryptographic exponents\nwhich frequently exhibit binary structure, such as those derived from\n[Solinas primes](https://en.wikipedia.org/wiki/Solinas_prime).\n\n> We have not yet found this method described in the literature, so it may be a new development.\n\n### Optimization\n\nClose inspection of addition chains produced by other algorithms revealed\ncases of redundant computation. This motivated a final optimization pass over\naddition chains to remove unecessary steps. The {{ pkg \"alg/opt\" }} package\nimplements the following optimization:\n\n1. Determine _all possible_ ways each element can be computed from those prior.\n2. Count how many times each element is used where it is the _only possible_ way of computing that entry.\n3. Prune elements that are always used in computations that have an alternative.\n\nThese micro-optimizations were vital in closing the gap between `addchain`'s\nautomated approaches and hand-optimized chains. This technique is reminiscent\nof basic passes in optimizing compilers, raising the question of whether\nother [compiler optimizations could apply to addition\nchains](https://github.com/mmcloughlin/addchain/issues/24)?\n\n> We have not yet found this method described in the literature, so it may be a new development.\n\n## Thanks\n\nThank you to [Tom Dean](https://web.stanford.edu/~trdean/), [Riad\nWahby](https://wahby.org/) and [Brian Smith](https://briansmith.org/) for\nadvice and encouragement.\n\n## Contributing\n\nContributions to `addchain` are welcome:\n\n* [Submit bug reports](https://github.com/mmcloughlin/addchain/issues/new) to\n  the issues page.\n* Suggest [test cases](https://github.com/mmcloughlin/addchain/blob/e6c070065205efcaa02627ab1b23e8ce6aeea1db/internal/results/results.go#L62)\n  or update best-known hand-optimized results.\n* Pull requests accepted. Please discuss in the [issues section](https://github.com/mmcloughlin/addchain/issues)\n  before starting significant work.\n\n## License\n\n`addchain` is available under the [BSD 3-Clause License](LICENSE).\n",
"/results.tmpl": "# Results\n\n{{ range .Results -}}\n* [{{ .Name }}](#{{ anchor .Name }})\n{{ end }}\n\n{{ range .Results -}}\n## {{ .Name }}\n\n| Property | Value |\n| --- | ----- |\n| _N_ | `{{ .N.String }}` |\n| _d_ | `{{ .D }}` |\n| _N_-_d_ | `{{ printf \"%x\" .Target }}` |\n| Length | {{ .Length }} |\n| Algorithm | `{{ .AlgorithmName }}` |\n{{- if gt .BestKnown 0 }}\n| Best Known | {{ .BestKnown }} |\n| Delta | {{ printf \"%+d\" .Delta }} |\n{{ end }}\n\nAddition chain produced by `addchain`:\n\n```go\n{{ include (printf \"internal/results/testdata/%s.golden\" .Slug) }}```\n\n{{ end }}\n",
}
)