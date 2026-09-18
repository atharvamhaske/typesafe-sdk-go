package typesafe

import (
	"context"
	"net/http"
)

// Model describes an available TypeSafe model.
type Model struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
}

// ListModels returns the available models via GET /v1/models.
func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
	var resp struct {
		Models []Model `json:"models"`
	}
	if err := c.do(ctx, http.MethodGet, "/v1/models", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Models, nil
}
