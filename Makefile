APP_NAME := subscriptions-service
MAIN_PKG := ./cmd/app
COVER_PROFILE := coverage.out

.PHONY: test coverage build run docker-up docker-down docker-reset clean

test:
	@go test ./...

coverage:
	@go test -coverpkg=./... -coverprofile=$(COVER_PROFILE) ./...
	@go tool cover -html=$(COVER_PROFILE)

build:
	@go build -o bin/$(APP_NAME) $(MAIN_PKG)

run:
	@go run $(MAIN_PKG)

docker-up:
	@docker compose up --build -d

docker-down:
	@docker compose down

clean:
	@rm -rf bin
