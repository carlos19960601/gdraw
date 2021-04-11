.PHONY: build

ICON="🔞"

# Project binaries
COMMANDS=gdraw

BINARIES=$(addprefix bin/,$(COMMANDS))

all: binaries

FORCE:
define BUILD_BINARY
@echo "$(ICON) $@"
@go build -o $@ ./$<
endef

build: cmd/gdraw
	@echo "$(ICON) $@"
	@go build -o bin/gdraw ./$<