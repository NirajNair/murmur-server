package router

import (
	pb "murmur/go-server/gen/go/inference"
	"murmur/go-server/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, inferenceGrpcClient pb.InferenceServiceClient) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	api := app.Group("/api")
	api.Post("/transcribe", handlers.TranscribeHandler(inferenceGrpcClient))
}
