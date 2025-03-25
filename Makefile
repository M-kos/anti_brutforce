ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: generate
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

.deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo $$PATH
	@echo $(go env GOPATH)
	export PATH="$$PATH:$(go env GOPATH)/bin"

.proto-generate:
	rm -rf internal/api/pb
	mkdir -p internal/api/pb

	protoc \
		--go_out=internal/api/pb \
		--go-grpc_out=internal/api/pb \
		internal/api/proto/*.proto

generate: .deps .proto-generate tidy

docker-redis-up:
	@echo $(REDIS_USER)
	@echo $(REDIS_USER_PASSWORD)
	@echo $(REDIS_PASSWORD)
	REDIS_PASSWORD=$(REDIS_PASSWORD) REDIS_USER=$(REDIS_USER) REDIS_USER_PASSWORD=$(REDIS_USER_PASSWORD) docker-compose -f ./docker/redis/docker-compose.yaml up -d

docker-redis-down:
	docker-compose -f ./docker/redis/docker-compose.yaml down

.PHONY: build run test lint clean tidy generate docker-redis-up docker-redis-down integration-test