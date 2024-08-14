export PATH := x'${PATH}:${PWD}/scripts:${PWD}/node_modules/.bin'

###

default:
	@just -l

###

build: build-backend build-frontend generate

build-backend: build-clean generate-backend
    @mkdir -p dist/backend
    go build -o ./dist/backend/e2clicker-backend ./cmd/e2clicker-backend

build-frontend: build-clean build-frontend-fix generate-frontend
    @mkdir -p dist/frontend
    vite build --logLevel error

# Fix a race condition with Vite running before Sveltekit can create its
# tsconfig.json.
[private]
build-frontend-fix:
    @mkdir .svelte-kit 2> /dev/null && echo "{}" > .svelte-kit/tsconfig.json || true

[private]
build-clean:
    @rm -rf dist

###

# dev:
#     exit 1 # Not implemented yet

dev-backend *FLAGS: generate-backend
    go run ./cmd/e2clicker {{ FLAGS }}

dev-frontend *FLAGS: generate-frontend
    vite {{ FLAGS }}

###

generate: openapi generate-backend generate-frontend generate-docs

[private]
generate-backend: openapi-backend
    go generate -x ./...
    go mod tidy
    go mod download

[private]
generate-frontend: openapi-frontend
    pnpm i

[private]
generate-docs: openapi-docs

###

openapi: openapi-schema openapi-backend openapi-frontend openapi-docs

[private]
openapi-schema:
    generate-openapi schema

[private]
openapi-backend:
    generate-openapi backend

[private]
openapi-frontend:
    generate-openapi frontend

[private]
openapi-docs:
    generate-openapi docs &> /dev/null

###

test: test-backend

[private]
test-backend: generate-backend
    go test -v ./...

[private]
test-frontend: generate-frontend
    # jest

###

format:
    @nixfmt $(find . -name '*.nix')
    @prettier --log-level warn -w .
    @goimports -w $(go list -f '{{{{.Dir}}' ./...)
    @just --unstable --fmt 2> /dev/null
