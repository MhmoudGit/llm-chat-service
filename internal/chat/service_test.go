package chat

import (
	"testing"
)

// MockLLM
type MockLLM struct {
	CapturedMessages []Message
	ResponseChunks   []string
	Err              error
}

func (m *MockLLM) StreamChat(messages []Message) (<-chan string, error) {
	m.CapturedMessages = messages
	if m.Err != nil {
		return nil, m.Err
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, chunk := range m.ResponseChunks {
			ch <- chunk
		}
	}()
	return ch, nil
}

func TestService_ProcessMessage(t *testing.T) {
	h := NewHistoryManager()
	mockLLM := &MockLLM{
		ResponseChunks: []string{"Hello", " ", "World"},
	}
	s := NewService(h, mockLLM)

	userContent := "Hi there"
	// Process
	stream, err := s.ProcessMessage(userContent)
	if err != nil {
		t.Fatalf("ProcessMessage failed: %v", err)
	}

	// Consume stream
	var fullResponse string
	for chunk := range stream {
		fullResponse += chunk
	}

	if fullResponse != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", fullResponse)
	}

	ctx := h.GetAll()
	if len(ctx) < 1 { // Actually should be 2 now (User + Assistant)
		t.Fatalf("History empty")
	}
	if ctx[0].Content != userContent {
		t.Errorf("First message should be user content")
	}

	if len(mockLLM.CapturedMessages) != 1 {
		t.Errorf("LLM called with wrong number of messages: %d", len(mockLLM.CapturedMessages))
	}

	ctx = h.GetAll()
	if len(ctx) != 2 {
		t.Errorf("Expected 2 messages in history, got %d", len(ctx))
	}
	if ctx[1].Role != RoleAssistant {
		t.Errorf("Expected 2nd message role to be assistant, got %s", ctx[1].Role)
	}
	if ctx[1].Content != "Hello World" {
		t.Errorf("Expected 2nd message content 'Hello World', got '%s'", ctx[1].Content)
	}
}
