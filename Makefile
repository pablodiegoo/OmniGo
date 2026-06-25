.PHONY: run test test-race lint migrate docker-up docker-down

run:
	go run ./cmd/omnigo

test:
	go test ./... -short

test-race:
	go test ./... -race -count=1

lint:
	golangci-lint run

migrate:
	go run ./cmd/omnigo

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
