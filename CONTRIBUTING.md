# Contributing

## Setup

```
git clone https://github.com/atharvamhaske/typesafe-sdk-go
cd typesafe-sdk-go
go test ./...
```

Tests run against a local `httptest.Server` and need no API key.

## Live tests

Set `TYPESAFE_API_KEY` to also run `TestLive` against the real API:

```
TYPESAFE_API_KEY=... go test ./...
```

## Before you open a PR

```
gofmt -w .
go vet ./...
go test ./...
```

Use conventional commit messages (`feat:`, `fix:`, `docs:`, `test:`, `ci:`).
