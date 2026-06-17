.PHONY: build run run-migrate migrate-up migrate-down migrate-new

build:
	docker compose build

run:
	docker compose up -d db app

run-migrate:
	docker compose up -d db
	docker compose run --rm migrate
	docker compose up -d app

migrate-up:
	docker compose run --rm migrate /app/migrate-tool up

migrate-down:
	docker compose run --rm migrate /app/migrate-tool down

migrate-new:
ifndef name
	$(error Usage: make migrate-new name=migration_name)
endif
	docker compose run --rm migrate /app/migrate-tool new $(name)
