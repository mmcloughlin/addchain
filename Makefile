REPO = github.com/mmcloughlin/addchain

.PHONY: fmt
fmt:
	find . -name '*.go' | xargs grep -L '// Code generated' | xargs sed -i.fmtbackup '/^import (/,/)/ { /^$$/ d; }'
	find . -name '*.fmtbackup' -delete
	find . -name '*.go' | xargs grep -L '// Code generated' | xargs gofumports -w -local $(REPO)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate -x ./...

.PHONY: bootstrap
bootstrap:
	GO111MODULE=off go get -u \
		mvdan.cc/gofumpt/gofumports \
		github.com/mna/pigeon
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin v1.19.1
