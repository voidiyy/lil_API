SERVICE_NAME = http_server

build:
	go build -o $(SERVICE_NAME) main.go


run: build
	go run ./$(SERVIVE_NAME)

test:
	go test ./... -v
