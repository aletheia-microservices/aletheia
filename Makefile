# Aletheia must run from the repository root: it loads the Blueprint apps through the replace
# directives in go.mod and reads registry/ and config/ relative to the working directory.

BIN  := bin/aletheia
APP  ?= postnotification
ARGS ?=

.DEFAULT_GOAL := help

# shell snippet that prints how to run aletheia as $$cmd, highlighted when stdout is a terminal
PRINT_USAGE = if [ -t 1 ]; then hl=$$(printf '\033[1;32m'); rs=$$(printf '\033[0m'); fi; \
	printf '\n  usage (from the repository root):\n\n      %s%s [flags] <app>%s\n      %s%s -h%s    lists the flags and registered apps\n\n' \
		"$$hl" "$$cmd" "$$rs" "$$hl" "$$cmd" "$$rs"

.PHONY: help
help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-18s %s\n", $$1, $$2}'

.PHONY: build
build: ## build the aletheia binary into bin/
	go build -o $(BIN) ./cmd/aletheia
	@cmd=./$(BIN); echo "built $(BIN)"; $(PRINT_USAGE)

.PHONY: install
install: ## install aletheia into $GOBIN or $GOPATH/bin (rerun after make registry)
	go install ./cmd/aletheia
	@dir=$$(go env GOBIN); [ -n "$$dir" ] || dir=$$(go env GOPATH | cut -d: -f1)/bin; \
	cmd=aletheia; echo "installed $$dir/aletheia"; \
	case ":$$PATH:" in *":$$dir:"*) ;; *) echo "note: $$dir is not on your PATH, add it with: export PATH=\"\$$PATH:$$dir\"";; esac; \
	$(PRINT_USAGE)

.PHONY: run
run: ## analyze APP (default: postnotification) with extra ARGS, e.g. make run APP=sockshop ARGS=-debug
	@go build -o $(BIN) ./cmd/aletheia
	./$(BIN) $(ARGS) $(APP)

.PHONY: test
test: ## run unit and integration tests
	go test ./tests/...

.PHONY: test-short
test-short: ## run unit tests only (skips the full pipeline)
	go test -short ./tests/...

.PHONY: verify
verify: ## compare warning counts of every app against scripts/verify/expected/
	scripts/verify/verify.sh

.PHONY: vet
vet: ## run go vet
	go vet ./cmd/... ./internal/... ./tests/... ./scripts/...

.PHONY: fmt
fmt: ## format all Go code
	gofmt -w cmd internal tests scripts

.PHONY: registry
registry: ## regenerate the app registry and go.mod entries from registry/*.yaml, then rebuild
	go run ./scripts/gen-registry
	@$(MAKE) --no-print-directory build

# the docs targets need uv, which the official install script puts in ~/.local/bin
docs-install docs-install-brew docs-serve docs-build: export PATH := $(HOME)/.local/bin:$(PATH)

.PHONY: docs-install
docs-install: ## install the docs dependencies (installs uv with its install script if missing)
	@command -v uv >/dev/null || curl -LsSf https://astral.sh/uv/install.sh | sh
	uv sync --group docs

.PHONY: docs-install-brew
docs-install-brew: ## install the docs dependencies (installs uv with Homebrew if missing)
	@command -v uv >/dev/null || brew install uv
	uv sync --group docs

.PHONY: docs-serve
docs-serve: ## serve the docs site with live reload and open it in the browser
	uv run --group docs mkdocs serve --open

.PHONY: docs-build
docs-build: ## build the docs site into site/ (fails on broken links)
	uv run --group docs mkdocs build --strict

.PHONY: clean
clean: ## remove build artifacts
	rm -rf bin site
