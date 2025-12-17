package chat

import (
	"fmt"
	"testing"
)

func TestHistoryManager(t *testing.T) {
	t.Run("Add and Retrieve", func(t *testing.T) {
		h := NewHistoryManager()
		msg := Message{Role: RoleUser, Content: "Hello"}
		h.AddMessage(msg)

		ctx := h.GetContext()
		if len(ctx) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(ctx))
		}
		if ctx[0] != msg {
			t.Errorf("Expected message %v, got %v", msg, ctx[0])
		}
	})

	t.Run("Max History Truncation", func(t *testing.T) {
		h := NewHistoryManager()
		for i := 0; i < 25; i++ {
			h.AddMessage(Message{Role: RoleUser, Content: fmt.Sprintf("msg %d", i)})
		}

		ctx := h.GetContext()
		if len(ctx) != maxHistory {
			t.Errorf("Expected %d messages, got %d", maxHistory, len(ctx))
		}

		// Check if it kept the *last* messages
		lastMsg := ctx[maxHistory-1]
		if lastMsg.Content != "msg 24" {
			t.Errorf("Expected last message 'msg 24', got '%s'", lastMsg.Content)
		}

		firstMsg := ctx[0]
		if firstMsg.Content != "msg 5" { // 0..24 is 25 items -> remove first 5 -> start at 5
			t.Errorf("Expected first message 'msg 5', got '%s'", firstMsg.Content)
		}
	})
}
