//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/cmd/dismiss"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/cmd/list"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/cmd/show"
)

// Strings carries the English user-facing text for the audit
// command and its subcommands. ctxctl supplies these from its
// own Go constants; the logic packages hold no copy of their
// own.
type Strings struct {
	// Use is the cobra Use string for `audit`.
	Use string
	// Short is the one-line `audit` description.
	Short string
	// Long is the multi-line `audit` help text.
	Long string
	// Example is the `audit` example-usage block.
	Example string
	// List supplies the list subcommand's text.
	List list.Strings
	// Show supplies the show subcommand's text.
	Show show.Strings
	// Dismiss supplies the dismiss subcommand's text.
	Dismiss dismiss.Strings
}
