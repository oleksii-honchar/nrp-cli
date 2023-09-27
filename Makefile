SHELL=/bin/bash
RED=\033[0;31m
GREEN=\033[0;32m
BG_GREY=\033[48;5;237m
YELLOW=\033[38;5;202m
NC=\033[0m # No Color
BOLD_ON=\033[1m
BOLD_OFF=\033[21m
CLEAR=\033[2J

LATEST_VERSION := $(shell cat ./pkg/latest-version/latest-version.txt)

.PHONY: help

help:
	@echo "nrp-cli" automation commands:
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run: ## run
	@go run main.go -log-level=debug -config=./configs/nrp.yaml

latest-version:
	@echo "latest-version.txt = $(LATEST_VERSION)"

build-all: build-linux build-darwin # build darwin & linux

build-linux: ## build for linux
	GOOS=linux CGO_ENABLED=0 go build -o nrp-cli-linux .

build-darwin: ## build for linux
	GOOS=darwin go build -o nrp-cli-darwin .

go-publish: ## publish blablo
	GOPROXY=proxy.golang.org go list -m github.com/oleksii-honchar/nrp-cli@$(LATEST_VERSION)

.ONESHELL:
run-all: ## run nrp-cli -> nrp
	@go run main.go
	@docker compose down
	@docker compose up --build --remove-orphans -d
	@docker compose logs --follow


.ONE-SHELL:
run-nginx: stop-nginx ## run nginx
	@docker run -d --rm -p 80:80 -p 443:443 \
		--name nginx-reverse-proxy \
		-v ./nginx-config:/etc/nginx \
		-v /etc/localtime:/etc/localtime:ro \
		-v ./letsencrypt:/etc/letsencrypt \
		tuiteraz/nginx-reverse-proxy:1.0;\
	docker logs nginx-reverse-proxy

stop-nginx: ## stop nginx
	@docker stop nginx-reverse-proxy || true

test-nginx: ## test nginx
	nginx -t -c $(PWD)/nginx-config/nginx.conf