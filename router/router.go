package router

import (
	pb "murmur/go-server/gen/go/inference"
	"murmur/go-server/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, inferenceClient pb.InferenceServiceClient) {
	api := app.Group("/api")

	api.Post("/transcribe", handlers.TranscribeHandler(inferenceClient))
}
