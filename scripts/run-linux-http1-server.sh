#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

FRAMEWORK="${1:-}"
SERVER_ADDR="${SERVER_ADDR:-0.0.0.0:19091}"
PAYLOAD="${PAYLOAD:-128}"
SERVER_CPU_SET="${SERVER_CPU_SET:-0,1,4,5}"
SERVER_GOMAXPROCS="${SERVER_GOMAXPROCS:-4}"
EVENT_LOOPS="${EVENT_LOOPS:-4}"
GNALLOY_WORKERS="${GNALLOY_WORKERS:-4}"
GNALLOY_HTTP1_MODE="${GNALLOY_HTTP1_MODE:-codec}"
GNALLOY_BOSS_CPU_SET="${GNALLOY_BOSS_CPU_SET:-4}"
GNALLOY_WORKER_CPU_SET="${GNALLOY_WORKER_CPU_SET:-0,1,4,5}"
GNALLOY_CPU_PROFILE="${GNALLOY_CPU_PROFILE:-}"
GNALLOY_RUNTIME_TRACE="${GNALLOY_RUNTIME_TRACE:-}"
NICE_LEVEL="${NICE_LEVEL:-0}"
SERVER_PID_FILE="${SERVER_PID_FILE:-/tmp/gnalloy-http1-cross-host.pid}"
GNALLOY_BENCH="${GNALLOY_BENCH:-${REPO_ROOT}/external/bin/gnalloy-bench}"
GNET_BENCH="${GNET_BENCH:-${REPO_ROOT}/external/bin/gnet-bench}"
NETPOLL_BENCH="${NETPOLL_BENCH:-${REPO_ROOT}/external/bin/netpoll-bench}"
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
  gnet)
    require_executable "${GNET_BENCH}"
    ;;
  netpoll)
    require_executable "${NETPOLL_BENCH}"
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
if [[ "${GNALLOY_HTTP1_MODE}" != "codec" && "${GNALLOY_HTTP1_MODE}" != "raw" ]]; then
  printf 'GNALLOY_HTTP1_MODE must be codec or raw: %s\n' "${GNALLOY_HTTP1_MODE}" >&2
  exit 1
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
    if [[ -n "${GNALLOY_RUNTIME_TRACE}" ]]; then
      mkdir -p "$(dirname "${GNALLOY_RUNTIME_TRACE}")"
      profile_args+=( -trace "${GNALLOY_RUNTIME_TRACE}" )
    fi
    exec taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${GNALLOY_BENCH}" \
      -protocol http1 -server-only=true -addr "${SERVER_ADDR}" -payload "${PAYLOAD}" -backend epoll -boss 1 \
      -workers "${GNALLOY_WORKERS}" -http1-mode "${GNALLOY_HTTP1_MODE}" -boss-cpus "${GNALLOY_BOSS_CPU_SET}" \
      -worker-cpus "${GNALLOY_WORKER_CPU_SET}" -reuseport=true -timeout 5m "${profile_args[@]}"
    ;;
  gnet)
    exec taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${GNET_BENCH}" \
      -protocol http1 -server-only=true -addr "${SERVER_ADDR}" -payload "${PAYLOAD}" \
      -event-loops "${EVENT_LOOPS}" -multicore=true -timeout 5m
    ;;
  netpoll)
    exec taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${NETPOLL_BENCH}" \
      -protocol http1 -server-only=true -addr "${SERVER_ADDR}" -payload "${PAYLOAD}" -timeout 5m
    ;;
  fasthttp)
    exec taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${FASTHTTP_BENCH}" \
      -protocol http1 -server-only=true -addr "${SERVER_ADDR}" -payload "${PAYLOAD}" -timeout 5m
    ;;
  netty)
    exec taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" "${JAVA_BIN}" -jar "${NETTY_BENCH_JAR}" \
      --protocol http1 --server-only true --addr "${SERVER_ADDR}" --payload "${PAYLOAD}" --backend epoll \
      --event-loops "${EVENT_LOOPS}" --timeout 5m
    ;;
esac
