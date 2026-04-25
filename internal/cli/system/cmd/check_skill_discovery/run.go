//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_skill_discovery

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/counter"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Run executes the skill discovery nudge hook logic.
//
// Reads hook input from stdin, checks the session prompt count,
// and fires a one-shot nudge at the configured threshold. The
// nudge surfaces mid-session skills that are easy to forget.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	_, sessionID, _, tmpDir, ok := coreCheck.FullPreamble(stdin)
	bailSilently := !ok
	if bailSilently {
		return nil
	}

	// One-shot guard: skip if already fired this session.
	guardFile := filepath.Join(
		tmpDir,
		cfgHook.PrefixSkillDiscoveryGuard+sessionID,
	)
	if _, statErr := os.Stat(guardFile); statErr == nil {
		return nil
	}

	// Read the prompt counter shared with check-context-size.
	counterFile := filepath.Join(
		tmpDir,
		stats.ContextSizeCounterPrefix+sessionID,
	)
	count := counter.Read(counterFile)
	if count < cfgHook.SkillDiscoveryThreshold {
		return nil
	}

	// Fire the nudge and mark as done.
	content := message.Load(
		cfgHook.CheckSkillDiscovery, "",
		nil,
		desc.Text(text.DescKeySkillDiscoveryContent),
	)
	if content == "" {
		internalIo.TouchFile(guardFile)
		return nil
	}

	box := message.NudgeBox(
		desc.Text(text.DescKeySkillDiscoveryPrefix),
		desc.Text(text.DescKeySkillDiscoveryBoxTitle),
		content,
	)
	writeSetup.NudgeBlock(cmd, box)
	internalIo.TouchFile(guardFile)

	return nil
}
