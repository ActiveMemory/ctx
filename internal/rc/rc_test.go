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
	"testing"
	"time"

	cfgAgent "github.com/ActiveMemory/ctx/internal/config/agent"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/env"
	cfgJournal "github.com/ActiveMemory/ctx/internal/config/journal"
	cfgLoadgate "github.com/ActiveMemory/ctx/internal/config/loadgate"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	"github.com/ActiveMemory/ctx/internal/i18n"
)

// declareContext sets up a tempDir layout with a .context/ directory
// and a .ctxrc at the project root, t.Chdir's into tempDir so that
// `$PWD/.context` resolves to the test's .context, and resets the
// rc singleton. Mirrors the cwd-anchored resolution model
// (spec: specs/cwd-anchored-context.md): rc.ContextDir() is
// `filepath.Join($PWD, ".context")`, and .ctxrc is read from
// `$PWD/.ctxrc`.
//
// Parameters:
//   - t: test handle for Chdir/TempDir/Cleanup wiring.
//   - content: YAML body to write into .ctxrc; empty for "no file".
//
// Returns:
//   - string: absolute path of the test .context/ directory.
func declareContext(t *testing.T, content string) string {
	t.Helper()
	tempDir := t.TempDir()
	ctxDir := filepath.Join(tempDir, dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0700); mkErr != nil {
		t.Fatalf("mkdir .context: %v", mkErr)
	}
	if content != "" {
		rcPath := filepath.Join(tempDir, ".ctxrc")
		if wrErr := os.WriteFile(rcPath, []byte(content), 0600); wrErr != nil {
			t.Fatalf("write .ctxrc: %v", wrErr)
		}
	}
	t.Chdir(tempDir)
	Reset()
	t.Cleanup(Reset)
	return ctxDir
}

func TestDefaultRC(t *testing.T) {
	rc := Default()

	if rc.TokenBudget != DefaultTokenBudget {
		t.Errorf("TokenBudget = %d, want %d", rc.TokenBudget, DefaultTokenBudget)
	}
	if rc.PriorityOrder != nil {
		t.Errorf("PriorityOrder = %v, want nil", rc.PriorityOrder)
	}
	if !rc.AutoArchive {
		t.Error("AutoArchive = false, want true")
	}
	if rc.ArchiveAfterDays != DefaultArchiveAfterDays {
		t.Errorf(
			"ArchiveAfterDays = %d, want %d",
			rc.ArchiveAfterDays, DefaultArchiveAfterDays,
		)
	}
}

// TestGetRC_NoContext: cwd has no .context/ → defaults apply (no
// .ctxrc to read).
func TestGetRC_NoContext(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)
	Reset()
	t.Cleanup(Reset)

	rc := RC()

	if rc.TokenBudget != DefaultTokenBudget {
		t.Errorf("TokenBudget = %d, want %d", rc.TokenBudget, DefaultTokenBudget)
	}
	if !rc.AutoArchive {
		t.Error("AutoArchive = false, want true (default)")
	}
}

// TestGetRC_WithFile: cwd has .context/ and .ctxrc adjacent →
// values picked up.
func TestGetRC_WithFile(t *testing.T) {
	declareContext(t, `token_budget: 4000
priority_order:
  - TASKS.md
  - DECISIONS.md
auto_archive: false
archive_after_days: 14
`)

	rc := RC()

	if rc.TokenBudget != 4000 {
		t.Errorf("TokenBudget = %d, want %d", rc.TokenBudget, 4000)
	}
	if len(rc.PriorityOrder) != 2 || rc.PriorityOrder[0] != "TASKS.md" {
		t.Errorf("PriorityOrder = %v, want [TASKS.md DECISIONS.md]", rc.PriorityOrder)
	}
	if rc.AutoArchive {
		t.Error("AutoArchive = true, want false")
	}
	if rc.ArchiveAfterDays != 14 {
		t.Errorf("ArchiveAfterDays = %d, want %d", rc.ArchiveAfterDays, 14)
	}
}

// TestGetRC_TokenBudgetEnvOverride: CTX_TOKEN_BUDGET beats .ctxrc.
func TestGetRC_TokenBudgetEnvOverride(t *testing.T) {
	declareContext(t, `token_budget: 4000`)
	t.Setenv(env.CtxTokenBudget, "2000")
	Reset()

	rc := RC()
	if rc.TokenBudget != 2000 {
		t.Errorf("TokenBudget = %d, want %d (env override)", rc.TokenBudget, 2000)
	}
}

// TestContextDir_NoDotContext: cwd has no .context/ →
// errCtx.ErrNoCtxHere.
func TestContextDir_NoDotContext(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if !errors.Is(err, errCtx.ErrNoCtxHere) {
		t.Errorf("ContextDir() err = %v, want ErrNoCtxHere", err)
	}
	if got != "" {
		t.Errorf("ContextDir() = %q, want \"\"", got)
	}
}

// TestContextDir_RejectsNotADirectory: cwd has .context as a
// regular file → ErrContextDirNotADirectory.
func TestContextDir_RejectsNotADirectory(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, dir.Context)
	if err := os.WriteFile(filePath, []byte("oops"), 0600); err != nil {
		t.Fatalf("seed regular file: %v", err)
	}
	t.Chdir(tempDir)
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if !errors.Is(err, errCtx.ErrContextDirNotADirectory) {
		t.Errorf("ContextDir() err = %v, want ErrContextDirNotADirectory", err)
	}
	if got != "" {
		t.Errorf("ContextDir() = %q, want \"\"", got)
	}
}

// TestContextDir_AcceptsCwdContext: cwd has .context/ → returns
// `$PWD/.context`.
func TestContextDir_AcceptsCwdContext(t *testing.T) {
	ctxDir := declareContext(t, "")

	got, err := ContextDir()
	if err != nil {
		t.Fatalf("ContextDir() err = %v, want nil", err)
	}
	gotResolved, _ := filepath.EvalSymlinks(got)
	wantResolved, _ := filepath.EvalSymlinks(ctxDir)
	if gotResolved != wantResolved {
		t.Errorf("ContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}

// TestContextDir_AcceptsSymlinkDir: symlink at $PWD/.context
// pointing at a real directory passes (Stat follows symlinks).
func TestContextDir_AcceptsSymlinkDir(t *testing.T) {
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "actual-target")
	if err := os.MkdirAll(target, 0700); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	link := filepath.Join(tempDir, dir.Context)
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}
	t.Chdir(tempDir)
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if err != nil {
		t.Fatalf("ContextDir() err = %v, want nil", err)
	}
	gotResolved, _ := filepath.EvalSymlinks(got)
	wantResolved, _ := filepath.EvalSymlinks(link)
	if gotResolved != wantResolved {
		t.Errorf("ContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}

// TestRequireContextDir_Present: cwd has .context/ → path + nil err.
func TestRequireContextDir_Present(t *testing.T) {
	ctxDir := declareContext(t, "")

	got, err := RequireContextDir()
	if err != nil {
		t.Fatalf("RequireContextDir() err = %v, want nil", err)
	}
	gotResolved, _ := filepath.EvalSymlinks(got)
	wantResolved, _ := filepath.EvalSymlinks(ctxDir)
	if gotResolved != wantResolved {
		t.Errorf("RequireContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}

// TestRequireContextDir_Absent: cwd has no .context/ → typed error
// with a non-empty user-facing message.
func TestRequireContextDir_Absent(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)
	Reset()
	t.Cleanup(Reset)

	got, err := RequireContextDir()
	if err == nil {
		t.Fatalf("RequireContextDir() err = nil, want non-nil")
	}
	if !errors.Is(err, errCtx.ErrNoCtxHere) {
		t.Errorf("RequireContextDir() err = %v, want ErrNoCtxHere", err)
	}
	if got != "" {
		t.Errorf("RequireContextDir() path = %q, want \"\" on error", got)
	}
	if msg := err.Error(); msg == "" {
		t.Error("RequireContextDir() returned empty error message")
	}
}

func TestGetTokenBudget(t *testing.T) {
	declareContext(t, "")
	budget := TokenBudget()
	if budget != DefaultTokenBudget {
		t.Errorf("TokenBudget() = %d, want %d", budget, DefaultTokenBudget)
	}
}

func TestGetRC_InvalidYAML(t *testing.T) {
	declareContext(t, "invalid: [yaml: content")
	rc := RC()
	if rc.TokenBudget != DefaultTokenBudget {
		t.Errorf(
			"TokenBudget = %d, want %d (defaults on invalid YAML)",
			rc.TokenBudget, DefaultTokenBudget,
		)
	}
}

func TestGetRC_PartialConfig(t *testing.T) {
	declareContext(t, `token_budget: 5000`)
	rc := RC()
	if rc.TokenBudget != 5000 {
		t.Errorf("TokenBudget = %d, want %d", rc.TokenBudget, 5000)
	}
	if rc.ArchiveAfterDays != DefaultArchiveAfterDays {
		t.Errorf("ArchiveAfterDays = %d, want default", rc.ArchiveAfterDays)
	}
}

func TestGetRC_InvalidEnvBudget(t *testing.T) {
	declareContext(t, "")
	t.Setenv(env.CtxTokenBudget, "not-a-number")
	Reset()

	rc := RC()
	if rc.TokenBudget != DefaultTokenBudget {
		t.Errorf(
			"TokenBudget = %d, want %d (default on invalid env)",
			rc.TokenBudget, DefaultTokenBudget,
		)
	}
}

func TestGetRC_NegativeEnvBudget(t *testing.T) {
	declareContext(t, "")
	t.Setenv(env.CtxTokenBudget, "-100")
	Reset()

	rc := RC()
	if rc.TokenBudget != DefaultTokenBudget {
		t.Errorf(
			"TokenBudget = %d, want %d (default on negative env)",
			rc.TokenBudget, DefaultTokenBudget,
		)
	}
}

func TestGetRC_Singleton(t *testing.T) {
	declareContext(t, "")
	rc1 := RC()
	rc2 := RC()
	if rc1 != rc2 {
		t.Error("RC() should return same instance")
	}
}

func TestPriorityOrder(t *testing.T) {
	declareContext(t, "")
	if order := PriorityOrder(); order != nil {
		t.Errorf("PriorityOrder() = %v, want nil", order)
	}
}

func TestPriorityOrder_Custom(t *testing.T) {
	declareContext(t, `priority_order:
  - TASKS.md
  - DECISIONS.md
  - LEARNINGS.md
`)

	order := PriorityOrder()
	if len(order) != 3 {
		t.Fatalf("PriorityOrder() len = %d, want 3", len(order))
	}
	if order[0] != "TASKS.md" {
		t.Errorf("PriorityOrder()[0] = %q, want %q", order[0], "TASKS.md")
	}
}

func TestAutoArchive(t *testing.T) {
	declareContext(t, "")
	if !AutoArchive() {
		t.Error("AutoArchive() = false, want true")
	}
}

func TestAutoArchive_Disabled(t *testing.T) {
	declareContext(t, `auto_archive: false`)
	if AutoArchive() {
		t.Error("AutoArchive() = true, want false")
	}
}

func TestArchiveAfterDays(t *testing.T) {
	declareContext(t, "")
	days := ArchiveAfterDays()
	if days != DefaultArchiveAfterDays {
		t.Errorf("ArchiveAfterDays() = %d, want %d", days, DefaultArchiveAfterDays)
	}
}

func TestArchiveAfterDays_Custom(t *testing.T) {
	declareContext(t, `archive_after_days: 30`)
	days := ArchiveAfterDays()
	if days != 30 {
		t.Errorf("ArchiveAfterDays() = %d, want %d", days, 30)
	}
}

func TestScratchpadEncrypt_Default(t *testing.T) {
	declareContext(t, "")
	if !ScratchpadEncrypt() {
		t.Error("ScratchpadEncrypt() = false, want true (default)")
	}
}

func TestScratchpadEncrypt_Explicit(t *testing.T) {
	declareContext(t, `scratchpad_encrypt: false`)
	if ScratchpadEncrypt() {
		t.Error("ScratchpadEncrypt() = true, want false")
	}
}

func TestScratchpadEncrypt_ExplicitTrue(t *testing.T) {
	declareContext(t, `scratchpad_encrypt: true`)
	if !ScratchpadEncrypt() {
		t.Error("ScratchpadEncrypt() = false, want true")
	}
}

func TestFilePriority_DefaultOrder(t *testing.T) {
	declareContext(t, "")

	if p := FilePriority(ctx.Constitution); p != 1 {
		t.Errorf("FilePriority(%q) = %d, want 1", ctx.Constitution, p)
	}
	if p := FilePriority(ctx.Task); p != 2 {
		t.Errorf("FilePriority(%q) = %d, want 2", ctx.Task, p)
	}
	if p := FilePriority("UNKNOWN.md"); p != 100 {
		t.Errorf("FilePriority(%q) = %d, want 100", "UNKNOWN.md", p)
	}
}

func TestFilePriority_CustomOrder(t *testing.T) {
	declareContext(t, `priority_order:
  - DECISIONS.md
  - TASKS.md
`)

	if p := FilePriority(ctx.Decision); p != 1 {
		t.Errorf("FilePriority(%q) = %d, want 1", ctx.Decision, p)
	}
	if p := FilePriority(ctx.Task); p != 2 {
		t.Errorf("FilePriority(%q) = %d, want 2", ctx.Task, p)
	}
	if p := FilePriority("UNKNOWN.md"); p != 100 {
		t.Errorf("FilePriority(%q) = %d, want 100", "UNKNOWN.md", p)
	}
}

func TestNotifyEvents_Default(t *testing.T) {
	declareContext(t, "")
	if events := NotifyEvents(); events != nil {
		t.Errorf("NotifyEvents() = %v, want nil", events)
	}
}

func TestNotifyEvents_Configured(t *testing.T) {
	declareContext(t, `notify:
  events:
    - loop
    - nudge
`)

	events := NotifyEvents()
	if len(events) != 2 || events[0] != "loop" || events[1] != "nudge" {
		t.Errorf("NotifyEvents() = %v, want [loop nudge]", events)
	}
}

func TestKeyRotationDays_Default(t *testing.T) {
	declareContext(t, "")
	if days := KeyRotationDays(); days != DefaultKeyRotationDays {
		t.Errorf("KeyRotationDays() = %d, want %d", days, DefaultKeyRotationDays)
	}
}

func TestKeyRotationDays_Custom(t *testing.T) {
	declareContext(t, `key_rotation_days: 30
`)
	if days := KeyRotationDays(); days != 30 {
		t.Errorf("KeyRotationDays() = %d, want %d", days, 30)
	}
}

func TestKeyRotationDays_LegacyNotify(t *testing.T) {
	declareContext(t, `notify:
  key_rotation_days: 45
`)
	if days := KeyRotationDays(); days != 45 {
		t.Errorf("KeyRotationDays() = %d, want %d (legacy notify fallback)", days, 45)
	}
}

func TestKeyRotationDays_TopLevelTakesPrecedence(t *testing.T) {
	declareContext(t, `key_rotation_days: 60
notify:
  key_rotation_days: 45
`)
	if days := KeyRotationDays(); days != 60 {
		t.Errorf(
			"KeyRotationDays() = %d, want %d (top-level takes precedence)",
			days, 60,
		)
	}
}

func TestSessionPrefixes_Default(t *testing.T) {
	declareContext(t, "")
	prefixes := SessionPrefixes()
	if len(prefixes) != 1 || prefixes[0] != "Session:" {
		t.Errorf("SessionPrefixes() = %v, want [Session:]", prefixes)
	}
}

func TestSessionPrefixes_Custom(t *testing.T) {
	declareContext(t, "session_prefixes:\n"+
		"  - \"Session:\"\n"+
		"  - \"セッション:\"\n"+
		"  - \"Sesión:\"\n")

	prefixes := SessionPrefixes()
	if len(prefixes) != 3 {
		t.Fatalf("SessionPrefixes() len = %d, want 3", len(prefixes))
	}
	if prefixes[0] != "Session:" || prefixes[1] != "セッション:" || prefixes[2] != "Sesión:" {
		t.Errorf("SessionPrefixes() = %v", prefixes)
	}
}

func TestSessionPrefixes_EmptyFallsBackToDefault(t *testing.T) {
	declareContext(t, "session_prefixes: []\n")

	prefixes := SessionPrefixes()
	if len(prefixes) != 1 || prefixes[0] != "Session:" {
		t.Errorf(
			"SessionPrefixes() with empty config = %v, want defaults [Session:]",
			prefixes,
		)
	}
}

func TestTool_Default(t *testing.T) {
	declareContext(t, "")
	if tool := Tool(); tool != "" {
		t.Errorf("Tool() = %q, want %q", tool, "")
	}
}

func TestTool_Configured(t *testing.T) {
	declareContext(t, `tool: kiro`)
	if tool := Tool(); tool != "kiro" {
		t.Errorf("Tool() = %q, want %q", tool, "kiro")
	}
}

func TestSteeringDir_Default(t *testing.T) {
	declareContext(t, "")
	if d := SteeringDir(); d != DefaultSteeringDir {
		t.Errorf("SteeringDir() = %q, want %q", d, DefaultSteeringDir)
	}
}

func TestSteeringDir_Configured(t *testing.T) {
	declareContext(t, `steering:
  dir: custom/steering
`)
	if d := SteeringDir(); d != "custom/steering" {
		t.Errorf("SteeringDir() = %q, want %q", d, "custom/steering")
	}
}

func TestHooksDir_Default(t *testing.T) {
	declareContext(t, "")
	if d := HooksDir(); d != DefaultHooksDir {
		t.Errorf("HooksDir() = %q, want %q", d, DefaultHooksDir)
	}
}

func TestHooksDir_Configured(t *testing.T) {
	declareContext(t, `hooks:
  dir: custom/hooks
`)
	if d := HooksDir(); d != "custom/hooks" {
		t.Errorf("HooksDir() = %q, want %q", d, "custom/hooks")
	}
}

func TestHookTimeout_Default(t *testing.T) {
	declareContext(t, "")
	if timeout := HookTimeout(); timeout != DefaultHookTimeout {
		t.Errorf("HookTimeout() = %d, want %d", timeout, DefaultHookTimeout)
	}
}

func TestHookTimeout_Configured(t *testing.T) {
	declareContext(t, `hooks:
  timeout: 30
`)
	if timeout := HookTimeout(); timeout != 30 {
		t.Errorf("HookTimeout() = %d, want %d", timeout, 30)
	}
}

func TestHooksEnabled_Default(t *testing.T) {
	declareContext(t, "")
	if !HooksEnabled() {
		t.Error("HooksEnabled() = false, want true (default)")
	}
}

func TestHooksEnabled_ExplicitFalse(t *testing.T) {
	declareContext(t, `hooks:
  enabled: false
`)
	if HooksEnabled() {
		t.Error("HooksEnabled() = true, want false")
	}
}

func TestPlaceholders_DefaultsOnly(t *testing.T) {
	declareContext(t, "")
	set, err := Placeholders()
	if err != nil {
		t.Fatalf("Placeholders(): %v", err)
	}
	// Shipped en.yaml has 9 entries.
	if len(set) != 9 {
		t.Errorf("len(set) = %d, want 9 (en defaults only)", len(set))
	}
	for _, want := range []string{"tbd", "n/a", "see chat", "to be done"} {
		if _, ok := set[want]; !ok {
			t.Errorf("set missing default %q", want)
		}
	}
}

func TestPlaceholders_UserExtendsDefaults(t *testing.T) {
	declareContext(t, "placeholders:\n"+
		"  - iptal\n"+
		"  - yapılacak\n")
	set, err := Placeholders()
	if err != nil {
		t.Fatalf("Placeholders(): %v", err)
	}
	if len(set) != 11 {
		t.Errorf("len(set) = %d, want 11 (9 defaults + 2 user)", len(set))
	}
	for _, want := range []string{"tbd", "iptal", "yapılacak"} {
		if _, ok := set[want]; !ok {
			t.Errorf("set missing %q (default+user merge)", want)
		}
	}
}

func TestPlaceholders_NormalizesUserEntriesDiacriticInsensitively(t *testing.T) {
	declareContext(t, "placeholders:\n"+
		"  - İPTAL\n")
	set, err := Placeholders()
	if err != nil {
		t.Fatalf("Placeholders(): %v", err)
	}
	// "İPTAL" MatchKey-normalizes to "iptal" (Fold gives
	// "i̇ptal"; strip-combining drops U+0307 → "iptal").
	// Validator inputs go through the same MatchKey, so
	// İPTAL / İptal / IPTAL / iptal all converge on the
	// same set key.
	for _, variant := range []string{"İPTAL", "İptal", "IPTAL", "iptal"} {
		if _, ok := set[i18n.MatchKey(variant)]; !ok {
			t.Errorf("set missing %q (MatchKey=%q)", variant, i18n.MatchKey(variant))
		}
	}
}

func TestPlaceholders_TrimsAndSkipsEmptyUserEntries(t *testing.T) {
	declareContext(t, "placeholders:\n"+
		"  - \"  spaced  \"\n"+
		"  - \"\"\n"+
		"  - \"   \"\n")
	set, err := Placeholders()
	if err != nil {
		t.Fatalf("Placeholders(): %v", err)
	}
	// 9 defaults + 1 real user entry (after trim+empty-skip).
	if len(set) != 10 {
		t.Errorf("len(set) = %d, want 10 (9 defaults + 1 trimmed user)", len(set))
	}
	if _, ok := set["spaced"]; !ok {
		t.Errorf("set missing trimmed user entry %q", "spaced")
	}
}

func TestPlaceholders_UserDuplicateOfDefaultDedupes(t *testing.T) {
	declareContext(t, "placeholders:\n"+
		"  - TBD\n"+
		"  - tbd\n"+
		"  - new-marker\n")
	set, err := Placeholders()
	if err != nil {
		t.Fatalf("Placeholders(): %v", err)
	}
	// 9 defaults + 1 new (TBD/tbd both fold to existing "tbd").
	if len(set) != 10 {
		t.Errorf("len(set) = %d, want 10 (defaults + 1 distinct user)", len(set))
	}
	if _, ok := set["new-marker"]; !ok {
		t.Errorf("set missing %q", "new-marker")
	}
}

func TestAutoPruneDays_Default(t *testing.T) {
	declareContext(t, "")
	if got := AutoPruneDays(); got != cfgLoadgate.AutoPruneStaleDays {
		t.Errorf("AutoPruneDays() = %d, want %d", got, cfgLoadgate.AutoPruneStaleDays)
	}
}

func TestAutoPruneDays_Custom(t *testing.T) {
	declareContext(t, `auto_prune_days: 14`)
	if got := AutoPruneDays(); got != 14 {
		t.Errorf("AutoPruneDays() = %d, want %d", got, 14)
	}
}

func TestAgentCooldownMinutes_Default(t *testing.T) {
	declareContext(t, "")
	want := int(cfgAgent.DefaultCooldown / time.Minute)
	if got := AgentCooldownMinutes(); got != want {
		t.Errorf("AgentCooldownMinutes() = %d, want %d", got, want)
	}
}

func TestAgentCooldownMinutes_Custom(t *testing.T) {
	declareContext(t, `agent_cooldown_minutes: 30`)
	if got := AgentCooldownMinutes(); got != 30 {
		t.Errorf("AgentCooldownMinutes() = %d, want %d", got, 30)
	}
}

func TestTaskBudgetPct_Default(t *testing.T) {
	declareContext(t, "")
	if got := TaskBudgetPct(); got != cfgAgent.TaskBudgetPct {
		t.Errorf("TaskBudgetPct() = %v, want %v", got, cfgAgent.TaskBudgetPct)
	}
}

func TestTaskBudgetPct_Custom(t *testing.T) {
	declareContext(t, `task_budget_pct: 0.5`)
	if got := TaskBudgetPct(); got != 0.5 {
		t.Errorf("TaskBudgetPct() = %v, want %v", got, 0.5)
	}
}

func TestConventionBudgetPct_Default(t *testing.T) {
	declareContext(t, "")
	if got := ConventionBudgetPct(); got != cfgAgent.ConventionBudgetPct {
		t.Errorf("ConventionBudgetPct() = %v, want %v", got, cfgAgent.ConventionBudgetPct)
	}
}

func TestConventionBudgetPct_Custom(t *testing.T) {
	declareContext(t, `convention_budget_pct: 0.25`)
	if got := ConventionBudgetPct(); got != 0.25 {
		t.Errorf("ConventionBudgetPct() = %v, want %v", got, 0.25)
	}
}

func TestTitleSlugMaxLen_Default(t *testing.T) {
	declareContext(t, "")
	if got := TitleSlugMaxLen(); got != cfgJournal.TitleSlugMaxLen {
		t.Errorf("TitleSlugMaxLen() = %d, want %d", got, cfgJournal.TitleSlugMaxLen)
	}
}

func TestTitleSlugMaxLen_Custom(t *testing.T) {
	declareContext(t, `title_slug_max_len: 80`)
	if got := TitleSlugMaxLen(); got != 80 {
		t.Errorf("TitleSlugMaxLen() = %d, want %d", got, 80)
	}
}

func TestRecallListLimit_Default(t *testing.T) {
	declareContext(t, "")
	if got := RecallListLimit(); got != cfgJournal.DefaultRecallListLimit {
		t.Errorf("RecallListLimit() = %d, want %d",
			got, cfgJournal.DefaultRecallListLimit)
	}
}

func TestRecallListLimit_Custom(t *testing.T) {
	declareContext(t, `recall_list_limit: 50`)
	if got := RecallListLimit(); got != 50 {
		t.Errorf("RecallListLimit() = %d, want %d", got, 50)
	}
}

// TestAutoPruneDays_NegativeFallsBack guards the data-loss fix: a
// negative day count must fall back to the default, never reach
// health.AutoPrune (where it would move the cutoff into the future and
// delete every state file).
func TestAutoPruneDays_NegativeFallsBack(t *testing.T) {
	declareContext(t, `auto_prune_days: -7`)
	if got := AutoPruneDays(); got != cfgLoadgate.AutoPruneStaleDays {
		t.Errorf("AutoPruneDays() = %d, want default %d",
			got, cfgLoadgate.AutoPruneStaleDays)
	}
}

// TestTitleSlugMaxLen_NegativeFallsBack guards the panic fix: a
// negative cap must fall back to the default, never reach the slice
// bound in slug.FromTitle.
func TestTitleSlugMaxLen_NegativeFallsBack(t *testing.T) {
	declareContext(t, `title_slug_max_len: -5`)
	if got := TitleSlugMaxLen(); got != cfgJournal.TitleSlugMaxLen {
		t.Errorf("TitleSlugMaxLen() = %d, want default %d",
			got, cfgJournal.TitleSlugMaxLen)
	}
}

// TestRecallListLimit_NegativeFallsBack: a non-positive limit is unset.
func TestRecallListLimit_NegativeFallsBack(t *testing.T) {
	declareContext(t, `recall_list_limit: -1`)
	if got := RecallListLimit(); got != cfgJournal.DefaultRecallListLimit {
		t.Errorf("RecallListLimit() = %d, want default %d",
			got, cfgJournal.DefaultRecallListLimit)
	}
}

// TestAgentCooldownMinutes_ZeroDisables: an explicit 0 is distinct from
// unset and disables the cooldown (parity with --cooldown 0), rather
// than being coerced back to the default.
func TestAgentCooldownMinutes_ZeroDisables(t *testing.T) {
	declareContext(t, `agent_cooldown_minutes: 0`)
	if got := AgentCooldownMinutes(); got != 0 {
		t.Errorf("AgentCooldownMinutes() = %d, want 0 (disabled)", got)
	}
}

// TestTaskBudgetPct_ZeroAllowed: an explicit 0 zeroes the section
// rather than falling back to the default.
func TestTaskBudgetPct_ZeroAllowed(t *testing.T) {
	declareContext(t, `task_budget_pct: 0.0`)
	if got := TaskBudgetPct(); got != 0 {
		t.Errorf("TaskBudgetPct() = %v, want 0", got)
	}
}

// TestTaskBudgetPct_ClampsOutOfRange: values outside [0,1] are coerced.
func TestTaskBudgetPct_ClampsOutOfRange(t *testing.T) {
	declareContext(t, `task_budget_pct: 1.5`)
	if got := TaskBudgetPct(); got != 1 {
		t.Errorf("TaskBudgetPct() = %v, want 1 (clamped)", got)
	}
	declareContext(t, `task_budget_pct: -0.5`)
	if got := TaskBudgetPct(); got != 0 {
		t.Errorf("TaskBudgetPct() = %v, want 0 (clamped)", got)
	}
}

// TestConventionBudgetPct_ZeroAllowed mirrors the task-pct case.
func TestConventionBudgetPct_ZeroAllowed(t *testing.T) {
	declareContext(t, `convention_budget_pct: 0.0`)
	if got := ConventionBudgetPct(); got != 0 {
		t.Errorf("ConventionBudgetPct() = %v, want 0", got)
	}
}
