ENV         ?= dev
GCP_PROJECT := mahking-$(ENV)
REGION      := asia-northeast1
SERVICE     := mahking-$(ENV)-api
REGISTRY    := $(REGION)-docker.pkg.dev/$(GCP_PROJECT)/mahking-$(ENV)
IMAGE       := $(REGISTRY)/api
GIT_SHA     := $(shell git rev-parse --short HEAD)
TAG         := $(IMAGE):$(GIT_SHA)

.PHONY: govalid sqlc vet lint docker-login build deploy release

govalid: ## govalid のバリデーター生成
	go tool govalid ./...

sqlc: ## sqlc のクエリコード生成
	sqlc generate

vet: ## 静的解析
	go vet ./...

lint: ## golangci-lint 実行
	go tool golangci-lint run ./...

docker-login: ## Artifact Registry へログイン
	gcloud auth configure-docker $(REGION)-docker.pkg.dev --quiet

build: docker-login ## イメージビルド & push
	docker build -t $(TAG) .
	docker push $(TAG)

deploy: ## Cloud Run へデプロイ
	gcloud run deploy $(SERVICE) \
		--project $(GCP_PROJECT) \
		--region $(REGION) \
		--image $(TAG)

release: build deploy ## ビルド + デプロイ
