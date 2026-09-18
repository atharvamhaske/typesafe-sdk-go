# API

Full reference: [pkg.go.dev](https://pkg.go.dev/github.com/atharvamhaske/typesafe-sdk-go)

## Client

- `typesafe.NewClient(opts ...Option) (*Client, error)`

Options:

- `typesafe.WithAPIKey(key string)`
- `typesafe.WithBaseURL(url string)`
- `typesafe.WithModel(model string)`
- `typesafe.WithHTTPClient(hc *http.Client)`

Environment variables: `TYPESAFE_API_KEY`, `TYPESAFE_BASE_URL`, `TYPESAFE_DEFAULT_MODEL`.

## SystemOne

`POST /v1/systemone`

- `client.SystemOne(ctx, state string, questions map[string]Question, opts ...SystemOneOption) (*SystemOneResponse, error)`
- `typesafe.WithRequestModel(model string)` — per-call model override

Question types (all implement `Question`):

- `typesafe.Choice{Instructions string, Criteria map[string]string}`
- `typesafe.Score{Instructions string, Criteria []string}`
- `typesafe.Noul{Instructions string, Criteria map[string]string}`

Response:

- `resp.Model string`
- `resp.Usage Usage` — `InputTokens`, `OutputTokens`
- `resp.Answers map[string]Answer` — union of the answer types below
- `resp.Choices() map[string]ChoiceAnswer` — `Choice`, `Confidence`, `Probabilities`
- `resp.Scores() map[string]ScoreAnswer` — `Score`, `Confidence`, `Legend`, `Probabilities`
- `resp.Nouls() map[string]NoulAnswer` — `Noul`

## Models

`GET /v1/models`

- `client.ListModels(ctx) ([]Model, error)` — `Name`, `Description`, `ReleaseDate`

## Errors

Non-2xx responses return `*typesafe.Error`:

- `err.StatusCode int`
- `err.Detail json.RawMessage` — validation detail, if any
- `err.Body []byte` — raw response body
