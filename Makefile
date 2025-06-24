run:
	@go run cmd/api/main.go
migrate-up:
	@go run cmd/api/main.go -m up
migrate-down:
	@go run cmd/api/main.go -m down
generate-docs:
	@swag init -g cmd/api/main.go