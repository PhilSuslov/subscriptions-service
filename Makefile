APP_NAME := subscriptions-service
MAIN_PKG := ./cmd/app
GOCACHE ?= /private/tmp/gocache

.PHONY: test build run docker-up docker-down clean

test:
	@mkdir -p $(GOCACHE)
	GOCACHE=$(GOCACHE) go test ./...

build:
	go build -o bin/$(APP_NAME) $(MAIN_PKG)

run:
	GOCACHE=$(GOCACHE) go run $(MAIN_PKG)

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

clean:
	rm -rf bin
