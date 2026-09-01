#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

GNALLOY_BENCH="${GNALLOY_BENCH:-${REPO_ROOT}/external/bin/gnalloy-bench}"
NETPOLL_BENCH="${NETPOLL_BENCH:-${REPO_ROOT}/external/bin/netpoll-bench}"
OUTPUT="${OUTPUT:-${REPO_ROOT}/reports/raw/linux-tcp-optimized.out}"
PAYLOADS="${PAYLOADS:-64 1024 16384}"
CONNECTIONS="${CONNECTIONS:-64}"
MESSAGES="${MESSAGES:-20000}"
WARMUP_MESSAGES="${WARMUP_MESSAGES:-1000}"
LATENCY_SAMPLE_RATE="${LATENCY_SAMPLE_RATE:-128}"
REPETITIONS="${REPETITIONS:-5}"
EXECUTOR_WORKERS="${EXECUTOR_WORKERS:-8}"
GOMAXPROCS_VALUE="${GOMAXPROCS_VALUE:-8}"
COOLDOWN_SECONDS="${COOLDOWN_SECONDS:-10}"
NICE_LEVEL="${NICE_LEVEL:-0}"

require_executable() {
  local path="$1"
  if [[ ! -x "${path}" ]]; then
    printf 'benchmark executable is missing: %s\n' "${path}" >&2
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

read_buffer_size() {
  local payload="$1"
  if ((payload < 4096)); then
    printf '4096\n'
    return
  fi
  printf '%d\n' "$((payload * 2))"
}

run_benchmark() {
  nice -n "${NICE_LEVEL}" env GOMAXPROCS="${GOMAXPROCS_VALUE}" "$@"
}

run_gnalloy() {
  local payload="$1"
  local run="$2"
  local read_buffer
  read_buffer="$(read_buffer_size "${payload}")"

  printf 'case=gnalloy payload=%d readBuffer=%d run=%d\n' "${payload}" "${read_buffer}" "${run}" | tee -a "${OUTPUT}"
  run_benchmark "${GNALLOY_BENCH}" \
    -protocol tcp-echo \
    -backend epoll \
    -payload "${payload}" \
    -connections "${CONNECTIONS}" \
    -messages "${MESSAGES}" \
    -boss 1 \
    -workers 1 \
    -read-buffer-size "${read_buffer}" \
    -tcp-echo-mode owner-executor \
    -flush-strategy event-loop-batch \
    -tcp-echo-executor-workers "${EXECUTOR_WORKERS}" \
    -tcp-echo-executor-queue-size 4096 \
    -latency-sample-rate "${LATENCY_SAMPLE_RATE}" \
    -warmup-messages "${WARMUP_MESSAGES}" \
    -timeout 5m | tee -a "${OUTPUT}"
}

run_netpoll() {
  local payload="$1"
  local run="$2"

  printf 'case=netpoll payload=%d run=%d\n' "${payload}" "${run}" | tee -a "${OUTPUT}"
  run_benchmark "${NETPOLL_BENCH}" \
    -protocol tcp-echo \
    -payload "${payload}" \
    -connections "${CONNECTIONS}" \
    -messages "${MESSAGES}" \
    -latency-sample-rate "${LATENCY_SAMPLE_RATE}" \
    -warmup-messages "${WARMUP_MESSAGES}" \
    -timeout 5m | tee -a "${OUTPUT}"
}

require_executable "${GNALLOY_BENCH}"
require_executable "${NETPOLL_BENCH}"
if [[ "$(uname -s)" != "Linux" ]]; then
  printf 'this benchmark requires Linux epoll\n' >&2
  exit 1
fi
require_positive_integer CONNECTIONS "${CONNECTIONS}"
require_positive_integer MESSAGES "${MESSAGES}"
require_positive_integer WARMUP_MESSAGES "${WARMUP_MESSAGES}"
require_positive_integer LATENCY_SAMPLE_RATE "${LATENCY_SAMPLE_RATE}"
require_positive_integer REPETITIONS "${REPETITIONS}"
require_positive_integer EXECUTOR_WORKERS "${EXECUTOR_WORKERS}"
require_positive_integer GOMAXPROCS_VALUE "${GOMAXPROCS_VALUE}"
require_non_negative_integer COOLDOWN_SECONDS "${COOLDOWN_SECONDS}"
require_nice_level

mkdir -p "$(dirname "${OUTPUT}")"
: >"${OUTPUT}"
{
  printf 'timestamp=%s\n' "$(date --iso-8601=seconds)"
  printf 'hostname=%s\n' "$(hostname)"
  printf 'kernel=%s\n' "$(uname -srmo)"
  printf 'cpuModel=%s\n' "$(awk -F ': ' '/^model name/{print $2; exit}' /proc/cpuinfo)"
  printf 'governor=%s\n' "$(cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor 2>/dev/null || printf unknown)"
  printf 'connections=%s messages=%s warmupMessages=%s latencySampleRate=%s repetitions=%s executorWorkers=%s gomaxprocs=%s cooldownSeconds=%s niceLevel=%s\n' \
    "${CONNECTIONS}" "${MESSAGES}" "${WARMUP_MESSAGES}" "${LATENCY_SAMPLE_RATE}" "${REPETITIONS}" "${EXECUTOR_WORKERS}" "${GOMAXPROCS_VALUE}" "${COOLDOWN_SECONDS}" "${NICE_LEVEL}"
  sha256sum "${GNALLOY_BENCH}" "${NETPOLL_BENCH}"
  if command -v go >/dev/null 2>&1; then
    go version "${GNALLOY_BENCH}"
    go version "${NETPOLL_BENCH}"
  fi
} | tee -a "${OUTPUT}"

read -r -a payload_values <<<"${PAYLOADS}"
for payload in "${payload_values[@]}"; do
  require_positive_integer PAYLOAD "${payload}"
  for ((run = 1; run <= REPETITIONS; run++)); do
    if ((run % 2 == 1)); then
      run_gnalloy "${payload}" "${run}"
      sleep "${COOLDOWN_SECONDS}"
      run_netpoll "${payload}" "${run}"
    else
      run_netpoll "${payload}" "${run}"
      sleep "${COOLDOWN_SECONDS}"
      run_gnalloy "${payload}" "${run}"
    fi
    sleep "${COOLDOWN_SECONDS}"
  done
done
