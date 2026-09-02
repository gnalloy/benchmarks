package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseCPUSet(value string) ([]int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parts := strings.Split(value, ",")
	cpus := make([]int, 0, len(parts))
	for _, part := range parts {
		cpu, err := parseCPU(part)
		if err != nil {
			return nil, err
		}
		cpus = append(cpus, cpu)
	}
	return cpus, nil
}

func parseCPU(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, fmt.Errorf("%w: empty CPU id", errInvalidConfig)
	}
	cpu, err := strconv.Atoi(value)
	if err != nil || cpu < 0 {
		return 0, fmt.Errorf("%w: invalid CPU id %q", errInvalidConfig, value)
	}
	return cpu, nil
}

func validateEventLoopTopology(bosses int, workers int, gomaxprocs int, bossCPUs []int, workerCPUs []int) error {
	if bosses > maxInt-workers || bosses+workers > gomaxprocs {
		return fmt.Errorf("%w: boss and worker event loops must not exceed GOMAXPROCS", errInvalidConfig)
	}
	if err := validateDedicatedCPUSet("boss", bosses, bossCPUs); err != nil {
		return err
	}
	if err := validateDedicatedCPUSet("worker", workers, workerCPUs); err != nil {
		return err
	}
	if overlapsCPUSet(bossCPUs, workerCPUs) {
		return fmt.Errorf("%w: boss and worker CPU sets must not overlap", errInvalidConfig)
	}
	return nil
}

func validateDedicatedCPUSet(name string, loops int, cpus []int) error {
	if len(cpus) == 0 {
		return nil
	}
	if len(cpus) < loops {
		return fmt.Errorf("%w: %s CPU set must provide at least one CPU per event loop", errInvalidConfig, name)
	}
	seen := make(map[int]struct{}, len(cpus))
	for _, cpu := range cpus {
		if _, exists := seen[cpu]; exists {
			return fmt.Errorf("%w: %s CPU set contains duplicate CPU %d", errInvalidConfig, name, cpu)
		}
		seen[cpu] = struct{}{}
	}
	return nil
}

func overlapsCPUSet(left []int, right []int) bool {
	if len(left) == 0 || len(right) == 0 {
		return false
	}
	seen := make(map[int]struct{}, len(left))
	for _, cpu := range left {
		seen[cpu] = struct{}{}
	}
	for _, cpu := range right {
		if _, exists := seen[cpu]; exists {
			return true
		}
	}
	return false
}
