
build:
	docker build -t trevatk/go-template:latest -f ./docker/Dockerfile .

deps:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run