## Player DB

### Overview
This repo contains one possible implementation of a rest server that populates a DB from a csv file and serves data from the database on two endpoints.

### Problem Statement
Your assignment is to create a microservice which serves the contents of Player.csv through a REST API.

The service should expose two REST endpoints:

```
GET /api/players - returns the list of all players.
GET /api/players/{playerID} - returns a single player by ID.
```

Please create unit tests that cover the core logic.

With time permitting, package the application for distribution. Some examples of this:

- Docker image (preferred)
- Tomcat WAR
- Static binary

### Current Status

The project builds a docker container (`make up_build`) and serves the API on port 8800. Thus, you can test it with curl, using something like this:

```
# Gets all players
curl -s http://localhost:8800/api/players     

# Gets one page (default is 250) players starting after the 5th page     
curl -s http://localhost:8800/api/players?page=5   

# Gets one 1000 players starting after the 3rd page
curl -s 'http://localhost:8800/api/players?page=3&page_size=1000'

```
Obviously, you can use postman or whatever tool you prefer to query endpoints.

The default page size is 250.

Paging is only used if the `page` query parameter is present and valid.


#### Database

If you run the solution in Docker, you will be using a MariaDB instance in the container.

If you run the code outside of Docker, you will be using a sqlite3 database (by default).

You can choose to use a local instance of MySQL or MariaDB by running the code like this:

```
USE_MYSQL=1 go run ./...
```

or

```
go build
USE_MYSQL=1 ./player_rest_server
```

### Opportunities for Improvement

The DB population code, and the translation from a DB record to a Player object could be much more elegant. There are probably some packages available to make this very slick and clean, but a more brute-force approach is currently embodied.