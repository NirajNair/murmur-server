package main

import (
	"log"
	pb "murmur/go-server/gen/go/inference"
	"murmur/go-server/router"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gRPCServerAddress = "localhost:50051"
)

func main() {
	conn, err := grpc.NewClient(gRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	inferenceClient := pb.NewInferenceServiceClient(conn)
	app := fiber.New()
	router.SetupRoutes(app, inferenceClient)
	log.Printf("Fiber API server listening on :8080")
	log.Fatal(app.Listen(":8080"))
}
