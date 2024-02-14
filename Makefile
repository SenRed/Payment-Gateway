.PHONY: start s test t lint unit_test ut integration_test it generate_doc generate_mock

# If exist, load environment config for .env file, otherwise from OS variables env
ifeq ($(shell test -s .env && echo -n yes),yes)
	include .env
	export $(shell sed 's/=.*//' .env)
endif

start s:
	docker compose up --build --remove-orphans

lint:
	golangci-lint run

test t: unit_test integration_test

ut unit_test:
	go test -v $(shell go list ./internal/... | grep -v e2e)

it integration_test:
	docker compose -f docker-compose.test.yml up --abort-on-container-exit --build --remove-orphans

generate_doc:
	swag init

generate_mock:
	@mockgen -source=internal/domain/port/repository.go -destination=internal/infrastructure/mock/mock_repository.go -package=mock
	@mockgen -source=internal/domain/port/paymentprocessor.go -destination=internal/infrastructure/mock/mock_paymentprocessor.go -package=mock