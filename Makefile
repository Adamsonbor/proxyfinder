MIGRATIONS_DIR = ./server/migrations/goose
DB_FILE = ./server/storage/local.db

SERVER_DIR = ./server

# docker
up: server-build frontend-install
	docker-compose up --build

down:
	docker-compose down

build: server-build frontend-build

server-build:
	make -C $(SERVER_DIR) build

frontend-build: frontend-install
	@cd frontend && yarn build

frontend-install:
	@cd frontend && yarn


# bench
apache-bench:
	ab -c 100 -n 10000 http://127.0.0.1:8080/api/v1/proxy


.PHONY:up down apache-bench
