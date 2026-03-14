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

	"github.com/ActiveMemory/ctx/internal/assets"
)

// ReadingStateDir wraps a state directory read failure.
//
// Parameters:
//   - cause: the underlying error from reading the directory.
//
// Returns:
//   - error: "reading state directory: <cause>"
func ReadingStateDir(cause error) error {
	return fmt.Errorf("reading state directory: %w", cause)
}

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

// AllWithSessionID returns a validation error when --all is used with a session ID.
//
// Returns:
//   - error: "cannot use --all with a session ID; use one or the other"
func AllWithSessionID() error {
	return errors.New("cannot use --all with a session ID; use one or the other")
}

// AllWithPattern returns a validation error when --all is used with a pattern.
//
// Returns:
//   - error: "cannot use --all with a pattern; use one or the other"
func AllWithPattern() error {
	return errors.New("cannot use --all with a pattern; use one or the other")
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

// ReadReminders returns a validation error for a malformed date flag.
//
// Parameters:
//   - flag: the flag name (e.g. "--since", "--until").
//   - value: the invalid date string.
//   - cause: the underlying parse error.
//
// Returns:
//   - error: formatted with the expected format hint
//
// ReadReminders wraps a failure to read the reminders file.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read reminders: <cause>"
func ReadReminders(cause error) error {
	return fmt.Errorf("read reminders: %w", cause)
}

// ParseReminders wraps a failure to parse the reminders file.
//
// Parameters:
//   - cause: the underlying parse error.
//
// Returns:
//   - error: "parse reminders: <cause>"
func ParseReminders(cause error) error {
	return fmt.Errorf("parse reminders: %w", cause)
}

// InvalidID returns an error for an unparseable ID string.
//
// Parameters:
//   - value: the invalid ID string.
//
// Returns:
//   - error: "invalid ID <value>"
func InvalidID(value string) error {
	return fmt.Errorf("invalid ID %q", value)
}

// ReminderNotFound returns an error when no reminder matches the given ID.
//
// Parameters:
//   - id: the ID that was not found.
//
// Returns:
//   - error: "no reminder with ID <id>"
func ReminderNotFound(id int) error {
	return fmt.Errorf("no reminder with ID %d", id)
}

// ReminderIDRequired returns an error when no reminder ID is provided.
//
// Returns:
//   - error: "provide a reminder ID or use --all"
func ReminderIDRequired() error {
	return errors.New("provide a reminder ID or use --all")
}

// InvalidDateValue returns an error for an invalid date string.
//
// Parameters:
//   - value: the invalid date string.
//
// Returns:
//   - error: "invalid date <value> (expected YYYY-MM-DD)"
func InvalidDateValue(value string) error {
	return fmt.Errorf("invalid date %q (expected YYYY-MM-DD)", value)
}

func InvalidDate(flag, value string, cause error) error {
	return fmt.Errorf(
		"invalid %s date %q (expected YYYY-MM-DD): %w", flag, value, cause,
	)
}

// MemoryDiscoverFailed wraps a MEMORY.md discovery failure.
//
// Parameters:
//   - cause: the underlying discovery error.
//
// Returns:
//   - error: "MEMORY.md not found: <cause>"
func MemoryDiscoverFailed(cause error) error {
	return fmt.Errorf("MEMORY.md not found: %w", cause)
}

// MemoryDiffFailed wraps a memory diff computation failure.
//
// Parameters:
//   - cause: the underlying diff error.
//
// Returns:
//   - error: "computing diff: <cause>"
func MemoryDiffFailed(cause error) error {
	return fmt.Errorf("computing diff: %w", cause)
}

// SelectContentFailed wraps a content selection failure.
//
// Parameters:
//   - cause: the underlying selection error.
//
// Returns:
//   - error: "selecting content: <cause>"
func SelectContentFailed(cause error) error {
	return fmt.Errorf("selecting content: %w", cause)
}

// PublishFailed wraps a publish operation failure.
//
// Parameters:
//   - cause: the underlying publish error.
//
// Returns:
//   - error: "publishing: <cause>"
func PublishFailed(cause error) error {
	return fmt.Errorf("publishing: %w", cause)
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

// DirNotFound returns an error when a directory does not exist.
//
// Parameters:
//   - dir: the missing directory path.
//
// Returns:
//   - error: "directory not found: <dir>"
func DirNotFound(dir string) error {
	return fmt.Errorf("directory not found: %s", dir)
}

// NoSiteConfig returns an error when the zensical config file is missing.
//
// Parameters:
//   - dir: directory where the config was expected.
//
// Returns:
//   - error: "no zensical.toml found in <dir>"
func NoSiteConfig(dir string) error {
	return fmt.Errorf("no zensical.toml found in %s", dir)
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
// The message is loaded from assets and includes guidance for the user.
//
// Returns:
//   - error: message from assets key parser.git-not-found
func GitNotFound() error {
	return fmt.Errorf("%s", assets.TextDesc(assets.TextDescKeyParserGitNotFound))
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

// NoInput returns an error for missing stdin input.
//
// Returns:
//   - error: "no input received"
func NoInput() error {
	return errors.New("no input received")
}

// WebhookEmpty returns an error for blank webhook URL input.
//
// Returns:
//   - error: "webhook URL cannot be empty"
func WebhookEmpty() error {
	return errors.New("webhook URL cannot be empty")
}

// SaveWebhook wraps a webhook save failure.
//
// Parameters:
//   - cause: the underlying error from the save operation.
//
// Returns:
//   - error: "save webhook: <cause>"
func SaveWebhook(cause error) error {
	return fmt.Errorf("save webhook: %w", cause)
}

// LoadWebhook wraps a webhook load failure.
//
// Parameters:
//   - cause: the underlying error from the load operation.
//
// Returns:
//   - error: "load webhook: <cause>"
func LoadWebhook(cause error) error {
	return fmt.Errorf("load webhook: %w", cause)
}

// MarshalPayload wraps a JSON marshal failure.
//
// Parameters:
//   - cause: the underlying marshal error.
//
// Returns:
//   - error: "marshal payload: <cause>"
func MarshalPayload(cause error) error {
	return fmt.Errorf("marshal payload: %w", cause)
}

// SendNotification wraps a notification send failure.
//
// Parameters:
//   - cause: the underlying HTTP error.
//
// Returns:
//   - error: "send test notification: <cause>"
func SendNotification(cause error) error {
	return fmt.Errorf("send test notification: %w", cause)
}

// FlagRequired returns an error for a missing required flag.
//
// Parameters:
//   - name: the flag name.
//
// Returns:
//   - error: "required flag \"<name>\" not set"
func FlagRequired(name string) error {
	return fmt.Errorf("required flag %q not set", name)
}

// ParserReadFile wraps a session file read failure.
//
// Parameters:
//   - cause: the underlying error from reading the file.
//
// Returns:
//   - error: "read file: <cause>"
func ParserReadFile(cause error) error {
	return fmt.Errorf("read file: %w", cause)
}

// ArgRequired returns an error for a missing required argument.
//
// Parameters:
//   - name: the argument name.
//
// Returns:
//   - error: "<name> argument is required"
func ArgRequired(name string) error {
	return fmt.Errorf("%s argument is required", name)
}

// ReadFile wraps a file read failure.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read file: <cause>"
func ReadFile(cause error) error {
	return fmt.Errorf("read file: %w", cause)
}

// FileTooLarge returns an error for a file exceeding the size limit.
//
// Parameters:
//   - size: actual file size in bytes.
//   - max: maximum allowed size in bytes.
//
// Returns:
//   - error: "file too large: <size> bytes (max <max>)"
func FileTooLarge(size, max int) error {
	return fmt.Errorf("file too large: %d bytes (max %d)", size, max)
}

// InvalidIndex returns an error for a non-numeric entry index.
//
// Parameters:
//   - value: the invalid index string.
//
// Returns:
//   - error: "invalid index: <value>"
func InvalidIndex(value string) error {
	return fmt.Errorf("invalid index: %s", value)
}

// EditBlobTextConflict returns an error when --file/--label and text
// editing flags are used together.
//
// Returns:
//   - error: describing the mutual exclusivity
func EditBlobTextConflict() error {
	return errors.New("--file/--label and positional text/--append/--prepend are mutually exclusive")
}

// EditTextConflict returns an error when multiple text editing modes
// are used together.
//
// Returns:
//   - error: describing the mutual exclusivity
func EditTextConflict() error {
	return errors.New("--append, --prepend, and positional text are mutually exclusive")
}

// EditNoMode returns an error when no editing mode was specified.
//
// Returns:
//   - error: prompting for a mode
func EditNoMode() error {
	return errors.New("provide replacement text, --append, or --prepend")
}

// BlobAppendNotAllowed returns an error for appending to a blob entry.
//
// Returns:
//   - error: "cannot append to a blob entry"
func BlobAppendNotAllowed() error {
	return errors.New("cannot append to a blob entry")
}

// BlobPrependNotAllowed returns an error for prepending to a blob entry.
//
// Returns:
//   - error: "cannot prepend to a blob entry"
func BlobPrependNotAllowed() error {
	return errors.New("cannot prepend to a blob entry")
}

// NotBlobEntry returns an error when a blob operation targets a non-blob.
//
// Parameters:
//   - n: the 1-based entry index.
//
// Returns:
//   - error: "entry <n> is not a blob entry"
func NotBlobEntry(n int) error {
	return fmt.Errorf("entry %d is not a blob entry", n)
}

// OpenFile wraps a file open failure.
//
// Parameters:
//   - path: the file path.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "open <path>: <cause>"
func OpenFile(path string, cause error) error {
	return fmt.Errorf("open %s: %w", path, cause)
}

// StatPath wraps a stat failure.
//
// Parameters:
//   - path: the path that failed.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "stat <path>: <cause>"
func StatPath(path string, cause error) error {
	return fmt.Errorf("stat %s: %w", path, cause)
}

// NotDirectory returns an error when a path is not a directory.
//
// Parameters:
//   - path: the path.
//
// Returns:
//   - error: "<path> is not a directory"
func NotDirectory(path string) error {
	return fmt.Errorf("%s is not a directory", path)
}

// ReadDirectory wraps a directory read failure.
//
// Parameters:
//   - path: the directory path.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "read directory <path>: <cause>"
func ReadDirectory(path string, cause error) error {
	return fmt.Errorf("read directory %s: %w", path, cause)
}

// ResolveNotEncrypted returns an error when resolve is used on an
// unencrypted scratchpad.
//
// Returns:
//   - error: "resolve is only needed for encrypted scratchpads"
func ResolveNotEncrypted() error {
	return errors.New("resolve is only needed for encrypted scratchpads")
}

// NoConflictFiles returns an error when no merge conflict files are found.
//
// Parameters:
//   - filename: the base scratchpad filename.
//
// Returns:
//   - error: "no conflict files found (<filename>.ours / <filename>.theirs)"
func NoConflictFiles(filename string) error {
	return fmt.Errorf("no conflict files found (%s.ours / %s.theirs)", filename, filename)
}

// WriteFileFailed wraps a file write failure.
//
// Parameters:
//   - cause: the underlying write error.
//
// Returns:
//   - error: "write file: <cause>"
func WriteFileFailed(cause error) error {
	return fmt.Errorf("write file: %w", cause)
}

// OutFlagRequiresBlob returns an error when --out is used on a non-blob entry.
//
// Returns:
//   - error: "--out can only be used with blob entries"
func OutFlagRequiresBlob() error {
	return errors.New("--out can only be used with blob entries")
}

// ReadJournalDir wraps a failure to read the journal directory.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "read journal directory: <cause>"
func ReadJournalDir(cause error) error {
	return fmt.Errorf("read journal directory: %w", cause)
}

// SettingsNotFound returns an error when settings.local.json is missing.
//
// Returns:
//   - error: "no .claude/settings.local.json found"
func SettingsNotFound() error {
	return errors.New("no .claude/settings.local.json found")
}

// GoldenNotFound returns an error when settings.golden.json is missing.
//
// Returns:
//   - error: advises the user to run 'ctx permissions snapshot' first
func GoldenNotFound() error {
	return errors.New(
		"no .claude/settings.golden.json found — run 'ctx permissions snapshot' first",
	)
}

// FileRead wraps a file read failure with path context.
//
// Parameters:
//   - path: file path that could not be read.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to read <path>: <cause>"
func FileRead(path string, cause error) error {
	return fmt.Errorf("failed to read %s: %w", path, cause)
}

// PromptExists returns an error when a prompt template already exists.
//
// Parameters:
//   - name: the prompt name that already exists.
//
// Returns:
//   - error: "prompt <name> already exists"
func PromptExists(name string) error {
	return fmt.Errorf("prompt %q already exists", name)
}

// PromptNotFound returns an error when a prompt template does not exist.
//
// Parameters:
//   - name: the prompt name that was not found.
//
// Returns:
//   - error: "prompt <name> not found"
func PromptNotFound(name string) error {
	return fmt.Errorf("prompt %q not found", name)
}

// RemovePrompt wraps a failure to remove a prompt template.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "remove prompt: <cause>"
func RemovePrompt(cause error) error {
	return fmt.Errorf("remove prompt: %w", cause)
}

// NoPromptTemplate returns an error when no embedded template exists.
//
// Parameters:
//   - name: the template name that was not found.
//
// Returns:
//   - error: advises the user to use --stdin
func NoPromptTemplate(name string) error {
	return fmt.Errorf(
		"no embedded template %q — use --stdin to provide content", name,
	)
}

// ReadScratchpad wraps a scratchpad read failure.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read scratchpad: <cause>"
func ReadScratchpad(cause error) error {
	return fmt.Errorf("read scratchpad: %w", cause)
}

// ContextNotInitialized returns an error when no .context/ directory is found.
//
// Returns:
//   - error: "no .context/ directory found. Run 'ctx init' first"
func ContextNotInitialized() error {
	return errors.New("no .context/ directory found. Run 'ctx init' first")
}

// InvalidBackupScope returns an error for an unrecognized backup scope value.
//
// Parameters:
//   - scope: the invalid scope string
//
// Returns:
//   - error: "invalid scope '<scope>': must be project, global, or all"
func InvalidBackupScope(scope string) error {
	return fmt.Errorf("invalid scope %q: must be project, global, or all", scope)
}

// BackupSMBConfig wraps an SMB configuration parse failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "parse SMB config: <cause>"
func BackupSMBConfig(cause error) error {
	return fmt.Errorf("parse SMB config: %w", cause)
}

// BackupProject wraps a project backup failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "project backup: <cause>"
func BackupProject(cause error) error {
	return fmt.Errorf("project backup: %w", cause)
}

// BackupGlobal wraps a global backup failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "global backup: <cause>"
func BackupGlobal(cause error) error {
	return fmt.Errorf("global backup: %w", cause)
}

// CreateArchive wraps an archive creation failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "create archive file: <cause>"
func CreateArchive(cause error) error {
	return fmt.Errorf("create archive file: %w", cause)
}

// ContextDirNotFound returns an error when the context directory does not exist.
//
// Parameters:
//   - dir: the missing context directory path.
//
// Returns:
//   - error: "context directory not found: <dir> — run 'ctx init'"
func ContextDirNotFound(dir string) error {
	return fmt.Errorf("context directory not found: %s — run 'ctx init'", dir)
}

// SourceNotFound returns an error when a backup source path is missing.
//
// Parameters:
//   - path: the missing source path
//
// Returns:
//   - error: "source not found: <path>"
func SourceNotFound(path string) error {
	return fmt.Errorf("source not found: %s", path)
}

// EmbeddedTemplateNotFound returns an error when an embedded hook
// message template cannot be located.
//
// Parameters:
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - error: "embedded template not found for <hook>/<variant>"
func EmbeddedTemplateNotFound(hook, variant string) error {
	return fmt.Errorf("embedded template not found for %s/%s", hook, variant)
}

// OverrideExists returns an error when a message override already
// exists and must be reset before editing.
//
// Parameters:
//   - path: existing override file path
//   - hook: hook name
//   - variant: template variant name
//
// Returns:
//   - error: "override already exists at <path>..."
func OverrideExists(path, hook, variant string) error {
	return fmt.Errorf("override already exists at %s\nEdit it directly or use `ctx system message reset %s %s` first",
		path, hook, variant)
}

// CreateDir wraps a directory creation failure.
//
// Parameters:
//   - dir: the directory path that could not be created
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to create directory <dir>: <cause>"
func CreateDir(dir string, cause error) error {
	return fmt.Errorf("failed to create directory %s: %w", dir, cause)
}

// WriteOverride wraps a message override write failure.
//
// Parameters:
//   - path: the override file path
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to write override <path>: <cause>"
func WriteOverride(path string, cause error) error {
	return fmt.Errorf("failed to write override %s: %w", path, cause)
}

// RemoveOverride wraps a message override removal failure.
//
// Parameters:
//   - path: the override file path
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to remove override <path>: <cause>"
func RemoveOverride(path string, cause error) error {
	return fmt.Errorf("failed to remove override %s: %w", path, cause)
}

// UnknownHook returns an error for an unrecognized hook name.
//
// Parameters:
//   - hook: the unknown hook name
//
// Returns:
//   - error: "unknown hook: <hook>..."
func UnknownHook(hook string) error {
	return fmt.Errorf("unknown hook: %s\nRun `ctx system message list` to see available hooks", hook)
}

// UnknownVariant returns an error for an unrecognized variant within
// a known hook.
//
// Parameters:
//   - variant: the unknown variant name
//   - hook: the parent hook name
//
// Returns:
//   - error: "unknown variant <variant> for hook <hook>..."
func UnknownVariant(variant, hook string) error {
	return fmt.Errorf("unknown variant %q for hook %q\nRun `ctx system message list` to see available variants", variant, hook)
}

// LoadJournalStateFailed wraps a journal state loading failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "load journal state: <cause>"
func LoadJournalStateFailed(cause error) error {
	return fmt.Errorf("load journal state: %w", cause)
}

// SaveJournalStateFailed wraps a journal state save failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "save journal state: <cause>"
func SaveJournalStateFailed(cause error) error {
	return fmt.Errorf("save journal state: %w", cause)
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
	return fmt.Errorf("unknown stage %q; valid: %s", stage, valid)
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
	return fmt.Errorf("%s: %s not set", filename, stage)
}

// EventLogRead wraps a failure to read the event log.
//
// Parameters:
//   - cause: the underlying error from the query operation.
//
// Returns:
//   - error: "reading event log: <cause>"
func EventLogRead(cause error) error {
	return fmt.Errorf("reading event log: %w", cause)
}

// StatsGlob wraps a failure to glob stats files.
//
// Parameters:
//   - cause: the underlying error from the glob operation.
//
// Returns:
//   - error: "globbing stats files: <cause>"
func StatsGlob(cause error) error {
	return fmt.Errorf("globbing stats files: %w", cause)
}

// CryptoCreateCipher wraps a failure to create an AES cipher.
//
// Parameters:
//   - cause: the underlying crypto error.
//
// Returns:
//   - error: "create cipher: <cause>"
func CryptoCreateCipher(cause error) error {
	return fmt.Errorf("create cipher: %w", cause)
}

// CryptoCreateGCM wraps a failure to create a GCM instance.
//
// Parameters:
//   - cause: the underlying crypto error.
//
// Returns:
//   - error: "create GCM: <cause>"
func CryptoCreateGCM(cause error) error {
	return fmt.Errorf("create GCM: %w", cause)
}

// CryptoGenerateNonce wraps a failure to generate a random nonce.
//
// Parameters:
//   - cause: the underlying IO error.
//
// Returns:
//   - error: "generate nonce: <cause>"
func CryptoGenerateNonce(cause error) error {
	return fmt.Errorf("generate nonce: %w", cause)
}

// CryptoGenerateKey wraps a failure to generate a random key.
//
// Parameters:
//   - cause: the underlying IO error.
//
// Returns:
//   - error: "generate key: <cause>"
func CryptoGenerateKey(cause error) error {
	return fmt.Errorf("generate key: %w", cause)
}

// CryptoCiphertextTooShort returns an error when ciphertext is shorter
// than the nonce size.
//
// Returns:
//   - error: "ciphertext too short"
func CryptoCiphertextTooShort() error {
	return errors.New("ciphertext too short")
}

// CryptoDecrypt wraps a decryption failure with cause.
//
// Parameters:
//   - cause: the underlying decryption error.
//
// Returns:
//   - error: "decrypt: <cause>"
func CryptoDecrypt(cause error) error {
	return fmt.Errorf("decrypt: %w", cause)
}

// CryptoReadKey wraps a failure to read a key file.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read key: <cause>"
func CryptoReadKey(cause error) error {
	return fmt.Errorf("read key: %w", cause)
}

// CryptoInvalidKeySize returns an error when a key file has the wrong size.
//
// Parameters:
//   - got: actual key size in bytes.
//   - want: expected key size in bytes.
//
// Returns:
//   - error: "invalid key size: got N bytes, want M"
func CryptoInvalidKeySize(got, want int) error {
	return fmt.Errorf("invalid key size: got %d bytes, want %d", got, want)
}

// CryptoWriteKey wraps a failure to write a key file.
//
// Parameters:
//   - cause: the underlying write error.
//
// Returns:
//   - error: "write key: <cause>"
func CryptoWriteKey(cause error) error {
	return fmt.Errorf("write key: %w", cause)
}

// SnapshotWrite wraps a failure to write a task snapshot file.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to write snapshot: <cause>"
func SnapshotWrite(cause error) error {
	return fmt.Errorf("failed to write snapshot: %w", cause)
}

// OpenLogFile wraps a failure to open a log file.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to open log file: <cause>"
func OpenLogFile(cause error) error {
	return fmt.Errorf("failed to open log file: %w", cause)
}

// UnknownUpdateType returns an error for an unrecognized context update type.
//
// Parameters:
//   - typeName: the update type that was not recognized.
//
// Returns:
//   - error: "unknown update type: <typeName>"
func UnknownUpdateType(typeName string) error {
	return fmt.Errorf("unknown update type: %s", typeName)
}

// NoTaskSpecified returns an error when no task query was provided.
//
// Returns:
//   - error: "no task specified"
func NoTaskSpecified() error {
	return errors.New("no task specified")
}

// NoTaskMatch returns an error when no task matches the search query.
//
// Parameters:
//   - query: the search string that matched nothing.
//
// Returns:
//   - error: "no task matching \"<query>\" found"
func NoTaskMatch(query string) error {
	return fmt.Errorf("no task matching %q found", query)
}

// ReadInputStream wraps a failure to read from the input stream.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "error reading input: <cause>"
func ReadInputStream(cause error) error {
	return fmt.Errorf("error reading input: %w", cause)
}

// ReindexFileNotFound returns an error when the file to reindex does not exist.
//
// Parameters:
//   - fileName: Display name (e.g., "DECISIONS.md")
//
// Returns:
//   - error: "<fileName> not found. Run 'ctx init' first"
func ReindexFileNotFound(fileName string) error {
	return fmt.Errorf("%s not found. Run 'ctx init' first", fileName)
}

// ReindexFileRead wraps a read failure during reindexing.
//
// Parameters:
//   - filePath: Path that could not be read
//   - cause: The underlying read error
//
// Returns:
//   - error: "failed to read <filePath>: <cause>"
func ReindexFileRead(filePath string, cause error) error {
	return fmt.Errorf("failed to read %s: %w", filePath, cause)
}

// ReindexFileWrite wraps a write failure during reindexing.
//
// Parameters:
//   - filePath: Path that could not be written
//   - cause: The underlying write error
//
// Returns:
//   - error: "failed to write <filePath>: <cause>"
func ReindexFileWrite(filePath string, cause error) error {
	return fmt.Errorf("failed to write %s: %w", filePath, cause)
}

// DiscoverResolveRoot wraps a project root resolution failure.
func DiscoverResolveRoot(cause error) error {
	return fmt.Errorf("resolving project root: %w", cause)
}

// DiscoverResolveHome wraps a home directory resolution failure.
func DiscoverResolveHome(cause error) error {
	return fmt.Errorf("resolving home directory: %w", cause)
}

// DiscoverNoMemory returns an error when no auto memory file exists.
func DiscoverNoMemory(path string) error {
	return fmt.Errorf("no auto memory found at %s", path)
}

// MemoryReadSource wraps a source file read failure during sync.
func MemoryReadSource(cause error) error {
	return fmt.Errorf("reading source: %w", cause)
}

// MemoryArchivePrevious wraps a failure to archive the previous mirror.
func MemoryArchivePrevious(cause error) error {
	return fmt.Errorf("archiving previous mirror: %w", cause)
}

// MemoryCreateDir wraps a failure to create the memory directory.
func MemoryCreateDir(cause error) error {
	return fmt.Errorf("creating memory directory: %w", cause)
}

// MemoryWriteMirror wraps a failure to write the mirror file.
func MemoryWriteMirror(cause error) error {
	return fmt.Errorf("writing mirror: %w", cause)
}

// MemoryReadMirrorArchive wraps a failure to read the mirror for archiving.
func MemoryReadMirrorArchive(cause error) error {
	return fmt.Errorf("reading mirror for archive: %w", cause)
}

// MemoryCreateArchiveDir wraps a failure to create the archive directory.
func MemoryCreateArchiveDir(cause error) error {
	return fmt.Errorf("creating archive directory: %w", cause)
}

// MemoryWriteArchive wraps a failure to write an archive file.
func MemoryWriteArchive(cause error) error {
	return fmt.Errorf("writing archive: %w", cause)
}

// MemoryReadMirror wraps a failure to read the mirror file.
func MemoryReadMirror(cause error) error {
	return fmt.Errorf("reading mirror: %w", cause)
}

// MemoryReadDiffSource wraps a failure to read the source for diff.
func MemoryReadDiffSource(cause error) error {
	return fmt.Errorf("reading source: %w", cause)
}

// MemorySelectContent wraps a failure to select publish content.
func MemorySelectContent(cause error) error {
	return fmt.Errorf("selecting content: %w", cause)
}

// MemoryWriteMemory wraps a failure to write MEMORY.md.
func MemoryWriteMemory(cause error) error {
	return fmt.Errorf("writing MEMORY.md: %w", cause)
}

// ParserOpenFile wraps a session file open failure.
//
// Parameters:
//   - cause: the underlying error from opening the file.
//
// Returns:
//   - error: "open file: <cause>"
func ParserOpenFile(cause error) error {
	return fmt.Errorf("open file: %w", cause)
}

// ParserNoMatch returns an error when no parser can handle a file.
//
// Parameters:
//   - path: the file path that no parser matched.
//
// Returns:
//   - error: "no parser found for file: <path>"
func ParserNoMatch(path string) error {
	return fmt.Errorf("no parser found for file: %s", path)
}

// ParserWalkDir wraps a directory walk failure during session scanning.
//
// Parameters:
//   - cause: the underlying error from filepath.Walk.
//
// Returns:
//   - error: "walk directory: <cause>"
func ParserWalkDir(cause error) error {
	return fmt.Errorf("walk directory: %w", cause)
}

// ParserFileError wraps a per-file parse failure with the file path.
//
// Parameters:
//   - path: the file path that failed to parse.
//   - cause: the underlying parse error.
//
// Returns:
//   - error: "<path>: <cause>"
func ParserFileError(path string, cause error) error {
	return fmt.Errorf("%s: %w", path, cause)
}

// ParserScanFile wraps a session file scan failure.
//
// Parameters:
//   - cause: the underlying error from scanning the file.
//
// Returns:
//   - error: "scan file: <cause>"
func ParserScanFile(cause error) error {
	return fmt.Errorf("scan file: %w", cause)
}

// ParserUnmarshal wraps a JSON unmarshal failure during session parsing.
//
// Parameters:
//   - cause: the underlying error from JSON unmarshaling.
//
// Returns:
//   - error: "unmarshal: <cause>"
func ParserUnmarshal(cause error) error {
	return fmt.Errorf("unmarshal: %w", cause)
}
