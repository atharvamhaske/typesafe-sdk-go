# typesafe-sdk-go

Unofficial Go client for the TypeSafe AI API.

```
go get github.com/atharvamhaske/typesafe-sdk-go
```

## Usage

```go
client, err := typesafe.NewClient() // reads TYPESAFE_API_KEY, TYPESAFE_BASE_URL, TYPESAFE_DEFAULT_MODEL
if err != nil {
    log.Fatal(err)
}

resp, err := client.SystemOne(ctx, "I was charged twice. Please refund the duplicate charge today.",
    map[string]typesafe.Question{
        "category": typesafe.Choice{Instructions: "Categorize the message", Criteria: map[string]string{
            "billing": "Billing issue", "technical": "Technical issue",
        }},
        "urgency": typesafe.Score{Instructions: "Rate urgency", Criteria: []string{"Can wait", "Needs attention"}},
        "is_dupe": typesafe.Noul{Instructions: "Is this a duplicate charge?", Criteria: map[string]string{
            "true": "Duplicate", "false": "Not a duplicate",
        }},
    })
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.Choices()["category"].Choice) // "billing"
fmt.Println(resp.Scores()["urgency"].Score)    // 1.7
fmt.Println(resp.Nouls()["is_dupe"].Noul)      // 0.98

models, err := client.ListModels(ctx)
```

Options: `WithAPIKey`, `WithBaseURL`, `WithModel`, `WithHTTPClient` on `NewClient`; `WithRequestModel` per call. Non-2xx responses return `*typesafe.Error` with `StatusCode` and validation `Detail`.
