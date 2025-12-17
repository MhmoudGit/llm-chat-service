package api

import (
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.HandleHealth)
	mux.HandleFunc("/chat", h.HandleChat)

	handler := CORSMiddleware(mux)
	handler = LoggerMiddleware(handler)

	return handler
}
