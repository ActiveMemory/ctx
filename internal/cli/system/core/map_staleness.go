//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/architecture"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// ReadMapTracking reads and parses the map-tracking.json file from the
// context directory.
//
// Returns:
//   - *MapTrackingInfo: parsed tracking info, or nil if not found or invalid
func ReadMapTracking() *MapTrackingInfo {
	data, readErr := io.SafeReadFile(rc.ContextDir(), architecture.MapTracking)
	if readErr != nil {
		return nil
	}

	var info MapTrackingInfo
	if jsonErr := json.Unmarshal(data, &info); jsonErr != nil {
		return nil
	}

	return &info
}

// CountModuleCommits counts git commits touching internal/ since the given date.
//
// Parameters:
//   - since: date string in YYYY-MM-DD format
//
// Returns:
//   - int: number of commits, or 0 on error or if git is unavailable
func CountModuleCommits(since string) int {
	if _, lookErr := exec.LookPath("git"); lookErr != nil {
		return 0
	}
	out, gitErr := exec.Command("git", "log", "--oneline", "--since="+since, "--", "internal/").Output() //nolint:gosec // date string from JSON
	if gitErr != nil {
		return 0
	}
	lines := strings.TrimSpace(string(out))
	if lines == "" {
		return 0
	}
	return len(strings.Split(lines, token.NewlineLF))
}

// EmitMapStalenessWarning builds and prints the architecture map staleness
// warning box.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionID: session identifier for notifications
//   - dateStr: last refresh date (YYYY-MM-DD)
//   - moduleCommits: number of commits touching modules since last refresh
func EmitMapStalenessWarning(cmd *cobra.Command, sessionID, dateStr string, moduleCommits int) {
	fallback := fmt.Sprintf(desc.TextDesc(text.DescKeyCheckMapStalenessFallback), dateStr, moduleCommits)
	content := LoadMessage(hook.CheckMapStaleness, hook.VariantStale,
		map[string]any{
			tpl.VarLastRefreshDate: dateStr,
			tpl.VarModuleCount:     moduleCommits,
		}, fallback)
	if content == "" {
		return
	}

	cmd.Println(NudgeBox(
		desc.TextDesc(text.DescKeyCheckMapStalenessRelayPrefix),
		desc.TextDesc(text.DescKeyCheckMapStalenessBoxTitle),
		content))

	ref := notify.NewTemplateRef(hook.CheckMapStaleness, hook.VariantStale,
		map[string]any{tpl.VarLastRefreshDate: dateStr, tpl.VarModuleCount: moduleCommits})
	notifyMsg := fmt.Sprintf(desc.TextDesc(text.DescKeyRelayPrefixFormat),
		hook.CheckMapStaleness, desc.TextDesc(text.DescKeyCheckMapStalenessRelayMessage))
	NudgeAndRelay(notifyMsg, sessionID, ref)
}
