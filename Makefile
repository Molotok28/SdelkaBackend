include .env
export

export PROJECT_ROOT=$(shell pwd)
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8

env-up:
	@docker compose up -d sdelka-postgres

env-down:
	@docker compose down sdelka-postgres

env-cleanup:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down sdelka-postgres && \
		docker compose down sdelka-port-forwarder && \
		rm -rf out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Операция отменена."; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Ошибка: Не указано имя миграции. Используйте 'make migrate-create seq=имя_миграции'"; \
		exit 1; \
	fi
	MSYS_NO_PATHCONV=1 docker compose run --rm sdelka-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Ошибка, не указан action"; \
		exit 1; \
	fi
	MSYS_NO_PATHCONV=1 docker compose run --rm sdelka-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@sdelka-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"${action}"

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Ошибка: укажи version. Пример: make migrate-force version=1"; \
		exit 1; \
	fi
	MSYS_NO_PATHCONV=1 docker compose run --rm sdelka-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@sdelka-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		force $(version)

env-port-forward:
	@docker compose up -d sdelka-port-forwarder

env-port-close:
	@docker compose down sdelka-port-forwarder

sdelka-run:
	@go run cmd/sdelka/main.go