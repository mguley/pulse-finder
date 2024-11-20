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

## install/goimports: Install goimports for formatting
.PHONY: install/goimports
install/goimports:
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest

## install/linter: Install GolangCI-Lint
.PHONY: install/linter
install/linter:
	@echo "Installing GolangCI-Lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

## lint: Run linter on all Go files in each module directory
.PHONY: lint
lint: install/linter
	@echo "Running GolangCI-Lint on all Go files in each module directory..."
	@find ./application ./cmd ./domain ./infrastructure ./interfaces -name '*.go' -exec dirname {} \; | sort -u | xargs $(shell go env GOPATH)/bin/golangci-lint run

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

	@echo 'Tidying tests module dependencies...'
	(cd ./tests && go mod tidy)
	@echo 'Verifying cmd module dependencies...'
	(cd ./tests && go mod verify)

	@echo 'Vendoring workspace dependencies...'
	go work vendor

# =============================================================================== #
# TESTING
# =============================================================================== #

## test/integration: Run integration tests (uses Go's caching mechanism)
.PHONY: test/integration
test/integration:
	@echo 'Running integration tests (with caching, sequentially)...'
	go test -v -p=1 ./tests/integration/...

## test/integration/no-cache: Run integration tests (bypass cache)
.PHONY: test/integration/no-cache
test/integration/no-cache:
	@echo 'Running integration tests (no cache, sequentially)...'
	go test -v -count=1 -p=1 ./tests/integration/...

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

# =============================================================================== #
# PRODUCTION DEPLOYMENT TASKS
# =============================================================================== #

## production/connect: Connect to the production server
.PHONY: production/connect
production/connect:
	ssh pulsefinder@${PRODUCTION_HOST_IP}

## production/tunnel: Create SSH tunnel from port 4000 on the production server to port 9999 locally
.PHONY: production/tunnel
production/tunnel:
	@echo 'Creating SSH tunnel to the production server...'
	ssh -L 9999:${PRODUCTION_HOST_IP}:${PORT} pulsefinder@${PRODUCTION_HOST_IP}

## production/deploy-files: Deploy new binary and service files
.PHONY: production/deploy-files
production/deploy-files:
	@echo 'Deploying new api binary and service files...'
	rsync -P ./bin/linux_amd64/api pulsefinder@${PRODUCTION_HOST_IP}:/tmp/api
	rsync -P ./remote/production/service/api.service pulsefinder@${PRODUCTION_HOST_IP}:/tmp/api.service
	rsync -P ./remote/production/service/nats.service pulsefinder@${PRODUCTION_HOST_IP}:/tmp/nats.service
	rsync -rP --delete ./infrastructure/migrations pulsefinder@${PRODUCTION_HOST_IP}:~
	ssh -t pulsefinder@${PRODUCTION_HOST_IP} 'set -e; \
	  sudo mkdir -p /opt/pulsefinder && \
	  sudo mv /tmp/api /opt/pulsefinder/ && \
	  sudo chown -R pulsefinder:pulsefinder /opt/pulsefinder && \
	  sudo chmod +x /opt/pulsefinder/api && \
	  sudo mv /tmp/api.service /etc/systemd/system/ && \
	  sudo mv /tmp/nats.service /etc/systemd/system/'

## production/migrate: Run database migrations
.PHONY: production/migrate
production/migrate:
	@echo 'Running database migrations on production...'
	ssh -t pulsefinder@${PRODUCTION_HOST_IP} 'set -e; \
	  migrate -path ~/migrations/jobs -database $$DB_DSN up'

## production/restart-service: Restart the API and NATS services
.PHONY: production/restart-service
production/restart-service:
	@echo 'Restarting the api and nats services...'
	ssh -t pulsefinder@${PRODUCTION_HOST_IP} 'set -e; \
	  sudo systemctl daemon-reload && \
	  sudo systemctl enable api && \
	  sudo systemctl restart api && \
	  sudo systemctl enable nats && \
	  sudo systemctl restart nats'

## production/check-status: Check the status of the API service
.PHONY: production/check-service-status
production/check-status:
	@echo 'Checking the api service status...'
	ssh -t pulsefinder@${PRODUCTION_HOST_IP} 'set -e; \
	  sudo systemctl status api --no-pager'

## production/check-nginx-status: Check the status of the Nginx service
.PHONY: production/check-nginx-status
production/check-nginx-status:
	@echo 'Checking the Nginx service status...'
	ssh -t pulsefinder@${PRODUCTION_HOST_IP} 'set -e; \
	  sudo systemctl status nginx --no-pager'

## production/upload-nginx-config: Upload Nginx config files and restart Nginx
.PHONY: production/upload-nginx-config
production/upload-nginx-config:
	@echo 'Uploading Nginx configuration files and restarting Nginx...'
	rsync -P ./remote/production/nginx/nginx.conf pulsefinder@${PRODUCTION_HOST_IP}:~
	rsync -P ./remote/production/nginx/api-nginx.conf pulsefinder@${PRODUCTION_HOST_IP}:~
	ssh -t pulsefinder@${PRODUCTION_HOST_IP} 'set -e; \
	  sudo mv ~/nginx.conf /etc/nginx/nginx.conf && \
	  sudo mv ~/api-nginx.conf /etc/nginx/conf.d/api-nginx.conf && \
	  sudo systemctl restart nginx'

## production/deploy/api: Deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	@$(MAKE) build/api/optimized
	@$(MAKE) production/deploy-files
	@$(MAKE) production/migrate
	@$(MAKE) production/restart-service
	@$(MAKE) production/upload-nginx-config
	@$(MAKE) production/check-service-status
	@$(MAKE) production/check-nginx-status
	@echo 'Deployment to production complete.'
