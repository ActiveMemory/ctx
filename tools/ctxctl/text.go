//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/cmd/dismiss"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/cmd/list"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/cmd/show"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/checkaudit"
)

// ctxctl owns its user-facing text as plain English Go
// constants, outside ctx's YAML localization and the
// desc/i18n engine: there is no French ctxctl. The values
// below were lifted verbatim from ctx's former audit
// descriptors so the migration preserves output exactly.
// See specs/ctxctl-bootstrap.md and DECISIONS.md (2026-05-27).

// Root command text.
const (
	rootUse   = "ctxctl"
	rootShort = "ctx maintainer and contributor tooling"
	rootLong  = "ctxctl houses maintainer-only tooling kept out of the " +
		"shipped ctx binary. First inhabitant: the out-of-band " +
		"audit channel."
)

// audit command + subcommand cobra Use strings.
const (
	useAudit        = "audit"
	useAuditList    = "list"
	useAuditShow    = "show ID"
	useAuditDismiss = "dismiss [ID...]"
	useAuditRelay   = "audit-relay"
)

// audit command + subcommand descriptions.
const (
	shortAudit = "Show and manage out-of-band audit reports"
	longAudit  = `Surface out-of-band audit reports dropped under .context/audit/.

The audit channel decouples discipline enforcement from the in-band
commit cadence: a separate Claude Code session runs an audit skill
(e.g. /ctx-surface-audit) against the current branch, drops a
structured report into .context/audit/<kind>.md, and the
UserPromptSubmit hook ` + "`ctxctl audit-relay`" + ` verbatim-relays it
on the next interactive turn.

` + "`ctxctl audit`" + ` (no subcommand) lists every report with status, age,
and dismissed-state. Subcommands manage the lifecycle:

  list     Show all reports (default action)
  show     Print one report's body
  dismiss  Mark one or more reports dismissed`

	shortAuditList    = "List all audit reports with status and age"
	shortAuditShow    = "Print the body of an audit report by id"
	shortAuditDismiss = "Dismiss one or more audit reports (or --all)"
	shortAuditRelay   = "Relay unread audit reports at the next prompt"
)

// audit example-usage blocks.
const (
	exampleAudit = "  ctxctl audit\n" +
		"  ctxctl audit show surface\n" +
		"  ctxctl audit dismiss surface"
	exampleAuditList    = "  ctxctl audit list"
	exampleAuditShow    = "  ctxctl audit show surface"
	exampleAuditDismiss = "  ctxctl audit dismiss surface\n" +
		"  ctxctl audit dismiss --all"
)

// --all flag description.
const flagAuditDismissAll = "Dismiss every audit report in .context/audit/"

// audit list / dismiss output format strings.
const (
	writeAuditNone         = "No audit reports."
	writeAuditListItem     = "%s  %-10s  %s  generated %s"
	writeAuditDismissed    = "Dismissed audit %s."
	writeAuditDismissedAll = "Dismissed %d audit report(s)."
)

// audit-relay box copy and provenance labels.
const (
	relayLabel           = "audit-relay"
	relayVariant         = "audits"
	relayBoxTitle        = "Audit Reports"
	relayDismissHint     = "Dismiss: ctxctl audit dismiss <id>"
	relayDismissAllHint  = "Dismiss all: ctxctl audit dismiss --all"
	relayNudgeFormat     = "You have %d unread audit report(s) at .context/audit/"
	relayPrefix          = "IMPORTANT: Relay these audit findings to the user VERBATIM before continuing."
	relayReportSeparator = "─────────────────────────────────"
	relayStalePrefix     = "STALE — %s (audited %s ago)"
)

// auditStrings assembles the English text for the audit
// command tree.
func auditStrings() audit.Strings {
	return audit.Strings{
		Use:     useAudit,
		Short:   shortAudit,
		Long:    longAudit,
		Example: exampleAudit,
		List: list.Strings{
			Use:      useAuditList,
			Short:    shortAuditList,
			Example:  exampleAuditList,
			None:     writeAuditNone,
			ListItem: writeAuditListItem,
		},
		Show: show.Strings{
			Use:     useAuditShow,
			Short:   shortAuditShow,
			Example: exampleAuditShow,
		},
		Dismiss: dismiss.Strings{
			Use:          useAuditDismiss,
			Short:        shortAuditDismiss,
			Example:      exampleAuditDismiss,
			AllFlag:      flagAuditDismissAll,
			Dismissed:    writeAuditDismissed,
			DismissedAll: writeAuditDismissedAll,
		},
	}
}

// relayStrings assembles the English text for the
// audit-relay hook and its box.
func relayStrings() checkaudit.Strings {
	return checkaudit.Strings{
		Use:             useAuditRelay,
		Short:           shortAuditRelay,
		Example:         "",
		RelayLabel:      relayLabel,
		RelayVariant:    relayVariant,
		BoxTitle:        relayBoxTitle,
		RelayPrefix:     relayPrefix,
		NudgeFormat:     relayNudgeFormat,
		DismissHint:     relayDismissHint,
		DismissAllHint:  relayDismissAllHint,
		ReportSeparator: relayReportSeparator,
		StalePrefix:     relayStalePrefix,
	}
}
