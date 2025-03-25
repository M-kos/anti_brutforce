ifneq (,$(wildcard ./.env))
    include .env
    export
endif

LOCAL_BIN := $(CURDIR)/bin
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

build: generate
	go build -o ./bin/rl ./cmd/limitter/main.go
	go build -o ./bin/rl-cli ./cmd/limitter-cli/main.go

build2: generate2
	go build -o ./bin/rl ./cmd/limitter/main.go
	go build -o ./bin/rl-cli ./cmd/limitter-cli/main.go

run:
	go run ./cmd/limitter/main.go

test:
	go test ./... -v -race -count 100

integration-test: docker-redis-up
	go test --tags=integration ./... -v

lint:
	golangci-lint run ./... -v 

clean:
	rm -rf ./bin

tidy:
	go mod tidy -v

.deps: export GOBIN := $(LOCAL_BIN)
.deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.proto-generate:
	rm -rf internal/api/pb
	mkdir -p internal/api/pb
	ls -la bin

	$(PROTOC) --proto_path=$(CURDIR) \
		--go_out=$(CURDIR)/internal/api/pb \
		--go-grpc_out=$(CURDIR)/internal/api/pb \
		$(CURDIR)/internal/api/proto/*.proto

	rm bin/protoc-gen-go
	rm bin/protoc-gen-go-grpc

.proto-generate2:
	rm -rf internal/api/pb
	mkdir -p internal/api/pb

	protoc --proto_path=$(CURDIR) \
		--go_out=$(CURDIR)/internal/api/pb \
		--go-grpc_out=$(CURDIR)/internal/api/pb \
		$(CURDIR)/internal/api/proto/*.proto

	rm bin/protoc-gen-go
	rm bin/protoc-gen-go-grpc

generate: .deps .proto-generate tidy

generate2: .proto-generate2 tidy

docker-redis-up:
	REDIS_PASSWORD=$(REDIS_PASSWORD) REDIS_USER=$(REDIS_USER) REDIS_USER_PASSWORD=$(REDIS_USER_PASSWORD) docker-compose -f ./docker/redis/docker-compose.yaml up -d

docker-redis-down:
	docker-compose -f ./docker/redis/docker-compose.yaml down

.PHONY: build run test lint clean tidy generate docker-redis-up docker-redis-down integration-test