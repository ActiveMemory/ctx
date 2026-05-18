//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gitmeta_test

import (
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
