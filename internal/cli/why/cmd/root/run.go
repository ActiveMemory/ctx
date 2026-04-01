//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/philosophy"
	errCli "github.com/ActiveMemory/ctx/internal/err/cli"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/why"
)

// Run dispatches to the interactive menu or direct document display.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: Command arguments; optional args[0] is the document alias
//
// Returns:
//   - error: Non-nil if the document is not found or input is invalid
func Run(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return ShowDoc(cmd, args[0])
	}
	return showMenu(cmd)
}

// ShowDoc loads an embedded document by alias, strips
// MkDocs syntax, and prints it.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - alias: Document alias (manifesto, about, or invariants)
//
// Returns:
//   - error: Non-nil if the alias is unknown or the document fails to load
func ShowDoc(cmd *cobra.Command, alias string) error {
	name, ok := DocAliases[alias]
	if !ok {
		return errCli.UnknownDocument(alias)
	}

	content, loadErr := philosophy.WhyDoc(name)
	if loadErr != nil {
		return errFs.FileRead(name, loadErr)
	}

	cleaned := StripMkDocs(string(content))
	why.Content(cmd, cleaned)

	return nil
}
