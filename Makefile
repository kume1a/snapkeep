APP_NAME=snapkeep
MAIN=./main.go
BUILD_DIR=bin
TMP_DIR=tmp

.PHONY: run build test clean zip dump backup

run:
	go run $(MAIN)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

clean:
	rm -rf $(TMP_DIR) $(BUILD_DIR)

qdash:
	asynq dash --uri "localhost:6380" --password ""

docker-dev:
	docker-compose -f docker-compose.development.yaml up

help:
	@echo "Usage:"
	@echo "  make run        Run the backup service"
	@echo "  make build      Build the binary"
	@echo "  make clean      Clean up tmp files and binaries"
	@echo "  make qdash      Run asynq dashboard (requires asynq installed)"
	@echo "  make docker-dev Run the development Docker Compose file"
