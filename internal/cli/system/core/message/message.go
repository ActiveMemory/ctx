//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/hook"
	"github.com/ActiveMemory/ctx/internal/config/box"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/resolve"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	"github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
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
	//
	// When no context dir is declared there is no place to look
	// for a user override, so we fall through to the embedded
	// default. Reading the template itself is best-effort: any
	// read error just means "no override", same outcome. A non-
	// declaration resolver failure gets logged because it points
	// at a real regression (bad .ctxrc, permission issue) that
	// would otherwise silently disable overrides.
	ctxDir, ctxErr := rc.ContextDir()
	switch {
	case ctxErr == nil:
		overrideDir := filepath.Join(ctxDir, dir.HooksMessages, hk)
		if data, readErr := io.SafeReadFile(overrideDir, filename); readErr == nil {
			return renderTemplate(string(data), vars, fallback)
		}
	case !errors.Is(ctxErr, errCtx.ErrDirNotDeclared):
		logWarn.Warn(warn.ContextDirResolve, ctxErr)
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
		box.Top + title + token.Space +
		strings.Repeat(box.BorderFill, pad) + token.NewlineLF
	msg += BoxLines(content)
	// Rendering-layer swallow: [NudgeBox] returns only a string, so
	// there is no channel to propagate a DirLine resolver error. The
	// noisy-TUI log inside DirLine already surfaces any unexpected
	// failure, and in hook contexts (the usual caller chain)
	// FullPreamble already gated ContextDir, so the error path is
	// unreachable in practice. If the line is empty or the resolver
	// errored, skip the footer.
	line, _ := ctxContext.DirLine()
	if line != "" {
		msg += box.LinePrefix + line + token.NewlineLF
	}
	msg += box.Bottom
	return msg
}
