//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_freshness

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/freshness"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the check-freshness hook logic.
//
// Reads tracked files from .ctxrc freshness_files config. For each
// file, stats it and warns if it has not been modified within the
// freshness window (~6 months). Files that do not exist are silently
// skipped. The hook is a no-op when no files are configured.
// Throttled to once per day.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	input, _, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	files := rc.FreshnessFiles()
	if len(files) == 0 {
		return nil
	}

	tmpDir := core.StateDir()
	throttleFile := filepath.Join(tmpDir, freshness.ThrottleID)

	if core.IsDailyThrottled(throttleFile) {
		return nil
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return nil
	}

	now := time.Now()
	var staleEntries []staleEntry

	for _, tf := range files {
		absPath := filepath.Join(cwd, tf.Path)

		info, statErr := os.Stat(absPath)
		if statErr != nil {
			continue
		}

		age := now.Sub(info.ModTime())
		if age <= freshness.StaleThreshold {
			continue
		}

		staleEntries = append(staleEntries, staleEntry{
			path:      tf.Path,
			desc:      tf.Desc,
			reviewURL: tf.ReviewURL,
			days:      int(age.Hours() / 24),
		})
	}

	if len(staleEntries) == 0 {
		return nil
	}

	staleText := formatStaleEntries(staleEntries)

	vars := map[string]any{tpl.VarStaleFiles: staleText}
	content := core.LoadMessage(hook.CheckFreshness, hook.VariantStale, vars, staleText)
	if content == "" {
		return nil
	}

	cmd.Println(core.NudgeBox(
		desc.Text(text.DescKeyFreshnessRelayPrefix),
		desc.Text(text.DescKeyFreshnessBoxTitle),
		content))

	ref := notify.NewTemplateRef(hook.CheckFreshness, hook.VariantStale, vars)
	core.NudgeAndRelay(hook.CheckFreshness+": "+
		desc.Text(text.DescKeyFreshnessRelayMessage),
		input.SessionID, ref,
	)

	core.TouchFile(throttleFile)

	return nil
}

type staleEntry struct {
	path      string
	desc      string
	reviewURL string
	days      int
}

// formatStaleEntries builds the display text for stale files.
//
// Groups entries by review URL. Entries without a URL get the generic
// "touch to mark as reviewed" footer. Entries with URLs get a
// "Review against: <url>" line after their group.
func formatStaleEntries(entries []staleEntry) string {
	// Partition: with URL (grouped) vs without URL
	byURL := make(map[string][]staleEntry)
	var noURL []staleEntry
	var urlOrder []string

	for _, e := range entries {
		if e.reviewURL == "" {
			noURL = append(noURL, e)
			continue
		}
		if _, seen := byURL[e.reviewURL]; !seen {
			urlOrder = append(urlOrder, e.reviewURL)
		}
		byURL[e.reviewURL] = append(byURL[e.reviewURL], e)
	}

	var b strings.Builder

	// Entries with review URLs, grouped
	for _, url := range urlOrder {
		group := byURL[url]
		for _, e := range group {
			_, err := fmt.Fprintf(&b, desc.Text(text.DescKeyFreshnessFileEntry),
				e.path, e.days, e.desc)
			if err != nil {
				return ""
			}
			b.WriteString(token.NewlineLF)
		}
		_, err := fmt.Fprintf(&b, desc.Text(text.DescKeyFreshnessReviewURL), url)
		if err != nil {
			return ""
		}
		b.WriteString(token.NewlineLF)
	}

	// Entries without review URLs
	for _, e := range noURL {
		_, err := fmt.Fprintf(&b, desc.Text(text.DescKeyFreshnessFileEntry),
			e.path, e.days, e.desc)
		if err != nil {
			return ""
		}
		b.WriteString(token.NewlineLF)
	}

	b.WriteString(desc.Text(text.DescKeyFreshnessTouchHint))

	return b.String()
}
