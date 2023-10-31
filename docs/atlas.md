# Atlas

After testing out several other options such as goose, sql-migrate, and golang-migrate, the database migration tool chosen for this project is Atlas, with versioned migrations rather than declarative.
Atlas integrates with sqlc, is itself easy to use (and actually functions correctly unlike some of the others), and is open source.

A guide on sqlc integration is even part of the Atlas docs!
https://atlasgo.io/guides/frameworks/sqlc-versioned


### Workflow to change the DB with Atlas+sqlc
1. Update the schema (schema.sql)
1b. Optionally update the queries (query.sql)
2. Run sqlc to generate the updated code
3. Run atlas migrate diff to compute and store the required changes
4. Run Atlas during the development and migration process, referencing the migration directory.