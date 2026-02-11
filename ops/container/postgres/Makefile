PG_MAJOR        := 18

GCP_PROJECT  := mahking-dev
IMAGE_NAME   := asia-northeast1-docker.pkg.dev/$(GCP_PROJECT)/mahking-dev/postgres
IMAGE_PATH   = $(IMAGE_NAME):$(PG_MAJOR)

.DEFAULT_GOAL := help

.PHONY: docker-login
docker-login: ## Artifact Registryへdocker loginを実行する
	docker login -u oauth2accesstoken -p "$$(gcloud auth print-access-token)" https://asia-northeast1-docker.pkg.dev

.PHONY: build
build: ## コンテナイメージをビルドし、Artifact Registryへプッシュする
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE_PATH) \
		--build-arg PG_MAJOR=$(PG_MAJOR) \
		--push \
		.

.PHONY: help
help: ## ヘルプ
	@grep -E '^[0-9a-zA-Z_/()$$-]+:.*?## .*$$' $(lastword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
