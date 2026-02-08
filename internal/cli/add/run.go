//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

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
		return errNoContentProvided(params.Type)
	}

	switch config.UserInputToEntry(params.Type) {
	case config.EntryDecision:
		if m := checkRequired([][2]string{
			{config.FieldContext, params.Context},
			{config.FieldRationale, params.Rationale},
			{config.FieldConsequence, params.Consequences},
		}); len(m) > 0 {
			return errMissingFields(config.EntryDecision, m)
		}

	case config.EntryLearning:
		if m := checkRequired([][2]string{
			{config.FieldContext, params.Context},
			{config.FieldLesson, params.Lesson},
			{config.FieldApplication, params.Application},
		}); len(m) > 0 {
			return errMissingFields(config.EntryLearning, m)
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
		return errUnknownType(fType)
	}

	filePath := filepath.Join(rc.ContextDir(), fileName)

	// Check if the file exists
	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return errFileNotFound(filePath)
	}

	// Read existing content
	existing, readErr := os.ReadFile(filepath.Clean(filePath))
	if readErr != nil {
		return errFileRead(filePath, readErr)
	}

	// Format the entry
	var entry string
	switch config.UserInputToEntry(fType) {
	case config.EntryDecision:
		entry = FormatDecision(
			params.Content, params.Context, params.Rationale, params.Consequences,
		)
	case config.EntryTask:
		entry = FormatTask(params.Content, params.Priority)
	case config.EntryLearning:
		entry = FormatLearning(
			params.Content, params.Context, params.Lesson, params.Application,
		)
	case config.EntryConvention:
		entry = FormatConvention(params.Content)
	default:
		return errUnknownType(fType)
	}

	// Append to file
	newContent := AppendEntry(existing, entry, fType, params.Section)

	if writeErr := os.WriteFile(filePath, newContent, config.PermFile); writeErr != nil {
		return errFileWrite(filePath, writeErr)
	}

	// Update index for decisions and learnings
	// (tasks/conventions don't have indexes)
	switch config.UserInputToEntry(fType) {
	case config.EntryDecision:
		indexed := index.UpdateDecisions(string(newContent))
		if indexErr := os.WriteFile(filePath, []byte(indexed), config.PermFile); indexErr != nil {
			return errIndexUpdate(filePath, indexErr)
		}
	case config.EntryLearning:
		indexed := index.UpdateLearnings(string(newContent))
		if indexErr := os.WriteFile(filePath, []byte(indexed), config.PermFile); indexErr != nil {
			return errIndexUpdate(filePath, indexErr)
		}
	case config.EntryTask, config.EntryConvention:
		// No index to update for these types
	}

	return nil
}

// runAdd executes the add command logic.
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
func runAdd(cmd *cobra.Command, args []string, flags addConfig) error {
	fType := strings.ToLower(args[0])

	// Determine the content source: args, --file, or stdin
	content, err := extractContent(args, flags)

	if err != nil || content == "" {
		return errNoContentProvided(fType)
	}

	// Build entry params
	params := EntryParams{
		Type:         fType,
		Content:      content,
		Section:      flags.section,
		Priority:     flags.priority,
		Context:      flags.context,
		Rationale:    flags.rationale,
		Consequences: flags.consequences,
		Lesson:       flags.lesson,
		Application:  flags.application,
	}

	// Validate required fields with CLI-friendly error messages
	switch config.UserInputToEntry(fType) {
	case config.EntryDecision:
		if m := checkRequired([][2]string{
			{"--context", flags.context},
			{"--rationale", flags.rationale},
			{"--consequences", flags.consequences},
		}); len(m) > 0 {
			return errMissingDecision(m)
		}
	case config.EntryLearning:
		if m := checkRequired([][2]string{
			{"--context", flags.context},
			{"--lesson", flags.lesson},
			{"--application", flags.application},
		}); len(m) > 0 {
			return errMissingLearning(m)
		}
	}

	// Validate type
	fName, ok := config.FileType[fType]
	if !ok {
		return errUnknownType(fType)
	}

	// Write the entry using the shared function
	if writeErr := WriteEntry(params); writeErr != nil {
		return writeErr
	}

	green := color.New(color.FgGreen).SprintFunc()
	cmd.Printf("%s Added to %s\n", green("✓"), fName)

	return nil
}
