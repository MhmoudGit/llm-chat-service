package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"chat-service/internal/chat"
	"chat-service/internal/config"
)

const groqAPIURL = "https://api.groq.com/openai/v1/chat/completions"

type Client struct {
	apiKey    string
	model     string
	maxTokens int
	client    *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		apiKey:    cfg.GroqAPIKey,
		model:     cfg.AppModel,
		maxTokens: cfg.MaxTokens,
		client:    &http.Client{},
	}
}

type groqRequest struct {
	Model     string         `json:"model"`
	Messages  []chat.Message `json:"messages"`
	MaxTokens int            `json:"max_tokens,omitempty"`
	Stream    bool           `json:"stream"`
}

type groqStreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// StreamChat sends messages to Groq and returns a channel of streamed content.
func (c *Client) StreamChat(messages []chat.Message) (<-chan string, error) {
	reqBody := groqRequest{
		Model:     c.model,
		Messages:  messages,
		MaxTokens: c.maxTokens,
		Stream:    true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", groqAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("groq api error (status %d): %s", resp.StatusCode, string(body))
	}

	streamChan := make(chan string)

	go func() {
		defer resp.Body.Close()
		defer close(streamChan)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if strings.TrimSpace(data) == "[DONE]" {
				return
			}

			var streamResp groqStreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}

			if len(streamResp.Choices) > 0 {
				content := streamResp.Choices[0].Delta.Content
				if content != "" {
					streamChan <- content
				}
			}
		}
	}()

	return streamChan, nil
}
