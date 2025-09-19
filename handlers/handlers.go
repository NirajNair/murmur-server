package handlers

import (
	"bytes"
	"io"
	"log"
	pb "murmur/go-server/gen/go/inference"

	"github.com/gofiber/fiber/v2"
)

func TranscribeHandler(grpcClient pb.InferenceServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		grpcStream, err := grpcClient.TranscribeAndFix(c.Context())
		if err != nil {
			log.Printf("Failed to start gRPC stream: %v", err)
			return internalServerError(c)
		}
		buffer := make([]byte, 32*1024)
		audioReader := bytes.NewReader(c.Body())
		for {
			n, err := audioReader.Read(buffer)
			if n > 0 {
				if err := grpcStream.Send(&pb.AudioChunk{AudioBytes: buffer[:n]}); err != nil {
					log.Printf("Error sending audio chunk: %v", err)
					return internalServerError(c)
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading audio data: %v", err)
				return internalServerError(c)
			}
		}
		response, err := grpcStream.CloseAndRecv()
		if err != nil {
			log.Printf("Error receiving response from gRPC stream: %v", err)
			return internalServerError(c)
		}
		return c.JSON(fiber.Map{"transcription": response.GetText()})
	}
}

func internalServerError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "It's not you, it's us. We are facing some issue, please try again later."})
}
