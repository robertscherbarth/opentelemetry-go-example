honeycomb-build: ## honeycomb build setup
	docker compose --env-file .env-honeycomb -f docker-compose-honeycomb.yaml build

honeycomb-up: ## honeycomb run setup
	docker compose --env-file .env-honeycomb -f docker-compose-honeycomb.yaml up

jaeger-build: ## jaeger build setup
	docker compose -f docker-compose-jaeger.yaml build

jaeger-up: ## jaeger run setup
	docker compose -f docker-compose-jaeger.yaml up

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