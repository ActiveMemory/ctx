//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// EntryParams contains all parameters needed to add an entry to a context file.
//
// Fields:
//   - Type: Entry type (decision, learning, task, convention)
//   - Content: Title or main content
//   - Section: Target section (for tasks)
//   - Priority: Priority level (for tasks)
//   - Context: Context field (for decisions/learnings)
//   - Rationale: Rationale (for decisions)
//   - Consequences: Consequences (for decisions)
//   - Lesson: Lesson (for learnings)
//   - Application: Application (for learnings)
type EntryParams struct {
	Type         string
	Content      string
	Section      string
	Priority     string
	Context      string
	Rationale    string
	Consequences string
	Lesson       string
	Application  string
}

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
		return fmt.Errorf("no content provided")
	}

	switch config.UserInputToEntry(params.Type) {
	case config.EntryDecision:
		var missing []string
		if params.Context == "" {
			missing = append(missing, config.FieldContext)
		}
		if params.Rationale == "" {
			missing = append(missing, config.FieldRationale)
		}
		if params.Consequences == "" {
			missing = append(missing, config.FieldConsequence)
		}
		if len(missing) > 0 {
			return fmt.Errorf(
				"decision missing required fields: %s", strings.Join(missing, ", "),
			)
		}

	case config.EntryLearning:
		var missing []string
		if params.Context == "" {
			missing = append(missing, config.FieldContext)
		}
		if params.Lesson == "" {
			missing = append(missing, config.FieldLesson)
		}
		if params.Application == "" {
			missing = append(missing, config.FieldApplication)
		}
		if len(missing) > 0 {
			return fmt.Errorf(
				"learning missing required fields: %s", strings.Join(missing, ", "),
			)
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
		return fmt.Errorf("unknown type %q", fType)
	}

	filePath := filepath.Join(rc.ContextDir(), fileName)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf(
			"context file %s not found. Run 'ctx init' first", filePath,
		)
	}

	// Read existing content
	existing, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
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
		return fmt.Errorf("unknown type %q", fType)
	}

	// Append to file
	newContent := AppendEntry(existing, entry, fType, params.Section)

	if err := os.WriteFile(filePath, newContent, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}

	// Update index for decisions and learnings
	// (tasks/conventions don't have indexes)
	switch config.UserInputToEntry(fType) {
	case config.EntryDecision:
		indexed := index.UpdateDecisions(string(newContent))
		if err := os.WriteFile(filePath, []byte(indexed), 0644); err != nil {
			return fmt.Errorf("failed to update index in %s: %w", filePath, err)
		}
	case config.EntryLearning:
		indexed := index.UpdateLearnings(string(newContent))
		if err := os.WriteFile(filePath, []byte(indexed), 0644); err != nil {
			return fmt.Errorf("failed to update index in %s: %w", filePath, err)
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
		var missing []string
		if flags.context == "" {
			missing = append(missing, "--context")
		}
		if flags.rationale == "" {
			missing = append(missing, "--rationale")
		}
		if flags.consequences == "" {
			missing = append(missing, "--consequences")
		}
		if len(missing) > 0 {
			return errMissingDecision(missing)
		}
	case config.EntryLearning:
		var missing []string
		if flags.context == "" {
			missing = append(missing, "--context")
		}
		if flags.lesson == "" {
			missing = append(missing, "--lesson")
		}
		if flags.application == "" {
			missing = append(missing, "--application")
		}
		if len(missing) > 0 {
			return errMissingLearning(missing)
		}
	}

	// Validate type
	fName, ok := config.FileType[fType]
	if !ok {
		return fmt.Errorf(
			"unknown type %q. Valid types: decision, task, learning, convention",
			fType,
		)
	}

	// Write the entry using the shared function
	if err := WriteEntry(params); err != nil {
		return err
	}

	green := color.New(color.FgGreen).SprintFunc()
	cmd.Printf("%s Added to %s\n", green("✓"), fName)

	return nil
}
