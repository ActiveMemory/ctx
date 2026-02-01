//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package serve

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the serve command.
//
// Serves a static site by invoking zensical serve on the specified directory.
//
// Returns:
//   - *cobra.Command: The serve command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve [directory]",
		Short: "Serve a static site locally via zensical",
		Long: `Serve a static site using zensical.

If no directory is specified, serves the journal site (.context/journal-site).

Requires zensical to be installed:
  pip install zensical

Examples:
  ctx serve                           # Serve journal site
  ctx serve .context/journal-site     # Serve specific directory
  ctx serve ./docs                    # Serve docs folder`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(cmd, args)
		},
	}

	return cmd
}

// runServe handles the serve command.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: Optional directory to serve
//
// Returns:
//   - error: Non-nil if zensical is not found or fails
func runServe(cmd *cobra.Command, args []string) error {
	var dir string

	if len(args) > 0 {
		dir = args[0]
	} else {
		// Default: journal site
		dir = filepath.Join(rc.ContextDir(), "journal-site")
	}

	// Verify directory exists
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("directory not found: %s", dir)
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	// Check zensical.toml exists
	tomlPath := filepath.Join(dir, "zensical.toml")
	if _, err := os.Stat(tomlPath); os.IsNotExist(err) {
		return fmt.Errorf("no zensical.toml found in %s", dir)
	}

	// Check if zensical is available
	_, err = exec.LookPath("zensical")
	if err != nil {
		return fmt.Errorf("zensical not found. Install with: pip install zensical")
	}

	// Run zensical serve
	zensical := exec.Command("zensical", "serve")
	zensical.Dir = dir
	zensical.Stdout = os.Stdout
	zensical.Stderr = os.Stderr
	zensical.Stdin = os.Stdin

	return zensical.Run()
}
