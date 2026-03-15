//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	fserr "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/validate"
	"github.com/ActiveMemory/ctx/internal/write"
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

// showMenu presents a numbered menu and reads user selection from stdin.
func showMenu(cmd *cobra.Command) error {
	write.WhyBanner(cmd)
	cmd.Println()
	for i, doc := range DocOrder {
		write.WhyMenuItem(cmd, i+1, doc.Label)
	}
	write.WhyMenuPrompt(cmd)

	reader := bufio.NewReader(os.Stdin)
	input, readErr := reader.ReadString('\n')
	if readErr != nil {
		return fserr.ReadInput(readErr)
	}

	input = strings.TrimSpace(input)
	choice, parseErr := strconv.Atoi(input)
	if parseErr != nil || choice < 1 || choice > len(DocOrder) {
		return ctxerr.InvalidSelection(input, len(DocOrder))
	}

	cmd.Println()
	return ShowDoc(cmd, DocOrder[choice-1].Alias)
}

// ShowDoc loads an embedded document by alias, strips MkDocs syntax, and prints it.
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
		return ctxerr.UnknownDocument(alias)
	}

	content, loadErr := assets.WhyDoc(name)
	if loadErr != nil {
		return fserr.FileRead(name, loadErr)
	}

	cleaned := StripMkDocs(string(content))
	cmd.Print(cleaned)

	return nil
}
