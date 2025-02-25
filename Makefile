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
