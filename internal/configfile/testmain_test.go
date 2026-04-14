package configfile

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Keep configfile tests deterministic even when the operator shell is
	// running in shared-server mode.
	for _, key := range []string{
		"BEADS_DOLT_SHARED_SERVER",
		"BEADS_DOLT_SERVER_MODE",
		"BEADS_SHARED_SERVER_DIR",
	} {
		_ = os.Unsetenv(key)
	}

	os.Exit(m.Run())
}
