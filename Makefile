MIGRATIONS_DIR = ./server/migrations/goose
DB_FILE = ./server/storage/local.db

# docker
up: build
	docker-compose up --build

down:
	docker-compose down


# bench
apache-bench:
	ab -c 100 -n 10000 http://127.0.0.1:8080/api/v1/proxy


.PHONY:up down apache-bench
