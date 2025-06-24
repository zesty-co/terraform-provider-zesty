GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
TF_PLUGIN_DOCS = $(shell pwd)/bin/tfplugindocs
GORELEASER = $(shell pwd)/bin/goreleaser

.PHONY: install
install:
	go install .

.PHONY: tools
tools:
	GOBIN=$(shell pwd)/bin go install tool

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run --output.code-climate.path=stdout | head -n 1 | tee gl-code-quality-report.json | jq -r '.[] | "\(.location.path):\(.location.lines.begin) \(.description)"'

.PHONY: lint-fix
lint-fix:
	$(GOLANGCI_LINT) run --fix

.PHONY: test
test:
	go test $(shell go list ./...) -coverprofile=coverage.tmp
	cat ./coverage.tmp | grep -v 'main.go' > ./coverage.out
	rm ./coverage.tmp
	go tool cover -html=coverage.out -o=coverage.html
	go tool cover -func=coverage.out

.PHONY: docs
docs:
	$(TF_PLUGIN_DOCS) generate --rendered-provider-name "Zesty" --ignore-deprecated --provider-name terraform-provider-zesty

.PHONY: format
format:
	go fmt ./...

.PHONY: generate
generate:
	GOBIN=$(shell pwd)/bin go generate ./...

.PHONY: generate-all
generate-all: format lint generate docs

.PHONY: build
build:
	go build -ldflags="-w -s" -o bin/terraform-provider-zesty main.go

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: terraformrc
terraformrc:
	@bin=$$(go env GOBIN); \
	if [ -z "$$bin" ]; then bin="$$HOME/go/bin"; fi; \
	printf 'provider_installation {\n\n  dev_overrides {\n      "registry.terraform.io/zesty/zesty" = "%s"\n  }\n\n  direct {}\n}\n' "$$bin" > $$HOME/.terraformrc

.PHONY: release
release:
	$(GORELEASER) release --clean --config=.github/.goreleaser.yaml
