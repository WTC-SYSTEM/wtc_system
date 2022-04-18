deploy:
	@cd app && go mod tidy && go mod vendor
	@docker-compose stop && docker-compose up --build -d