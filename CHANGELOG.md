# Changelog

All notable changes to this project are documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [v0.2.0] - 2026-09-21

### Added
- Retries with backoff for network errors, 429s, and 5xxs, honoring `Retry-After` (`WithMaxRetries`, default 3)
- Opt-in in-memory cache for `SystemOne` responses, keyed on the request body (`WithCache`)

### Docs
- Hosted docs site restructured to match the official Python SDK and other Go SDKs' conventions
- References section linking the official TS/Python SDKs and the OpenAPI spec
- Braintrust tracing guide, using the real `trace/contrib/typesafe` middleware
- Community-pattern examples (spam filter, PR labeler, game move picker), each marked as unrelated to this SDK where the linked project isn't Go
- Live demo embeds and a personal-favorite Jev use cases section, fully credited

## [v0.1.0] - 2026-09-18

Initial public release.

### Added
- `Client` with `TYPESAFE_API_KEY`, `TYPESAFE_BASE_URL`, `TYPESAFE_DEFAULT_MODEL` env support and functional options (`WithAPIKey`, `WithBaseURL`, `WithModel`, `WithHTTPClient`)
- `SystemOne` — POST `/v1/systemone` with typed `Choice`, `Score`, and `Noul` questions
- `ListModels` — GET `/v1/models`
- Typed answer union (`ChoiceAnswer`, `ScoreAnswer`, `NoulAnswer`) via `resp.Choices()`, `resp.Scores()`, `resp.Nouls()`
- Typed `*Error` for non-2xx responses, including validation detail
- Runnable examples in `examples/`
- Upstream OpenAPI spec committed at `openapi.json`

[v0.2.0]: https://github.com/atharvamhaske/typesafe-sdk-go/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/atharvamhaske/typesafe-sdk-go/releases/tag/v0.1.0
