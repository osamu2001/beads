package userpaths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUserFormulasDir_DefaultsToXDGPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_DATA_HOME", filepath.Join(home, "xdg-data"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(home, "xdg-state"))
	t.Setenv("BEADS_SHARED_SERVER_DIR", "")

	resolved, err := UserFormulasDir()
	if err != nil {
		t.Fatalf("UserFormulasDir: %v", err)
	}

	want := filepath.Join(home, "xdg-data", "beads", "formulas")
	if resolved.Path != want {
		t.Fatalf("Path = %q, want %q", resolved.Path, want)
	}
	if resolved.Source != SourceXDGDefault {
		t.Fatalf("Source = %q, want %q", resolved.Source, SourceXDGDefault)
	}
}

func TestUserFormulasDir_FallsBackToLegacyWhenOnlyLegacyExists(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_DATA_HOME", filepath.Join(home, "xdg-data"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(home, "xdg-state"))

	legacy := filepath.Join(home, ".beads", "formulas")
	if err := os.MkdirAll(legacy, 0o755); err != nil {
		t.Fatalf("MkdirAll legacy: %v", err)
	}

	resolved, err := UserFormulasDir()
	if err != nil {
		t.Fatalf("UserFormulasDir: %v", err)
	}

	if resolved.Path != legacy {
		t.Fatalf("Path = %q, want %q", resolved.Path, legacy)
	}
	if resolved.Source != SourceLegacyFallback {
		t.Fatalf("Source = %q, want %q", resolved.Source, SourceLegacyFallback)
	}
}

func TestUserMoleculesPath_PrefersXDGWhenBothExist(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_DATA_HOME", filepath.Join(home, "xdg-data"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(home, "xdg-state"))

	xdg := filepath.Join(home, "xdg-data", "beads", "molecules.jsonl")
	if err := os.MkdirAll(filepath.Dir(xdg), 0o755); err != nil {
		t.Fatalf("MkdirAll xdg: %v", err)
	}
	if err := os.WriteFile(xdg, []byte(""), 0o644); err != nil {
		t.Fatalf("WriteFile xdg: %v", err)
	}

	legacy := filepath.Join(home, ".beads", "molecules.jsonl")
	if err := os.MkdirAll(filepath.Dir(legacy), 0o755); err != nil {
		t.Fatalf("MkdirAll legacy: %v", err)
	}
	if err := os.WriteFile(legacy, []byte(""), 0o644); err != nil {
		t.Fatalf("WriteFile legacy: %v", err)
	}

	resolved, err := UserMoleculesPath()
	if err != nil {
		t.Fatalf("UserMoleculesPath: %v", err)
	}

	if resolved.Path != xdg {
		t.Fatalf("Path = %q, want %q", resolved.Path, xdg)
	}
	if resolved.Source != SourceXDGExisting {
		t.Fatalf("Source = %q, want %q", resolved.Source, SourceXDGExisting)
	}
	if !resolved.HasLegacy || !resolved.HasXDG {
		t.Fatalf("expected both XDG and legacy markers, got %+v", resolved)
	}
}

func TestSharedServerPaths_UseCombinedOverride(t *testing.T) {
	home := t.TempDir()
	override := filepath.Join(home, "custom-shared-server")
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_DATA_HOME", filepath.Join(home, "xdg-data"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(home, "xdg-state"))
	t.Setenv("BEADS_SHARED_SERVER_DIR", override)

	stateDir, err := SharedServerStateDir()
	if err != nil {
		t.Fatalf("SharedServerStateDir: %v", err)
	}
	dataDir, err := SharedServerDoltDir()
	if err != nil {
		t.Fatalf("SharedServerDoltDir: %v", err)
	}

	if stateDir.Path != override {
		t.Fatalf("state Path = %q, want %q", stateDir.Path, override)
	}
	if dataDir.Path != filepath.Join(override, "dolt") {
		t.Fatalf("data Path = %q, want %q", dataDir.Path, filepath.Join(override, "dolt"))
	}
	if stateDir.Source != SourceOverride || dataDir.Source != SourceOverride {
		t.Fatalf("override sources = %q / %q, want %q", stateDir.Source, dataDir.Source, SourceOverride)
	}
}

func TestContributorPlanningRepo_DefaultsToXDGDataRoot(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_DATA_HOME", filepath.Join(home, "xdg-data"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(home, "xdg-state"))

	resolved, err := ContributorPlanningRepo()
	if err != nil {
		t.Fatalf("ContributorPlanningRepo: %v", err)
	}

	want := filepath.Join(home, "xdg-data", "beads", "planning")
	if resolved.Path != want {
		t.Fatalf("Path = %q, want %q", resolved.Path, want)
	}
	if resolved.Source != SourceXDGDefault {
		t.Fatalf("Source = %q, want %q", resolved.Source, SourceXDGDefault)
	}
}

func TestContributorPlanningRepo_FallsBackToLegacyWhenOnlyLegacyExists(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_DATA_HOME", filepath.Join(home, "xdg-data"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(home, "xdg-state"))

	legacy := filepath.Join(home, ".beads-planning")
	if err := os.MkdirAll(legacy, 0o755); err != nil {
		t.Fatalf("MkdirAll legacy: %v", err)
	}

	resolved, err := ContributorPlanningRepo()
	if err != nil {
		t.Fatalf("ContributorPlanningRepo: %v", err)
	}

	if resolved.Path != legacy {
		t.Fatalf("Path = %q, want %q", resolved.Path, legacy)
	}
	if resolved.Source != SourceLegacyFallback {
		t.Fatalf("Source = %q, want %q", resolved.Source, SourceLegacyFallback)
	}
}
