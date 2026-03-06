//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import "fmt"

// MemoryNotFound returns an error indicating that MEMORY.md was not
// discovered. Used by all memory subcommands (sync, status, diff).
//
// Returns:
//   - error: "MEMORY.md not found"
func MemoryNotFound() error {
	return fmt.Errorf("MEMORY.md not found")
}

// SyncFailed wraps a sync operation failure.
//
// Parameters:
//   - cause: the underlying error from the sync operation.
//
// Returns:
//   - error: "sync failed: <cause>"
func SyncFailed(cause error) error {
	return fmt.Errorf("sync failed: %w", cause)
}

// LoadState wraps a state-loading failure.
//
// Parameters:
//   - cause: the underlying error from loading state.
//
// Returns:
//   - error: "loading state: <cause>"
func LoadState(cause error) error {
	return fmt.Errorf("loading state: %w", cause)
}

// SaveState wraps a state-saving failure.
//
// Parameters:
//   - cause: the underlying error from saving state.
//
// Returns:
//   - error: "saving state: <cause>"
func SaveState(cause error) error {
	return fmt.Errorf("saving state: %w", cause)
}

// CtxNotInPath returns an error indicating that ctx was not found in PATH.
//
// Returns:
//   - error: "ctx not found in PATH"
func CtxNotInPath() error {
	return fmt.Errorf("ctx not found in PATH")
}

// WorkingDirectory wraps a failure to determine the working directory.
//
// Parameters:
//   - cause: the underlying error from os.Getwd.
//
// Returns:
//   - error: "failed to get working directory: <cause>"
func WorkingDirectory(cause error) error {
	return fmt.Errorf("failed to get working directory: %w", cause)
}

// FindSessions wraps a session-scanning failure.
//
// Parameters:
//   - cause: the underlying error from the parser.
//
// Returns:
//   - error: "failed to find sessions: <cause>"
func FindSessions(cause error) error {
	return fmt.Errorf("failed to find sessions: %w", cause)
}

// SessionNotFound returns an error for an unresolved session query.
//
// Parameters:
//   - query: the session ID or slug that was not found.
//
// Returns:
//   - error: "session not found: <query>"
func SessionNotFound(query string) error {
	return fmt.Errorf("session not found: %s", query)
}

// AmbiguousQuery returns an error when a session query matches
// multiple results.
//
// Returns:
//   - error: "ambiguous query, use a more specific ID"
func AmbiguousQuery() error {
	return fmt.Errorf("ambiguous query, use a more specific ID")
}

// Mkdir wraps a directory creation failure.
//
// Parameters:
//   - desc: human description of the directory (e.g. "journal directory").
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to create <desc>: <cause>"
func Mkdir(desc string, cause error) error {
	return fmt.Errorf("failed to create %s: %w", desc, cause)
}

// ReadInput wraps a failure to read user input.
//
// Parameters:
//   - cause: the underlying error from the read operation.
//
// Returns:
//   - error: "failed to read input: <cause>"
func ReadInput(cause error) error {
	return fmt.Errorf("failed to read input: %w", cause)
}

// NoSessionsFound returns an error when no sessions exist.
//
// Parameters:
//   - hint: additional guidance (e.g. "use --all-projects to search all").
//     Empty string omits the hint.
//
// Returns:
//   - error: "no sessions found" with optional hint
func NoSessionsFound(hint string) error {
	if hint != "" {
		return fmt.Errorf("no sessions found; %s", hint)
	}
	return fmt.Errorf("no sessions found")
}

// SessionIDRequired returns an error when no session ID was provided.
//
// Returns:
//   - error: "please provide a session ID or use --latest"
func SessionIDRequired() error {
	return fmt.Errorf("please provide a session ID or use --latest")
}

// AllWithArgument returns a validation error when --all is used alongside
// a positional argument.
//
// Parameters:
//   - argType: what the argument represents (e.g. "a session ID", "a pattern").
//
// Returns:
//   - error: "cannot use --all with <argType>; use one or the other"
func AllWithArgument(argType string) error {
	return fmt.Errorf(
		"cannot use --all with %s; use one or the other", argType,
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
	return fmt.Errorf("no journal entries match: %s", patterns)
}

// LoadJournalState wraps a journal state loading failure.
//
// Parameters:
//   - cause: the underlying error.
//
// Returns:
//   - error: "load journal state: <cause>"
func LoadJournalState(cause error) error {
	return fmt.Errorf("load journal state: %w", cause)
}

// SaveJournalState wraps a journal state saving failure.
//
// Parameters:
//   - cause: the underlying error.
//
// Returns:
//   - error: "save journal state: <cause>"
func SaveJournalState(cause error) error {
	return fmt.Errorf("save journal state: %w", cause)
}

// ReadDir wraps a directory read failure.
//
// Parameters:
//   - desc: human description of the directory (e.g. "journal directory").
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "read <desc>: <cause>"
func ReadDir(desc string, cause error) error {
	return fmt.Errorf("read %s: %w", desc, cause)
}

// RegenerateRequiresAll returns a validation error when --regenerate
// is used without --all.
//
// Returns:
//   - error: explains the flag dependency
func RegenerateRequiresAll() error {
	return fmt.Errorf(
		"--regenerate requires --all (single-session export always writes)",
	)
}

// InvalidDate returns a validation error for a malformed date flag.
//
// Parameters:
//   - flag: the flag name (e.g. "--since", "--until").
//   - value: the invalid date string.
//   - cause: the underlying parse error.
//
// Returns:
//   - error: formatted with the expected format hint
func InvalidDate(flag, value string, cause error) error {
	return fmt.Errorf(
		"invalid %s date %q (expected YYYY-MM-DD): %w", flag, value, cause,
	)
}

// ReadMemory wraps a failure to read MEMORY.md.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "reading MEMORY.md: <cause>"
func ReadMemory(cause error) error {
	return fmt.Errorf("reading MEMORY.md: %w", cause)
}

// WriteMemory wraps a failure to write MEMORY.md.
//
// Parameters:
//   - cause: the underlying write error.
//
// Returns:
//   - error: "writing MEMORY.md: <cause>"
func WriteMemory(cause error) error {
	return fmt.Errorf("writing MEMORY.md: %w", cause)
}

// FileWrite wraps a file write failure.
//
// Parameters:
//   - path: file path that could not be written.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to write <path>: <cause>"
func FileWrite(path string, cause error) error {
	return fmt.Errorf("failed to write %s: %w", path, cause)
}

// NoJournalDir returns an error when the journal directory does not exist.
//
// Parameters:
//   - path: absolute path to the missing journal directory.
//
// Returns:
//   - error: includes a hint to run 'ctx recall export --all'
func NoJournalDir(path string) error {
	return fmt.Errorf(
		"no journal directory found at %s\nRun 'ctx recall export --all' first",
		path,
	)
}

// ScanJournal wraps a journal scanning failure.
//
// Parameters:
//   - cause: the underlying scan error.
//
// Returns:
//   - error: "failed to scan journal: <cause>"
func ScanJournal(cause error) error {
	return fmt.Errorf("failed to scan journal: %w", cause)
}

// NoJournalEntries returns an error when the journal directory has no entries.
//
// Parameters:
//   - path: path to the empty journal directory.
//
// Returns:
//   - error: includes a hint to run 'ctx recall export --all'
func NoJournalEntries(path string) error {
	return fmt.Errorf(
		"no journal entries found in %s\nRun 'ctx recall export --all' first",
		path,
	)
}

// ZensicalNotFound returns an error when zensical is not installed.
//
// Returns:
//   - error: includes installation instructions
func ZensicalNotFound() error {
	return fmt.Errorf(
		"zensical not found. Install with: pipx install zensical (requires Python >= 3.10)",
	)
}
