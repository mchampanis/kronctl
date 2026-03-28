BINARY    := kronctl.exe
CMD       := ./cmd/kronctl
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS   := -X main.version=$(VERSION)

.PHONY: all release vet clean

all: ## Build debug binary
	CGO_ENABLED=1 go build -race -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

release: ## Build release binary (stripped)
	CGO_ENABLED=1 go build -ldflags "-s -w $(LDFLAGS)" -o $(BINARY) $(CMD)

vet: ## Run static analysis
	go vet ./...

clean: ## Remove build artifacts
	rm -f $(BINARY)
