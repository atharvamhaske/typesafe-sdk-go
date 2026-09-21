// This example shows how to manually log typesafe-sdk-go calls to Braintrust
// as llm spans, without a dedicated trace/contrib/typesafe middleware.
//
// Field reference: https://www.braintrust.dev/docs/integrations/opentelemetry#manual-tracing
package main

import (
	"context"
	"encoding/json"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"

	typesafe "github.com/atharvamhaske/typesafe-sdk-go"
	"github.com/braintrustdata/braintrust-sdk-go"
)

func main() {
	tp := trace.NewTracerProvider()
	defer tp.Shutdown(context.Background()) //nolint:errcheck
	otel.SetTracerProvider(tp)

	bt, err := braintrust.New(tp,
		braintrust.WithProject("typesafe-sdk-go-examples"),
		braintrust.WithBlockingLogin(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	_ = bt

	client, err := typesafe.NewClient() // reads TYPESAFE_API_KEY
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	state := "I was charged twice. Please refund the duplicate charge today."
	questions := map[string]typesafe.Question{
		"category": typesafe.Choice{
			Instructions: "Categorize the message",
			Criteria:     map[string]string{"billing": "Billing issue", "technical": "Technical issue"},
		},
		"urgency": typesafe.Score{
			Instructions: "Rate urgency",
			Criteria:     []string{"Can wait", "Needs attention"},
		},
	}

	tracer := otel.Tracer("typesafe-sdk-go")
	spanCtx, span := tracer.Start(ctx, "typesafe.systemone")
	defer span.End()

	setJSONAttr(span, "braintrust.input_json", map[string]any{"state": state, "questions": questions})
	setJSONAttr(span, "braintrust.span_attributes", map[string]string{"type": "llm"})

	resp, err := client.SystemOne(spanCtx, state, questions)
	if err != nil {
		span.RecordError(err)
		log.Fatal(err)
	}

	setJSONAttr(span, "braintrust.output_json", resp.Answers)
	setJSONAttr(span, "braintrust.metadata", map[string]any{
		"model":    resp.Model,
		"provider": "typesafe",
	})
	setJSONAttr(span, "braintrust.metrics", map[string]any{
		"prompt_tokens":     resp.Usage.InputTokens,
		"completion_tokens": resp.Usage.OutputTokens,
		"tokens":            resp.Usage.InputTokens + resp.Usage.OutputTokens,
	})
}

func setJSONAttr(span oteltrace.Span, key string, value any) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Printf("warning: failed to marshal %s: %v", key, err)
		return
	}
	span.SetAttributes(attribute.String(key, string(jsonBytes)))
}
