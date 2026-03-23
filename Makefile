PREFIX     := oxyl

OUTPUT_DIR := ./builds
LDFLAGS    := -ldflags "-s -w"
GCFLAGS    := -gcflags "all=-trimpath=$(shell pwd)"

MODULES := $(bash -c 'find . -type f -iname "main.go" -exec dirname {} \;')

all: build

# we should always generate the corresponding generated values for the code. (Consistency)
generate:
	$(go generate ./shared/pkg/version/version.go)

build: clean generate
	@mkdir -p $(OUTPUT_DIR)
	@for pkg in $(MODULES); do \
		name=$$(basename $$pkg); \
		GOOS=linux GOARCH=amd64 go build $(LDFLAGS) $(GCFLAGS) -o $(OUTPUT_DIR)/$(PREFIX)-$$name $$pkg; \
	done

clean:
	$(shell rm -rf $(OUTPUT_DIR))


.PHONY: all build clean