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
# TODO: figure out how to pass this into the dev VM

export BACKEND_HTTP_ADDRESS := "http://localhost:8000"

dev:
    nix run .#e2clicker-dev

dev-vm: generate
    nix run .#nixosConfigurations.dev-vm.config.system.build.nixos-shell

dev-backend: generate-backend
    go run ./cmd/e2clicker --port 8000

dev-frontend: generate-frontend
    vite --port 8080

###

generate: openapi generate-backend generate-frontend generate-docs

[private]
generate-go:
	go generate -x ./...

[private]
generate-backend: openapi-backend generate-go
    go mod tidy
    go mod download
    gomod2nix --outdir nix

[private]
generate-frontend: openapi-frontend

[private]
generate-docs: openapi-docs

###

openapi: openapi-schema openapi-backend openapi-frontend openapi-docs

[private]
openapi-schema:
    generate-openapi schema

[private]
openapi-backend: openapi-schema
    generate-openapi backend

[private]
openapi-frontend: openapi-schema
    pnpm i --prefer-offline --prefer-frozen-lockfile --use-stderr
    generate-openapi frontend

[private]
openapi-docs: openapi-schema
    pnpm i --prefer-offline --prefer-frozen-lockfile --use-stderr
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
