//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// LoadState wraps a journal state loading failure.
//
// Parameters:
//   - cause: the underlying error.
//
// Returns:
//   - error: "load journal state: <cause>"
func LoadState(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalLoadJournalState), cause,
	)
}

// SaveState wraps a journal state saving failure.
//
// Parameters:
//   - cause: the underlying error.
//
// Returns:
//   - error: "save journal state: <cause>"
func SaveState(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalSaveJournalState), cause,
	)
}

// LoadStateErr wraps a failure to load journal processing state.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "load journal state: <cause>"
func LoadStateErr(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalLoadJournalState), cause,
	)
}

// LoadStateFailed wraps a journal state loading failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "load journal state: <cause>"
func LoadStateFailed(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalLoadJournalState), cause,
	)
}

// SaveStateFailed wraps a journal state save failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "save journal state: <cause>"
func SaveStateFailed(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalSaveJournalState), cause,
	)
}

// NoDir returns an error when the journal directory does not exist.
//
// Parameters:
//   - path: absolute path to the missing journal directory.
//
// Returns:
//   - error: includes a hint to run 'ctx recall export --all'
func NoDir(path string) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalNoJournalDir), path,
	)
}

// Scan wraps a journal scanning failure.
//
// Parameters:
//   - cause: the underlying scan error.
//
// Returns:
//   - error: "failed to scan journal: <cause>"
func Scan(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalScanJournal), cause,
	)
}

// NoEntries returns an error when the journal directory has no entries.
//
// Parameters:
//   - path: path to the empty journal directory.
//
// Returns:
//   - error: includes a hint to run 'ctx recall export --all'
func NoEntries(path string) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalNoJournalEntries), path,
	)
}

// NoEntriesMatch returns an error when a pattern matches nothing.
//
// Parameters:
//   - patterns: the patterns that matched nothing.
//
// Returns:
//   - error: "no journal entries match: <patterns>"
func NoEntriesMatch(patterns string) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalNoEntriesMatch), patterns,
	)
}

// ReadDir wraps a failure to read the journal directory.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "read journal directory: <cause>"
func ReadDir(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalReadJournalDir), cause,
	)
}

// UnknownStage returns an error for an unrecognized journal stage.
//
// Parameters:
//   - stage: the unknown stage name
//   - valid: comma-separated list of valid stage names
//
// Returns:
//   - error: "unknown stage <stage>; valid: <valid>"
func UnknownStage(stage, valid string) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalUnknownStage), stage, valid,
	)
}

// StageNotSet returns an error when a journal stage has not been set.
//
// Parameters:
//   - filename: the journal filename
//   - stage: the stage name
//
// Returns:
//   - error: "<filename>: <stage> not set"
func StageNotSet(filename, stage string) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrJournalStageNotSet), filename, stage,
	)
}

// RegenerateRequiresAll returns a validation error when --regenerate
// is used without --all.
//
// Returns:
//   - error: explains the flag dependency
func RegenerateRequiresAll() error {
	return errors.New(
		desc.TextDesc(text.DescKeyErrJournalRegenerateRequiresAll),
	)
}
