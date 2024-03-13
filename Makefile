.PHONY: all build test lint fmt clean run deps

SHELL=/bin/bash -o pipefail

VERSION=0.1.0
GOBUILD=env GOOS=linux GOARCH=amd64 go build -v
GOCLEAN=go clean
LINT_VERSION=v1.42.1

# Reports
LINT_REPORT="lint_report.xml"
VET_REPORT="vet_report.json"
SEC_REPORT="sec_report.json"


importer: generate
	$(GOBUILD) -o importer ./cmd/importer

api: generate
	$(GOBUILD) -o api ./cmd/api

test-api:
	go test -v ./cmd/api

test-importer:
	MONGODB_DATABASE_NAME=geolocation_test go test -v ./cmd/importer

lint: install_lint
	golangci-lint run --tests=false --out-format=checkstyle ./... > $(LINT_REPORT)

fmt:
	go fmt ./...

clean:
	$(GOCLEAN) ./...
	rm -f api importer


deps:
	go mod tidy

install_lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin $(LINT_VERSION)

vet:
	go vet -json ./... 2> $(VET_REPORT)

sec: tools
	gosec -fmt sonarqube -out $(SEC_REPORT) ./...

generate:
	go generate ./...

version:
	@echo $(VERSION)
