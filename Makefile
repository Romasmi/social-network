.PHONY: run docker-up

api:
	go run cmd/api/main.go

worker:
	go run cmd/worker/main.go

docker-up:
	docker compose up -d

run: docker-up