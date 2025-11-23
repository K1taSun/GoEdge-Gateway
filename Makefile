.PHONY: all build test clean run-infra stop-infra gen-proto

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

gen-proto:
	docker run --rm -v $(PWD):/workspace -w /workspace rvolosatovs/protoc:4.0.0 \
		--proto_path=. --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/gateway.proto
