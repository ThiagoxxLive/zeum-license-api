.PHONY: run start stop restart logs build vet tests tidy

run:
	docker compose up -d
	@echo "✅ zeum-license-api no ar em http://localhost:8090"

start:
	docker compose up -d app

stop:
	docker compose down

restart:
	docker compose restart app

logs:
	docker compose logs -f app

build:
	docker compose exec app go build ./...

vet:
	docker compose exec app go vet ./...

tests:
	docker compose exec app go test ./... -v

tidy:
	docker compose exec app go mod tidy
