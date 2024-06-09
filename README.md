## Player Server

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

The default `page_size` is 250. The max `page_size` is 1000.

Paging is only used if the `page` query parameter is present and valid. In other words, specifying `page_size` without `page` has no effect.

### Makefile

You can use the makefile to build the Docker container and start it. Here are some of the make commands you can use.

```
make run         # Builds the executable and runs it

make runr        # Builds the executable and runs with gin in release mode

make up_build    # Builds the executables, composes a Docker image, and starts it in Docker

make down        # Stops the Docker images

make up          # Does docker-compose up -d without building executables
```


#### Database

If you run the solution in Docker, you will be using a MariaDB instance in the container. Otherwise, a local MySQL database should be running with a table named `rest_server`. There should be a user named `rest_api_user` with password `rest_api_pw`, and this user should have all privileges on the `rest_server` database.

The database data is stored in `./data.nobackup`. That directory will be created if it doesn't exist when the docker container is started.

### Opportunities for Improvement

- Translation from a DB record to a Player object could be much more elegant. There are probably some packages available to make this very slick and clean, but a more brute-force approach is currently embodied.

- For a production service, some validation of the input csv data should be added.

- There are currently no unit or integration tests.

- I've been unsuccessful in mounting the csv directory into the container as of yet. So, the Player.csv file is just copied into the container. The code checks every minute for a change in the file (and updates the DB if a change is detected), but it will never be triggered because of the static csv file in the container.
