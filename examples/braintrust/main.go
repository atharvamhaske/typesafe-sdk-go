// This example shows automatic tracing of typesafe-sdk-go calls with
// braintrust-sdk-go's trace/contrib/typesafe middleware. The middleware
// wraps the client's *http.Client and traces every POST /v1/systemone
// call, so no manual span code is needed in application code.
package main

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/braintrustdata/braintrust-sdk-go"
	bttypesafe "github.com/braintrustdata/braintrust-sdk-go/trace/contrib/typesafe"

	typesafe "github.com/atharvamhaske/typesafe-sdk-go"
)

func main() {
	tp := trace.NewTracerProvider()
	defer tp.Shutdown(context.Background()) //nolint:errcheck
	otel.SetTracerProvider(tp)

	if _, err := braintrust.New(tp,
		braintrust.WithProject("typesafe-sdk-go-examples"),
		braintrust.WithBlockingLogin(true),
	); err != nil {
		log.Fatal(err)
	}

	client, err := typesafe.NewClient(
		typesafe.WithHTTPClient(bttypesafe.Client()), // reads TYPESAFE_API_KEY
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.SystemOne(context.Background(),
		"I was charged twice. Please refund the duplicate charge today.",
		map[string]typesafe.Question{
			"category": typesafe.Choice{
				Instructions: "Categorize the message",
				Criteria:     map[string]string{"billing": "Billing issue", "technical": "Technical issue"},
			},
			"urgency": typesafe.Score{
				Instructions: "Rate urgency",
				Criteria:     []string{"Can wait", "Needs attention"},
			},
		})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Choices()["category"].Choice)
	fmt.Println(resp.Scores()["urgency"].Score)
}
