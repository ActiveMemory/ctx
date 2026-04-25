//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package health

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/config/architecture"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/project"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	execGit "github.com/ActiveMemory/ctx/internal/exec/git"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// ReadMapTracking reads and parses the map-tracking.json file from the
// context directory.
//
// ctxDir is supplied by the caller (typically a FullPreamble-gated
// hook) so this function does not re-resolve it; a second resolution
// would be dead code today and would pair an ambiguous (nil, err)
// return with the genuine "no tracking yet" result.
//
// Returns (nil, nil) when the tracking file is simply absent: ordinary
// "nothing to track yet" state. A read or parse failure is propagated
// so the caller can distinguish "no tracking yet" from "tracking data
// is corrupt".
//
// Parameters:
//   - ctxDir: absolute path to the context directory
//
// Returns:
//   - *MapTrackingInfo: parsed tracking info, or nil if absent
//   - error: non-nil on I/O failure or JSON parse failure
func ReadMapTracking(ctxDir string) (*MapTrackingInfo, error) {
	data, readErr := io.SafeReadFile(ctxDir, architecture.MapTracking)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, nil
		}
		return nil, readErr
	}

	var info MapTrackingInfo
	if jsonErr := json.Unmarshal(data, &info); jsonErr != nil {
		return nil, jsonErr
	}

	return &info, nil
}

// CountModuleCommits counts git commits touching internal/
// since the given date.
//
// Parameters:
//   - since: date string in YYYY-MM-DD format
//
// Returns:
//   - int: number of commits, or 0 on error or if git is unavailable
func CountModuleCommits(since string) int {
	// Validate since as a date to prevent command injection.
	t, parseErr := time.Parse(cfgTime.DateFormat, since)
	if parseErr != nil {
		return 0
	}
	out, gitErr := execGit.LogSince(t,
		cfgGit.FlagOneline,
		cfgGit.FlagPathSep, project.DirInternalSlash,
	)
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
//   - string: formatted nudge box, or empty string if silenced.
//   - error: propagated from [nudge.EmitAndRelay] so callers can
//     honor the log-first principle: if the relay audit entry or
//     webhook fails, the nudge box should not be printed.
func EmitMapStalenessWarning(
	sessionID, dateStr string, moduleCommits int,
) (string, error) {
	fallback := fmt.Sprintf(
		desc.Text(text.DescKeyCheckMapStalenessFallback),
		dateStr, moduleCommits,
	)
	content := message.Load(hook.CheckMapStaleness, hook.VariantStale,
		map[string]any{
			architecture.VarLastRefreshDate: dateStr,
			architecture.VarModuleCount:     moduleCommits,
		}, fallback)
	if content == "" {
		return "", nil
	}

	box := message.NudgeBox(
		desc.Text(text.DescKeyCheckMapStalenessRelayPrefix),
		desc.Text(text.DescKeyCheckMapStalenessBoxTitle),
		content)

	ref := notify.NewTemplateRef(hook.CheckMapStaleness, hook.VariantStale,
		map[string]any{
			architecture.VarLastRefreshDate: dateStr,
			architecture.VarModuleCount:     moduleCommits,
		},
	)
	notifyMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckMapStaleness,
		desc.Text(text.DescKeyCheckMapStalenessRelayMessage),
	)
	if err := nudge.EmitAndRelay(notifyMsg, sessionID, ref); err != nil {
		return "", err
	}
	return box, nil
}
