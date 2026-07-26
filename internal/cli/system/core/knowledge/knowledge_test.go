//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package knowledge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFile writes content under dir/name, creating parent dirs.
func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// unmigratedLearnings builds an un-migrated LEARNINGS root with n staged
// entries (no ## Themes) — the whole entry region is staging.
func unmigratedLearnings(n int) string {
	var b strings.Builder
	b.WriteString("# Learnings\n\n<!-- guide -->\n\n")
	for i := 0; i < n; i++ {
		b.WriteString("## [2026-07-15-1200")
		b.WriteByte(byte('0' + i%10))
		b.WriteString("0] entry\n\nbody.\n\n")
	}
	return b.String()
}

// M1: an un-migrated root over threshold yields a foldable finding whose
// count is the full staged-entry count.
func TestHealth_FoldableUnmigrated(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "LEARNINGS.md", unmigratedLearnings(40))

	got := Health(dir, Thresholds{Learnings: 30})
	if len(got) != 1 {
		t.Fatalf("findings = %d, want 1: %+v", len(got), got)
	}
	if got[0].Kind != foldable || got[0].Count != 40 {
		t.Errorf("finding = %+v, want foldable count 40", got[0])
	}
}

// M2: a folded root (few staged, bulk in themes) stays under threshold
// and does not foldable-nudge.
func TestHealth_FoldedQuiet(t *testing.T) {
	dir := t.TempDir()
	// 3 staged entries above ## Themes; theme gists are bullets, not "## [".
	root := "# Learnings\n\n" +
		"## [2026-07-15-120000] a\n\nx.\n\n" +
		"## [2026-07-15-120001] b\n\nx.\n\n" +
		"## [2026-07-15-120002] c\n\nx.\n\n" +
		"## Themes\n\n- hooks — gist → [hooks](learnings/hooks.md)\n"
	writeFile(t, dir, "LEARNINGS.md", root)

	got := Health(dir, Thresholds{Learnings: 30})
	for _, f := range got {
		if f.Kind == foldable {
			t.Errorf("unexpected foldable finding on a folded root: %+v", f)
		}
	}
}

// M3: a theme file over the byte ceiling yields a heavy finding on that
// file, even when the root is lean.
func TestHealth_HeavyThemeFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "CONVENTIONS.md",
		"# Conventions\n\n## Themes\n\n- x — g → [x](conventions/x.md)\n")
	writeFile(t, dir, "conventions/x.md",
		"# x\n\n"+strings.Repeat("padding line\n", 8000)) // ~100 KB

	got := Health(dir, Thresholds{Conventions: 12, PageBytes: 65536})
	var heavyFile string
	for _, f := range got {
		if f.Kind == heavy && strings.Contains(f.File, "conventions/x.md") {
			heavyFile = f.File
		}
	}
	if heavyFile == "" {
		t.Errorf("no heavy finding on the fat theme file: %+v", got)
	}
}

// M4: a large un-migrated root trips both signals; foldable is ordered
// before heavy (folding reduces both).
func TestHealth_BothFire_FoldableFirst(t *testing.T) {
	dir := t.TempDir()
	// 40 entries and > 64 KB of body.
	big := "# Learnings\n\n"
	for i := 0; i < 40; i++ {
		big += "## [2026-07-15-12000" + string(rune('0'+i%10)) +
			"] e\n\n" + strings.Repeat("x", 2000) + "\n\n"
	}
	writeFile(t, dir, "LEARNINGS.md", big)

	got := Health(dir, Thresholds{Learnings: 30, PageBytes: 65536})
	if len(got) < 2 {
		t.Fatalf("want both signals, got %+v", got)
	}
	if got[0].Kind != foldable {
		t.Errorf("first finding = %+v, want foldable to lead", got[0])
	}
	sawHeavy := false
	for _, f := range got {
		if f.Kind == heavy {
			sawHeavy = true
		}
	}
	if !sawHeavy {
		t.Error("no heavy finding on an oversized root")
	}
}

// M5 (measure): conventions are counted by section, not lines — a short-
// section-count but long root does not foldable-nudge.
func TestHealth_ConventionMeasureIsSections(t *testing.T) {
	dir := t.TempDir()
	// 5 sections, each padded long (many lines) but only 5 sections.
	root := "# Conventions\n\n"
	for i := 0; i < 5; i++ {
		root += "## Section " + string(rune('A'+i)) + "\n\n" +
			strings.Repeat("- a rule line\n", 20) + "\n"
	}
	writeFile(t, dir, "CONVENTIONS.md", root)

	got := Health(dir, Thresholds{Conventions: 12}) // byte check disabled
	for _, f := range got {
		if f.Kind == foldable {
			t.Errorf("5 sections should not foldable-nudge at 12: %+v", f)
		}
	}
}

// M6/M7: threshold fires at > N (not >=), and 0 disables the check.
func TestHealth_BoundaryAndDisable(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "LEARNINGS.md", unmigratedLearnings(30))

	// Exactly at threshold → no finding.
	if got := Health(dir, Thresholds{Learnings: 30}); len(got) != 0 {
		t.Errorf("count == threshold should not fire: %+v", got)
	}
	// One over → fires.
	writeFile(t, dir, "LEARNINGS.md", unmigratedLearnings(31))
	if got := Health(dir, Thresholds{Learnings: 30}); len(got) != 1 {
		t.Errorf("count > threshold should fire: %+v", got)
	}
	// Disabled (0) → never fires.
	if got := Health(dir, Thresholds{Learnings: 0}); len(got) != 0 {
		t.Errorf("threshold 0 should disable: %+v", got)
	}
}
