# Tech Stack

## Language & runtime
- Go (module `vincehpicton/click`, Go 1.26), located in `goapp/`.

## Key libraries
- `github.com/gorilla/mux` - HTTP router.
- `github.com/jackc/pgx/v5` - Postgres driver.
- `github.com/google/uuid` - UUID types (primary keys are UUIDs).
- `github.com/golang-jwt/jwt/v5` - JWT access tokens.
- `github.com/stretchr/testify` - assertions and suites.
- `github.com/testcontainers/testcontainers-go` (+ postgres module) - ephemeral DB for tests.

## Database tooling
- Postgres is the datastore.
- sqlc generates the Go data-access layer from `db/queries/*.sql` and `db/schema`.
  Never hand-edit `goapp/internal/db/sqlc`; regenerate instead.
- Atlas manages schema migrations (`db/migrations`, `atlas.hcl`).
- pgadmin is available for inspection (`pgadmin_servers.json`).

## Local development
- Copy `.env.example` to `.env` and fill values.
- `sqlc generate` - regenerate DB code (also `make sqlc`, which wipes and regenerates).
- `docker compose build` then `docker compose up` - run the stack locally.
- Scripts in `scripts/` wrap common tasks: `build.sh`, `migrate.sh`, `localrun.sh`,
  `psqlcli.sh`, `atlasdiff.sh`, `rebuild.sh`, `clean.sh`.

## Testing
- Use testify (assert/require and suites) for Go tests.
- Do not write Go unit tests for generated sqlc code; it is regenerated frequently.
  Instead, test the SQL queries themselves under `db/sqlc_test/` against a real
  Postgres via testcontainers.
- Run tests with `go test ./...` from `goapp/`. Prefer a single run (no watch mode).

## Command notes
- The environment shell is bash on Windows. Do not run long-lived processes
  (dev servers, watchers) as blocking commands; hand those to the user to run.
