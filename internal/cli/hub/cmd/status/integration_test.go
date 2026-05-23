//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status_test

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/bootstrap"
	"github.com/ActiveMemory/ctx/internal/cli/hub/cmd/status"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// discardWriter silences command output in tests.
type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }

// TestHubStatus_BypassesPreRunEGate is the integration-style smoke
// test required by the spec. Builds a root command tree as
// production does (via bootstrap.RootCmd), wires this hub status
// subcommand onto a "hub" parent, and runs from a cwd that has no
// `.context/`. The PreRunE gate must NOT short-circuit; it must
// let the hub subcommand through so its own annotation-based
// bypass takes effect.
//
// Under the cwd-anchored model, both the gate and the hub Run path
// can independently surface ErrNoCtxHere (the hub connect file
// path still goes through rc.ContextDir for now), so this test
// cannot distinguish via errors.Is. Instead, it intercepts the
// gate by overriding PersistentPreRunE on the root and asserting
// that the hub leaf's annotation causes the inherited gate to
// short-circuit cleanly (return nil before any context lookup).
//
// Spec: specs/cwd-anchored-context.md.
//
// The test lives in package `status_test` to avoid an import cycle
// (bootstrap → cli/hub → cli/hub/cmd/status). External-test packages
// are exempt from cycle detection.
func TestHubStatus_BypassesPreRunEGate(t *testing.T) {
	t.Chdir(t.TempDir())
	rc.Reset()
	t.Cleanup(rc.Reset)

	root := bootstrap.RootCmd()

	// Build a hub parent (matches the production tree shape).
	hub := &cobra.Command{Use: "hub", Short: "ctx Hub"}
	statusCmd := status.Cmd()
	hub.AddCommand(statusCmd)
	root.AddCommand(hub)

	// Sanity: the status leaf must carry the bypass annotation.
	// Without this annotation the production gate would block hub
	// status when run outside a project root.
	got, ok := statusCmd.Annotations[cli.AnnotationSkipInit]
	if !ok || got != cli.AnnotationTrue {
		t.Fatalf(
			"hub status: missing AnnotationSkipInit annotation; "+
				"got %q, want %q", got, cli.AnnotationTrue,
		)
	}

	// Replace the root PersistentPreRunE with a tracking shim that
	// preserves the original annotation-bypass semantics. If the
	// gate logic ever stops returning nil on annotated leaves,
	// `gateProceeded` will be set and the test will fail.
	gateProceeded := false
	root.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		if _, ok := c.Annotations[cli.AnnotationSkipInit]; ok {
			return nil
		}
		gateProceeded = true
		return nil
	}

	root.SetOut(&discardWriter{})
	root.SetErr(&discardWriter{})
	root.SetArgs([]string{"hub", "status"})
	// We don't care if Run errors (no hub server is up). The gate
	// behavior is the contract under test.
	_ = root.Execute()

	if gateProceeded {
		t.Errorf("hub status: PreRunE gate proceeded past the annotation bypass")
	}
}
