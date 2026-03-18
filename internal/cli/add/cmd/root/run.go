//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"strings"

	entry2 "github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/write/add"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add/core"
	"github.com/ActiveMemory/ctx/internal/entry"
)

// Config is an alias for core.Config.
type Config = core.Config

// Run executes the add command logic.
//
// Reads content from the specified source (argument, file, or stdin),
// validates the entry, and writes it to the appropriate context file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Command arguments; args[0] is the entry type, args[1:] is content
//   - flags: All flag values from the command
//
// Returns:
//   - error: Non-nil if content is missing, type is invalid, required flags
//     are missing, or file operations fail
func Run(cmd *cobra.Command, args []string, flags Config) error {
	fType := strings.ToLower(args[0])

	content, extractErr := core.ExtractContent(args, flags)
	if extractErr != nil || content == "" {
		return add.ErrNoContentProvided(fType, core.ExamplesForType(fType))
	}

	params := entry.Params{
		Type:        fType,
		Content:     content,
		Section:     flags.Section,
		Priority:    flags.Priority,
		Context:     flags.Context,
		Rationale:   flags.Rationale,
		Consequence: flags.Consequence,
		Lesson:      flags.Lesson,
		Application: flags.Application,
	}

	if validateErr := entry.Validate(params, core.ExamplesForType); validateErr != nil {
		return validateErr
	}

	fName, ok := entry2.ToCtxFile[fType]
	if !ok {
		return add.ErrUnknownType(fType)
	}

	if writeErr := entry.Write(params); writeErr != nil {
		return writeErr
	}

	add.InfoAddedTo(cmd, fName)

	return nil
}
