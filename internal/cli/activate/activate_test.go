//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package activate_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/activate"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/env"
)

// runActivate invokes `ctx activate` with the given args and returns
// (stdout, stderr, error) as separate buffers. Stream separation is
// load-bearing: the eval-bindable shell content goes to stdout, the
// human-readable advisories ("ctx activated at:", "ctx: also
// visible upward:") go to stderr. Tests that conflate the two miss
// regressions where an advisory leaks into the eval stream.
//
// The command inherits the test process's env; use t.Setenv /
// t.Chdir to scope state.
func runActivate(t *testing.T, args []string) (stdout, stderr string, err error) {
	t.Helper()
	c := activate.Cmd()
	c.SetArgs(args)
	var so, se bytes.Buffer
	c.SetOut(&so)
	c.SetErr(&se)
	err = c.Execute()
	return so.String(), se.String(), err
}

// TestActivate_NoArgs_NoCandidates: cwd with no .context/ anywhere →
// NoCandidates error, no stdout emit, no advisory either.
func TestActivate_NoArgs_NoCandidates(t *testing.T) {
	t.Setenv(env.CtxDir, "")
	t.Chdir(t.TempDir())

	stdout, _, err := runActivate(t, nil)
	if err == nil {
		t.Fatalf("expected NoCandidates error, got nil (stdout=%q)", stdout)
	}
	if stdout != "" {
		t.Errorf("stdout must be empty on error path: %q", stdout)
	}
}

// TestActivate_NoArgs_OneCandidate: exactly one .context/ upward →
// stdout carries the export, stderr carries the
// `ctx: activated at:<path>` advisory.
func TestActivate_NoArgs_OneCandidate(t *testing.T) {
	t.Setenv(env.CtxDir, "")

	projectRoot := t.TempDir()
	ctxPath := filepath.Join(projectRoot, dir.Context)
	if err := os.MkdirAll(ctxPath, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	t.Chdir(projectRoot)
	t.Setenv("SHELL", "/bin/bash")

	stdout, stderr, err := runActivate(t, nil)
	if err != nil {
		t.Fatalf("expected success, got err=%v", err)
	}
	if !strings.HasPrefix(stdout, "export CTX_DIR=") {
		t.Errorf("stdout must start with export, got %q", stdout)
	}
	if !strings.Contains(stdout, ctxPath) {
		t.Errorf("stdout missing path %q: %q", ctxPath, stdout)
	}
	// Activated-at advisory always announces the bind, even in
	// the single-candidate case.
	wantActivated := "ctx: activated at: " + ctxPath
	if !strings.Contains(stderr, wantActivated) {
		t.Errorf("stderr missing activated-at advisory %q: %q",
			wantActivated, stderr)
	}
}

// TestActivate_ErrorPath_StdoutEmpty guards the eval-recursion
// trap surfaced by the smoke test: if any error path lets cobra
// print Usage / Flags / Examples to stdout, `eval "$(ctx
// activate)"` captures the Examples block (which literally
// contains `eval "$(ctx activate)"`) and re-executes activate,
// looping until the captured text mangles past the parser.
//
// Stdout MUST stay empty on every error path. Stderr can carry
// the human-readable error; the eval shell never sees stderr.
//
// Uses the no-candidates case (zero `.context/` visible upward)
// since multi-candidate is no longer an error case under the
// innermost-wins policy.
func TestActivate_ErrorPath_StdoutEmpty(t *testing.T) {
	t.Setenv(env.CtxDir, "")
	t.Chdir(t.TempDir())

	stdout, stderr, err := runActivate(t, nil)
	if err == nil {
		t.Fatalf("expected NoCandidates error, got nil")
	}
	if stdout != "" {
		t.Errorf("stdout must be empty on error path, got %q", stdout)
	}
	if !strings.Contains(stderr, "no .context/ directory found") {
		t.Errorf("stderr should describe the error, got %q", stderr)
	}
}

// TestActivate_NoArgs_ManyCandidates: two `.context/` dirs on the
// upward path → innermost wins on stdout (eval-bindable),
// stderr carries both the `ctx activated at:` line and one
// `ctx: also visible upward:` line per other candidate. Matches
// git/make innermost-project semantics.
//
// The split-stream assertion is load-bearing: putting any
// advisory on stdout (the eval-captured stream) makes it
// invisible to anyone running `eval "$(ctx activate)"`.
func TestActivate_NoArgs_ManyCandidates(t *testing.T) {
	t.Setenv(env.CtxDir, "")
	t.Setenv("SHELL", "/bin/bash")

	tempDir := t.TempDir()
	outerCtx := filepath.Join(tempDir, dir.Context)
	innerDir := filepath.Join(tempDir, "inner")
	innerCtx := filepath.Join(innerDir, dir.Context)
	startDir := filepath.Join(innerDir, "deep")

	for _, d := range []string{outerCtx, innerCtx, startDir} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}
	t.Chdir(startDir)

	stdout, stderr, err := runActivate(t, nil)
	if err != nil {
		t.Fatalf("expected success (innermost wins), got err=%v", err)
	}

	// stdout: only the export line for the innermost candidate.
	if !strings.HasPrefix(stdout, "export CTX_DIR=") {
		t.Errorf("stdout must start with export, got %q", stdout)
	}
	if !strings.Contains(stdout, innerCtx) {
		t.Errorf("export should bind the inner candidate %q: %q",
			innerCtx, stdout)
	}
	if strings.Contains(stdout, "also visible") ||
		strings.Contains(stdout, "activated at") {
		t.Errorf("stdout must NOT carry advisories (eval invisibility): %q",
			stdout)
	}

	// stderr: activated-at line for the inner, also-visible line for
	// the outer.
	wantActivated := "ctx: activated at: " + innerCtx
	if !strings.Contains(stderr, wantActivated) {
		t.Errorf("stderr missing %q: %q", wantActivated, stderr)
	}
	wantAdvisory := "ctx: also visible upward: " + outerCtx
	if !strings.Contains(stderr, wantAdvisory) {
		t.Errorf("stderr missing %q: %q", wantAdvisory, stderr)
	}
}

// TestActivate_RejectsArgs guards the spec contract: `ctx activate
// <path>` is removed under the single-source-anchor model. Any
// positional argument must be rejected (either as cobra's
// "accepts 0 arg(s)" or "unknown command", whichever cobra picks
// for the literal value) and emit nothing on stdout.
func TestActivate_RejectsArgs(t *testing.T) {
	t.Setenv(env.CtxDir, "")
	t.Chdir(t.TempDir())

	stdout, _, err := runActivate(t, []string{"some-explicit-path"})
	if err == nil {
		t.Fatalf("expected cobra args rejection, got nil (stdout=%q)", stdout)
	}
	if strings.Contains(stdout, "export CTX_DIR") {
		t.Errorf("stdout should not contain export on error: %q", stdout)
	}
}

// TestActivate_StaleReplacementComment: parent shell has a stale
// CTX_DIR pointing at a different project; activate emits a
// `# ctx: replacing stale CTX_DIR=<old>` comment before the export
// so the user can see the change in `eval` output.
func TestActivate_StaleReplacementComment(t *testing.T) {
	stale := filepath.Join(t.TempDir(), "old", dir.Context)
	if err := os.MkdirAll(stale, 0700); err != nil {
		t.Fatalf("mkdir stale: %v", err)
	}
	t.Setenv(env.CtxDir, stale)

	projectRoot := t.TempDir()
	ctxPath := filepath.Join(projectRoot, dir.Context)
	if err := os.MkdirAll(ctxPath, 0700); err != nil {
		t.Fatalf("mkdir new: %v", err)
	}
	t.Chdir(projectRoot)
	t.Setenv("SHELL", "/bin/bash")

	stdout, _, err := runActivate(t, nil)
	if err != nil {
		t.Fatalf("expected success, got err=%v", err)
	}
	wantPrefix := fmt.Sprintf("# ctx: replacing stale %s=%s\n",
		env.CtxDir, stale)
	if !strings.HasPrefix(stdout, wantPrefix) {
		t.Errorf("stdout missing stale-replacement comment.\n got: %q\nwant prefix: %q",
			stdout, wantPrefix)
	}
	if !strings.Contains(stdout, "export CTX_DIR=") {
		t.Errorf("stdout missing export: %q", stdout)
	}
	if !strings.Contains(stdout, ctxPath) {
		t.Errorf("stdout missing new path %q: %q", ctxPath, stdout)
	}
}

// TestActivate_NoStaleCommentOnFirstActivate: when CTX_DIR is unset
// or matches the resolved value, the comment is suppressed.
func TestActivate_NoStaleCommentOnFirstActivate(t *testing.T) {
	t.Setenv(env.CtxDir, "")

	projectRoot := t.TempDir()
	ctxPath := filepath.Join(projectRoot, dir.Context)
	if err := os.MkdirAll(ctxPath, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	t.Chdir(projectRoot)
	t.Setenv("SHELL", "/bin/bash")

	stdout, _, err := runActivate(t, nil)
	if err != nil {
		t.Fatalf("expected success, got err=%v", err)
	}
	if strings.Contains(stdout, "replacing stale") {
		t.Errorf("stdout should not contain stale comment: %q", stdout)
	}
}

// TestActivate_ShellFlag: --shell zsh uses POSIX export syntax
// (same output shape as bash; flag is just a dispatch key).
func TestActivate_ShellFlag(t *testing.T) {
	t.Setenv(env.CtxDir, "")

	projectRoot := t.TempDir()
	ctxPath := filepath.Join(projectRoot, dir.Context)
	if err := os.MkdirAll(ctxPath, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	t.Chdir(projectRoot)

	stdout, _, err := runActivate(t, []string{"--shell", "zsh"})
	if err != nil {
		t.Fatalf("expected success, got err=%v", err)
	}
	if !strings.HasPrefix(stdout, "export CTX_DIR=") {
		t.Errorf("expected export prefix, got %q", stdout)
	}
	if !strings.HasSuffix(strings.TrimSpace(stdout), "'") {
		t.Errorf("expected trailing single quote (shell quoting), got %q", stdout)
	}
}
