//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/philosophy"
	errCli "github.com/ActiveMemory/ctx/internal/err/cli"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/why"
)

// ShowDoc loads an embedded document by alias, strips
// MkDocs syntax, and prints it.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - alias: Document alias (manifesto, about, invariants)
//
// Returns:
//   - error: Non-nil if alias unknown or document load fails
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
