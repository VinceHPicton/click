env "local" {
    url = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
    dev = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" // TODO: This will become different in future
}