package userpaths

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Source identifies how a user-level path was resolved.
type Source string

const (
	SourceOverride       Source = "override"
	SourceXDGExisting    Source = "xdg-existing"
	SourceXDGDefault     Source = "xdg-default"
	SourceLegacyFallback Source = "legacy-fallback"
)

// ResolvedPath captures the chosen path alongside candidate roots so callers
// can surface migration hints without reimplementing precedence logic.
type ResolvedPath struct {
	Path       string
	XDGPath    string
	LegacyPath string
	Source     Source
	HasXDG     bool
	HasLegacy  bool
}

// BeadsDataRoot returns the user-level beads data root.
func BeadsDataRoot() (string, error) {
	dataHome, err := userDataHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataHome, "beads"), nil
}

// BeadsStateRoot returns the user-level beads state root.
func BeadsStateRoot() (string, error) {
	stateHome, err := userStateHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(stateHome, "beads"), nil
}

// UserFormulasDir resolves the user-level formulas directory.
func UserFormulasDir() (ResolvedPath, error) {
	dataRoot, err := BeadsDataRoot()
	if err != nil {
		return ResolvedPath{}, err
	}
	legacyRoot, err := legacyBeadsRoot()
	if err != nil {
		return ResolvedPath{}, err
	}
	return resolve("", filepath.Join(dataRoot, "formulas"), filepath.Join(legacyRoot, "formulas")), nil
}

// UserMoleculesPath resolves the user-level molecules catalog path.
func UserMoleculesPath() (ResolvedPath, error) {
	dataRoot, err := BeadsDataRoot()
	if err != nil {
		return ResolvedPath{}, err
	}
	legacyRoot, err := legacyBeadsRoot()
	if err != nil {
		return ResolvedPath{}, err
	}
	return resolve("", filepath.Join(dataRoot, "molecules.jsonl"), filepath.Join(legacyRoot, "molecules.jsonl")), nil
}

// SharedServerStateDir resolves the shared-server runtime-state directory.
func SharedServerStateDir() (ResolvedPath, error) {
	override := strings.TrimSpace(os.Getenv("BEADS_SHARED_SERVER_DIR"))
	stateRoot, err := BeadsStateRoot()
	if err != nil {
		return ResolvedPath{}, err
	}
	legacyRoot, err := legacyBeadsRoot()
	if err != nil {
		return ResolvedPath{}, err
	}
	return resolve(override, filepath.Join(stateRoot, "shared-server"), filepath.Join(legacyRoot, "shared-server")), nil
}

// SharedServerDoltDir resolves the shared-server durable Dolt data directory.
func SharedServerDoltDir() (ResolvedPath, error) {
	override := strings.TrimSpace(os.Getenv("BEADS_SHARED_SERVER_DIR"))
	if override != "" {
		override = filepath.Join(override, "dolt")
	}
	dataRoot, err := BeadsDataRoot()
	if err != nil {
		return ResolvedPath{}, err
	}
	legacyRoot, err := legacyBeadsRoot()
	if err != nil {
		return ResolvedPath{}, err
	}
	return resolve(override, filepath.Join(dataRoot, "shared-server", "dolt"), filepath.Join(legacyRoot, "shared-server", "dolt")), nil
}

// ContributorPlanningRepo resolves the default contributor planning repository.
func ContributorPlanningRepo() (ResolvedPath, error) {
	dataRoot, err := BeadsDataRoot()
	if err != nil {
		return ResolvedPath{}, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ResolvedPath{}, fmt.Errorf("cannot determine home directory: %w", err)
	}
	return resolve("", filepath.Join(dataRoot, "planning"), filepath.Join(home, ".beads-planning")), nil
}

func resolve(override, xdgPath, legacyPath string) ResolvedPath {
	if override != "" {
		return ResolvedPath{
			Path:       override,
			XDGPath:    xdgPath,
			LegacyPath: legacyPath,
			Source:     SourceOverride,
		}
	}

	hasXDG := pathExists(xdgPath)
	hasLegacy := pathExists(legacyPath)

	switch {
	case hasLegacy && !hasXDG:
		return ResolvedPath{
			Path:       legacyPath,
			XDGPath:    xdgPath,
			LegacyPath: legacyPath,
			Source:     SourceLegacyFallback,
			HasXDG:     hasXDG,
			HasLegacy:  hasLegacy,
		}
	case hasXDG:
		return ResolvedPath{
			Path:       xdgPath,
			XDGPath:    xdgPath,
			LegacyPath: legacyPath,
			Source:     SourceXDGExisting,
			HasXDG:     hasXDG,
			HasLegacy:  hasLegacy,
		}
	default:
		return ResolvedPath{
			Path:       xdgPath,
			XDGPath:    xdgPath,
			LegacyPath: legacyPath,
			Source:     SourceXDGDefault,
			HasXDG:     hasXDG,
			HasLegacy:  hasLegacy,
		}
	}
}

func userDataHome() (string, error) {
	if dir := strings.TrimSpace(os.Getenv("XDG_DATA_HOME")); dir != "" {
		return dir, nil
	}
	switch runtime.GOOS {
	case "windows":
		if dir := strings.TrimSpace(os.Getenv("LOCALAPPDATA")); dir != "" {
			return dir, nil
		}
		if dir := strings.TrimSpace(os.Getenv("APPDATA")); dir != "" {
			return dir, nil
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".local", "share"), nil
}

func userStateHome() (string, error) {
	if dir := strings.TrimSpace(os.Getenv("XDG_STATE_HOME")); dir != "" {
		return dir, nil
	}
	switch runtime.GOOS {
	case "windows":
		if dir := strings.TrimSpace(os.Getenv("LOCALAPPDATA")); dir != "" {
			return dir, nil
		}
		if dir := strings.TrimSpace(os.Getenv("APPDATA")); dir != "" {
			return dir, nil
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".local", "state"), nil
}

func legacyBeadsRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".beads"), nil
}

func pathExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}
