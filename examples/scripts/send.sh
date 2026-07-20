#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="${SCRIPT_DIR}/../../.env"

if [ -f "$ENV_FILE" ]; then
    set -a
    # shellcheck disable=SC1090
    source "$ENV_FILE"
    set +a
fi

NOTICOEL_URL="${NOTICOEL_URL:-http://localhost:8080}"
NOTICOEL_TOKEN="${NOTICOEL_TOKEN:-${NOTICOEL_AUTH_TOKEN:-change-me}}"

EVENT_FILE="${1:-../events/workflow-success.json}"

if [ ! -f "$EVENT_FILE" ]; then
    echo "Event file not found: $EVENT_FILE"
    echo ""
    echo "Usage:"
    echo "  ./send.sh ../events/workflow-success.json"
    echo "  ./send.sh ../events/workflow-failure.json"
    echo "  ./send.sh ../events/release.json"
    exit 1
fi

echo "Sending event: $EVENT_FILE"

PRETTY=(cat)
command -v jq >/dev/null 2>&1 && PRETTY=(jq .)

curl \
    --fail \
    --silent \
    --show-error \
    --request POST \
    --url "${NOTICOEL_URL}/api/v1/events/create" \
    --header "Authorization: Bearer ${NOTICOEL_TOKEN}" \
    --header "Content-Type: application/json" \
    --data "@${EVENT_FILE}" \
    | "${PRETTY[@]}"

echo "✓ Event sent successfully"
