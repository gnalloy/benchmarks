[CmdletBinding()]
param(
    [string]$ServerHost = "172.16.8.171",
    [string]$ClientHost = "172.16.8.172",
    [string]$SSHUser = "root",
    [int]$SSHPort = 22,
    [string]$ServerRepo = "/opt/test/gnalloy/http2-server-20260903/benchmarks",
    [string]$ServerBindAddress = "0.0.0.0:19092",
    [string]$TargetAddress = "172.16.8.171:19092",
    [string]$ServerCPUSet = "0,1,2,3",
    [string]$ClientCPUSet = "0,1,2,3",
    [int]$ServerGoMaxProcs = 4,
    [int]$ClientGoMaxProcs = 4,
    [int]$EventLoops = 3,
    [int]$GnalloyWorkers = 3,
    [int]$GnalloyMaxMessagesPerRead = 32,
    [ValidateSet("http2", "https2-tls12", "https2-tls13")]
    [string[]]$Scenarios = @("http2", "https2-tls12", "https2-tls13"),
    [int[]]$Payloads = @(128, 1024),
    [int]$Connections = 64,
    [int]$Messages = 20000,
    [int]$WarmupMessages = 1000,
    [int]$LatencySampleRate = 64,
    [int]$Repetitions = 5,
    [int]$CooldownSeconds = 5,
    [bool]$SetPerformanceGovernor = $true,
    [ValidateSet("gnalloy", "hertz", "netty", "gnet", "netpoll", "fasthttp")]
    [string[]]$Frameworks = @("gnalloy", "hertz", "netty", "gnet", "netpoll", "fasthttp"),
    [string]$GnalloyBench = "/opt/test/gnalloy/http2-server-20260903/gnalloy-bench",
    [string]$HertzBench = "/opt/test/gnalloy/http2-server-20260903/hertz-bench",
    [string]$NettyBenchJar = "/opt/test/gnalloy/http2-server-20260903/netty-bench.jar",
    [string]$ClientBench = "/opt/test/gnalloy/http2-client-20260903/gnalloy-bench",
    [string]$Output = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($Output)) {
    $Output = Join-Path $repoRoot "reports/raw/linux-cross-host-http2-saturation.out"
}
$remotePIDFile = "/tmp/gnalloy-http2-cross-host.pid"
$remoteLog = "/tmp/gnalloy-http2-cross-host.log"
$supportedFrameworks = @("gnalloy", "hertz", "netty")

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
    param([string]$Framework, [string]$Protocol, [string]$TLSVersion, [int]$Payload)
    Stop-Server
    $command = @'
cd '__SERVER_REPO__'
: >'__LOG__'
nohup env PROTOCOL='__PROTOCOL__' TLS_VERSION='__TLS_VERSION__' PAYLOAD='__PAYLOAD__' SERVER_ADDR='__BIND_ADDR__' SERVER_CPU_SET='__SERVER_CPUS__' SERVER_GOMAXPROCS='__SERVER_PROCS__' EVENT_LOOPS='__EVENT_LOOPS__' GNALLOY_WORKERS='__WORKERS__' GNALLOY_MAX_MESSAGES_PER_READ='__MAX_MESSAGES_PER_READ__' GNALLOY_BOSS_CPU_SET='3' GNALLOY_WORKER_CPU_SET='0,1,2' SERVER_PID_FILE='__PID_FILE__' GNALLOY_BENCH='__GNALLOY_BENCH__' HERTZ_BENCH='__HERTZ_BENCH__' NETTY_BENCH_JAR='__NETTY_JAR__' bash ./scripts/run-linux-http2-server.sh '__FRAMEWORK__' >'__LOG__' 2>&1 </dev/null &
'@
    $command = $command.Replace("__SERVER_REPO__", $ServerRepo).
        Replace("__LOG__", $remoteLog).
        Replace("__PROTOCOL__", $Protocol).
        Replace("__TLS_VERSION__", $TLSVersion).
        Replace("__PAYLOAD__", $Payload.ToString()).
        Replace("__BIND_ADDR__", $ServerBindAddress).
        Replace("__SERVER_CPUS__", $ServerCPUSet).
        Replace("__SERVER_PROCS__", $ServerGoMaxProcs.ToString()).
        Replace("__EVENT_LOOPS__", $EventLoops.ToString()).
        Replace("__WORKERS__", $GnalloyWorkers.ToString()).
        Replace("__MAX_MESSAGES_PER_READ__", $GnalloyMaxMessagesPerRead.ToString()).
        Replace("__PID_FILE__", $remotePIDFile).
        Replace("__GNALLOY_BENCH__", $GnalloyBench).
        Replace("__HERTZ_BENCH__", $HertzBench).
        Replace("__NETTY_JAR__", $NettyBenchJar).
        Replace("__FRAMEWORK__", $Framework)
    Invoke-SSH -HostName $ServerHost -Command $command | Out-Null

    $readyPrefix = "serverReady=true framework=$Framework protocol=$Protocol addr="
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
    param([string]$Protocol, [string]$TLSVersion, [int]$Payload)
    $command = "taskset -c '$ClientCPUSet' env GOMAXPROCS='$ClientGoMaxProcs' '$ClientBench' -protocol '$Protocol' -client-only=true -addr '$TargetAddress' -payload '$Payload' -connections '$Connections' -messages '$Messages' -warmup-messages '$WarmupMessages' -latency-sample-rate '$LatencySampleRate' -backend epoll -boss 1 -workers 3 -timeout 5m"
    if ($Protocol -eq "https2") {
        $command += " -tls-version '$TLSVersion' -alpn h2"
    }
    $result = Invoke-SSH -HostName $ClientHost -Command $command
    Write-Output $result
    Add-Content -LiteralPath $Output -Value $result -Encoding utf8
}

function Get-ProtocolCases {
    if ($Scenarios -contains "http2") {
        [pscustomobject]@{ Protocol = "http2"; TLSVersion = "" }
    }
    if ($Scenarios -contains "https2-tls12") {
        [pscustomobject]@{ Protocol = "https2"; TLSVersion = "1.2" }
    }
    if ($Scenarios -contains "https2-tls13") {
        [pscustomobject]@{ Protocol = "https2"; TLSVersion = "1.3" }
    }
}

function Get-RotatedFrameworks {
    param([int]$Run)
    $offset = ($Run - 1) % $supportedFrameworks.Count
    for ($index = 0; $index -lt $supportedFrameworks.Count; $index++) {
        $supportedFrameworks[($index + $offset) % $supportedFrameworks.Count]
    }
}

foreach ($value in @($ServerHost, $ClientHost, $SSHUser, $ServerRepo, $ServerBindAddress, $TargetAddress, $ServerCPUSet, $ClientCPUSet, $GnalloyBench, $HertzBench, $NettyBenchJar, $ClientBench)) {
    Assert-SafeValue -Name "parameter" -Value $value
}
if ($Scenarios.Count -eq 0 -or $Connections -le 0 -or $Messages -le 0 -or $WarmupMessages -lt 0 -or $LatencySampleRate -lt 0 -or $Repetitions -le 0 -or $CooldownSeconds -lt 0 -or $GnalloyMaxMessagesPerRead -le 0) {
    throw "Load parameters are out of range"
}
foreach ($payload in $Payloads) {
    if ($payload -le 0) {
        throw "Payload must be positive: $payload"
    }
}

$parent = Split-Path -Parent $Output
New-Item -ItemType Directory -Path $parent -Force | Out-Null
$serverCheck = "test -f '$ServerRepo/scripts/run-linux-http2-server.sh' && test -x '$GnalloyBench' && test -x '$HertzBench' && test -f '$NettyBenchJar' && taskset -c '$ServerCPUSet' true && printf ready"
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
        "scenarios=$($Scenarios -join ',') payloads=$($Payloads -join ',') connections=$Connections messages=$Messages warmupMessages=$WarmupMessages latencySampleRate=$LatencySampleRate repetitions=$Repetitions cooldownSeconds=$CooldownSeconds",
        "serverCPUSet=$ServerCPUSet serverGOMAXPROCS=$ServerGoMaxProcs clientCPUSet=$ClientCPUSet clientGOMAXPROCS=$ClientGoMaxProcs gnalloyMaxMessagesPerRead=$GnalloyMaxMessagesPerRead commonClient=gnalloy-tcp+handler-tls+codec-http2 performanceGovernor=$SetPerformanceGovernor",
        "unsupportedFrameworks=gnet,netpoll,fasthttp status=N/A reason=no-native-http2-codec"
    ) -Encoding utf8

    $serverMetadata = Invoke-SSH -HostName $ServerHost -Command "hostname; uname -srmo; awk -F ': ' '/^model name/{print `"serverCPU=`" `$2; exit}' /proc/cpuinfo; sha256sum '$GnalloyBench' '$HertzBench' '$NettyBenchJar'"
    $clientMetadata = Invoke-SSH -HostName $ClientHost -Command "hostname; uname -srmo; awk -F ': ' '/^model name/{print `"clientCPU=`" `$2; exit}' /proc/cpuinfo; sha256sum '$ClientBench'; ping -c 10 '$ServerHost'"
    Add-Content -LiteralPath $Output -Value @($serverMetadata, $clientMetadata) -Encoding utf8

    foreach ($protocolCase in (Get-ProtocolCases)) {
        foreach ($payload in $Payloads) {
            foreach ($unsupported in @("gnet", "netpoll", "fasthttp")) {
                if ($Frameworks -contains $unsupported) {
                    Add-Content -LiteralPath $Output -Value "case=$unsupported protocol=$($protocolCase.Protocol) tlsVersion=$($protocolCase.TLSVersion) payload=$payload status=N/A reason=no-native-http2-codec" -Encoding utf8
                }
            }
            for ($run = 1; $run -le $Repetitions; $run++) {
                foreach ($framework in (Get-RotatedFrameworks -Run $run)) {
                    if ($Frameworks -notcontains $framework) {
                        continue
                    }
                    $case = "case=$framework protocol=$($protocolCase.Protocol) tlsVersion=$($protocolCase.TLSVersion) payload=$payload run=$run"
                    Write-Output $case
                    Add-Content -LiteralPath $Output -Value $case -Encoding utf8
                    Start-Server -Framework $framework -Protocol $protocolCase.Protocol -TLSVersion $protocolCase.TLSVersion -Payload $payload
                    try {
                        Invoke-Client -Protocol $protocolCase.Protocol -TLSVersion $protocolCase.TLSVersion -Payload $payload
                    } finally {
                        Stop-Server
                    }
                    if ($CooldownSeconds -gt 0) {
                        Start-Sleep -Seconds $CooldownSeconds
                    }
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
