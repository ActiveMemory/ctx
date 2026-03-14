//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package assets provides embedded assets for ctx: .context/ templates
// stamped by "ctx init" and the Claude Code plugin (skills, hooks,
// manifest) served directly from claude/.
package assets

import (
	"embed"
	"encoding/json"
	"strings"
	"sync"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"gopkg.in/yaml.v3"
)

//go:embed claude/.claude-plugin/plugin.json claude/CLAUDE.md claude/skills/*/references/*.md claude/skills/*/SKILL.md context/*.md project/* entry-templates/*.md hooks/*.md hooks/messages/*/*.txt hooks/messages/registry.yaml prompt-templates/*.md ralph/*.md schema/*.json why/*.md permissions/*.txt commands/*.yaml journal/*.css
var FS embed.FS

const (
	FlagDescKeyAddApplication              = "add.application"
	FlagDescKeyAddConsequences             = "add.consequences"
	FlagDescKeyAddContext                  = "add.context"
	FlagDescKeyAddFile                     = "add.file"
	FlagDescKeyAddLesson                   = "add.lesson"
	FlagDescKeyAddPriority                 = "add.priority"
	FlagDescKeyAddRationale                = "add.rationale"
	FlagDescKeyAddSection                  = "add.section"
	FlagDescKeyAgentBudget                 = "agent.budget"
	FlagDescKeyAgentCooldown               = "agent.cooldown"
	FlagDescKeyAgentFormat                 = "agent.format"
	FlagDescKeyAgentSession                = "agent.session"
	FlagDescKeyChangesSince                = "changes.since"
	FlagDescKeyCompactArchive              = "compact.archive"
	FlagDescKeyDepsExternal                = "deps.external"
	FlagDescKeyDepsFormat                  = "deps.format"
	FlagDescKeyDepsType                    = "deps.type"
	FlagDescKeyDoctorJson                  = "doctor.json"
	FlagDescKeyDriftFix                    = "drift.fix"
	FlagDescKeyDriftJson                   = "drift.json"
	FlagDescKeyGuideCommands               = "guide.commands"
	FlagDescKeyGuideSkills                 = "guide.skills"
	FlagDescKeyHookWrite                   = "hook.write"
	FlagDescKeyInitializeForce             = "initialize.force"
	FlagDescKeyInitializeMerge             = "initialize.merge"
	FlagDescKeyInitializeMinimal           = "initialize.minimal"
	FlagDescKeyInitializeNoPluginEnable    = "initialize.no-plugin-enable"
	FlagDescKeyInitializeRalph             = "initialize.ralph"
	FlagDescKeyJournalObsidianOutput       = "journal.obsidian.output"
	FlagDescKeyJournalSiteBuild            = "journal.site.build"
	FlagDescKeyJournalSiteOutput           = "journal.site.output"
	FlagDescKeyJournalSiteServe            = "journal.site.serve"
	FlagDescKeyLoadBudget                  = "load.budget"
	FlagDescKeyLoadRaw                     = "load.raw"
	FlagDescKeyLoopCompletion              = "loop.completion"
	FlagDescKeyLoopMaxIterations           = "loop.max-iterations"
	FlagDescKeyLoopOutput                  = "loop.output"
	FlagDescKeyLoopPrompt                  = "loop.prompt"
	FlagDescKeyLoopTool                    = "loop.tool"
	FlagDescKeyMemoryImportDryRun          = "memory.import.dry-run"
	FlagDescKeyMemoryPublishBudget         = "memory.publish.budget"
	FlagDescKeyMemoryPublishDryRun         = "memory.publish.dry-run"
	FlagDescKeyMemorySyncDryRun            = "memory.sync.dry-run"
	FlagDescKeyNotifyEvent                 = "notify.event"
	FlagDescKeyNotifyHook                  = "notify.hook"
	FlagDescKeyNotifySessionId             = "notify.session-id"
	FlagDescKeyNotifyVariant               = "notify.variant"
	FlagDescKeyPadAddFile                  = "pad.add.file"
	FlagDescKeyPadEditAppend               = "pad.edit.append"
	FlagDescKeyPadEditFile                 = "pad.edit.file"
	FlagDescKeyPadEditLabel                = "pad.edit.label"
	FlagDescKeyPadEditPrepend              = "pad.edit.prepend"
	FlagDescKeyPadExportDryRun             = "pad.export.dry-run"
	FlagDescKeyPadExportForce              = "pad.export.force"
	FlagDescKeyPadImpBlobs                 = "pad.imp.blobs"
	FlagDescKeyPadMergeDryRun              = "pad.merge.dry-run"
	FlagDescKeyPadMergeKey                 = "pad.merge.key"
	FlagDescKeyPadShowOut                  = "pad.show.out"
	FlagDescKeyPauseSessionId              = "pause.session-id"
	FlagDescKeyPromptAddStdin              = "prompt.add.stdin"
	FlagDescKeyRecallExportAll             = "recall.export.all"
	FlagDescKeyRecallExportAllProjects     = "recall.export.all-projects"
	FlagDescKeyRecallExportDryRun          = "recall.export.dry-run"
	FlagDescKeyRecallExportKeepFrontmatter = "recall.export.keep-frontmatter"
	FlagDescKeyRecallExportRegenerate      = "recall.export.regenerate"
	FlagDescKeyRecallExportSkipExisting    = "recall.export.skip-existing"
	FlagDescKeyRecallExportYes             = "recall.export.yes"
	FlagDescKeyRecallListAllProjects       = "recall.list.all-projects"
	FlagDescKeyRecallListLimit             = "recall.list.limit"
	FlagDescKeyRecallListProject           = "recall.list.project"
	FlagDescKeyRecallListSince             = "recall.list.since"
	FlagDescKeyRecallListTool              = "recall.list.tool"
	FlagDescKeyRecallListUntil             = "recall.list.until"
	FlagDescKeyRecallLockAll               = "recall.lock.all"
	FlagDescKeyRecallShowAllProjects       = "recall.show.all-projects"
	FlagDescKeyRecallShowFull              = "recall.show.full"
	FlagDescKeyRecallShowLatest            = "recall.show.latest"
	FlagDescKeyRecallUnlockAll             = "recall.unlock.all"
	FlagDescKeyRemindAddAfter              = "remind.add.after"
	FlagDescKeyRemindAfter                 = "remind.after"
	FlagDescKeyRemindDismissAll            = "remind.dismiss.all"
	FlagDescKeyResumeSessionId             = "resume.session-id"
	FlagDescKeySiteFeedBaseUrl             = "site.feed.base-url"
	FlagDescKeySiteFeedOut                 = "site.feed.out"
	FlagDescKeyStatusJson                  = "status.json"
	FlagDescKeyStatusVerbose               = "status.verbose"
	FlagDescKeySyncDryRun                  = "sync.dry-run"
	FlagDescKeySystemBackupJson            = "system.backup.json"
	FlagDescKeySystemBackupScope           = "system.backup.scope"
	FlagDescKeySystemBootstrapJson         = "system.bootstrap.json"
	FlagDescKeySystemBootstrapQuiet        = "system.bootstrap.quiet"
	FlagDescKeySystemEventsAll             = "system.events.all"
	FlagDescKeySystemEventsEvent           = "system.events.event"
	FlagDescKeySystemEventsHook            = "system.events.hook"
	FlagDescKeySystemEventsJson            = "system.events.json"
	FlagDescKeySystemEventsLast            = "system.events.last"
	FlagDescKeySystemEventsSession         = "system.events.session"
	FlagDescKeySystemMarkjournalCheck      = "system.markjournal.check"
	FlagDescKeySystemMessageJson           = "system.message.json"
	FlagDescKeySystemPauseSessionId        = "system.pause.session-id"
	FlagDescKeySystemPruneDays             = "system.prune.days"
	FlagDescKeySystemPruneDryRun           = "system.prune.dry-run"
	FlagDescKeySystemResourcesJson         = "system.resources.json"
	FlagDescKeySystemResumeSessionId       = "system.resume.session-id"
	FlagDescKeySystemStatsFollow           = "system.stats.follow"
	FlagDescKeySystemStatsJson             = "system.stats.json"
	FlagDescKeySystemStatsLast             = "system.stats.last"
	FlagDescKeySystemStatsSession          = "system.stats.session"
	FlagDescKeyTaskArchiveDryRun           = "task.archive.dry-run"
	FlagDescKeyWatchDryRun                 = "watch.dry-run"
	FlagDescKeyWatchLog                    = "watch.log"
)

const (
	CmdDescKeyAdd                          = "add"
	CmdDescKeyAgent                        = "agent"
	CmdDescKeyChanges                      = "changes"
	CmdDescKeyCompact                      = "compact"
	CmdDescKeyComplete                     = "complete"
	CmdDescKeyConfig                       = "config"
	CmdDescKeyConfigSchema                 = "config.schema"
	CmdDescKeyConfigStatus                 = "config.status"
	CmdDescKeyConfigSwitch                 = "config.switch"
	CmdDescKeyCtx                          = "ctx"
	CmdDescKeyDecision                     = "decision"
	CmdDescKeyDecisionReindex              = "decision.reindex"
	CmdDescKeyDeps                         = "deps"
	CmdDescKeyDoctor                       = "doctor"
	CmdDescKeyDrift                        = "drift"
	CmdDescKeyGuide                        = "guide"
	CmdDescKeyHook                         = "hook"
	CmdDescKeyInitialize                   = "initialize"
	CmdDescKeyJournal                      = "journal"
	CmdDescKeyJournalObsidian              = "journal.obsidian"
	CmdDescKeyJournalSite                  = "journal.site"
	CmdDescKeyLearnings                    = "learnings"
	CmdDescKeyLearningsReindex             = "learnings.reindex"
	CmdDescKeyLoad                         = "load"
	CmdDescKeyLoop                         = "loop"
	CmdDescKeyMcp                          = "mcp"
	CmdDescKeyMcpServe                     = "mcp.serve"
	CmdDescKeyMemory                       = "memory"
	CmdDescKeyMemoryDiff                   = "memory.diff"
	CmdDescKeyMemoryImport                 = "memory.import"
	CmdDescKeyMemoryPublish                = "memory.publish"
	CmdDescKeyMemoryStatus                 = "memory.status"
	CmdDescKeyMemorySync                   = "memory.sync"
	CmdDescKeyMemoryUnpublish              = "memory.unpublish"
	CmdDescKeyNotify                       = "notify"
	CmdDescKeyNotifySetup                  = "notify.setup"
	CmdDescKeyNotifyTest                   = "notify.test"
	CmdDescKeyPad                          = "pad"
	CmdDescKeyPadAdd                       = "pad.add"
	CmdDescKeyPadEdit                      = "pad.edit"
	CmdDescKeyPadExport                    = "pad.export"
	CmdDescKeyPadImp                       = "pad.imp"
	CmdDescKeyPadMerge                     = "pad.merge"
	CmdDescKeyPadMv                        = "pad.mv"
	CmdDescKeyPadResolve                   = "pad.resolve"
	CmdDescKeyPadRm                        = "pad.rm"
	CmdDescKeyPadShow                      = "pad.show"
	CmdDescKeyPause                        = "pause"
	CmdDescKeyPermissions                  = "permissions"
	CmdDescKeyPermissionsRestore           = "permissions.restore"
	CmdDescKeyPermissionsSnapshot          = "permissions.snapshot"
	CmdDescKeyPrompt                       = "prompt"
	CmdDescKeyPromptAdd                    = "prompt.add"
	CmdDescKeyPromptList                   = "prompt.list"
	CmdDescKeyPromptRm                     = "prompt.rm"
	CmdDescKeyPromptShow                   = "prompt.show"
	CmdDescKeyRecall                       = "recall"
	CmdDescKeyRecallExport                 = "recall.export"
	CmdDescKeyRecallList                   = "recall.list"
	CmdDescKeyRecallLock                   = "recall.lock"
	CmdDescKeyRecallShow                   = "recall.show"
	CmdDescKeyRecallSync                   = "recall.sync"
	CmdDescKeyRecallUnlock                 = "recall.unlock"
	CmdDescKeyReindex                      = "reindex"
	CmdDescKeyRemind                       = "remind"
	CmdDescKeyRemindAdd                    = "remind.add"
	CmdDescKeyRemindDismiss                = "remind.dismiss"
	CmdDescKeyRemindList                   = "remind.list"
	CmdDescKeyResume                       = "resume"
	CmdDescKeyServe                        = "serve"
	CmdDescKeySite                         = "site"
	CmdDescKeySiteFeed                     = "site.feed"
	CmdDescKeyStatus                       = "status"
	CmdDescKeySync                         = "sync"
	CmdDescKeySystem                       = "system"
	CmdDescKeySystemBackup                 = "system.backup"
	CmdDescKeySystemBlockDangerousCommands = "system.blockdangerouscommands"
	CmdDescKeySystemBlockNonPathCtx        = "system.blocknonpathctx"
	CmdDescKeySystemBootstrap              = "system.bootstrap"
	CmdDescKeySystemCheckBackupAge         = "system.checkbackupage"
	CmdDescKeySystemCheckCeremonies        = "system.checkceremonies"
	CmdDescKeySystemCheckContextSize       = "system.checkcontextsize"
	CmdDescKeySystemCheckJournal           = "system.checkjournal"
	CmdDescKeySystemCheckKnowledge         = "system.checkknowledge"
	CmdDescKeySystemCheckMapStaleness      = "system.checkmapstaleness"
	CmdDescKeySystemCheckMemoryDrift       = "system.checkmemorydrift"
	CmdDescKeySystemCheckPersistence       = "system.checkpersistence"
	CmdDescKeySystemCheckReminders         = "system.checkreminders"
	CmdDescKeySystemCheckResources         = "system.checkresources"
	CmdDescKeySystemCheckTaskCompletion    = "system.checktaskcompletion"
	CmdDescKeySystemCheckVersion           = "system.checkversion"
	CmdDescKeySystemContextLoadGate        = "system.contextloadgate"
	CmdDescKeySystemEvents                 = "system.events"
	CmdDescKeySystemHeartbeat              = "system.heartbeat"
	CmdDescKeySystemMarkJournal            = "system.markjournal"
	CmdDescKeySystemMarkWrappedUp          = "system.markwrappedup"
	CmdDescKeySystemMessage                = "system.message"
	CmdDescKeySystemMessageEdit            = "system.message.edit"
	CmdDescKeySystemMessageList            = "system.message.list"
	CmdDescKeySystemMessageReset           = "system.message.reset"
	CmdDescKeySystemMessageShow            = "system.message.show"
	CmdDescKeySystemPause                  = "system.pause"
	CmdDescKeySystemPostCommit             = "system.postcommit"
	CmdDescKeySystemPrune                  = "system.prune"
	CmdDescKeySystemQaReminder             = "system.qareminder"
	CmdDescKeySystemResources              = "system.resources"
	CmdDescKeySystemResume                 = "system.resume"
	CmdDescKeySystemSpecsNudge             = "system.specsnudge"
	CmdDescKeySystemStats                  = "system.stats"
	CmdDescKeyTask                         = "task"
	CmdDescKeyTaskArchive                  = "task.archive"
	CmdDescKeyTaskSnapshot                 = "task.snapshot"
	CmdDescKeyWatch                        = "watch"
	CmdDescKeyWhy                          = "why"
)

const (
	TextDescKeyAgentInstruction                 = "agent.instruction"
	TextDescKeyBackupBoxTitle                   = "backup.box-title"
	TextDescKeyBackupNoMarker                   = "backup.no-marker"
	TextDescKeyBackupRelayMessage               = "backup.relay-message"
	TextDescKeyBackupRelayPrefix                = "backup.relay-prefix"
	TextDescKeyBackupRunHint                    = "backup.run-hint"
	TextDescKeyBackupSMBNotMounted              = "backup.smb-not-mounted"
	TextDescKeyBackupSMBUnavailable             = "backup.smb-unavailable"
	TextDescKeyBackupStale                      = "backup.stale"
	TextDescKeyBootstrapNextSteps               = "bootstrap.next-steps"
	TextDescKeyBootstrapNone                    = "bootstrap.none"
	TextDescKeyBootstrapPluginWarning           = "bootstrap.plugin-warning"
	TextDescKeyBootstrapRules                   = "bootstrap.rules"
	TextDescKeyContextLoadGateFileHeader        = "context-load-gate.file-header"
	TextDescKeyContextLoadGateFooter            = "context-load-gate.footer"
	TextDescKeyContextLoadGateHeader            = "context-load-gate.header"
	TextDescKeyContextLoadGateIndexFallback     = "context-load-gate.index-fallback"
	TextDescKeyContextLoadGateIndexHeader       = "context-load-gate.index-header"
	TextDescKeyContextLoadGateOversizeAction    = "context-load-gate.oversize-action"
	TextDescKeyContextLoadGateOversizeBreakdown = "context-load-gate.oversize-breakdown"
	TextDescKeyContextLoadGateOversizeFileEntry = "context-load-gate.oversize-file-entry"
	TextDescKeyContextLoadGateOversizeHeader    = "context-load-gate.oversize-header"
	TextDescKeyContextLoadGateOversizeInjected  = "context-load-gate.oversize-injected"
	TextDescKeyContextLoadGateOversizeTimestamp = "context-load-gate.oversize-timestamp"
	TextDescKeyContextLoadGateWebhook           = "context-load-gate.webhook"

	TextDescKeyHeartbeatLogTokens    = "heartbeat.log-tokens"
	TextDescKeyHeartbeatLogPlain     = "heartbeat.log-plain"
	TextDescKeyHeartbeatNotifyTokens = "heartbeat.notify-tokens"
	TextDescKeyHeartbeatNotifyPlain  = "heartbeat.notify-plain"

	TextDescKeyEventsEmpty       = "events.empty"
	TextDescKeyEventsHumanFormat = "events.human-format"

	TextDescKeyStatsEmpty        = "stats.empty"
	TextDescKeyStatsHeaderFormat = "stats.header-format"
	TextDescKeyStatsLineFormat   = "stats.line-format"

	TextDescKeyCheckContextSizeBillingBoxTitle       = "check-context-size.billing-box-title"
	TextDescKeyCheckContextSizeBillingFallback       = "check-context-size.billing-fallback"
	TextDescKeyCheckContextSizeBillingRelayFormat    = "check-context-size.billing-relay-format"
	TextDescKeyCheckContextSizeBillingRelayPrefix    = "check-context-size.billing-relay-prefix"
	TextDescKeyCheckContextSizeCheckpointBoxTitle    = "check-context-size.checkpoint-box-title"
	TextDescKeyCheckContextSizeCheckpointFallback    = "check-context-size.checkpoint-fallback"
	TextDescKeyCheckContextSizeCheckpointRelayFormat = "check-context-size.checkpoint-relay-format"
	TextDescKeyCheckContextSizeOversizeFallback      = "check-context-size.oversize-fallback"
	TextDescKeyCheckContextSizeRelayPrefix           = "check-context-size.relay-prefix"
	TextDescKeyCheckContextSizeRunningLowSuffix      = "check-context-size.running-low-suffix"
	TextDescKeyCheckContextSizeSilentLogFormat       = "check-context-size.silent-log-format"
	TextDescKeyCheckContextSizeSilencedCheckpointLog = "check-context-size.silenced-checkpoint-log"
	TextDescKeyCheckContextSizeCheckpointLogFormat   = "check-context-size.checkpoint-log-format"
	TextDescKeyCheckContextSizeSuppressedLogFormat   = "check-context-size.suppressed-log-format"
	TextDescKeyCheckContextSizeSilencedWindowLog     = "check-context-size.silenced-window-log"
	TextDescKeyCheckContextSizeWindowLogFormat       = "check-context-size.window-log-format"
	TextDescKeyCheckContextSizeSilencedBillingLog    = "check-context-size.silenced-billing-log"
	TextDescKeyCheckContextSizeBillingLogFormat      = "check-context-size.billing-log-format"
	TextDescKeyCheckContextSizeTokenLow              = "check-context-size.token-low"
	TextDescKeyCheckContextSizeTokenNormal           = "check-context-size.token-normal"
	TextDescKeyCheckContextSizeTokenUsage            = "check-context-size.token-usage"
	TextDescKeyCheckContextSizeWindowBoxTitle        = "check-context-size.window-box-title"
	TextDescKeyCheckContextSizeWindowFallback        = "check-context-size.window-fallback"
	TextDescKeyCheckContextSizeWindowRelayFormat     = "check-context-size.window-relay-format"
	TextDescKeyCheckJournalBoxTitle                  = "check-journal.box-title"
	TextDescKeyCheckJournalFallbackBoth              = "check-journal.fallback-both"
	TextDescKeyCheckJournalFallbackUnenriched        = "check-journal.fallback-unenriched"
	TextDescKeyCheckJournalFallbackUnexported        = "check-journal.fallback-unexported"
	TextDescKeyCheckJournalRelayFormat               = "check-journal.relay-format"
	TextDescKeyCheckJournalRelayPrefix               = "check-journal.relay-prefix"
	TextDescKeyCheckKnowledgeBoxTitle                = "check-knowledge.box-title"
	TextDescKeyCheckKnowledgeFallback                = "check-knowledge.fallback"
	TextDescKeyCheckKnowledgeFindingFormat           = "check-knowledge.finding-format"
	TextDescKeyCheckKnowledgeRelayMessage            = "check-knowledge.relay-message"
	TextDescKeyCheckKnowledgeRelayPrefix             = "check-knowledge.relay-prefix"
	TextDescKeyCheckPersistenceBoxTitle              = "check-persistence.box-title"
	TextDescKeyCheckPersistenceBoxTitleFormat        = "check-persistence.box-title-format"
	TextDescKeyCheckPersistenceCheckpointFormat      = "check-persistence.checkpoint-format"
	TextDescKeyCheckPersistenceFallback              = "check-persistence.fallback"
	TextDescKeyCheckPersistenceInitLogFormat         = "check-persistence.init-log-format"
	TextDescKeyCheckPersistenceModifiedLogFormat     = "check-persistence.modified-log-format"
	TextDescKeyCheckPersistenceRelayFormat           = "check-persistence.relay-format"
	TextDescKeyCheckPersistenceRelayPrefix           = "check-persistence.relay-prefix"
	TextDescKeyCheckPersistenceSilencedLogFormat     = "check-persistence.silenced-log-format"
	TextDescKeyCheckPersistenceSilentLogFormat       = "check-persistence.silent-log-format"
	TextDescKeyCheckPersistenceStateFormat           = "check-persistence.state-format"
	TextDescKeyCheckVersionBoxTitle                  = "check-version.box-title"
	TextDescKeyCheckVersionFallback                  = "check-version.fallback"
	TextDescKeyCheckVersionKeyBoxTitle               = "check-version.key-box-title"
	TextDescKeyCheckVersionKeyFallback               = "check-version.key-fallback"
	TextDescKeyCheckVersionKeyRelayFormat            = "check-version.key-relay-format"
	TextDescKeyCheckVersionKeyRelayPrefix            = "check-version.key-relay-prefix"
	TextDescKeyCheckVersionMismatchRelayFormat       = "check-version.mismatch-relay-format"
	TextDescKeyCheckVersionRelayPrefix               = "check-version.relay-prefix"
	TextDescKeyCheckMapStalenessBoxTitle             = "check-map-staleness.box-title"
	TextDescKeyCheckMapStalenessFallback             = "check-map-staleness.fallback"
	TextDescKeyCheckMapStalenessRelayMessage         = "check-map-staleness.relay-message"
	TextDescKeyCheckMapStalenessRelayPrefix          = "check-map-staleness.relay-prefix"
	TextDescKeyCheckMemoryDriftBoxTitle              = "check-memory-drift.box-title"
	TextDescKeyCheckMemoryDriftContent               = "check-memory-drift.content"
	TextDescKeyCheckMemoryDriftRelayMessage          = "check-memory-drift.relay-message"
	TextDescKeyCheckMemoryDriftRelayPrefix           = "check-memory-drift.relay-prefix"
	TextDescKeyCeremonyBoxBoth                       = "ceremony.box-both"
	TextDescKeyCeremonyBoxRemember                   = "ceremony.box-remember"
	TextDescKeyCeremonyBoxWrapup                     = "ceremony.box-wrapup"
	TextDescKeyCeremonyFallbackBoth                  = "ceremony.fallback-both"
	TextDescKeyCeremonyFallbackRemember              = "ceremony.fallback-remember"
	TextDescKeyCeremonyFallbackWrapup                = "ceremony.fallback-wrapup"
	TextDescKeyCeremonyRelayMessage                  = "ceremony.relay-message"
	TextDescKeyCeremonyRelayPrefix                   = "ceremony.relay-prefix"
	TextDescKeyCheckRemindersBoxTitle                = "check-reminders.box-title"
	TextDescKeyCheckRemindersDismissHint             = "check-reminders.dismiss-hint"
	TextDescKeyCheckRemindersDismissAllHint          = "check-reminders.dismiss-all-hint"
	TextDescKeyCheckRemindersItemFormat              = "check-reminders.item-format"
	TextDescKeyCheckRemindersNudgeFormat             = "check-reminders.nudge-format"
	TextDescKeyCheckRemindersRelayPrefix             = "check-reminders.relay-prefix"

	TextDescKeyCheckResourcesBoxTitle        = "check-resources.box-title"
	TextDescKeyCheckResourcesFallbackLow     = "check-resources.fallback-low"
	TextDescKeyCheckResourcesFallbackPersist = "check-resources.fallback-persist"
	TextDescKeyCheckResourcesFallbackEnd     = "check-resources.fallback-end"
	TextDescKeyCheckResourcesRelayMessage    = "check-resources.relay-message"
	TextDescKeyCheckResourcesRelayPrefix     = "check-resources.relay-prefix"

	TextDescKeyCheckTaskCompletionFallback     = "check-task-completion.fallback"
	TextDescKeyCheckTaskCompletionNudgeMessage = "check-task-completion.nudge-message"

	TextDescKeyVersionDriftRelayMessage = "version-drift.relay-message"

	TextDescKeyChangesFallbackLabel              = "changes.fallback-label"
	TextDescKeyChangesSincePrefix                = "changes.since-prefix"
	TextDescKeyDoctorContextFileFormat           = "doctor.context-file.format"
	TextDescKeyDoctorContextInitializedError     = "doctor.context-initialized.error"
	TextDescKeyDoctorContextInitializedOk        = "doctor.context-initialized.ok"
	TextDescKeyDoctorContextSizeFormat           = "doctor.context-size.format"
	TextDescKeyDoctorContextSizeWarningSuffix    = "doctor.context-size.warning-suffix"
	TextDescKeyDoctorCtxrcValidationError        = "doctor.ctxrc-validation.error"
	TextDescKeyDoctorCtxrcValidationOk           = "doctor.ctxrc-validation.ok"
	TextDescKeyDoctorCtxrcValidationOkNoFile     = "doctor.ctxrc-validation.ok-no-file"
	TextDescKeyDoctorCtxrcValidationWarning      = "doctor.ctxrc-validation.warning"
	TextDescKeyDoctorDriftDetected               = "doctor.drift.detected"
	TextDescKeyDoctorDriftOk                     = "doctor.drift.ok"
	TextDescKeyDoctorDriftViolations             = "doctor.drift.violations"
	TextDescKeyDoctorDriftWarningLoad            = "doctor.drift.warning-load"
	TextDescKeyDoctorDriftWarnings               = "doctor.drift.warnings"
	TextDescKeyDoctorEventLoggingInfo            = "doctor.event-logging.info"
	TextDescKeyDoctorEventLoggingOk              = "doctor.event-logging.ok"
	TextDescKeyDoctorOutputHeader                = "doctor.output.header"
	TextDescKeyDoctorOutputResultLine            = "doctor.output.result-line"
	TextDescKeyDoctorOutputSeparator             = "doctor.output.separator"
	TextDescKeyDoctorOutputSummary               = "doctor.output.summary"
	TextDescKeyDoctorPluginEnabledGlobalOk       = "doctor.plugin-enabled-global.ok"
	TextDescKeyDoctorPluginEnabledLocalOk        = "doctor.plugin-enabled-local.ok"
	TextDescKeyDoctorPluginEnabledWarning        = "doctor.plugin-enabled.warning"
	TextDescKeyDoctorPluginInstalledInfo         = "doctor.plugin-installed.info"
	TextDescKeyDoctorPluginInstalledOk           = "doctor.plugin-installed.ok"
	TextDescKeyDoctorRecentEventsInfo            = "doctor.recent-events.info"
	TextDescKeyDoctorRecentEventsOk              = "doctor.recent-events.ok"
	TextDescKeyDoctorRemindersInfo               = "doctor.reminders.info"
	TextDescKeyDoctorRemindersOk                 = "doctor.reminders.ok"
	TextDescKeyDoctorRequiredFilesError          = "doctor.required-files.error"
	TextDescKeyDoctorRequiredFilesOk             = "doctor.required-files.ok"
	TextDescKeyDoctorResourceDiskFormat          = "doctor.resource-disk.format"
	TextDescKeyDoctorResourceLoadFormat          = "doctor.resource-load.format"
	TextDescKeyDoctorResourceMemoryFormat        = "doctor.resource-memory.format"
	TextDescKeyDoctorResourceSwapFormat          = "doctor.resource-swap.format"
	TextDescKeyDoctorTaskCompletionFormat        = "doctor.task-completion.format"
	TextDescKeyDoctorTaskCompletionWarningSuffix = "doctor.task-completion.warning-suffix"
	TextDescKeyDoctorWebhookInfo                 = "doctor.webhook.info"
	TextDescKeyDoctorWebhookOk                   = "doctor.webhook.ok"
	TextDescKeyHookAider                         = "hook.aider"
	TextDescKeyImportCountConvention             = "import.count-convention"
	TextDescKeyImportCountDecision               = "import.count-decision"
	TextDescKeyImportCountLearning               = "import.count-learning"
	TextDescKeyImportCountTask                   = "import.count-task"
	TextDescKeyHookClaude                        = "hook.claude"
	TextDescKeyHookCopilot                       = "hook.copilot"
	TextDescKeyHookCursor                        = "hook.cursor"
	TextDescKeyHookSupportedTools                = "hook.supported-tools"
	TextDescKeyHookWindsurf                      = "hook.windsurf"
	TextDescKeyTimeAgo                           = "time.ago"
	TextDescKeyTimeDay                           = "time.day"
	TextDescKeyTimeHour                          = "time.hour"
	TextDescKeyTimeJustNow                       = "time.just-now"
	TextDescKeyTimeMinute                        = "time.minute"

	TextDescKeyConfirmProceed           = "confirm.proceed"
	TextDescKeySyncDepsDescription      = "sync.deps.description"
	TextDescKeySyncDepsSuggestion       = "sync.deps.suggestion"
	TextDescKeySyncConfigDescription    = "sync.config.description"
	TextDescKeySyncConfigSuggestion     = "sync.config.suggestion"
	TextDescKeySyncDirDescription       = "sync.dir.description"
	TextDescKeySyncDirSuggestion        = "sync.dir.suggestion"
	TextDescKeyBlockNonPathRelayMessage = "block.non-path-relay-message"
	TextDescKeyBlockConstitutionSuffix  = "block.constitution-suffix"
	TextDescKeyBlockMidSudo             = "block.mid-sudo"
	TextDescKeyBlockMidGitPush          = "block.mid-git-push"
	TextDescKeyBlockCpToBin             = "block.cp-to-bin"
	TextDescKeyBlockInstallToLocalBin   = "block.install-to-local-bin"
	TextDescKeyBlockDotSlash            = "block.dot-slash"
	TextDescKeyBlockGoRun               = "block.go-run"
	TextDescKeyBlockAbsolutePath        = "block.absolute-path"
	TextDescKeyPadKeyCreated            = "pad.key-created"
	TextDescKeyParserGitNotFound        = "parser.git-not-found"
	TextDescKeyParserSessionPrefix      = "parser.session_prefix"
	TextDescKeyParserSessionPrefixAlt   = "parser.session_prefix_alt"
	TextDescKeyPauseConfirmed           = "pause.confirmed"
	TextDescKeyPostCommitFallback       = "post-commit.fallback"
	TextDescKeyPostCommitRelayMessage   = "post-commit.relay-message"

	TextDescKeyPruneDryRunLine    = "prune.dry-run-line"
	TextDescKeyPruneDryRunSummary = "prune.dry-run-summary"
	TextDescKeyPruneErrorLine     = "prune.error-line"
	TextDescKeyPruneSummary       = "prune.summary"

	TextDescKeyMarkJournalChecked     = "mark-journal.checked"
	TextDescKeyMarkJournalMarked      = "mark-journal.marked"
	TextDescKeyMarkWrappedUpConfirmed = "mark-wrapped-up.confirmed"

	TextDescKeyMessageCtxSpecificWarning = "message.ctx-specific-warning"
	TextDescKeyMessageEditHint           = "message.edit-hint"
	TextDescKeyMessageListHeaderCategory = "message.list-header-category"
	TextDescKeyMessageListHeaderHook     = "message.list-header-hook"
	TextDescKeyMessageListHeaderOverride = "message.list-header-override"
	TextDescKeyMessageListHeaderVariant  = "message.list-header-variant"
	TextDescKeyMessageNoOverride         = "message.no-override"
	TextDescKeyMessageOverrideCreated    = "message.override-created"
	TextDescKeyMessageOverrideLabel      = "message.override-label"
	TextDescKeyMessageOverrideRemoved    = "message.override-removed"
	TextDescKeyMessageSourceDefault      = "message.source-default"
	TextDescKeyMessageSourceOverride     = "message.source-override"
	TextDescKeyMessageTemplateVarsLabel  = "message.template-vars-label"
	TextDescKeyMessageTemplateVarsNone   = "message.template-vars-none"

	TextDescKeySpecsNudgeFallback     = "specs-nudge.fallback"
	TextDescKeySpecsNudgeNudgeMessage = "specs-nudge.nudge-message"

	TextDescKeyQaReminderFallback     = "qa-reminder.fallback"
	TextDescKeyQaReminderRelayMessage = "qa-reminder.relay-message"

	TextDescKeyResourcesAlertDisk    = "resources.alert-disk"
	TextDescKeyResourcesAlertLoad    = "resources.alert-load"
	TextDescKeyResourcesAlertMemory  = "resources.alert-memory"
	TextDescKeyResourcesAlertSwap    = "resources.alert-swap"
	TextDescKeyResourcesAlertDanger  = "resources.alert-danger"
	TextDescKeyResourcesAlertWarning = "resources.alert-warning"
	TextDescKeyResourcesAlerts       = "resources.alerts"
	TextDescKeyResourcesAllClear     = "resources.all-clear"
	TextDescKeyResourcesHeader       = "resources.header"
	TextDescKeyResourcesSeparator    = "resources.separator"
	TextDescKeyResourcesStatusDanger = "resources.status-danger"
	TextDescKeyResourcesStatusOk     = "resources.status-ok"
	TextDescKeyResourcesStatusWarn   = "resources.status-warn"
	TextDescKeyResumeConfirmed       = "resume.confirmed"

	TextDescKeyRcParseWarning = "rc.parse_warning"

	TextDescKeySummaryActive     = "summary.active"
	TextDescKeySummaryCompleted  = "summary.completed"
	TextDescKeySummaryDecision   = "summary.decision"
	TextDescKeySummaryDecisions  = "summary.decisions"
	TextDescKeySummaryEmpty      = "summary.empty"
	TextDescKeySummaryInvariants = "summary.invariants"
	TextDescKeySummaryLoaded     = "summary.loaded"
	TextDescKeySummaryTerm       = "summary.term"
	TextDescKeySummaryTerms      = "summary.terms"

	TextDescKeyTaskArchiveContentPreview = "task-archive.content-preview"
	TextDescKeyTaskArchiveDryRunHeader   = "task-archive.dry-run-header"
	TextDescKeyTaskArchiveDryRunSummary  = "task-archive.dry-run-summary"
	TextDescKeyTaskArchiveNoCompleted    = "task-archive.no-completed"
	TextDescKeyTaskArchivePendingRemain  = "task-archive.pending-remain"
	TextDescKeyTaskArchiveSkipIncomplete = "task-archive.skip-incomplete"
	TextDescKeyTaskArchiveSkipping       = "task-archive.skipping"
	TextDescKeyTaskArchiveSuccess        = "task-archive.success"
	TextDescKeyTaskArchiveSuccessWithAge = "task-archive.success-with-age"
	TextDescKeyTaskSnapshotHeaderFormat  = "task-snapshot.header-format"
	TextDescKeyTaskSnapshotCreatedFormat = "task-snapshot.created-format"
	TextDescKeyTaskSnapshotSaved         = "task-snapshot.saved"
	TextDescKeyWatchCloseLogError        = "watch.close-log-error"
	TextDescKeyWatchDryRun               = "watch.dry-run"
	TextDescKeyWatchStopHint             = "watch.stop-hint"
	TextDescKeyWhyAdmonitionFormat       = "why.admonition-format"
	TextDescKeyWhyBanner                 = "why.banner"
	TextDescKeyWhyBlockquotePrefix       = "why.blockquote-prefix"
	TextDescKeyWhyBoldFormat             = "why.bold-format"
	TextDescKeyWhyMenuItemFormat         = "why.menu-item-format"
	TextDescKeyWhyMenuPrompt             = "why.menu-prompt"

	TextDescKeyWatchApplyFailed   = "watch.apply-failed"
	TextDescKeyWatchApplySuccess  = "watch.apply-success"
	TextDescKeyWatchDryRunPreview = "watch.dry-run-preview"
	TextDescKeyWatchWatching      = "watch.watching"

	TextDescKeyDriftCleared = "drift.cleared"

	TextDescKeyMemoryDiffOldFormat = "memory.diff-old-format"
	TextDescKeyMemoryDiffNewFormat = "memory.diff-new-format"
	TextDescKeyMemoryImportSource  = "memory.import-source"
	TextDescKeyMemoryPublishTitle  = "memory.publish-title"
	TextDescKeyMemoryPublishTasks  = "memory.publish-tasks"
	TextDescKeyMemoryPublishDec    = "memory.publish-decisions"
	TextDescKeyMemoryPublishConv   = "memory.publish-conventions"
	TextDescKeyMemoryPublishLrn    = "memory.publish-learnings"
	TextDescKeyMemorySelectContent = "memory.select-content"
	TextDescKeyMemoryWriteMemory   = "memory.write-memory"
	TextDescKeyMemoryImportReview  = "memory.import-review"

	TextDescKeyMCPResConstitution = "mcp.res-constitution"
	TextDescKeyMCPResTasks        = "mcp.res-tasks"
	TextDescKeyMCPResConventions  = "mcp.res-conventions"
	TextDescKeyMCPResArchitecture = "mcp.res-architecture"
	TextDescKeyMCPResDecisions    = "mcp.res-decisions"
	TextDescKeyMCPResLearnings    = "mcp.res-learnings"
	TextDescKeyMCPResGlossary     = "mcp.res-glossary"
	TextDescKeyMCPResPlaybook     = "mcp.res-playbook"
	TextDescKeyMCPResAgent        = "mcp.res-agent"
	TextDescKeyMCPFailedMarshal   = "mcp.failed-marshal"
	TextDescKeyMCPLoadContext     = "mcp.load-context"
	TextDescKeyMCPMethodNotFound  = "mcp.method-not-found"
	TextDescKeyMCPPacketHeader    = "mcp.packet-header"
	TextDescKeyMCPParseError      = "mcp.parse-error"
	TextDescKeyMCPFileNotFound    = "mcp.file-not-found"
	TextDescKeyMCPInvalidParams   = "mcp.invalid-params"
	TextDescKeyMCPUnknownResource = "mcp.unknown-resource"
	TextDescKeyMCPUnknownTool     = "mcp.unknown-tool"

	TextDescKeyMCPToolStatusDesc      = "mcp.tool-status-desc"
	TextDescKeyMCPToolAddDesc         = "mcp.tool-add-desc"
	TextDescKeyMCPToolCompleteDesc    = "mcp.tool-complete-desc"
	TextDescKeyMCPToolDriftDesc       = "mcp.tool-drift-desc"
	TextDescKeyMCPToolPropType        = "mcp.tool-prop-type"
	TextDescKeyMCPToolPropContent     = "mcp.tool-prop-content"
	TextDescKeyMCPToolPropPriority    = "mcp.tool-prop-priority"
	TextDescKeyMCPToolPropContext     = "mcp.tool-prop-context"
	TextDescKeyMCPToolPropRationale   = "mcp.tool-prop-rationale"
	TextDescKeyMCPToolPropConseq      = "mcp.tool-prop-consequences"
	TextDescKeyMCPToolPropLesson      = "mcp.tool-prop-lesson"
	TextDescKeyMCPToolPropApplication = "mcp.tool-prop-application"
	TextDescKeyMCPToolPropQuery       = "mcp.tool-prop-query"
	TextDescKeyMCPTypeContentRequired = "mcp.type-content-required"
	TextDescKeyMCPQueryRequired       = "mcp.query-required"
	TextDescKeyMCPWriteFailed         = "mcp.write-failed"
	TextDescKeyMCPAddedFormat         = "mcp.added-format"
	TextDescKeyMCPCompletedFormat     = "mcp.completed-format"
	TextDescKeyMCPStatusContextFormat = "mcp.status-context-format"
	TextDescKeyMCPStatusFilesFormat   = "mcp.status-files-format"
	TextDescKeyMCPStatusTokensFormat  = "mcp.status-tokens-format"
	TextDescKeyMCPStatusFileFormat    = "mcp.status-file-format"
	TextDescKeyMCPStatusOK            = "mcp.status-ok"
	TextDescKeyMCPStatusEmpty         = "mcp.status-empty"
	TextDescKeyMCPDriftStatusFormat   = "mcp.drift-status-format"
	TextDescKeyMCPDriftViolations     = "mcp.drift-violations"
	TextDescKeyMCPDriftWarnings       = "mcp.drift-warnings"
	TextDescKeyMCPDriftPassed         = "mcp.drift-passed"
	TextDescKeyMCPDriftIssueFormat    = "mcp.drift-issue-format"
	TextDescKeyMCPDriftPassedFormat   = "mcp.drift-passed-format"
	TextDescKeyMCPSectionFormat       = "mcp.section-format"
	TextDescKeyMCPAlsoNoted           = "mcp.also-noted"
	TextDescKeyMCPOmittedFormat       = "mcp.omitted-format"
	TextDescKeyDriftDeadPath          = "drift.dead-path"
	TextDescKeyDriftEntryCount        = "drift.entry-count"
	TextDescKeyDriftMissingFile       = "drift.missing-file"
	TextDescKeyDriftRegenerated       = "drift.regenerated"
	TextDescKeyDriftMissingPackage    = "drift.missing-package"
	TextDescKeyDriftSecret            = "drift.secret"
	TextDescKeyDriftStaleAge          = "drift.stale-age"
	TextDescKeyDriftStaleness         = "drift.staleness"

	TextDescKeyJournalMocSessionLink    = "journal.moc.session-link"
	TextDescKeyJournalMocNavDescription = "journal.moc.nav-description"
	TextDescKeyJournalMocBrowseBy       = "journal.moc.browse-by"
	TextDescKeyJournalMocTopicsDesc     = "journal.moc.topics-description"
	TextDescKeyJournalMocFilesDesc      = "journal.moc.files-description"
	TextDescKeyJournalMocTypesDesc      = "journal.moc.types-description"

	TextDescKeyWriteAddedTo                    = "write.added-to"
	TextDescKeyWriteArchived                   = "write.archived"
	TextDescKeyWriteBackupResult               = "write.backup-result"
	TextDescKeyWriteBackupSMBDest              = "write.backup-smb-dest"
	TextDescKeyWriteBootstrapDir               = "write.bootstrap-dir"
	TextDescKeyWriteBootstrapFiles             = "write.bootstrap-files"
	TextDescKeyWriteBootstrapNextSteps         = "write.bootstrap-next-steps"
	TextDescKeyWriteBootstrapNumbered          = "write.bootstrap-numbered"
	TextDescKeyWriteBootstrapRules             = "write.bootstrap-rules"
	TextDescKeyWriteBootstrapSep               = "write.bootstrap-sep"
	TextDescKeyWriteBootstrapTitle             = "write.bootstrap-title"
	TextDescKeyWriteBootstrapWarning           = "write.bootstrap-warning"
	TextDescKeyWriteCompletedTask              = "write.completed-task"
	TextDescKeyWriteConfigProfileBase          = "write.config-profile-base"
	TextDescKeyWriteConfigProfileDev           = "write.config-profile-dev"
	TextDescKeyWriteConfigProfileNone          = "write.config-profile-none"
	TextDescKeyWriteDepsLookingFor             = "write.deps-looking-for"
	TextDescKeyWriteDepsNoDeps                 = "write.deps-no-deps"
	TextDescKeyWriteDepsNoProject              = "write.deps-no-project"
	TextDescKeyWriteDepsUseType                = "write.deps-use-type"
	TextDescKeyWriteDryRun                     = "write.dry-run"
	TextDescKeyWriteExistsWritingAsAlternative = "write.exists-writing-as-alternative"
	TextDescKeyWriteHookCopilotCreated         = "write.hook-copilot-created"
	TextDescKeyWriteHookCopilotForceHint       = "write.hook-copilot-force-hint"
	TextDescKeyWriteHookCopilotMerged          = "write.hook-copilot-merged"
	TextDescKeyWriteHookCopilotSessionsDir     = "write.hook-copilot-sessions-dir"
	TextDescKeyWriteHookCopilotSkipped         = "write.hook-copilot-skipped"
	TextDescKeyWriteHookCopilotSummary         = "write.hook-copilot-summary"
	TextDescKeyWriteHookUnknownTool            = "write.hook-unknown-tool"
	TextDescKeyWriteImportAdded                = "write.import-added"
	TextDescKeyWriteImportClassified           = "write.import-classified"
	TextDescKeyWriteImportClassifiedSkip       = "write.import-classified-skip"
	TextDescKeyWriteImportDuplicates           = "write.import-duplicates"
	TextDescKeyWriteImportEntry                = "write.import-entry"
	TextDescKeyWriteImportFound                = "write.import-found"
	TextDescKeyWriteImportNoEntries            = "write.import-no-entries"
	TextDescKeyWriteImportScanning             = "write.import-scanning"
	TextDescKeyWriteImportSkipped              = "write.import-skipped"
	TextDescKeyWriteImportSummary              = "write.import-summary"
	TextDescKeyWriteImportSummaryDryRun        = "write.import-summary-dry-run"
	TextDescKeyWriteInitAborted                = "write.init-aborted"
	TextDescKeyWriteInitBackup                 = "write.init-backup"
	TextDescKeyWriteInitCreatedDir             = "write.init-created-dir"
	TextDescKeyWriteInitCreatedWith            = "write.init-created-with"
	TextDescKeyWriteInitCreatingRootFiles      = "write.init-creating-root-files"
	TextDescKeyWriteInitCtxContentExists       = "write.init-ctx-content-exists"
	TextDescKeyWriteInitExistsSkipped          = "write.init-exists-skipped"
	TextDescKeyWriteInitFileCreated            = "write.init-file-created"
	TextDescKeyWriteInitFileExistsNoCtx        = "write.init-file-exists-no-ctx"
	TextDescKeyWriteInitGitignoreReview        = "write.init-gitignore-review"
	TextDescKeyWriteInitGitignoreUpdated       = "write.init-gitignore-updated"
	TextDescKeyWriteInitMakefileAppended       = "write.init-makefile-appended"
	TextDescKeyWriteInitMakefileCreated        = "write.init-makefile-created"
	TextDescKeyWriteInitMakefileIncludes       = "write.init-makefile-includes"
	TextDescKeyWriteInitMerged                 = "write.init-merged"
	TextDescKeyWriteInitNextSteps              = "write.init-next-steps"
	TextDescKeyWriteInitNoChanges              = "write.init-no-changes"
	TextDescKeyWriteInitOverwritePrompt        = "write.init-overwrite-prompt"
	TextDescKeyWriteInitPermsAllow             = "write.init-perms-allow"
	TextDescKeyWriteInitPermsAllowDeny         = "write.init-perms-allow-deny"
	TextDescKeyWriteInitPermsDeduped           = "write.init-perms-deduped"
	TextDescKeyWriteInitPermsDeny              = "write.init-perms-deny"
	TextDescKeyWriteInitPermsMergedDeduped     = "write.init-perms-merged-deduped"
	TextDescKeyWriteInitPluginAlreadyEnabled   = "write.init-plugin-already-enabled"
	TextDescKeyWriteInitPluginEnabled          = "write.init-plugin-enabled"
	TextDescKeyWriteInitPluginInfo             = "write.init-plugin-info"
	TextDescKeyWriteInitPluginNote             = "write.init-plugin-note"
	TextDescKeyWriteInitPluginSkipped          = "write.init-plugin-skipped"
	TextDescKeyWriteInitScratchpadKeyCreated   = "write.init-scratchpad-key-created"
	TextDescKeyWriteInitScratchpadNoKey        = "write.init-scratchpad-no-key"
	TextDescKeyWriteInitScratchpadPlaintext    = "write.init-scratchpad-plaintext"
	TextDescKeyWriteInitSettingUpPermissions   = "write.init-setting-up-permissions"
	TextDescKeyWriteInitSkippedDir             = "write.init-skipped-dir"
	TextDescKeyWriteInitSkippedPlain           = "write.init-skipped-plain"
	TextDescKeyWriteInitUpdatedCtxSection      = "write.init-updated-ctx-section"
	TextDescKeyWriteInitUpdatedPlanSection     = "write.init-updated-plan-section"
	TextDescKeyWriteInitUpdatedPromptSection   = "write.init-updated-prompt-section"
	TextDescKeyWriteInitWarnNonFatal           = "write.init-warn-non-fatal"
	TextDescKeyWriteInitialized                = "write.initialized"
	TextDescKeyWriteJournalOrphanRemoved       = "write.journal-orphan-removed"
	TextDescKeyWriteJournalSiteAlt             = "write.journal-site-alt"
	TextDescKeyWriteJournalSiteBuilding        = "write.journal-site-building"
	TextDescKeyWriteJournalSiteGenerated       = "write.journal-site-generated"
	TextDescKeyWriteJournalSiteNextSteps       = "write.journal-site-next-steps"
	TextDescKeyWriteJournalSiteStarting        = "write.journal-site-starting"
	TextDescKeyWriteJournalSyncLocked          = "write.journal-sync-locked"
	TextDescKeyWriteJournalSyncLockedCount     = "write.journal-sync-locked-count"
	TextDescKeyWriteJournalSyncMatch           = "write.journal-sync-match"
	TextDescKeyWriteJournalSyncNone            = "write.journal-sync-none"
	TextDescKeyWriteJournalSyncUnlocked        = "write.journal-sync-unlocked"
	TextDescKeyWriteJournalSyncUnlockedCount   = "write.journal-sync-unlocked-count"
	TextDescKeyWriteLines                      = "write.lines"
	TextDescKeyWriteLinesPrevious              = "write.lines-previous"
	TextDescKeyWriteLockUnlockEntry            = "write.lock-unlock-entry"
	TextDescKeyWriteLockUnlockNoChanges        = "write.lock-unlock-no-changes"
	TextDescKeyWriteLockUnlockSummary          = "write.lock-unlock-summary"
	TextDescKeyWriteLoopCompletion             = "write.loop-completion"
	TextDescKeyWriteLoopGenerated              = "write.loop-generated"
	TextDescKeyWriteLoopMaxIterations          = "write.loop-max-iterations"
	TextDescKeyWriteLoopPrompt                 = "write.loop-prompt"
	TextDescKeyWriteLoopRunCmd                 = "write.loop-run-cmd"
	TextDescKeyWriteLoopTool                   = "write.loop-tool"
	TextDescKeyWriteLoopUnlimited              = "write.loop-unlimited"
	TextDescKeyWriteMemoryArchives             = "write.memory-archives"
	TextDescKeyWriteMemoryBridgeHeader         = "write.memory-bridge-header"
	TextDescKeyWriteMemoryDriftDetected        = "write.memory-drift-detected"
	TextDescKeyWriteMemoryDriftNone            = "write.memory-drift-none"
	TextDescKeyWriteMemoryLastSync             = "write.memory-last-sync"
	TextDescKeyWriteMemoryLastSyncNever        = "write.memory-last-sync-never"
	TextDescKeyWriteMemoryMirror               = "write.memory-mirror"
	TextDescKeyWriteMemoryMirrorLines          = "write.memory-mirror-lines"
	TextDescKeyWriteMemoryMirrorNotSynced      = "write.memory-mirror-not-synced"
	TextDescKeyWriteMemoryNoChanges            = "write.memory-no-changes"
	TextDescKeyWriteMemorySource               = "write.memory-source"
	TextDescKeyWriteMemorySourceLines          = "write.memory-source-lines"
	TextDescKeyWriteMemorySourceLinesDrift     = "write.memory-source-lines-drift"
	TextDescKeyWriteMemorySourceNotActive      = "write.memory-source-not-active"
	TextDescKeyWriteMirror                     = "write.mirror"
	TextDescKeyWriteMovingTask                 = "write.moving-task"
	TextDescKeyWriteNewContent                 = "write.new-content"
	TextDescKeyWriteObsidianGenerated          = "write.obsidian-generated"
	TextDescKeyWriteObsidianNextSteps          = "write.obsidian-next-steps"
	TextDescKeyWritePadBlobWritten             = "write.pad-blob-written"
	TextDescKeyWritePadEmpty                   = "write.pad-empty"
	TextDescKeyWritePadEntryAdded              = "write.pad-entry-added"
	TextDescKeyWritePadEntryMoved              = "write.pad-entry-moved"
	TextDescKeyWritePadEntryRemoved            = "write.pad-entry-removed"
	TextDescKeyWritePadEntryUpdated            = "write.pad-entry-updated"
	TextDescKeyWritePadExportDone              = "write.pad-export-done"
	TextDescKeyWritePadExportNone              = "write.pad-export-none"
	TextDescKeyWritePadExportPlan              = "write.pad-export-plan"
	TextDescKeyWritePadExportSummary           = "write.pad-export-summary"
	TextDescKeyWritePadExportVerbDone          = "write.pad-export-verb-done"
	TextDescKeyWritePadExportVerbDryRun        = "write.pad-export-verb-dry-run"
	TextDescKeyWritePadExportWriteFailed       = "write.pad-export-write-failed"
	TextDescKeyWritePadImportBlobAdded         = "write.pad-import-blob-added"
	TextDescKeyWritePadImportBlobNone          = "write.pad-import-blob-none"
	TextDescKeyWritePadImportBlobSkipped       = "write.pad-import-blob-skipped"
	TextDescKeyWritePadImportBlobSummary       = "write.pad-import-blob-summary"
	TextDescKeyWritePadImportBlobTooLarge      = "write.pad-import-blob-too-large"
	TextDescKeyWritePadImportCloseWarning      = "write.pad-import-close-warning"
	TextDescKeyWritePadImportDone              = "write.pad-import-done"
	TextDescKeyWritePadImportNone              = "write.pad-import-none"
	TextDescKeyWritePadKeyCreated              = "write.pad-key-created"
	TextDescKeyWritePadMergeAdded              = "write.pad-merge-added"
	TextDescKeyWritePadMergeBinaryWarning      = "write.pad-merge-binary-warning"
	TextDescKeyWritePadMergeBlobConflict       = "write.pad-merge-blob-conflict"
	TextDescKeyWritePadMergeDone               = "write.pad-merge-done"
	TextDescKeyWritePadMergeDryRun             = "write.pad-merge-dry-run"
	TextDescKeyWritePadMergeDupe               = "write.pad-merge-dupe"
	TextDescKeyWritePadMergeNone               = "write.pad-merge-none"
	TextDescKeyWritePadMergeNoneNew            = "write.pad-merge-none-new"
	TextDescKeyWritePadResolveEntry            = "write.pad-resolve-entry"
	TextDescKeyWritePadResolveHeader           = "write.pad-resolve-header"
	TextDescKeyWritePathExists                 = "write.path-exists"
	TextDescKeyWritePaused                     = "write.paused"
	TextDescKeyWritePrefixError                = "write.prefix-error"
	TextDescKeyWritePromptCreated              = "write.prompt-created"
	TextDescKeyWritePromptItem                 = "write.prompt-item"
	TextDescKeyWritePromptNone                 = "write.prompt-none"
	TextDescKeyWritePromptRemoved              = "write.prompt-removed"
	TextDescKeyWritePublishBlock               = "write.publish-block"
	TextDescKeyWritePublishBudget              = "write.publish-budget"
	TextDescKeyWritePublishConventions         = "write.publish-conventions"
	TextDescKeyWritePublishDecisions           = "write.publish-decisions"
	TextDescKeyWritePublishDone                = "write.publish-done"
	TextDescKeyWritePublishDryRun              = "write.publish-dry-run"
	TextDescKeyWritePublishHeader              = "write.publish-header"
	TextDescKeyWritePublishLearnings           = "write.publish-learnings"
	TextDescKeyWritePublishSourceFiles         = "write.publish-source-files"
	TextDescKeyWritePublishTasks               = "write.publish-tasks"
	TextDescKeyWritePublishTotal               = "write.publish-total"
	TextDescKeyWriteReminderAdded              = "write.reminder-added"
	TextDescKeyWriteReminderAfterSuffix        = "write.reminder-after-suffix"
	TextDescKeyWriteReminderDismissed          = "write.reminder-dismissed"
	TextDescKeyWriteReminderDismissedAll       = "write.reminder-dismissed-all"
	TextDescKeyWriteReminderItem               = "write.reminder-item"
	TextDescKeyWriteReminderNone               = "write.reminder-none"
	TextDescKeyWriteReminderNotDue             = "write.reminder-not-due"
	TextDescKeyWriteRestoreAdded               = "write.restore-added"
	TextDescKeyWriteRestoreDenyDroppedHeader   = "write.restore-deny-dropped-header"
	TextDescKeyWriteRestoreDenyRestoredHeader  = "write.restore-deny-restored-header"
	TextDescKeyWriteRestoreDone                = "write.restore-done"
	TextDescKeyWriteRestoreDroppedHeader       = "write.restore-dropped-header"
	TextDescKeyWriteRestoreMatch               = "write.restore-match"
	TextDescKeyWriteRestoreNoLocal             = "write.restore-no-local"
	TextDescKeyWriteRestorePermMatch           = "write.restore-perm-match"
	TextDescKeyWriteRestoreRemoved             = "write.restore-removed"
	TextDescKeyWriteRestoreRestoredHeader      = "write.restore-restored-header"
	TextDescKeyWriteResumed                    = "write.resumed"
	TextDescKeyWriteSetupDone                  = "write.setup-done"
	TextDescKeyWriteSetupPrompt                = "write.setup-prompt"
	TextDescKeyWriteSkillLine                  = "write.skill-line"
	TextDescKeyWriteSkillsHeader               = "write.skills-header"
	TextDescKeyWriteSnapshotSaved              = "write.snapshot-saved"
	TextDescKeyWriteSnapshotUpdated            = "write.snapshot-updated"
	TextDescKeyWriteSource                     = "write.source"
	TextDescKeyWriteStatusActivityHeader       = "write.status-activity-header"
	TextDescKeyWriteStatusActivityItem         = "write.status-activity-item"
	TextDescKeyWriteStatusDir                  = "write.status-dir"
	TextDescKeyWriteStatusDrift                = "write.status-drift"
	TextDescKeyWriteStatusFileCompact          = "write.status-file-compact"
	TextDescKeyWriteStatusFileVerbose          = "write.status-file-verbose"
	TextDescKeyWriteStatusFiles                = "write.status-files"
	TextDescKeyWriteStatusFilesHeader          = "write.status-files-header"
	TextDescKeyWriteStatusNoDrift              = "write.status-no-drift"
	TextDescKeyWriteStatusPreviewLine          = "write.status-preview-line"
	TextDescKeyWriteStatusSeparator            = "write.status-separator"
	TextDescKeyWriteStatusTitle                = "write.status-title"
	TextDescKeyWriteStatusTokens               = "write.status-tokens"
	TextDescKeyWriteSynced                     = "write.synced"
	TextDescKeyWriteSyncAction                 = "write.sync-action"
	TextDescKeyWriteSyncDryRun                 = "write.sync-dry-run"
	TextDescKeyWriteSyncDryRunSummary          = "write.sync-dry-run-summary"
	TextDescKeyWriteSyncHeader                 = "write.sync-header"
	TextDescKeyWriteSyncInSync                 = "write.sync-in-sync"
	TextDescKeyWriteSyncSeparator              = "write.sync-separator"
	TextDescKeyWriteSyncSuggestion             = "write.sync-suggestion"
	TextDescKeyWriteSyncSummary                = "write.sync-summary"
	TextDescKeyWriteTestFiltered               = "write.test-filtered"
	TextDescKeyWriteTestNoWebhook              = "write.test-no-webhook"
	TextDescKeyWriteTestResult                 = "write.test-result"
	TextDescKeyWriteTestWorking                = "write.test-working"
	TextDescKeyWriteTimeDayAgo                 = "write.time-day-ago"
	TextDescKeyWriteTimeDaysAgo                = "write.time-days-ago"
	TextDescKeyWriteTimeHourAgo                = "write.time-hour-ago"
	TextDescKeyWriteTimeHoursAgo               = "write.time-hours-ago"
	TextDescKeyWriteTimeJustNow                = "write.time-just-now"
	TextDescKeyWriteTimeMinuteAgo              = "write.time-minute-ago"
	TextDescKeyWriteTimeMinutesAgo             = "write.time-minutes-ago"
	TextDescKeyWriteUnpublishDone              = "write.unpublish-done"
	TextDescKeyWriteUnpublishNotFound          = "write.unpublish-not-found"
)

// Template reads a template file by name from the embedded filesystem.
//
// Parameters:
//   - name: Template filename (e.g., "TASKS.md")
//
// Returns:
//   - []byte: Template content
//   - error: Non-nil if the file is not found or read fails
func Template(name string) ([]byte, error) {
	return FS.ReadFile("context/" + name)
}

// List returns all available template file names.
//
// Returns:
//   - []string: List of template filenames in the root templates directory
//   - error: Non-nil if directory read fails
func List() ([]string, error) {
	entries, err := FS.ReadDir("context")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// ListEntry returns available entry template file names.
//
// Returns:
//   - []string: List of template filenames in entry-templates/
//   - error: Non-nil if directory read fails
func ListEntry() ([]string, error) {
	entries, err := FS.ReadDir("entry-templates")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// Entry reads an entry template by name.
//
// Parameters:
//   - name: Template filename (e.g., "decision.md")
//
// Returns:
//   - []byte: Template content from entry-templates/
//   - error: Non-nil if the file is not found or read fails
func Entry(name string) ([]byte, error) {
	return FS.ReadFile("entry-templates/" + name)
}

// ListPromptTemplates returns available prompt template file names.
//
// Returns:
//   - []string: List of template filenames in prompt-templates/
//   - error: Non-nil if directory read fails
func ListPromptTemplates() ([]string, error) {
	entries, err := FS.ReadDir("prompt-templates")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// PromptTemplate reads a prompt template by name.
//
// Parameters:
//   - name: Template filename (e.g., "code-review.md")
//
// Returns:
//   - []byte: Template content from prompt-templates/
//   - error: Non-nil if the file is not found or read fails
func PromptTemplate(name string) ([]byte, error) {
	return FS.ReadFile("prompt-templates/" + name)
}

// ListSkills returns available skill directory names.
//
// Each skill is a directory containing a SKILL.md file following the
// Agent Skills specification (https://agentskills.io/specification).
//
// Returns:
//   - []string: List of skill directory names in claude/skills/
//   - error: Non-nil if directory read fails
func ListSkills() ([]string, error) {
	entries, err := FS.ReadDir("claude/skills")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// SkillContent reads a skill's SKILL.md file by skill name.
//
// Parameters:
//   - name: Skill directory name (e.g., "ctx-status")
//
// Returns:
//   - []byte: SKILL.md content from claude/skills/<name>/
//   - error: Non-nil if the file not found or read fails
func SkillContent(name string) ([]byte, error) {
	return FS.ReadFile("claude/skills/" + name + "/SKILL.md")
}

// SkillReference reads a reference file from a skill's references/ directory.
//
// Parameters:
//   - skill: Skill directory name (e.g., "ctx-skill-audit")
//   - filename: Reference filename (e.g., "anthropic-best-practices.md")
//
// Returns:
//   - []byte: Reference file content
//   - error: Non-nil if the file is not found or read fails
func SkillReference(skill, filename string) ([]byte, error) {
	return FS.ReadFile("claude/skills/" + skill + "/references/" + filename)
}

// ListSkillReferences returns available reference filenames for a skill.
//
// Parameters:
//   - skill: Skill directory name (e.g., "ctx-skill-audit")
//
// Returns:
//   - []string: List of reference filenames
//   - error: Non-nil if the references directory is not found or read fails
func ListSkillReferences(skill string) ([]string, error) {
	entries, readErr := FS.ReadDir("claude/skills/" + skill + "/references")
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// MakefileCtx reads the ctx-owned Makefile include template.
//
// Returns:
//   - []byte: Makefile.ctx content
//   - error: Non-nil if the file is not found or read fails
func MakefileCtx() ([]byte, error) {
	return FS.ReadFile("project/Makefile.ctx")
}

// ProjectFile reads a project-root file by name from the embedded filesystem.
//
// These files are deployed to the project root (not .context/) by dedicated
// handlers during initialization.
//
// Parameters:
//   - name: Filename (e.g., "IMPLEMENTATION_PLAN.md")
//
// Returns:
//   - []byte: File content
//   - error: Non-nil if the file is not found or read fails
func ProjectFile(name string) ([]byte, error) {
	return FS.ReadFile("project/" + name)
}

// ProjectReadme reads a project directory README template by directory name.
//
// Templates are stored as project/<dir>-README.md in the embedded filesystem.
//
// Parameters:
//   - dir: Directory name (e.g., "specs", "ideas")
//
// Returns:
//   - []byte: README.md content for the directory
//   - error: Non-nil if the file is not found or read fails
func ProjectReadme(dir string) ([]byte, error) {
	return FS.ReadFile("project/" + dir + "-README.md")
}

// ClaudeMd reads the CLAUDE.md template from the embedded filesystem.
//
// CLAUDE.md is deployed to the project root by a dedicated handler
// during initialization, separate from the .context/ templates.
//
// Returns:
//   - []byte: CLAUDE.md content
//   - error: Non-nil if the file is not found or read fails
func ClaudeMd() ([]byte, error) {
	return FS.ReadFile("claude/CLAUDE.md")
}

// RalphTemplate reads a Ralph-mode template file by name.
//
// Ralph mode templates are designed for autonomous loop operation,
// with instructions for one-task-per-iteration, completion signals,
// and no clarifying questions.
//
// Parameters:
//   - name: Template filename (e.g., "PROMPT.md")
//
// Returns:
//   - []byte: Template content from ralph/
//   - error: Non-nil if the file is not found or read fails
func RalphTemplate(name string) ([]byte, error) {
	return FS.ReadFile("ralph/" + name)
}

// HookMessage reads a hook message template by hook name and filename.
//
// Parameters:
//   - hook: Hook directory name (e.g., "qa-reminder")
//   - filename: Template filename (e.g., "gate.txt")
//
// Returns:
//   - []byte: Template content from hooks/messages/<hook>/
//   - error: Non-nil if the file is not found or read fails
func HookMessage(hook, filename string) ([]byte, error) {
	return FS.ReadFile("hooks/messages/" + hook + "/" + filename)
}

// HookMessageRegistry reads the embedded registry.yaml that describes
// all hook message templates.
//
// Returns:
//   - []byte: Raw YAML content
//   - error: Non-nil if the file is not found or read fails
func HookMessageRegistry() ([]byte, error) {
	return FS.ReadFile("hooks/messages/registry.yaml")
}

// CopilotInstructions reads the embedded Copilot instructions template.
//
// Returns:
//   - []byte: Template content from hooks/copilot-instructions.md
//   - error: Non-nil if the file is not found or read fails
func CopilotInstructions() ([]byte, error) {
	return FS.ReadFile("hooks/copilot-instructions.md")
}

// JournalExtraCSS reads the embedded extra.css for journal site generation.
//
// Returns:
//   - []byte: CSS content
//   - error: Non-nil if the file is not found or read fails
func JournalExtraCSS() ([]byte, error) {
	return FS.ReadFile("journal/extra.css")
}

// ListHookMessages returns available hook message directory names.
//
// Each hook is a directory under hooks/messages/ containing one or
// more variant .txt template files.
//
// Returns:
//   - []string: List of hook directory names
//   - error: Non-nil if directory read fails
func ListHookMessages() ([]string, error) {
	entries, readErr := FS.ReadDir("hooks/messages")
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// ListHookVariants returns available variant filenames for a hook.
//
// Parameters:
//   - hook: Hook directory name (e.g., "qa-reminder")
//
// Returns:
//   - []string: List of variant filenames (e.g., "gate.txt")
//   - error: Non-nil if the hook directory is not found or read fails
func ListHookVariants(hook string) ([]string, error) {
	entries, readErr := FS.ReadDir("hooks/messages/" + hook)
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// WhyDoc reads a "why" document by name from the embedded filesystem.
//
// Parameters:
//   - name: Document name (e.g., "manifesto", "about", "design-invariants")
//
// Returns:
//   - []byte: Document content from why/
//   - error: Non-nil if the file is not found or read fails
func WhyDoc(name string) ([]byte, error) {
	return FS.ReadFile("why/" + name + file.ExtMarkdown)
}

// ListWhyDocs returns available "why" document names (without extension).
//
// Returns:
//   - []string: List of document names in why/
//   - error: Non-nil if directory read fails
func ListWhyDocs() ([]string, error) {
	entries, readErr := FS.ReadDir("why")
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if len(name) > 3 && name[len(name)-3:] == file.ExtMarkdown {
				names = append(names, name[:len(name)-3])
			}
		}
	}
	return names, nil
}

// Schema reads the embedded JSON Schema for .ctxrc.
//
// Returns:
//   - []byte: JSON Schema content
//   - error: Non-nil if the file is not found or read fails
func Schema() ([]byte, error) {
	return FS.ReadFile("schema/ctxrc.schema.json")
}

var (
	commandsOnce sync.Once
	commandsMap  map[string]commandEntry
	flagsOnce    sync.Once
	flagsMap     map[string]commandEntry
	textOnce     sync.Once
	textMap      map[string]commandEntry
	examplesOnce sync.Once
	examplesMap  map[string]commandEntry
)

type commandEntry struct {
	Short string `yaml:"short"`
	Long  string `yaml:"long"`
}

// loadYAML parses an embedded YAML file into a commandEntry map.
func loadYAML(path string) map[string]commandEntry {
	data, readErr := FS.ReadFile(path)
	if readErr != nil {
		return make(map[string]commandEntry)
	}
	m := make(map[string]commandEntry)
	if parseErr := yaml.Unmarshal(data, &m); parseErr != nil {
		return make(map[string]commandEntry)
	}
	return m
}

func loadCommands() {
	commandsOnce.Do(func() { commandsMap = loadYAML("commands/commands.yaml") })
}

func loadFlags() {
	flagsOnce.Do(func() { flagsMap = loadYAML("commands/flags.yaml") })
}

func loadText() {
	textOnce.Do(func() { textMap = loadYAML("commands/text.yaml") })
}

func loadExamples() {
	examplesOnce.Do(func() { examplesMap = loadYAML("commands/examples.yaml") })
}

// CommandDesc returns the Short and Long descriptions for a command.
//
// Keys use dot notation: "pad", "pad.show", "system.bootstrap".
// Returns empty strings if the key is not found.
//
// Parameters:
//   - key: Command key in dot notation
//
// Returns:
//   - short: One-line description
//   - long: Multi-line help text (may be empty)
func CommandDesc(key string) (short, long string) {
	loadCommands()
	entry, ok := commandsMap[key]
	if !ok {
		return "", ""
	}
	return entry.Short, entry.Long
}

// FlagDesc returns the description for a flag.
//
// Keys use dot notation: "add.file", "context-dir".
// Returns an empty string if the key is not found.
//
// Parameters:
//   - name: Flag key in dot notation
//
// Returns:
//   - string: Flag description
func FlagDesc(name string) string {
	loadFlags()
	entry, ok := flagsMap[name]
	if !ok {
		return ""
	}
	return entry.Short
}

// ExampleDesc returns example usage text for a given key.
//
// Keys match entry types: "decision", "learning", "task", "convention".
// Returns an empty string if the key is not found.
//
// Parameters:
//   - name: Entry type key
//
// Returns:
//   - string: Example text
func ExampleDesc(name string) string {
	loadExamples()
	entry, ok := examplesMap[name]
	if !ok {
		return ""
	}
	return entry.Short
}

// TextDesc returns a user-facing text string by key.
//
// Keys use dot notation: "agent.instruction", "backup.run-hint".
// Returns an empty string if the key is not found.
//
// Parameters:
//   - name: Text key in dot notation
//
// Returns:
//   - string: Text content
func TextDesc(name string) string {
	loadText()
	entry, ok := textMap[name]
	if !ok {
		return ""
	}
	return entry.Short
}

var (
	stopWordsOnce sync.Once
	stopWordsMap  map[string]bool
)

// StopWords returns the default set of stop words for keyword extraction.
//
// Loaded from the embedded text.yaml asset under "stopwords".
// The result is cached after the first call.
//
// Returns:
//   - map[string]bool: Set of lowercase stop words
func StopWords() map[string]bool {
	stopWordsOnce.Do(func() {
		raw := TextDesc("stopwords")
		words := strings.Fields(raw)
		stopWordsMap = make(map[string]bool, len(words))
		for _, w := range words {
			stopWordsMap[w] = true
		}
	})
	return stopWordsMap
}

var (
	allowOnce  sync.Once
	allowPerms []string

	denyOnce  sync.Once
	denyPerms []string
)

// parsePermissions splits a text file into permission entries.
//
// Lines are trimmed; empty lines and lines starting with '#' are skipped.
func parsePermissions(data []byte) []string {
	var result []string
	for _, line := range strings.Split(string(data), token.NewlineLF) {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		result = append(result, line)
	}
	return result
}

// DefaultAllowPermissions returns the default allow permissions for ctx
// commands and skills, parsed from the embedded permissions/allow.txt.
func DefaultAllowPermissions() []string {
	allowOnce.Do(func() {
		data, readErr := FS.ReadFile("permissions/allow.txt")
		if readErr != nil {
			return
		}
		allowPerms = parsePermissions(data)
	})
	return allowPerms
}

// DefaultDenyPermissions returns the default deny permissions that block
// dangerous operations, parsed from the embedded permissions/deny.txt.
func DefaultDenyPermissions() []string {
	denyOnce.Do(func() {
		data, readErr := FS.ReadFile("permissions/deny.txt")
		if readErr != nil {
			return
		}
		denyPerms = parsePermissions(data)
	})
	return denyPerms
}

// PluginVersion returns the version string from the embedded plugin.json.
func PluginVersion() (string, error) {
	data, err := FS.ReadFile("claude/.claude-plugin/plugin.json")
	if err != nil {
		return "", err
	}
	var manifest struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", err
	}
	return manifest.Version, nil
}
