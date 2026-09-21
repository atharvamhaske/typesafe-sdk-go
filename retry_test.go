package typesafe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetrySucceedsAfterTransientFailures(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&attempts, 1) <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"model":"jev-latest","usage":{},"answers":{}}`))
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(WithAPIKey("k"), WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.SystemOne(context.Background(), "s", nil); err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Errorf("attempts = %d, want 3", got)
	}
}

func TestRetryExhaustsAndReturnsError(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(WithAPIKey("k"), WithBaseURL(srv.URL), WithMaxRetries(2))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.SystemOne(context.Background(), "s", nil)
	apiErr, ok := err.(*Error)
	if !ok || apiErr.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("err = %v, want *Error 503", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 { // initial + 2 retries
		t.Errorf("attempts = %d, want 3", got)
	}
}

func TestRetryHonorsContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(WithAPIKey("k"), WithBaseURL(srv.URL), WithMaxRetries(5))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		_, _ = c.SystemOne(ctx, "s", nil)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("retry loop did not honor context cancellation")
	}
}

func TestRetryDoesNotRetryOnClientError(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"detail":"bad request"}`))
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(WithAPIKey("k"), WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.SystemOne(context.Background(), "s", nil); err == nil {
		t.Fatal("expected error")
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Errorf("attempts = %d, want 1 (4xx should not retry)", got)
	}
}
