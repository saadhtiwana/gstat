package process

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Info struct {
	PID        int
	Name       string
	CPUPercent float64
	MemPercent float64
}

// Get reads /proc and returns the top 10 processes by CPU ticks.
func Get() []Info {
	totalMem := getTotalMem()

	// os.ReadDir reads a directory and returns its entries sorted by name
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}

	var result []Info

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Only directories with numeric names are processes
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		name := getProcessName(pid)
		memKB := getProcessMem(pid)
		cpuTicks := getProcessCPU(pid)

		var memPct float64
		if totalMem > 0 {
			// memKB is in kilobytes, totalMem is in bytes — align units first
			memPct = float64(memKB*1024) / float64(totalMem) * 100
		}

		result = append(result, Info{
			PID:        pid,
			Name:       name,
			CPUPercent: float64(cpuTicks), // raw ticks, relative comparison only
			MemPercent: memPct,
		})
	}

	// Sort descending by CPUPercent using a simple bubble sort
	// No imports needed — good for a beginner exercise
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].CPUPercent > result[i].CPUPercent {
				result[i], result[j] = result[j], result[i] // swap
			}
		}
	}

	if len(result) > 10 {
		return result[:10]
	}
	return result
}

func getProcessName(pid int) string {
	// /proc/<pid>/comm = just the process name, one line
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "comm"))
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func getProcessMem(pid int) uint64 {
	// VmRSS in /proc/<pid>/status = physical RAM currently used by this process
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "status"))
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "VmRSS:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				val, _ := strconv.ParseUint(fields[1], 10, 64)
				return val // kB
			}
		}
	}
	return 0
}

func getProcessCPU(pid int) uint64 {
	// /proc/<pid>/stat field 14 = utime (user CPU ticks), field 15 = stime (kernel CPU ticks)
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return 0
	}
	fields := strings.Fields(string(data))
	if len(fields) < 16 {
		return 0
	}
	utime, _ := strconv.ParseUint(fields[13], 10, 64)
	stime, _ := strconv.ParseUint(fields[14], 10, 64)
	return utime + stime
}

func getTotalMem() uint64 {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				val, _ := strconv.ParseUint(fields[1], 10, 64)
				return val * 1024 // kB → bytes
			}
		}
	}
	return 0
}
