//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gitmeta_test

import (
	"os/exec"
	"strings"
	"testing"

	cfgGitmeta "github.com/ActiveMemory/ctx/internal/config/gitmeta"
	"github.com/ActiveMemory/ctx/internal/gitmeta"
)

func TestResolveHead_CtxTaskCommitOverrideUsedVerbatim(t *testing.T) {
	t.Setenv(cfgGitmeta.EnvCtxTaskCommit, "deadbee")
	t.Setenv(cfgGitmeta.EnvGithubActions, "")
	t.Setenv(cfgGitmeta.EnvGithubSHA, "")

	ref, err := gitmeta.ResolveHead(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.SHA != "deadbee" {
		t.Errorf("SHA: want deadbee; got %q", ref.SHA)
	}
	// Branch falls through to git; tempdir has no git, so the
	// branch resolver returns "detached".
	if ref.Branch != cfgGitmeta.BranchDetached {
		t.Errorf("Branch: want %q; got %q", cfgGitmeta.BranchDetached, ref.Branch)
	}
}

func TestResolveHead_GithubShaTruncatedToShort(t *testing.T) {
	t.Setenv(cfgGitmeta.EnvCtxTaskCommit, "")
	t.Setenv(cfgGitmeta.EnvGithubActions, cfgGitmeta.GithubActionsTrue)
	t.Setenv(cfgGitmeta.EnvGithubSHA, "abcdef0123456789abcdef0123456789abcdef01")

	ref, err := gitmeta.ResolveHead(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := len(ref.SHA), cfgGitmeta.ShortLen; got != want {
		t.Errorf("SHA length: want %d; got %d (%q)", want, got, ref.SHA)
	}
	if ref.SHA != "abcdef0" {
		t.Errorf("SHA: want abcdef0; got %q", ref.SHA)
	}
}

func TestResolveHead_GithubShaIgnoredWithoutActionsFlag(t *testing.T) {
	t.Setenv(cfgGitmeta.EnvCtxTaskCommit, "")
	t.Setenv(cfgGitmeta.EnvGithubActions, "")
	t.Setenv(cfgGitmeta.EnvGithubSHA, "abcdef0123456789abcdef0123456789abcdef01")

	_, err := gitmeta.ResolveHead(t.TempDir())
	// No git in tempdir + no env override → resolution fails.
	// We don't care about the exact error text, only that the
	// function did NOT silently take the GITHUB_SHA value.
	if err == nil {
		t.Fatal("want error (no override + no git tree); got nil")
	}
}

func TestResolveHead_NoOverridesAndNoGitFails(t *testing.T) {
	t.Setenv(cfgGitmeta.EnvCtxTaskCommit, "")
	t.Setenv(cfgGitmeta.EnvGithubActions, "")
	t.Setenv(cfgGitmeta.EnvGithubSHA, "")

	_, err := gitmeta.ResolveHead(t.TempDir())
	if err == nil {
		t.Fatal("want error in dir without git; got nil")
	}
}

// TestResolveHead_RealRepoReturnsBranchName guards against
// the `git rev-parse --show-current` regression where the
// branch resolver was emitting the literal flag string
// "--show-current" because rev-parse echoes unknown args
// rather than erroring. The correct invocation is
// `git branch --show-current`.
func TestResolveHead_RealRepoReturnsBranchName(t *testing.T) {
	t.Setenv(cfgGitmeta.EnvCtxTaskCommit, "")
	t.Setenv(cfgGitmeta.EnvGithubActions, "")
	t.Setenv(cfgGitmeta.EnvGithubSHA, "")

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not on PATH")
	}

	tmp := t.TempDir()
	runGit := func(args ...string) {
		t.Helper()
		//nolint:gosec // G204: test fixture, args are hardcoded literals
		cmd := exec.Command("git", append([]string{"-C", tmp}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	runGit("init", "--initial-branch=trunk")
	runGit("config", "user.email", "test@example.com")
	runGit("config", "user.name", "Test")
	runGit("commit", "--allow-empty", "-m", "init")

	ref, err := gitmeta.ResolveHead(tmp)
	if err != nil {
		t.Fatalf("ResolveHead: %v", err)
	}
	if ref.Branch != "trunk" {
		t.Errorf("Branch: want %q; got %q", "trunk", ref.Branch)
	}
	if strings.Contains(ref.Branch, "--") {
		t.Errorf("Branch leaks a flag literal: %q", ref.Branch)
	}
}
