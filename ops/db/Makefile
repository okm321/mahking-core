SHELL := /bin/bash

.PHONY: help schema-export proxy db-password db-connect migrate-status migrate-apply migrate-diff

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Cloud SQL 接続設定
ENV             ?= dev
GCP_PROJECT     := mahking-$(ENV)
INSTANCE        := mahking-$(ENV)-db
REGION          := asia-northeast1
PROXY_PORT      := 15432
CONNECTION_NAME := $(GCP_PROJECT):$(REGION):$(INSTANCE)

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

proxy:          ## Cloud SQL Auth Proxy を起動（別ターミナルで実行）
	cloud-sql-proxy $(CONNECTION_NAME) \
		--address 127.0.0.1 \
		--port $(PROXY_PORT)

db-password:    ## Secret Manager から DB パスワードを取得
	@gcloud secrets versions access latest \
		--secret=mahking-$(ENV)-db-password \
		--project=$(GCP_PROJECT)

db-connect:     ## pgcli で Cloud SQL に接続
	@pgcli "postgres://app_user:$$($(MAKE) -s db-password)@localhost:$(PROXY_PORT)/mahking_$(ENV)?sslmode=disable"

migrate-status: ## マイグレーション状態を確認
	@atlas migrate status --env $(ENV) \
		--var "db_password=$$($(MAKE) -s db-password)"

migrate-status-history: ## マイグレーション状態を確認
	@atlas-migrate-status --url "postgres://app_user:$$($(MAKE) -s db-password)@localhost:15432/mahking_dev?search_path=public&sslmode=disable"

migrate-apply:  ## マイグレーションを適用
	@atlas migrate apply --env $(ENV) \
		--var "db_password=$$($(MAKE) -s db-password)"

migrate-diff:   ## スキーマ差分からマイグレーションファイルを生成
	@atlas migrate diff --env $(ENV) \
		--var "db_password=$$($(MAKE) -s db-password)"
