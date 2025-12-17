package integration

import (
	"os"
	"testing"
	"time"

	"chat-service/internal/chat"
	"chat-service/internal/config"
	"chat-service/internal/llm"
)

func TestGroqIntegration(t *testing.T) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: GROQ_API_KEY not set")
	}

	cfg := &config.Config{
		GroqAPIKey: apiKey,
		AppModel:   "llama-3.3-70b-versatile",
		MaxTokens:  50,
	}

	client := llm.NewClient(cfg)

	messages := []chat.Message{
		{Role: chat.RoleUser, Content: "Say hello in one word."},
	}

	stream, err := client.StreamChat(messages)
	if err != nil {
		t.Fatalf("Failed to call Groq API: %v", err)
	}

	var response string
	done := make(chan bool)
	go func() {
		for chunk := range stream {
			response += chunk
		}
		done <- true
	}()

	select {
	case <-done:
		if response == "" {
			t.Error("Received empty response from Groq API")
		}
		t.Logf("Groq Response: %s", response)
	case <-time.After(10 * time.Second):
		t.Fatal("Timeout waiting for Groq API response")
	}
}
