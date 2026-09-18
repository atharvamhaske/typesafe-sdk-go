# typesafe-sdk-go

[![ci](https://github.com/atharvamhaske/typesafe-sdk-go/actions/workflows/ci.yml/badge.svg)](https://github.com/atharvamhaske/typesafe-sdk-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/atharvamhaske/typesafe-sdk-go.svg)](https://pkg.go.dev/github.com/atharvamhaske/typesafe-sdk-go)

Unofficial Go client for the TypeSafe AI API. See [api.md](api.md) for the method index and [pkg.go.dev](https://pkg.go.dev/github.com/atharvamhaske/typesafe-sdk-go) for full reference docs.

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

## Examples

Runnable examples live in [examples/](examples/):

```
TYPESAFE_API_KEY=... go run ./examples/systemone
TYPESAFE_API_KEY=... go run ./examples/listmodels
```

## Verified against the live API

The raw endpoint via curl:

```console
$ curl -s https://api.typesafe.ai/v1/systemone \
    -H "Authorization: Bearer $TYPESAFE_API_KEY" \
    -H "Content-Type: application/json" \
    -d '{"state":"I was charged twice. Please refund the duplicate charge today.",
         "model":"jev-latest",
         "questions":{"category":{"type":"choice","instructions":"Categorize the message",
           "criteria":{"billing":"Billing issue","technical":"Technical issue"}}}}'
{
  "model": "jev-1.13.0",
  "answers": {
    "category": {
      "type": "choice",
      "choice": "billing",
      "confidence": 1.0,
      "probabilities": {"billing": 1.0, "technical": 0.0}
    }
  },
  "usage": {"input_tokens": 319, "output_tokens": 31}
}
```

The same call through the SDK:

```console
$ TYPESAFE_API_KEY=... go run ./examples/systemone
category: billing (confidence 1.00)
urgency:  1.00
is_dupe:  0.66
usage:    381 in / 64 out
```

`live_test.go` runs this against the real API in CI whenever `TYPESAFE_API_KEY` is set, and skips otherwise.

## Screenshot 


