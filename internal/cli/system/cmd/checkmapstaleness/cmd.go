//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkmapstaleness

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

const mapStaleDays = 30

// Cmd returns the "ctx system check-map-staleness" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-map-staleness subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-map-staleness",
		Short: "Architecture map staleness nudge",
		Long: `Checks whether map-tracking.json is stale (>30 days) and there are
commits touching internal/ since the last map refresh. Outputs a VERBATIM
relay nudge suggesting /ctx-map when both conditions are met.

Hook event: UserPromptSubmit
Output: VERBATIM relay (when stale and modules changed), silent otherwise
Silent when: map-tracking.json missing or fresh, opted out, no module
commits, already nudged today, or uninitialized`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckMapStaleness(cmd, os.Stdin)
		},
	}
}

// mapTrackingInfo holds the minimal fields needed from map-tracking.json.
type mapTrackingInfo struct {
	OptedOut bool   `json:"opted_out"`
	LastRun  string `json:"last_run"`
}

func runCheckMapStaleness(cmd *cobra.Command, stdin *os.File) error {
	if !core.IsInitialized() {
		return nil
	}

	input := core.ReadInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}
	if core.Paused(sessionID) > 0 {
		return nil
	}
	markerPath := filepath.Join(core.StateDir(), "check-map-staleness")
	if core.IsDailyThrottled(markerPath) {
		return nil
	}

	contextDir := rc.ContextDir()
	trackingPath := filepath.Join(contextDir, config.FileMapTracking)

	data, readErr := os.ReadFile(trackingPath) //nolint:gosec // project-local path
	if readErr != nil {
		return nil // no tracking file — nothing to nudge about
	}

	var info mapTrackingInfo
	if jsonErr := json.Unmarshal(data, &info); jsonErr != nil {
		return nil
	}

	if info.OptedOut {
		return nil
	}

	lastRun, parseErr := time.Parse("2006-01-02", info.LastRun)
	if parseErr != nil {
		return nil
	}

	if time.Since(lastRun) < time.Duration(mapStaleDays)*24*time.Hour {
		return nil
	}

	// Count commits touching internal/ since last run
	moduleCommits := countModuleCommits(info.LastRun)
	if moduleCommits == 0 {
		return nil
	}

	// Emit VERBATIM nudge
	dateStr := lastRun.Format("2006-01-02")
	fallback := fmt.Sprintf("ARCHITECTURE.md hasn't been refreshed since %s\n", dateStr) +
		fmt.Sprintf("and there are commits touching %d modules.\n", moduleCommits) +
		"/ctx-map keeps architecture docs drift-free.\n" +
		config.NewlineLF +
		"Want me to run /ctx-map to refresh?"
	content := core.LoadMessage("check-map-staleness", "stale",
		map[string]any{
			"LastRefreshDate": dateStr,
			"ModuleCount":     moduleCommits,
		}, fallback)
	if content == "" {
		return nil
	}

	msg := "IMPORTANT: Relay this architecture map notice to the user VERBATIM before answering their question.\n\n" +
		"\u250c\u2500 Architecture Map Stale \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\n"
	msg += core.BoxLines(content)
	if line := core.ContextDirLine(); line != "" {
		msg += "\u2502 " + line + config.NewlineLF
	}
	msg += "\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500"
	cmd.Println(msg)

	ref := notify.NewTemplateRef("check-map-staleness", "stale",
		map[string]any{"LastRefreshDate": dateStr, "ModuleCount": moduleCommits})
	_ = notify.Send("nudge", "check-map-staleness: Architecture map stale", input.SessionID, ref)
	_ = notify.Send("relay", "check-map-staleness: Architecture map stale", input.SessionID, ref)
	eventlog.Append("relay", "check-map-staleness: Architecture map stale", input.SessionID, ref)

	core.TouchFile(markerPath)

	return nil
}

// countModuleCommits counts git commits touching internal/ since the given date.
func countModuleCommits(since string) int {
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
	return len(strings.Split(lines, config.NewlineLF))
}
