package typesafe

import (
	"context"
	"os"
	"testing"
)

// TestLive hits the real API. It is skipped unless TYPESAFE_API_KEY is set.
func TestLive(t *testing.T) {
	if os.Getenv("TYPESAFE_API_KEY") == "" {
		t.Skip("TYPESAFE_API_KEY not set")
	}
	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	models, err := c.ListModels(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(models) == 0 {
		t.Error("no models returned")
	}

	resp, err := c.SystemOne(ctx, "I was charged twice. Please refund the duplicate charge today.",
		map[string]Question{
			"category": Choice{Instructions: "Categorize the message", Criteria: map[string]string{
				"billing": "Billing issue", "technical": "Technical issue",
			}},
			"urgency": Score{Instructions: "Rate urgency", Criteria: []string{"Can wait", "Needs attention"}},
			"is_dupe": Noul{Instructions: "Is this a duplicate charge?", Criteria: map[string]string{
				"true": "Duplicate", "false": "Not a duplicate",
			}},
		})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := resp.Choices()["category"]; !ok {
		t.Error("missing choice answer for category")
	}
	if _, ok := resp.Scores()["urgency"]; !ok {
		t.Error("missing score answer for urgency")
	}
	if _, ok := resp.Nouls()["is_dupe"]; !ok {
		t.Error("missing noul answer for is_dupe")
	}
	if resp.Usage.InputTokens == 0 {
		t.Error("usage.input_tokens is zero")
	}
}
