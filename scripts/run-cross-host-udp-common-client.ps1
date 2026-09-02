[CmdletBinding()]
param(
    [string]$ClientHost = "172.16.8.172",
    [string]$ClientAddress = "172.16.8.172",
    [string]$ServerHost = "172.16.8.171",
    [ValidateRange(1, 65535)]
    [int]$SSHPort = 22,
    [string]$SSHUser = "root",
    [string]$ServerRepo = "/opt/test/gnalloy/benchmarks-e481c9a",
    [string]$ClientRepo = "/opt/test/gnalloy/benchmarks-cross-host",
    [string]$ServerScript = "./scripts/run-linux-udp-server.sh",
    [string]$ServerBindAddress = "0.0.0.0:19090",
    [string]$TargetAddress = "172.16.8.171:19090",
    [int[]]$Payloads = @(128, 1024),
    [ValidateRange(1, [int]::MaxValue)]
    [int]$Connections = 64,
    [ValidateRange(1, [int]::MaxValue)]
    [int]$Messages = 20000,
    [ValidateRange(0, [int]::MaxValue)]
    [int]$WarmupMessages = 1000,
    [ValidateRange(0, [int]::MaxValue)]
    [int]$LatencySampleRate = 64,
    [ValidateRange(0, [long]::MaxValue)]
    [long]$TargetRate = 0,
    [ValidateRange(1, 100)]
    [int]$Repetitions = 5,
    [ValidateRange(0, 300)]
    [int]$CooldownSeconds = 10,
    [bool]$SetPerformanceGovernor = $true,
    [ValidateRange(1, 64)]
    [int]$ClientGoMaxProcs = 4,
    [string]$ClientCPUSet = "0,1,2,4",
    [string]$ServerCPUSet = "0,1,4,5",
    [ValidateRange(1, 64)]
    [int]$ServerGoMaxProcs = 4,
    [ValidateRange(1, 64)]
    [int]$EventLoops = 4,
    [ValidateRange(1, 64)]
    [int]$GnalloyWorkers = 4,
    [ValidateRange(1, 1024)]
    [int]$GnalloyMaxMessagesPerRead = 64,
    [string]$GnalloyBossCPUSet = "4",
    [string]$GnalloyWorkerCPUSet = "0,1,4,5",
    [string]$ClientUdpLoad = "",
    [string]$Output = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($ClientUdpLoad)) {
    $ClientUdpLoad = "$ClientRepo/external/bin/udp-load"
}
if ([string]::IsNullOrWhiteSpace($Output)) {
    $mode = if ($TargetRate -gt 0) { "offered-$TargetRate" } else { "saturation" }
    $Output = Join-Path $repoRoot "reports/raw/linux-cross-host-udp-$mode.out"
}

$remotePidFile = "/tmp/gnalloy-udp-cross-host.pid"
$remoteLog = "/tmp/gnalloy-udp-cross-host.log"

function Assert-SafeRemoteValue {
    param([string]$Name, [string]$Value)
    if ([string]::IsNullOrWhiteSpace($Value) -or $Value -notmatch '^[A-Za-z0-9_./,:=-]+$') {
        throw "$Name contains unsupported characters: $Value"
    }
}

function Invoke-SSH {
    param([string]$HostName, [string]$Command, [switch]$IgnoreStandardError)
    $ssh = (Get-Command ssh -ErrorAction Stop).Source
    $startInfo = [System.Diagnostics.ProcessStartInfo]::new()
    $startInfo.FileName = $ssh
    $startInfo.UseShellExecute = $false
    $startInfo.RedirectStandardOutput = $true
    $startInfo.RedirectStandardError = $true
    foreach ($argument in @("-p", $SSHPort.ToString(), "-o", "BatchMode=yes", "-o", "ConnectTimeout=10", "${SSHUser}@${HostName}", $Command)) {
        $startInfo.ArgumentList.Add($argument)
    }
    $process = [System.Diagnostics.Process]::new()
    $process.StartInfo = $startInfo
    if (-not $process.Start()) {
        throw "Failed to start SSH"
    }
    $stdoutTask = $process.StandardOutput.ReadToEndAsync()
    $stderrTask = $process.StandardError.ReadToEndAsync()
    $process.WaitForExit()
    $stdout = $stdoutTask.GetAwaiter().GetResult().TrimEnd()
    $stderr = $stderrTask.GetAwaiter().GetResult().TrimEnd()
    if ($process.ExitCode -ne 0) {
        throw "SSH command failed with exit code $($process.ExitCode): $stderr"
    }
    if (-not $IgnoreStandardError -and -not [string]::IsNullOrWhiteSpace($stderr)) {
        [Console]::Error.WriteLine($stderr)
    }
    return $stdout
}

function Stop-RemoteServer {
    $template = @'
if test -r '__PID_FILE__'; then
  pid="$(cat '__PID_FILE__')"
  if kill -0 "$pid" 2>/dev/null; then
    kill -TERM "$pid" 2>/dev/null || true
    attempt=0
    while kill -0 "$pid" 2>/dev/null && test "$attempt" -lt 100; do
      sleep 0.1
      attempt=$((attempt + 1))
    done
    if kill -0 "$pid" 2>/dev/null; then
      kill -KILL "$pid" 2>/dev/null || true
    fi
  fi
  rm -f -- '__PID_FILE__'
fi
'@
    $command = $template.Replace("__PID_FILE__", $remotePidFile)
    Invoke-SSH -HostName $ServerHost -Command $command -IgnoreStandardError | Out-Null
}

function Set-RemotePerformanceGovernor {
    param([string]$HostName)
    if (-not $SetPerformanceGovernor) {
        return ""
    }
    $command = @'
for path in /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor; do
  if test ! -w "$path"; then
    printf 'CPU governor is not writable: %s\n' "$path" >&2
    exit 1
  fi
  printf '%s=%s\n' "$path" "$(cat "$path")"
  printf performance >"$path"
done
'@
    return Invoke-SSH -HostName $HostName -Command $command -IgnoreStandardError
}

function Restore-RemoteGovernors {
    param([string]$HostName, [string]$Snapshot)
    if ([string]::IsNullOrWhiteSpace($Snapshot)) {
        return
    }
    $commands = foreach ($line in $Snapshot -split "`n") {
        $parts = $line.Split("=", 2)
        if ($parts.Count -ne 2 -or $parts[0] -notmatch '^/sys/devices/system/cpu/cpu[0-9]+/cpufreq/scaling_governor$' -or $parts[1] -notmatch '^[a-z0-9_-]+$') {
            throw "Invalid governor snapshot from ${HostName}: $line"
        }
        "printf '$($parts[1])' >'$($parts[0])'"
    }
    Invoke-SSH -HostName $HostName -Command ($commands -join "`n") -IgnoreStandardError | Out-Null
}

function Start-RemoteServer {
    param([ValidateSet("gnalloy", "gnet", "netty")][string]$Framework)
    Stop-RemoteServer
    $template = @'
cd '__REMOTE_REPO__'
: >'__REMOTE_LOG__'
nohup env SERVER_ADDR='__BIND_ADDR__' SERVER_CPU_SET='__SERVER_CPUS__' SERVER_GOMAXPROCS='__SERVER_GOMAXPROCS__' EVENT_LOOPS='__EVENT_LOOPS__' GNALLOY_WORKERS='__GNALLOY_WORKERS__' GNALLOY_MAX_MESSAGES_PER_READ='__READ_BATCH__' GNALLOY_BOSS_CPU_SET='__BOSS_CPUS__' GNALLOY_WORKER_CPU_SET='__WORKER_CPUS__' SERVER_PID_FILE='__PID_FILE__' GNALLOY_BENCH='__REMOTE_REPO__/external/bin/gnalloy-bench' GNET_BENCH='__REMOTE_REPO__/external/bin/gnet-bench' NETTY_BENCH_JAR='__REMOTE_REPO__/external/bin/netty-bench.jar' bash '__SERVER_SCRIPT__' '__FRAMEWORK__' >'__REMOTE_LOG__' 2>&1 </dev/null &
'@
    $command = $template.Replace("__REMOTE_REPO__", $ServerRepo).
        Replace("__REMOTE_LOG__", $remoteLog).
        Replace("__BIND_ADDR__", $ServerBindAddress).
        Replace("__SERVER_CPUS__", $ServerCPUSet).
        Replace("__SERVER_GOMAXPROCS__", $ServerGoMaxProcs.ToString()).
        Replace("__EVENT_LOOPS__", $EventLoops.ToString()).
        Replace("__GNALLOY_WORKERS__", $GnalloyWorkers.ToString()).
        Replace("__READ_BATCH__", $GnalloyMaxMessagesPerRead.ToString()).
        Replace("__BOSS_CPUS__", $GnalloyBossCPUSet).
        Replace("__WORKER_CPUS__", $GnalloyWorkerCPUSet).
        Replace("__PID_FILE__", $remotePidFile).
        Replace("__SERVER_SCRIPT__", $ServerScript).
        Replace("__FRAMEWORK__", $Framework)
    Invoke-SSH -HostName $ServerHost -Command $command -IgnoreStandardError | Out-Null

    $readyPrefix = "serverReady=true framework=$Framework protocol=udp-echo addr="
    for ($attempt = 0; $attempt -lt 60; $attempt++) {
        Start-Sleep -Milliseconds 250
        $log = Invoke-SSH -HostName $ServerHost -Command "cat '$remoteLog' 2>/dev/null || true" -IgnoreStandardError
        $readyLine = $log -split "`n" | Where-Object { $_.StartsWith($readyPrefix) } | Select-Object -First 1
        if ($null -ne $readyLine) {
            Add-Content -LiteralPath $Output -Value $readyLine -Encoding utf8
            return
        }
        $alive = Invoke-SSH -HostName $ServerHost -Command "test -r '$remotePidFile' && kill -0 `$(cat '$remotePidFile') 2>/dev/null && printf alive || true" -IgnoreStandardError
        if ($alive -ne "alive") {
            throw "Remote $Framework server exited before readiness:`n$log"
        }
    }
    $lastLog = Invoke-SSH -HostName $ServerHost -Command "cat '$remoteLog' 2>/dev/null || true" -IgnoreStandardError
    throw "Remote $Framework server readiness timed out:`n$lastLog"
}

function Invoke-Client {
    param([string]$Framework, [int]$Payload)
    $command = "taskset -c '$ClientCPUSet' env GOMAXPROCS='$ClientGoMaxProcs' '$ClientUdpLoad' -server-framework '$Framework' -addr '$TargetAddress' -payload '$Payload' -connections '$Connections' -messages '$Messages' -warmup-messages '$WarmupMessages' -latency-sample-rate '$LatencySampleRate' -target-rate '$TargetRate' -timeout 5m"
    $result = Invoke-SSH -HostName $ClientHost -Command $command -IgnoreStandardError
    Write-Output $result
    Add-Content -LiteralPath $Output -Value $result -Encoding utf8
}

function Invoke-Case {
    param([string]$Framework, [int]$Payload, [int]$Run)
    $case = "case=$Framework payload=$Payload run=$Run"
    Write-Output $case
    Add-Content -LiteralPath $Output -Value $case -Encoding utf8
    Start-RemoteServer -Framework $Framework
    try {
        Invoke-Client -Framework $Framework -Payload $Payload
    } finally {
        Stop-RemoteServer
    }
    if ($CooldownSeconds -gt 0) {
        Start-Sleep -Seconds $CooldownSeconds
    }
}

foreach ($value in @($ClientHost, $ClientAddress, $ServerHost, $SSHUser, $ServerRepo, $ClientRepo, $ServerScript, $ServerBindAddress, $TargetAddress, $ClientCPUSet, $ServerCPUSet, $GnalloyBossCPUSet, $GnalloyWorkerCPUSet, $ClientUdpLoad)) {
    Assert-SafeRemoteValue -Name "parameter" -Value $value
}
foreach ($payload in $Payloads) {
    if ($payload -le 0) {
        throw "Payload must be positive: $payload"
    }
}
$clientCheck = Invoke-SSH -HostName $ClientHost -Command "test -x '$ClientUdpLoad' && ip -o -4 address show | grep -F ' $ClientAddress/' >/dev/null && taskset -c '$ClientCPUSet' true && printf ready" -IgnoreStandardError
if ($clientCheck -ne "ready") {
    throw "Linux client prerequisites are not satisfied on $ClientHost"
}

$parent = Split-Path -Parent $Output
New-Item -ItemType Directory -Path $parent -Force | Out-Null
$serverGovernorSnapshot = ""
$clientGovernorSnapshot = ""
try {
    $serverGovernorSnapshot = Set-RemotePerformanceGovernor -HostName $ServerHost
    $clientGovernorSnapshot = Set-RemotePerformanceGovernor -HostName $ClientHost
    Set-Content -LiteralPath $Output -Value @(
    "timestamp=$([DateTimeOffset]::Now.ToString('o'))",
    "crossHost=true clientHost=$ClientHost clientAddress=$ClientAddress serverHost=$ServerHost serverAddress=$TargetAddress",
    "connections=$Connections messages=$Messages warmupMessages=$WarmupMessages targetRate=$TargetRate latencySampleRate=$LatencySampleRate repetitions=$Repetitions cooldownSeconds=$CooldownSeconds",
    "clientCPUSet=$ClientCPUSet clientGOMAXPROCS=$ClientGoMaxProcs serverCPUSet=$ServerCPUSet serverGOMAXPROCS=$ServerGoMaxProcs eventLoops=$EventLoops gnalloyWorkers=$GnalloyWorkers gnalloyBossCPUSet=$GnalloyBossCPUSet gnalloyWorkerCPUSet=$GnalloyWorkerCPUSet gnalloyMaxMessagesPerRead=$GnalloyMaxMessagesPerRead performanceGovernor=$SetPerformanceGovernor",
    "excludedFrameworks=netpoll,fasthttp reason=no-comparable-udp-server"
) -Encoding utf8
$metadataTemplate = @'
printf 'serverHostname=%s\n' "$(hostname)"
uname -srmo
awk -F ': ' '/^model name/{print "serverCPU=" $2; exit}' /proc/cpuinfo
printf 'serverGovernor=%s\n' "$(cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor 2>/dev/null || printf unknown)"
sha256sum '__REMOTE_REPO__/external/bin/gnalloy-bench' '__REMOTE_REPO__/external/bin/gnet-bench' '__REMOTE_REPO__/external/bin/netty-bench.jar'
'@
$remoteMetadata = Invoke-SSH -HostName $ServerHost -Command $metadataTemplate.Replace("__REMOTE_REPO__", $ServerRepo) -IgnoreStandardError
Add-Content -LiteralPath $Output -Value $remoteMetadata -Encoding utf8
$clientMetadataTemplate = @'
printf 'clientHostname=%s\n' "$(hostname)"
uname -srmo
awk -F ': ' '/^model name/{print "clientCPU=" $2; exit}' /proc/cpuinfo
printf 'clientGovernor=%s\n' "$(cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor 2>/dev/null || printf unknown)"
sha256sum '__UDP_LOAD__'
'@
$clientMetadata = Invoke-SSH -HostName $ClientHost -Command $clientMetadataTemplate.Replace("__UDP_LOAD__", $ClientUdpLoad) -IgnoreStandardError
Add-Content -LiteralPath $Output -Value $clientMetadata -Encoding utf8
Add-Content -LiteralPath $Output -Value "pingBaseline:" -Encoding utf8
$pingBaseline = Invoke-SSH -HostName $ClientHost -Command "ping -c 10 '$ServerHost'" -IgnoreStandardError
Write-Output $pingBaseline
Add-Content -LiteralPath $Output -Value $pingBaseline -Encoding utf8

    foreach ($payload in $Payloads) {
        for ($run = 1; $run -le $Repetitions; $run++) {
            switch (($run - 1) % 3) {
                0 {
                    Invoke-Case -Framework "gnalloy" -Payload $payload -Run $run
                    Invoke-Case -Framework "gnet" -Payload $payload -Run $run
                    Invoke-Case -Framework "netty" -Payload $payload -Run $run
                }
                1 {
                    Invoke-Case -Framework "gnet" -Payload $payload -Run $run
                    Invoke-Case -Framework "netty" -Payload $payload -Run $run
                    Invoke-Case -Framework "gnalloy" -Payload $payload -Run $run
                }
                2 {
                    Invoke-Case -Framework "netty" -Payload $payload -Run $run
                    Invoke-Case -Framework "gnalloy" -Payload $payload -Run $run
                    Invoke-Case -Framework "gnet" -Payload $payload -Run $run
                }
            }
        }
    }
} finally {
    try {
        Stop-RemoteServer
    } finally {
        try {
            Restore-RemoteGovernors -HostName $ClientHost -Snapshot $clientGovernorSnapshot
        } finally {
            Restore-RemoteGovernors -HostName $ServerHost -Snapshot $serverGovernorSnapshot
        }
    }
}
