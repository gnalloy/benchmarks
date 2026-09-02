[CmdletBinding()]
param(
    [string]$ClientHost = "172.16.8.172",
    [string]$ClientAddress = "172.16.8.172",
    [string]$ServerHost = "172.16.8.171",
    [ValidateRange(1, 65535)]
    [int]$SSHPort = 22,
    [string]$SSHUser = "root",
    [string]$ServerRepo = "/opt/test/gnalloy/benchmarks-http1-server",
    [string]$ClientRepo = "/opt/test/gnalloy/benchmarks-http1-client",
    [string]$ServerScript = "./scripts/run-linux-http1-server.sh",
    [string]$ServerBindAddress = "0.0.0.0:19091",
    [string]$TargetAddress = "172.16.8.171:19091",
    [ValidateSet("http1", "https1")]
    [string[]]$Protocols = @("http1", "https1"),
    [ValidateSet("1.1", "1.2", "1.3")]
    [string[]]$TLSVersions = @("1.1", "1.2", "1.3"),
    [string]$HTTP1ALPN = "http/1.1",
    [string]$TLS11Cipher = "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
    [string]$TLS12Cipher = "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
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
    [ValidateSet("gnalloy", "gnet", "netpoll", "fasthttp", "netty")]
    [string[]]$Frameworks = @("gnalloy", "gnet", "netpoll", "fasthttp", "netty"),
    [ValidateRange(0, 300)]
    [int]$CooldownSeconds = 10,
    [bool]$SetPerformanceGovernor = $true,
    [ValidateRange(1, 64)]
    [int]$ClientGoMaxProcs = 4,
    [string]$ClientCPUSet = "0,1,2,4",
    [string]$ServerCPUSet = "0,1,2,3",
    [ValidateRange(1, 64)]
    [int]$ServerGoMaxProcs = 4,
    [ValidateRange(1, 64)]
    [int]$EventLoops = 4,
    [ValidateRange(1, 64)]
    [int]$GnalloyWorkers = 4,
    [ValidateSet("immediate", "read-complete", "event-loop-batch")]
    [string]$GnalloyFlushStrategy = "read-complete",
    [string]$GnalloyBossCPUSet = "3",
    [string]$GnalloyWorkerCPUSet = "0,1,2,3",
    [switch]$CaptureCPUProfile,
    [switch]$CaptureAllocProfile,
    [switch]$CaptureRuntimeTrace,
    [string]$ProfileOutputDirectory = "",
    [string]$ClientHTTP1Load = "",
    [string]$Output = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($ClientHTTP1Load)) {
    $ClientHTTP1Load = "$ClientRepo/external/bin/http1-load"
}
if ([string]::IsNullOrWhiteSpace($Output)) {
    $mode = if ($TargetRate -gt 0) { "offered-$TargetRate" } else { "saturation" }
    $Output = Join-Path $repoRoot "reports/raw/linux-cross-host-http1-family-$mode.out"
}
if (($CaptureCPUProfile -or $CaptureAllocProfile -or $CaptureRuntimeTrace) -and [string]::IsNullOrWhiteSpace($ProfileOutputDirectory)) {
    $ProfileOutputDirectory = Join-Path $repoRoot "reports/raw/cross-host-http1-profiles"
}
if (@($CaptureCPUProfile, $CaptureAllocProfile, $CaptureRuntimeTrace).Where({ $_ }).Count -gt 1) {
    throw "CPU, allocation, and runtime trace profiles must run separately"
}

$remotePidFile = "/tmp/gnalloy-http1-cross-host.pid"
$remoteLog = "/tmp/gnalloy-http1-cross-host.log"
$remoteProfileDirectory = "$ServerRepo/reports/raw/cross-host-http1-profiles"

function Assert-SafeRemoteValue {
    param([string]$Name, [string]$Value)
    if ([string]::IsNullOrWhiteSpace($Value) -or $Value -notmatch '^[A-Za-z0-9_./,:=-]+$') {
        throw "$Name contains unsupported characters: $Value"
    }
}

function Invoke-SSH {
    param([string]$HostName, [string]$Command, [switch]$IgnoreStandardError)
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
        throw "SSH command failed with exit code $($process.ExitCode): $stderr"
    }
    if (-not $IgnoreStandardError -and -not [string]::IsNullOrWhiteSpace($stderr)) {
        [Console]::Error.WriteLine($stderr)
    }
    return $stdout
}

function Receive-RemoteFile {
    param([string]$RemotePath, [string]$LocalPath)
    $startInfo = [System.Diagnostics.ProcessStartInfo]::new()
    $startInfo.FileName = (Get-Command scp -ErrorAction Stop).Source
    $startInfo.UseShellExecute = $false
    $startInfo.RedirectStandardOutput = $true
    $startInfo.RedirectStandardError = $true
    foreach ($argument in @("-P", $SSHPort.ToString(), "-o", "BatchMode=yes", "-o", "ConnectTimeout=10", "${SSHUser}@${ServerHost}:$RemotePath", $LocalPath)) {
        $startInfo.ArgumentList.Add($argument)
    }
    $process = [System.Diagnostics.Process]::new()
    $process.StartInfo = $startInfo
    if (-not $process.Start()) {
        throw "Failed to start SCP"
    }
    $stdoutTask = $process.StandardOutput.ReadToEndAsync()
    $stderrTask = $process.StandardError.ReadToEndAsync()
    $process.WaitForExit()
    if ($process.ExitCode -ne 0) {
        throw "SCP failed with exit code $($process.ExitCode): $($stderrTask.GetAwaiter().GetResult().TrimEnd())"
    }
    $null = $stdoutTask.GetAwaiter().GetResult()
}

function Stop-RemoteServer {
    $command = @'
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
'@.Replace("__PID_FILE__", $remotePidFile)
    Invoke-SSH -HostName $ServerHost -Command $command -IgnoreStandardError | Out-Null
}

function Set-RemotePerformanceGovernor {
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
    param(
        [string]$Framework,
        [string]$Protocol,
        [string]$TLSVersion,
        [string]$CipherSuites,
        [int]$Payload,
        [string]$CPUProfile = "",
        [string]$AllocProfile = "",
        [string]$RuntimeTrace = ""
    )
    Stop-RemoteServer
    $command = @'
cd '__REMOTE_REPO__'
: >'__REMOTE_LOG__'
nohup env SERVER_ADDR='__BIND_ADDR__' PROTOCOL='__PROTOCOL__' TLS_VERSION='__TLS_VERSION__' ALPN='__ALPN__' CIPHER_SUITES='__CIPHER_SUITES__' PAYLOAD='__PAYLOAD__' SERVER_CPU_SET='__SERVER_CPUS__' SERVER_GOMAXPROCS='__SERVER_GOMAXPROCS__' EVENT_LOOPS='__EVENT_LOOPS__' GNALLOY_WORKERS='__GNALLOY_WORKERS__' GNALLOY_FLUSH_STRATEGY='__FLUSH_STRATEGY__' GNALLOY_BOSS_CPU_SET='__BOSS_CPUS__' GNALLOY_WORKER_CPU_SET='__WORKER_CPUS__' GNALLOY_CPU_PROFILE='__CPU_PROFILE__' GNALLOY_ALLOC_PROFILE='__ALLOC_PROFILE__' GNALLOY_RUNTIME_TRACE='__RUNTIME_TRACE__' FASTHTTP_CPU_PROFILE='__CPU_PROFILE__' SERVER_PID_FILE='__PID_FILE__' GNALLOY_BENCH='__REMOTE_REPO__/external/bin/gnalloy-bench' FASTHTTP_BENCH='__REMOTE_REPO__/external/bin/fasthttp-bench' NETTY_BENCH_JAR='__REMOTE_REPO__/external/bin/netty-bench.jar' bash '__SERVER_SCRIPT__' '__FRAMEWORK__' >'__REMOTE_LOG__' 2>&1 </dev/null &
'@
    $command = $command.Replace("__REMOTE_REPO__", $ServerRepo).
        Replace("__REMOTE_LOG__", $remoteLog).
        Replace("__BIND_ADDR__", $ServerBindAddress).
        Replace("__PROTOCOL__", $Protocol).
        Replace("__TLS_VERSION__", $TLSVersion).
        Replace("__ALPN__", $HTTP1ALPN).
        Replace("__CIPHER_SUITES__", $CipherSuites).
        Replace("__PAYLOAD__", $Payload.ToString()).
        Replace("__SERVER_CPUS__", $ServerCPUSet).
        Replace("__SERVER_GOMAXPROCS__", $ServerGoMaxProcs.ToString()).
        Replace("__EVENT_LOOPS__", $EventLoops.ToString()).
        Replace("__GNALLOY_WORKERS__", $GnalloyWorkers.ToString()).
        Replace("__FLUSH_STRATEGY__", $GnalloyFlushStrategy).
        Replace("__BOSS_CPUS__", $GnalloyBossCPUSet).
        Replace("__WORKER_CPUS__", $GnalloyWorkerCPUSet).
        Replace("__CPU_PROFILE__", $CPUProfile).
        Replace("__ALLOC_PROFILE__", $AllocProfile).
        Replace("__RUNTIME_TRACE__", $RuntimeTrace).
        Replace("__PID_FILE__", $remotePidFile).
        Replace("__SERVER_SCRIPT__", $ServerScript).
        Replace("__FRAMEWORK__", $Framework)
    Invoke-SSH -HostName $ServerHost -Command $command -IgnoreStandardError | Out-Null

    $readyPrefix = "serverReady=true framework=$Framework protocol=$Protocol addr="
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
    throw "Remote $Framework server readiness timed out"
}

function Get-ProfilePaths {
    param([string]$Framework, [string]$Protocol, [string]$TLSVersion, [int]$Payload, [int]$Run)
    $paths = @{}
    if ($Framework -ne "gnalloy" -and $Framework -ne "fasthttp") {
        return $paths
    }
    $mode = if ($TargetRate -gt 0) { "offered-$TargetRate" } else { "saturation" }
    $tlsLabel = if ([string]::IsNullOrWhiteSpace($TLSVersion)) { "plain" } else { "tls$($TLSVersion.Replace('.', ''))" }
    $baseName = "$Protocol-$Framework-$mode-$tlsLabel-p$Payload-r$Run-codec"
    if ($CaptureCPUProfile) {
        $paths.CPUProfile = "$remoteProfileDirectory/$baseName.cpu.pprof"
    }
    if ($Framework -eq "fasthttp") {
        return $paths
    }
    if ($CaptureAllocProfile) {
        $paths.AllocProfile = "$remoteProfileDirectory/$baseName.alloc.pprof"
    }
    if ($CaptureRuntimeTrace) {
        $paths.RuntimeTrace = "$remoteProfileDirectory/$baseName.trace"
    }
    return $paths
}

function Receive-Profiles {
    param([hashtable]$Paths)
    if ($Paths.Count -eq 0) {
        return
    }
    New-Item -ItemType Directory -Path $ProfileOutputDirectory -Force | Out-Null
    foreach ($remotePath in $Paths.Values) {
        $remoteCheck = Invoke-SSH -HostName $ServerHost -Command "test -s '$remotePath' && printf ready" -IgnoreStandardError
        if ($remoteCheck -ne "ready") {
            throw "Remote profile is missing or empty: $remotePath"
        }
        $localPath = Join-Path $ProfileOutputDirectory ([System.IO.Path]::GetFileName($remotePath))
        Receive-RemoteFile -RemotePath $remotePath -LocalPath $localPath
    }
}

function Invoke-Client {
    param([string]$Framework, [string]$Protocol, [string]$TLSVersion, [string]$CipherSuites, [int]$Payload)
    $command = "taskset -c '$ClientCPUSet' env GOMAXPROCS='$ClientGoMaxProcs' '$ClientHTTP1Load' -server-framework '$Framework' -addr '$TargetAddress' -payload '$Payload' -connections '$Connections' -messages '$Messages' -warmup-messages '$WarmupMessages' -latency-sample-rate '$LatencySampleRate' -target-rate '$TargetRate' -timeout 5m"
    if ($Protocol -eq "https1") {
        $command += " -tls -tls-version '$TLSVersion' -alpn '$HTTP1ALPN' -insecure-skip-verify"
        if (-not [string]::IsNullOrWhiteSpace($CipherSuites)) {
            $command += " -cipher-suites '$CipherSuites'"
        }
    }
    $result = Invoke-SSH -HostName $ClientHost -Command $command -IgnoreStandardError
    Write-Output $result
    Add-Content -LiteralPath $Output -Value $result -Encoding utf8
}

function Invoke-Case {
    param([string]$Framework, [string]$Protocol, [string]$TLSVersion, [string]$CipherSuites, [int]$Payload, [int]$Run)
    if ($Framework -eq "gnet" -or $Framework -eq "netpoll") {
        $reason = if ($Framework -eq "gnet") {
            "gnet does not provide an HTTP codec; benchmark-owned parsing is prohibited"
        } else {
            "CloudWeGo netpoll does not provide an HTTP codec; benchmark-owned parsing is prohibited"
        }
        $case = "case=$Framework protocol=$Protocol tlsVersion=$TLSVersion payload=$Payload run=$Run status=N/A reason=$reason"
        Write-Output $case
        Add-Content -LiteralPath $Output -Value $case -Encoding utf8
        return
    }
    $case = "case=$Framework protocol=$Protocol tlsVersion=$TLSVersion cipherSuites=$CipherSuites payload=$Payload run=$Run"
    Write-Output $case
    Add-Content -LiteralPath $Output -Value $case -Encoding utf8
    $profiles = Get-ProfilePaths -Framework $Framework -Protocol $Protocol -TLSVersion $TLSVersion -Payload $Payload -Run $Run
    foreach ($remotePath in $profiles.Values) {
        Invoke-SSH -HostName $ServerHost -Command "rm -f -- '$remotePath'" -IgnoreStandardError | Out-Null
    }
    Start-RemoteServer -Framework $Framework -Protocol $Protocol -TLSVersion $TLSVersion -CipherSuites $CipherSuites -Payload $Payload -CPUProfile $profiles["CPUProfile"] -AllocProfile $profiles["AllocProfile"] -RuntimeTrace $profiles["RuntimeTrace"]
    try {
        Invoke-Client -Framework $Framework -Protocol $Protocol -TLSVersion $TLSVersion -CipherSuites $CipherSuites -Payload $Payload
    } finally {
        try {
            Stop-RemoteServer
        } finally {
            Receive-Profiles -Paths $profiles
        }
    }
    if ($CooldownSeconds -gt 0) {
        Start-Sleep -Seconds $CooldownSeconds
    }
}

function Get-CipherSuites {
    param([string]$Protocol, [string]$TLSVersion)
    if ($Protocol -ne "https1") {
        return ""
    }
    switch ($TLSVersion) {
        "1.1" { return $TLS11Cipher }
        "1.2" { return $TLS12Cipher }
        "1.3" { return "" }
        default { throw "unsupported TLS version: $TLSVersion" }
    }
}

function Get-ProtocolCases {
    foreach ($protocol in $Protocols) {
        if ($protocol -eq "http1") {
            [pscustomobject]@{ Protocol = $protocol; TLSVersion = ""; CipherSuites = "" }
            continue
        }
        foreach ($tlsVersion in $TLSVersions) {
            [pscustomobject]@{
                Protocol     = $protocol
                TLSVersion   = $tlsVersion
                CipherSuites = Get-CipherSuites -Protocol $protocol -TLSVersion $tlsVersion
            }
        }
    }
}

function Get-RotatedFrameworks {
    param([int]$Run)
    $base = @("gnalloy", "gnet", "netpoll", "fasthttp", "netty")
    $offset = ($Run - 1) % $base.Count
    for ($index = 0; $index -lt $base.Count; $index++) {
        $base[($index + $offset) % $base.Count]
    }
}

foreach ($value in @($ClientHost, $ClientAddress, $ServerHost, $SSHUser, $ServerRepo, $ClientRepo, $ServerScript, $ServerBindAddress, $TargetAddress, $HTTP1ALPN, $TLS11Cipher, $TLS12Cipher, $ClientCPUSet, $ServerCPUSet, $GnalloyFlushStrategy, $GnalloyBossCPUSet, $GnalloyWorkerCPUSet, $ClientHTTP1Load, $remoteProfileDirectory)) {
    Assert-SafeRemoteValue -Name "parameter" -Value $value
}
foreach ($payload in $Payloads) {
    if ($payload -le 0) {
        throw "Payload must be positive: $payload"
    }
}
$clientCheck = Invoke-SSH -HostName $ClientHost -Command "test -x '$ClientHTTP1Load' && ip -o -4 address show | grep -F ' $ClientAddress/' >/dev/null && taskset -c '$ClientCPUSet' true && printf ready" -IgnoreStandardError
if ($clientCheck -ne "ready") {
    throw "Linux client prerequisites are not satisfied on $ClientHost"
}
$serverCheck = Invoke-SSH -HostName $ServerHost -Command "test -f '$ServerRepo/$ServerScript' && taskset -c '$ServerCPUSet' true && printf ready" -IgnoreStandardError
if ($serverCheck -ne "ready") {
    throw "Linux server prerequisites are not satisfied on $ServerHost"
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
        "protocols=$($Protocols -join ',') tlsVersions=$($TLSVersions -join ',') alpn=$HTTP1ALPN tls11Cipher=$TLS11Cipher tls12Cipher=$TLS12Cipher connections=$Connections messages=$Messages warmupMessages=$WarmupMessages targetRate=$TargetRate latencySampleRate=$LatencySampleRate repetitions=$Repetitions cooldownSeconds=$CooldownSeconds frameworks=$($Frameworks -join ',')",
        "clientCPUSet=$ClientCPUSet clientGOMAXPROCS=$ClientGoMaxProcs serverCPUSet=$ServerCPUSet serverGOMAXPROCS=$ServerGoMaxProcs eventLoops=$EventLoops gnalloyWorkers=$GnalloyWorkers gnalloyFlushStrategy=$GnalloyFlushStrategy gnalloyHTTP1Pipeline=tcp+channel+codec-http1+handler performanceGovernor=$SetPerformanceGovernor",
        "unsupportedFrameworks=gnet,netpoll status=N/A reason=no-framework-http-codec"
    ) -Encoding utf8
    $serverMetadataCommand = @'
printf 'serverHostname=%s\n' "$(hostname)"
uname -srmo
awk -F ': ' '/^model name/{print "serverCPU=" $2; exit}' /proc/cpuinfo
printf 'serverGovernor=%s\n' "$(cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor 2>/dev/null || printf unknown)"
printf 'serverCPUTopology='; lscpu -p=CPU,CORE,SOCKET | grep -v '^#' | tr '\n' ';'; printf '\n'
sha256sum '__SERVER_REPO__/external/bin/gnalloy-bench' '__SERVER_REPO__/external/bin/fasthttp-bench' '__SERVER_REPO__/external/bin/netty-bench.jar'
'@.Replace("__SERVER_REPO__", $ServerRepo)
    $serverMetadata = Invoke-SSH -HostName $ServerHost -Command $serverMetadataCommand -IgnoreStandardError
    Add-Content -LiteralPath $Output -Value $serverMetadata -Encoding utf8
    $clientMetadataCommand = @'
printf 'clientHostname=%s\n' "$(hostname)"
uname -srmo
awk -F ': ' '/^model name/{print "clientCPU=" $2; exit}' /proc/cpuinfo
printf 'clientGovernor=%s\n' "$(cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor 2>/dev/null || printf unknown)"
printf 'clientCPUTopology='; lscpu -p=CPU,CORE,SOCKET | grep -v '^#' | tr '\n' ';'; printf '\n'
sha256sum '__HTTP1_LOAD__'
'@.Replace("__HTTP1_LOAD__", $ClientHTTP1Load)
    $clientMetadata = Invoke-SSH -HostName $ClientHost -Command $clientMetadataCommand -IgnoreStandardError
    Add-Content -LiteralPath $Output -Value $clientMetadata -Encoding utf8
    $pingBaseline = Invoke-SSH -HostName $ClientHost -Command "ping -c 10 '$ServerHost'" -IgnoreStandardError
    Add-Content -LiteralPath $Output -Value @("pingBaseline:", $pingBaseline) -Encoding utf8

    foreach ($protocolCase in (Get-ProtocolCases)) {
        foreach ($payload in $Payloads) {
            for ($run = 1; $run -le $Repetitions; $run++) {
                foreach ($framework in (Get-RotatedFrameworks -Run $run)) {
                    if ($Frameworks -contains $framework) {
                        Invoke-Case -Framework $framework -Protocol $protocolCase.Protocol -TLSVersion $protocolCase.TLSVersion -CipherSuites $protocolCase.CipherSuites -Payload $payload -Run $run
                    }
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
