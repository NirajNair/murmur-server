package main

import (
	"log"
	"murmur/go-server/config"
	pb "murmur/go-server/gen/go/inference"
	"murmur/go-server/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()
	conn, err := grpc.NewClient(cfg.InferenceServiceGRPCUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to Inference Service GRPC server: %v", err)
	}
	defer conn.Close()
	inferenceClient := pb.NewInferenceServiceClient(conn)
	app := fiber.New()
	app.Use(cors.New())
	router.SetupRoutes(app, inferenceClient)
	log.Printf("Fiber API server listening on %s", cfg.GetHTTPAddress())
	log.Fatal(app.Listen(cfg.GetHTTPAddress()))
}
