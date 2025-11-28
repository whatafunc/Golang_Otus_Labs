BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
COMPOSE_FILE := deployments/docker-compose.yaml

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# build:
# 	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

# run: build
# 	$(BIN) -config ./configs/config.yaml

# build-img:
# 	docker build \
# 		--build-arg=LDFLAGS="$(LDFLAGS)" \
# 		-t $(DOCKER_IMG) \
# 		-f build/Dockerfile .

# run-img: build-img
# 	docker run -p 8888:8081 -p 50051:50051 $(DOCKER_IMG)

# version: build
# 	$(BIN) version

# test:
# #	go test -race ./internal/... ./pkg/...
# 	go test -race ./internal/...

# install-lint-deps:
# 	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.8

# lint: install-lint-deps
# 	golangci-lint run ./...

# Start all services in detached mode
up:
	docker-compose -f $(COMPOSE_FILE) up -d

# Stop and remove all containers, networks
down:
	docker-compose -f $(COMPOSE_FILE) down

# Build all services without starting them
build:
	docker-compose -f $(COMPOSE_FILE) build

# Follow logs from all services
logs:
	docker-compose -f $(COMPOSE_FILE) logs -f

# Stop services without removing containers
stop:
	docker-compose -f $(COMPOSE_FILE) stop

# Start stopped services
start:
	docker-compose -f $(COMPOSE_FILE) start

# Restart all services
restart:
	docker-compose -f $(COMPOSE_FILE) restart

# View running containers status
ps:
	docker-compose -f $(COMPOSE_FILE) ps

# Clean up all containers, volumes, and networks
clean: down
	docker-compose -f $(COMPOSE_FILE) down -v --remove-orphans
	docker system prune

# Rebuild and restart services
rebuild: down build up

.PHONY: build run build-img run-img version test lint up down build logs clean
