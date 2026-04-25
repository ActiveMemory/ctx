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

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/env"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
)

// declareContext sets up a tempDir layout with a .context/ directory
// and a .ctxrc at the project root (the parent of CTX_DIR), declares
// CTX_DIR via t.Setenv, and resets the rc singleton. The helper
// matches the single-source-anchor resolution model
// (spec: specs/single-source-context-anchor.md): .ctxrc is read from
// filepath.Dir(ContextDir())/.ctxrc, not CWD.
//
// Parameters:
//   - t: test handle for Setenv/TempDir/Cleanup wiring.
//   - content: YAML body to write into .ctxrc; empty for "no file".
//
// Returns:
//   - string: absolute path of the declared .context/ directory.
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
	t.Setenv(env.CtxDir, ctxDir)
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

// TestGetRC_NoFile: no CTX_DIR declared and no .ctxrc anywhere →
// defaults apply.
func TestGetRC_NoFile(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	// Ensure no env leak from other tests.
	t.Setenv(env.CtxDir, "")
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

// TestGetRC_WithFile: CTX_DIR declared, .ctxrc adjacent → values
// picked up.
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

// TestContextDir_RejectsUnset: CTX_DIR unset → ErrDirNotDeclared.
func TestContextDir_RejectsUnset(t *testing.T) {
	t.Setenv(env.CtxDir, "")
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if !errors.Is(err, errCtx.ErrDirNotDeclared) {
		t.Errorf("ContextDir() err = %v, want ErrDirNotDeclared", err)
	}
	if got != "" {
		t.Errorf("ContextDir() = %q, want \"\"", got)
	}
}

// TestContextDir_RejectsEmpty: CTX_DIR set to empty string is
// treated as unset. Spec contract: declared-or-not, no
// in-between.
func TestContextDir_RejectsEmpty(t *testing.T) {
	t.Setenv(env.CtxDir, "")
	Reset()
	t.Cleanup(Reset)

	_, err := ContextDir()
	if !errors.Is(err, errCtx.ErrDirNotDeclared) {
		t.Errorf("ContextDir() err = %v, want ErrDirNotDeclared", err)
	}
}

// TestContextDir_RejectsRelative_DotContext: critical regression
// guard against silent cwd-dependency. Without IsAbs check,
// CTX_DIR=.context would be cwd-absolutized via filepath.Abs and
// pass the basename guard, defeating the resolver.
func TestContextDir_RejectsRelative_DotContext(t *testing.T) {
	t.Setenv(env.CtxDir, ".context")
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if !errors.Is(err, errCtx.ErrRelativeNotAllowed) {
		t.Errorf("ContextDir() err = %v, want ErrRelativeNotAllowed", err)
	}
	if got != "" {
		t.Errorf("ContextDir() = %q, want \"\"", got)
	}
}

// TestContextDir_RejectsRelative_DotSlashContext: another shape of
// relative path, same expected error.
func TestContextDir_RejectsRelative_DotSlashContext(t *testing.T) {
	t.Setenv(env.CtxDir, "./.context")
	Reset()
	t.Cleanup(Reset)

	_, err := ContextDir()
	if !errors.Is(err, errCtx.ErrRelativeNotAllowed) {
		t.Errorf("ContextDir() err = %v, want ErrRelativeNotAllowed", err)
	}
}

// TestContextDir_RejectsRelative_DotDot: dot-dot relative path
// also rejected.
func TestContextDir_RejectsRelative_DotDot(t *testing.T) {
	t.Setenv(env.CtxDir, "../foo/.context")
	Reset()
	t.Cleanup(Reset)

	_, err := ContextDir()
	if !errors.Is(err, errCtx.ErrRelativeNotAllowed) {
		t.Errorf("ContextDir() err = %v, want ErrRelativeNotAllowed", err)
	}
}

// TestContextDir_RejectsNonCanonicalBasename: catches the common
// `export CTX_DIR=$(pwd)` footgun on first use rather than
// letting init deposit canonical files in the project root.
func TestContextDir_RejectsNonCanonicalBasename(t *testing.T) {
	t.Setenv(env.CtxDir, "/tmp/notdotcontext")
	Reset()
	t.Cleanup(Reset)

	_, err := ContextDir()
	if !errors.Is(err, errCtx.ErrNonCanonicalBasename) {
		t.Errorf("ContextDir() err = %v, want ErrNonCanonicalBasename", err)
	}
	if err != nil && !contains(err.Error(), "notdotcontext") {
		t.Errorf("err message %q should include offending basename", err.Error())
	}
}

// TestContextDir_RejectsRoot: filepath.Base("/") returns "/", not
// ".context", so root path is rejected by the basename guard.
func TestContextDir_RejectsRoot(t *testing.T) {
	t.Setenv(env.CtxDir, "/")
	Reset()
	t.Cleanup(Reset)

	_, err := ContextDir()
	if !errors.Is(err, errCtx.ErrNonCanonicalBasename) {
		t.Errorf("ContextDir() err = %v, want ErrNonCanonicalBasename", err)
	}
}

// TestContextDir_AcceptsCanonical: canonical absolute `.context`
// path is the happy path.
func TestContextDir_AcceptsCanonical(t *testing.T) {
	t.Setenv(env.CtxDir, "/tmp/.context")
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if err != nil {
		t.Fatalf("ContextDir() err = %v, want nil", err)
	}
	if got != "/tmp/.context" {
		t.Errorf("ContextDir() = %q, want %q", got, "/tmp/.context")
	}
}

// TestContextDir_NormalizesTrailingSlash: filepath.Clean strips
// trailing slash; basename guard still passes.
func TestContextDir_NormalizesTrailingSlash(t *testing.T) {
	t.Setenv(env.CtxDir, "/tmp/.context/")
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if err != nil {
		t.Fatalf("ContextDir() err = %v, want nil", err)
	}
	if got != "/tmp/.context" {
		t.Errorf("ContextDir() = %q, want %q", got, "/tmp/.context")
	}
}

// TestContextDir_NormalizesDotSegments: filepath.Clean
// canonicalizes dot segments.
func TestContextDir_NormalizesDotSegments(t *testing.T) {
	t.Setenv(env.CtxDir, "/tmp/./.context")
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if err != nil {
		t.Fatalf("ContextDir() err = %v, want nil", err)
	}
	if got != "/tmp/.context" {
		t.Errorf("ContextDir() = %q, want %q", got, "/tmp/.context")
	}
}

// TestContextDir_AcceptsSymlinkNamedDotContext: a symlink whose
// basename is `.context` (regardless of where it points) passes
// the basename guard. The resolver checks the *declared* name,
// not the symlink target name.
func TestContextDir_AcceptsSymlinkNamedDotContext(t *testing.T) {
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "actual-target")
	if err := os.MkdirAll(target, 0700); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	link := filepath.Join(tempDir, dir.Context)
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}
	t.Setenv(env.CtxDir, link)
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if err != nil {
		t.Fatalf("ContextDir() err = %v, want nil", err)
	}
	if got != link {
		t.Errorf("ContextDir() = %q, want %q (declared symlink path)", got, link)
	}
}

// contains is a small helper for substring checks in error
// messages. Avoids pulling strings.Contains everywhere.
func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

// TestContextDir_Unset: no env declaration → errCtx.ErrDirNotDeclared.
// Under the single-source-anchor model, unset is a valid signal used by
// exempt commands and rc.RequireContextDir's error path.
func TestContextDir_Unset(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)
	t.Setenv(env.CtxDir, "")
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if err == nil {
		t.Errorf("ContextDir() err = nil, want errCtx.ErrDirNotDeclared")
	}
	if got != "" {
		t.Errorf("ContextDir() = %q, want \"\" (unset)", got)
	}
}

// TestContextDir_EnvOnly: CTX_DIR env set with canonical absolute
// `.context` path → resolves to that path.
func TestContextDir_EnvOnly(t *testing.T) {
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, dir.Context)
	_ = os.MkdirAll(target, 0700)
	t.Setenv(env.CtxDir, target)
	Reset()
	t.Cleanup(Reset)

	got, err := ContextDir()
	if err != nil {
		t.Fatalf("ContextDir() err = %v, want nil", err)
	}
	if !filepath.IsAbs(got) {
		t.Errorf("ContextDir() = %q, want absolute path", got)
	}
	gotResolved, _ := filepath.EvalSymlinks(got)
	wantResolved, _ := filepath.EvalSymlinks(target)
	if gotResolved != wantResolved {
		t.Errorf("ContextDir() = %q, want %q (env)", gotResolved, wantResolved)
	}
}

// TestRequireContextDir_Declared: a declared CTX_DIR yields the
// path and no error.
func TestRequireContextDir_Declared(t *testing.T) {
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

// TestRequireContextDir_Undeclared: no override, no env → error
// with a tailored, non-empty message.
func TestRequireContextDir_Undeclared(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)
	t.Setenv(env.CtxDir, "")
	Reset()
	t.Cleanup(Reset)

	got, err := RequireContextDir()
	if err == nil {
		t.Fatalf("RequireContextDir() err = nil, want non-nil")
	}
	if got != "" {
		t.Errorf("RequireContextDir() path = %q, want \"\" on error", got)
	}
	if msg := err.Error(); msg == "" {
		t.Error("RequireContextDir() returned empty error message")
	}
}

// TestScanCandidates_NoMatches: empty tree → empty slice.
func TestScanCandidates_NoMatches(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	got := ScanCandidates(tempDir)
	if len(got) != 0 {
		t.Errorf("ScanCandidates() = %v, want []", got)
	}
}

// TestScanCandidates_SelfMatch: .context/ exists at start dir →
// one candidate, same path.
func TestScanCandidates_SelfMatch(t *testing.T) {
	tempDir := t.TempDir()
	ctxPath := filepath.Join(tempDir, dir.Context)
	_ = os.MkdirAll(ctxPath, 0700)

	got := ScanCandidates(tempDir)
	if len(got) != 1 {
		t.Fatalf("ScanCandidates() len = %d, want 1", len(got))
	}

	wantResolved, _ := filepath.EvalSymlinks(ctxPath)
	gotResolved, _ := filepath.EvalSymlinks(got[0])
	if gotResolved != wantResolved {
		t.Errorf("ScanCandidates()[0] = %q, want %q", gotResolved, wantResolved)
	}
}

// TestScanCandidates_ManyAncestors: nested .context/ dirs upward
// are all returned, innermost first.
func TestScanCandidates_ManyAncestors(t *testing.T) {
	tempDir := t.TempDir()
	inner := filepath.Join(tempDir, "inner", "deep")
	innerCtx := filepath.Join(tempDir, "inner", dir.Context)
	outerCtx := filepath.Join(tempDir, dir.Context)

	for _, d := range []string{inner, innerCtx, outerCtx} {
		if mkErr := os.MkdirAll(d, 0700); mkErr != nil {
			t.Fatalf("mkdir %s: %v", d, mkErr)
		}
	}

	got := ScanCandidates(inner)
	if len(got) < 2 {
		t.Fatalf("ScanCandidates() len = %d, want >= 2", len(got))
	}

	// Innermost first: the first candidate must be in the parent of
	// the start dir (i.e., inner/.context).
	innerResolved, _ := filepath.EvalSymlinks(innerCtx)
	gotInner, _ := filepath.EvalSymlinks(got[0])
	if gotInner != innerResolved {
		t.Errorf("ScanCandidates()[0] = %q, want %q (innermost)", gotInner, innerResolved)
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
