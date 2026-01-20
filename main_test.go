package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	lengths := []int{4, 6, 8, 10}
	for _, l := range lengths {
		slug := generateSlug(l)
		if len(slug) != l {
			t.Errorf("Expected slug length %d, got %d", l, len(slug))
		}
	}
}

func TestRateLimit(t *testing.T) {
	s := NewServer(nil, "")
	handler := s.rateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, 1, 2) // 1 request per second, burst of 2

	server := httptest.NewServer(handler)
	defer server.Close()

	// First two requests should pass (burst = 2)
	for i := 0; i < 2; i++ {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Request %d: expected 200, got %d", i, resp.StatusCode)
		}
	}

	// Third request should be rate limited
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected 429 Too Many Requests, got %d", resp.StatusCode)
	}
}
