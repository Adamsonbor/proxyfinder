GOOSE_PATH=build/goose
SOURCE_DIR=cmd
SOURCE=$(shell find cmd -name "*.go")
BUILD_DIR=build
BUILDFILES=$(patsubst $(SOURCE_DIR)/%,$(BUILD_DIR)/%, $(dir $(SOURCE)))


build: $(BUILDFILES)
	rm -r $(BUILD_DIR)/*
	@echo $(BUILDFILES)
	@echo $(SOURCE)

#build/geonode: cmd/collector/geonode/main.go
$(BUILD_DIR)/%: $(SOURCE_DIR)/%/main.go
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -ldflags "-linkmode external -extldflags '-static'" -o $@ $<

geonode: $(BUILD_DIR)/geonode


goose-up:
	$(GOOSE_PATH) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_FILE) up

goose-reset:
	$(GOOSE_PATH) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_FILE) reset
	$(GOOSE_PATH) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_FILE) up

goose-down:
	$(GOOSE_PATH) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_FILE) down


.PHONY: build goose-up goose-reset goose-down
