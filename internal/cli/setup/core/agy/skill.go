//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/agy"
	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	coreAgents "github.com/ActiveMemory/ctx/internal/cli/setup/core/agents"
	cfgAsset "github.com/ActiveMemory/ctx/internal/config/asset"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deploySkills creates .agents/skills/<name>/SKILL.md for each
// embedded Antigravity skill and deploys its references.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func deploySkills(cmd *cobra.Command) error {
	skills, readErr := agent.AgySkills()
	if readErr != nil {
		return readErr
	}
	refs, refErr := agent.SkillReferences(cfgAsset.DirAgySkills)
	if refErr != nil {
		return refErr
	}

	names := make([]string, 0, len(skills))
	for name := range skills {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		content := skills[name]
		skillDir := filepath.Join(cfgSetup.SkillsPathAgy, name)
		target := filepath.Join(skillDir, cfgHook.FileSKILLMd)
		if _, validateErr := validateManagedTarget(target); validateErr != nil {
			return validateErr
		}

		if existing, statErr := ctxIo.SafeReadUserFile(target); statErr == nil {
			if bytes.Equal(existing, content) {
				writeSetup.InfoAgySkipped(cmd, target)
				if refDeployErr := coreAgents.DeployReferences(
					cmd, cfgSetup.SkillsPathAgy, name, refs[name],
					validateManagedTarget, writeSetup.InfoAgyCreated,
				); refDeployErr != nil {
					return refDeployErr
				}
				continue
			}
			if !agy.SkillManaged(existing, name) {
				writeSetup.InfoAgyRejected(cmd, target)
				continue
			}
		} else if !os.IsNotExist(statErr) {
			return errFs.FileRead(target, statErr)
		}

		if mkErr := ctxIo.SafeMkdirAll(
			skillDir, fs.PermExec,
		); mkErr != nil {
			return errFs.Mkdir(skillDir, mkErr)
		}

		if wErr := ctxIo.SafeWriteFile(
			target, content, fs.PermFile,
		); wErr != nil {
			return errFs.FileWrite(target, wErr)
		}
		writeSetup.InfoAgyCreated(cmd, target)
		if refDeployErr := coreAgents.DeployReferences(
			cmd, cfgSetup.SkillsPathAgy, name, refs[name],
			validateManagedTarget, writeSetup.InfoAgyCreated,
		); refDeployErr != nil {
			return refDeployErr
		}
	}

	return nil
}
