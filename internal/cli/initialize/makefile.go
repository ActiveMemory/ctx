//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// includeDirective is the line appended to the user's Makefile to pull
// in ctx targets. The leading dash suppresses errors when the file is
// absent.
const includeDirective = "-include Makefile.ctx"

// handleMakefileCtx deploys Makefile.ctx and ensures the user's
// Makefile includes it.
//
// Makefile.ctx is fully owned by ctx and always overwritten.
// The user's Makefile is only amended (never overwritten): if it
// exists, the include directive is appended when missing; if it does
// not exist, a minimal Makefile is created.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if template read or file operations fail
func handleMakefileCtx(cmd *cobra.Command) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Deploy Makefile.ctx (always overwrite — ctx-owned)
	content, err := assets.MakefileCtx()
	if err != nil {
		return fmt.Errorf("failed to read Makefile.ctx template: %w", err)
	}

	if err = os.WriteFile(
		config.FileMakefileCtx, content, config.PermFile,
	); err != nil {
		return fmt.Errorf(
			"failed to write %s: %w", config.FileMakefileCtx, err,
		)
	}
	cmd.Printf("  %s %s\n", green("✓"), config.FileMakefileCtx)

	// Ensure the user's Makefile includes Makefile.ctx
	existing, err := os.ReadFile("Makefile")
	if err != nil {
		// No Makefile — create a minimal one
		minimal := includeDirective + "\n"
		if err := os.WriteFile(
			"Makefile", []byte(minimal), config.PermFile,
		); err != nil {
			return fmt.Errorf("failed to create Makefile: %w", err)
		}
		cmd.Printf("  %s Makefile (created with ctx include)\n", green("✓"))
		return nil
	}

	// Makefile exists — check if it already includes Makefile.ctx
	if strings.Contains(string(existing), includeDirective) {
		cmd.Printf(
			"  %s Makefile (already includes %s)\n",
			yellow("○"), config.FileMakefileCtx,
		)
		return nil
	}

	// Append the include directive
	amended := string(existing)
	if !strings.HasSuffix(amended, "\n") {
		amended += "\n"
	}
	amended += "\n" + includeDirective + "\n"

	if err := os.WriteFile(
		"Makefile", []byte(amended), config.PermFile,
	); err != nil {
		return fmt.Errorf("failed to amend Makefile: %w", err)
	}
	cmd.Printf(
		"  %s Makefile (appended %s include)\n",
		green("✓"), config.FileMakefileCtx,
	)

	return nil
}
