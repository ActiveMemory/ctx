//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package why

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// runWhy dispatches to the interactive menu or direct document display.
func runWhy(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return showDoc(cmd, args[0])
	}
	return showMenu(cmd)
}

// showMenu presents a numbered menu and reads user selection from stdin.
func showMenu(cmd *cobra.Command) error {
	bt := "`"
	cmd.Println(`
   /    ctx:                         https://ctx.ist
 ,'` + bt + `./    do you remember?
 ` + bt + `.,'\
   \
      {}  -> what
      ctx -> why`)
	cmd.Println()
	for i, doc := range docOrder {
		cmd.Printf("  [%d] %s\n", i+1, doc.label)
	}
	cmd.Print("\nSelect a document (1-3): ")

	reader := bufio.NewReader(os.Stdin)
	input, readErr := reader.ReadString('\n')
	if readErr != nil {
		return fmt.Errorf("reading selection: %w", readErr)
	}

	input = strings.TrimSpace(input)
	choice, parseErr := strconv.Atoi(input)
	if parseErr != nil || choice < 1 || choice > len(docOrder) {
		return fmt.Errorf("invalid selection: %q (expected 1-%d)", input, len(docOrder))
	}

	cmd.Println()
	return showDoc(cmd, docOrder[choice-1].alias)
}

// showDoc loads an embedded document by alias, strips MkDocs syntax, and prints it.
func showDoc(cmd *cobra.Command, alias string) error {
	name, ok := docAliases[alias]
	if !ok {
		return fmt.Errorf("unknown document %q (available: manifesto, about, invariants)", alias)
	}

	content, loadErr := assets.WhyDoc(name)
	if loadErr != nil {
		return fmt.Errorf("loading document %q: %w", name, loadErr)
	}

	cleaned := StripMkDocs(string(content))
	cmd.Print(cleaned)

	return nil
}
