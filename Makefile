SERVER_DIR = ./server
SERVER_BUILD_DIR = ${SERVER_DIR}/build
SERVER_CONFIG_DIR = ${SERVER_DIR}/config
SERVER_MIGRATIONS_DIR = ${SERVER_DIR}/migrations

FRONTEND_DIR = ./frontend
FRONTEND_BUILD_DIR = ${FRONTEND_DIR}/dist

ADMIN_DIR = ./admin
ADMIN_BUILD_DIR = ${ADMIN_DIR}/dist

BUILD_DIRS = ${SERVER_BUILD_DIR} ${FRONTEND_BUILD_DIR} ${ADMIN_BUILD_DIR}

# docker
up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans

up-dev:
	docker-compose -f docker-compose-dev.yaml up

down-dev:
	docker-compose -f docker-compose-dev.yaml down --remove-orphans

up-build: re
	docker-compose up --build

build: ${BUILD_DIRS}

clean:
	sudo rm -rf ${BUILD_DIRS}

re: clean build

${SERVER_BUILD_DIR}:
	make -C ${SERVER_DIR} build

${FRONTEND_BUILD_DIR}:
	@cd ${FRONTEND_DIR} && yarn && yarn build:prod

${ADMIN_BUILD_DIR}:
	@cd ${ADMIN_DIR} && yarn && yarn build


prod: re
	sudo rm -rf prod
	mkdir -p \
		prod/server\
		prod/frontend\
		prod/admin
	cp -r ${SERVER_DIR} prod
	cp -r ${FRONTEND_BUILD_DIR} prod/frontend
	cp -r ${ADMIN_BUILD_DIR} prod/admin
	cp -r ./nginx.conf prod
	cp -r ./docker-compose.yaml prod
	cp -r ./Makefile prod
	tar -czf prod.tar.gz prod

# bench
apache-bench:
	ab -c 100 -n 10000 http://127.0.0.1:8080/api/v1/proxy


.PHONY:up down apache-bench build prod
