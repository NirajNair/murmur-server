# Builder Stage
FROM golang:1.24.2-alpine AS builder
RUN apk add --no-cache git ca-certificates protobuf-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN go install \
    google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="/go/bin:${PATH}"
COPY /protos ./protos
RUN protoc --go_out=. --go_opt=module=murmur/go-server \
           --go-grpc_out=. --go-grpc_opt=module=murmur/go-server \
           protos/inference.proto

RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix static -o main .

# Final Stage
FROM alpine:3.19
RUN apk add --no-cache git ca-certificates
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/main .
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]