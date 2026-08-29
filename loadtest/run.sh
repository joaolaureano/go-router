#!/usr/bin/env bash
# Drives the vegeta CLI against the load server: starts it, runs every target
# file, prints a report per scenario, and shuts down.
#
# The in-process suite (go test ./loadtest/...) is the one that asserts; this
# is for longer runs, higher rates, and vegeta's own reporters -- histograms,
# HDR plots, JSON for a dashboard.
#
#   ./run.sh                 # every scenario at the default rate
#   RATE=20000 DURATION=30s ./run.sh
#   ./run.sh mixed static    # only these scenarios
#   RATE=0 ./run.sh catalog  # burst: no pacing, every route in the generated table
set -euo pipefail

cd "$(dirname "$0")"

RATE="${RATE:-5000}"
DURATION="${DURATION:-10s}"
ADDR="${ADDR:-:8080}"
PORT="${ADDR##*:}"
OUT="${OUT:-results}"

command -v vegeta >/dev/null || {
  echo "vegeta not found. Install it with:" >&2
  echo "  go install github.com/tsenart/vegeta/v12@latest" >&2
  echo "  # or: brew install vegeta" >&2
  exit 1
}

scenarios=("$@")
if [ ${#scenarios[@]} -eq 0 ]; then
  scenarios=(static params catchall middleware notfound methods mixed)
fi

mkdir -p "$OUT"

# "catalog" is the whole generated table: three API versions of nested
# subroutes, mounted subtrees and all. Its target file is generated from the
# same catalog the server registers, so it can neither drift nor miss a route,
# and the server has to be started with -catalog to actually serve it.
serve_catalog=""
for scenario in "${scenarios[@]}"; do
  if [ "$scenario" = "catalog" ]; then serve_catalog="-catalog"; fi
done

echo "building load server..."
go build -o "$OUT/loadserver" ./cmd/loadserver

if [ -n "$serve_catalog" ]; then
  "$OUT/loadserver" -dump-targets targets/catalog.txt -base "http://localhost:${PORT}"
fi

"$OUT/loadserver" -addr "$ADDR" $serve_catalog &
server_pid=$!
trap 'kill "$server_pid" 2>/dev/null || true' EXIT

# Wait for the socket rather than sleeping a guessed amount: a warm-up request
# that misses would be recorded as a routing failure.
probe="/ping"
[ -n "$serve_catalog" ] && probe="/api/v1/users"
for _ in $(seq 1 50); do
  if curl -sf "http://localhost:${PORT}${probe}" >/dev/null; then break; fi
  sleep 0.1
done
curl -sf "http://localhost:${PORT}${probe}" >/dev/null || { echo "server did not come up" >&2; exit 1; }

for scenario in "${scenarios[@]}"; do
  targets="targets/${scenario}.txt"
  [ -f "$targets" ] || { echo "no such scenario: $scenario ($targets)" >&2; exit 1; }

  echo
  if [ "$RATE" = "0" ]; then
    echo "=== $scenario  (burst: unpaced, ${DURATION}, ${MAX_WORKERS:-200} connections) ==="
  else
    echo "=== $scenario  (${RATE}/s for ${DURATION}) ==="
  fi
  vegeta attack \
    -targets="$targets" \
    -rate="$RATE" \
    -duration="$DURATION" \
    -keepalive=true \
    -max-workers="${MAX_WORKERS:-200}" \
    > "$OUT/${scenario}.bin"

  vegeta report "$OUT/${scenario}.bin"
  vegeta report -type=json "$OUT/${scenario}.bin" > "$OUT/${scenario}.json"
  vegeta plot "$OUT/${scenario}.bin" > "$OUT/${scenario}.html"
done

echo
echo "raw results, JSON reports and HTML plots in $OUT/"
