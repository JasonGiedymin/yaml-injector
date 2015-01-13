# Global Settings
NO_COLOR=\033[0m
TEXT_COLOR=\033[1m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

VERSION=0.2.0

# Go Opts
GO_EXEC=go

help:
	@echo "$(OK_COLOR)-----------------------Commands:----------------------$(NO_COLOR)"
	@echo "$(TEXT_COLOR) help            : Test help listing $(NO_COLOR)"

	@echo "$(OK_COLOR)\n                   Helpers $(NO_COLOR)"
	@echo "$(TEXT_COLOR) benchmark       : Run benchmark $(NO_COLOR)"
	@echo "$(TEXT_COLOR) clean           : Cleans project dirs (runs all clean cmds) $(NO_COLOR)"
	
	@echo "$(OK_COLOR)\n                   Tests $(NO_COLOR)"
	@echo "$(TEXT_COLOR) clean test      : Cleans test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) tests           : Runs all unit tests $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-datafile   : Test with datafile to stdout $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-preview    : Test preview functionality to stdout $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-stdin      : Stdin JSON usage with debug output $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-stdin      : Run with stdin JSON usage to stdout $(NO_COLOR)"

	@echo "$(OK_COLOR)\n                   Integration Tests $(NO_COLOR)"
	@echo "$(TEXT_COLOR) itest-stdin     : Integration test for stdin. $(NO_COLOR)"

	@echo "$(OK_COLOR)\n                   Builds $(NO_COLOR)"
	@echo "$(TEXT_COLOR) build           : Build all binaries. $(NO_COLOR)"
	@echo "$(TEXT_COLOR) build-linux64   : Build linux 64 bit binary. $(NO_COLOR)"
	@echo "$(TEXT_COLOR) build-darwin64  : Build darwin 64 bit binary. $(NO_COLOR)"
	@echo "$(TEXT_COLOR) clean-target    : Cleans build targets dir. $(NO_COLOR)"
	@echo "$(OK_COLOR)------------------------------------------------------$(NO_COLOR)"

all: help

clean-test:
	@echo "$(TEXT_COLOR)==> Cleaning up files: $(NO_COLOR)"
	@cp test/original.input.yaml test/input.yaml
	@if [ -e test/input.yaml.* ]; then\
		rm test/input.yaml.* ;\
	fi;
	@echo Done.

clean-target:
	@echo "$(TEXT_COLOR)==> Cleaning up target directory: $(NO_COLOR)"
	@rm -R target/

clean: clean-test clean-target

tests:
	$(GO_EXEC) test ./... -v

benchmark:
	$(GO_EXEC) test ./... -bench=. -benchmem -benchtime 1s

test-datafile:
	@echo "$(TEXT_COLOR)==> Data file usage: $(NO_COLOR)"
	$(GO_EXEC) run yaml-injector.go \
	  --debug \
	  --test \
	  --file test/input.yaml \
	  --using test/data.yaml \
	  --key "a" \
	  inject into "a"

test-preview:
	@echo "$(TEXT_COLOR)==> Data file usage but using test run to see expected output on stdout $(NO_COLOR)"
	$(GO_EXEC) run yaml-injector.go \
    --debug \
    --test \
    --file test/input.yaml \
    --using test/data.yaml \
    --key "a" \
    inject into "a"

test-stdin:
	@echo "$(TEXT_COLOR)==> Stdin JSON usage with debug output $(NO_COLOR)"
	(echo '{"a":1}') | $(GO_EXEC) run yaml-injector.go \
	  --debug \
	  --test \
    --file test/input.yaml \
    --using test/data.yaml \
    --key "a" inject into "a"

_itest-stdin:
	@echo "$(TEXT_COLOR)==> Integration test for stdin: $(NO_COLOR)"

	@echo Original File content prior to yaml-injector:
	@cat test/input.yaml
	@echo

	$(GO_EXEC) run yaml-injector.go \
	  --debug \
	  --file test/input.yaml \
	  --using test/data.yaml \
	  --key "a" \
	  inject into "a"

	@echo File content post yaml-injector:
	@cat test/input.yaml
	@echo

	@echo Backup Files created:
	@ls -l test/input.yaml.*
	@echo

	@echo "$(TEXT_COLOR)==>Testing checksum of output file... $(NO_COLOR)"
	@echo "$(WARN_COLOR)Note: Failures will need users to run \"make clean\" manually. $(NO_COLOR)"
	@echo
	@test $(shell md5 -q test/input.yaml) == "641fcad0eeab2f71d25f19c994a9f793"
	@echo "$(OK_COLOR)PASS $(NO_COLOR)"
	@echo

itest-stdin: _itest-stdin clean

build-linux64:
	@echo "$(TEXT_COLOR)==> Building Linux 64bit binary... $(NO_COLOR)"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO_EXEC) build -o ./target/yaml-injector-$(VERSION)-linux64
	@tar -cvjf ./target/yaml-injector-$(VERSION)-linux64.tar.bz2 ./target/yaml-injector-$(VERSION)-linux64

build-darwin64:
	@echo "$(TEXT_COLOR)==> Building Darwin 64bit binary... $(NO_COLOR)"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GO_EXEC) build -o ./target/yaml-injector-$(VERSION)-darwin64
	@tar -cvjf ./target/yaml-injector-$(VERSION)-darwin64.tar.bz2 ./target/yaml-injector-$(VERSION)-darwin64

build: build-linux64 build-darwin64
