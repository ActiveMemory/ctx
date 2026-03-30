//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ceremony

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/config/ceremony"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/io"
)

// RecentJournalFiles returns the n most recent Markdown files in the given
// journal directory, sorted by filename descending (newest first).
//
// Parameters:
//   - dir: absolute path to the journal directory
//   - n: maximum number of files to return
//
// Returns:
//   - []string: absolute paths to the most recent journal files, or nil on
//     read error or empty directory
func RecentJournalFiles(dir string, n int) []string {
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return nil
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), file.ExtMarkdown) {
			continue
		}
		names = append(names, e.Name())
	}

	sort.Sort(sort.Reverse(sort.StringSlice(names)))

	if len(names) > n {
		names = names[:n]
	}

	paths := make([]string, len(names))
	for i, name := range names {
		paths[i] = filepath.Join(dir, name)
	}
	return paths
}

// ScanJournalsForCeremonies checks whether the given journal files contain
// references to /ctx-remember and /ctx-wrap-up ceremony commands.
//
// Parameters:
//   - files: absolute paths to journal files to scan
//
// Returns:
//   - remember: true if any file contains "ctx-remember"
//   - wrapUp: true if any file contains "ctx-wrap-up"
func ScanJournalsForCeremonies(files []string) (remember, wrapUp bool) {
	for _, path := range files {
		data, readErr := io.SafeReadUserFile(path)
		if readErr != nil {
			continue
		}
		content := string(data)
		if !remember && strings.Contains(content, ceremony.RememberCmd) {
			remember = true
		}
		if !wrapUp && strings.Contains(content, ceremony.WrapUpCmd) {
			wrapUp = true
		}
		if remember && wrapUp {
			return
		}
	}
	return
}

// Emit builds a ceremony nudge message box based on which
// ceremonies (remember, wrapUp) are missing from recent sessions.
//
// Parameters:
//   - remember: whether /ctx-remember was found in recent journals
//   - wrapUp: whether /ctx-wrap-up was found in recent journals
//
// Returns:
//   - msg: the formatted nudge message, or empty string if no content
//   - variant: the selected variant string for notifications
func Emit(remember, wrapUp bool) (msg, variant string) {
	var boxTitleKey, fallbackKey string

	switch {
	case !remember && !wrapUp:
		variant = hook.VariantBoth
		boxTitleKey = text.DescKeyCeremonyBoxBoth
		fallbackKey = text.DescKeyCeremonyFallbackBoth
	case !remember:
		variant = hook.VariantRemember
		boxTitleKey = text.DescKeyCeremonyBoxRemember
		fallbackKey = text.DescKeyCeremonyFallbackRemember
	case !wrapUp:
		variant = hook.VariantWrapup
		boxTitleKey = text.DescKeyCeremonyBoxWrapup
		fallbackKey = text.DescKeyCeremonyFallbackWrapup
	}

	boxTitle := desc.Text(boxTitleKey)
	fallback := desc.Text(fallbackKey)

	content := message.Load(hook.CheckCeremonies, variant, nil, fallback)
	if content == "" {
		return "", variant
	}

	relayPrefix := desc.Text(text.DescKeyCeremonyRelayPrefix)

	msg = message.NudgeBox(relayPrefix, boxTitle, content)
	return msg, variant
}
