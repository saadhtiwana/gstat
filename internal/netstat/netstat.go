package netstat

import (
	"os/exec"
	"strconv"
	"strings"
)

type Info struct {
	Name      string
	BytesSent uint64
	BytesRecv uint64
}

func Get() []Info {
	// /proc/net/dev has per-interface stats since boot
	out, err := exec.Command("cat", "/proc/net/dev").Output()
	if err != nil {
		return nil
	}

	var result []Info
	lines := strings.Split(string(out), "\n")

	// First two lines are headers — skip them
	for _, line := range lines[2:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Format: "eth0:  1234 ..."
		// Split on ":" to separate interface name from data
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		if name == "lo" {
			continue // loopback is not useful
		}

		fields := strings.Fields(parts[1])
		if len(fields) < 9 {
			continue
		}

		// Column 0 = bytes received, column 8 = bytes sent
		recv, _ := strconv.ParseUint(fields[0], 10, 64)
		sent, _ := strconv.ParseUint(fields[8], 10, 64)

		result = append(result, Info{
			Name:      name,
			BytesRecv: recv,
			BytesSent: sent,
		})
	}
	return result
}
