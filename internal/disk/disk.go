package disk

import (
	"os/exec"
	"strconv"
	"strings"
)

type Info struct {
	Mountpoint  string
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func Get() []Info {
	// df -B1 = disk usage in bytes
	// --output=target,size,used,avail,pcent = only the columns we need
	out, err := exec.Command("df", "-B1", "--output=target,size,used,avail,pcent").Output()
	if err != nil {
		return nil
	}

	var result []Info
	lines := strings.Split(string(out), "\n")

	// lines[0] is the header row — skip it
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		mount := fields[0]

		// These are virtual filesystems — not real disk partitions
		if strings.HasPrefix(mount, "/proc") ||
			strings.HasPrefix(mount, "/sys") ||
			strings.HasPrefix(mount, "/dev/pts") ||
			mount == "tmpfs" || mount == "devtmpfs" {
			continue
		}

		total, _ := strconv.ParseUint(fields[1], 10, 64)
		used, _ := strconv.ParseUint(fields[2], 10, 64)
		free, _ := strconv.ParseUint(fields[3], 10, 64)

		// pcent looks like "42%" — strip the percent sign
		pctStr := strings.TrimSuffix(fields[4], "%")
		pct, _ := strconv.ParseFloat(pctStr, 64)

		result = append(result, Info{
			Mountpoint:  mount,
			Total:       total,
			Used:        used,
			Free:        free,
			UsedPercent: pct,
		})
	}
	return result
}
