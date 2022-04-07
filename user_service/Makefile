redeploy:
	@docker-compose down && docker-compose build && docker-compose up -d

migrateup:
	@migrate -path app/db/migration -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up

migratedown:
	@migrate -path app/db/migration -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose down
