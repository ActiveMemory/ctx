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
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// skillTrees lists every embedded directory under which each
// immediate subdirectory is a skill containing a SKILL.md.
var skillTrees = []string{
	asset.DirClaudeSkills,
	asset.DirIntegrationsOpenCodeSkill,
	asset.DirIntegrationsCopilotSkill,
}

// skillFrontmatter is the minimum frontmatter contract every
// embedded SKILL.md must satisfy. Per-surface extras
// (claude's `allowed-tools`, etc.) are intentionally not
// validated here — see specs/test-skill-frontmatter.md for
// scope.
type skillFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// TestSkillFrontmatter walks every embedded SKILL.md across
// the three tool trees and asserts the minimum frontmatter
// contract: `name` matches the containing directory's
// basename, and `description` is a non-empty string. All
// violations are reported in a single pass.
func TestSkillFrontmatter(t *testing.T) {
	var checked int
	for _, tree := range skillTrees {
		entries, err := fs.ReadDir(assets.FS, tree)
		if err != nil {
			t.Errorf("read skill tree %q: %v", tree, err)
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			skillPath := path.Join(tree, entry.Name(), asset.FileSKILLMd)
			body, readErr := fs.ReadFile(assets.FS, skillPath)
			if readErr != nil {
				t.Errorf("%s: read: %v", skillPath, readErr)
				continue
			}
			fm, fmErr := extractFrontmatter(body)
			if fmErr != nil {
				t.Errorf("%s: %v", skillPath, fmErr)
				continue
			}
			var parsed skillFrontmatter
			if err := yaml.Unmarshal(fm, &parsed); err != nil {
				t.Errorf("%s: yaml: %v", skillPath, err)
				continue
			}
			if parsed.Name == "" {
				t.Errorf("%s: missing or empty `name`", skillPath)
			} else if parsed.Name != entry.Name() {
				t.Errorf(
					"%s: name %q does not match directory %q",
					skillPath, parsed.Name, entry.Name(),
				)
			}
			if strings.TrimSpace(parsed.Description) == "" {
				t.Errorf("%s: missing or empty `description`", skillPath)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no SKILL.md files discovered — embed glob or tree constants regressed")
	}
	t.Logf("validated %d SKILL.md files across %d trees", checked, len(skillTrees))
}

// extractFrontmatter returns the YAML body between the first
// pair of `---` delimiter lines. Returns an error if either
// delimiter is missing.
func extractFrontmatter(body []byte) ([]byte, error) {
	const delim = "---"
	lines := bytes.Split(body, []byte("\n"))
	if len(lines) == 0 || string(bytes.TrimSpace(lines[0])) != delim {
		return nil, errMissingOpen
	}
	for i := 1; i < len(lines); i++ {
		if string(bytes.TrimSpace(lines[i])) == delim {
			return bytes.Join(lines[1:i], []byte("\n")), nil
		}
	}
	return nil, errMissingClose
}

var (
	errMissingOpen  = &frontmatterErr{msg: "missing opening `---` delimiter"}
	errMissingClose = &frontmatterErr{msg: "missing closing `---` delimiter"}
)

type frontmatterErr struct{ msg string }

func (e *frontmatterErr) Error() string { return e.msg }
