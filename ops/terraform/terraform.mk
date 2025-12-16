ENV ?=
CI  ?=

# ifeq ($(ENV),prd)
# 	else
# endif

.DEFAULT_GOAL := help

.PHONY: fmt
fmt:
	terraform fmt -recursive

.PHONY: check
check: .check-env
	terraform fmt -recursive
	@cd $(ENV) && \
		terraform init -upgrade && \
		terraform validate

.PHONY: plan
plan: .check-env check
	@cd $(ENV) && terraform plan -lock=false --parallelism=30

.PHONY: apply
apply: .check-env check
	@cd $(ENV) && terraform apply -auto-approve

.PHONY: .check_env
.check-env:
ifndef ENV
	$(error ENV is required.)
endif
