//go:build cgo

package fix

import (
	"fmt"
	"os"
	"testing"

	"github.com/steveyegge/beads/internal/testutil"
)

// TestMain starts an isolated Dolt server so fix tests don't hit the
// production server on port 3307.
func TestMain(m *testing.M) {
	for _, key := range []string{
		"BEADS_DOLT_SHARED_SERVER",
		"BEADS_DOLT_SERVER_MODE",
		"BEADS_SHARED_SERVER_DIR",
		"BEADS_DOLT_DATA_DIR",
	} {
		_ = os.Unsetenv(key)
	}

	os.Setenv("BEADS_TEST_MODE", "1")
	if err := testutil.EnsureDoltContainerForTestMain(); err != nil {
		fmt.Fprintf(os.Stderr, "WARN: %v, skipping Dolt tests\n", err)
	} else {
		defer testutil.TerminateDoltContainer()
	}

	code := m.Run()

	os.Unsetenv("BEADS_DOLT_PORT")
	os.Unsetenv("BEADS_TEST_MODE")
	os.Exit(code)
}
