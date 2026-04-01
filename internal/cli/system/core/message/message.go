//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/hook"
	"github.com/ActiveMemory/ctx/internal/config/box"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/resolve"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Load loads a hook message template by hook name and variant.
//
// Priority:
//  1. .context/hooks/messages/{hook}/{variant}.txt (user override)
//  2. internal/assets/hooks/messages/{hook}/{variant}.txt (embedded default)
//  3. fallback string (hardcoded, belt and suspenders)
//
// Returns empty string if the resolved template is empty or whitespace-only
// (intentional silence). The vars map provides template variables;
// nil is valid when no dynamic content is needed.
//
// Parameters:
//   - hook: Hook name
//   - variant: Template variant name
//   - vars: Template variables (nil for static messages)
//   - fallback: Hardcoded fallback string
//
// Returns:
//   - string: Rendered message or empty string for intentional silence
func Load(hk, variant string, vars map[string]any, fallback string) string {
	filename := variant + file.ExtTxt

	// 1. User override in .context/
	overrideDir := filepath.Join(rc.ContextDir(), dir.HooksMessages, hk)
	if data, readErr := io.SafeReadFile(overrideDir, filename); readErr == nil {
		return renderTemplate(string(data), vars, fallback)
	}

	// 2. Embedded default
	if data, readErr := hook.Message(hk, filename); readErr == nil {
		return renderTemplate(string(data), vars, fallback)
	}

	// 3. Hardcoded fallback
	return renderTemplate(fallback, vars, fallback)
}

// BoxLines wraps each line of content with the │ box-drawing prefix.
// Trailing newlines on content are trimmed before splitting to avoid
// an empty trailing box line.
//
// Parameters:
//   - content: Multi-line string to wrap
//
// Returns:
//   - string: Box-wrapped content
func BoxLines(content string) string {
	var b strings.Builder
	trimmed := strings.TrimRight(content, token.NewlineLF)
	for _, line := range strings.Split(trimmed, token.NewlineLF) {
		b.WriteString(box.LinePrefix)
		b.WriteString(line)
		b.WriteString(token.NewlineLF)
	}
	return b.String()
}

// NudgeBox builds a complete nudge box with relay prefix, titled top
// border, box-wrapped content, optional context directory footer, and
// bottom border.
//
// Parameters:
//   - relayPrefix: VERBATIM relay instruction line
//   - title: box title (e.g., "Backup Warning")
//   - content: multi-line body text
//
// Returns:
//   - string: fully formatted nudge box
func NudgeBox(relayPrefix, title, content string) string {
	pad := box.NudgeBoxWidth - len(title)
	if pad < 0 {
		pad = 0
	}
	msg := relayPrefix + token.NewlineLF + token.NewlineLF +
		box.Top + title + " " + strings.Repeat(box.BorderFill, pad) + token.NewlineLF
	msg += BoxLines(content)
	if line := ctxContext.DirLine(); line != "" {
		msg += box.LinePrefix + line + token.NewlineLF
	}
	msg += box.Bottom
	return msg
}
