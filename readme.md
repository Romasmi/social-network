# Social network mockup

FYI: check researches in docs folder, it's interesting!

## Goal
Create social network which allows to create and view user cards

## How to start
Ensure port 8000 is available and you have docker installed.
After start Api will be available at http://localhost:8000

### if you have goland 1.25.1^ locally

```shell
make local && go run cmd/api/main.go
```

### If you don't have go 1.25.1^ locally
```shell
make local-fullы
```

## Requirements
### Functional
- basic/bearer authorization
- create user data: name, surname, birthdate, gender, interests, city.

### Non-functional
- monolithic architecture 

## Endpoints
Specification : https://github.com/OtusTeam/highload/blob/master/homework/openapi.json
- /login  (Returns bearer token)
- /user/register
- /user/get/{id}   (Protected: use bearer token)


## Development

### Migrations
#### Create migration
``` bash
migrate create -ext sql -dir migrations -seq create_roles_stable
```

### Import users
Go to root folder and run:
```shell
go run cmd/cli/main.go import -l <link to file>

```

### Import posts
````shell
go run ./cmd/cli importPosts <link> <profile ids to assign each post randomly>
````

Example:
````shell
go run ./cmd/cli importPosts https://raw.githubusercontent.com/OtusTeam/highload/refs/heads/master/homework/posts.txt 019ad582-424c-7316-b98c-eee5bf0f8d08 019ad582-42a7-7d0f-a287-67d3dc0040ae 019ad582-42b3-7340-a5b8-6025ce042c01 019ad582-426a-756e-ad6e-7341623185c4 019ad582-4292-78be-889c-027d95fcc30
````


File structure:
Name | birdate | Profile

### Load testing via K6

```shell
 k6 run --vus 1000 --duration 30s -e TOKEN=<token> load_testing/user_search.js
```
