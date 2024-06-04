#!make

include .env
export $(shell sed 's/=.*//' .env)

DOCKER_COMPOSE_FILE ?= docker-compose.yml

#========================#
#== DEVELOPMENT ==#
#========================#

up:
	docker compose -f ${DOCKER_COMPOSE_FILE} up -d --remove-orphans

down:
	docker compose -f ${DOCKER_COMPOSE_FILE} down


install:
	go mod download && \
	go mod tidy

update-codegen:
	mkdir -p pkg/generated && \
	protoc --proto_path=./proto:/usr/local/include --go_out=./pkg/generated --go-grpc_out=./pkg/generated $(shell find proto -type f -name '*.proto')


#========================#
#== BUILD & RUN ==#
#========================#
build:
	go build -o bin/server cmd/main.go

run:
	go run cmd/main.go

run-mirrord:
	mirrord exec -f .mirrord/mirrord.json go run cmd/main.go