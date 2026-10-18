---
inclusion: fileMatch
fileMatchPattern: '**/*.go'
---

# Go Coding Standards

These conventions are derived from the existing codebase (see `internal/auth/auth.go`
as a reference for the expected style). Apply them when writing or editing Go code.

## Comments
- Do not add comments to your code, your naming of variables and functions should be used to allow the code to explain itself
- Comments should be brief and only used for docstrings for packages and methods/functions, docstrings should only be used where they genuinely add value which good naming does not already provide
- Where you add a docstring comment to explain to me (the human engineer) something you've done, start your comment with "-- AI EXPLANATION FOR DEV --"

## Errors
- Wrap errors with context using `fmt.Errorf("<action>: %w", err)` so the chain is
  preserved and greppable, but only do this where it adds value, otherwise just return err.

## Data access
- All DB access goes through generated sqlc `*Queries`. Do not hand-write SQL execution
  in Go and do not edit `internal/db/sqlc`.
- Never write SQL queries as string literals, only edit or add to /db/migrations, you may also add to factories when needed after which you must run `make` in root to regenerate Go code.
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
- Only add tests which genuinely add value
