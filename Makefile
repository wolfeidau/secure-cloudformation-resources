APPNAME := cf-security-transform
STAGE ?= dev
BRANCH ?= master

GOLANGCI_VERSION = 1.31.0

GIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

default: clean build archive deploy

ci: clean lint test
.PHONY: ci

LDFLAGS := -ldflags="-s -w -X main.buildDate=${BUILD_DATE} -X main.commit=${GIT_HASH}"

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

bin/mockgen:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/golang/mock/mockgen

bin/gcov2lcov:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov

bin/rain:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/aws-cloudformation/rain

clean:
	@echo "--- clean all the things"
	@rm -rf ./dist
	@rm -f ./handler.zip
	@rm -f ./*.out.yaml
.PHONY: clean

lint: bin/golangci-lint
	@echo "--- lint all the things"
	@bin/golangci-lint run
.PHONY: lint

test: bin/gcov2lcov
	@echo "--- test all the things"
	@go test -v -covermode=count -coverprofile=coverage.txt ./ ./pkg/... ./internal/...
	@bin/gcov2lcov -infile=coverage.txt -outfile=coverage.lcov
.PHONY: test

build:
	@echo "--- build all the things"
	@mkdir -p dist
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -trimpath -o dist ./cmd/...
.PHONY: build

archive:
	@echo "--- build an archive"
	@cd dist && zip -X -9 -r ../handler.zip *-lambda
.PHONY: archive

deploy:
	@echo "--- deploy db into aws"
	bin/rain deploy sam/transform/transform.yaml $(APPNAME)-$(STAGE)-$(BRANCH) \
		--params AppName=$(APPNAME),Stage=$(STAGE),Branch=$(BRANCH) --force
.PHONY: deploy

deploy-test:
	@echo "--- deploy db into aws"
	bin/rain deploy sam/test/test.yaml $(APPNAME)-$(STAGE)-$(BRANCH)-test \
		--params AppName=$(APPNAME),Stage=$(STAGE),Branch=$(BRANCH) --force
.PHONY: deploy-test
