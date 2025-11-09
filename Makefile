.PHONY: run docker-up

run:
	cd cmd/api && go run main.go

docker-up:
	cd deployment/local && docker compose up -d

local: docker-up run