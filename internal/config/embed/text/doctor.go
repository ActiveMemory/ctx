//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for doctor diagnostics.
const (
	// DescKeyDoctorCheckDidNotRun is the text key for the failure
	// result emitted by the doctor runner when a check returns an
	// error it could not handle itself.
	DescKeyDoctorCheckDidNotRun = "doctor.check-did-not-run"
	// DescKeyDoctorCheckDidNotRunCascade is emitted once for the
	// first context-dependent check that fails with
	// [errCtx.ErrDirNotDeclared]; later dependent checks are
	// silently skipped so the report shows one loud line instead
	// of the same message N times.
	DescKeyDoctorCheckDidNotRunCascade = "doctor.check-did-not-run-cascade"
	// DescKeyDoctorContextFileFormat is the text key for doctor context file
	// format messages.
	DescKeyDoctorContextFileFormat = "doctor.context-file.format"
	// DescKeyDoctorContextInitializedError is the text key for doctor context
	// initialized error messages.
	DescKeyDoctorContextInitializedError = "doctor.context-initialized.error"
	// DescKeyDoctorContextInitializedOk is the text key for doctor context
	// initialized ok messages.
	DescKeyDoctorContextInitializedOk = "doctor.context-initialized.ok"
	// DescKeyDoctorContextSizeFormat is the text key for doctor context size
	// format messages.
	DescKeyDoctorContextSizeFormat = "doctor.context-size.format"
	// DescKeyDoctorContextSizeWarningSuffix is the text key for doctor context
	// size warning suffix messages.
	DescKeyDoctorContextSizeWarningSuffix = "doctor.context-size.warning-suffix"
	// DescKeyDoctorCtxrcValidationError is the text key for doctor ctxrc
	// validation error messages.
	DescKeyDoctorCtxrcValidationError = "doctor.ctxrc-validation.error"
	// DescKeyDoctorCtxrcValidationOk is the text key for doctor ctxrc validation
	// ok messages.
	DescKeyDoctorCtxrcValidationOk = "doctor.ctxrc-validation.ok"
	// DescKeyDoctorCtxrcValidationOkNoFile is the text key for doctor ctxrc
	// validation ok no file messages.
	DescKeyDoctorCtxrcValidationOkNoFile = "doctor.ctxrc-validation.ok-no-file"
	// DescKeyDoctorCtxrcValidationWarning is the text key for doctor ctxrc
	// validation warning messages.
	DescKeyDoctorCtxrcValidationWarning = "doctor.ctxrc-validation.warning"
	// DescKeyDoctorDriftDetected is the text key for doctor drift detected
	// messages.
	DescKeyDoctorDriftDetected = "doctor.drift.detected"
	// DescKeyDoctorDriftOk is the text key for doctor drift ok messages.
	DescKeyDoctorDriftOk = "doctor.drift.ok"
	// DescKeyDoctorDriftViolations is the text key for doctor drift violations
	// messages.
	DescKeyDoctorDriftViolations = "doctor.drift.violations"
	// DescKeyDoctorDriftWarningLoad is the text key for doctor drift warning load
	// messages.
	DescKeyDoctorDriftWarningLoad = "doctor.drift.warning-load"
	// DescKeyDoctorDriftWarnings is the text key for doctor drift warnings
	// messages.
	DescKeyDoctorDriftWarnings = "doctor.drift.warnings"
	// DescKeyDoctorEventLoggingInfo is the text key for doctor event logging info
	// messages.
	DescKeyDoctorEventLoggingInfo = "doctor.event-logging.info"
	// DescKeyDoctorEventLoggingOk is the text key for doctor event logging ok
	// messages.
	DescKeyDoctorEventLoggingOk = "doctor.event-logging.ok"
	// DescKeyDoctorOutputHeader is the text key for doctor output header messages.
	DescKeyDoctorOutputHeader = "doctor.output.header"
	// DescKeyDoctorOutputResultLine is the text key for doctor output result line
	// messages.
	DescKeyDoctorOutputResultLine = "doctor.output.result-line"
	// DescKeyDoctorOutputSeparator is the text key for doctor output separator
	// messages.
	DescKeyDoctorOutputSeparator = "doctor.output.separator"
	// DescKeyDoctorOutputSummary is the text key for doctor output summary
	// messages.
	DescKeyDoctorOutputSummary = "doctor.output.summary"
	// DescKeyDoctorCompanionConfigOk is the text key for doctor companion config
	// ok messages.
	DescKeyDoctorCompanionConfigOk = "doctor.companion-config.ok"
	// DescKeyDoctorCompanionConfigInfo is the text key for doctor companion
	// config info messages.
	DescKeyDoctorCompanionConfigInfo = "doctor.companion-config.info"
	// DescKeyDoctorPluginEnabledGlobalOk is the text key for doctor plugin
	// enabled global ok messages.
	DescKeyDoctorPluginEnabledGlobalOk = "doctor.plugin-enabled-global.ok"
	// DescKeyDoctorPluginEnabledLocalOk is the text key for doctor plugin enabled
	// local ok messages.
	DescKeyDoctorPluginEnabledLocalOk = "doctor.plugin-enabled-local.ok"
	// DescKeyDoctorPluginEnabledWarning is the text key for doctor plugin enabled
	// warning messages.
	DescKeyDoctorPluginEnabledWarning = "doctor.plugin-enabled.warning"
	// DescKeyDoctorPluginInstalledInfo is the text key for doctor plugin
	// installed info messages.
	DescKeyDoctorPluginInstalledInfo = "doctor.plugin-installed.info"
	// DescKeyDoctorPluginInstalledOk is the text key for doctor plugin installed
	// ok messages.
	DescKeyDoctorPluginInstalledOk = "doctor.plugin-installed.ok"
	// DescKeyDoctorRecentEventsInfo is the text key for doctor recent events info
	// messages.
	DescKeyDoctorRecentEventsInfo = "doctor.recent-events.info"
	// DescKeyDoctorRecentEventsOk is the text key for doctor recent events ok
	// messages.
	DescKeyDoctorRecentEventsOk = "doctor.recent-events.ok"
	// DescKeyDoctorRemindersInfo is the text key for doctor reminders info
	// messages.
	DescKeyDoctorRemindersInfo = "doctor.reminders.info"
	// DescKeyDoctorRemindersOk is the text key for doctor reminders ok messages.
	DescKeyDoctorRemindersOk = "doctor.reminders.ok"
	// DescKeyDoctorRequiredFilesError is the text key for doctor required files
	// error messages.
	DescKeyDoctorRequiredFilesError = "doctor.required-files.error"
	// DescKeyDoctorRequiredFilesOk is the text key for doctor required files ok
	// messages.
	DescKeyDoctorRequiredFilesOk = "doctor.required-files.ok"
	// DescKeyDoctorResourceDiskFormat is the text key for doctor resource disk
	// format messages.
	DescKeyDoctorResourceDiskFormat = "doctor.resource-disk.format"
	// DescKeyDoctorResourceLoadFormat is the text key for doctor resource load
	// format messages.
	DescKeyDoctorResourceLoadFormat = "doctor.resource-load.format"
	// DescKeyDoctorResourceMemoryFormat is the text key for doctor resource
	// memory format messages.
	DescKeyDoctorResourceMemoryFormat = "doctor.resource-memory.format"
	// DescKeyDoctorResourceSwapFormat is the text key for doctor resource swap
	// format messages.
	DescKeyDoctorResourceSwapFormat = "doctor.resource-swap.format"
	// DescKeyDoctorTaskCompletionFormat is the text key for doctor task
	// completion format messages.
	DescKeyDoctorTaskCompletionFormat = "doctor.task-completion.format"
	// DescKeyDoctorTaskCompletionWarningSuffix is the text key for doctor task
	// completion warning suffix messages.
	DescKeyDoctorTaskCompletionWarningSuffix = "doctor.task-completion.warning-suffix"
	// DescKeyDoctorWebhookInfo is the text key for doctor webhook info messages.
	DescKeyDoctorWebhookInfo = "doctor.webhook.info"
	// DescKeyDoctorWebhookOk is the text key for doctor webhook ok messages.
	DescKeyDoctorWebhookOk = "doctor.webhook.ok"
)
