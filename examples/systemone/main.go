// Command systemone asks typed questions about a support message.
package main

import (
	"context"
	"fmt"
	"log"

	typesafe "github.com/atharvamhaske/typesafe-sdk-go"
)

func main() {
	client, err := typesafe.NewClient() // reads TYPESAFE_API_KEY
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
			"is_dupe": typesafe.Noul{
				Instructions: "Is this a duplicate charge?",
				Criteria:     map[string]string{"true": "Duplicate", "false": "Not a duplicate"},
			},
		})
	if err != nil {
		log.Fatal(err)
	}

	cat := resp.Choices()["category"]
	fmt.Printf("category: %s (confidence %.2f)\n", cat.Choice, cat.Confidence)
	fmt.Printf("urgency:  %.2f\n", resp.Scores()["urgency"].Score)
	fmt.Printf("is_dupe:  %.2f\n", resp.Nouls()["is_dupe"].Noul)
	fmt.Printf("usage:    %d in / %d out\n", resp.Usage.InputTokens, resp.Usage.OutputTokens)
}
