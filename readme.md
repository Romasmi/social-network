# Social network mockup

## Goal
Create social network which allows to create and view user cards

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