package chat

import (
	"fmt"
	"strings"
)

// LLMClient interface to decouple from concrete implementation (useful for testing)
type LLMClient interface {
	StreamChat(messages []Message) (<-chan string, error)
}

type Service struct {
	history *HistoryManager
	llm     LLMClient
}

func NewService(h *HistoryManager, llm LLMClient) *Service {
	return &Service{
		history: h,
		llm:     llm,
	}
}

// ProcessMessage handles a new user message, updates history, and streams the response.
// It returns a channel that emits chunks of the assistant's response.
func (s *Service) ProcessMessage(userContent string) (<-chan string, error) {
	userMsg := Message{Role: RoleUser, Content: userContent}
	s.history.AddMessage(userMsg)

	messages := s.history.GetContext()

	stream, err := s.llm.StreamChat(messages)
	if err != nil {
		return nil, fmt.Errorf("llm call failed: %w", err)
	}

	outChan := make(chan string)

	go func() {
		defer close(outChan)
		var sb strings.Builder

		for chunk := range stream {
			sb.WriteString(chunk)
			outChan <- chunk
		}

		fullResponse := sb.String()
		if fullResponse != "" {
			s.history.AddMessage(Message{
				Role:    RoleAssistant,
				Content: fullResponse,
			})
		}
	}()

	return outChan, nil
}

func (s *Service) GetHistory() []Message {
	return s.history.GetAll()
}
