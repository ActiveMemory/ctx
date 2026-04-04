//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config/token"
	errSkill "github.com/ActiveMemory/ctx/internal/err/skill"
)

// parseManifest extracts YAML frontmatter and markdown body from a
// SKILL.md file. The frontmatter is delimited by --- lines.
func parseManifest(data []byte, name, dir string) (*Skill, error) {
	raw, body, splitErr := splitFrontmatter(data)
	if splitErr != nil {
		return nil, errSkill.Load(name, splitErr)
	}

	sk := &Skill{}
	if unmarshalErr := yaml.Unmarshal(raw, sk); unmarshalErr != nil {
		return nil, errSkill.InvalidYAML(name, unmarshalErr)
	}

	sk.Body = body
	sk.Dir = dir
	return sk, nil
}

// splitFrontmatter separates YAML frontmatter from the markdown body.
// Frontmatter must start with a --- line and end with a second --- line.
func splitFrontmatter(
	data []byte,
) (frontmatter []byte, body string, err error) {
	content := strings.TrimLeft(string(data), trimCR)

	if !strings.HasPrefix(content, frontmatterDelimiter) {
		return nil, "", errSkill.MissingOpeningDelimiter()
	}

	// Skip the opening delimiter line.
	rest := content[len(frontmatterDelimiter):]
	rest = strings.TrimPrefix(rest, token.NewlineLF)

	needle := token.NewlineLF + frontmatterDelimiter
	idx := strings.Index(rest, needle)
	if idx < 0 {
		return nil, "", errSkill.MissingClosingDelimiter()
	}

	fm := rest[:idx]

	// Skip past the closing delimiter line.
	after := rest[idx+1+len(frontmatterDelimiter):]
	// Trim exactly one leading newline from the body if present.
	after = strings.TrimPrefix(after, token.NewlineLF)

	return []byte(fm), after, nil
}
