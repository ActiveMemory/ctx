//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill_test

import (
	"bytes"
	"io/fs"
	"path"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// presenceSkillTrees lists generated trees with opt-in enrollment:
// a skill syncs iff its directory exists in the tree, and a
// directory with no Claude counterpart is tool-only and exempt
// (the Copilot tree carries native wrapper skills).
var presenceSkillTrees = []string{
	asset.DirIntegrationsCopilotSkill,
}

// mirrorSkillTrees lists generated trees that fully mirror the
// canonical Claude tree: every canonical skill outside
// claudeOnlySkills must be present, orphans are removed by the
// sync script, and references/ directories are copied along.
var mirrorSkillTrees = []string{
	asset.DirIntegrationsOpenCodeSkill,
	asset.DirCodexSkills,
}

// claudeOnlySkills mirrors the EXCLUDE lists in
// hack/sync-opencode-skills.sh and hack/sync-codex-skills.sh:
// skills whose body only makes sense inside Claude Code. A skill
// excluded in a script but not here fails the completeness check,
// so the lists cannot silently drift apart in that direction.
var claudeOnlySkills = map[string]bool{
	"ctx-permission-sanitize": true, // audits .claude/settings.local.json
	"ctx-plan-import":         true, // imports ~/.claude/plans/
	"ctx-dream":               true, // headless `claude -p` cron + guard.sh
	"ctx-skill-create":        true, // authors Claude Code skills/plugins
}

// TestSyncedSkillParity asserts the sync contract at the test layer,
// where CI enforces it (`make check-*-skills` only runs via `make
// audit` on developer machines): every skill in a generated tree that
// has a canonical Claude counterpart must be byte-identical to that
// counterpart minus `allowed-tools:` lines (the Claude Code-specific
// frontmatter key the sync scripts strip). Mirror trees are
// additionally held to completeness (every canonical skill outside
// claudeOnlySkills is present), no orphans, and reference parity
// (every embedded canonical references/ file is a byte-copy;
// non-.md references sit outside the canonical embed glob and are
// covered by the sync scripts, not this test).
func TestSyncedSkillParity(t *testing.T) {
	var checked, exempt, refsChecked int
	allTrees := make([]string, 0, len(mirrorSkillTrees)+len(presenceSkillTrees))
	allTrees = append(allTrees, mirrorSkillTrees...)
	allTrees = append(allTrees, presenceSkillTrees...)

	mirror := make(map[string]bool, len(mirrorSkillTrees))
	for _, tree := range mirrorSkillTrees {
		mirror[tree] = true
	}

	for _, tree := range allTrees {
		entries, dirErr := fs.ReadDir(assets.FS, tree)
		if dirErr != nil {
			t.Errorf("read skill tree %q: %v", tree, dirErr)
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			canonical, readErr := fs.ReadFile(assets.FS, path.Join(
				asset.DirClaudeSkills, entry.Name(), asset.FileSKILLMd,
			))
			if readErr != nil {
				if mirror[tree] {
					t.Errorf(
						"%s: orphaned — no canonical claude counterpart; "+
							"run the tree's sync script and commit",
						path.Join(tree, entry.Name()),
					)
					continue
				}
				// No ctx counterpart — tool-only skill, left
				// untouched by the presence-based sync script.
				exempt++
				continue
			}
			if mirror[tree] && claudeOnlySkills[entry.Name()] {
				t.Errorf(
					"%s: Claude-only skill must not ship in a mirror "+
						"tree; run the tree's sync script and commit",
					path.Join(tree, entry.Name()),
				)
				continue
			}
			generatedPath := path.Join(tree, entry.Name(), asset.FileSKILLMd)
			generated, genErr := fs.ReadFile(assets.FS, generatedPath)
			if genErr != nil {
				t.Errorf("%s: read: %v", generatedPath, genErr)
				continue
			}
			if want := stripAllowedTools(canonical); !bytes.Equal(generated, want) {
				t.Errorf(
					"%s: drifted from canonical claude source minus "+
						"allowed-tools — run 'make sync-opencode-skills "+
						"sync-codex-skills sync-copilot-skills' and commit",
					generatedPath,
				)
			}
			checked++
		}
	}

	// Completeness and reference parity for mirror trees.
	canonicalEntries, canonErr := fs.ReadDir(assets.FS, asset.DirClaudeSkills)
	if canonErr != nil {
		t.Fatalf("read canonical skill tree: %v", canonErr)
	}
	for _, entry := range canonicalEntries {
		if !entry.IsDir() || claudeOnlySkills[entry.Name()] {
			continue
		}
		for _, tree := range mirrorSkillTrees {
			generatedPath := path.Join(tree, entry.Name(), asset.FileSKILLMd)
			if _, genErr := fs.ReadFile(assets.FS, generatedPath); genErr != nil {
				t.Errorf(
					"%s: canonical skill missing from mirror tree — "+
						"run the tree's sync script and commit",
					generatedPath,
				)
			}
		}
		refDir := path.Join(
			asset.DirClaudeSkills, entry.Name(), asset.DirReferences)
		refEntries, refErr := fs.ReadDir(assets.FS, refDir)
		if refErr != nil {
			// No embedded references for this skill.
			continue
		}
		for _, ref := range refEntries {
			if ref.IsDir() {
				continue
			}
			canonicalRef, refReadErr := fs.ReadFile(
				assets.FS, path.Join(refDir, ref.Name()))
			if refReadErr != nil {
				t.Errorf("%s: read: %v", path.Join(refDir, ref.Name()), refReadErr)
				continue
			}
			for _, tree := range mirrorSkillTrees {
				generatedRefPath := path.Join(
					tree, entry.Name(), asset.DirReferences, ref.Name())
				generatedRef, genRefErr := fs.ReadFile(assets.FS, generatedRefPath)
				if genRefErr != nil {
					t.Errorf(
						"%s: canonical reference missing from mirror "+
							"tree — run the tree's sync script and commit",
						generatedRefPath,
					)
					continue
				}
				if !bytes.Equal(generatedRef, canonicalRef) {
					t.Errorf(
						"%s: drifted from canonical reference — run "+
							"the tree's sync script and commit",
						generatedRefPath,
					)
				}
				refsChecked++
			}
		}
	}

	if checked == 0 {
		t.Fatal("no synced SKILL.md files discovered — embed glob or tree constants regressed")
	}
	t.Logf(
		"verified %d synced skills (%d tool-only exempt, %d references) across %d trees",
		checked, exempt, refsChecked,
		len(mirrorSkillTrees)+len(presenceSkillTrees),
	)
}

// TestAllowedToolsConfinedToFrontmatter guards the sync transform's
// over-reach: `sed '/^allowed-tools:/d'` deletes every column-0
// `allowed-tools:` line anywhere in a file, so a canonical body that
// ever gains one (most plausibly a fenced frontmatter example in a
// skill about writing skills) would be silently corrupted in every
// generated tree — and TestSyncedSkillParity would still pass,
// because stripAllowedTools replicates the same transform. Assert
// every canonical skill carries at most one such line and that it
// sits inside the leading `---` frontmatter block.
func TestAllowedToolsConfinedToFrontmatter(t *testing.T) {
	fence := []byte("---")
	entries, dirErr := fs.ReadDir(assets.FS, asset.DirClaudeSkills)
	if dirErr != nil {
		t.Fatalf("read canonical skill tree: %v", dirErr)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		p := path.Join(asset.DirClaudeSkills, entry.Name(), asset.FileSKILLMd)
		src, readErr := fs.ReadFile(assets.FS, p)
		if readErr != nil {
			t.Errorf("%s: read: %v", p, readErr)
			continue
		}
		lines := bytes.Split(src, []byte("\n"))
		fmEnd := -1
		if len(lines) > 0 && bytes.Equal(lines[0], fence) {
			for i := 1; i < len(lines); i++ {
				if bytes.Equal(lines[i], fence) {
					fmEnd = i
					break
				}
			}
		}
		count := 0
		for i, line := range lines {
			if !bytes.HasPrefix(line, []byte("allowed-tools:")) {
				continue
			}
			count++
			if len(bytes.TrimSpace(line)) == len("allowed-tools:") {
				t.Errorf(
					"%s:%d: block-form allowed-tools: — the sync "+
						"transform strips only the key line and would "+
						"orphan its sequence items in every generated "+
						"tree; use the inline form",
					p, i+1,
				)
			}
			if fmEnd == -1 || i > fmEnd {
				t.Errorf(
					"%s:%d: allowed-tools: line outside the leading "+
						"frontmatter block — the sync transform would strip "+
						"it from generated bodies while parity still passes",
					p, i+1,
				)
			}
		}
		if count > 1 {
			t.Errorf("%s: %d allowed-tools: lines (want at most 1)", p, count)
		}
	}
}

// stripAllowedTools removes every line starting with
// `allowed-tools:`, mirroring the sync scripts'
// `sed '/^allowed-tools:/d'` transform exactly.
//
// Parameters:
//   - src: Canonical SKILL.md content
//
// Returns:
//   - []byte: Content with allowed-tools lines removed
func stripAllowedTools(src []byte) []byte {
	lines := bytes.Split(src, []byte("\n"))
	kept := make([][]byte, 0, len(lines))
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("allowed-tools:")) {
			continue
		}
		kept = append(kept, line)
	}
	return bytes.Join(kept, []byte("\n"))
}
