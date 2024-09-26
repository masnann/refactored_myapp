# Targets
test-unit-repository:
	go test ./test/unit/repository/... -coverpkg=./repository/... -coverprofile=coverage_repository.out
	go tool cover -html=coverage_repository.out -o coverage_repository.html

test-unit-service:
	go test ./test/unit/service/... -coverpkg=./service/... -coverprofile=coverage_service.out
	go tool cover -html=coverage_service.out -o coverage_service.html

test-integrations:
	go test ./test/integrations/... -coverpkg=./handler/...,./service/...,./repository/... -coverprofile=coverage_integrations.out
	go tool cover -html=coverage_integrations.out -o coverage_integrations.html

test-all: 
	go test ./... -coverpkg=./handler/...,./service/...,./repository/... -coverprofile=coverage_all.out
	go tool cover -html=coverage_all.out -o coverage_all.html

run:
	go run main.go

migrations:
	goose -dir db/migrations postgres "postgres://postgres:mkpmobile2024@localhost:5432/goauth?sslmode=disable" up