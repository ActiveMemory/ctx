//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// resolveCmd returns the pad resolve subcommand.
//
// Returns:
//   - *cobra.Command: Configured resolve subcommand
func resolveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resolve",
		Short: "Show both sides of a merge conflict",
		Long: `Decrypt and display both sides of a merge conflict for the scratchpad.

Git stores conflict versions as .context/scratchpad.enc.ours and
.context/scratchpad.enc.theirs during a merge conflict. This command
decrypts both and displays them for manual resolution.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResolve(cmd)
		},
	}
}

// runResolve reads and prints both sides of a merge conflict.
func runResolve(cmd *cobra.Command) error {
	if !rc.ScratchpadEncrypt() {
		return errors.New("resolve is only needed for encrypted scratchpads")
	}

	key, err := crypto.LoadKey(keyPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New(errNoKey)
		}
		return fmt.Errorf("load key: %w", err)
	}

	dir := rc.ContextDir()
	oursPath := filepath.Join(dir, config.FileScratchpadEnc+".ours")
	theirsPath := filepath.Join(dir, config.FileScratchpadEnc+".theirs")

	ours, errOurs := decryptFile(key, oursPath)
	theirs, errTheirs := decryptFile(key, theirsPath)

	if errOurs != nil && errTheirs != nil {
		return fmt.Errorf("no conflict files found (%s.ours / %s.theirs)",
			config.FileScratchpadEnc, config.FileScratchpadEnc)
	}

	if errOurs == nil {
		cmd.Println("=== OURS ===")
		for i, entry := range ours {
			cmd.Printf("  %d. %s\n", i+1, displayEntry(entry))
		}
	}

	if errTheirs == nil {
		cmd.Println("=== THEIRS ===")
		for i, entry := range theirs {
			cmd.Printf("  %d. %s\n", i+1, displayEntry(entry))
		}
	}

	return nil
}

// decryptFile reads and decrypts a single file, returning its entries.
func decryptFile(key []byte, path string) ([]string, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is constructed from config constants
	if err != nil {
		return nil, err
	}

	plaintext, err := crypto.Decrypt(key, data)
	if err != nil {
		return nil, errors.New(errDecryptFail)
	}

	return parseEntries(plaintext), nil
}
