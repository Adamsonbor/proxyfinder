SRC_DIR=cmd
SRC=$(shell find cmd -name "*.go")
BIN_DIR=build
BIN=$(addprefix build/,$(patsubst cmd/%/main.go,%,$(SRC)))

DB_PATH=./storage/local.db
DB_PROD_PATH=./storage/prod.db
MIGRATIONS_DIR=./migrations/goose
GOMIGRATION_PATH=$(BIN_DIR)/initdb
GOMIGRATION_COMMAND=$(GOMIGRATION_PATH) --dir $(MIGRATIONS_DIR)
GOOSE_PATH=$(BIN_DIR)/goose
GOOSE_COMMAND=$(GOOSE_PATH) -dir $(MIGRATIONS_DIR) sqlite3

# PRODUCTION
up-api:
	build/api --config config/prod.yaml

up-scheduler:
	build/scheduler --config config/prod.yaml

build: $(BIN_DIR) $(BIN)

migrate:
	$(GOOSE_COMMAND) $(DB_PROD_PATH) up
	$(GOMIGRATION_COMMAND) --db $(DB_PROD_PATH) up


# DEVELOPMENT
migrate-dev:
	$(GOOSE_COMMAND) $(DB_PATH) up
	$(GOMIGRATION_COMMAND) --db $(DB_PATH) up

up-api-dev:
	air -c .air-api.toml

up-scheduler-dev:
	air -c .air-scheduler.toml

up-mail-dev:
	air -c .air-mail.toml


# HELPERS
$(BIN_DIR)/%: $(SRC_DIR)/%/main.go
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -ldflags "-linkmode external -extldflags '-static'" -o $@ $<

$(BIN_DIR):
	mkdir -p $@

clean:
	rm -rf $(BIN_DIR)

re: clean build

goose-up:
	$(GOOSE_COMMAND) up

goose-reset:
	$(GOOSE_COMMAND) reset
	$(GOOSE_COMMAND) up

goose-down:
	$(GOOSE_COMMAND) down

goose-status:
	$(GOOSE_COMMAND) status

goose-create:
	$(GOOSE_COMMAND) create $(ARG)

.PHONY: build goose-up goose-reset goose-down goose-status clean re
