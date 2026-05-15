package netstat

import "testing"

func TestGet(t *testing.T) {
	ifaces := Get()

	// Most machines have at least one non-loopback interface
	if len(ifaces) == 0 {
		t.Log("no network interfaces found — this may be expected in some environments")
		return
	}

	for _, iface := range ifaces {
		if iface.Name == "" {
			t.Error("interface has empty name")
		}
		// lo should be filtered out
		if iface.Name == "lo" {
			t.Error("loopback interface should be filtered")
		}
	}
}
