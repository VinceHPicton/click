sqlc:
	rm -rf goapp/internal/db/sqlc
	sqlc generate

db:
	rm -rf goapp/internal/db/sqlc
	sqlc generate
	docker compose build db