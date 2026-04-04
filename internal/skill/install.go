//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgSkill "github.com/ActiveMemory/ctx/internal/config/skill"
	errSkill "github.com/ActiveMemory/ctx/internal/err/skill"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Install copies a skill from the source directory into skillsDir/<name>/.
// The source must contain a valid SKILL.md with parseable YAML frontmatter.
// The skill name is derived from the parsed frontmatter.
func Install(source, skillsDir string) (*Skill, error) {
	manifestPath := filepath.Join(source, cfgSkill.SkillManifest)

	data, readErr := ctxIo.SafeReadUserFile(manifestPath)
	if readErr != nil {
		return nil, errSkill.NotValidSource(readErr)
	}

	sk, parseErr := parseManifest(data, filepath.Base(source), source)
	if parseErr != nil {
		return nil, errSkill.InvalidManifest(cfgSkill.SkillManifest, parseErr)
	}

	if sk.Name == "" {
		return nil, errSkill.MissingName(cfgSkill.SkillManifest)
	}

	destDir := filepath.Join(skillsDir, sk.Name)
	if mkdirErr := ctxIo.SafeMkdirAll(
		destDir, fs.PermRestrictedDir,
	); mkdirErr != nil {
		return nil, errSkill.CreateDest(mkdirErr)
	}

	if copyErr := copyDir(source, destDir); copyErr != nil {
		// Clean up partial copy on failure.
		_ = os.RemoveAll(destDir)
		return nil, errSkill.Install(sk.Name, copyErr)
	}

	sk.Dir = destDir
	return sk, nil
}
