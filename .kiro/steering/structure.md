# Project Structure

Repository root is a Go backend. The Go module lives under `goapp/` with module path
`vincehpicton/click`.

## Top-level layout
- `goapp/` - the Go application (module root, all Go code).
- `db/` - all database assets: `migrations/` (Atlas), `queries/` (raw SQL for sqlc),
  `schema/` (declarative schema sqlc reads from).
- `docs/` - rationale docs for tooling choices (atlas, docker, pgadmin, sqlc, testing) and layout.
- `scripts/` - bash scripts supporting build, migrate, local run, and psql access.
- `build/` - image build assets.
- `atlas.hcl`, `docker-compose.yml`, `Makefile` - project-level config and entry points.
- `tmp/pgdata/` - generated local Postgres data files (not committed).

## Go internal packages (`goapp/internal/`)
- `auth/` - authentication business rules (domain layer). See layering rule below.
- `tokens/` - access/refresh token generation and hashing.
- `db/sqlc/` - GENERATED sqlc code. Do not hand-edit; it is regenerated from `db/queries`
  and `db/schema`.
- `db/factory/` - test data factories.
- `db/sqlc_test/` - tests for the SQL queries themselves.
- `httpserver/` - HTTP transport layer: handlers, routes, middleware, response helpers.
- `testdb/` - test database helpers (testcontainers).
- `examplesuite/` - example testify suite for reference.

## Layering rules (important)
- The domain layer (`auth`, `tokens`) must not import `net/http`. It takes plain inputs,
  returns domain types and sentinel errors, and leaves status codes, headers, and JSON
  encoding to the transport layer. `auth/auth.go` documents this explicitly.
- The transport layer (`httpserver`) owns HTTP concerns and calls into the domain layer.
- Generated DB code (`db/sqlc`) is the only data-access path; business logic composes
  sqlc queries, opening transactions when a flow needs multiple atomic writes.

## Naming conventions observed in the codebase
- Handler files: `<feature>_handler.go`; matching tests: `<feature>_test.go`.
- Shared test setup files use a `_SETUP_test.go` suffix.
- Sentinel domain errors live in an `errors.go` within the package that owns them.
