//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	coreEntry "github.com/ActiveMemory/ctx/internal/cli/add/core/entry"
	"github.com/ActiveMemory/ctx/internal/cli/add/core/example"
	"github.com/ActiveMemory/ctx/internal/cli/add/core/extract"
	corePub "github.com/ActiveMemory/ctx/internal/cli/connect/core/publish"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/entry"
	errAdd "github.com/ActiveMemory/ctx/internal/err/add"
	"github.com/ActiveMemory/ctx/internal/hub"
	"github.com/ActiveMemory/ctx/internal/trace"
	writeAdd "github.com/ActiveMemory/ctx/internal/write/add"
)

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
func Run(cmd *cobra.Command, args []string, flags entity.AddConfig) error {
	fType := strings.ToLower(args[0])

	content, extractErr := extract.Content(args, flags)
	if extractErr != nil || content == "" {
		return errAdd.NoContentProvided(fType, example.ForType(fType))
	}

	params := entity.EntryParams{
		Type:        fType,
		Content:     content,
		Section:     flags.Section,
		Priority:    flags.Priority,
		SessionID:   flags.SessionID,
		Branch:      flags.Branch,
		Commit:      flags.Commit,
		Context:     flags.Context,
		Rationale:   flags.Rationale,
		Consequence: flags.Consequence,
		Lesson:      flags.Lesson,
		Application: flags.Application,
	}

	if validateErr := entry.Validate(
		params, example.ForType,
	); validateErr != nil {
		return validateErr
	}

	fName, ok := cfgEntry.CtxFile(fType)
	if !ok {
		return errAdd.UnknownType(fType)
	}

	if writeErr := entry.Write(params); writeErr != nil {
		return writeErr
	}

	writeAdd.Added(cmd, fName)

	// Best-effort: publish to shared hub if --share is set.
	if flags.Share {
		pubEntry := hub.PublishEntry{
			Type:    fType,
			Content: content,
			Origin:  filepath.Base(state.Dir()),
		}
		_ = corePub.Run(
			cmd, []hub.PublishEntry{pubEntry},
		)
	}

	if fType == cfgEntry.Task && coreEntry.NeedsSpec(content) {
		writeAdd.SpecNudge(cmd)
	}

	// Best-effort: record pending context for commit tracing.
	// Decisions and learnings are prepended (see insert.AppendEntry),
	// so the new entry is always #1 in file order. This coupling is
	// intentional: if the prepend logic changes, this must be updated.
	if fType == cfgEntry.Decision || fType == cfgEntry.Learning {
		_ = trace.Record(fType+cfgTrace.RefFirstEntry, state.Dir())
	}

	return nil
}
