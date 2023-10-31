# Docker

Locally this app runs with 3 different containers, one for the Go program, one for the postgres database server, and one for pgadmin. A 4th container for psql 
commands can optionally be spun up.

When deployed, likely the db server will be run on a service like RDS, not within a container, but decisions are still to be made around deployment.