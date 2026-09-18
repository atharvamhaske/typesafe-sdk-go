package typesafe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const systemOneFixture = `{
  "model": "jev-latest",
  "usage": {"input_tokens": 120, "output_tokens": 12},
  "answers": {
    "category": {"type": "choice", "choice": "billing", "confidence": 0.9, "probabilities": {"billing": 0.8, "technical": 0.1}},
    "urgency":  {"type": "score", "score": 1.7, "confidence": 0.9, "legend": {"0": "Can wait"}, "probabilities": {"0": 0.1, "1": 0.9}},
    "is_dupe":  {"type": "noul", "noul": 0.98}
  }
}`

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c, err := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSystemOne(t *testing.T) {
	var gotBody map[string]json.RawMessage
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/systemone" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization = %q", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(systemOneFixture))
	})

	resp, err := c.SystemOne(context.Background(), "I was charged twice.", map[string]Question{
		"category": Choice{Instructions: "categorize", Criteria: map[string]string{"billing": "b", "technical": "t"}},
		"urgency":  Score{Instructions: "rate", Criteria: []string{"Can wait", "Needs attention"}},
		"is_dupe":  Noul{Instructions: "dupe?", Criteria: map[string]string{"true": "y", "false": "n"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	var sent struct {
		State     string                    `json:"state"`
		Model     string                    `json:"model"`
		Questions map[string]map[string]any `json:"questions"`
	}
	full, _ := json.Marshal(gotBody)
	if err := json.Unmarshal(full, &sent); err != nil {
		t.Fatal(err)
	}
	if sent.State != "I was charged twice." || sent.Model != "jev-latest" {
		t.Errorf("sent state=%q model=%q", sent.State, sent.Model)
	}
	for key, want := range map[string]string{"category": "choice", "urgency": "score", "is_dupe": "noul"} {
		if got := sent.Questions[key]["type"]; got != want {
			t.Errorf("question %s type = %v, want %s", key, got, want)
		}
	}

	if resp.Usage.InputTokens != 120 || resp.Usage.OutputTokens != 12 {
		t.Errorf("usage = %+v", resp.Usage)
	}
	if got := resp.Choices()["category"]; got.Choice != "billing" || got.Confidence != 0.9 || got.Probabilities["billing"] != 0.8 {
		t.Errorf("category = %+v", got)
	}
	if got := resp.Scores()["urgency"]; got.Score != 1.7 || got.Legend["0"] != "Can wait" {
		t.Errorf("urgency = %+v", got)
	}
	if got := resp.Nouls()["is_dupe"]; got.Noul != 0.98 {
		t.Errorf("is_dupe = %+v", got)
	}
}

func TestSystemOneModelOverride(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Model string `json:"model"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Model != "jev-2" {
			t.Errorf("model = %q, want jev-2", body.Model)
		}
		w.Write([]byte(`{"model":"jev-2","usage":{},"answers":{}}`))
	})
	if _, err := c.SystemOne(context.Background(), "s", nil, WithRequestModel("jev-2")); err != nil {
		t.Fatal(err)
	}
}

func TestListModels(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/models" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Write([]byte(`{"models":[{"name":"jev-latest","description":"d","release_date":"2026-01-01"}]}`))
	})
	models, err := c.ListModels(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(models) != 1 || models[0].Name != "jev-latest" {
		t.Errorf("models = %+v", models)
	}
}

func TestValidationError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		w.Write([]byte(`{"detail":[{"loc":["body","state"],"msg":"field required","type":"missing"}]}`))
	})
	_, err := c.SystemOne(context.Background(), "", nil)
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("err = %v, want *Error", err)
	}
	if apiErr.StatusCode != 422 || len(apiErr.Detail) == 0 {
		t.Errorf("apiErr = %+v", apiErr)
	}
}

func TestNewClientMissingKey(t *testing.T) {
	t.Setenv("TYPESAFE_API_KEY", "")
	if _, err := NewClient(); err == nil {
		t.Error("expected error for missing API key")
	}
}
