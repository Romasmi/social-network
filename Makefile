.PHONY: run docker-up

run:
	go run cmd/api/main.go

docker-up:
	cd deployment/local && docker compose up -d

local: docker-up
	@sleep 5
	@$(MAKE) run

local-full:
	cd deployment/local_full && docker compose up -d