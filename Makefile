server:
	@go run ./cmd/http/ .

exp:
	@go run ./cmd/exp/ .


migrate-up:
	@migrate -database $(DB_CONN_URL) -path db/migrations up
migrate-down:
	@migrate -database $(DB_CONN_URL) -path db/migrations down

.PHONY: server exp
