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
	@echo "$(TEXT_COLOR) benchmark       : Run benchmark $(NO_COLOR)"
	@echo "$(TEXT_COLOR) tests           : Runs all unit tests $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-datafile   : Test with datafile to stdout $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-preview    : Test preview functionality to stdout $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-stdin      : Stdin JSON usage with debug output $(NO_COLOR)"
	@echo "$(TEXT_COLOR) run-stdin       : Run with stdin JSON usage to stdout $(NO_COLOR)"
	@echo "$(OK_COLOR)------------------------------------------------------$(NO_COLOR)"

all: help

tests:
	$(GO_EXEC) test ./... -v

benchmark:
	$(GO_EXEC) test ./... -bench=. -benchmem -benchtime 1s

test-datafile:
	@echo "$(OK_COLOR)==> Data file usage: $(NO_COLOR)"
	$(GO_EXEC) run yaml-injector.go \
	  --debug 
	  --file test/input.yaml \
	  --using test/data.yaml \
	  --key "a" \
	  inject into "a"

test-preview:
	@echo "$(OK_COLOR)==> Data file usage but using test run to see expected output on stdout $(NO_COLOR)"
	$(GO_EXEC) run yaml-injector.go \
    --debug \
    --test \
    --file test/input.yaml \
    --using test/data.yaml \
    --key "a" \
    inject into "a"

test-stdin:
	@echo "$(OK_COLOR)==> Stdin JSON usage with debug output $(NO_COLOR)"
	(echo '{"a":1}') | $(GO_EXEC) run yaml-injector.go \
	  --debug \
	  --test \
    --file test/input.yaml \
    --using test/data.yaml \
    --key "a" inject into "a"

run-stdin:
	@echo "$(OK_COLOR)==> Run with stdin JSON usage to stdout $(NO_COLOR)"
	@(echo '{"a":1}') | $(GO_EXEC) run yaml-injector.go \
		--test \
    --file test/input.yaml \
    --using test/data.yaml \
    --key "a" inject into "a"
