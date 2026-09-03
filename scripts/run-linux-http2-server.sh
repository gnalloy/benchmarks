#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

FRAMEWORK="${1:-}"
PROTOCOL="${PROTOCOL:-http2}"
SERVER_ADDR="${SERVER_ADDR:-0.0.0.0:19092}"
PAYLOAD="${PAYLOAD:-128}"
TLS_VERSION="${TLS_VERSION:-1.3}"
SERVER_CPU_SET="${SERVER_CPU_SET:-0,1,2,3}"
SERVER_GOMAXPROCS="${SERVER_GOMAXPROCS:-4}"
EVENT_LOOPS="${EVENT_LOOPS:-3}"
GNALLOY_WORKERS="${GNALLOY_WORKERS:-3}"
GNALLOY_BOSS_CPU_SET="${GNALLOY_BOSS_CPU_SET:-3}"
GNALLOY_WORKER_CPU_SET="${GNALLOY_WORKER_CPU_SET:-0,1,2}"
GNALLOY_READ_BUFFER_SIZE="${GNALLOY_READ_BUFFER_SIZE:-4096}"
GNALLOY_MAX_MESSAGES_PER_READ="${GNALLOY_MAX_MESSAGES_PER_READ:-32}"
GNALLOY_CPU_PROFILE="${GNALLOY_CPU_PROFILE:-}"
SERVER_PID_FILE="${SERVER_PID_FILE:-/tmp/gnalloy-http2-cross-host.pid}"
GNALLOY_BENCH="${GNALLOY_BENCH:-${REPO_ROOT}/external/bin/gnalloy-bench}"
HERTZ_BENCH="${HERTZ_BENCH:-${REPO_ROOT}/external/bin/hertz-bench}"
NETTY_BENCH_JAR="${NETTY_BENCH_JAR:-${REPO_ROOT}/external/bin/netty-bench.jar}"
JAVA_BIN="${JAVA_BIN:-/opt/software/java21/bin/java}"

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

if [[ "${PROTOCOL}" != "http2" && "${PROTOCOL}" != "https2" ]]; then
  printf 'unsupported HTTP/2 protocol: %s\n' "${PROTOCOL}" >&2
  exit 1
fi
if [[ "${PROTOCOL}" == "https2" && "${TLS_VERSION}" != "1.2" && "${TLS_VERSION}" != "1.3" ]]; then
  printf 'HTTP/2 over TLS requires TLS 1.2 or 1.3: %s\n' "${TLS_VERSION}" >&2
  exit 1
fi
require_positive_integer PAYLOAD "${PAYLOAD}"
require_positive_integer SERVER_GOMAXPROCS "${SERVER_GOMAXPROCS}"
require_positive_integer EVENT_LOOPS "${EVENT_LOOPS}"
require_positive_integer GNALLOY_WORKERS "${GNALLOY_WORKERS}"
require_positive_integer GNALLOY_READ_BUFFER_SIZE "${GNALLOY_READ_BUFFER_SIZE}"
require_positive_integer GNALLOY_MAX_MESSAGES_PER_READ "${GNALLOY_MAX_MESSAGES_PER_READ}"
taskset -c "${SERVER_CPU_SET}" true
if [[ -e "${SERVER_PID_FILE}" ]]; then
  printf 'server PID file already exists: %s\n' "${SERVER_PID_FILE}" >&2
  exit 1
fi
printf '%s\n' "$$" >"${SERVER_PID_FILE}"

tls_args=()
netty_tls_args=()
if [[ "${PROTOCOL}" == "https2" ]]; then
  tls_args+=( -tls-version "${TLS_VERSION}" -alpn h2 )
  netty_tls_args+=( --tls-version "${TLS_VERSION}" --alpn h2 )
fi

case "${FRAMEWORK}" in
  gnalloy)
    require_executable "${GNALLOY_BENCH}"
    profile_args=()
    if [[ -n "${GNALLOY_CPU_PROFILE}" ]]; then
      mkdir -p "$(dirname "${GNALLOY_CPU_PROFILE}")"
      profile_args+=( -cpuprofile "${GNALLOY_CPU_PROFILE}" )
    fi
    exec taskset -c "${SERVER_CPU_SET}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${GNALLOY_BENCH}" \
      -protocol "${PROTOCOL}" -server-only=true -addr "${SERVER_ADDR}" -payload "${PAYLOAD}" \
      -backend epoll -boss 1 -workers "${GNALLOY_WORKERS}" -boss-cpus "${GNALLOY_BOSS_CPU_SET}" \
      -worker-cpus "${GNALLOY_WORKER_CPU_SET}" -read-buffer-size "${GNALLOY_READ_BUFFER_SIZE}" \
      -max-messages-per-read "${GNALLOY_MAX_MESSAGES_PER_READ}" \
      -timeout 5m "${tls_args[@]}" "${profile_args[@]}"
    ;;
  hertz)
    require_executable "${HERTZ_BENCH}"
    exec taskset -c "${SERVER_CPU_SET}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${HERTZ_BENCH}" \
      -protocol "${PROTOCOL}" -addr "${SERVER_ADDR}" -payload "${PAYLOAD}" -timeout 5m "${tls_args[@]}"
    ;;
  netty)
    if [[ ! -f "${NETTY_BENCH_JAR}" ]]; then
      printf 'benchmark file is missing: %s\n' "${NETTY_BENCH_JAR}" >&2
      exit 1
    fi
    require_executable "${JAVA_BIN}"
    exec taskset -c "${SERVER_CPU_SET}" "${JAVA_BIN}" -XX:ActiveProcessorCount="${SERVER_GOMAXPROCS}" \
      -jar "${NETTY_BENCH_JAR}" --protocol "${PROTOCOL}" --server-only true --addr "${SERVER_ADDR}" \
      --payload "${PAYLOAD}" --backend epoll --event-loops "${EVENT_LOOPS}" --timeout 5m "${netty_tls_args[@]}"
    ;;
  *)
    printf 'unsupported HTTP/2 framework: %s\n' "${FRAMEWORK}" >&2
    exit 1
    ;;
esac
