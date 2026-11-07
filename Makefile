# datacosmos fork-only verb surface (never proposed upstream).
#
# Each verb runs exactly the command of the matching upstream CI step in
# .github/workflows/development.yml; when upstream changes a step, change the
# verb in the same commit. Tool versions come from mise.toml / mise.lock.

GO_BUILD_TAGS := goolm
BINARY        := matterbridge

.PHONY: fmt lint test build check

# CI: go fmt ./...; git status && git diff --exit-code
fmt:
	go fmt ./...
	git diff --exit-code

# CI: golangci-lint (new issues) and golangci-lint (entire codebase).
# Locally both run on the whole codebase, so every finding is reported.
lint:
	golangci-lint run --build-tags=$(GO_BUILD_TAGS)
	golangci-lint run --build-tags=$(GO_BUILD_TAGS) --config=.golangci-incremental.yaml

# CI: go test -tags goolm ./...
test:
	go test -tags $(GO_BUILD_TAGS) ./...

# CI: build-upload (CGO_ENABLED=0, GitHash stamped into version).
build:
	CGO_ENABLED=0 go build -tags $(GO_BUILD_TAGS) \
		-ldflags "-s -X github.com/matterbridge-org/matterbridge/version.GitHash=$$(git log --pretty=format:'%h' -n 1)" \
		-o $(BINARY)

check: fmt lint test build
