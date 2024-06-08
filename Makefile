API_SERVER_BINARY=player_server

.PHONY : tidy
.PHONY : build
.PHONY : clean_image


## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

build: build_api_service

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

build_api_service:
	@echo "Building api server binary..."
	cd api-service && GOOS=linux CGO_ENABLED=0 go build -o ${API_SERVER_BINARY}
	@echo "Done!"


## stop: stop the server
stop:
	@echo "Stopping api service..."
	@-pkill -SIGTERM -f "./${API_SERVER_BINARY}"
	@echo "Stopped api service!"


clean: down
	@echo "deleting linux binaries..."
	cd api-service && rm -f ${API_SERVER_BINARY}

clean_image:
	@echo "deleting docker image..."
	# @-docker image rm `docker image ls | grep project-broker-service | awk '{print $3}'`
	@-bin/delete_images.sh

tidy:
	@echo running tidy in authentication-service
	go mod tidy
