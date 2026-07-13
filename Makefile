.PHONY: start
start: ## Start repl

	go run ./cmd/repl/main.go

.PHONY: test
test: ## Test all

	go test ./lexer
	go test ./ast
	go test ./parser
	go test ./evaluator

.PHONY: help
help:

	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        sort | \
         awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'