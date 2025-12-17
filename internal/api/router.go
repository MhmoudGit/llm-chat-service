package api

import (
	"net/http"

	"chat-service/internal/config"
)

func NewRouter(h *Handler, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	authMw := AuthMiddleware(cfg)
	rateMw := RateLimitMiddleware(cfg)
	chain := func(h http.Handler) http.Handler {
		return rateMw(authMw(h))
	}

	mux.HandleFunc("/health", h.HandleHealth)

	mux.Handle("/chat", chain(http.HandlerFunc(h.HandleChat)))
	mux.Handle("/history", chain(http.HandlerFunc(h.HandleHistory)))

	mux.HandleFunc("/web", h.HandleWeb)

	handler := CORSMiddleware(mux)
	handler = LoggerMiddleware(handler)

	return handler
}
