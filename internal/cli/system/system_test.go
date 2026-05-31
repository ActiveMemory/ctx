//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/cli/system"
	embedCmd "github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

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

func TestSystemUnknownSubcommandFailsLoud(t *testing.T) {
	root, out := newRoot(system.Cmd())
	root.SetArgs([]string{"system", "no-such-verb"})

	if err := root.Execute(); err == nil {
		t.Fatal("want non-nil error for an unknown `ctx system` subcommand")
	}
	got := out.String()
	if !strings.Contains(got, "no-such-verb") {
		t.Errorf("want a relay box naming the verb; got:\n%s", got)
	}
	// The whole point: no ~51-line Long-help dump on the unknown path.
	if strings.Contains(got, "Available Commands:") {
		t.Errorf("unknown subcommand must not dump help; got:\n%s", got)
	}
}

func TestSystemBareStillPrintsHelp(t *testing.T) {
	root, out := newRoot(system.Cmd())
	root.SetArgs([]string{"system"})

	if err := root.Execute(); err != nil {
		t.Fatalf("bare `ctx system`: want nil error (help), got %v", err)
	}
	if out.Len() == 0 {
		t.Error("bare `ctx system` should print help")
	}
}

// TestParentCmdScopeUnchanged documents that the unknown-subcommand
// relay is an explicit per-group opt-in (system, hook), not a default:
// the shared parent.Cmd still produces a group with no RunE, so groups
// that do not opt in keep cobra's default (help + exit 0) on an unknown
// subcommand. If someone moves the relay into parent.Cmd, this fails.
func TestParentCmdScopeUnchanged(t *testing.T) {
	if system.Cmd().RunE == nil {
		t.Fatal("system.Cmd() must set a RunE (the unknown-subcommand fix)")
	}

	grp := parent.Cmd(embedCmd.DescKeySystem, "grouplike",
		&cobra.Command{
			Use:  "real",
			RunE: func(*cobra.Command, []string) error { return nil },
		},
	)
	if grp.RunE != nil {
		t.Error("parent.Cmd must not set RunE; fix is scoped to system")
	}

	root, _ := newRoot(grp)
	root.SetArgs([]string{"grouplike", "bogus"})
	if err := root.Execute(); err != nil {
		t.Errorf(
			"a parent.Cmd group should exit 0 (help) on an unknown "+
				"subcommand, got %v", err,
		)
	}
}
