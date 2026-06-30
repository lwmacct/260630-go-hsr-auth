# go-hsr-auth

Reusable HSR auth feature module.

## Layout

- `pkg/auth`: public library API for other projects.
- `internal/repository`, `internal/service`, `internal/handler`: private auth implementation.
- Runtime concerns such as database opening, request context middleware, HTTP server wiring, and schema helpers live in `github.com/lwmacct/260630-go-hsr-shared`.

## Checks

```bash
go test -count=1 ./...
go test -count=1 ./internal/testutil/tddcheck
```
