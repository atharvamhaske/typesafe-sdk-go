package typesafe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Error is a non-2xx API response.
type Error struct {
	StatusCode int
	Detail     json.RawMessage // HTTPValidationError detail, if any
	Body       []byte
}

func (e *Error) Error() string {
	if len(e.Detail) > 0 {
		return fmt.Sprintf("typesafe: HTTP %d: %s", e.StatusCode, e.Detail)
	}
	return fmt.Sprintf("typesafe: HTTP %d: %s", e.StatusCode, e.Body)
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("typesafe: marshal request: %w", err)
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, rd)
	if err != nil {
		return fmt.Errorf("typesafe: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("typesafe: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("typesafe: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &Error{StatusCode: resp.StatusCode, Body: data}
		var ve struct {
			Detail json.RawMessage `json:"detail"`
		}
		if json.Unmarshal(data, &ve) == nil {
			apiErr.Detail = ve.Detail
		}
		return apiErr
	}
	if out != nil {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("typesafe: decode response: %w", err)
		}
	}
	return nil
}
