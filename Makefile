lint:
	@echo "Running revive..."
	@revive ./...

webserver:
	@go run cmd/webserver/main.go
