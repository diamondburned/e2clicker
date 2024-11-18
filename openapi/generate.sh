#!/usr/bin/env bash
set -e

main() {
	f=openapi::"${1:-all}"
	if command -v "$f" > /dev/null; then
		"$f"
	else
		echo "Usage: $0 {schema|backend|frontend|docs}" >&2
		exit 1
	fi
}

openapi::schema() {
	pass1=$(mktemp -p openapi ".tmp.XXXXXXXXXX.json.tmp")
	pass2=$(mktemp -p openapi ".tmp.XXXXXXXXXX.json.tmp")

	mergeSchema openapi/_base.yml > "$pass1"
	mergeSchema "$pass1" openapi/_*.yml > "$pass2"

	redocly bundle -o openapi/openapi.gen.json "$pass2" --ext json

	rm "$pass1" "$pass2"
}

mergeSchema() {
	schemas=()
	for path in "$@"; do
		name="$(basename "$path" .yml)"
		if [[ $name == _* ]]; then
			name=${name/_/}
			schema="$(yq e "(.paths // (.paths = {}))[].[].tags += [\"$name\"]" "$path")"
		else
			schema="$(cat "$path")"
		fi
		schemas+=( "$schema" )
	done

	printf '%s\n---\n' "${schemas[@]}" | \
		yq -pauto -ojson ea '. as $item ireduce ({}; . *n+ $item )' -
}

log() {
	echo "$@" >&2
}

openapi::backend() {
	# Special case for the API package.
	#
	# TODO: In the future, it would be more ideal if we can decouple
	# services/api/handler.go into multiple packages like:
	#   - services/dosage/openapi/handler.go
	#   - services/notification/openapi/handler.go
	#   - services/user/openapi/handler.go
	oapi-codegen \
		-config ./openapi/oapi-codegen-api.yml \
		-import-mapping "" \
		-o ./services/api/openapi/openapi.gen.go \
		./openapi/openapi.gen.json

	for oapiFile in ./openapi/_*.yml; do
		service=$(basename "$oapiFile" .yml)
		service=${service/_/}

		genPath="./services/$service/openapi/openapi.gen.go"
		mkdir -p "$(dirname "$genPath")"

		oapi-codegen \
			-config ./openapi/oapi-codegen-models.yml \
			-o "$genPath" \
			"$oapiFile"
	done
}

openapi::frontend() {
	pnpx oazapfts --optimistic \
		./openapi/openapi.gen.json \
		./frontend/lib/openapi.gen.ts
}

openapi::docs() {
	pnpx widdershins \
		--code \
		--summary \
		--omitHeader \
		--shallowSchemas \
		./openapi/openapi.gen.json \
		./docs/openapi/README.md
}

main "$@"
