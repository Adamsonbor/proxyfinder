SRC_DIR=cmd
SRC=$(shell find cmd -name "*.go")
BIN_DIR=build
BIN=$(addprefix build/,$(patsubst cmd/%/main.go,%,$(SRC)))

DB_PATH=./storage/local.db
MIGRATIONS_DIR=./migrations/goose
GOOSE_PATH=$(BIN_DIR)/goose
GOOSE_COMMAND=$(GOOSE_PATH) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH)


build: $(BIN_DIR) $(BIN)

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
