default:
	just --list

build: build-backend build-frontend

build-backend: generate
	@mkdir -p dist/backend
	@echo "Building Go backend..."
	go build -o ./dist/backend/hrtclicker ./cmd/hrtclicker

build-frontend:
	@mkdir -p dist
	@echo "Building Svelte frontend..."
	vite build --logLevel error

dev-frontend:
	vite

test: generate
	go test -v ./...

generate: format
	go generate ./...

format:
	@nixfmt $(find . -name '*.nix')
	@goimports -w $(go list -f '{{{{.Dir}}' ./...)
