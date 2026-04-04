//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"gopkg.in/yaml.v3"

	errSkill "github.com/ActiveMemory/ctx/internal/err/skill"
	"github.com/ActiveMemory/ctx/internal/parse"
)

// parseManifest extracts YAML frontmatter and markdown
// body from a SKILL.md file.
func parseManifest(
	data []byte, name, dir string,
) (*Skill, error) {
	raw, body, splitErr := parse.SplitFrontmatter(
		data,
	)
	if splitErr != nil {
		return nil, errSkill.Load(name, splitErr)
	}

	sk := &Skill{}
	if unmarshalErr := yaml.Unmarshal(
		raw, sk,
	); unmarshalErr != nil {
		return nil, errSkill.InvalidYAML(
			name, unmarshalErr,
		)
	}

	sk.Body = body
	sk.Dir = dir
	return sk, nil
}
