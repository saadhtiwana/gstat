package disk

import "testing"

func TestGet(t *testing.T) {
	disks := Get()

	// On any Linux machine there should be at least one real mount
	if len(disks) == 0 {
		t.Error("expected at least one disk entry")
	}

	for _, d := range disks {
		if d.Mountpoint == "" {
			t.Error("disk has empty mountpoint")
		}

		if d.Used > d.Total {
			t.Errorf("disk %s: used (%d) > total (%d)", d.Mountpoint, d.Used, d.Total)
		}

		if d.UsedPercent < 0 || d.UsedPercent > 100 {
			t.Errorf("disk %s: UsedPercent out of range: %f", d.Mountpoint, d.UsedPercent)
		}
	}
}
