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

build-image:
	docker build -t murmur-server:latest .
	slim build --target murmur-server:latest --tag murmur-server:slim --continue-after 10
	docker rmi murmur-server:latest
	
build-image-prod:
	docker buildx build --platform linux/amd64 -t murmur-server-prod:latest .
	slim build --target murmur-server:latest --tag murmur-server-prod:slim --continue-after 10
	docker rmi murmur-server-prod:latest
	
start-container:
	docker compose up -d
	
stop-container:
	docker compose down

restart-container: stop-container start-container

test:
	go test ./...
