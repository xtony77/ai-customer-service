-include .env

MIGRATE ?= migrate
MIGRATIONS_DIR ?= migrations
DB_HOST ?= localhost
DB_PORT ?= $(if $(POSTGRESQL_PORT),$(POSTGRESQL_PORT),5432)
DB_USER ?= $(if $(POSTGRESQL_USER),$(POSTGRESQL_USER),postgres)
DB_PASSWORD ?= $(if $(POSTGRESQL_PASSWORD),$(POSTGRESQL_PASSWORD),postgres)
DB_NAME ?= $(if $(POSTGRESQL_DATABASE),$(POSTGRESQL_DATABASE),postgres)
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: migrate-up wire

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

wire:
	wire ./internal/wire/wire.go
