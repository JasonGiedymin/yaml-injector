# Global Settings
NO_COLOR=\033[0m
TEXT_COLOR=\033[1m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Go Opts
GO_EXEC=go

help:
	@echo "$(OK_COLOR)-----------------------Commands:----------------------$(NO_COLOR)"
	@echo "$(TEXT_COLOR) help            : Test help listing $(NO_COLOR)"

	@echo "$(OK_COLOR)\n                   Helpers $(NO_COLOR)"
	@echo "$(TEXT_COLOR) benchmark       : Run benchmark $(NO_COLOR)"
	@echo "$(TEXT_COLOR) cleans          : Cleans test dir $(NO_COLOR)"
	
	@echo "$(OK_COLOR)\n                   Tests $(NO_COLOR)"
	@echo "$(TEXT_COLOR) tests           : Runs all unit tests $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-datafile   : Test with datafile to stdout $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-preview    : Test preview functionality to stdout $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-stdin      : Stdin JSON usage with debug output $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-stdin      : Run with stdin JSON usage to stdout $(NO_COLOR)"

	@echo "$(OK_COLOR)\n                   Integration Tests $(NO_COLOR)"
	@echo "$(TEXT_COLOR) itest-stdin     : Integration test for stdin. $(NO_COLOR)"
	@echo "$(OK_COLOR)------------------------------------------------------$(NO_COLOR)"

all: help

clean:
	@echo "$(TEXT_COLOR)==> Cleaning up files: $(NO_COLOR)"
	@cp test/original.input.yaml test/input.yaml
	@rm test/input.yaml.*
	@echo Done.

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
