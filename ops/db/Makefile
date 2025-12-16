SHELL := /bin/bash

.PHONY: help migrate-create schema-export

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Goose commands
MIGRATION_DIR = ./migration_postgres

# Schema export commands (single schema: mahking_local)
INITDB_DIR ?= ../../go/postgres-initdb.d
POSTGRES_CONTAINER ?= go-postgres-db-1
POSTGRES_DB ?= postgres
SCHEMA_NAME ?= mahking_local
SCHEMA_FILE ?= 01_schema.sql

schema-export: ## 現在のDBから mahking_local スキーマを抽出して go/postgres-initdb.d へ保存
	@mkdir -p $(INITDB_DIR)
	@echo "Exporting schema '$(SCHEMA_NAME)' from $(POSTGRES_CONTAINER)";
	docker exec $(POSTGRES_CONTAINER) pg_dump -U postgres -d $(POSTGRES_DB) -n $(SCHEMA_NAME) -s > $(INITDB_DIR)/$(SCHEMA_FILE)
	sed -i '' '/^SECURITY LABEL FOR anon /d' $(INITDB_DIR)/$(SCHEMA_FILE)
	sed -i '' "s/$(SCHEMA_NAME)\.//g" $(INITDB_DIR)/$(SCHEMA_FILE)
	awk '/CREATE SCHEMA $(SCHEMA_NAME);/ {print; print "SET search_path TO $(SCHEMA_NAME);"; next} 1' $(INITDB_DIR)/$(SCHEMA_FILE) > $(INITDB_DIR)/.tmp_schema && mv $(INITDB_DIR)/.tmp_schema $(INITDB_DIR)/$(SCHEMA_FILE)
	sed -i '' 's/CREATE SCHEMA $(SCHEMA_NAME);/CREATE SCHEMA IF NOT EXISTS $(SCHEMA_NAME);/g' $(INITDB_DIR)/$(SCHEMA_FILE)
