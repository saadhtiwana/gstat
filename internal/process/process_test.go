package process

import "testing"

func TestGet(t *testing.T) {
	procs := Get()

	if len(procs) == 0 {
		t.Error("expected at least one process")
	}

	// Should never return more than 10
	if len(procs) > 10 {
		t.Errorf("expected max 10 processes, got %d", len(procs))
	}

	for _, p := range procs {
		if p.PID <= 0 {
			t.Errorf("invalid PID: %d", p.PID)
		}
		if p.MemPercent < 0 || p.MemPercent > 100 {
			t.Errorf("MemPercent out of range for %s: %f", p.Name, p.MemPercent)
		}
	}
}

func TestGetTotalMem(t *testing.T) {
	mem := getTotalMem()
	if mem == 0 {
		t.Error("getTotalMem() returned 0")
	}
}
