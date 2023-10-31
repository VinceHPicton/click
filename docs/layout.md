# Project Layout

This project is as follows:

* [./build](../build/) - contains scripts to build various images
* [./db](../db/) - contains all db files (including those for sqlc and Atlas), such as migrations, seed data, raw sql queries, declarative schema
* [./scripts](../scripts/) - contains bash scripts to support project operations such as building code.

The following are generated folders:

* [./tmp/pgdata](../tmp/pgdata) -- contains files to hold local postgres data files 