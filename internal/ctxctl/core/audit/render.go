//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"fmt"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/token"
	auditStore "github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/core/store"
	cfgAudit "github.com/ActiveMemory/ctx/internal/ctxctl/config/audit"
)

// RenderReports formats a slice of due reports for the
// audit-relay hook's relay box. Each report body is
// prefixed with its id and commit-range header; reports
// older than [cfgAudit.StalenessAge] get a STALE prefix
// line. Multiple reports are joined with the supplied
// separator.
//
// Parameters:
//   - reports: due reports (post-filter)
//   - now: wall-clock anchor for staleness comparison
//   - separator: rule drawn between multiple reports
//   - stalePrefixFmt: STALE prefix format string
//     (verbs: commit-range, age)
//
// Returns:
//   - string: rendered body suitable for direct injection
//     into the verbatim-relay envelope
func RenderReports(
	reports []auditStore.Report, now time.Time,
	separator, stalePrefixFmt string,
) string {
	parts := make([]string, 0, len(reports))
	for _, r := range reports {
		var b strings.Builder
		header := fmt.Sprintf(
			cfgAudit.FmtReportHeader,
			r.ID, r.CommitRange, token.NewlineLF,
		)
		b.WriteString(header)
		age := now.Sub(r.GeneratedAt)
		if cfgAudit.StalenessAge > 0 && age > cfgAudit.StalenessAge {
			staleBody := fmt.Sprintf(stalePrefixFmt,
				r.CommitRange,
				AgeString(age),
			)
			staleLine := fmt.Sprintf(
				cfgAudit.FmtStaleLine,
				staleBody, token.NewlineLF,
			)
			b.WriteString(staleLine)
		}
		b.WriteString(r.Body)
		parts = append(parts, b.String())
	}
	return strings.Join(parts,
		token.NewlineLF+separator+token.NewlineLF)
}

// AgeString rounds a duration to whole days when the
// duration spans at least one day, otherwise hours.
// Keeps the STALE prefix concise.
//
// Parameters:
//   - d: positive duration since the report was generated
//
// Returns:
//   - string: humanized age suffix (e.g. "9d", "3h")
func AgeString(d time.Duration) string {
	dayLen := time.Duration(cfgAudit.HoursPerDay) * time.Hour
	if d >= dayLen {
		return fmt.Sprintf(
			cfgAudit.FmtAgeDays,
			int(d.Hours())/cfgAudit.HoursPerDay,
		)
	}
	return fmt.Sprintf(cfgAudit.FmtAgeHours, int(d.Hours()))
}
