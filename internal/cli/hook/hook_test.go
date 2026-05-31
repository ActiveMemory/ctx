//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	"github.com/ActiveMemory/ctx/internal/cli/hook"
	"github.com/ActiveMemory/ctx/internal/config/cli"
)

func TestMain(m *testing.M) {
	lookup.Init()
	m.Run()
}

// newRoot wraps a command in a minimal parent so the wrapped group
// reports HasParent()==true — the condition under which cobra lets an
// unmatched subcommand reach the group's RunE (rather than raising
// "unknown command" as it does only for the root).
func newRoot(child *cobra.Command) (*cobra.Command, *bytes.Buffer) {
	root := &cobra.Command{Use: "ctx"}
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.AddCommand(child)
	return root, &out
}

func TestHookUnknownSubcommandFailsLoud(t *testing.T) {
	root, out := newRoot(hook.Cmd())
	root.SetArgs([]string{"hook", "no-such-verb"})

	if err := root.Execute(); err == nil {
		t.Fatal("want non-nil error for an unknown `ctx hook` subcommand")
	}
	got := out.String()
	if !strings.Contains(got, "no-such-verb") {
		t.Errorf("want a relay box naming the verb; got:\n%s", got)
	}
	if !strings.Contains(got, "Unknown Hook Subcommand") {
		t.Errorf("want the hook-specific box title; got:\n%s", got)
	}
	// The whole point: no Long-help dump on the unknown path.
	if strings.Contains(got, "Available Commands:") {
		t.Errorf("unknown subcommand must not dump help; got:\n%s", got)
	}
}

func TestHookBareStillPrintsHelp(t *testing.T) {
	root, out := newRoot(hook.Cmd())
	root.SetArgs([]string{"hook"})

	if err := root.Execute(); err != nil {
		t.Fatalf("bare `ctx hook`: want nil error (help), got %v", err)
	}
	if out.Len() == 0 {
		t.Error("bare `ctx hook` should print help")
	}
}

// TestHookOptsIntoUnknownRelay guards that the hook group carries both
// the RunE (the relay) and the AnnotationSkipInit that keeps RootCmd's
// PersistentPreRunE from newly imposing context/git preconditions on a
// group that previously had no RunE.
func TestHookOptsIntoUnknownRelay(t *testing.T) {
	c := hook.Cmd()
	if c.RunE == nil {
		t.Error("hook.Cmd() must set a RunE (the unknown-subcommand relay)")
	}
	if _, ok := c.Annotations[cli.AnnotationSkipInit]; !ok {
		t.Error("hook.Cmd() must carry AnnotationSkipInit so the relay is reachable without an initialized context")
	}
}
