MIGRATIONS_DIR = ./server/migrations/goose
DB_FILE = ./server/storage/local.db

SERVER_DIR = ./server

# docker
up:
	docker-compose up

down:
	docker-compose down --remove-orphans

up-dev:
	docker-compose -f docker-compose-dev.yaml up

down-dev:
	docker-compose -f docker-compose-dev.yaml down --remove-orphans

up-build: build
	docker-compose up --build

down:
	docker-compose down

build: server-build frontend-build admin-build

server-build:
	make -C $(SERVER_DIR) build

frontend-build: frontend-install
	@cd frontend && sudo yarn build

frontend-install:
	@cd frontend && yarn

admin-build: admin-install
	@cd admin && sudo yarn build

admin-install:
	@cd admin && yarn


# bench
apache-bench:
	ab -c 100 -n 10000 http://127.0.0.1:8080/api/v1/proxy


.PHONY:up down apache-bench build
