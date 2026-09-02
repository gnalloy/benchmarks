#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

GNALLOY_BENCH="${GNALLOY_BENCH:-${REPO_ROOT}/external/bin/gnalloy-bench}"
GNET_BENCH="${GNET_BENCH:-${REPO_ROOT}/external/bin/gnet-bench}"
NETTY_BENCH_JAR="${NETTY_BENCH_JAR:-${REPO_ROOT}/external/bin/netty-bench.jar}"
UDP_LOAD="${UDP_LOAD:-${REPO_ROOT}/external/bin/udp-load}"
JAVA_BIN="${JAVA_BIN:-java}"
OUTPUT="${OUTPUT:-${REPO_ROOT}/reports/raw/linux-udp-common-client.out}"
SERVER_ADDR="${SERVER_ADDR:-127.0.0.1:19090}"
PAYLOADS="${PAYLOADS:-128 1024}"
CONNECTIONS="${CONNECTIONS:-64}"
MESSAGES="${MESSAGES:-20000}"
WARMUP_MESSAGES="${WARMUP_MESSAGES:-1000}"
LATENCY_SAMPLE_RATE="${LATENCY_SAMPLE_RATE:-64}"
TARGET_RATE="${TARGET_RATE:-0}"
REPETITIONS="${REPETITIONS:-5}"
EVENT_LOOPS="${EVENT_LOOPS:-4}"
GNALLOY_WORKERS="${GNALLOY_WORKERS:-3}"
GNALLOY_MAX_MESSAGES_PER_READ="${GNALLOY_MAX_MESSAGES_PER_READ:-1}"
GNALLOY_BOSS_CPU_SET="${GNALLOY_BOSS_CPU_SET:-3}"
GNALLOY_WORKER_CPU_SET="${GNALLOY_WORKER_CPU_SET:-0,1,2}"
SERVER_GOMAXPROCS="${SERVER_GOMAXPROCS:-4}"
CLIENT_GOMAXPROCS="${CLIENT_GOMAXPROCS:-4}"
SERVER_CPU_SET="${SERVER_CPU_SET:-0-3}"
CLIENT_CPU_SET="${CLIENT_CPU_SET:-4-7}"
COOLDOWN_SECONDS="${COOLDOWN_SECONDS:-10}"
NICE_LEVEL="${NICE_LEVEL:-0}"
SET_PERFORMANCE_GOVERNOR="${SET_PERFORMANCE_GOVERNOR:-1}"
READY_TIMEOUT_SECONDS="${READY_TIMEOUT_SECONDS:-15}"

declare -a GOVERNOR_FILES=()
declare -a GOVERNOR_VALUES=()
SERVER_PID=""
SERVER_LOG=""

require_executable() {
  local path="$1"
  if [[ ! -x "${path}" ]]; then
    printf 'benchmark executable is missing: %s\n' "${path}" >&2
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

require_positive_integer() {
  local name="$1"
  local value="$2"
  if [[ ! "${value}" =~ ^[1-9][0-9]*$ ]]; then
    printf '%s must be a positive integer: %s\n' "${name}" "${value}" >&2
    exit 1
  fi
}

require_non_negative_integer() {
  local name="$1"
  local value="$2"
  if [[ ! "${value}" =~ ^[0-9]+$ ]]; then
    printf '%s must be a non-negative integer: %s\n' "${name}" "${value}" >&2
    exit 1
  fi
}

expand_cpu_set() {
  local value="$1"
  local segment first last cpu
  local -a segments=()
  IFS=',' read -r -a segments <<<"${value}"
  for segment in "${segments[@]}"; do
    if [[ "${segment}" =~ ^([0-9]+)-([0-9]+)$ ]]; then
      first="${BASH_REMATCH[1]}"
      last="${BASH_REMATCH[2]}"
      if ((first > last)); then
        printf 'invalid CPU range: %s\n' "${segment}" >&2
        exit 1
      fi
      for ((cpu = first; cpu <= last; cpu++)); do
        printf '%d\n' "${cpu}"
      done
    elif [[ "${segment}" =~ ^[0-9]+$ ]]; then
      printf '%d\n' "${segment}"
    else
      printf 'invalid CPU set: %s\n' "${value}" >&2
      exit 1
    fi
  done
}

validate_cpu_sets() {
  local cpu
  declare -A server_cpus=()
  while read -r cpu; do
    server_cpus["${cpu}"]=1
  done < <(expand_cpu_set "${SERVER_CPU_SET}")
  while read -r cpu; do
    if [[ -n "${server_cpus[${cpu}]+present}" ]]; then
      printf 'server and client CPU sets overlap at CPU %s\n' "${cpu}" >&2
      exit 1
    fi
  done < <(expand_cpu_set "${CLIENT_CPU_SET}")
  local tuned_set
  for tuned_set in "${GNALLOY_BOSS_CPU_SET}" "${GNALLOY_WORKER_CPU_SET}"; do
    while read -r cpu; do
      if [[ -z "${server_cpus[${cpu}]+present}" ]]; then
        printf 'Gnalloy CPU %s is outside the server CPU set\n' "${cpu}" >&2
        exit 1
      fi
    done < <(expand_cpu_set "${tuned_set}")
  done
  taskset -c "${SERVER_CPU_SET}" true
  taskset -c "${CLIENT_CPU_SET}" true
}

restore_governors() {
  local index
  for index in "${!GOVERNOR_FILES[@]}"; do
    printf '%s' "${GOVERNOR_VALUES[index]}" >"${GOVERNOR_FILES[index]}" || true
  done
}

configure_governors() {
  if [[ "${SET_PERFORMANCE_GOVERNOR}" != "1" ]]; then
    return
  fi
  local path
  for path in /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor; do
    if [[ ! -w "${path}" ]]; then
      printf 'CPU governor is not writable: %s\n' "${path}" >&2
      exit 1
    fi
    GOVERNOR_FILES+=("${path}")
    GOVERNOR_VALUES+=("$(<"${path}")")
  done
  for path in "${GOVERNOR_FILES[@]}"; do
    printf 'performance' >"${path}"
  done
}

stop_server() {
  if [[ -z "${SERVER_PID}" ]]; then
    return
  fi
  if kill -0 "${SERVER_PID}" 2>/dev/null; then
    kill -TERM "${SERVER_PID}" 2>/dev/null || true
    local attempt
    for ((attempt = 0; attempt < 100; attempt++)); do
      if ! kill -0 "${SERVER_PID}" 2>/dev/null; then
        break
      fi
      sleep 0.1
    done
  fi
  if kill -0 "${SERVER_PID}" 2>/dev/null; then
    kill -KILL "${SERVER_PID}" 2>/dev/null || true
  fi
  wait "${SERVER_PID}" 2>/dev/null || true
  SERVER_PID=""
}

cleanup() {
  stop_server
  restore_governors
  if [[ -n "${SERVER_LOG}" && -f "${SERVER_LOG}" ]]; then
    rm -f -- "${SERVER_LOG}"
  fi
}

wait_for_server() {
  local framework="$1"
  local attempts=$((READY_TIMEOUT_SECONDS * 10))
  local attempt
  for ((attempt = 0; attempt < attempts; attempt++)); do
    if grep -Fq "serverReady=true framework=${framework} protocol=udp-echo addr=${SERVER_ADDR}" "${SERVER_LOG}"; then
      return
    fi
    if ! kill -0 "${SERVER_PID}" 2>/dev/null; then
      printf 'server exited before readiness: %s\n' "${framework}" >&2
      cat "${SERVER_LOG}" >&2
      exit 1
    fi
    sleep 0.1
  done
  printf 'server readiness timed out: %s\n' "${framework}" >&2
  cat "${SERVER_LOG}" >&2
  exit 1
}

start_server() {
  local framework="$1"
  : >"${SERVER_LOG}"
  case "${framework}" in
    gnalloy)
      taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${GNALLOY_BENCH}" \
        -protocol udp-echo -server-only=true -addr "${SERVER_ADDR}" -backend epoll -boss 1 \
		-workers "${GNALLOY_WORKERS}" -boss-cpus "${GNALLOY_BOSS_CPU_SET}" -worker-cpus "${GNALLOY_WORKER_CPU_SET}" \
		-reuseport=true -max-messages-per-read "${GNALLOY_MAX_MESSAGES_PER_READ}" -timeout 5m >"${SERVER_LOG}" 2>&1 &
      ;;
    gnet)
      taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" env GOMAXPROCS="${SERVER_GOMAXPROCS}" "${GNET_BENCH}" \
        -protocol udp-echo -server-only=true -addr "${SERVER_ADDR}" -event-loops "${EVENT_LOOPS}" \
        -multicore=true -timeout 5m >"${SERVER_LOG}" 2>&1 &
      ;;
    netty)
      taskset -c "${SERVER_CPU_SET}" nice -n "${NICE_LEVEL}" "${JAVA_BIN}" -jar "${NETTY_BENCH_JAR}" \
        --protocol udp-echo --server-only true --addr "${SERVER_ADDR}" --backend epoll \
        --event-loops "${EVENT_LOOPS}" --timeout 5m >"${SERVER_LOG}" 2>&1 &
      ;;
    *)
      printf 'unsupported framework: %s\n' "${framework}" >&2
      exit 1
      ;;
  esac
  SERVER_PID="$!"
  wait_for_server "${framework}"
}

run_case() {
  local framework="$1"
  local payload="$2"
  local run="$3"
  printf 'case=%s payload=%d run=%d\n' "${framework}" "${payload}" "${run}" | tee -a "${OUTPUT}"
  start_server "${framework}"
  tee -a "${OUTPUT}" <"${SERVER_LOG}"
  taskset -c "${CLIENT_CPU_SET}" nice -n "${NICE_LEVEL}" env GOMAXPROCS="${CLIENT_GOMAXPROCS}" "${UDP_LOAD}" \
    -server-framework "${framework}" -addr "${SERVER_ADDR}" -payload "${payload}" \
    -connections "${CONNECTIONS}" -messages "${MESSAGES}" -warmup-messages "${WARMUP_MESSAGES}" \
    -latency-sample-rate "${LATENCY_SAMPLE_RATE}" -target-rate "${TARGET_RATE}" -timeout 5m | tee -a "${OUTPUT}"
  stop_server
  sleep "${COOLDOWN_SECONDS}"
}

run_rotation() {
  local payload="$1"
  local run="$2"
  case $(((run - 1) % 3)) in
    0)
      run_case gnalloy "${payload}" "${run}"
      run_case gnet "${payload}" "${run}"
      run_case netty "${payload}" "${run}"
      ;;
    1)
      run_case gnet "${payload}" "${run}"
      run_case netty "${payload}" "${run}"
      run_case gnalloy "${payload}" "${run}"
      ;;
    2)
      run_case netty "${payload}" "${run}"
      run_case gnalloy "${payload}" "${run}"
      run_case gnet "${payload}" "${run}"
      ;;
  esac
}

require_executable "${GNALLOY_BENCH}"
require_executable "${GNET_BENCH}"
require_executable "${UDP_LOAD}"
require_file "${NETTY_BENCH_JAR}"
if [[ "$(uname -s)" != "Linux" ]]; then
  printf 'this benchmark requires Linux epoll\n' >&2
  exit 1
fi
for command in "${JAVA_BIN}" taskset sha256sum; do
  if ! command -v "${command}" >/dev/null 2>&1; then
    printf 'required command is missing: %s\n' "${command}" >&2
    exit 1
  fi
done
require_positive_integer CONNECTIONS "${CONNECTIONS}"
require_positive_integer MESSAGES "${MESSAGES}"
require_non_negative_integer WARMUP_MESSAGES "${WARMUP_MESSAGES}"
require_non_negative_integer LATENCY_SAMPLE_RATE "${LATENCY_SAMPLE_RATE}"
require_non_negative_integer TARGET_RATE "${TARGET_RATE}"
require_positive_integer REPETITIONS "${REPETITIONS}"
require_positive_integer EVENT_LOOPS "${EVENT_LOOPS}"
require_positive_integer GNALLOY_WORKERS "${GNALLOY_WORKERS}"
require_positive_integer GNALLOY_MAX_MESSAGES_PER_READ "${GNALLOY_MAX_MESSAGES_PER_READ}"
require_positive_integer SERVER_GOMAXPROCS "${SERVER_GOMAXPROCS}"
require_positive_integer CLIENT_GOMAXPROCS "${CLIENT_GOMAXPROCS}"
require_non_negative_integer COOLDOWN_SECONDS "${COOLDOWN_SECONDS}"
require_positive_integer READY_TIMEOUT_SECONDS "${READY_TIMEOUT_SECONDS}"
if [[ ! "${NICE_LEVEL}" =~ ^-?[0-9]+$ ]] || ((NICE_LEVEL < -20 || NICE_LEVEL > 19)); then
  printf 'NICE_LEVEL must be between -20 and 19: %s\n' "${NICE_LEVEL}" >&2
  exit 1
fi
validate_cpu_sets
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
configure_governors

mkdir -p "$(dirname "${OUTPUT}")"
SERVER_LOG="$(mktemp "${TMPDIR:-/tmp}/gnalloy-udp-server.XXXXXX")"
: >"${OUTPUT}"
{
  printf 'timestamp=%s\n' "$(date --iso-8601=seconds)"
  printf 'hostname=%s\n' "$(hostname)"
  printf 'kernel=%s\n' "$(uname -srmo)"
  printf 'cpuModel=%s\n' "$(awk -F ': ' '/^model name/{print $2; exit}' /proc/cpuinfo)"
  printf 'cpus=%s governor=%s\n' "$(nproc)" "$(cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor 2>/dev/null || printf unknown)"
  printf 'serverCPUSet=%s clientCPUSet=%s serverGOMAXPROCS=%s clientGOMAXPROCS=%s eventLoops=%s\n' \
    "${SERVER_CPU_SET}" "${CLIENT_CPU_SET}" "${SERVER_GOMAXPROCS}" "${CLIENT_GOMAXPROCS}" "${EVENT_LOOPS}"
  printf 'gnalloyWorkers=%s gnalloyBossCPUSet=%s gnalloyWorkerCPUSet=%s gnalloyMaxMessagesPerRead=%s\n' \
    "${GNALLOY_WORKERS}" "${GNALLOY_BOSS_CPU_SET}" "${GNALLOY_WORKER_CPU_SET}" "${GNALLOY_MAX_MESSAGES_PER_READ}"
  printf 'connections=%s messages=%s warmupMessages=%s targetRate=%s latencySampleRate=%s repetitions=%s cooldownSeconds=%s niceLevel=%s\n' \
    "${CONNECTIONS}" "${MESSAGES}" "${WARMUP_MESSAGES}" "${TARGET_RATE}" "${LATENCY_SAMPLE_RATE}" "${REPETITIONS}" "${COOLDOWN_SECONDS}" "${NICE_LEVEL}"
  printf 'excludedFrameworks=netpoll,fasthttp reason=no-comparable-udp-server\n'
  sha256sum "${GNALLOY_BENCH}" "${GNET_BENCH}" "${NETTY_BENCH_JAR}" "${UDP_LOAD}"
  go version "${GNALLOY_BENCH}" 2>/dev/null || true
  go version "${GNET_BENCH}" 2>/dev/null || true
  go version "${UDP_LOAD}" 2>/dev/null || true
  "${JAVA_BIN}" -version 2>&1 | head -3
} | tee -a "${OUTPUT}"

read -r -a payload_values <<<"${PAYLOADS}"
for payload in "${payload_values[@]}"; do
  require_positive_integer PAYLOAD "${payload}"
  for ((run = 1; run <= REPETITIONS; run++)); do
    run_rotation "${payload}" "${run}"
  done
done
