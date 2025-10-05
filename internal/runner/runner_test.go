package runner

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Testa se o runner respeita total de requests e mede distribuição de status.
func TestRunBasic(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// alterna status conforme path
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	cfg := Config{
		URL:         srv.URL + "/ok",
		Requests:    50,
		Concurrency: 5,
		Timeout:     5 * time.Second,
		Method:      http.MethodGet,
		Headers:     make(http.Header),
		Format:      "text",
	}

	sum, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if sum.TotalRequests != 50 {
		t.Fatalf("TotalRequests=%d want=50", sum.TotalRequests)
	}
	if sum.StatusDist[200] == 0 {
		t.Fatalf("expected some 200 responses, got dist=%v", sum.StatusDist)
	}

	// Agora testa distribuição com outro path
	cfg.URL = srv.URL + "/notfound"
	sum, err = Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if sum.StatusDist[404] == 0 {
		t.Fatalf("expected 404 responses, got dist=%v", sum.StatusDist)
	}
}

// Testa latência registrando um pequeno delay no servidor.
func TestRunLatency(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := Config{
		URL:         srv.URL,
		Requests:    10,
		Concurrency: 2,
		Timeout:     2 * time.Second,
		Method:      http.MethodGet,
		Headers:     make(http.Header),
		Format:      "text",
	}

	sum, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if sum.TotalRequests != 10 {
		t.Fatalf("TotalRequests=%d want=10", sum.TotalRequests)
	}
	// Espera alguma média de latência > 0
	if sum.AvgLatency <= 0 {
		t.Fatalf("expected AvgLatency > 0, got %v", sum.AvgLatency)
	}
}
