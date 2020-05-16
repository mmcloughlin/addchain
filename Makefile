REPO = github.com/mmcloughlin/addchain

.PHONY: fmt
fmt:
	find . -name '*.go' | xargs grep -L '// Code generated' | xargs sed -i.fmtbackup '/^import (/,/)/ { /^$$/ d; }'
	find . -name '*.fmtbackup' -delete
	find . -name '*.go' | xargs grep -L '// Code generated' | xargs gofumports -w -local $(REPO)
	find . -name '*.go' | grep -v _test | xargs grep -L '// Code generated' | xargs mathfmt -w
	find . -name '*.go' | xargs grep -L '// Code generated' | xargs bib process -bib doc/references.bib -w

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate -x ./...
