#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec env PROTOCOL=http3 bash "${SCRIPT_DIR}/run-linux-quic-server.sh" "$@"
