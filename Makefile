BINARY_NAME=main

build:
	go build -o $(BINARY_NAME) main.go

run:
	clear && go run main.go

clean:
	go clean
	rm -f $(BINARY_NAME)

deps:
	go mod download

run-docker:
	docker-compose up -d

stop-docker:
	docker-compose down

test:
	go test ./...
