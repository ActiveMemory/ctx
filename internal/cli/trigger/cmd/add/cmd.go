//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errTrigger "github.com/ActiveMemory/ctx/internal/err/trigger"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trigger"
	writeTrigger "github.com/ActiveMemory/ctx/internal/write/trigger"
)

// Cmd returns the "ctx hook add" subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyTriggerAdd)

	return &cobra.Command{
		Use:     cmd.UseTriggerAdd,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyTriggerAdd),
		Args:    cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			return Run(c, args[0], args[1])
		},
	}
}

// Run creates a new hook script with a template.
//
// Parameters:
//   - c: The cobra command for output
//   - hookType: The hook type (e.g., "pre-tool-use")
//   - name: The hook script name (without .sh extension)
func Run(c *cobra.Command, hookType, name string) error {
	// Validate hook type.
	ht := hookType
	valid := trigger.ValidTypes()

	found := false
	for _, v := range valid {
		if v == ht {
			found = true
			break
		}
	}

	if !found {
		names := make([]string, len(valid))
		copy(names, valid)
		return errTrigger.InvalidType(hookType, strings.Join(names, token.CommaSpace))
	}

	hooksDir := rc.HooksDir()
	typeDir := filepath.Join(hooksDir, hookType)

	// Ensure the type directory exists.
	if mkdirErr := ctxIo.SafeMkdirAll(
		typeDir, fs.PermRestrictedDir,
	); mkdirErr != nil {
		return errTrigger.CreateDir(mkdirErr)
	}

	filePath := filepath.Join(typeDir, name+file.ExtSh)

	// Error if file already exists.
	if _, statErr := ctxIo.SafeStat(filePath); statErr == nil {
		return errTrigger.ScriptExists(filePath)
	}

	content := fmt.Sprintf(tpl.TriggerScript, name, hookType)
	writeErr := ctxIo.SafeWriteFile(
		filePath, []byte(content), fs.PermExec,
	)
	if writeErr != nil {
		return errTrigger.WriteScript(writeErr)
	}

	writeTrigger.Created(c, filePath)
	return nil
}
