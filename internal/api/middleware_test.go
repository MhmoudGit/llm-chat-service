package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"chat-service/internal/config"
)

func TestAuthMiddleware(t *testing.T) {
	cfg := &config.Config{
		APIKey: "secret",
	}

	middleware := AuthMiddleware(cfg)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name           string
		headers        map[string]string
		expectedStatus int
	}{
		{
			name:           "No Key",
			headers:        map[string]string{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong Key",
			headers:        map[string]string{"X-API-Key": "wrong"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid Key Header",
			headers:        map[string]string{"X-API-Key": "secret"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid Key Bearer",
			headers:        map[string]string{"Authorization": "Bearer secret"},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	cfg := &config.Config{
		RateLimitRPS:   10,
		RateLimitBurst: 2,
	}

	middleware := RateLimitMiddleware(cfg)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// 1st request - OK
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("1st req: expected 200, got %d", rr.Code)
	}

	// 2nd request - OK
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("2nd req: expected 200, got %d", rr.Code)
	}

	// 3rd request - Should fail (Burst is 2) - Assuming immediate consecutive calls exceed rate
	// Wait, rate limiter allows N events.
	// If Burst is 2, we can do 2 immediately.
	// If RPS is 10, that's 1 token every 100ms.
	// We might need to send more than burst to trigger 429 safely.

	// Let's force excessive requests
	blocked := false
	for i := 0; i < 5; i++ {
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code == http.StatusTooManyRequests {
			blocked = true
			break
		}
	}

	if !blocked {
		t.Errorf("Rate limiter did not block excessive requests")
	}
}
