.DEFAULT_GOAL := help
NAMESPACE=default

build: ## Build docker image with tag mainak90/job-admission:latest
	@docker build . --tag mainak90/job-admission:latest

certificates: ## Generate certs/manifest.yaml with the ValidatingWebhookConfiguration and a randomly generated cert
	@./generate_certs.sh job-admission $(NAMESPACE)

kind: build certificates ## Build image and upload to king
	@kind load docker-image mainak90/job-admission:latest --name kind

skaffold: certificates ## Generate certificates and start skaffold
	@skaffold dev

help:
	@echo "Use make NAMESPACE=override to change the target namespace\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'