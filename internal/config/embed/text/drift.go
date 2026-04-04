//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for drift detection.
const (
	// DescKeyDriftDeadPath is the text key for drift dead path messages.
	DescKeyDriftDeadPath = "drift.dead-path"
	// DescKeyDriftEntryCount is the text key for drift entry count messages.
	DescKeyDriftEntryCount = "drift.entry-count"
	// DescKeyDriftMissingFile is the text key for drift missing file messages.
	DescKeyDriftMissingFile = "drift.missing-file"
	// DescKeyDriftRegenerated is the text key for drift regenerated messages.
	DescKeyDriftRegenerated = "drift.regenerated"
	// DescKeyDriftMissingPackage is the text key for drift missing package
	// messages.
	DescKeyDriftMissingPackage = "drift.missing-package"
	// DescKeyDriftSecret is the text key for drift secret messages.
	DescKeyDriftSecret = "drift.secret"
	// DescKeyDriftStaleAge is the text key for drift stale age messages.
	DescKeyDriftStaleAge = "drift.stale-age"
	// DescKeyDriftStaleness is the text key for drift staleness messages.
	DescKeyDriftStaleness = "drift.staleness"
	// DescKeyDriftCleared is the text key for drift cleared messages.
	DescKeyDriftCleared = "drift.cleared"
	// DescKeyDriftApplying is the text key for drift applying messages.
	DescKeyDriftApplying = "drift.applying"
	// DescKeyDriftFixedCount is the text key for drift fixed count messages.
	DescKeyDriftFixedCount = "drift.fixed-count"
	// DescKeyDriftSkippedCount is the text key for drift skipped count messages.
	DescKeyDriftSkippedCount = "drift.skipped-count"
	// DescKeyDriftFixError is the text key for drift fix error messages.
	DescKeyDriftFixError = "drift.fix-error"
	// DescKeyDriftRechecking is the text key for drift rechecking messages.
	DescKeyDriftRechecking = "drift.rechecking"
	// DescKeyDriftFixStaleness is the text key for drift fix staleness messages.
	DescKeyDriftFixStaleness = "drift.fix-staleness"
	// DescKeyDriftFixStalenessErr is the text key for drift fix staleness err
	// messages.
	DescKeyDriftFixStalenessErr = "drift.fix-staleness-err"
	// DescKeyDriftFixMissing is the text key for drift fix missing messages.
	DescKeyDriftFixMissing = "drift.fix-missing"
	// DescKeyDriftFixMissingErr is the text key for drift fix missing err
	// messages.
	DescKeyDriftFixMissingErr = "drift.fix-missing-err"
	// DescKeyDriftSkipDeadPath is the text key for drift skip dead path messages.
	DescKeyDriftSkipDeadPath = "drift.skip-dead-path"
	// DescKeyDriftSkipStaleAge is the text key for drift skip stale age messages.
	DescKeyDriftSkipStaleAge = "drift.skip-stale-age"
	// DescKeyDriftSkipSensitiveFile is the text key for drift skip sensitive file
	// messages.
	DescKeyDriftSkipSensitiveFile = "drift.skip-sensitive-file"
	// DescKeyDriftArchived is the text key for drift archived messages.
	DescKeyDriftArchived = "drift.archived"
	// DescKeyDriftReportHeading is the text key for drift report heading messages.
	DescKeyDriftReportHeading = "drift.report-heading"
	// DescKeyDriftReportSeparator is the text key for drift report separator
	// messages.
	DescKeyDriftReportSeparator = "drift.report-separator"
	// DescKeyDriftViolationsHeading is the text key for drift violations heading
	// messages.
	DescKeyDriftViolationsHeading = "drift.violations-heading"
	// DescKeyDriftViolationLine is the text key for drift violation line messages.
	DescKeyDriftViolationLine = "drift.violation-line"
	// DescKeyDriftViolationLineLoc is the text key for drift violation line loc
	// messages.
	DescKeyDriftViolationLineLoc = "drift.violation-line-loc"
	// DescKeyDriftViolationRule is the text key for drift violation rule messages.
	DescKeyDriftViolationRule = "drift.violation-rule"
	// DescKeyDriftWarningsHeading is the text key for drift warnings heading
	// messages.
	DescKeyDriftWarningsHeading = "drift.warnings-heading"
	// DescKeyDriftPathRefsLabel is the text key for drift path refs label
	// messages.
	DescKeyDriftPathRefsLabel = "drift.path-refs-label"
	// DescKeyDriftPathRefLine is the text key for drift path ref line messages.
	DescKeyDriftPathRefLine = "drift.path-ref-line"
	// DescKeyDriftStalenessLabel is the text key for drift staleness label
	// messages.
	DescKeyDriftStalenessLabel = "drift.staleness-label"
	// DescKeyDriftStalenessLine is the text key for drift staleness line messages.
	DescKeyDriftStalenessLine = "drift.staleness-line"
	// DescKeyDriftOtherLabel is the text key for drift other label messages.
	DescKeyDriftOtherLabel = "drift.other-label"
	// DescKeyDriftOtherLine is the text key for drift other line messages.
	DescKeyDriftOtherLine = "drift.other-line"
	// DescKeyDriftPassedHeading is the text key for drift passed heading messages.
	DescKeyDriftPassedHeading = "drift.passed-heading"
	// DescKeyDriftPassedLine is the text key for drift passed line messages.
	DescKeyDriftPassedLine = "drift.passed-line"
	// DescKeyDriftStatusViolation is the text key for drift status violation
	// messages.
	DescKeyDriftStatusViolation = "drift.status-violation"
	// DescKeyDriftStatusWarning is the text key for drift status warning messages.
	DescKeyDriftStatusWarning = "drift.status-warning"
	// DescKeyDriftStatusOK is the text key for drift status ok messages.
	DescKeyDriftStatusOK = "drift.status-ok"
	// DescKeyDriftCheckPathRefs is the text key for drift check path refs
	// messages.
	DescKeyDriftCheckPathRefs = "drift.check-path-refs"
	// DescKeyDriftCheckStaleness is the text key for drift check staleness
	// messages.
	DescKeyDriftCheckStaleness = "drift.check-staleness"
	// DescKeyDriftCheckConstitution is the text key for drift check constitution
	// messages.
	DescKeyDriftCheckConstitution = "drift.check-constitution"
	// DescKeyDriftCheckRequired is the text key for drift check required messages.
	DescKeyDriftCheckRequired = "drift.check-required"
	// DescKeyDriftCheckFileAge is the text key for drift check file age messages.
	DescKeyDriftCheckFileAge = "drift.check-file-age"
	// DescKeyDriftStaleHeader is the text key for drift stale header messages.
	DescKeyDriftStaleHeader = "drift.stale-header"
	// DescKeyDriftCheckTemplateHeader is the text key for drift check template
	// header messages.
	DescKeyDriftCheckTemplateHeader = "drift.check-template-header"
	// DescKeyDriftInvalidTool is the text key for drift invalid tool messages.
	DescKeyDriftInvalidTool = "drift.invalid-tool"
	// DescKeyDriftHookNoExec is the text key for drift hook no exec messages.
	DescKeyDriftHookNoExec = "drift.hook-no-exec"
	// DescKeyDriftStaleSyncFile is the text key for drift stale sync file
	// messages.
	DescKeyDriftStaleSyncFile = "drift.stale-sync-file"
	// DescKeyDriftToolSuffix is the text key for drift tool suffix messages.
	DescKeyDriftToolSuffix = "drift.tool-suffix"
	// DescKeyVersionDriftRelayMessage is the text key for version drift relay
	// message messages.
	DescKeyVersionDriftRelayMessage = "version-drift.relay-message"
	// DescKeyWriteVersionDriftFallback is the text key for write version drift
	// fallback messages.
	DescKeyWriteVersionDriftFallback = "write.version-drift-fallback"
)
