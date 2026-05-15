package main

// Every Go file starts with package declaration.
// main is special — it's the entry point of the program.

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/saadhtiwana/gstat/internal/cpu"
	"github.com/saadhtiwana/gstat/internal/disk"
	"github.com/saadhtiwana/gstat/internal/mem"
	"github.com/saadhtiwana/gstat/internal/netstat"
	"github.com/saadhtiwana/gstat/internal/process"
	"github.com/saadhtiwana/gstat/internal/render"
)

func main() {
	// flag.Int defines a flag called "interval".
	// Returns a pointer *int, not an int directly.
	// Default is 3. The string is shown when user runs --help.
	interval := flag.Int("interval", 3, "seconds between refreshes in monitor mode")

	// This actually reads os.Args and fills in all flags you defined above.
	flag.Parse()

	// flag.Args() returns everything that isn't a flag.
	// So `gstat --interval 2 monitor` → args = ["monitor"]
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: gstat <snapshot|monitor|export>")
		fmt.Println("       gstat --interval 2 monitor")
		os.Exit(1)
	}

	// Switch on the first argument — our subcommand
	switch args[0] {
	case "snapshot":
		runSnapshot()
	case "monitor":
		// *interval dereferences the pointer — gives us the actual int value
		runMonitor(*interval)
	case "export":
		runExport()
	default:
		fmt.Printf("unknown command: %s\n", args[0])
		os.Exit(1)
	}
}

// collectAll calls every package's Get() and bundles the results.
// All three subcommands use this — one place to change if we add a new stat.
func collectAll() render.Stats {
	return render.Stats{
		CPU:       cpu.Get(),
		Memory:    mem.Get(),
		Disk:      disk.Get(),
		Network:   netstat.Get(),
		Processes: process.Get(),
	}
}

// runSnapshot collects stats once and prints them.
func runSnapshot() {
	stats := collectAll()
	render.Print(stats)
}

// runMonitor loops forever, clearing the screen and reprinting every N seconds.
func runMonitor(interval int) {
	// time.NewTicker fires on its channel every N seconds.
	// Think of it as a recurring alarm clock.
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	// defer runs when the surrounding function returns.
	// Stop() cleans up the ticker goroutine — always do this.
	defer ticker.Stop()

	for {
		// ANSI escape: clear screen and move cursor to top-left.
		// \033[2J = clear, \033[H = move cursor home
		fmt.Print("\033[2J\033[H")

		stats := collectAll()
		render.Print(stats)

		// Block here until the ticker sends — then loop again.
		<-ticker.C
	}
}

// runExport collects stats and writes them to a JSON file.
func runExport() {
	stats := collectAll()
	err := render.ExportJSON(stats)
	if err != nil {
		// Stderr for errors, not stdout
		fmt.Fprintf(os.Stderr, "export error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("written to gstat_export.json")
}