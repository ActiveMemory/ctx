//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// EntryParams is an alias for core.EntryParams, kept for backward
// compatibility with callers that reference cmd.EntryParams directly.
type EntryParams = core.EntryParams

// Config is an alias for core.Config, kept for backward compatibility
// with callers that reference cmd.Config directly.
type Config = core.Config

// ValidateEntry checks that required fields are present for the given
// entry type.
//
// Parameters:
//   - params: Entry parameters to validate
//
// Returns:
//   - error: Non-nil with details about missing fields, nil if valid
func ValidateEntry(params EntryParams) error {
	if params.Content == "" {
		return core.ErrNoContentProvided(params.Type)
	}

	switch config.UserInputToEntry(params.Type) {
	case config.EntryDecision:
		if m := core.CheckRequired([][2]string{
			{config.FieldContext, params.Context},
			{config.FieldRationale, params.Rationale},
			{config.FieldConsequence, params.Consequences},
		}); len(m) > 0 {
			return core.ErrMissingFields(config.EntryDecision, m)
		}

	case config.EntryLearning:
		if m := core.CheckRequired([][2]string{
			{config.FieldContext, params.Context},
			{config.FieldLesson, params.Lesson},
			{config.FieldApplication, params.Application},
		}); len(m) > 0 {
			return core.ErrMissingFields(config.EntryLearning, m)
		}
	}

	return nil
}

// WriteEntry formats and writes an entry to the appropriate context file.
//
// This function handles the complete write cycle: read existing content,
// format the entry, append it, write back, and update the index if needed.
//
// Parameters:
//   - params: EntryParams containing type, content, and optional fields
//
// Returns:
//   - error: Non-nil if type is unknown, the file doesn't exist, or write fails
func WriteEntry(params EntryParams) error {
	fType := strings.ToLower(params.Type)

	fileName, ok := config.FileType[fType]
	if !ok {
		return core.ErrUnknownType(fType)
	}

	contextDir := params.ContextDir
	if contextDir == "" {
		contextDir = rc.ContextDir()
	}
	filePath := filepath.Join(contextDir, fileName)

	// Check if the file exists
	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return core.ErrFileNotFound(filePath)
	}

	// Read existing content
	existing, readErr := os.ReadFile(filepath.Clean(filePath))
	if readErr != nil {
		return core.ErrFileRead(filePath, readErr)
	}

	// Format the entry
	var entry string
	switch config.UserInputToEntry(fType) {
	case config.EntryDecision:
		entry = core.FormatDecision(
			params.Content, params.Context, params.Rationale, params.Consequences,
		)
	case config.EntryTask:
		entry = core.FormatTask(params.Content, params.Priority)
	case config.EntryLearning:
		entry = core.FormatLearning(
			params.Content, params.Context, params.Lesson, params.Application,
		)
	case config.EntryConvention:
		entry = core.FormatConvention(params.Content)
	default:
		return core.ErrUnknownType(fType)
	}

	// Append to file
	newContent := core.AppendEntry(existing, entry, fType, params.Section)

	if writeErr := os.WriteFile(filePath, newContent, config.PermFile); writeErr != nil {
		return core.ErrFileWrite(filePath, writeErr)
	}

	// Update index for decisions and learnings
	// (tasks/conventions don't have indexes)
	switch config.UserInputToEntry(fType) {
	case config.EntryDecision:
		indexed := index.UpdateDecisions(string(newContent))
		if indexErr := os.WriteFile(filePath, []byte(indexed), config.PermFile); indexErr != nil {
			return core.ErrIndexUpdate(filePath, indexErr)
		}
	case config.EntryLearning:
		indexed := index.UpdateLearnings(string(newContent))
		if indexErr := os.WriteFile(filePath, []byte(indexed), config.PermFile); indexErr != nil {
			return core.ErrIndexUpdate(filePath, indexErr)
		}
	case config.EntryTask, config.EntryConvention:
		// No index to update for these types
	}

	return nil
}

// Run executes the add command logic.
//
// It reads content from the specified source (argument, file, or stdin),
// validates the entry type, formats the entry, and appends it to the
// appropriate context file.
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

	// Determine the content source: args, --file, or stdin
	content, err := core.ExtractContent(args, flags)

	if err != nil || content == "" {
		return core.ErrNoContentProvided(fType)
	}

	// Build entry params
	params := EntryParams{
		Type:         fType,
		Content:      content,
		Section:      flags.Section,
		Priority:     flags.Priority,
		Context:      flags.Context,
		Rationale:    flags.Rationale,
		Consequences: flags.Consequences,
		Lesson:       flags.Lesson,
		Application:  flags.Application,
	}

	// Validate required fields with CLI-friendly error messages
	switch config.UserInputToEntry(fType) {
	case config.EntryDecision:
		if m := core.CheckRequired([][2]string{
			{"--context", flags.Context},
			{"--rationale", flags.Rationale},
			{"--consequences", flags.Consequences},
		}); len(m) > 0 {
			return core.ErrMissingDecision(m)
		}
	case config.EntryLearning:
		if m := core.CheckRequired([][2]string{
			{"--context", flags.Context},
			{"--lesson", flags.Lesson},
			{"--application", flags.Application},
		}); len(m) > 0 {
			return core.ErrMissingLearning(m)
		}
	}

	// Validate type
	fName, ok := config.FileType[fType]
	if !ok {
		return core.ErrUnknownType(fType)
	}

	// Write the entry using the shared function
	if writeErr := WriteEntry(params); writeErr != nil {
		return writeErr
	}

	green := color.New(color.FgGreen).SprintFunc()
	cmd.Println(fmt.Sprintf("%s Added to %s", green("✓"), fName))

	return nil
}
