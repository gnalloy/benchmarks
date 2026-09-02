#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

GNALLOY_BENCH="${GNALLOY_BENCH:-${REPO_ROOT}/external/bin/gnalloy-bench}"
GNET_BENCH="${GNET_BENCH:-${REPO_ROOT}/external/bin/gnet-bench}"
NETTY_BENCH_JAR="${NETTY_BENCH_JAR:-${REPO_ROOT}/external/bin/netty-bench.jar}"
JAVA_BIN="${JAVA_BIN:-java}"
OUTPUT="${OUTPUT:-${REPO_ROOT}/reports/raw/linux-udp-optimized.out}"
PAYLOADS="${PAYLOADS:-128 1024}"
CONNECTIONS="${CONNECTIONS:-64}"
MESSAGES="${MESSAGES:-20000}"
WARMUP_MESSAGES="${WARMUP_MESSAGES:-1000}"
LATENCY_SAMPLE_RATE="${LATENCY_SAMPLE_RATE:-64}"
REPETITIONS="${REPETITIONS:-5}"
EVENT_LOOPS="${EVENT_LOOPS:-4}"
GOMAXPROCS_VALUE="${GOMAXPROCS_VALUE:-8}"
COOLDOWN_SECONDS="${COOLDOWN_SECONDS:-10}"
NICE_LEVEL="${NICE_LEVEL:-0}"
SET_PERFORMANCE_GOVERNOR="${SET_PERFORMANCE_GOVERNOR:-1}"

declare -a GOVERNOR_FILES=()
declare -a GOVERNOR_VALUES=()

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

require_nice_level() {
  if [[ ! "${NICE_LEVEL}" =~ ^-?[0-9]+$ ]] || ((NICE_LEVEL < -20 || NICE_LEVEL > 19)); then
    printf 'NICE_LEVEL must be between -20 and 19: %s\n' "${NICE_LEVEL}" >&2
    exit 1
  fi
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
  trap restore_governors EXIT INT TERM
  for path in "${GOVERNOR_FILES[@]}"; do
    printf 'performance' >"${path}"
  done
}

run_command() {
  nice -n "${NICE_LEVEL}" "$@"
}

run_gnalloy() {
  local payload="$1"
  local run="$2"
  printf 'case=gnalloy payload=%d run=%d\n' "${payload}" "${run}" | tee -a "${OUTPUT}"
  run_command env GOMAXPROCS="${GOMAXPROCS_VALUE}" "${GNALLOY_BENCH}" \
    -protocol udp-echo \
    -backend epoll \
    -payload "${payload}" \
    -connections "${CONNECTIONS}" \
    -messages "${MESSAGES}" \
    -warmup-messages "${WARMUP_MESSAGES}" \
    -latency-sample-rate "${LATENCY_SAMPLE_RATE}" \
    -boss 1 \
    -workers "${EVENT_LOOPS}" \
    -reuseport=true \
    -max-messages-per-read 1 \
    -timeout 5m | tee -a "${OUTPUT}"
}

run_gnet() {
  local payload="$1"
  local run="$2"
  printf 'case=gnet payload=%d run=%d\n' "${payload}" "${run}" | tee -a "${OUTPUT}"
  run_command env GOMAXPROCS="${GOMAXPROCS_VALUE}" "${GNET_BENCH}" \
    -protocol udp-echo \
    -payload "${payload}" \
    -connections "${CONNECTIONS}" \
    -messages "${MESSAGES}" \
    -warmup-messages "${WARMUP_MESSAGES}" \
    -latency-sample-rate "${LATENCY_SAMPLE_RATE}" \
    -event-loops "${EVENT_LOOPS}" \
    -multicore=true \
    -timeout 5m | tee -a "${OUTPUT}"
}

run_netty() {
  local payload="$1"
  local run="$2"
  printf 'case=netty payload=%d run=%d\n' "${payload}" "${run}" | tee -a "${OUTPUT}"
  run_command "${JAVA_BIN}" -jar "${NETTY_BENCH_JAR}" \
    --protocol udp-echo \
    --backend epoll \
    --payload "${payload}" \
    --connections "${CONNECTIONS}" \
    --messages "${MESSAGES}" \
    --warmup-messages "${WARMUP_MESSAGES}" \
    --latency-sample-rate "${LATENCY_SAMPLE_RATE}" \
    --event-loops "${EVENT_LOOPS}" \
    --timeout 5m | tee -a "${OUTPUT}"
}

cooldown() {
  sleep "${COOLDOWN_SECONDS}"
}

run_rotation() {
  local payload="$1"
  local run="$2"
  case $(((run - 1) % 3)) in
    0)
      run_gnalloy "${payload}" "${run}"
      cooldown
      run_gnet "${payload}" "${run}"
      cooldown
      run_netty "${payload}" "${run}"
      ;;
    1)
      run_gnet "${payload}" "${run}"
      cooldown
      run_netty "${payload}" "${run}"
      cooldown
      run_gnalloy "${payload}" "${run}"
      ;;
    2)
      run_netty "${payload}" "${run}"
      cooldown
      run_gnalloy "${payload}" "${run}"
      cooldown
      run_gnet "${payload}" "${run}"
      ;;
  esac
  cooldown
}

require_executable "${GNALLOY_BENCH}"
require_executable "${GNET_BENCH}"
require_file "${NETTY_BENCH_JAR}"
if ! command -v "${JAVA_BIN}" >/dev/null 2>&1; then
  printf 'Java executable is missing: %s\n' "${JAVA_BIN}" >&2
  exit 1
fi
if [[ "$(uname -s)" != "Linux" ]]; then
  printf 'this benchmark requires Linux epoll\n' >&2
  exit 1
fi
require_positive_integer CONNECTIONS "${CONNECTIONS}"
require_positive_integer MESSAGES "${MESSAGES}"
require_positive_integer WARMUP_MESSAGES "${WARMUP_MESSAGES}"
require_positive_integer LATENCY_SAMPLE_RATE "${LATENCY_SAMPLE_RATE}"
require_positive_integer REPETITIONS "${REPETITIONS}"
require_positive_integer EVENT_LOOPS "${EVENT_LOOPS}"
require_positive_integer GOMAXPROCS_VALUE "${GOMAXPROCS_VALUE}"
require_non_negative_integer COOLDOWN_SECONDS "${COOLDOWN_SECONDS}"
require_nice_level
configure_governors

mkdir -p "$(dirname "${OUTPUT}")"
: >"${OUTPUT}"
{
  printf 'timestamp=%s\n' "$(date --iso-8601=seconds)"
  printf 'hostname=%s\n' "$(hostname)"
  printf 'kernel=%s\n' "$(uname -srmo)"
  printf 'cpuModel=%s\n' "$(awk -F ': ' '/^model name/{print $2; exit}' /proc/cpuinfo)"
  printf 'cpus=%s\n' "$(nproc)"
  printf 'governor=%s\n' "$(cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor 2>/dev/null || printf unknown)"
  printf 'connections=%s messages=%s warmupMessages=%s latencySampleRate=%s repetitions=%s eventLoops=%s gomaxprocs=%s cooldownSeconds=%s niceLevel=%s\n' \
    "${CONNECTIONS}" "${MESSAGES}" "${WARMUP_MESSAGES}" "${LATENCY_SAMPLE_RATE}" "${REPETITIONS}" "${EVENT_LOOPS}" "${GOMAXPROCS_VALUE}" "${COOLDOWN_SECONDS}" "${NICE_LEVEL}"
  sha256sum "${GNALLOY_BENCH}" "${GNET_BENCH}" "${NETTY_BENCH_JAR}"
  go version "${GNALLOY_BENCH}" 2>/dev/null || true
  go version "${GNET_BENCH}" 2>/dev/null || true
  "${JAVA_BIN}" -version 2>&1 | head -3
} | tee -a "${OUTPUT}"

read -r -a payload_values <<<"${PAYLOADS}"
for payload in "${payload_values[@]}"; do
  require_positive_integer PAYLOAD "${payload}"
  for ((run = 1; run <= REPETITIONS; run++)); do
    run_rotation "${payload}" "${run}"
  done
done
