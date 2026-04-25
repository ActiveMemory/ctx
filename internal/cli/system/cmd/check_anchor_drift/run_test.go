//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_anchor_drift

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/env"
)

// runDrift invokes Run with both env vars set to the given values
// and returns whatever the hook printed to stdout.
func runDrift(t *testing.T, inherited, injected string) string {
	t.Helper()
	t.Setenv(env.CtxDirInherited, inherited)
	t.Setenv(env.CtxDir, injected)

	// Provide an empty stdin (no JSON envelope). ReadInput
	// gracefully returns a zero-value HookInput when the body is
	// empty, which the hook treats as "no session ID".
	stdinPath := filepath.Join(t.TempDir(), "stdin")
	if err := os.WriteFile(stdinPath, []byte("{}"), 0o600); err != nil {
		t.Fatalf("write stdin: %v", err)
	}
	stdin, err := os.Open(stdinPath)
	if err != nil {
		t.Fatalf("open stdin: %v", err)
	}
	t.Cleanup(func() { _ = stdin.Close() })

	c := &cobra.Command{}
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&out)
	if runErr := Run(c, stdin); runErr != nil {
		t.Fatalf("Run() err = %v, want nil", runErr)
	}
	return out.String()
}

// TestCheckAnchorDrift_Match: inherited and injected values equal
// after filepath.Clean → silent (no banner emitted).
func TestCheckAnchorDrift_Match(t *testing.T) {
	out := runDrift(t, "/project-a/.context", "/project-a/.context")
	if out != "" {
		t.Errorf("hook should be silent on match, got %q", out)
	}
}

// TestCheckAnchorDrift_MatchAfterClean: trailing slash on one
// side normalizes via filepath.Clean and comparison matches.
func TestCheckAnchorDrift_MatchAfterClean(t *testing.T) {
	out := runDrift(t, "/project-a/.context/", "/project-a/.context")
	if out != "" {
		t.Errorf("hook should be silent on match after Clean, got %q", out)
	}
}

// TestCheckAnchorDrift_Mismatch: inherited points at project A,
// injected points at project B → emit banner naming both.
func TestCheckAnchorDrift_Mismatch(t *testing.T) {
	out := runDrift(t, "/project-a/.context", "/project-b/.context")
	if out == "" {
		t.Fatal("hook should emit banner on mismatch")
	}
	if !strings.Contains(out, "/project-a/.context") {
		t.Errorf("banner should name inherited path, got %q", out)
	}
	if !strings.Contains(out, "/project-b/.context") {
		t.Errorf("banner should name injected path, got %q", out)
	}
	if !strings.Contains(out, "Anchor Drift") {
		t.Errorf("banner should carry the box title, got %q", out)
	}
}

// TestCheckAnchorDrift_InheritedEmpty: user has not run
// `ctx activate`; no shell-level declaration to drift from →
// silent regardless of injected value.
func TestCheckAnchorDrift_InheritedEmpty(t *testing.T) {
	out := runDrift(t, "", "/project-a/.context")
	if out != "" {
		t.Errorf("hook should be silent when inherited is empty, got %q", out)
	}
}

// TestCheckAnchorDrift_AcceptsNonCanonicalInherited: the hook is a
// diagnostic — it must accept any inherited value (including
// non-canonical) so it can describe reality, not impose policy.
// Verifies the hook bypasses rc.ContextDir's basename guard.
func TestCheckAnchorDrift_AcceptsNonCanonicalInherited(t *testing.T) {
	out := runDrift(t,
		"/some/random/path", "/project-a/.context",
	)
	if out == "" {
		t.Fatal("hook should emit banner on mismatch even with non-canonical inherited")
	}
	if !strings.Contains(out, "/some/random/path") {
		t.Errorf("banner should name inherited path verbatim, got %q", out)
	}
}

// TestCheckAnchorDrift_SymlinkEquivalent: paths that differ
// byte-for-byte but resolve to the same directory via a symlink
// (the canonical macOS case: `/tmp` → `/private/tmp`) must NOT
// trip the drift alarm. The smoke-test surfacing this case
// blocked step 9; without symlink resolution the banner fires
// every prompt for any session run from `/tmp/*` on macOS, and
// for any user with a symlinked workspace path elsewhere.
func TestCheckAnchorDrift_SymlinkEquivalent(t *testing.T) {
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "target", ".context")
	if err := os.MkdirAll(target, 0o700); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	link := filepath.Join(tempDir, "link")
	if err := os.Symlink(filepath.Join(tempDir, "target"), link); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}
	linkedContext := filepath.Join(link, ".context")

	// `target` and `linkedContext` differ as strings but resolve
	// to the same directory. Hook should be silent.
	out := runDrift(t, linkedContext, target)
	if out != "" {
		t.Errorf("hook should be silent for symlink-equivalent paths, got %q", out)
	}
}

// TestCheckAnchorDrift_SymlinkResolutionFails_FallsBackToString:
// when the inherited path can't be resolved (e.g. it points at a
// deleted directory), genuine drift should still fire. Defends
// against an over-eager symlink fix that would silently swallow
// real misalignment.
func TestCheckAnchorDrift_SymlinkResolutionFails_FallsBackToString(t *testing.T) {
	out := runDrift(t,
		"/definitely/does/not/exist/.context",
		"/project-b/.context",
	)
	if out == "" {
		t.Fatal("hook should still fire when inherited resolution fails")
	}
}
