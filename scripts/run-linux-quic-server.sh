#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

FRAMEWORK="${1:-}"
PROTOCOL="${PROTOCOL:-http3}"
SERVER_ADDR="${SERVER_ADDR:-0.0.0.0:19093}"
PAYLOAD="${PAYLOAD:-128}"
SERVER_CPU_SET="${SERVER_CPU_SET:-0,1,2,3}"
SERVER_GOMAXPROCS="${SERVER_GOMAXPROCS:-4}"
EVENT_LOOPS="${EVENT_LOOPS:-3}"
GNALLOY_BENCH="${GNALLOY_BENCH:-${REPO_ROOT}/external/bin/gnalloy-bench}"
NETTY_BENCH_JAR="${NETTY_BENCH_JAR:-${REPO_ROOT}/external/bin/netty-bench.jar}"
JAVA_BIN="${JAVA_BIN:-/opt/software/java21/bin/java}"
SERVER_PID_FILE="${SERVER_PID_FILE:-/tmp/gnalloy-http3-cross-host.pid}"

case "${PROTOCOL}" in
  http3)
    ALPN="h3"
    ;;
  quic-stream)
    ALPN="gnalloy-quic"
    ;;
  *)
    printf 'unsupported QUIC protocol: %s\n' "${PROTOCOL}" >&2
    exit 1
    ;;
esac

require_positive_integer() {
  local name="$1"
  local value="$2"
  if [[ ! "${value}" =~ ^[1-9][0-9]*$ ]]; then
    printf '%s must be a positive integer: %s\n' "${name}" "${value}" >&2
    exit 1
  fi
}

require_executable() {
  local path="$1"
  if [[ ! -x "${path}" ]]; then
    printf 'benchmark executable is missing: %s\n' "${path}" >&2
    exit 1
  fi
}

require_positive_integer PAYLOAD "${PAYLOAD}"
require_positive_integer SERVER_GOMAXPROCS "${SERVER_GOMAXPROCS}"
require_positive_integer EVENT_LOOPS "${EVENT_LOOPS}"
taskset -c "${SERVER_CPU_SET}" true
if [[ -e "${SERVER_PID_FILE}" ]]; then
  printf 'server PID file already exists: %s\n' "${SERVER_PID_FILE}" >&2
  exit 1
fi
printf '%s\n' "$$" >"${SERVER_PID_FILE}"

case "${FRAMEWORK}" in
  gnalloy)
    require_executable "${GNALLOY_BENCH}"
    exec taskset -c "${SERVER_CPU_SET}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${GNALLOY_BENCH}" \
      -protocol "${PROTOCOL}" -server-only=true -addr "${SERVER_ADDR}" -payload "${PAYLOAD}" \
      -tls-version 1.3 -alpn "${ALPN}" -timeout 5m
    ;;
  netty)
    if [[ ! -f "${NETTY_BENCH_JAR}" ]]; then
      printf 'benchmark file is missing: %s\n' "${NETTY_BENCH_JAR}" >&2
      exit 1
    fi
    require_executable "${JAVA_BIN}"
    exec taskset -c "${SERVER_CPU_SET}" "${JAVA_BIN}" -XX:ActiveProcessorCount="${SERVER_GOMAXPROCS}" \
      -jar "${NETTY_BENCH_JAR}" --protocol "${PROTOCOL}" --server-only true --addr "${SERVER_ADDR}" \
      --payload "${PAYLOAD}" --backend epoll --event-loops "${EVENT_LOOPS}" \
      --tls-version 1.3 --alpn "${ALPN}" --timeout 5m
    ;;
  *)
    printf 'unsupported HTTP/3 framework: %s\n' "${FRAMEWORK}" >&2
    exit 1
    ;;
esac
