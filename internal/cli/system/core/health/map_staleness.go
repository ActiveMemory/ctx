//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package health

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/config/architecture"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"

	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// MapTrackingInfo holds the minimal fields needed from map-tracking.json.
type MapTrackingInfo struct {
	OptedOut bool   `json:"opted_out"`
	LastRun  string `json:"last_run"`
}

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

// EmitMapStalenessWarning builds the architecture map staleness warning box.
//
// Parameters:
//   - sessionID: session identifier for notifications
//   - dateStr: last refresh date (YYYY-MM-DD)
//   - moduleCommits: number of commits touching modules since last refresh
//
// Returns:
//   - string: formatted nudge box, or empty string if silenced
func EmitMapStalenessWarning(sessionID, dateStr string, moduleCommits int) string {
	fallback := fmt.Sprintf(desc.Text(text.DescKeyCheckMapStalenessFallback), dateStr, moduleCommits)
	content := message.LoadMessage(hook.CheckMapStaleness, hook.VariantStale,
		map[string]any{
			architecture.VarLastRefreshDate: dateStr,
			architecture.VarModuleCount:     moduleCommits,
		}, fallback)
	if content == "" {
		return ""
	}

	box := message.NudgeBox(
		desc.Text(text.DescKeyCheckMapStalenessRelayPrefix),
		desc.Text(text.DescKeyCheckMapStalenessBoxTitle),
		content)

	ref := notify.NewTemplateRef(hook.CheckMapStaleness, hook.VariantStale,
		map[string]any{architecture.VarLastRefreshDate: dateStr, architecture.VarModuleCount: moduleCommits})
	notifyMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckMapStaleness, desc.Text(text.DescKeyCheckMapStalenessRelayMessage))
	nudge.NudgeAndRelay(notifyMsg, sessionID, ref)
	return box
}
