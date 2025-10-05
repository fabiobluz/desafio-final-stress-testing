.PHONY: build test run docker-build

build:
	go build -o bin/gload ./cmd/gload

test:
	go test ./... -v

run:
	go run ./cmd/gload --url=http://google.com --requests=1000 --concurrency=10 --format=json

docker-build:
	docker build -t acme/gload:latest .
