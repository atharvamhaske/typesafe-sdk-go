package typesafe_test

import (
	"context"
	"fmt"
	"log"

	typesafe "github.com/atharvamhaske/typesafe-sdk-go"
)

func ExampleClient_SystemOne() {
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

	fmt.Println(resp.Choices()["category"].Choice)
	fmt.Println(resp.Scores()["urgency"].Score)
	fmt.Println(resp.Nouls()["is_dupe"].Noul)
}

func ExampleClient_ListModels() {
	client, err := typesafe.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	models, err := client.ListModels(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range models {
		fmt.Println(m.Name, m.Description)
	}
}
