package typesafe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Usage reports token consumption for a request.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Answer is one of ChoiceAnswer, ScoreAnswer, or NoulAnswer.
type Answer interface {
	answerType() string
}

// ChoiceAnswer answers a Choice question.
type ChoiceAnswer struct {
	Choice        string             `json:"choice"`
	Confidence    float64            `json:"confidence"`
	Probabilities map[string]float64 `json:"probabilities"`
}

func (ChoiceAnswer) answerType() string { return "choice" }

// ScoreAnswer answers a Score question.
type ScoreAnswer struct {
	Score         float64            `json:"score"`
	Confidence    float64            `json:"confidence"`
	Legend        map[string]string  `json:"legend"`
	Probabilities map[string]float64 `json:"probabilities"`
}

func (ScoreAnswer) answerType() string { return "score" }

// NoulAnswer answers a Noul question.
type NoulAnswer struct {
	Noul float64 `json:"noul"`
}

func (NoulAnswer) answerType() string { return "noul" }

// SystemOneResponse is the result of a SystemOne call.
type SystemOneResponse struct {
	Model   string            `json:"model"`
	Usage   Usage             `json:"usage"`
	Answers map[string]Answer `json:"-"`
}

func (r *SystemOneResponse) UnmarshalJSON(data []byte) error {
	var wire struct {
		Model   string                     `json:"model"`
		Usage   Usage                      `json:"usage"`
		Answers map[string]json.RawMessage `json:"answers"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	r.Model = wire.Model
	r.Usage = wire.Usage
	r.Answers = make(map[string]Answer, len(wire.Answers))
	for key, raw := range wire.Answers {
		var head struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &head); err != nil {
			return fmt.Errorf("answer %q: %w", key, err)
		}
		var (
			ans Answer
			err error
		)
		switch head.Type {
		case "choice":
			var a ChoiceAnswer
			err, ans = json.Unmarshal(raw, &a), a
		case "score":
			var a ScoreAnswer
			err, ans = json.Unmarshal(raw, &a), a
		case "noul":
			var a NoulAnswer
			err, ans = json.Unmarshal(raw, &a), a
		default:
			return fmt.Errorf("answer %q: unknown type %q", key, head.Type)
		}
		if err != nil {
			return fmt.Errorf("answer %q: %w", key, err)
		}
		r.Answers[key] = ans
	}
	return nil
}

// Choices returns the choice-typed answers by question key.
func (r *SystemOneResponse) Choices() map[string]ChoiceAnswer {
	out := make(map[string]ChoiceAnswer)
	for k, a := range r.Answers {
		if c, ok := a.(ChoiceAnswer); ok {
			out[k] = c
		}
	}
	return out
}

// Scores returns the score-typed answers by question key.
func (r *SystemOneResponse) Scores() map[string]ScoreAnswer {
	out := make(map[string]ScoreAnswer)
	for k, a := range r.Answers {
		if s, ok := a.(ScoreAnswer); ok {
			out[k] = s
		}
	}
	return out
}

// Nouls returns the noul-typed answers by question key.
func (r *SystemOneResponse) Nouls() map[string]NoulAnswer {
	out := make(map[string]NoulAnswer)
	for k, a := range r.Answers {
		if n, ok := a.(NoulAnswer); ok {
			out[k] = n
		}
	}
	return out
}

// SystemOneOption configures a single SystemOne call.
type SystemOneOption func(*systemOneRequest)

// WithRequestModel overrides the client's default model for this call.
func WithRequestModel(model string) SystemOneOption {
	return func(r *systemOneRequest) { r.Model = model }
}

type systemOneRequest struct {
	State     string              `json:"state"`
	Model     string              `json:"model"`
	Questions map[string]Question `json:"questions"`
}

// SystemOne asks questions about a state via POST /v1/systemone. When the
// client was built with WithCache, an identical request within the cache's
// TTL is served without a network call.
func (c *Client) SystemOne(ctx context.Context, state string, questions map[string]Question, opts ...SystemOneOption) (*SystemOneResponse, error) {
	req := systemOneRequest{State: state, Model: c.model, Questions: questions}
	for _, opt := range opts {
		opt(&req)
	}

	var key string
	if c.cache != nil {
		body, err := json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("typesafe: marshal request: %w", err)
		}
		key = cacheKey(body)
		if cached, ok := c.cache.get(key); ok {
			var resp SystemOneResponse
			if err := json.Unmarshal(cached, &resp); err != nil {
				return nil, fmt.Errorf("typesafe: decode cached response: %w", err)
			}
			return &resp, nil
		}
	}

	data, err := c.doRaw(ctx, http.MethodPost, "/v1/systemone", req)
	if err != nil {
		return nil, err
	}
	var resp SystemOneResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("typesafe: decode response: %w", err)
	}
	if c.cache != nil {
		c.cache.set(key, data)
	}
	return &resp, nil
}
