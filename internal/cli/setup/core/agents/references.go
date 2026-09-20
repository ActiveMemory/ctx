//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agents

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	cfgAsset "github.com/ActiveMemory/ctx/internal/config/asset"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// DeployReferences writes a skill's reference files under its
// deployed directory. Call it only after the skill's SKILL.md was
// deployed or confirmed identical, so references never land in an
// unmanaged skill directory. Identical files are left untouched.
//
// Every tool that ships skill references deploys them the same way;
// the per-tool parts are the managed-target check and the "created"
// notice, which callers supply.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - skillsBase: Root of the tool's deployed skills tree
//   - name: Skill directory name
//   - files: Reference file name -> content (may be nil)
//   - validate: Tool's managed-target check for a path
//   - notify: Tool's "created" notice for a written path
//
// Returns:
//   - error: Non-nil if validation, directory creation, or a write fails
func DeployReferences(
	cmd *cobra.Command,
	skillsBase string,
	name string,
	files map[string][]byte,
	validate func(string) (bool, error),
	notify func(*cobra.Command, string),
) error {
	if len(files) == 0 {
		return nil
	}
	refDir := filepath.Join(skillsBase, name, cfgAsset.DirReferences)

	refNames := make([]string, 0, len(files))
	for refName := range files {
		refNames = append(refNames, refName)
	}
	sort.Strings(refNames)

	for _, refName := range refNames {
		target := filepath.Join(refDir, refName)
		if _, validateErr := validate(target); validateErr != nil {
			return validateErr
		}
		content := files[refName]
		if existing, statErr := ctxIo.SafeReadUserFile(target); statErr == nil {
			if bytes.Equal(existing, content) {
				continue
			}
		} else if !os.IsNotExist(statErr) {
			return errFs.FileRead(target, statErr)
		}
		if mkErr := ctxIo.SafeMkdirAll(refDir, fs.PermExec); mkErr != nil {
			return errFs.Mkdir(refDir, mkErr)
		}
		if wErr := ctxIo.SafeWriteFile(
			target, content, fs.PermFile,
		); wErr != nil {
			return errFs.FileWrite(target, wErr)
		}
		notify(cmd, target)
	}
	return nil
}
