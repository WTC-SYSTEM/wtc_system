deploy:
	@cd app && go mod tidy && go mod vendor
	@docker-compose stop && docker-compose up --build -d

migrateup:
	@migrate -path app/db/migration -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up

migratedown:
	@migrate -path app/db/migration -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose down
