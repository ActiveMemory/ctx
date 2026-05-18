//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout

import (
	"strings"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// buildFilename derives a closeout's on-disk name from its
// frontmatter.
//
// Shape: `<TS>-<mode>-closeout.md` where `<TS>` is the UTC
// compact RFC-3339 form (`20260517T021837Z`; colons stripped
// for filesystem safety) and `<mode>` is the frontmatter's
// Mode field.
//
// Parameters:
//   - fm: frontmatter (uses Mode + GeneratedAt).
//
// Returns:
//   - string: filename portion only (no directory).
func buildFilename(fm Frontmatter) string {
	var sb strings.Builder
	sb.WriteString(fm.GeneratedAt.UTC().Format(cfgTime.RFC3339Compact))
	sb.WriteString(token.Dash)
	sb.WriteString(fm.Mode)
	sb.WriteString(cfgKB.CloseoutSuffix)
	return sb.String()
}
