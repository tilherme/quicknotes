server:
	go run ./cmd/http/ .

exp:
	go run ./cmd/exp/ .

.PHONY: server exp