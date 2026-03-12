package utils

import (
	"strings"
	"unicode"
)

// issueTypeAliases maps shorthand type names to canonical types
var issueTypeAliases = map[string]string{
	"mr":          "merge-request",
	"feat":        "feature",
	"mol":         "molecule",
	"enhancement": "feature",
	"dec":         "decision",
	"adr":         "decision",
}

// NormalizeIssueType expands type aliases to their canonical forms.
// For example: "mr" -> "merge-request", "feat" -> "feature", "mol" -> "molecule"
// Returns the input unchanged if it's not an alias.
func NormalizeIssueType(t string) string {
	if canonical, ok := issueTypeAliases[strings.ToLower(t)]; ok {
		return canonical
	}
	return t
}

// NormalizeLabels trims whitespace, removes empty strings, and deduplicates labels
// while preserving order.
func NormalizeLabels(ss []string) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// NormalizeDetectedIssuePrefix converts an auto-detected directory name into a
// safe, canonical issue prefix. Explicit --prefix values are intentionally not
// rewritten by this helper.
func NormalizeDetectedIssuePrefix(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))

	var b strings.Builder
	b.Grow(len(name))
	lastWasDash := false
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastWasDash = false
		default:
			if !lastWasDash {
				b.WriteByte('-')
				lastWasDash = true
			}
		}
	}

	prefix := strings.Trim(b.String(), "-")
	if prefix == "" {
		return "bd"
	}
	if unicode.IsDigit(rune(prefix[0])) {
		return "bd-" + prefix
	}
	return prefix
}

// DatabaseNameFromPrefix derives the SQL database name from an issue prefix.
func DatabaseNameFromPrefix(prefix string) string {
	return strings.ReplaceAll(prefix, "-", "_")
}
