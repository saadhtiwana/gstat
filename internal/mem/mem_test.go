package mem

import "testing"

func TestGet(t *testing.T) {
	info := Get()

	if info.Total == 0 {
		t.Error("Total memory should not be 0")
	}

	if info.Used > info.Total {
		t.Errorf("Used (%d) cannot exceed Total (%d)", info.Used, info.Total)
	}

	if info.UsedPercent < 0 || info.UsedPercent > 100 {
		t.Errorf("UsedPercent out of range: %f", info.UsedPercent)
	}

	// Available + Used should roughly equal Total
	// We allow a small delta because kernel uses some memory internally
	delta := int64(info.Total) - int64(info.Available) - int64(info.Used)
	if delta < 0 {
		delta = -delta
	}
	if delta > 1024*1024*100 { // more than 100MB discrepancy is suspicious
		t.Errorf("memory values don't add up: total=%d avail=%d used=%d", info.Total, info.Available, info.Used)
	}
}
