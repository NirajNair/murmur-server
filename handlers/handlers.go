package handlers

import (
	"io"
	"log"
	pb "murmur/go-server/gen/go/inference"

	"github.com/gofiber/fiber/v2"
)

func TranscribeHandler(grpcClient pb.InferenceServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodPost {
			return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"error": "Only POST method is allowed"})
		}

		log.Println("Transcribe endpoint hit")

		grpcStream, err := grpcClient.TranscribeStream(c.Context())
		if err != nil {
			log.Printf("Failed to start gRPC stream: %v", err)
			return internalServerError(c)
		}

		buffer := make([]byte, 32*1024)

		fileHeader, err := c.FormFile("audio")
		if err != nil {
			log.Printf("Error getting audio file from form: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "audio file not found in form"})
		}

		audioFile, err := fileHeader.Open()
		if err != nil {
			log.Printf("Error opening audio file: %v", err)
			return internalServerError(c)
		}
		defer audioFile.Close()

		for {
			n, err := audioFile.Read(buffer)
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
				log.Printf("Error reading audio file: %v", err)
				return internalServerError(c)
			}
		}

		response, err := grpcStream.CloseAndRecv()
		if err != nil {
			log.Printf("Error receiving response from gRPC stream: %v", err)
			return internalServerError(c)
		}

		log.Printf("Response: %v", response.GetText())

		return c.JSON(fiber.Map{"transcribed_text": response.GetText()})
	}
}

func internalServerError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "It's not you, it's us. We are facing some issue, please try again later."})
}
