//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// SystemContextUpdate matches context-update XML tags.
//
// Groups:
//   - 1: opening tag attributes (e.g., ` type="task" context="..."`)
//   - 2: content between tags
var SystemContextUpdate = regexp.MustCompile(`<context-update(\s+[^>]+)>([^<]+)</context-update>`)

// SystemClaudeTag matches Claude Code internal markup tags that leak into
// session titles via the first user message. This MUST remain an allowlist
// of known Claude Code tags — do NOT replace with a blanket regex.
var SystemClaudeTag = regexp.MustCompile(`</?(?:command-message|command-name|local-command-caveat)>`)

// SystemReminder matches <system-reminder>...</system-reminder> blocks.
// These are injected by Claude Code into tool results.
// Groups:
//   - 1: content between tags
var SystemReminder = regexp.MustCompile(`(?s)<system-reminder>\s*(.*?)\s*</system-reminder>`)
