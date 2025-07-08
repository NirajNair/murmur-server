package handlers

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	pb "murmur/go-server/gen/go/inference"

	"github.com/gofiber/fiber/v2"
)

type APIRequest struct {
	Audio string `json:"audio"`
}

type APIResponse struct {
	RawText   string `json:"raw_text"`
	FinalText string `json:"final_text"`
}

func TranscribeHandler(client pb.InferenceServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodPost {
			return c.Status(fiber.StatusMethodNotAllowed).SendString("Only POST method is allowed")
		}

		var req APIRequest
		if err := c.BodyParser(&req); err != nil {
			log.Printf("Error parsing request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		audioBytes, err := base64.StdEncoding.DecodeString(req.Audio)
		if err != nil {
			log.Printf("Error decoding base64 audio: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid base64 audio data"})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
		defer cancel()

		gRPCReq := &pb.AudioRequest{
			AudioPcm: audioBytes,
		}

		gRPCRes, err := client.TranscribeAndFix(ctx, gRPCReq)
		if err != nil {
			log.Printf("gRPC call failed: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to perform inference"})
		}

		apiRes := APIResponse{
			RawText:   gRPCRes.GetText(),
			FinalText: gRPCRes.GetText(),
		}

		return c.Status(fiber.StatusOK).JSON(apiRes)
	}
}
