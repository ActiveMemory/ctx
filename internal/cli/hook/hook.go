//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/event"
	"github.com/ActiveMemory/ctx/internal/cli/message"
	"github.com/ActiveMemory/ctx/internal/cli/notify"
	"github.com/ActiveMemory/ctx/internal/cli/pause"
	"github.com/ActiveMemory/ctx/internal/cli/resume"
	"github.com/ActiveMemory/ctx/internal/cli/unknown"
	cliCfg "github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx hook" parent command.
//
// Consolidates hook-related user-facing commands: message,
// notify, pause, resume, and event.
//
// An unknown `ctx hook <verb>` fails loud (verbatim relay box + best-
// effort relay event + non-zero exit) instead of printing help at
// exit 0 — a silent failure for the skills and loop scripts that call
// `ctx hook` verbs by name when one drifts out of the binary. A bare
// `ctx hook` still prints help and exits 0.
//
// The AnnotationSkipInit is required: RootCmd's PersistentPreRunE
// exempts grouping commands that have no Run/RunE, but this group now
// has one. Without the annotation, bare `ctx hook` and an unknown verb
// would newly require an initialized context + git tree. The annotation
// is evaluated against the target command, so it exempts only the
// group-level invocation; the real subcommands keep their own
// preconditions. See specs/unknown-subcommand-relay-generalization.md.
//
// Returns:
//   - *cobra.Command: Parent command with hook subcommands
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyHook)
	c := &cobra.Command{
		Use:     cmd.UseHook,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyHook),
		Annotations: map[string]string{
			cliCfg.AnnotationSkipInit: cliCfg.AnnotationTrue,
		},
		RunE: unknown.HandlerFor(unknown.HookConfig),
	}
	c.AddCommand(
		event.Cmd(),
		message.Cmd(),
		notify.Cmd(),
		pause.Cmd(),
		resume.Cmd(),
	)
	return c
}
