//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/rc"
)

// Setup configures the scratchpad key or plaintext file.
//
// When encryption is enabled (default):
//   - Generates a 256-bit key at ~/.ctx/ if not present
//   - Warns if .enc exists but no key
//
// When encryption is disabled:
//   - Creates empty .context/scratchpad.md if not present
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: The .context/ directory path
//
// Returns:
//   - error: Non-nil if key generation or file operations fail
func Setup(cmd *cobra.Command, contextDir string) error {
	if !rc.ScratchpadEncrypt() {
		return setupPlaintext(cmd, contextDir)
	}
	return setupEncrypted(cmd, contextDir)
}
