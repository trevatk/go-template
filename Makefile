build:
	docker build -t trevatk/go-template:latest .

deps:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run
