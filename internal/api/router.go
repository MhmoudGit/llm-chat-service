package api

import (
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.HandleHealth)
	mux.HandleFunc("/chat", h.HandleChat)
	mux.HandleFunc("/history", h.HandleHistory)
	mux.HandleFunc("/web", h.HandleWeb)

	handler := CORSMiddleware(mux)
	handler = LoggerMiddleware(handler)

	return handler
}
