# Include variables
include .envrc

## help: Print this help message
.PHONY: help
help:
	@echo 'Usage':
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# =============================================================================== #
# DEVELOPMENT
# =============================================================================== #

## run/api: Run the application
.PHONY: run/api
run/api:
	go run ./cmd/main

# =============================================================================== #
# QUALITY CONTROL
# =============================================================================== #

## tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...

	@echo 'Tidying root module dependencies...'
	(cd ./ && go mod tidy)
	@echo 'Verifying root module dependencies...'
	(cd ./ && go mod verify)

	@echo 'Tidying application module dependencies...'
	(cd ./application && go mod tidy)
	@echo 'Verifying application module dependencies...'
	(cd ./application && go mod verify)

	@echo 'Tidying domain module dependencies...'
	(cd ./domain && go mod tidy)
	@echo 'Verifying domain module dependencies...'
	(cd ./domain && go mod verify)

	@echo 'Tidying infrastructure module dependencies...'
	(cd ./infrastructure && go mod tidy)
	@echo 'Verifying infrastructure module dependencies...'
	(cd ./infrastructure && go mod verify)

	@echo 'Tidying interfaces module dependencies...'
	(cd ./interfaces && go mod tidy)
	@echo 'Verifying interfaces module dependencies...'
	(cd ./interfaces && go mod verify)

	@echo 'Tidying cmd module dependencies...'
	(cd ./cmd && go mod tidy)
	@echo 'Verifying cmd module dependencies...'
	(cd ./cmd && go mod verify)

	@echo 'Vendoring workspace dependencies...'
	go work vendor

# =============================================================================== #
# BUILD
# =============================================================================== #

## build/api: Build application without optimizations
.PHONY: build/api
build/api:
	@echo 'Building application without optimizations...'
	@mkdir -p ./bin
	GOARCH=amd64 GOOS=linux go build -o=./bin/api ./cmd/main
	@echo 'Build for Linux (amd64) complete.'

## build/api/optimized: Build application with optimizations
.PHONY: build/api/optimized
build/api/optimized:
	@echo 'Building application with optimizations...'
	@mkdir -p ./bin
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -ldflags="-s -w" -o=./bin/linux_amd64/api ./cmd/main
	@echo 'Build for Linux (amd64) complete (optimized).'