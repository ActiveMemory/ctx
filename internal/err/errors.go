//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"
	"os"
)

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

// LoadKey classifies a key-loading failure.
//
// If the underlying error is os.ErrNotExist, returns NoKeyAt(keyPath).
// Otherwise wraps the cause as a generic load-key error.
//
// Parameters:
//   - cause: the underlying error from crypto.LoadKey
//   - keyPath: the resolved key path that was checked
//
// Returns:
//   - error: NoKeyAt or "load key: <cause>"
func LoadKey(cause error, keyPath string) error {
	if errors.Is(cause, os.ErrNotExist) {
		return NoKeyAt(keyPath)
	}
	return fmt.Errorf("load key: %w", cause)
}

// EncryptFailed wraps an encryption failure.
//
// Parameters:
//   - cause: the underlying error from crypto.Encrypt.
//
// Returns:
//   - error: "encrypt: <cause>"
func EncryptFailed(cause error) error {
	return fmt.Errorf("encrypt: %w", cause)
}

// DecryptFailed returns an error indicating decryption failure.
//
// Returns:
//   - error: "decryption failed: wrong key?"
func DecryptFailed() error {
	return fmt.Errorf("decryption failed: wrong key?")
}

// NoKeyAt returns an error indicating a missing encryption key.
//
// Parameters:
//   - path: the resolved key path that was checked.
//
// Returns:
//   - error: "encrypted scratchpad found but no key at <path>"
func NoKeyAt(path string) error {
	return fmt.Errorf("encrypted scratchpad found but no key at %s", path)
}

// EntryRange returns an error for an out-of-range scratchpad entry.
//
// Parameters:
//   - n: the requested entry number.
//   - total: the total number of entries.
//
// Returns:
//   - error: "entry <n> does not exist, scratchpad has <total> entries"
func EntryRange(n, total int) error {
	return fmt.Errorf("entry %d does not exist, scratchpad has %d entries", n, total)
}

// BoundaryViolation wraps a boundary validation error with a hint
// to use --allow-outside-cwd.
//
// Parameters:
//   - cause: the underlying validation error
//
// Returns:
//   - error: "<cause>\nUse --allow-outside-cwd to override this check"
func BoundaryViolation(cause error) error {
	return fmt.Errorf("%w\nUse --allow-outside-cwd to override this check", cause)
}

// NotInitialized returns an error indicating ctx has not been initialized.
//
// Returns:
//   - error: "ctx: not initialized — run \"ctx init\" first"
func NotInitialized() error {
	return fmt.Errorf("ctx: not initialized — run \"ctx init\" first")
}

// SkillList wraps a failure to list embedded skill directories.
//
// Parameters:
//   - cause: the underlying error from the list operation
//
// Returns:
//   - error: "failed to list skills: <cause>"
func SkillList(cause error) error {
	return fmt.Errorf("failed to list skills: %w", cause)
}

// SkillRead wraps a failure to read a skill's content.
//
// Parameters:
//   - name: Skill directory name that failed to read
//   - cause: the underlying error from the read operation
//
// Returns:
//   - error: "failed to read skill <name>: <cause>"
func SkillRead(name string, cause error) error {
	return fmt.Errorf("failed to read skill %s: %w", name, cause)
}

// DetectReferenceTime wraps a failure to detect the reference time for changes.
//
// Parameters:
//   - cause: the underlying detection error
//
// Returns:
//   - error: "detecting reference time: <cause>"
func DetectReferenceTime(cause error) error {
	return fmt.Errorf("detecting reference time: %w", cause)
}

// CreateArchiveDir wraps a failure to create the archive directory.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create archive directory: <cause>"
func CreateArchiveDir(cause error) error {
	return fmt.Errorf("failed to create archive directory: %w", cause)
}

// WriteArchive wraps a failure to write an archive file.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to write archive: <cause>"
func WriteArchive(cause error) error {
	return fmt.Errorf("failed to write archive: %w", cause)
}

// TaskFileNotFound returns an error when TASKS.md does not exist.
//
// Returns:
//   - error: "TASKS.md not found"
func TaskFileNotFound() error {
	return fmt.Errorf("TASKS.md not found")
}

// TaskFileRead wraps a failure to read TASKS.md.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "failed to read TASKS.md: <cause>"
func TaskFileRead(cause error) error {
	return fmt.Errorf("failed to read TASKS.md: %w", cause)
}

// TaskFileWrite wraps a failure to write TASKS.md.
//
// Parameters:
//   - cause: the underlying write error
//
// Returns:
//   - error: "failed to write TASKS.md: <cause>"
func TaskFileWrite(cause error) error {
	return fmt.Errorf("failed to write TASKS.md: %w", cause)
}

// TaskMultipleMatches returns an error when a query matches more than one task.
//
// Parameters:
//   - query: the search string that matched multiple tasks
//
// Returns:
//   - error: "multiple tasks match <query>; be more specific or use task number"
func TaskMultipleMatches(query string) error {
	return fmt.Errorf(
		"multiple tasks match %q; be more specific or use task number",
		query,
	)
}

// TaskNotFound returns an error when no task matches the query.
//
// Parameters:
//   - query: the search string that matched nothing
//
// Returns:
//   - error: "no task matching <query> found"
func TaskNotFound(query string) error {
	return fmt.Errorf("no task matching %q found", query)
}

// ReadEmbeddedSchema wraps a failure to read the embedded JSON Schema.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "read embedded schema: <cause>"
func ReadEmbeddedSchema(cause error) error {
	return fmt.Errorf("read embedded schema: %w", cause)
}

// LoadJournalStateErr wraps a failure to load journal processing state.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "load journal state: <cause>"
func LoadJournalStateErr(cause error) error {
	return fmt.Errorf("load journal state: %w", cause)
}

// UnknownProfile returns an error for an unrecognized config profile name.
//
// Parameters:
//   - name: the profile name that was not recognized
//
// Returns:
//   - error: "unknown profile <name>: must be dev, base, or prod"
func UnknownProfile(name string) error {
	return fmt.Errorf("unknown profile %q: must be dev, base, or prod", name)
}

// ReadProfile wraps a failure to read a profile file.
//
// Parameters:
//   - name: profile filename
//   - cause: the underlying read error
//
// Returns:
//   - error: "read <name>: <cause>"
func ReadProfile(name string, cause error) error {
	return fmt.Errorf("read %s: %w", name, cause)
}

// GitNotFound returns an error when git is not installed.
//
// Returns:
//   - error: "git not found in PATH"
func GitNotFound() error {
	return fmt.Errorf("git not found in PATH")
}

// NotInGitRepo wraps a failure from git rev-parse.
//
// Parameters:
//   - cause: the underlying exec error
//
// Returns:
//   - error: "not in a git repository: <cause>"
func NotInGitRepo(cause error) error {
	return fmt.Errorf("not in a git repository: %w", cause)
}

// UnknownFormat returns an error for an unsupported output format.
//
// Parameters:
//   - format: the format string that was not recognized
//   - supported: list of valid formats
//
// Returns:
//   - error: "unknown format <format> (supported: <list>)"
func UnknownFormat(format, supported string) error {
	return fmt.Errorf("unknown format %q (supported: %s)", format, supported)
}

// UnknownProjectType returns an error for an unsupported project type.
//
// Parameters:
//   - projType: the type string that was not recognized
//   - supported: list of valid types
//
// Returns:
//   - error: "unknown project type <type> (supported: <list>)"
func UnknownProjectType(projType, supported string) error {
	return fmt.Errorf("unknown project type %q (supported: %s)", projType, supported)
}

// InvalidTool returns an error for an unsupported AI tool name.
//
// Parameters:
//   - tool: the tool name that was not recognized
//
// Returns:
//   - error: "invalid tool <tool>: must be claude, aider, or generic"
func InvalidTool(tool string) error {
	return fmt.Errorf("invalid tool %q: must be claude, aider, or generic", tool)
}

// NoCompletedTasks returns an error when there are no completed tasks to archive.
//
// Returns:
//   - error: "no completed tasks to archive"
func NoCompletedTasks() error {
	return fmt.Errorf("no completed tasks to archive")
}

// NoTemplate wraps a failure to find an embedded template.
//
// Parameters:
//   - filename: Name of the file without a template
//   - cause: the underlying read error
//
// Returns:
//   - error: "no template available for <filename>: <cause>"
func NoTemplate(filename string, cause error) error {
	return fmt.Errorf("no template available for %s: %w", filename, cause)
}

// UnsupportedTool returns an error for an unrecognized AI tool name.
//
// Parameters:
//   - tool: the tool name that was not recognized
//
// Returns:
//   - error: "unsupported tool: <tool>"
func UnsupportedTool(tool string) error {
	return fmt.Errorf("unsupported tool: %s", tool)
}

// DriftViolations returns an error when drift detection found violations.
//
// Returns:
//   - error: "drift detection found violations"
func DriftViolations() error {
	return fmt.Errorf("drift detection found violations")
}

// ListTemplates wraps a failure to list embedded templates.
//
// Parameters:
//   - cause: the underlying error from the list operation
//
// Returns:
//   - error: "failed to list templates: <cause>"
func ListTemplates(cause error) error {
	return fmt.Errorf("failed to list templates: %w", cause)
}

// ReadTemplate wraps a failure to read an embedded template.
//
// Parameters:
//   - name: template name that failed to read
//   - cause: the underlying error from the read operation
//
// Returns:
//   - error: "failed to read template <name>: <cause>"
func ReadTemplate(name string, cause error) error {
	return fmt.Errorf("failed to read template %s: %w", name, cause)
}

// GenerateKey wraps a failure to generate an encryption key.
//
// Parameters:
//   - cause: the underlying error from key generation
//
// Returns:
//   - error: "failed to generate scratchpad key: <cause>"
func GenerateKey(cause error) error {
	return fmt.Errorf("failed to generate scratchpad key: %w", cause)
}

// SaveKey wraps a failure to save an encryption key.
//
// Parameters:
//   - cause: the underlying error from key saving
//
// Returns:
//   - error: "failed to save scratchpad key: <cause>"
func SaveKey(cause error) error {
	return fmt.Errorf("failed to save scratchpad key: %w", cause)
}

// MkdirKeyDir wraps a failure to create the key directory.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create key dir: <cause>"
func MkdirKeyDir(cause error) error {
	return fmt.Errorf("failed to create key dir: %w", cause)
}

// CreateBackup wraps a failure to create a backup file.
//
// Parameters:
//   - name: backup filename that could not be created
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create backup <name>: <cause>"
func CreateBackup(name string, cause error) error {
	return fmt.Errorf("failed to create backup %s: %w", name, cause)
}

// CreateBackupGeneric wraps a generic backup creation failure.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create backup: <cause>"
func CreateBackupGeneric(cause error) error {
	return fmt.Errorf("failed to create backup: %w", cause)
}

// WriteMerged wraps a failure to write a merged file.
//
// Parameters:
//   - path: file path that could not be written
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to write merged <path>: <cause>"
func WriteMerged(path string, cause error) error {
	return fmt.Errorf("failed to write merged %s: %w", path, cause)
}

// MarkerNotFound returns an error when a section marker is missing.
//
// Parameters:
//   - kind: marker kind (e.g. "ctx", "plan", "prompt")
//
// Returns:
//   - error: "<kind> start marker not found"
func MarkerNotFound(kind string) error {
	return fmt.Errorf("%s start marker not found", kind)
}

// TemplateMissingMarkers returns an error when a template lacks markers.
//
// Parameters:
//   - kind: marker kind (e.g. "ctx", "plan", "prompt")
//
// Returns:
//   - error: "template missing <kind> markers"
func TemplateMissingMarkers(kind string) error {
	return fmt.Errorf("template missing %s markers", kind)
}

// FileUpdate wraps a failure to update a file.
//
// Parameters:
//   - path: file path that could not be updated
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to update <path>: <cause>"
func FileUpdate(path string, cause error) error {
	return fmt.Errorf("failed to update %s: %w", path, cause)
}

// ParseFile wraps a failure to parse a file.
//
// Parameters:
//   - path: file path that could not be parsed
//   - cause: the underlying parse error
//
// Returns:
//   - error: "failed to parse %s: <cause>"
func ParseFile(path string, cause error) error {
	return fmt.Errorf("failed to parse %s: %w", path, cause)
}

// MarshalSettings wraps a failure to marshal settings JSON.
//
// Parameters:
//   - cause: the underlying marshal error
//
// Returns:
//   - error: "failed to marshal settings: <cause>"
func MarshalSettings(cause error) error {
	return fmt.Errorf("failed to marshal settings: %w", cause)
}

// ListPromptTemplates wraps a failure to list prompt templates.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to list prompt templates: <cause>"
func ListPromptTemplates(cause error) error {
	return fmt.Errorf("failed to list prompt templates: %w", cause)
}

// ReadPromptTemplate wraps a failure to read a prompt template.
//
// Parameters:
//   - name: template name that failed to read
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to read prompt template <name>: <cause>"
func ReadPromptTemplate(name string, cause error) error {
	return fmt.Errorf("failed to read prompt template %s: %w", name, cause)
}

// ListEntryTemplates wraps a failure to list entry templates.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to list entry templates: <cause>"
func ListEntryTemplates(cause error) error {
	return fmt.Errorf("failed to list entry templates: %w", cause)
}

// ReadEntryTemplate wraps a failure to read an entry template.
//
// Parameters:
//   - name: template name that failed to read
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to read entry template <name>: <cause>"
func ReadEntryTemplate(name string, cause error) error {
	return fmt.Errorf("failed to read entry template %s: %w", name, cause)
}

// HomeDir wraps a failure to determine the home directory.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "cannot determine home directory: <cause>"
func HomeDir(cause error) error {
	return fmt.Errorf("cannot determine home directory: %w", cause)
}

// MarshalPlugins wraps a failure to marshal enabledPlugins JSON.
//
// Parameters:
//   - cause: the underlying marshal error
//
// Returns:
//   - error: "failed to marshal enabledPlugins: <cause>"
func MarshalPlugins(cause error) error {
	return fmt.Errorf("failed to marshal enabledPlugins: %w", cause)
}

// FileAmend wraps a failure to amend an existing file.
//
// Parameters:
//   - path: file path that could not be amended
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to amend <path>: <cause>"
func FileAmend(path string, cause error) error {
	return fmt.Errorf("failed to amend %s: %w", path, cause)
}

// ReadProjectReadme wraps a failure to read a project README template.
//
// Parameters:
//   - dir: directory name whose README failed to read
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to read <dir> README template: <cause>"
func ReadProjectReadme(dir string, cause error) error {
	return fmt.Errorf("failed to read %s README template: %w", dir, cause)
}

// ReadInitTemplate wraps a failure to read an init template file.
//
// Parameters:
//   - name: template filename that failed to read
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to read <name> template: <cause>"
func ReadInitTemplate(name string, cause error) error {
	return fmt.Errorf("failed to read %s template: %w", name, cause)
}

// CreateMakefile wraps a failure to create a new Makefile.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create Makefile: <cause>"
func CreateMakefile(cause error) error {
	return fmt.Errorf("failed to create Makefile: %w", cause)
}
