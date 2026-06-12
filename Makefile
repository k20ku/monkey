.PHONY: start
start: ## Start repl

	go run ./cmd/repl/main.go

.PHONY: help
help:

        @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
         sort | \
         awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'