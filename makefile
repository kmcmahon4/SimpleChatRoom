run:
	go run ./cmd/main.go

build:
	go build -o /bin/chatroom cmd/main.go

deps:
	go mod download

fmt:
	go fmt ./...

test:
	go test ./...