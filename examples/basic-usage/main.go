package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}
	ctx := context.Background()
	llm, err := googleai.New(ctx, googleai.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
		googleai.WithDefaultModel("gemini-2.0-flash"))
	if err != nil {
		log.Fatal(err)
	}
	prompt := "What would be a good company name for a company that makes colorful socks? Reply with a company name only, nothing else."
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt, llms.WithTemperature(0.9))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)
}
