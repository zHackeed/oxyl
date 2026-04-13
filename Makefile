PREFIX     := oxyl

OUTPUT_DIR := ./builds

MODULES := $(shell bash -c 'find . -type f -iname "main.go" -exec dirname {} \;')

all: build

generate:
	@echo "Generating buildable dependencies..."
	@go generate ./shared/pkg/version/version.go

	@if command -v buf >/dev/null 2>&1; then \
  		echo "buf found, generating protobuf files..."; \
	else \
		echo "buf not found, please install buf to generate protobuf files"; \
		exit 1; \
	fi

build: clean generate
	@mkdir -p $(OUTPUT_DIR)
	@echo "Building modules: $(MODULES)"
	@for pkg in $(MODULES); do \
		name=$$(basename $$pkg); \
		echo "Building $(PREFIX)-$$name"; \
		GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o $(OUTPUT_DIR)/$(PREFIX)-$$name $$pkg; \
	done

clean:
	$(shell rm -rf $(OUTPUT_DIR))


.PHONY: all build clean
