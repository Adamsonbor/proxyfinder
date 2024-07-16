SRC_DIR=cmd
SRC=$(shell find cmd -name "*.go")
BIN_DIR=build
BIN=$(addprefix build/,$(patsubst cmd/%/main.go,%,$(SRC)))

DB_PATH=./storage/local.db
MIGRATIONS_DIR=./migrations/goose
GOOSE_PATH=build/goose
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

goose: goose-rm $(BIN_DIR)/goose

goose-rm:
	rm $(BIN_DIR)/goose

goose-up: goose
	$(GOOSE_COMMAND) up

goose-reset: goose
	$(GOOSE_COMMAND) reset
	$(GOOSE_COMMAND) up

goose-down: goose
	$(GOOSE_COMMAND) down

goose-status: goose
	$(GOOSE_COMMAND) status



.PHONY: build goose-up goose-reset goose-down goose-status clean re
