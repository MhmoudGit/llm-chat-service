package chat

import (
	"sync"
)

const maxHistory = 20

type HistoryManager struct {
	mu       sync.RWMutex
	messages []Message
}

func NewHistoryManager() *HistoryManager {
	return &HistoryManager{
		messages: make([]Message, 0),
	}
}

func (h *HistoryManager) AddMessage(msg Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.messages = append(h.messages, msg)
}

func (h *HistoryManager) GetContext() []Message {
	h.mu.RLock()
	defer h.mu.RUnlock()

	total := len(h.messages)
	start := 0
	if total > maxHistory {
		start = total - maxHistory
	}

	result := make([]Message, total-start)
	copy(result, h.messages[start:])
	return result
}

func (h *HistoryManager) GetAll() []Message {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]Message, len(h.messages))
	copy(result, h.messages)
	return result
}
