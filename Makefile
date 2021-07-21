.DEFAULT_GOAL := help

# Load .env. see https://lithic.tech/blog/2020-05/makefile-dot-env
ifneq (,$(wildcard ./.env))
  include .env
  export
endif

# Environment
GO111MODULE := on
PATH := $(CURDIR)/.go-tools/bin:$(PATH)
SHELL := bash

# This is a magic code to output help message at default
# see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[/0-9a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-35s\033[0m %s\n", $$1, $$2}'

.PHONY: format
format: ## Run `go fmt`
	@go fmt $$(go list ./... | grep -v vendor)

.PHONY: test
test: ## Run all tests
	@go test -cover $$(go list ./... | grep -v vendor)

.PHONY: testv
testv: ## Run all tests with verbose outputing.
	@go test -v -cover $$(go list ./... | grep -v vendor)

.PHONY: bindata
bindata: ## bindata
	@go-bindata -o ./assets/bindata.go -ignore bindata.go -prefix assets -pkg assets assets/...

.PHONY: serve
serve: bindata ## Run dev process
	@go run main.go

.PHONY: build-bundle
build-bundle: bindata ## Build macOS Application bundle
	@./build-bundle

.PHONY: clean
clean: ## clean built outputs
	@rm -rf dist

.PHONY: testcov
testcov: ## Run all tests and outputs coverage report.
	@gocov test $$(go list ./... | grep -v vendor) | gocov-html > coverage-report.html

.PHONY: deps
deps: ## Install dependences.
	@go mod tidy

.PHONY: installtools
installtools: ## Install dev tools
	@export GO111MODULE=off && GOPATH=$(CURDIR)/.go-tools && \
      go get -u github.com/mattn/go-bindata/... && \
      go get -u github.com/axw/gocov/gocov && \
      go get -u gopkg.in/matm/v1/gocov-html
	rm -rf $(CURDIR)/.go-tools/pkg
	rm -rf $(CURDIR)/.go-tools/src

.PHONY: cleantools
cleantools:
	@export GO111MODULE=off && GOPATH=$(CURDIR)/.go-tools && rm -rf $(CURDIR)/.go-tools



