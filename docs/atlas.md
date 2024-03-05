# Atlas

After testing out several other options such as goose, sql-migrate, and golang-migrate, the database migration tool chosen for this project is Atlas, with versioned migrations rather than declarative.
Atlas integrates with sqlc, is itself easy to use and is open source.

A guide on sqlc integration is even part of the Atlas docs!
https://atlasgo.io/guides/frameworks/sqlc-versioned


### Workflow to change the DB with Atlas + sqlc
1. Change the schema (0_schema.sql) to your newly desired database state (adding/deleting columns, tables etc etc)
1b. Optionally update the queries (query.sql)
2. Run `sqlc generate` to ensure the Go database layer reflects your changes
3. Run the `atlas migrate diff` command with required flags (premade `atlasdiff.sh`) generate the required new migration file which will bring about the desired DB state
2b. Optionally edit this file how you see fit if needed and run `atlas migrate hash` to regenerate the hash file (atlas.sum)
4. Run `atlas schema apply` with required flags (premade in `migrate.sh`) to execute the new migration file against the database
