mock:
	mockgen -source=./pkg/log/logger.go -destination=./pkg/log/mocks/logger.go
	mockgen -source=./internal/batch/repository.go -destination=./internal/batch/mocks/repository.go
	mockgen -source=./internal/batch/service.go -destination=./internal/batch/mocks/service.go
	mockgen -source=./internal/transaction/service.go -destination=./internal/transaction/mocks/service.go
	mockgen -source=./internal/transaction/repository.go -destination=./internal/transaction/mocks/repository.go
build:
	go build -mod=vendor -o bin/api ./cmd/api/main.go
fmt:
	go fmt ./...
test: fmt
	go test -mod=vendor -race -cover -coverprofile=cover.out `go list ./... | grep -v mocks`
	go tool cover -func cover.out | grep total:
cov: test
	go tool cover -html=cover.out
start-api: build
	export $$(cat .env | grep -v ^\# | xargs) && ./bin/api
run-api:
	export $$(cat .env | grep -v ^\# | xargs) && go run ./cmd/api/main.go
start-worker: build
	export $$(cat .env | grep -v ^\# | xargs) && ./bin/worker
run-worker:
	export $$(cat .env | grep -v ^\# | xargs) && go run ./cmd/worker/main.go
