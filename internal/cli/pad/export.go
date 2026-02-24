//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// exportCmd returns the pad export subcommand.
func exportCmd() *cobra.Command {
	var force, dryRun bool

	cmd := &cobra.Command{
		Use:   "export [DIR]",
		Short: "Export blob entries to a directory as files",
		Long: `Export all blob entries from the scratchpad to a directory as files.
Each blob's label becomes the filename. Non-blob entries are skipped.

When a file already exists, a unix timestamp is prepended to avoid
collisions. Use --force to overwrite instead.

Examples:
  ctx pad export
  ctx pad export ./ideas
  ctx pad export --dry-run
  ctx pad export --force ./backup`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			return runExport(cmd, dir, force, dryRun)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing files instead of timestamping")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print what would be exported without writing")

	return cmd
}

// runExport exports blob entries to the given directory.
func runExport(cmd *cobra.Command, dir string, force, dryRun bool) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	if !dryRun {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}

	var count int
	for _, entry := range entries {
		label, data, ok := splitBlob(entry)
		if !ok {
			continue
		}

		outPath := filepath.Join(dir, label)

		if !force {
			if _, err := os.Stat(outPath); err == nil {
				ts := fmt.Sprintf("%d", time.Now().Unix())
				newName := ts + "-" + label
				if dryRun {
					cmd.Printf("  %s → %s (exists)\n", label, filepath.Join(dir, newName))
					count++
					continue
				}
				outPath = filepath.Join(dir, newName)
				cmd.Printf("  ! %s exists, writing as %s\n", label, newName)
			}
		}

		if dryRun {
			cmd.Printf("  %s → %s\n", label, outPath)
			count++
			continue
		}

		if err := os.WriteFile(outPath, data, 0o600); err != nil {
			cmd.PrintErrf("  ! failed to write %s: %v\n", label, err)
			continue
		}

		cmd.Printf("  + %s\n", label)
		count++
	}

	if count == 0 {
		cmd.Println("No blob entries to export.")
		return nil
	}

	verb := "Exported"
	if dryRun {
		verb = "Would export"
	}
	cmd.Printf("%s %d blobs.\n", verb, count)
	return nil
}
