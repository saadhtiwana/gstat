package mem

import (
	"os/exec"
	"strconv"
	"strings"
)

type Info struct {
	Total       uint64
	Available   uint64
	Used        uint64
	UsedPercent float64
}

func Get() Info {
	// /proc/meminfo has all memory fields, values in kB
	out, err := exec.Command("cat", "/proc/meminfo").Output()
	if err != nil {
		return Info{}
	}

	// Build a map from field name → value in bytes
	// e.g. fields["MemTotal"] = 16777216000
	fields := make(map[string]uint64)

	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSuffix(parts[0], ":") // "MemTotal:" → "MemTotal"
		val, _ := strconv.ParseUint(parts[1], 10, 64)
		fields[key] = val * 1024 // kB → bytes
	}

	total := fields["MemTotal"]
	available := fields["MemAvailable"]
	used := total - available

	var usedPct float64
	if total > 0 {
		usedPct = float64(used) / float64(total) * 100
	}

	return Info{
		Total:       total,
		Available:   available,
		Used:        used,
		UsedPercent: usedPct,
	}
}