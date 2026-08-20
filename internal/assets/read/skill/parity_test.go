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

// syncedSkillTrees lists the embedded skill trees that are generated
// from the canonical Claude tree by a hack/sync-*-skills.sh script.
// Enrollment is opt-in by directory presence: a skill directory with
// no Claude counterpart is tool-only and exempt.
var syncedSkillTrees = []string{
	asset.DirIntegrationsOpenCodeSkill,
	asset.DirIntegrationsCopilotSkill,
}

// TestSyncedSkillParity asserts the sync contract at the test layer,
// where CI enforces it (`make check-*-skills` only runs via `make
// audit` on developer machines): every skill in a generated tree that
// has a canonical Claude counterpart must be byte-identical to that
// counterpart minus `allowed-tools:` lines (the Claude Code-specific
// frontmatter key the sync scripts strip).
func TestSyncedSkillParity(t *testing.T) {
	var checked, exempt int
	for _, tree := range syncedSkillTrees {
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
				// No ctx counterpart — tool-only skill, left
				// untouched by the sync script.
				exempt++
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
						"sync-copilot-skills' and commit",
					generatedPath,
				)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no synced SKILL.md files discovered — embed glob or tree constants regressed")
	}
	t.Logf("verified %d synced skills (%d tool-only exempt) across %d trees",
		checked, exempt, len(syncedSkillTrees))
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
