//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Command ctxctl is the ctx maintainer/contributor binary. It houses
// tooling that must never ship in the user-facing ctx binary; its
// first inhabitant is the out-of-band audit channel (audit
// list/show/dismiss + audit-relay).
//
// ctxctl is a separate Go module from ctx by design: ctx's go.mod
// does not require it, so ctx can never import ctxctl. ctxctl reuses
// ctx's internal/ packages via the repo-root go.work workspace.
//
// ctxctl owns its user-facing text as plain English Go constants
// (see text.go), outside ctx's YAML i18n: there is no French ctxctl.
//
// See specs/ctxctl-bootstrap.md and DECISIONS.md (2026-05-27).
package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/checkaudit"
)

// main is the entry point for the ctxctl maintainer binary.
func main() {
	lookup.Init()

	root := newRoot()
	if err := root.Execute(); err != nil {
		printErr(root, err)
		os.Exit(1)
	}
}

// newRoot assembles the ctxctl root command and its subtree (audit
// list/show/dismiss + audit-relay).
//
// Cobra auto-prints returned errors to stderr by default; ctxctl
// prints them itself via [printErr] so the audit channel's typed
// errors render in ctxctl's English voice. SilenceErrors makes
// printErr the sole printer (mirrors ctx's internal/write/err.With);
// SilenceUsage stays per-return in the audit Run functions.
//
// Returns:
//   - *cobra.Command: the configured ctxctl root command
func newRoot() *cobra.Command {
	root := &cobra.Command{
		Use:           rootUse,
		Short:         rootShort,
		Long:          rootLong,
		SilenceErrors: true,
	}
	root.AddCommand(audit.Cmd(auditStrings()))
	root.AddCommand(checkaudit.Cmd(relayStrings()))
	return root
}
