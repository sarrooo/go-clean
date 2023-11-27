ENV	?=	dev
PORT ?= 8080

.PHONY: run
run:
	@echo "Running..."
	docker-compose up

.PHONY: build
build:
	@echo "Building & running..."
	docker-compose up --build

.PHONY: stop
stop:
	@echo "Stopping..."
	docker-compose down

.PHONY: unit-test
unit-test:
	@echo "Testing..."
	GIN_MODE="release" go test ./...  -coverprofile=coverage.out
	go tool cover -html=coverage.out
	@echo "Total Coverage : $$(go tool cover -func coverage.out | grep -E '^total:' | awk '{print $$NF}' | tr -d '%')%"

.PHONY: func-test
func-test:
	@echo "Testing..."
	tests/server.sh $(PORT)

.PHONY: lint
lint:
	@echo "Linting..."
	goimports -w .
	golangci-lint run

.PHONY: install-swagger
install-swagger:
	which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger && go install github.com/go-swagger/go-swagger/cmd/swagger)

.PHONY: swagger
swagger: install-swagger
	swagger generate spec -o ./docs/swagger.yaml --scan-models

.PHONY: serve-swagger
serve-swagger: install-swagger
	swagger serve -F=swagger ./docs/swagger.yaml

.PHONY: generate-mocks
generate-mocks:
	mockery
