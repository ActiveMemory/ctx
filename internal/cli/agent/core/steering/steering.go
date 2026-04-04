//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	errSkill "github.com/ActiveMemory/ctx/internal/err/skill"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/skill"
	"github.com/ActiveMemory/ctx/internal/steering"
)

// LoadBodies loads and filters steering files,
// returning their bodies as strings. Returns nil
// when the steering directory does not exist or
// contains no applicable files.
func LoadBodies() []string {
	steeringDir := rc.SteeringDir()

	files, loadErr := steering.LoadAll(steeringDir)
	if loadErr != nil {
		return nil
	}

	filtered := steering.Filter(
		files, "", nil, rc.Tool(),
	)

	var bodies []string
	for _, sf := range filtered {
		if sf.Body != "" {
			bodies = append(bodies, sf.Body)
		}
	}
	return bodies
}

// LoadSkill loads a named skill and returns its body
// content. Returns an error if the skill is not found.
func LoadSkill(name string) (string, error) {
	skillsDir := filepath.Join(
		rc.ContextDir(), dir.Skills,
	)

	sk, loadErr := skill.Load(skillsDir, name)
	if loadErr != nil {
		if errors.Is(loadErr, os.ErrNotExist) {
			return "", errSkill.NotFound(name)
		}
		return "", errSkill.LoadQuoted(name, loadErr)
	}
	return sk.Body, nil
}
