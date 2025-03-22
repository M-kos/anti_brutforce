build:
	go build -o ./bin/rl ./cmd/main.go

run:
	go run ./cmd/limitter/main.go

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -rf ./bin

tidy:
	go mod tidy -v

.deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	export PATH="$$PATH:$(go env GOPATH)/bin"

.proto-generate:
	rm -rf internal/api/pb
	mkdir -p internal/api/pb

	protoc \
		--go_out=internal/api/pb \
		--go-grpc_out=internal/api/pb \
		internal/api/proto/*.proto

generate: .deps .proto-generate tidy
