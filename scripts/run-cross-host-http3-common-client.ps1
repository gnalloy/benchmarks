[CmdletBinding()]
param(
    [string]$ServerHost = "172.16.8.171",
    [string]$ClientHost = "172.16.8.172",
    [string]$SSHUser = "root",
    [int]$SSHPort = 22,
    [string]$ServerRepo = "/opt/test/gnalloy/http3-server-20260903/benchmarks",
    [string]$ServerBindAddress = "0.0.0.0:19093",
    [string]$TargetAddress = "172.16.8.171:19093",
    [string]$ServerCPUSet = "0,1,2,3",
    [string]$ClientCPUSet = "0,1,2,3",
    [int]$ServerGoMaxProcs = 4,
    [int]$ClientGoMaxProcs = 4,
    [int]$EventLoops = 3,
    [int[]]$Payloads = @(128, 1024),
    [int]$Connections = 64,
    [int]$Messages = 20000,
    [int]$WarmupMessages = 1000,
    [int]$LatencySampleRate = 64,
    [long]$TargetRate = 0,
    [int]$Repetitions = 5,
    [int]$CooldownSeconds = 5,
    [bool]$SetPerformanceGovernor = $true,
    [ValidateSet("gnalloy", "netty", "hertz", "gnet", "netpoll", "fasthttp")]
    [string[]]$Frameworks = @("gnalloy", "netty", "hertz", "gnet", "netpoll", "fasthttp"),
    [string]$GnalloyBench = "/opt/test/gnalloy/http3-server-20260903/gnalloy-bench",
    [string]$NettyBenchJar = "/opt/test/gnalloy/http3-server-20260903/netty-bench.jar",
    [string]$ClientBench = "/opt/test/gnalloy/http3-client-20260903/gnalloy-bench",
    [string]$Output = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($Output)) {
    $mode = if ($TargetRate -gt 0) { "rate-$TargetRate" } else { "saturation" }
    $Output = Join-Path $repoRoot "reports/raw/linux-cross-host-http3-$mode.out"
}
$remotePIDFile = "/tmp/gnalloy-http3-cross-host.pid"
$remoteLog = "/tmp/gnalloy-http3-cross-host.log"
$supportedFrameworks = @("gnalloy", "netty")
$unsupportedFrameworks = @("hertz", "gnet", "netpoll", "fasthttp")

function Assert-SafeValue {
    param([string]$Name, [string]$Value)
    if ([string]::IsNullOrWhiteSpace($Value) -or $Value -notmatch '^[A-Za-z0-9_./,:=-]+$') {
        throw "$Name contains unsupported characters: $Value"
    }
}

function Invoke-SSH {
    param([string]$HostName, [string]$Command)
    $startInfo = [System.Diagnostics.ProcessStartInfo]::new()
    $startInfo.FileName = (Get-Command ssh -ErrorAction Stop).Source
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
        throw "SSH command failed on ${HostName}: $stderr"
    }
    return $stdout
}

function Stop-Server {
    $command = @'
if test -r '__PID_FILE__'; then
  pid="$(cat '__PID_FILE__')"
  kill -TERM "$pid" 2>/dev/null || true
  attempt=0
  while kill -0 "$pid" 2>/dev/null && test "$attempt" -lt 100; do
    sleep 0.1
    attempt=$((attempt + 1))
  done
  kill -KILL "$pid" 2>/dev/null || true
  rm -f -- '__PID_FILE__'
fi
'@.Replace("__PID_FILE__", $remotePIDFile)
    Invoke-SSH -HostName $ServerHost -Command $command | Out-Null
}

function Set-PerformanceGovernor {
    param([string]$HostName)
    if (-not $SetPerformanceGovernor) {
        return ""
    }
    $command = @'
for path in /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor; do
  test -w "$path" || { printf 'CPU governor is not writable: %s\n' "$path" >&2; exit 1; }
  printf '%s=%s\n' "$path" "$(cat "$path")"
  printf performance >"$path"
done
'@
    return Invoke-SSH -HostName $HostName -Command $command
}

function Restore-Governor {
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
    Invoke-SSH -HostName $HostName -Command ($commands -join "`n") | Out-Null
}

function Start-Server {
    param([string]$Framework, [int]$Payload)
    Stop-Server
    $command = @'
cd '__SERVER_REPO__'
: >'__LOG__'
nohup env PAYLOAD='__PAYLOAD__' SERVER_ADDR='__BIND_ADDR__' SERVER_CPU_SET='__SERVER_CPUS__' SERVER_GOMAXPROCS='__SERVER_PROCS__' EVENT_LOOPS='__EVENT_LOOPS__' SERVER_PID_FILE='__PID_FILE__' GNALLOY_BENCH='__GNALLOY_BENCH__' NETTY_BENCH_JAR='__NETTY_JAR__' bash ./scripts/run-linux-http3-server.sh '__FRAMEWORK__' >'__LOG__' 2>&1 </dev/null &
'@
    $command = $command.Replace("__SERVER_REPO__", $ServerRepo).
        Replace("__LOG__", $remoteLog).
        Replace("__PAYLOAD__", $Payload.ToString()).
        Replace("__BIND_ADDR__", $ServerBindAddress).
        Replace("__SERVER_CPUS__", $ServerCPUSet).
        Replace("__SERVER_PROCS__", $ServerGoMaxProcs.ToString()).
        Replace("__EVENT_LOOPS__", $EventLoops.ToString()).
        Replace("__PID_FILE__", $remotePIDFile).
        Replace("__GNALLOY_BENCH__", $GnalloyBench).
        Replace("__NETTY_JAR__", $NettyBenchJar).
        Replace("__FRAMEWORK__", $Framework)
    Invoke-SSH -HostName $ServerHost -Command $command | Out-Null

    $readyPrefix = "serverReady=true framework=$Framework protocol=http3 addr="
    for ($attempt = 0; $attempt -lt 80; $attempt++) {
        Start-Sleep -Milliseconds 250
        $log = Invoke-SSH -HostName $ServerHost -Command "cat '$remoteLog' 2>/dev/null || true"
        $ready = $log -split "`n" | Where-Object { $_.StartsWith($readyPrefix) } | Select-Object -First 1
        if ($null -ne $ready) {
            Add-Content -LiteralPath $Output -Value $ready -Encoding utf8
            return
        }
        $alive = Invoke-SSH -HostName $ServerHost -Command "test -r '$remotePIDFile' && kill -0 `$(cat '$remotePIDFile') 2>/dev/null && printf alive || true"
        if ($alive -ne "alive") {
            throw "Remote $Framework server exited before readiness:`n$log"
        }
    }
    throw "Remote $Framework server readiness timed out"
}

function Invoke-Client {
    param([int]$Payload)
    $command = "taskset -c '$ClientCPUSet' env GOMAXPROCS='$ClientGoMaxProcs' '$ClientBench' -protocol http3 -client-only=true -addr '$TargetAddress' -payload '$Payload' -connections '$Connections' -messages '$Messages' -warmup-messages '$WarmupMessages' -latency-sample-rate '$LatencySampleRate' -target-rate '$TargetRate' -tls-version 1.3 -alpn h3 -timeout 5m"
    $result = Invoke-SSH -HostName $ClientHost -Command $command
    Write-Output $result
    Add-Content -LiteralPath $Output -Value $result -Encoding utf8
}

function Get-RotatedFrameworks {
    param([int]$Run)
    $offset = ($Run - 1) % $supportedFrameworks.Count
    for ($index = 0; $index -lt $supportedFrameworks.Count; $index++) {
        $supportedFrameworks[($index + $offset) % $supportedFrameworks.Count]
    }
}

foreach ($value in @($ServerHost, $ClientHost, $SSHUser, $ServerRepo, $ServerBindAddress, $TargetAddress, $ServerCPUSet, $ClientCPUSet, $GnalloyBench, $NettyBenchJar, $ClientBench)) {
    Assert-SafeValue -Name "parameter" -Value $value
}
if ($Payloads.Count -eq 0 -or $Connections -le 0 -or $Messages -le 0 -or $WarmupMessages -lt 0 -or $LatencySampleRate -lt 0 -or $TargetRate -lt 0 -or $Repetitions -le 0 -or $CooldownSeconds -lt 0) {
    throw "Load parameters are out of range"
}
foreach ($payload in $Payloads) {
    if ($payload -le 0) {
        throw "Payload must be positive: $payload"
    }
}

$parent = Split-Path -Parent $Output
New-Item -ItemType Directory -Path $parent -Force | Out-Null
$serverCheck = "test -f '$ServerRepo/scripts/run-linux-http3-server.sh' && test -x '$GnalloyBench' && test -f '$NettyBenchJar' && test -x /opt/software/java21/bin/java && taskset -c '$ServerCPUSet' true && printf ready"
if ((Invoke-SSH -HostName $ServerHost -Command $serverCheck) -ne "ready") {
    throw "Server prerequisites are not satisfied"
}
$clientCheck = "test -x '$ClientBench' && taskset -c '$ClientCPUSet' true && printf ready"
if ((Invoke-SSH -HostName $ClientHost -Command $clientCheck) -ne "ready") {
    throw "Client prerequisites are not satisfied"
}

$serverGovernorSnapshot = ""
$clientGovernorSnapshot = ""
try {
    $serverGovernorSnapshot = Set-PerformanceGovernor -HostName $ServerHost
    $clientGovernorSnapshot = Set-PerformanceGovernor -HostName $ClientHost
    Set-Content -LiteralPath $Output -Value @(
        "timestamp=$([DateTimeOffset]::Now.ToString('o'))",
        "crossHost=true serverHost=$ServerHost clientHost=$ClientHost target=$TargetAddress",
        "protocol=http3 tlsVersion=1.3 alpn=h3 payloads=$($Payloads -join ',') connections=$Connections messages=$Messages warmupMessages=$WarmupMessages latencySampleRate=$LatencySampleRate targetRate=$TargetRate repetitions=$Repetitions cooldownSeconds=$CooldownSeconds",
        "serverCPUSet=$ServerCPUSet serverGOMAXPROCS=$ServerGoMaxProcs clientCPUSet=$ClientCPUSet clientGOMAXPROCS=$ClientGoMaxProcs commonClient=gnalloy-quic+transport-http3+codec-http3 performanceGovernor=$SetPerformanceGovernor",
        "unsupportedFrameworks=hertz,gnet,netpoll,fasthttp status=N/A reason=no-native-http3-codec"
    ) -Encoding utf8

    $serverMetadata = Invoke-SSH -HostName $ServerHost -Command "hostname; uname -srmo; awk -F ': ' '/^model name/{print `"serverCPU=`" `$2; exit}' /proc/cpuinfo; sha256sum '$GnalloyBench' '$NettyBenchJar'"
    $clientMetadata = Invoke-SSH -HostName $ClientHost -Command "hostname; uname -srmo; awk -F ': ' '/^model name/{print `"clientCPU=`" `$2; exit}' /proc/cpuinfo; sha256sum '$ClientBench'; ping -c 10 '$ServerHost'"
    Add-Content -LiteralPath $Output -Value @($serverMetadata, $clientMetadata) -Encoding utf8

    foreach ($payload in $Payloads) {
        foreach ($unsupported in $unsupportedFrameworks) {
            if ($Frameworks -contains $unsupported) {
                Add-Content -LiteralPath $Output -Value "case=$unsupported protocol=http3 tlsVersion=1.3 payload=$payload status=N/A reason=no-native-http3-codec" -Encoding utf8
            }
        }
        for ($run = 1; $run -le $Repetitions; $run++) {
            foreach ($framework in (Get-RotatedFrameworks -Run $run)) {
                if ($Frameworks -notcontains $framework) {
                    continue
                }
                $case = "case=$framework protocol=http3 tlsVersion=1.3 payload=$payload run=$run"
                Write-Output $case
                Add-Content -LiteralPath $Output -Value $case -Encoding utf8
                Start-Server -Framework $framework -Payload $payload
                try {
                    Invoke-Client -Payload $payload
                } finally {
                    Stop-Server
                }
                if ($CooldownSeconds -gt 0) {
                    Start-Sleep -Seconds $CooldownSeconds
                }
            }
        }
    }
} finally {
    try {
        Stop-Server
    } finally {
        try {
            Restore-Governor -HostName $ClientHost -Snapshot $clientGovernorSnapshot
        } finally {
            Restore-Governor -HostName $ServerHost -Snapshot $serverGovernorSnapshot
        }
    }
}
