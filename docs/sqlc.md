# sqlc

sqlc is a go package which generates DB access Go code directly from raw SQL statements.
This not only removes the need to write this Go code, but also the tests for it, while also giving granular ability to optimize the SQL queries.
This will save a huge amount of dev time and improve the app's speed.

An ORM like GORM could be used, but we would lose that granular control of the SQL queries, resulting eventually in a very slow app due to inefficient generated queries.

### Key files for sqlc
sqlc reads the db/schema/schema.sql file and the files in db/queries (locations defined in sqlc.yaml config file) to generate Go code.

A quick guide to how sqlc works is below:
https://www.youtube.com/watch?v=uBPXNREhZZw