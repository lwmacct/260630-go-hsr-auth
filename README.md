# go-hsr-auth

Reusable auth module plus an example HTTP server.

## Layout

- `pkg/auth`: public library API for other projects.
- `internal/repository`, `internal/service`, `internal/handler`: private auth implementation.
- `internal/appcmd/server`: example server wiring for config, database, TLS, and HTTP serving.

## Checks

```bash
go test -count=1 ./...
go test -count=1 ./internal/testutil/tddcheck
```
