package router

import (
	pb "murmur/go-server/gen/go/inference"
	"murmur/go-server/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, inferenceGrpcClient pb.InferenceServiceClient) {
	api := app.Group("/api/v1")

	api.Post("/transcribe", handlers.TranscribeHandler(inferenceGrpcClient))
}
