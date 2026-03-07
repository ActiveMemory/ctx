//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// IncludeDirective is the line appended to the user's Makefile to pull
// in ctx targets. The leading dash suppresses errors when the file is absent.
const IncludeDirective = "-include Makefile.ctx"

// HandleMakefileCtx deploys Makefile.ctx and amends the user Makefile.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if file operations fail
func HandleMakefileCtx(cmd *cobra.Command) error {
	content, err := assets.MakefileCtx()
	if err != nil {
		return fmt.Errorf("failed to read Makefile.ctx template: %w", err)
	}
	if err = os.WriteFile(config.FileMakefileCtx, content, config.PermFile); err != nil {
		return fmt.Errorf("failed to write %s: %w", config.FileMakefileCtx, err)
	}
	cmd.Println(fmt.Sprintf("  ✓ %s", config.FileMakefileCtx))
	existing, err := os.ReadFile("Makefile")
	if err != nil {
		minimal := IncludeDirective + config.NewlineLF
		if err := os.WriteFile("Makefile", []byte(minimal), config.PermFile); err != nil {
			return fmt.Errorf("failed to create Makefile: %w", err)
		}
		cmd.Println("  ✓ Makefile (created with ctx include)")
		return nil
	}
	if strings.Contains(string(existing), IncludeDirective) {
		cmd.Println(fmt.Sprintf("  ○ Makefile (already includes %s)\n", config.FileMakefileCtx))
		return nil
	}
	amended := string(existing)
	if !strings.HasSuffix(amended, config.NewlineLF) {
		amended += config.NewlineLF
	}
	amended += config.NewlineLF + IncludeDirective + config.NewlineLF
	if err := os.WriteFile("Makefile", []byte(amended), config.PermFile); err != nil {
		return fmt.Errorf("failed to amend Makefile: %w", err)
	}
	cmd.Println(fmt.Sprintf("  ✓ Makefile (appended %s include)\n", config.FileMakefileCtx))
	return nil
}
