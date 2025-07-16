lint:
	@echo "Running revive..."
	@revive ./...

webserver:
	@go run cmd/webserver/main.go

mockgen:
	@echo "Creating mocks..."
	@mockgen \
    -source=internal/domain/quotes/frete_rapido_repository.go \
    -destination=mocks/quotesmocks/frete_rapido_repository_mock.go \
    -package=quotesmocks \
    -typed
