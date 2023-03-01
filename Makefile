BASE_PATH=$(shell pwd)

test:
	cd ./backend; make test
cov:
	cd ./backend; make cov
lint:
	docker run --rm -v $(BASE_PATH)/backend:/app -w /app golangci/golangci-lint:v1.51.2 golangci-lint run -v --timeout 5m0s
sec:
	docker run --rm -v $(BASE_PATH)/backend:/backend -w /backend securego/gosec ./...
