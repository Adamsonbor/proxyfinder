build:
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -C server -ldflags "-linkmode external -extldflags '-static'" -o build ./...

up: build
	docker-compose up --build

down:
	docker-compose down

apache-bench:
	ab -c 100 -n 10000 http://127.0.0.1:8080/api/v1/proxy


.PHONY: build up down
