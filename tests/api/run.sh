#!/usr/bin/env bash
set -euo pipefail

export TZ=UTC
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
WORK="$(mktemp -d "${TMPDIR:-/tmp}/e2clicker-api-test.XXXXXX")"
PG_PID=""
BACKEND_PID=""

cleanup() {
	set +e
	if [ -n "$BACKEND_PID" ]; then
		kill "$BACKEND_PID" 2>/dev/null
		wait "$BACKEND_PID" 2>/dev/null
	fi
	if [ -d "$WORK/pgdata" ]; then
		pg_ctl -D "$WORK/pgdata" -m immediate stop >/dev/null 2>&1 || true
	fi
	rm -rf "$WORK"
}
trap cleanup EXIT

mkdir -p "$WORK/pgdata" "$WORK/pgsock" "$WORK/bin"

initdb -D "$WORK/pgdata" --username=postgres --auth=trust --no-sync >"$WORK/initdb.log" 2>&1
pg_ctl -D "$WORK/pgdata" -l "$WORK/pg.log" -o "-k $WORK/pgsock -h '' -p 5432" start >/dev/null
createdb -h "$WORK/pgsock" -U postgres e2clicker

(
	cd "$ROOT"
	go build -o "$WORK/bin/e2clicker-backend" ./cmd/e2clicker-backend
	go build -o "$WORK/bin/vapid-generate" ./cmd/vapid-generate
)

"$WORK/bin/vapid-generate" -subscriber "mailto:test@e2clicker.local" >"$WORK/vapid-keys.json"

PORT="$(node -e 'const net = require("net"); const s = net.createServer(); s.listen(0, "127.0.0.1", () => { console.log(s.address().port); s.close(); });')"

cat >"$WORK/config.json" <<EOF
{
  "enable": true,
  "package": "local",
  "api": {
    "listenAddress": "127.0.0.1:${PORT}",
    "debugRequests": false
  },
  "postgresql": {
    "databaseURI": "postgres://postgres@/e2clicker?host=${WORK}/pgsock"
  },
  "notification": {
    "clientTimeout": "5s",
    "webPush": "${WORK}/vapid-keys.json",
    "email": null
  },
  "debug": true,
  "logFormat": "text"
}
EOF

"$WORK/bin/e2clicker-backend" -c "$WORK/config.json" >"$WORK/backend.log" 2>&1 &
BACKEND_PID="$!"

ready=0
for _ in $(seq 1 50); do
	if node -e "fetch('http://127.0.0.1:${PORT}/api/delivery-methods').then(r => process.exit(r.ok ? 0 : 1)).catch(() => process.exit(1))" 2>/dev/null; then
		ready=1
		break
	fi
	if ! kill -0 "$BACKEND_PID" 2>/dev/null; then
		echo "e2clicker-backend exited unexpectedly:" >&2
		cat "$WORK/backend.log" >&2
		exit 1
	fi
	sleep 0.2
done

if [ "$ready" -ne 1 ]; then
	echo "e2clicker-backend did not become ready in time:" >&2
	cat "$WORK/backend.log" >&2
	exit 1
fi

if ! hurl --test --variable "base_url=http://127.0.0.1:${PORT}" "$ROOT"/tests/api/*.hurl; then
	echo "--- backend log ---" >&2
	tail -n 50 "$WORK/backend.log" >&2
	exit 1
fi
