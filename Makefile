# variables
APP_NAME := mywebapp
SRC_DIR := ./cmd
CONFIG_FILE := configs/api.toml
FLAG := -conf

.PHONY: build and run

build:
	@echo "==> Building ${APP_NAME}.."
	go build ${SRC_DIR}/${APP_NAME}

run:
	@echo "==> Running ${APP_NAME} with config ${CONFIG_FILE}..."
	./${APP_NAME} {FLAG} ${CONFIG_FILE}

build-and-run: build run
	@echo "==> Build completed. Running ${APP_NAME}..."

.PHONY: migrate
migrate:
	migrate -path migrations -database "postgres://localhost:5432/mywebapp?sslmode=disable&user=postgres&password=postgres" up