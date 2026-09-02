package main

import "gnalloy.org/gnalloy/transport"

const (
	linuxNativePollerAutoWorkerLimit = 4
	windowsIOCPAutoWorkerLimit       = 8
)

type workerSizingInput struct {
	GOOS       string
	Backend    transport.BackendKind
	GOMAXPROCS int
	Bosses     int
}

func defaultWorkerCount(input workerSizingInput) int {
	workers := availableWorkerBudget(input.GOMAXPROCS, input.Bosses)
	if input.GOOS == "linux" && isLinuxNativePoller(input.Backend) && workers > linuxNativePollerAutoWorkerLimit {
		return linuxNativePollerAutoWorkerLimit
	}
	if input.GOOS == "windows" && input.Backend == transport.BackendIOCP && workers > windowsIOCPAutoWorkerLimit {
		return windowsIOCPAutoWorkerLimit
	}
	return workers
}

func availableWorkerBudget(gomaxprocs int, bosses int) int {
	if gomaxprocs <= 0 {
		gomaxprocs = 1
	}
	if bosses < 0 {
		bosses = 0
	}
	workers := gomaxprocs - bosses
	if workers <= 0 {
		return 1
	}
	return workers
}

func isLinuxNativePoller(backend transport.BackendKind) bool {
	return backend == transport.BackendEpoll || backend == transport.BackendIOUring
}
