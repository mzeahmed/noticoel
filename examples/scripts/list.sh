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

LIMIT="${1:-20}"
OFFSET="${2:-0}"

echo "Listing events (limit=${LIMIT}, offset=${OFFSET})"

PRETTY=(cat)
command -v jq >/dev/null 2>&1 && PRETTY=(jq .)

curl \
    --fail \
    --silent \
    --show-error \
    --request GET \
    --url "${NOTICOEL_URL}/api/v1/events/list?limit=${LIMIT}&offset=${OFFSET}" \
    --header "Authorization: Bearer ${NOTICOEL_TOKEN}" \
    | "${PRETTY[@]}"
