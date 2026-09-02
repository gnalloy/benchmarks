#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

FRAMEWORK="${1:-}"
SERVER_ADDR="${SERVER_ADDR:-0.0.0.0:19091}"
PROTOCOL="${PROTOCOL:-http1}"
TLS_VERSION="${TLS_VERSION:-1.3}"
ALPN="${ALPN:-http/1.1}"
CIPHER_SUITES="${CIPHER_SUITES:-}"
PAYLOAD="${PAYLOAD:-128}"
SERVER_CPU_SET="${SERVER_CPU_SET:-0,1,2,3}"
SERVER_GOMAXPROCS="${SERVER_GOMAXPROCS:-4}"
EVENT_LOOPS="${EVENT_LOOPS:-4}"
GNALLOY_WORKERS="${GNALLOY_WORKERS:-4}"
GNALLOY_READ_BUFFER_SIZE="${GNALLOY_READ_BUFFER_SIZE:-4096}"
GNALLOY_MAX_MESSAGES_PER_READ="${GNALLOY_MAX_MESSAGES_PER_READ:-32}"
GNALLOY_EVENT_BATCH_SIZE="${GNALLOY_EVENT_BATCH_SIZE:-0}"
GNALLOY_FLUSH_STRATEGY="${GNALLOY_FLUSH_STRATEGY:-read-complete}"
GNALLOY_BOSS_CPU_SET="${GNALLOY_BOSS_CPU_SET:-3}"
GNALLOY_WORKER_CPU_SET="${GNALLOY_WORKER_CPU_SET:-0,1,2,3}"
GNALLOY_CPU_PROFILE="${GNALLOY_CPU_PROFILE:-}"
GNALLOY_ALLOC_PROFILE="${GNALLOY_ALLOC_PROFILE:-}"
GNALLOY_RUNTIME_TRACE="${GNALLOY_RUNTIME_TRACE:-}"
FASTHTTP_CPU_PROFILE="${FASTHTTP_CPU_PROFILE:-}"
NICE_LEVEL="${NICE_LEVEL:-0}"
SERVER_PID_FILE="${SERVER_PID_FILE:-/tmp/gnalloy-http1-cross-host.pid}"
GNALLOY_BENCH="${GNALLOY_BENCH:-${REPO_ROOT}/external/bin/gnalloy-bench}"
FASTHTTP_BENCH="${FASTHTTP_BENCH:-${REPO_ROOT}/external/bin/fasthttp-bench}"
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

require_nonnegative_integer() {
  local name="$1"
  local value="$2"
  if [[ ! "${value}" =~ ^[0-9]+$ ]]; then
    printf '%s must be a non-negative integer: %s\n' "${name}" "${value}" >&2
    exit 1
  fi
}

require_file() {
  local path="$1"
  if [[ ! -f "${path}" ]]; then
    printf 'benchmark file is missing: %s\n' "${path}" >&2
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

case "${FRAMEWORK}" in
  gnalloy)
    require_executable "${GNALLOY_BENCH}"
    ;;
  fasthttp)
    require_executable "${FASTHTTP_BENCH}"
    ;;
  netty)
    require_file "${NETTY_BENCH_JAR}"
    require_executable "${JAVA_BIN}"
    ;;
  *)
    printf 'unsupported HTTP/1 framework: %s\n' "${FRAMEWORK}" >&2
    exit 1
    ;;
esac

require_positive_integer PAYLOAD "${PAYLOAD}"
require_positive_integer SERVER_GOMAXPROCS "${SERVER_GOMAXPROCS}"
require_positive_integer EVENT_LOOPS "${EVENT_LOOPS}"
require_positive_integer GNALLOY_WORKERS "${GNALLOY_WORKERS}"
require_positive_integer GNALLOY_READ_BUFFER_SIZE "${GNALLOY_READ_BUFFER_SIZE}"
require_positive_integer GNALLOY_MAX_MESSAGES_PER_READ "${GNALLOY_MAX_MESSAGES_PER_READ}"
require_nonnegative_integer GNALLOY_EVENT_BATCH_SIZE "${GNALLOY_EVENT_BATCH_SIZE}"
if [[ "${PROTOCOL}" != "http1" && "${PROTOCOL}" != "https1" ]]; then
  printf 'unsupported HTTP/1 family protocol: %s\n' "${PROTOCOL}" >&2
  exit 1
fi
case "${GNALLOY_FLUSH_STRATEGY}" in
  immediate|read-complete|event-loop-batch) ;;
  *)
    printf 'unsupported Gnalloy flush strategy: %s\n' "${GNALLOY_FLUSH_STRATEGY}" >&2
    exit 1
    ;;
esac
tls_args=()
netty_tls_args=()
if [[ "${PROTOCOL}" == "https1" ]]; then
  tls_args+=( -tls-version "${TLS_VERSION}" -alpn "${ALPN}" )
  netty_tls_args+=( --tls-version "${TLS_VERSION}" --alpn "${ALPN}" )
  if [[ -n "${CIPHER_SUITES}" ]]; then
    tls_args+=( -cipher-suites "${CIPHER_SUITES}" )
    netty_tls_args+=( --cipher-suites "${CIPHER_SUITES}" )
  fi
fi
if [[ ! "${NICE_LEVEL}" =~ ^-?[0-9]+$ ]] || ((NICE_LEVEL < -20 || NICE_LEVEL > 19)); then
  printf 'NICE_LEVEL must be between -20 and 19: %s\n' "${NICE_LEVEL}" >&2
  exit 1
fi
if [[ -e "${SERVER_PID_FILE}" ]]; then
  printf 'server PID file already exists: %s\n' "${SERVER_PID_FILE}" >&2
  exit 1
fi
taskset -c "${SERVER_CPU_SET}" true
printf '%s\n' "$$" >"${SERVER_PID_FILE}"

case "${FRAMEWORK}" in
  gnalloy)
    profile_args=()
    if [[ -n "${GNALLOY_CPU_PROFILE}" ]]; then
      mkdir -p "$(dirname "${GNALLOY_CPU_PROFILE}")"
      profile_args+=( -cpuprofile "${GNALLOY_CPU_PROFILE}" )
    fi
    if [[ -n "${GNALLOY_ALLOC_PROFILE}" ]]; then
      mkdir -p "$(dirname "${GNALLOY_ALLOC_PROFILE}")"
      profile_args+=( -allocprofile "${GNALLOY_ALLOC_PROFILE}" )
    fi
    if [[ -n "${GNALLOY_RUNTIME_TRACE}" ]]; then
      mkdir -p "$(dirname "${GNALLOY_RUNTIME_TRACE}")"
      profile_args+=( -trace "${GNALLOY_RUNTIME_TRACE}" )
    fi
    exec taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${GNALLOY_BENCH}" \
      -protocol "${PROTOCOL}" -server-only=true -addr "${SERVER_ADDR}" -payload "${PAYLOAD}" -backend epoll -boss 1 \
      -workers "${GNALLOY_WORKERS}" -boss-cpus "${GNALLOY_BOSS_CPU_SET}" \
      -worker-cpus "${GNALLOY_WORKER_CPU_SET}" -reuseport=true -read-buffer-size "${GNALLOY_READ_BUFFER_SIZE}" \
      -max-messages-per-read "${GNALLOY_MAX_MESSAGES_PER_READ}" -event-batch-size "${GNALLOY_EVENT_BATCH_SIZE}" \
      -flush-strategy "${GNALLOY_FLUSH_STRATEGY}" \
      -timeout 5m "${tls_args[@]}" "${profile_args[@]}"
    ;;
  fasthttp)
    profile_args=()
    if [[ -n "${FASTHTTP_CPU_PROFILE}" ]]; then
      mkdir -p "$(dirname "${FASTHTTP_CPU_PROFILE}")"
      profile_args+=( -cpuprofile "${FASTHTTP_CPU_PROFILE}" )
    fi
    exec taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${FASTHTTP_BENCH}" \
      -protocol "${PROTOCOL}" -server-only=true -addr "${SERVER_ADDR}" -payload "${PAYLOAD}" -timeout 5m "${tls_args[@]}" "${profile_args[@]}"
    ;;
  netty)
    exec taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" "${JAVA_BIN}" -jar "${NETTY_BENCH_JAR}" \
      --protocol "${PROTOCOL}" --server-only true --addr "${SERVER_ADDR}" --payload "${PAYLOAD}" --backend epoll \
      --event-loops "${EVENT_LOOPS}" --timeout 5m "${netty_tls_args[@]}"
    ;;
esac
