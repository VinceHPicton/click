---
inclusion: fileMatch
fileMatchPattern: '**/*.go'
---

# Go Coding Standards

These conventions are derived from the existing codebase (see `internal/auth/auth.go`
as a reference for the expected style). Apply them when writing or editing Go code.

## Package & documentation
- Every package has a package-level doc comment on the `package` clause explaining its
  responsibility and any hard constraints (e.g. "must not import net/http").
- Exported types and functions have doc comments. Where a decision is non-obvious,
  explain the "why", not just the "what" — the codebase favors comments that capture
  design rationale and known trade-offs.

## Errors
- Wrap errors with context using `fmt.Errorf("<action>: %w", err)` so the chain is
  preserved and greppable.
- Define sentinel errors (e.g. `ErrInvalidCode`, `ErrUserNotFound`) in an `errors.go`
  within the owning package. Compare with `errors.Is`.
- Domain functions return sentinel/domain errors; the transport layer maps them to
  HTTP status codes. Never leak `net/http` concerns into domain packages.

## Data access
- All DB access goes through generated sqlc `*Queries`. Do not hand-write SQL execution
  in Go and do not edit `internal/db/sqlc`.
- For multi-write flows that must be atomic, open a transaction and pass a
  transaction-scoped `*sqlc.Queries` into helpers, matching the pattern in `auth`.
- When ordering matters for correctness or security (e.g. consume-then-check for
  one-time codes), document the ordering and why it must not change.

## Types
- Use `uuid.UUID` for identifiers.
- Keep secrets out of stored state: store only hashes of refresh tokens; return
  plaintext to the caller once and never persist it.

## Testing
- Use testify. Test SQL queries under `db/sqlc_test/` by following the pattern already used in the package.
- Do not add unit tests for generated sqlc Go code.
- Name tests `<feature>_test.go` and shared setup `*_SETUP_test.go`.

## Formatting
- Run `gofmt`/`go vet` clean. Match existing import grouping (std, module-internal,
  third-party).
