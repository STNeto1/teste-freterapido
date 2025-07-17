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
	@mockgen \
    -source=internal/domain/quotes/clickhouse_repository.go \
    -destination=mocks/quotesmocks/clickhouse_repository_mock.go \
    -package=quotesmocks \
    -typed
	@mockgen \
    -source=internal/domain/analytics/clickhouse_repository.go \
    -destination=mocks/analyticsmocks/clickhouse_repository_mock.go \
    -package=analyticsmocks \
    -typed
