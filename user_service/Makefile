deploy:
	@cd app && go mod tidy && go mod vendor
	@docker-compose stop && docker-compose up --build -d
godoc:
	@cd app && swag init -d cmd/main, internal/user -g app.go --output ../../docs