//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ActiveMemory/ctx/internal/assets/read/placeholders"
	"github.com/ActiveMemory/ctx/internal/config/asset"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
	cfgMemory "github.com/ActiveMemory/ctx/internal/config/memory"
	"github.com/ActiveMemory/ctx/internal/config/parser"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/crypto"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	"github.com/ActiveMemory/ctx/internal/i18n"
)

// Default returns a new CtxRC with hardcoded default values.
//
// Returns:
//   - *CtxRC: Configuration with defaults
//     (8000 token budget, 7-day archive, etc.)
func Default() *CtxRC {
	return &CtxRC{
		TokenBudget:         DefaultTokenBudget,
		PriorityOrder:       nil, // nil means use config.ReadOrder
		AutoArchive:         true,
		ArchiveAfterDays:    DefaultArchiveAfterDays,
		EntryCountLearnings: DefaultEntryCountLearnings,
		EntryCountDecisions: DefaultEntryCountDecisions,
		ConventionLineCount: DefaultConventionLineCount,
		InjectionTokenWarn:  DefaultInjectionTokenWarn,
		ContextWindow:       DefaultContextWindow,
		TaskNudgeInterval:   DefaultTaskNudgeInterval,
		StaleAgeDays:        DefaultStaleAgeDays,
	}
}

// RC returns the loaded configuration, initializing it on the first
// call.
//
// Under the cwd-anchored resolution model
// (spec: specs/cwd-anchored-context.md), `.ctxrc` is read from
// `$PWD/.ctxrc`: the project root, which by contract is the parent
// of [ContextDir]. When `$PWD/.context/` is absent, `.ctxrc` is not
// read and defaults apply. Environment overrides (CTX_TOKEN_BUDGET)
// are applied afterward. The result is cached for subsequent calls.
//
// Returns:
//   - *CtxRC: The loaded and cached configuration
func RC() *CtxRC {
	rcOnce.Do(func() {
		rc = load()
	})
	return rc
}

// ContextDir returns the project's context directory.
//
// Under the cwd-anchored resolution model
// (spec: specs/cwd-anchored-context.md), the answer is always
// `$PWD/.context/`: ctx anchors to its working directory the way
// `zensical` anchors to `zensical.toml` or Claude Code anchors to
// `$CLAUDE_PROJECT_DIR`. There is no env-var channel, no upward
// walk, no candidate scan. A single [os.Stat] checks the
// directory exists; absence and wrong-type are typed errors.
//
// Rejection conditions, in order:
//
//  1. [os.Getwd] failure: wrapped via [errCtx.StatFailed](".", err)
//     so callers can match [errCtx.ErrContextDirStat] uniformly.
//     Rare; usually an unlinked or permission-locked working
//     directory.
//  2. `$PWD/.context/` does not exist:
//     [errCtx.NoCtxHere](cwd) wrapping [errCtx.ErrNoCtxHere].
//     Callers that can proceed without a project (init, bootstrap
//     diagnostics) check with [errors.Is]; everyone else
//     propagates.
//  3. `$PWD/.context` exists but is a regular file (or other
//     non-directory): [errCtx.NotADir](path) wrapping
//     [errCtx.ErrContextDirNotADirectory].
//  4. Stat failed for another reason (permission, I/O):
//     [errCtx.StatFailed](path, cause) wrapping
//     [errCtx.ErrContextDirStat].
//
// Symlinks at `$PWD/.context` resolve transparently to their
// targets: a symlink-to-directory passes, a symlink-to-file is
// rejected as not-a-directory.
//
// Returns:
//   - string: absolute path to `$PWD/.context` when present.
//   - error: typed errCtx error depending on which check failed.
func ContextDir() (string, error) {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		// Surface the raw [os.Getwd] failure through the
		// [errCtx.ErrContextDirStat] sentinel so callers (and
		// tests) can match it with [errors.Is] alongside any
		// other stat-class diagnostic. The placeholder "." is
		// used because the actual cwd path is unknown by
		// definition once Getwd has failed (typically: process's
		// cwd was unlinked or chmod'd to deny lookup).
		return "", errCtx.StatFailed(token.Dot, cwdErr)
	}
	candidate := filepath.Join(cwd, dir.Context)
	info, statErr := os.Stat(candidate)
	if statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return "", errCtx.NoCtxHere(cwd)
		}
		return "", errCtx.StatFailed(candidate, statErr)
	}
	if !info.IsDir() {
		return "", errCtx.NotADir(candidate)
	}
	return candidate, nil
}

// TokenBudget returns the configured default token budget.
//
// Priority: env var > .ctxrc > default (8000).
//
// Returns:
//   - int: The token budget for context assembly
func TokenBudget() int {
	return RC().TokenBudget
}

// PriorityOrder returns the configured file priority order.
//
// Returns:
//   - []string: File names in priority order, or nil if not configured
//     (callers should fall back to config.ReadOrder)
func PriorityOrder() []string {
	return RC().PriorityOrder
}

// AutoArchive returns whether auto-archiving is enabled.
//
// Returns:
//   - bool: True if completed tasks should be auto-archived
func AutoArchive() bool {
	return RC().AutoArchive
}

// ArchiveAfterDays returns the configured days before archiving.
//
// Returns:
//   - int: Number of days after which completed tasks are archived (default 7)
func ArchiveAfterDays() int {
	return RC().ArchiveAfterDays
}

// ScratchpadEncrypt returns whether the scratchpad should be encrypted.
//
// Returns true (default) when the field is not set in .ctxrc.
//
// Returns:
//   - bool: True if scratchpad encryption is enabled (default true)
func ScratchpadEncrypt() bool {
	v := RC().ScratchpadEncrypt
	if v == nil {
		return true
	}
	return *v
}

// EntryCountLearnings returns the entry count threshold for LEARNINGS.md.
//
// Returns 0 if the check is disabled. Default: 30.
//
// Returns:
//   - int: Threshold above which a drift warning is emitted
func EntryCountLearnings() int {
	return RC().EntryCountLearnings
}

// EntryCountDecisions returns the entry count threshold for DECISIONS.md.
//
// Returns 0 if the check is disabled. Default: 20.
//
// Returns:
//   - int: Threshold above which a drift warning is emitted
func EntryCountDecisions() int {
	return RC().EntryCountDecisions
}

// ConventionLineCount returns the line count threshold for CONVENTIONS.md.
//
// Returns 0 if the check is disabled. Default: 200.
//
// Returns:
//   - int: Threshold above which a drift warning is emitted
func ConventionLineCount() int {
	return RC().ConventionLineCount
}

// InjectionTokenWarn returns the token threshold for
// oversize injection warning.
//
// Returns 0 if the check is disabled. Default: 15000.
//
// Returns:
//   - int: Threshold above which an oversize flag is written
func InjectionTokenWarn() int {
	return RC().InjectionTokenWarn
}

// BillingTokenWarn returns the absolute token threshold for billing warnings.
//
// Returns 0 (default, disabled). When set to a positive value, the
// check-context-size hook emits a one-shot VERBATIM warning the first
// time session tokens exceed this threshold.
//
// Returns:
//   - int: Token threshold, or 0 if disabled
func BillingTokenWarn() int {
	return RC().BillingTokenWarn
}

// ContextWindow returns the configured context window size in tokens.
//
// Returns 200000 (default). For Claude Code users this value is a no-op:
// the system hook auto-detects 200k vs 1M from ~/.claude/settings.json.
// Only useful as a manual override for non-Claude AI tools.
//
// Returns:
//   - int: Context window size in tokens
func ContextWindow() int {
	v := RC().ContextWindow
	if v <= 0 {
		return DefaultContextWindow
	}
	return v
}

// NotifyEvents returns the configured event filter list for notifications.
//
// Returns nil if Notify is nil (no filtering: all events pass).
//
// Returns:
//   - []string: Event names to allow, or nil for all
func NotifyEvents() []string {
	n := RC().Notify
	if n == nil {
		return nil
	}
	return n.Events
}

// KeyPath returns the resolved encryption key file path.
//
// Under the cwd-anchored model the caller must be at a project
// root. The previous implementation silently handed "" to
// [crypto.ResolveKeyPath] when ContextDir failed, which either
// filepath.Join'd a CWD-relative `.ctx.key` path or fell through
// to the global `~/.ctx/.ctx.key`: exactly the class of
// silent-wrong-location / wrong-key-rotation bug this branch aims
// to eliminate. The error is propagated instead so callers handle
// the absence of a project rather than rotating encryption
// against a surprise key.
//
// Within ResolveKeyPath the existing priority still applies:
// key_path in .ctxrc (explicit) > project-local
// (.context/.ctx.key) > global (~/.ctx/.ctx.key).
//
// Returns:
//   - string: Resolved path to the encryption key file
//   - error: [errCtx.ErrNoCtxHere] or any other ContextDir
//     resolver failure, propagated unchanged
func KeyPath() (string, error) {
	ctxDir, err := ContextDir()
	if err != nil {
		return "", err
	}
	return crypto.ResolveKeyPath(ctxDir, RC().KeyPathOverride), nil
}

// KeyRotationDays returns the configured key rotation threshold in days.
//
// The encryption key is shared by both ctx pad and ctx hook notify, so the
// rotation threshold is a project-wide setting.
//
// Priority: top-level key_rotation_days >
//
//	notify.key_rotation_days (legacy) > default (90).
//
// Returns:
//   - int: Number of days before a key rotation nudge
func KeyRotationDays() int {
	cfg := RC()
	if cfg.KeyRotationDays > 0 {
		return cfg.KeyRotationDays
	}
	if cfg.Notify != nil && cfg.Notify.KeyRotationDays > 0 {
		return cfg.Notify.KeyRotationDays
	}
	return DefaultKeyRotationDays
}

// TaskNudgeInterval returns the number of Edit/Write calls between task
// completion nudges. Returns 0 if disabled.
//
// Returns:
//   - int: Interval between nudges, or 0 if disabled
func TaskNudgeInterval() int {
	return RC().TaskNudgeInterval
}

// StaleAgeDays returns the number of days before a context file is
// flagged as stale by drift detection. Returns 0 if disabled.
//
// Returns:
//   - int: Days threshold, or 0 to disable the check
func StaleAgeDays() int {
	return RC().StaleAgeDays
}

// SessionPrefixes returns the list of recognized session header prefixes
// for the Markdown parser. Falls back to parser.DefaultSessionPrefixes
// when unconfigured or empty in .ctxrc.
//
// Returns:
//   - []string: Recognized prefixes (e.g., ["Session:"])
func SessionPrefixes() []string {
	prefixes := RC().SessionPrefixes
	if len(prefixes) == 0 {
		return parser.DefaultSessionPrefixes
	}
	return prefixes
}

// ClassifyRules returns the keyword rules for memory entry classification.
// Returns user-configured rules from .ctxrc if set, otherwise the built-in
// defaults from config/memory.
//
// Returns:
//   - []cfgMemory.ClassifyRule: Classification rules in priority order
func ClassifyRules() []cfgMemory.ClassifyRule {
	rules := RC().ClassifyRules
	if len(rules) == 0 {
		return cfgMemory.DefaultClassifyRules
	}
	return rules
}

// SpecSignalWords returns the terms that trigger a spec nudge
// when adding tasks. Returns user-configured words from .ctxrc
// if set, otherwise the built-in defaults from config/entry.
//
// Returns:
//   - []string: Signal words in lowercase
func SpecSignalWords() []string {
	words := RC().SpecSignalWords
	if len(words) == 0 {
		return cfgEntry.DefaultSpecSignalWords
	}
	return words
}

// SpecNudgeMinLen returns the task content length threshold for
// spec nudges. Returns user-configured value from .ctxrc if set,
// otherwise the built-in default from config/entry.
//
// Returns:
//   - int: Minimum content length to trigger a spec nudge
func SpecNudgeMinLen() int {
	n := RC().SpecNudgeMinLen
	if n == 0 {
		return cfgEntry.SpecNudgeMinLen
	}
	return n
}

// Placeholders returns the active rejected-placeholder set
// for body-flag validators (`ctx decision add` /
// `ctx learning add`). EXTEND semantics: the shipped
// defaults (loaded from
// `internal/assets/i18n/placeholders/<locale>.yaml`) are
// always present; any entries listed under `placeholders:`
// in `.ctxrc` are appended after normalization via
// [i18n.MatchKey] (Unicode case fold + diacritic strip)
// and whitespace trimming. Empty user entries are
// skipped; duplicates collapse to a single set
// membership.
//
// The set keys are MatchKey-normalized; callers compare
// against `i18n.MatchKey(strings.TrimSpace(input))`. The
// diacritic-insensitive contract means casual keyboard
// variation between vocabulary entries and user input
// matches transparently: `İPTAL` hits `iptal`, `Straße`
// hits `strasse`, `café` hits `cafe`.
//
// Returns:
//   - map[string]struct{}: normalized placeholder set
//     ready for O(1) lookup.
//   - error: non-nil only if the embedded defaults YAML
//     fails to load (build-time invariant violation).
func Placeholders() (map[string]struct{}, error) {
	defaults, loadErr := placeholders.Load(asset.LocaleEN)
	if loadErr != nil {
		return nil, loadErr
	}
	// Copy so caller mutations don't leak into the
	// loader's memoized cache.
	merged := make(map[string]struct{}, len(defaults)+len(RC().Placeholders))
	for k := range defaults {
		merged[k] = struct{}{}
	}
	for _, raw := range RC().Placeholders {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		merged[i18n.MatchKey(trimmed)] = struct{}{}
	}
	return merged, nil
}

// FreshnessFiles returns the configured list of files to track for
// freshness. Returns nil if no files are configured: the hook is
// a no-op when the list is empty.
//
// Returns:
//   - []FreshnessFile: Tracked files, or nil if unconfigured
func FreshnessFiles() []FreshnessFile {
	return RC().FreshnessFiles
}

// EventLog returns whether local hook event logging is enabled.
//
// Returns false (default) when the field is not set in .ctxrc.
//
// Returns:
//   - bool: True if hook events should be logged to .context/state/events.jsonl
func EventLog() bool {
	return RC().EventLog
}

// CompanionCheck returns whether the companion tool availability check
// should run during /ctx-remember. Returns true (default) unless
// explicitly set to false in .ctxrc.
//
// NOTE: No Go callers yet. The /ctx-remember skill currently reads
// this via ctx config status. This accessor exists for the planned
// hook-based companion check (see TASKS.md). Do not delete.
//
// Returns:
//   - bool: True if companion tools should be checked at the session start
func CompanionCheck() bool {
	cfg := RC()
	if cfg.CompanionCheck == nil {
		return true
	}
	return *cfg.CompanionCheck
}

// Tool returns the configured AI tool identifier (e.g., "claude", "cursor",
// "cline", "kiro", "codex").
//
// Returns an empty string when no tool is configured in .ctxrc.
//
// Returns:
//   - string: The tool identifier, or "" if not set
func Tool() string {
	return RC().Tool
}

// ProvenanceSessionRequired reports whether --session-id is
// required when adding tasks, decisions, and learnings.
// Returns true (default) unless explicitly disabled in .ctxrc.
//
// Returns:
//   - bool: True if --session-id is required
func ProvenanceSessionRequired() bool {
	cfg := RC()
	if cfg.ProvenanceRequired == nil ||
		cfg.ProvenanceRequired.SessionID == nil {
		return true
	}
	return *cfg.ProvenanceRequired.SessionID
}

// ProvenanceBranchRequired reports whether --branch is
// required when adding tasks, decisions, and learnings.
// Returns true (default) unless explicitly disabled in .ctxrc.
//
// Returns:
//   - bool: True if --branch is required
func ProvenanceBranchRequired() bool {
	cfg := RC()
	if cfg.ProvenanceRequired == nil ||
		cfg.ProvenanceRequired.Branch == nil {
		return true
	}
	return *cfg.ProvenanceRequired.Branch
}

// ProvenanceCommitRequired reports whether --commit is
// required when adding tasks, decisions, and learnings.
// Returns true (default) unless explicitly disabled in .ctxrc.
//
// Returns:
//   - bool: True if --commit is required
func ProvenanceCommitRequired() bool {
	cfg := RC()
	if cfg.ProvenanceRequired == nil ||
		cfg.ProvenanceRequired.Commit == nil {
		return true
	}
	return *cfg.ProvenanceRequired.Commit
}

// SteeringDir returns the configured steering directory path.
//
// Returns the value from .ctxrc steering.dir, or the default
// ".context/steering" when not configured.
//
// Returns:
//   - string: The steering directory path
func SteeringDir() string {
	cfg := RC()
	if cfg.Steering != nil && cfg.Steering.Dir != "" {
		return cfg.Steering.Dir
	}
	return DefaultSteeringDir
}

// HooksDir returns the configured hooks directory path.
//
// Returns the value from .ctxrc hooks.dir, or the default
// ".context/hooks" when not configured.
//
// Returns:
//   - string: The hooks directory path
func HooksDir() string {
	cfg := RC()
	if cfg.Hooks != nil && cfg.Hooks.Dir != "" {
		return cfg.Hooks.Dir
	}
	return DefaultHooksDir
}

// HookTimeout returns the configured per-hook execution timeout in seconds.
//
// Returns the value from .ctxrc hooks.timeout, or the default 10 seconds
// when not configured or set to zero.
//
// Returns:
//   - int: Timeout in seconds
func HookTimeout() int {
	cfg := RC()
	if cfg.Hooks != nil && cfg.Hooks.Timeout > 0 {
		return cfg.Hooks.Timeout
	}
	return DefaultHookTimeout
}

// HooksEnabled returns whether hook execution is enabled.
//
// Returns true (default) when the hooks section is not configured or
// when the enabled field is not explicitly set. Returns false only when
// hooks.enabled is explicitly set to false in .ctxrc.
//
// Returns:
//   - bool: True if hooks are enabled
func HooksEnabled() bool {
	cfg := RC()
	if cfg.Hooks != nil && cfg.Hooks.Enabled != nil {
		return *cfg.Hooks.Enabled
	}
	return true
}

// Reset clears the cached configuration, forcing
// reload on the next access.
func Reset() {
	rcMu.Lock()
	defer rcMu.Unlock()
	rcOnce = sync.Once{}
	rc = nil
}

// FilePriority returns the priority of a context file.
//
// If a priority_order is configured in .ctxrc, that order is used.
// Otherwise, the default config.ReadOrder is used.
//
// Lower numbers indicate higher priority (1 = highest).
// Unknown files return 100.
//
// Parameters:
//   - name: Filename to look up (e.g., "TASKS.md")
//
// Returns:
//   - int: Priority value (1-9 for known files, 100 for unknown)
func FilePriority(name string) int {
	// Check for .ctxrc override first
	if order := PriorityOrder(); order != nil {
		for i, fName := range order {
			if fName == name {
				return i + 1
			}
		}
		// File not in custom order gets the lowest priority
		return ctx.UnknownFilePriority
	}

	// Use the default priority from config.ReadOrder
	for i, fName := range ctx.ReadOrder {
		if fName == name {
			return i + 1
		}
	}
	return ctx.UnknownFilePriority
}
