export PATH := x'${PATH}:${PWD}/node_modules/.bin'

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

dev: generate
    zellij delete-session e2clicker || true
    zellij --session e2clicker --new-session-with-layout ./nix/dev/zellij-layout.kdl

[private]
dev-vm: generate
    # Force a path to the flake instead of a git+path to include the
    # vapid-keys.json file, which is gitignore'd.
    nix run "path://$PWD"'#nixosConfigurations.dev-vm.config.system.build.nixos-shell'

[private]
dev-backend:
    go run ./cmd/e2clicker --port 8000

[private]
dev-frontend:
    vite --port 8080

###

generate: openapi generate-backend generate-frontend generate-docs generate-vapid

[private]
generate-go:
    go generate -x ./...

[private]
generate-backend: openapi-backend generate-go generate-backend-config
    go mod tidy
    go mod download
    gomod2nix --outdir nix

[private]
generate-backend-config:
    #!/bin/sh
    flakePath=$(jq -n --arg path "$PWD" '$path')
    nixmod2go \
        -P e2clickermodule \
        -T BackendConfig \
        -O services.e2clicker.backend \
        -c ./nix/modules/nixmod2go.json \
        .#nixosModules.e2clicker \
        ./nix/modules/e2clicker/config.go

[private]
generate-frontend: openapi-frontend

[private]
generate-docs: openapi-docs

[private]
generate-vapid:
    #!/bin/sh
    if [ ! -f vapid-keys.json ]; then
        go run ./cmd/vapid-generate > vapid-keys.json
    fi

###

openapi: openapi-schema openapi-backend openapi-frontend openapi-docs

[private]
openapi-schema:
    ./openapi/generate.sh schema

[private]
openapi-backend: openapi-schema
    ./openapi/generate.sh backend

[private]
openapi-frontend: openapi-schema pnpm
    ./openapi/generate.sh frontend

[private]
openapi-docs: openapi-schema pnpm
    ./openapi/generate.sh docs &> /dev/null

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

###

[private]
pnpm:
    pnpm i --prefer-offline --prefer-frozen-lockfile --use-stderr
