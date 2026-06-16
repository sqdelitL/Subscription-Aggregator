.PHONY: run run-migrate migrate-up migrate-down migrate-new

run:
	docker compose up --build

run-migrate:
	docker compose run --rm --build --service-ports app ./main-app -migrate

migrate-up:
	docker compose build app
	docker compose run --rm app ./migrate-tool up

migrate-down:
	docker compose build app
	docker compose run --rm app ./migrate-tool down

migrate-new:
