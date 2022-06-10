up:
	@docker-compose down && docker-compose up --build -d
godoc:
	@cd app && swag init -d cmd/main, internal/user -g app.go --output ../../docs