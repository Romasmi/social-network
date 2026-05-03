# Social network mockup

FYI: check researches in docs folder, it's interesting!

## Goal
Create social network which allows to create and view user cards

## How to start
It provides a fully prepared docker env including api for demo and workers.

Run 
```shell
make up
```

Generate user and data
This command registers a user, add 10 friends and creates post for each friend.
```shell
make dummy-data
```
it returns registered user data like:
```json
{
  "id": "019b7b4e-ef4e-7a12-8d6d-da73e00c701d",
  "password": "password"
}
```

Service URLs

| service    | url                             |                          
|------------|---------------------------------|
| API        | http://api.localhost            |
| Redis UI   | http://redis-insight.localhost  |
| Kafka UI   | http://kafka-ui.localhost       |
| Traefik    | http://localhost:8080           |
| Prometheus | http://prometheus.localhost     |
| Grafana    | http://grafana.localhost        |


## Initial requirements
### Functional
- basic/bearer authorization
- create user data: name, surname, birthdate, gender, interests, city.

### Non-functional
- monolithic architecture 

### Endpoints
Specification : https://github.com/OtusTeam/highload/blob/master/homework/openapi.json
- /login  (Returns bearer token)
- /user/register
- /user/get/{id}   (Protected: use bearer token)


## Development
### Run API
```shell
make api
```

### Migrations
#### Create migration
``` bash
migrate create -ext sql -dir migrations -seq create_roles_table
```

### Import users
Go to root folder and run:
```shell
go run cmd/cli/main.go import <link to file> <max count of users to import>
```
```shell
go run cmd/cli/main.go import https://raw.githubusercontent.com/OtusTeam/highload/master/homework/people.v2.csv 1000
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
#### Search user
```shell
 k6 run --vus 1000 --duration 30s -e TOKEN=<token> load_testing/user_search.js
```
#### Feed cache
```shell
k6 run --vus 1000 --duration 30s -e TOKEN="019b6a77-8c9f-7cf8-9a48-f8bd6a714be8" load_testing/get_feed.js
```

## Caching
### User feed
`/feed?offset=0&limit=20` returns latest posts of friends.
It caches latest 1000 posts.

#### On what event to invalidate cache
| Event       | Action description                                  | Condition                         
|-------------|-----------------------------------------------------|-----------------------------------|
| postCreated | push post id to fiends feeds                        |                                   |
| postUpdated | update post data in cache posts:{postId} {postData} |                                   |
| postDeleted | delete post from user friends feeds                 |
| friendAdded | push users posts to feeds                           | user has posts,  friend has posts |
| friendDeleted | remove posts from feed of user and friend           | user has posts,  friend has posts                  |

#### Caching strategy
I use `fan-out on write` approach.
For the sake of simplicity, some steps are omitted below.
##### Get feed
````mermaid
sequenceDiagram
    Client->>API: GET /feed?offset={N}&limit={M}
    API->>Redis: ZRANGE feed:{userId} N N+M-1 REV
    Redis-->>API: PostsIds
    API->>Redis: Get posts by ids
    Redis-->>API: Posts
    API->>DB: Get missed posts by ids
    API->>Redis: Save missed posts to cache
    API-->>Client: Posts
    
````

##### Create feed 
````mermaid
sequenceDiagram
    API->>DB: Get user friends latest posts 
    DB->>API: User friends posts ids LIMIT 1000
    API->>Redis: ZADD feed:{userId} {postId1} {postId2} ... {postIdN} 
    API->>Redis: ZREMRANGEBYRANK feed:{userId} 0 -1001 
````

##### Create/push_to feed onPostCreated
I'll use Redis pipeline for batch operation.
````mermaid
sequenceDiagram
    Client->>API: POST /post/create
    API->>DB: Create post
    API->>DB: Get user friends
    DB->>API: User friends ids
    API->>Redis: ZADD feed:{friendId1} {postId}
    API->>Redis: ZREMRANGEBYRANK feed:{friendId1} 0 -1001
    API->>Redis: ZADD feed:{friendId2} {postId}
    API->>Redis: ZREMRANGEBYRANK feed:{friendId2} 0 -1001
    API->>Redis: ZADD feed:{friendIdN} {postId}
    API->>Redis: ZREMRANGEBYRANK feed:{friendIdN} 0 -1001
    API->>Redis: SET posts:{postId} {postContent}
    API-->>Client: response
````

##### Invalidate cache on friendDeleted

````mermaid
sequenceDiagram
    Client->>API: DELETE /friend/add/{friendId}
    API->>DB: Delete friend from user friends
    API->>DB: Get user's posts
    API->>Redis: Delete user posts from ex-friend feed 
    API->>DB: Get ex-friend's posts
    API->>Redis: Delete ex-friend posts from ex-friend feed
````