run-crud: ## start the crud service
	go run ./cmd/crud/.

run-caller: ## start the caller
	go run ./cmd/caller/.

users-get-all: ## curl all users from the crud service
	curl -i -X GET http://localhost:8081/users

user-delete: ## curl delete single user
	curl -i -X DELETE http://localhost:8081/users/acabbaba-f65f-4a29-9091-19b3264dafff

fmt: ## Run go fmt against code
	go fmt ./pkg/... ./cmd/...

vet: ## Run go vet against code
	go vet ./pkg/... ./cmd/...

test: ## Runs the tests
	go test -cover -race -short -v $(shell go list ./... | grep -v /vendor/ )

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: fmt vet test help users-get-all run-caller