APP_NAME := subscriptions-service
MAIN_PKG := ./cmd/app
GOCACHE ?= /private/tmp/gocache
COVERPKG ?= ./...
COVERPROFILE ?= coverage.out

.PHONY: test coverage build run docker-up docker-down docker-reset clean

test:
	@mkdir -p $(GOCACHE)
	GOCACHE=$(GOCACHE) go test ./...

coverage:
	@mkdir -p $(GOCACHE)
	GOCACHE=$(GOCACHE) go test -coverpkg=$(COVERPKG) -coverprofile=$(COVERPROFILE) ./...
	go tool cover -html=$(COVERPROFILE)

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
