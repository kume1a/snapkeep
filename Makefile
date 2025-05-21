APP_NAME=snapkeep
MAIN=./cmd/$(APP_NAME)/main.go
BUILD_DIR=bin
TMP_DIR=tmp

.PHONY: run build test clean zip dump backup

run:
	go run $(MAIN) \
		-db "postgres://user:pass@localhost/db" \
		-public ./public \
		-backend ../backend \
		-upload s3

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

clean:
	rm -rf $(TMP_DIR) $(BUILD_DIR)

dump:
	./scripts/dump.sh "postgres://user:pass@localhost/db" "$(TMP_DIR)/manual_dump.sql"

backup:
	go run $(MAIN)

help:
	@echo "Usage:"
	@echo "  make run        Run the backup service"
	@echo "  make build      Build the binary"
	@echo "  make test       Run all Go tests"
	@echo "  make clean      Clean up tmp files and binaries"
	@echo "  make dump       Dump database using script"
	@echo "  make backup     Run backup logic directly"
