package cpu

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Info holds CPU data. Exported fields (capital letter) appear in JSON output.
type Info struct {
	UsagePercent float64
	Cores        int
	Model        string
}

// Get is the only function other packages call.
// It builds and returns an Info struct.
func Get() Info {
	return Info{
		UsagePercent: getUsage(),
		Cores:        runtime.NumCPU(), // Go's runtime knows the CPU count directly
		Model:        getModel(),
	}
}

// getUsage reads /proc/stat twice with a gap and calculates CPU usage %.
// /proc/stat is a virtual file — the Linux kernel writes live data into it.
// It doesn't sit on disk. Reading it is essentially a kernel call.
func getUsage() float64 {
	// readStat is a helper function defined inline using a variable.
	// This is called a function literal or closure in Go.
	readStat := func() (idle, total uint64) {
		out, err := exec.Command("cat", "/proc/stat").Output()
		if err != nil {
			return 0, 0
		}

		// First line: "cpu  264Reid 0 119951 3435168 ..."
		// Fields after "cpu": user, nice, system, idle, iowait, irq, softirq, steal
		line := strings.Split(string(out), "\n")[0]
		fields := strings.Fields(line)[1:] // skip the "cpu" label at index 0

		var vals []uint64
		for _, f := range fields {
			n, _ := strconv.ParseUint(f, 10, 64)
			vals = append(vals, n)
		}

		if len(vals) < 4 {
			return 0, 0
		}

		// idle is the 4th field (index 3)
		idle = vals[3]
		for _, v := range vals {
			total += v // total = sum of all time fields
		}
		return // named returns — Go returns idle and total automatically
	}

	idle1, total1 := readStat()
	time.Sleep(250 * time.Millisecond) // sample over a short window
	idle2, total2 := readStat()

	totalDelta := total2 - total1
	idleDelta := idle2 - idle1

	if totalDelta == 0 {
		return 0.0
	}

	// Usage = time spent NOT idle / total time
	return float64(totalDelta-idleDelta) / float64(totalDelta) * 100
}

// getModel reads the CPU model name from /proc/cpuinfo.
func getModel() string {
	// grep -m1 stops after the first match — faster than reading the whole file
	out, err := exec.Command("grep", "-m1", "model name", "/proc/cpuinfo").Output()
	if err != nil {
		return "unknown"
	}

	// Line looks like: "model name	: Intel(R) Core(TM) i7-..."
	parts := strings.SplitN(string(out), ":", 2)
	if len(parts) < 2 {
		return "unknown"
	}
	return strings.TrimSpace(parts[1])
}