lint:
	golangci-lint cache clean
	golangci-lint run -c .golangci.yml

build:
	go build -o task-manager cmd/main.go

run:
	go run cmd/main.go
