default:
	just --list

###

build: build-backend build-frontend

build-backend: openapi-backend modules generate
	@mkdir -p dist/backend
	@echo "Building Go backend..."
	go build -o ./dist/backend/e2clicker ./cmd/e2clicker

build-frontend: openapi-frontend
	@mkdir -p dist
	@echo "Building Svelte frontend..."
	pnpx vite build --logLevel error

###

dev:
	exit 1 # Not implemented yet

dev-frontend:
	pnpx vite

###

openapi: \
	openapi-schema \
	openapi-backend \
	openapi-frontend \
	openapi-docs \
	modules

openapi-schema:
	./internal/gen-openapi schema

openapi-backend:
	./internal/gen-openapi backend

openapi-frontend:
	./internal/gen-openapi frontend

openapi-docs:
	./internal/gen-openapi docs &> /dev/null

###

test: test-backend

test-backend: modules generate
	go test -v ./...

###

generate: openapi format
	go generate -x ./...

format:
	@nixfmt $(find . -name '*.nix')
	@goimports -w $(go list -f '{{{{.Dir}}' ./...)

modules:
	go mod tidy
	go mod download
