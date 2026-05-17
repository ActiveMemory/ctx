//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

import (
	"strings"

	cfgKbCli "github.com/ActiveMemory/ctx/internal/config/kb/cli"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// RenderBlock builds the managed block content that the
// reindex command splices between the CTX:KB:TOPICS START and
// END markers in `.context/kb/index.md`.
//
// Parameters:
//   - slugs: sorted topic slugs.
//
// Returns:
//   - string: managed block with delimiters.
func RenderBlock(slugs []string) string {
	var sb strings.Builder
	sb.WriteString(cfgKbCli.ManagedKBTopicsStart)
	sb.WriteString(token.NewlineLF)
	if len(slugs) == 0 {
		sb.WriteString(cfgKbCli.ManagedKBTopicsEmpty)
	}
	for _, s := range slugs {
		sb.WriteString(cfgKbCli.TopicEntryPrefix)
		sb.WriteString(s)
		sb.WriteString(cfgKbCli.TopicEntryMiddle)
		sb.WriteString(s)
		sb.WriteString(cfgKbCli.TopicEntrySuffix)
	}
	sb.WriteString(cfgKbCli.ManagedKBTopicsEnd)
	return sb.String()
}
