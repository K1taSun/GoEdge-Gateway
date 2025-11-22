.PHONY: all build test clean run-infra stop-infra

all: build

build:
	go build -o bin/mqtt-ingestor ./cmd/mqtt-ingestor
	go build -o bin/storage-service ./cmd/storage-service

test:
	go test ./...

clean:
	rm -rf bin/

run-infra:
	docker-compose up -d postgres mqtt-broker

stop-infra:
	docker-compose down

