package render

import (
	"testing"

	"github.com/saadhtiwana/gstat/internal/cpu"
	"github.com/saadhtiwana/gstat/internal/disk"
	"github.com/saadhtiwana/gstat/internal/mem"
	"github.com/saadhtiwana/gstat/internal/netstat"
	"github.com/saadhtiwana/gstat/internal/process"
)

func TestColorFor(t *testing.T) {
	// Table-driven tests — idiomatic Go.
	// Instead of writing 3 separate test functions, we define cases as a slice
	// and loop over them. Cleaner and easier to extend.
	tests := []struct {
		input    float64
		expected string
	}{
		{30.0, green},
		{65.0, yellow},
		{90.0, red},
	}

	for _, tt := range tests {
		got := colorFor(tt.input)
		if got != tt.expected {
			t.Errorf("colorFor(%.1f) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		got := formatBytes(tt.input)
		if got != tt.expected {
			t.Errorf("formatBytes(%d) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestExportJSON(t *testing.T) {
	// Build a minimal Stats with zero values and make sure export doesn't error
	s := Stats{
		CPU:       cpu.Info{UsagePercent: 10.0, Cores: 4, Model: "Test CPU"},
		Memory:    mem.Info{Total: 8000, Used: 4000, Available: 4000, UsedPercent: 50},
		Disk:      []disk.Info{},
		Network:   []netstat.Info{},
		Processes: []process.Info{},
	}

	err := ExportJSON(s)
	if err != nil {
		t.Errorf("ExportJSON returned error: %v", err)
	}
}
