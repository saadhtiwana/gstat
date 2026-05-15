package render

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/saadhtiwana/gstat/internal/cpu"
	"github.com/saadhtiwana/gstat/internal/disk"
	"github.com/saadhtiwana/gstat/internal/mem"
	"github.com/saadhtiwana/gstat/internal/netstat"
	"github.com/saadhtiwana/gstat/internal/process"
)

// Stats is the single container for all collected data.
// render is the only package that imports all others — clean dependency direction.
type Stats struct {
	CPU       cpu.Info
	Memory    mem.Info
	Disk      []disk.Info
	Network   []netstat.Info
	Processes []process.Info
}

// ANSI escape codes for terminal colors
const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	cyan   = "\033[36m"
)

// Print renders all stats to stdout with color formatting.
func Print(s Stats) {
	fmt.Printf("%sgstat%s — %s\n\n", bold, reset, time.Now().Format("Mon 02 Jan 2006 15:04:05"))
	printCPU(s.CPU)
	printMemory(s.Memory)
	printDisk(s.Disk)
	printNetwork(s.Network)
	printProcesses(s.Processes)
}

func printCPU(c cpu.Info) {
	fmt.Printf("%sCPU%s\n", bold, reset)
	fmt.Printf("  %-20s %s%.2f%%%s\n", "Usage", colorFor(c.UsagePercent), c.UsagePercent, reset)
	fmt.Printf("  %-20s %d\n", "Cores", c.Cores)
	fmt.Printf("  %-20s %s\n\n", "Model", c.Model)
}

func printMemory(m mem.Info) {
	fmt.Printf("%sMemory%s\n", bold, reset)
	fmt.Printf("  %-20s %s%.2f%%%s\n", "Used", colorFor(m.UsedPercent), m.UsedPercent, reset)
	fmt.Printf("  %-20s %s\n", "Total", formatBytes(m.Total))
	fmt.Printf("  %-20s %s\n", "Used", formatBytes(m.Used))
	fmt.Printf("  %-20s %s\n\n", "Available", formatBytes(m.Available))
}

func printDisk(disks []disk.Info) {
	fmt.Printf("%sDisk%s\n", bold, reset)
	fmt.Printf("  %-18s %-10s %-10s %-10s %s\n", "Mount", "Total", "Used", "Free", "Usage")
	fmt.Println("  " + line(58))
	for _, d := range disks {
		fmt.Printf("  %-18s %-10s %-10s %-10s %s%.1f%%%s\n",
			d.Mountpoint,
			formatBytes(d.Total),
			formatBytes(d.Used),
			formatBytes(d.Free),
			colorFor(d.UsedPercent), d.UsedPercent, reset,
		)
	}
	fmt.Println()
}

func printNetwork(ifaces []netstat.Info) {
	fmt.Printf("%sNetwork%s\n", bold, reset)
	fmt.Printf("  %-15s %-15s %-15s\n", "Interface", "Received", "Sent")
	fmt.Println("  " + line(48))
	for _, n := range ifaces {
		fmt.Printf("  %-15s %-15s %-15s\n",
			n.Name,
			formatBytes(n.BytesRecv),
			formatBytes(n.BytesSent),
		)
	}
	fmt.Println()
}

func printProcesses(procs []process.Info) {
	fmt.Printf("%sTop Processes%s\n", bold, reset)
	fmt.Printf("  %-8s %-22s %-12s %-10s\n", "PID", "Name", "CPU Ticks", "Mem%")
	fmt.Println("  " + line(55))
	for _, p := range procs {
		fmt.Printf("  %-8d %-22s %-12.0f %-10.2f\n",
			p.PID, p.Name, p.CPUPercent, p.MemPercent,
		)
	}
	fmt.Println()
}

// colorFor returns a terminal color based on usage percentage.
// green < 50%, yellow < 80%, red >= 80%
func colorFor(pct float64) string {
	if pct < 50 {
		return green
	} else if pct < 80 {
		return yellow
	}
	return red
}

// formatBytes converts raw bytes to a human-readable string.
func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGT"[exp])
}

// line returns a string of dashes of length n — used for table separators.
func line(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "-"
	}
	return s
}

// ExportJSON serializes Stats to a JSON file.
func ExportJSON(s Stats) error {
	// MarshalIndent = pretty-printed JSON with 2-space indent
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		// %w wraps the error — caller can use errors.Is() to inspect it
		return fmt.Errorf("json marshal failed: %w", err)
	}
	// 0644 = rw-r--r-- file permissions
	return os.WriteFile("gstat_export.json", data, 0644)
}