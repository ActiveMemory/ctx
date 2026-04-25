//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/bootstrap"
	"github.com/ActiveMemory/ctx/internal/cli/hub/cmd/status"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/env"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// discardWriter silences command output in tests.
type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }

// TestHubStatus_BypassesPreRunEGate is the integration-style smoke
// test required by the spec. Builds a root command tree as
// production does (via bootstrap.RootCmd), wires this hub status
// subcommand onto a "hub" parent, and runs with CTX_DIR pointing at
// a deliberately-non-existent path. The PreRunE gate must NOT
// short-circuit with ErrDirNotDeclared.
//
// Without this guard, a future refactor that breaks PreRunE's
// annotation handling could leave the annotation in place but
// regress the actual bypass behavior.
//
// Spec: specs/single-source-context-anchor.md.
//
// The test lives in package `status_test` to avoid an import cycle
// (bootstrap → cli/hub → cli/hub/cmd/status). External-test packages
// are exempt from cycle detection.
func TestHubStatus_BypassesPreRunEGate(t *testing.T) {
	// Wire CTX_DIR to a deliberately-non-existent shape-valid path
	// so RequireContextDir would fail loud if PreRunE actually ran.
	t.Setenv(env.CtxDir, filepath.Join(t.TempDir(), "absent", dir.Context))
	rc.Reset()
	t.Cleanup(rc.Reset)

	root := bootstrap.RootCmd()

	// Build a hub parent (matches the production tree shape).
	hub := &cobra.Command{
		Use:   "hub",
		Short: "ctx Hub",
	}
	hub.AddCommand(status.Cmd())
	root.AddCommand(hub)

	root.SetOut(&discardWriter{})
	root.SetErr(&discardWriter{})
	root.SetArgs([]string{"hub", "status"})

	err := root.Execute()
	// Server is not running so coreStatus.Run will return its own
	// connect error — that's fine. The contract: the error must
	// NOT be the gate's "context dir not declared" sentinel.
	if errors.Is(err, errCtx.ErrDirNotDeclared) {
		t.Errorf("hub status: PreRunE gate short-circuited with ErrDirNotDeclared (annotation bypass broken)")
	}
}
