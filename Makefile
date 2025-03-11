build:
	go build -o ./bin/rl ./cmd/main.go

run:
	go run ./cmd/main.go

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -rf ./bin

tidy:
	go mod tidy -v

generate:
	rm -rf internal/api/pb
	mkdir -p internal/api/pb

	protoc \
		--go_out=internal/api/pb \
		--go-grpc_out=internal/api/pb \
		internal/api/proto/*.proto
		