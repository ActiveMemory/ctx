//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// LoadJournalState wraps a journal state loading failure.
//
// Parameters:
//   - cause: the underlying error.
//
// Returns:
//   - error: "load journal state: <cause>"
func LoadJournalState(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalLoadJournalState), cause)
}

// SaveJournalState wraps a journal state saving failure.
//
// Parameters:
//   - cause: the underlying error.
//
// Returns:
//   - error: "save journal state: <cause>"
func SaveJournalState(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalSaveJournalState), cause)
}

// LoadJournalStateErr wraps a failure to load journal processing state.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "load journal state: <cause>"
func LoadJournalStateErr(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalLoadJournalState), cause)
}

// LoadJournalStateFailed wraps a journal state loading failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "load journal state: <cause>"
func LoadJournalStateFailed(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalLoadJournalState), cause)
}

// SaveJournalStateFailed wraps a journal state save failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "save journal state: <cause>"
func SaveJournalStateFailed(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalSaveJournalState), cause)
}

// NoJournalDir returns an error when the journal directory does not exist.
//
// Parameters:
//   - path: absolute path to the missing journal directory.
//
// Returns:
//   - error: includes a hint to run 'ctx recall export --all'
func NoJournalDir(path string) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalNoJournalDir), path)
}

// ScanJournal wraps a journal scanning failure.
//
// Parameters:
//   - cause: the underlying scan error.
//
// Returns:
//   - error: "failed to scan journal: <cause>"
func ScanJournal(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalScanJournal), cause)
}

// NoJournalEntries returns an error when the journal directory has no entries.
//
// Parameters:
//   - path: path to the empty journal directory.
//
// Returns:
//   - error: includes a hint to run 'ctx recall export --all'
func NoJournalEntries(path string) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalNoJournalEntries), path)
}

// NoEntriesMatch returns an error when a pattern matches nothing.
//
// Parameters:
//   - patterns: the patterns that matched nothing.
//
// Returns:
//   - error: "no journal entries match: <patterns>"
func NoEntriesMatch(patterns string) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalNoEntriesMatch), patterns)
}

// ReadJournalDir wraps a failure to read the journal directory.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "read journal directory: <cause>"
func ReadJournalDir(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalReadJournalDir), cause)
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
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalUnknownStage), stage, valid)
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
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrJournalStageNotSet), filename, stage)
}

// RegenerateRequiresAll returns a validation error when --regenerate
// is used without --all.
//
// Returns:
//   - error: explains the flag dependency
func RegenerateRequiresAll() error {
	return errors.New(assets.TextDesc(assets.TextDescKeyErrJournalRegenerateRequiresAll))
}
