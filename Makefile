ENV	?=	dev
PORT ?= 8080
COLOR_RESET=\033[0m
COLOR_GREEN=\033[32m
COLOR_RED=\033[31m

.PHONY: run
run:
	@echo "$(COLOR_GREEN)Running...$(COLOR_RESET)"
	docker-compose up

.PHONY: build
build:
	@echo "$(COLOR_GREEN)Building & running...$(COLOR_RESET)"
	docker-compose up --build

.PHONY: stop
stop:
	@echo "$(COLOR_RED)Stopping...$(COLOR_RESET)"
	docker-compose down

.PHONY: unit-test
unit-test:
	@echo "$(COLOR_GREEN)Testing...$(COLOR_RESET)"
	GIN_MODE="release" go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
	@echo "$(COLOR_GREEN)Total Coverage: $$(go tool cover -func coverage.out | grep -E '^total:' | awk '{print $$NF}' | tr -d '%')%$(COLOR_RESET)"

.PHONY: lint
lint:
	@echo "$(COLOR_GREEN)Linting...$(COLOR_RESET)"
	goimports -w .
	golangci-lint run

.PHONY: install-swagger
install-swagger:
	@which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger && go install github.com/go-swagger/go-swagger/cmd/swagger)

.PHONY: swagger
swagger: install-swagger
	@echo "$(COLOR_GREEN)Generating Swagger...$(COLOR_RESET)"
	swagger generate spec -o ./docs/swagger.yaml --scan-models

.PHONY: serve-swagger
serve-swagger: install-swagger
	@echo "$(COLOR_GREEN)Serving Swagger...$(COLOR_RESET)"
	swagger serve -F=swagger ./docs/swagger.yaml

.PHONY: generate-mocks
generate-mocks:
	@echo "$(COLOR_GREEN)Generating Mocks...$(COLOR_RESET)"
	mockery
