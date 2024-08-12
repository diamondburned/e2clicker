#!/usr/bin/env bash

BASE_DIR=$(dirname "${BASH_SOURCE[0]}")
GO_MODULE=$(go list -m)
GO_OPENAPI_PKG="${GO_MODULE}/openapi"
INITIALISMS=(
	API
	URL
	ID
)

main() {
	cd "$BASE_DIR"
	
	importMappings=()
	declare -A packages
	
	for name in api*.yml; do
		base="${name%.yml}"

		importMappings+=( "./${name}:${GO_OPENAPI_PKG}/${base}" )
		packages["$name"]="$base"
	done
	
	printf -v importMapping "%s," "${importMappings[@]}"
	importMapping="${importMapping%,}"

	printf -v initialisms "%s," "${INITIALISMS[@]}"
	initialisms="${initialisms%,}"
	
	for name in api*.yml; do
		base="${packages["$name"]}"

		mkdir -p "$base"

		oapi-codegen \
			--package "${base}" \
			--generate types,chi-server,strict-server \
			--alias-types \
			--import-mapping "${importMapping}" \
			--response-type-suffix Response \
			-o "${base}/${base}.gen.go" \
			"${name}"
	done
}

main
