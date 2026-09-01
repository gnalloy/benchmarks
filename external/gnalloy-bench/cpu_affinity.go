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
