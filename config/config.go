package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost              string
	ServerPort              string
	InferenceServiceGRPCUrl string
	AudioBufferSizeInKb     int
	Environment             string
}

var (
	instance *Config
	once     sync.Once
)

func GetInstance() *Config {
	once.Do(func() {
		instance = load()
	})
	return instance
}

func Load() *Config {
	return GetInstance()
}

func load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}
	config := &Config{
		ServerHost:              getEnv("SERVER_HOST"),
		ServerPort:              getEnv("SERVER_PORT"),
		InferenceServiceGRPCUrl: getEnv("INFERENCE_SERVICE_GRPC_URL"),
		AudioBufferSizeInKb:     getEnvAsInt("AUDIO_BUFFER_SIZE_IN_KB"),
		Environment:             getEnv("ENVIRONMENT"),
	}

	return config
}

func (c *Config) GetHTTPAddress() string {
	return c.ServerHost + ":" + c.ServerPort
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set or is empty", key)
	}
	return value
}

func getEnvAsInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set or is empty", key)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid integer value for %s: %s", key, value)
	}
	return intValue
}
