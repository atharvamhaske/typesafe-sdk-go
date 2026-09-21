package typesafe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheServesRepeatedRequestWithoutNetworkCall(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Write([]byte(`{"model":"jev-latest","usage":{"input_tokens":1,"output_tokens":1},"answers":{}}`))
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(WithAPIKey("k"), WithBaseURL(srv.URL), WithCache(10, time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	for range 3 {
		if _, err := c.SystemOne(context.Background(), "same state", nil); err != nil {
			t.Fatal(err)
		}
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("network calls = %d, want 1 (cache should serve the other 2)", got)
	}
}

func TestCacheMissesOnDifferentState(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Write([]byte(`{"model":"jev-latest","usage":{},"answers":{}}`))
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(WithAPIKey("k"), WithBaseURL(srv.URL), WithCache(10, time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.SystemOne(context.Background(), "state a", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := c.SystemOne(context.Background(), "state b", nil); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("network calls = %d, want 2 (different requests should not share a cache entry)", got)
	}
}

func TestCacheExpiresAfterTTL(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Write([]byte(`{"model":"jev-latest","usage":{},"answers":{}}`))
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(WithAPIKey("k"), WithBaseURL(srv.URL), WithCache(10, 10*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.SystemOne(context.Background(), "s", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
	if _, err := c.SystemOne(context.Background(), "s", nil); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("network calls = %d, want 2 (entry should have expired)", got)
	}
}

func TestNoCacheByDefault(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Write([]byte(`{"model":"jev-latest","usage":{},"answers":{}}`))
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(WithAPIKey("k"), WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if _, err := c.SystemOne(context.Background(), "s", nil); err != nil {
			t.Fatal(err)
		}
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("network calls = %d, want 2 (caching must be opt-in)", got)
	}
}
