.PHONY: run docker-up

api:
	go run cmd/api/main.go

worker:
	go run cmd/worker/main.go

docker-up:
	docker compose up -d

run: docker-up

import-users:
	go run cmd/cli/main.go import https://raw.githubusercontent.com/OtusTeam/highload/master/homework/people.v2.csv 10

dummy-data: import-users
	@echo "Creating dummy user..."
	@ID=$$(curl -s -X POST 'http://api.localhost/user/register' \
		--header 'Content-Type: application/json' \
		--data '{"first_name":"Имя","second_name":"Фамилия","birthdate":"2017-02-01","biography":"Хобби, интересы и т.п.","city":"Москва","password":"password"}' \
		| grep -o '"id":"[^"]*"' | cut -d'"' -f4); \
	echo ""; \
	go run ./cmd/cli addRandomFriendsWithPosts $$ID 10; \
	echo "{"; \
	echo "  \"id\": \"$$ID\","; \
	echo "  \"password\": \"password\""; \
	echo "}"; \
	echo "";

k6-search:
	k6 run --vus 1000 --duration 5m -e TOKEN=$(TOKEN) load_testing/user_search.js

k6-get-user:
	k6 run --vus 1000 --duration 5m -e TOKEN=$(TOKEN) load_testing/user_get.js

k6-feed:
	k6 run --vus 1000 --duration 5m -e TOKEN=$(TOKEN) load_testing/get_feed.js