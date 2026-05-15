package cpu

// Tests live in the same package so they can test unexported functions too.
// File must end in _test.go — Go only compiles these when running `go test`.

import (
	"testing"
)

// TestGet checks that Get() returns something sensible.
// Every test function must start with Test and take *testing.T.
func TestGet(t *testing.T) {
	info := Get()

	// t.Errorf marks the test as failed but keeps running.
	// t.Fatalf would stop immediately.
	if info.Cores <= 0 {
		t.Errorf("expected Cores > 0, got %d", info.Cores)
	}

	if info.UsagePercent < 0 || info.UsagePercent > 100 {
		t.Errorf("UsagePercent out of range: %f", info.UsagePercent)
	}

	if info.Model == "" {
		t.Errorf("expected non-empty Model")
	}
}

// TestGetUsage checks the usage calculation stays in bounds.
func TestGetUsage(t *testing.T) {
	usage := getUsage()
	if usage < 0 || usage > 100 {
		t.Errorf("getUsage() returned %f, want 0-100", usage)
	}
}

// TestGetModel checks we get a non-empty model string.
func TestGetModel(t *testing.T) {
	model := getModel()
	if model == "" {
		t.Error("getModel() returned empty string")
	}
}
