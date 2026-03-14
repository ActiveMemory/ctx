//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/time"
)

// TplBootstrapTitle is the heading for bootstrap output.
var TplBootstrapTitle = assets.TextDesc(assets.TextDescKeyWriteBootstrapTitle)

// TplBootstrapSep is the visual separator under the bootstrap heading.
var TplBootstrapSep = assets.TextDesc(assets.TextDescKeyWriteBootstrapSep)

// TplBootstrapDir is a format template for the context directory.
// Arguments: context directory path.
var TplBootstrapDir = assets.TextDesc(assets.TextDescKeyWriteBootstrapDir)

// TplBootstrapFiles is the heading for the file list section.
var TplBootstrapFiles = assets.TextDesc(assets.TextDescKeyWriteBootstrapFiles)

// TplBootstrapRules is the heading for the rules section.
var TplBootstrapRules = assets.TextDesc(assets.TextDescKeyWriteBootstrapRules)

// TplBootstrapNextSteps is the heading for the next steps section.
var TplBootstrapNextSteps = assets.TextDesc(assets.TextDescKeyWriteBootstrapNextSteps)

// TplBootstrapNumbered is a format template for a numbered list item.
// Arguments: index, text.
var TplBootstrapNumbered = assets.TextDesc(assets.TextDescKeyWriteBootstrapNumbered)

// TplBootstrapWarning is a format template for a warning line.
// Arguments: warning text.
var TplBootstrapWarning = assets.TextDesc(assets.TextDescKeyWriteBootstrapWarning)

// PrefixError is prepended to all error messages written to stderr.
var PrefixError = assets.TextDesc(assets.TextDescKeyWritePrefixError)

// TplPathExists is a format template for reporting that a destination path
// already exists. Arguments: original path, resolved destination path.
var TplPathExists = assets.TextDesc(assets.TextDescKeyWritePathExists)

// TplExistsWritingAsAlternative is a format template for reporting that a
// file exists and content was written to an alternative filename instead.
// Arguments: original path, alternative path.
var TplExistsWritingAsAlternative = assets.TextDesc(assets.TextDescKeyWriteExistsWritingAsAlternative)

// TplDryRun is printed when a command runs in dry-run mode.
var TplDryRun = assets.TextDesc(assets.TextDescKeyWriteDryRun)

// TplSource is a format template for reporting a source path.
// Arguments: path.
var TplSource = assets.TextDesc(assets.TextDescKeyWriteSource)

// TplMirror is a format template for reporting a mirror path.
// Arguments: relative mirror path.
var TplMirror = assets.TextDesc(assets.TextDescKeyWriteMirror)

// TplStatusDrift is printed when drift is detected.
var TplStatusDrift = assets.TextDesc(assets.TextDescKeyWriteStatusDrift)

// TplStatusNoDrift is printed when no drift is detected.
var TplStatusNoDrift = assets.TextDesc(assets.TextDescKeyWriteStatusNoDrift)

// TplArchived is a format template for reporting an archived file.
// Arguments: archive filename.
var TplArchived = assets.TextDesc(assets.TextDescKeyWriteArchived)

// TplSynced is a format template for reporting a successful sync.
// Arguments: source label, destination relative path.
var TplSynced = assets.TextDesc(assets.TextDescKeyWriteSynced)

// TplLines is a format template for reporting line counts.
// Arguments: line count.
var TplLines = assets.TextDesc(assets.TextDescKeyWriteLines)

// TplLinesPrevious is a format template appended to line counts when a
// previous count is available. Arguments: previous line count.
var TplLinesPrevious = assets.TextDesc(assets.TextDescKeyWriteLinesPrevious)

// TplNewContent is a format template for reporting new content since last sync.
// Arguments: line count.
var TplNewContent = assets.TextDesc(assets.TextDescKeyWriteNewContent)

// TplAddedTo is a format template for confirming an entry was added.
// Arguments: filename.
var TplAddedTo = assets.TextDesc(assets.TextDescKeyWriteAddedTo)

// TplMovingTask is a format template for a task being moved to completed.
// Arguments: truncated task text.
var TplMovingTask = assets.TextDesc(assets.TextDescKeyWriteMovingTask)

// TplCompletedTask is a format template for a task marked complete.
// Arguments: task text.
var TplCompletedTask = assets.TextDesc(assets.TextDescKeyWriteCompletedTask)

// TplConfigProfileDev is the status output for dev profile.
var TplConfigProfileDev = assets.TextDesc(assets.TextDescKeyWriteConfigProfileDev)

// TplConfigProfileBase is the status output for base profile.
var TplConfigProfileBase = assets.TextDesc(assets.TextDescKeyWriteConfigProfileBase)

// TplConfigProfileNone is the status output when no profile exists.
// Arguments: ctxrc filename.
var TplConfigProfileNone = assets.TextDesc(assets.TextDescKeyWriteConfigProfileNone)

// TplDepsNoProject is printed when no supported project is detected.
var TplDepsNoProject = assets.TextDesc(assets.TextDescKeyWriteDepsNoProject)

// TplDepsLookingFor is printed with the list of files checked.
var TplDepsLookingFor = assets.TextDesc(assets.TextDescKeyWriteDepsLookingFor)

// TplDepsUseType hints at the --type flag.
// Arguments: comma-separated list of builder names.
var TplDepsUseType = assets.TextDesc(assets.TextDescKeyWriteDepsUseType)

// TplDepsNoDeps is printed when no dependencies are found.
var TplDepsNoDeps = assets.TextDesc(assets.TextDescKeyWriteDepsNoDeps)

// TplSkillsHeader is the heading for the skills list.
var TplSkillsHeader = assets.TextDesc(assets.TextDescKeyWriteSkillsHeader)

// TplSkillLine formats a single skill entry.
// Arguments: name, description.
var TplSkillLine = assets.TextDesc(assets.TextDescKeyWriteSkillLine)

// TplHookCopilotSkipped reports that copilot instructions were skipped.
// Arguments: target file path.
var TplHookCopilotSkipped = assets.TextDesc(assets.TextDescKeyWriteHookCopilotSkipped)

// TplHookCopilotForceHint tells the user about the --force flag.
var TplHookCopilotForceHint = assets.TextDesc(assets.TextDescKeyWriteHookCopilotForceHint)

// TplHookCopilotMerged reports that copilot instructions were merged.
// Arguments: target file path.
var TplHookCopilotMerged = assets.TextDesc(assets.TextDescKeyWriteHookCopilotMerged)

// TplHookCopilotCreated reports that copilot instructions were created.
// Arguments: target file path.
var TplHookCopilotCreated = assets.TextDesc(assets.TextDescKeyWriteHookCopilotCreated)

// TplHookCopilotSessionsDir reports that the sessions directory was created.
// Arguments: sessions directory path.
var TplHookCopilotSessionsDir = assets.TextDesc(assets.TextDescKeyWriteHookCopilotSessionsDir)

// TplHookCopilotSummary is the post-write summary for copilot.
var TplHookCopilotSummary = assets.TextDesc(assets.TextDescKeyWriteHookCopilotSummary)

// TplHookUnknownTool reports an unrecognized tool name.
// Arguments: tool name.
var TplHookUnknownTool = assets.TextDesc(assets.TextDescKeyWriteHookUnknownTool)

// TplInitOverwritePrompt prompts the user before overwriting .context/.
// Arguments: context directory path.
var TplInitOverwritePrompt = assets.TextDesc(assets.TextDescKeyWriteInitOverwritePrompt)

// TplInitAborted is printed when the user declines overwriting.
var TplInitAborted = assets.TextDesc(assets.TextDescKeyWriteInitAborted)

// TplInitExistsSkipped reports a file that was skipped because it exists.
// Arguments: filename.
var TplInitExistsSkipped = assets.TextDesc(assets.TextDescKeyWriteInitExistsSkipped)

// TplInitFileCreated reports a file that was successfully created.
// Arguments: filename.
var TplInitFileCreated = assets.TextDesc(assets.TextDescKeyWriteInitFileCreated)

// TplInitialized reports successful context initialization.
// Arguments: context directory path.
var TplInitialized = assets.TextDesc(assets.TextDescKeyWriteInitialized)

// TplInitWarnNonFatal reports a non-fatal warning during init.
// Arguments: label, error.
var TplInitWarnNonFatal = assets.TextDesc(assets.TextDescKeyWriteInitWarnNonFatal)

// TplInitScratchpadPlaintext reports a plaintext scratchpad was created.
// Arguments: path.
var TplInitScratchpadPlaintext = assets.TextDesc(assets.TextDescKeyWriteInitScratchpadPlaintext)

// TplInitScratchpadNoKey warns about a missing key for an encrypted scratchpad.
// Arguments: key path.
var TplInitScratchpadNoKey = assets.TextDesc(assets.TextDescKeyWriteInitScratchpadNoKey)

// TplInitScratchpadKeyCreated reports a scratchpad key was generated.
// Arguments: key path.
var TplInitScratchpadKeyCreated = assets.TextDesc(assets.TextDescKeyWriteInitScratchpadKeyCreated)

// TplInitCreatingRootFiles is the heading before project root file creation.
var TplInitCreatingRootFiles = assets.TextDesc(assets.TextDescKeyWriteInitCreatingRootFiles)

// TplInitSettingUpPermissions is the heading before permissions setup.
var TplInitSettingUpPermissions = assets.TextDesc(assets.TextDescKeyWriteInitSettingUpPermissions)

// TplInitGitignoreUpdated reports .gitignore entries were added.
// Arguments: count of entries added.
var TplInitGitignoreUpdated = assets.TextDesc(assets.TextDescKeyWriteInitGitignoreUpdated)

// TplInitGitignoreReview hints how to review the .gitignore changes.
var TplInitGitignoreReview = assets.TextDesc(assets.TextDescKeyWriteInitGitignoreReview)

// TplInitNextSteps is the next-steps guidance block after init completes.
var TplInitNextSteps = assets.TextDesc(assets.TextDescKeyWriteInitNextSteps)

// TplInitPluginInfo is the plugin installation guidance block.
var TplInitPluginInfo = assets.TextDesc(assets.TextDescKeyWriteInitPluginInfo)

// TplInitPluginNote is the note about local plugin enabling.
var TplInitPluginNote = assets.TextDesc(assets.TextDescKeyWriteInitPluginNote)

// TplInitCtxContentExists reports a file skipped because ctx content exists.
// Arguments: path.
var TplInitCtxContentExists = assets.TextDesc(
	assets.TextDescKeyWriteInitCtxContentExists,
)

// TplInitMerged reports a file merged during init.
// Arguments: path.
var TplInitMerged = assets.TextDesc(assets.TextDescKeyWriteInitMerged)

// TplInitBackup reports a backup file created.
// Arguments: backup path.
var TplInitBackup = assets.TextDesc(assets.TextDescKeyWriteInitBackup)

// TplInitUpdatedCtxSection reports a file whose ctx section was updated.
// Arguments: path.
var TplInitUpdatedCtxSection = assets.TextDesc(
	assets.TextDescKeyWriteInitUpdatedCtxSection,
)

// TplInitUpdatedPlanSection reports a file whose plan section was updated.
// Arguments: path.
var TplInitUpdatedPlanSection = assets.TextDesc(
	assets.TextDescKeyWriteInitUpdatedPlanSection,
)

// TplInitUpdatedPromptSection reports a file whose prompt section was updated.
// Arguments: path.
var TplInitUpdatedPromptSection = assets.TextDesc(
	assets.TextDescKeyWriteInitUpdatedPromptSection,
)

// TplInitFileExistsNoCtx reports a file exists without ctx content.
// Arguments: path.
var TplInitFileExistsNoCtx = assets.TextDesc(
	assets.TextDescKeyWriteInitFileExistsNoCtx,
)

// TplInitNoChanges reports a settings file with no changes needed.
// Arguments: path.
var TplInitNoChanges = assets.TextDesc(assets.TextDescKeyWriteInitNoChanges)

// TplInitPermsMergedDeduped reports permissions merged and deduped.
// Arguments: path.
var TplInitPermsMergedDeduped = assets.TextDesc(
	assets.TextDescKeyWriteInitPermsMergedDeduped,
)

// TplInitPermsDeduped reports duplicate permissions removed.
// Arguments: path.
var TplInitPermsDeduped = assets.TextDesc(
	assets.TextDescKeyWriteInitPermsDeduped,
)

// TplInitPermsAllowDeny reports allow+deny permissions added.
// Arguments: path.
var TplInitPermsAllowDeny = assets.TextDesc(
	assets.TextDescKeyWriteInitPermsAllowDeny,
)

// TplInitPermsDeny reports deny permissions added.
// Arguments: path.
var TplInitPermsDeny = assets.TextDesc(assets.TextDescKeyWriteInitPermsDeny)

// TplInitPermsAllow reports ctx permissions added.
// Arguments: path.
var TplInitPermsAllow = assets.TextDesc(assets.TextDescKeyWriteInitPermsAllow)

// TplInitMakefileCreated is printed when a new Makefile is created.
var TplInitMakefileCreated = assets.TextDesc(
	assets.TextDescKeyWriteInitMakefileCreated,
)

// TplInitMakefileIncludes reports Makefile already includes the directive.
// Arguments: filename.
var TplInitMakefileIncludes = assets.TextDesc(
	assets.TextDescKeyWriteInitMakefileIncludes,
)

// TplInitMakefileAppended reports an include appended to Makefile.
// Arguments: filename.
var TplInitMakefileAppended = assets.TextDesc(
	assets.TextDescKeyWriteInitMakefileAppended,
)

// TplInitPluginSkipped is printed when plugin enablement is skipped.
var TplInitPluginSkipped = assets.TextDesc(
	assets.TextDescKeyWriteInitPluginSkipped,
)

// TplInitPluginAlreadyEnabled is printed when plugin is already enabled.
var TplInitPluginAlreadyEnabled = assets.TextDesc(
	assets.TextDescKeyWriteInitPluginAlreadyEnabled,
)

// TplInitPluginEnabled reports plugin enabled globally.
// Arguments: settings path.
var TplInitPluginEnabled = assets.TextDesc(
	assets.TextDescKeyWriteInitPluginEnabled,
)

// TplInitSkippedDir reports a directory skipped because it exists.
// Arguments: dir.
var TplInitSkippedDir = assets.TextDesc(
	assets.TextDescKeyWriteInitSkippedDir,
)

// TplInitCreatedDir reports a directory created during init.
// Arguments: dir.
var TplInitCreatedDir = assets.TextDesc(
	assets.TextDescKeyWriteInitCreatedDir,
)

// TplInitCreatedWith reports a file created with a qualifier.
// Arguments: path, qualifier.
var TplInitCreatedWith = assets.TextDesc(
	assets.TextDescKeyWriteInitCreatedWith,
)

// TplInitSkippedPlain reports a file skipped without detail.
// Arguments: path.
var TplInitSkippedPlain = assets.TextDesc(
	assets.TextDescKeyWriteInitSkippedPlain,
)

// TplObsidianGenerated reports successful Obsidian vault generation.
// Arguments: entry count, output directory.
var TplObsidianGenerated = assets.TextDesc(
	assets.TextDescKeyWriteObsidianGenerated,
)

// TplObsidianNextSteps is the post-generation guidance.
// Arguments: output directory.
var TplObsidianNextSteps = assets.TextDesc(
	assets.TextDescKeyWriteObsidianNextSteps,
)

// TplJournalOrphanRemoved reports a removed orphan file.
// Arguments: filename.
var TplJournalOrphanRemoved = assets.TextDesc(
	assets.TextDescKeyWriteJournalOrphanRemoved,
)

// TplJournalSiteGenerated reports successful site generation.
// Arguments: entry count, output directory.
var TplJournalSiteGenerated = assets.TextDesc(
	assets.TextDescKeyWriteJournalSiteGenerated,
)

// TplJournalSiteStarting reports the server is starting.
var TplJournalSiteStarting = assets.TextDesc(
	assets.TextDescKeyWriteJournalSiteStarting,
)

// TplJournalSiteBuilding reports a build is in progress.
var TplJournalSiteBuilding = assets.TextDesc(
	assets.TextDescKeyWriteJournalSiteBuilding,
)

// TplJournalSiteNextSteps shows post-generation guidance.
// Arguments: output directory, zensical binary name.
var TplJournalSiteNextSteps = assets.TextDesc(
	assets.TextDescKeyWriteJournalSiteNextSteps,
)

// TplJournalSiteAlt is the alternative command hint.
var TplJournalSiteAlt = assets.TextDesc(
	assets.TextDescKeyWriteJournalSiteAlt,
)

// TplLoopGenerated reports successful loop script generation.
// Arguments: output file path.
var TplLoopGenerated = assets.TextDesc(
	assets.TextDescKeyWriteLoopGenerated,
)

// TplLoopRunCmd shows how to run the generated script.
// Arguments: output file path.
var TplLoopRunCmd = assets.TextDesc(
	assets.TextDescKeyWriteLoopRunCmd,
)

// TplLoopTool shows the selected tool.
// Arguments: tool name.
var TplLoopTool = assets.TextDesc(assets.TextDescKeyWriteLoopTool)

// TplLoopPrompt shows the prompt file.
// Arguments: prompt file path.
var TplLoopPrompt = assets.TextDesc(assets.TextDescKeyWriteLoopPrompt)

// TplLoopMaxIterations shows the max iterations setting.
// Arguments: count.
var TplLoopMaxIterations = assets.TextDesc(
	assets.TextDescKeyWriteLoopMaxIterations,
)

// TplLoopUnlimited shows unlimited iterations.
var TplLoopUnlimited = assets.TextDesc(assets.TextDescKeyWriteLoopUnlimited)

// TplLoopCompletion shows the completion signal.
// Arguments: signal string.
var TplLoopCompletion = assets.TextDesc(assets.TextDescKeyWriteLoopCompletion)

// TplUnpublishNotFound reports no published block was found.
// Arguments: source filename.
var TplUnpublishNotFound = assets.TextDesc(
	assets.TextDescKeyWriteUnpublishNotFound,
)

// TplUnpublishDone reports the published block was removed.
// Arguments: source filename.
var TplUnpublishDone = assets.TextDesc(assets.TextDescKeyWriteUnpublishDone)

// TplPublishHeader reports publishing has started.
var TplPublishHeader = assets.TextDesc(assets.TextDescKeyWritePublishHeader)

// TplPublishSourceFiles lists the source files used for publishing.
var TplPublishSourceFiles = assets.TextDesc(
	assets.TextDescKeyWritePublishSourceFiles,
)

// TplPublishBudget reports the line budget.
// Arguments: budget.
var TplPublishBudget = assets.TextDesc(assets.TextDescKeyWritePublishBudget)

// TplPublishBlock is the heading for the published block detail.
var TplPublishBlock = assets.TextDesc(assets.TextDescKeyWritePublishBlock)

// TplPublishTasks reports pending tasks count.
// Arguments: count.
var TplPublishTasks = assets.TextDesc(assets.TextDescKeyWritePublishTasks)

// TplPublishDecisions reports recent decisions count.
// Arguments: count.
var TplPublishDecisions = assets.TextDesc(
	assets.TextDescKeyWritePublishDecisions,
)

// TplPublishConventions reports key conventions count.
// Arguments: count.
var TplPublishConventions = assets.TextDesc(
	assets.TextDescKeyWritePublishConventions,
)

// TplPublishLearnings reports recent learnings count.
// Arguments: count.
var TplPublishLearnings = assets.TextDesc(
	assets.TextDescKeyWritePublishLearnings,
)

// TplPublishTotal reports the total line count within budget.
// Arguments: total lines, budget.
var TplPublishTotal = assets.TextDesc(assets.TextDescKeyWritePublishTotal)

// TplPublishDryRun reports a publish dry run.
var TplPublishDryRun = assets.TextDesc(assets.TextDescKeyWritePublishDryRun)

// TplPublishDone reports successful publishing with marker info.
var TplPublishDone = assets.TextDesc(assets.TextDescKeyWritePublishDone)

// TplImportNoEntries reports no entries found in MEMORY.md.
var TplImportNoEntries = assets.TextDesc(assets.TextDescKeyWriteImportNoEntries)

// TplImportScanning reports scanning has started.
// Arguments: source filename.
var TplImportScanning = assets.TextDesc(assets.TextDescKeyWriteImportScanning)

// TplImportFound reports the number of entries found.
// Arguments: count.
var TplImportFound = assets.TextDesc(assets.TextDescKeyWriteImportFound)

// TplImportEntry reports an entry being processed.
// Arguments: truncated title (already quoted).
var TplImportEntry = assets.TextDesc(assets.TextDescKeyWriteImportEntry)

// TplImportClassifiedSkip reports an entry classified as skip.
var TplImportClassifiedSkip = assets.TextDesc(
	assets.TextDescKeyWriteImportClassifiedSkip,
)

// TplImportClassified reports an entry classification.
// Arguments: target file, comma-joined keywords.
var TplImportClassified = assets.TextDesc(assets.TextDescKeyWriteImportClassified)

// TplImportAdded reports an entry added to a target file.
// Arguments: target filename.
var TplImportAdded = assets.TextDesc(assets.TextDescKeyWriteImportAdded)

// TplImportSummaryDryRun is the dry-run summary prefix.
// Arguments: count.
var TplImportSummaryDryRun = assets.TextDesc(
	assets.TextDescKeyWriteImportSummaryDryRun,
)

// TplImportSummary is the import summary prefix.
// Arguments: count.
var TplImportSummary = assets.TextDesc(assets.TextDescKeyWriteImportSummary)

// TplImportSkipped reports skipped entries.
// Arguments: count.
var TplImportSkipped = assets.TextDesc(assets.TextDescKeyWriteImportSkipped)

// TplImportDuplicates reports duplicate entries.
// Arguments: count.
var TplImportDuplicates = assets.TextDesc(assets.TextDescKeyWriteImportDuplicates)

// TplMemoryNoChanges reports no changes since last sync.
var TplMemoryNoChanges = assets.TextDesc(assets.TextDescKeyWriteMemoryNoChanges)

// TplMemoryBridgeHeader is the heading for memory status output.
var TplMemoryBridgeHeader = assets.TextDesc(
	assets.TextDescKeyWriteMemoryBridgeHeader,
)

// TplMemorySourceNotActive reports that auto memory is not active.
var TplMemorySourceNotActive = assets.TextDesc(
	assets.TextDescKeyWriteMemorySourceNotActive,
)

// TplMemorySource is a format template for the source path.
// Arguments: path.
var TplMemorySource = assets.TextDesc(assets.TextDescKeyWriteMemorySource)

// TplMemoryMirror is a format template for the mirror relative path.
// Arguments: relative path.
var TplMemoryMirror = assets.TextDesc(assets.TextDescKeyWriteMemoryMirror)

// TplMemoryLastSync is a format template for the last sync time.
// Arguments: formatted time, human-readable duration.
var TplMemoryLastSync = assets.TextDesc(assets.TextDescKeyWriteMemoryLastSync)

// TplMemoryLastSyncNever reports no sync has occurred.
var TplMemoryLastSyncNever = assets.TextDesc(
	assets.TextDescKeyWriteMemoryLastSyncNever,
)

// TplMemorySourceLines is a format template for MEMORY.md line count.
// Arguments: line count.
var TplMemorySourceLines = assets.TextDesc(assets.TextDescKeyWriteMemorySourceLines)

// TplMemorySourceLinesDrift is a format template for MEMORY.md line count
// when drift is detected. Arguments: line count.
var TplMemorySourceLinesDrift = assets.TextDesc(
	assets.TextDescKeyWriteMemorySourceLinesDrift,
)

// TplMemoryMirrorLines is a format template for mirror line count.
// Arguments: line count.
var TplMemoryMirrorLines = assets.TextDesc(
	assets.TextDescKeyWriteMemoryMirrorLines,
)

// TplMemoryMirrorNotSynced reports the mirror has not been synced.
var TplMemoryMirrorNotSynced = assets.TextDesc(
	assets.TextDescKeyWriteMemoryMirrorNotSynced,
)

// TplMemoryDriftDetected reports drift was detected.
var TplMemoryDriftDetected = assets.TextDesc(
	assets.TextDescKeyWriteMemoryDriftDetected,
)

// TplMemoryDriftNone reports no drift.
var TplMemoryDriftNone = assets.TextDesc(assets.TextDescKeyWriteMemoryDriftNone)

// TplMemoryArchives is a format template for archive snapshot count.
// Arguments: count, archive directory name.
var TplMemoryArchives = assets.TextDesc(assets.TextDescKeyWriteMemoryArchives)

// TplPadEntryAdded is a format template for pad entry confirmation.
// Arguments: entry number.
var TplPadEntryAdded = assets.TextDesc(assets.TextDescKeyWritePadEntryAdded)

// TplPadEntryUpdated is a format template for pad entry update confirmation.
// Arguments: entry number.
var TplPadEntryUpdated = assets.TextDesc(assets.TextDescKeyWritePadEntryUpdated)

// TplPadExportPlan is a format template for a dry-run export line.
// Arguments: label, output path.
var TplPadExportPlan = assets.TextDesc(assets.TextDescKeyWritePadExportPlan)

// TplPadExportDone is a format template for a successfully exported blob.
// Arguments: label.
var TplPadExportDone = assets.TextDesc(assets.TextDescKeyWritePadExportDone)

// TplPadExportWriteFailed is a format template for a failed blob write (stderr).
// Arguments: label, error.
var TplPadExportWriteFailed = assets.TextDesc(
	assets.TextDescKeyWritePadExportWriteFailed,
)

// TplPadExportNone is the message when no blob entries exist to export.
var TplPadExportNone = assets.TextDesc(assets.TextDescKeyWritePadExportNone)

// TplPadExportSummary is a format template for the export summary.
// Arguments: verb ("Exported"/"Would export"), count.
var TplPadExportSummary = assets.TextDesc(assets.TextDescKeyWritePadExportSummary)

// TplPadExportVerbDone is the past-tense verb for export summary.
var TplPadExportVerbDone = assets.TextDesc(assets.TextDescKeyWritePadExportVerbDone)

// TplPadExportVerbDryRun is the dry-run verb for export summary.
var TplPadExportVerbDryRun = assets.TextDesc(
	assets.TextDescKeyWritePadExportVerbDryRun,
)

// TplPadImportNone is the message when no entries were found to import.
var TplPadImportNone = assets.TextDesc(assets.TextDescKeyWritePadImportNone)

// TplPadImportDone is a format template for successful line import.
// Arguments: count.
var TplPadImportDone = assets.TextDesc(assets.TextDescKeyWritePadImportDone)

// TplPadImportBlobAdded is a format template for a successfully imported blob.
// Arguments: filename.
var TplPadImportBlobAdded = assets.TextDesc(
	assets.TextDescKeyWritePadImportBlobAdded,
)

// TplPadImportBlobSkipped is a format template for a skipped blob (stderr).
// Arguments: filename, reason.
var TplPadImportBlobSkipped = assets.TextDesc(
	assets.TextDescKeyWritePadImportBlobSkipped,
)

// TplPadImportBlobTooLarge is a format template for a blob exceeding the size limit (stderr).
// Arguments: filename, max bytes.
var TplPadImportBlobTooLarge = assets.TextDesc(
	assets.TextDescKeyWritePadImportBlobTooLarge,
)

// TplPadImportBlobNone is the message when no files were found to import.
var TplPadImportBlobNone = assets.TextDesc(
	assets.TextDescKeyWritePadImportBlobNone,
)

// TplPadImportBlobSummary is a format template for blob import summary.
// Arguments: added count, skipped count.
var TplPadImportBlobSummary = assets.TextDesc(
	assets.TextDescKeyWritePadImportBlobSummary,
)

// TplPadImportCloseWarning is a format template for file close warning (stderr).
// Arguments: filename, error.
var TplPadImportCloseWarning = assets.TextDesc(
	assets.TextDescKeyWritePadImportCloseWarning,
)

// TplPaused is a format template for the pause confirmation.
// Arguments: session ID.
var TplPaused = assets.TextDesc(assets.TextDescKeyWritePaused)

// TplRestoreNoLocal is the message when golden is restored with no local file.
var TplRestoreNoLocal = assets.TextDesc(assets.TextDescKeyWriteRestoreNoLocal)

// TplRestoreMatch is the message when settings already match golden.
var TplRestoreMatch = assets.TextDesc(assets.TextDescKeyWriteRestoreMatch)

// TplRestoreDroppedHeader is a format template for dropped permissions header.
// Arguments: count.
var TplRestoreDroppedHeader = assets.TextDesc(
	assets.TextDescKeyWriteRestoreDroppedHeader,
)

// TplRestoreRestoredHeader is a format template for restored permissions header.
// Arguments: count.
var TplRestoreRestoredHeader = assets.TextDesc(
	assets.TextDescKeyWriteRestoreRestoredHeader,
)

// TplRestoreDenyDroppedHeader is a format template for dropped deny rules header.
// Arguments: count.
var TplRestoreDenyDroppedHeader = assets.TextDesc(
	assets.TextDescKeyWriteRestoreDenyDroppedHeader,
)

// TplRestoreDenyRestoredHeader is a format template for restored deny rules header.
// Arguments: count.
var TplRestoreDenyRestoredHeader = assets.TextDesc(
	assets.TextDescKeyWriteRestoreDenyRestoredHeader,
)

// TplRestoreRemoved is a format template for a removed permission line.
// Arguments: permission string.
var TplRestoreRemoved = assets.TextDesc(assets.TextDescKeyWriteRestoreRemoved)

// TplRestoreAdded is a format template for an added permission line.
// Arguments: permission string.
var TplRestoreAdded = assets.TextDesc(assets.TextDescKeyWriteRestoreAdded)

// TplRestorePermMatch is the message when only non-permission settings differ.
var TplRestorePermMatch = assets.TextDesc(assets.TextDescKeyWriteRestorePermMatch)

// TplRestoreDone is the message after successful restore.
var TplRestoreDone = assets.TextDesc(assets.TextDescKeyWriteRestoreDone)

// TplSnapshotSaved is a format template for golden image save.
// Arguments: golden file path.
var TplSnapshotSaved = assets.TextDesc(assets.TextDescKeyWriteSnapshotSaved)

// TplSnapshotUpdated is a format template for golden image update.
// Arguments: golden file path.
var TplSnapshotUpdated = assets.TextDesc(assets.TextDescKeyWriteSnapshotUpdated)

// TplResumed is a format template for the resume confirmation.
// Arguments: session ID.
var TplResumed = assets.TextDesc(assets.TextDescKeyWriteResumed)

// TplPadEmpty is the message when the scratchpad has no entries.
var TplPadEmpty = assets.TextDesc(assets.TextDescKeyWritePadEmpty)

// TplPadKeyCreated is a format template for key creation notice (stderr).
// Arguments: key file path.
var TplPadKeyCreated = assets.TextDesc(assets.TextDescKeyWritePadKeyCreated)

// TplPadBlobWritten is a format template for blob file write confirmation.
// Arguments: byte count, output path.
var TplPadBlobWritten = assets.TextDesc(assets.TextDescKeyWritePadBlobWritten)

// TplPadEntryRemoved is a format template for pad entry removal confirmation.
// Arguments: entry number.
var TplPadEntryRemoved = assets.TextDesc(assets.TextDescKeyWritePadEntryRemoved)

// TplPadResolveHeader is a format template for a conflict side header.
// Arguments: side label ("OURS"/"THEIRS").
var TplPadResolveHeader = assets.TextDesc(assets.TextDescKeyWritePadResolveHeader)

// TplPadResolveEntry is a format template for a numbered conflict entry.
// Arguments: 1-based index, display string.
var TplPadResolveEntry = assets.TextDesc(assets.TextDescKeyWritePadResolveEntry)

// TplPadEntryMoved is a format template for pad entry move confirmation.
// Arguments: source position, destination position.
var TplPadEntryMoved = assets.TextDesc(assets.TextDescKeyWritePadEntryMoved)

// TplPadMergeDupe is a format template for a duplicate entry during merge.
// Arguments: display string.
var TplPadMergeDupe = assets.TextDesc(assets.TextDescKeyWritePadMergeDupe)

// TplPadMergeAdded is a format template for a newly added entry during merge.
// Arguments: display string, source file.
var TplPadMergeAdded = assets.TextDesc(assets.TextDescKeyWritePadMergeAdded)

// TplPadMergeBlobConflict is a format template for a blob label conflict warning.
// Arguments: label.
var TplPadMergeBlobConflict = assets.TextDesc(
	assets.TextDescKeyWritePadMergeBlobConflict,
)

// TplPadMergeBinaryWarning is a format template for a binary data warning.
// Arguments: filename.
var TplPadMergeBinaryWarning = assets.TextDesc(
	assets.TextDescKeyWritePadMergeBinaryWarning,
)

// TplPadMergeNone is the message when no entries were found to merge.
var TplPadMergeNone = assets.TextDesc(assets.TextDescKeyWritePadMergeNone)

// TplPadMergeNoneNew is a format template when all entries are duplicates.
// Arguments: dupe count, pluralized "duplicate".
var TplPadMergeNoneNew = assets.TextDesc(assets.TextDescKeyWritePadMergeNoneNew)

// TplPadMergeDryRun is a format template for dry-run merge summary.
// Arguments: added count, pluralized "entry", dupe count, pluralized "duplicate".
var TplPadMergeDryRun = assets.TextDesc(assets.TextDescKeyWritePadMergeDryRun)

// TplPadMergeDone is a format template for successful merge summary.
// Arguments: added count, pluralized "entry", dupe count, pluralized "duplicate".
var TplPadMergeDone = assets.TextDesc(assets.TextDescKeyWritePadMergeDone)

// TplSetupPrompt is the interactive prompt for webhook URL entry.
var TplSetupPrompt = assets.TextDesc(assets.TextDescKeyWriteSetupPrompt)

// TplSetupDone is a format template for successful webhook configuration.
// Arguments: masked URL, encrypted file path.
var TplSetupDone = assets.TextDesc(assets.TextDescKeyWriteSetupDone)

// TplTestNoWebhook is the message when no webhook is configured.
var TplTestNoWebhook = assets.TextDesc(assets.TextDescKeyWriteTestNoWebhook)

// TplTestFiltered is the notice when the test event is filtered.
var TplTestFiltered = assets.TextDesc(assets.TextDescKeyWriteTestFiltered)

// TplTestResult is a format template for webhook test response.
// Arguments: HTTP status code, status text.
var TplTestResult = assets.TextDesc(assets.TextDescKeyWriteTestResult)

// TplTestWorking is the success message after a webhook test.
// Arguments: encrypted file path.
var TplTestWorking = assets.TextDesc(assets.TextDescKeyWriteTestWorking)

// TplPromptCreated is the confirmation after creating a prompt template.
// Arguments: prompt name.
var TplPromptCreated = assets.TextDesc(assets.TextDescKeyWritePromptCreated)

// TplPromptNone is printed when no prompts are found.
var TplPromptNone = assets.TextDesc(assets.TextDescKeyWritePromptNone)

// TplPromptItem is a format template for listing a prompt name.
// Arguments: prompt name.
var TplPromptItem = assets.TextDesc(assets.TextDescKeyWritePromptItem)

// TplPromptRemoved is the confirmation after removing a prompt template.
// Arguments: prompt name.
var TplPromptRemoved = assets.TextDesc(assets.TextDescKeyWritePromptRemoved)

// TplReminderAdded is the confirmation for a newly added reminder.
// Arguments: id, message, suffix (e.g. "  (after 2026-03-10)" or "").
var TplReminderAdded = assets.TextDesc(assets.TextDescKeyWriteReminderAdded)

// TplReminderDismissed is the confirmation for a dismissed reminder.
// Arguments: id, message.
var TplReminderDismissed = assets.TextDesc(assets.TextDescKeyWriteReminderDismissed)

// TplReminderNone is printed when there are no reminders.
var TplReminderNone = assets.TextDesc(assets.TextDescKeyWriteReminderNone)

// TplReminderDismissedAll is the summary after dismissing all reminders.
// Arguments: count.
var TplReminderDismissedAll = assets.TextDesc(
	assets.TextDescKeyWriteReminderDismissedAll,
)

// TplReminderItem is a format template for listing a reminder.
// Arguments: id, message, annotation.
var TplReminderItem = assets.TextDesc(assets.TextDescKeyWriteReminderItem)

// TplReminderNotDue is the annotation for reminders not yet due.
// Arguments: date string.
var TplReminderNotDue = assets.TextDesc(assets.TextDescKeyWriteReminderNotDue)

// TplReminderAfterSuffix formats the date-gate suffix for a reminder.
// Arguments: date string.
var TplReminderAfterSuffix = assets.TextDesc(
	assets.TextDescKeyWriteReminderAfterSuffix,
)

// TplLockUnlockEntry is the confirmation for a single locked/unlocked entry.
// Arguments: filename, verb ("locked" or "unlocked").
var TplLockUnlockEntry = assets.TextDesc(assets.TextDescKeyWriteLockUnlockEntry)

// TplLockUnlockNoChanges is printed when all entries already have the target state.
// Arguments: verb.
var TplLockUnlockNoChanges = assets.TextDesc(
	assets.TextDescKeyWriteLockUnlockNoChanges,
)

// TplLockUnlockSummary is the summary after locking/unlocking entries.
// Arguments: capitalized verb, count.
var TplLockUnlockSummary = assets.TextDesc(assets.TextDescKeyWriteLockUnlockSummary)

// TplBackupResult is a format template for a backup result line.
// Arguments: scope, archive path, formatted size.
var TplBackupResult = assets.TextDesc(assets.TextDescKeyWriteBackupResult)

// TplBackupSMBDest is a format template for the SMB destination suffix.
// Arguments: SMB destination path.
var TplBackupSMBDest = assets.TextDesc(assets.TextDescKeyWriteBackupSMBDest)

// TplStatusTitle is the heading for the status output.
var TplStatusTitle = assets.TextDesc(assets.TextDescKeyWriteStatusTitle)

// TplStatusSeparator is the visual separator under the heading.
var TplStatusSeparator = assets.TextDesc(assets.TextDescKeyWriteStatusSeparator)

// TplStatusDir is a format template for the context directory.
// Arguments: context directory path.
var TplStatusDir = assets.TextDesc(assets.TextDescKeyWriteStatusDir)

// TplStatusFiles is a format template for the total file count.
// Arguments: count.
var TplStatusFiles = assets.TextDesc(assets.TextDescKeyWriteStatusFiles)

// TplStatusTokens is a format template for the token estimate.
// Arguments: formatted token count.
var TplStatusTokens = assets.TextDesc(assets.TextDescKeyWriteStatusTokens)

// TplStatusFilesHeader is the heading for the file list section.
var TplStatusFilesHeader = assets.TextDesc(
	assets.TextDescKeyWriteStatusFilesHeader,
)

// TplStatusFileVerbose is a format template for a verbose file entry.
// Arguments: indicator, name, status, formatted tokens, formatted size.
var TplStatusFileVerbose = assets.TextDesc(
	assets.TextDescKeyWriteStatusFileVerbose,
)

// TplStatusFileCompact is a format template for a compact file entry.
// Arguments: indicator, name, status.
var TplStatusFileCompact = assets.TextDesc(
	assets.TextDescKeyWriteStatusFileCompact,
)

// TplStatusPreviewLine is a format template for a content preview line.
// Arguments: line text.
var TplStatusPreviewLine = assets.TextDesc(
	assets.TextDescKeyWriteStatusPreviewLine,
)

// TplStatusActivityHeader is the heading for the recent activity section.
var TplStatusActivityHeader = assets.TextDesc(
	assets.TextDescKeyWriteStatusActivityHeader,
)

// TplStatusActivityItem is a format template for a recent activity entry.
// Arguments: filename, relative time string.
var TplStatusActivityItem = assets.TextDesc(
	assets.TextDescKeyWriteStatusActivityItem,
)

// TplTimeJustNow is the display string for "just now" relative time.
var TplTimeJustNow = assets.TextDesc(assets.TextDescKeyWriteTimeJustNow)

// TplTimeMinuteAgo is the display string for "1 minute ago".
var TplTimeMinuteAgo = assets.TextDesc(assets.TextDescKeyWriteTimeMinuteAgo)

// TplTimeMinutesAgo is a format template for minutes ago.
// Arguments: count.
var TplTimeMinutesAgo = assets.TextDesc(assets.TextDescKeyWriteTimeMinutesAgo)

// TplTimeHourAgo is the display string for "1 hour ago".
var TplTimeHourAgo = assets.TextDesc(assets.TextDescKeyWriteTimeHourAgo)

// TplTimeHoursAgo is a format template for hours ago.
// Arguments: count.
var TplTimeHoursAgo = assets.TextDesc(assets.TextDescKeyWriteTimeHoursAgo)

// TplTimeDayAgo is the display string for "1 day ago".
var TplTimeDayAgo = assets.TextDesc(assets.TextDescKeyWriteTimeDayAgo)

// TplTimeDaysAgo is a format template for days ago.
// Arguments: count.
var TplTimeDaysAgo = assets.TextDesc(assets.TextDescKeyWriteTimeDaysAgo)

// TplTimeOlderFormat is the Go time layout for dates older than a week.
// Exported because callers must format the fallback date before calling FormatTimeAgo.
//
// Deprecated: Use config.OlderFormat instead.
const TplTimeOlderFormat = time.OlderFormat

// TplSyncInSync is printed when context is fully in sync.
var TplSyncInSync = assets.TextDesc(assets.TextDescKeyWriteSyncInSync)

// TplSyncHeader is the heading for the sync analysis output.
var TplSyncHeader = assets.TextDesc(assets.TextDescKeyWriteSyncHeader)

// TplSyncSeparator is the visual separator under the heading.
var TplSyncSeparator = assets.TextDesc(assets.TextDescKeyWriteSyncSeparator)

// TplSyncDryRun is printed when running in dry-run mode.
var TplSyncDryRun = assets.TextDesc(assets.TextDescKeyWriteSyncDryRun)

// TplSyncAction is a format template for a sync action item.
// Arguments: index, type, description.
var TplSyncAction = assets.TextDesc(assets.TextDescKeyWriteSyncAction)

// TplSyncSuggestion is a format template for a suggestion under an action.
// Arguments: suggestion text.
var TplSyncSuggestion = assets.TextDesc(assets.TextDescKeyWriteSyncSuggestion)

// TplSyncDryRunSummary is a format template for dry-run summary.
// Arguments: count.
var TplSyncDryRunSummary = assets.TextDesc(
	assets.TextDescKeyWriteSyncDryRunSummary,
)

// TplSyncSummary is a format template for the sync summary.
// Arguments: count.
var TplSyncSummary = assets.TextDesc(assets.TextDescKeyWriteSyncSummary)

// TplJournalSyncNone is printed when no journal entries are found.
var TplJournalSyncNone = assets.TextDesc(assets.TextDescKeyWriteJournalSyncNone)

// TplJournalSyncLocked is a format template for a newly locked entry.
// Arguments: filename.
var TplJournalSyncLocked = assets.TextDesc(
	assets.TextDescKeyWriteJournalSyncLocked,
)

// TplJournalSyncUnlocked is a format template for a newly unlocked entry.
// Arguments: filename.
var TplJournalSyncUnlocked = assets.TextDesc(
	assets.TextDescKeyWriteJournalSyncUnlocked,
)

// TplJournalSyncMatch is printed when state already matches frontmatter.
var TplJournalSyncMatch = assets.TextDesc(
	assets.TextDescKeyWriteJournalSyncMatch,
)

// TplJournalSyncLockedCount is a format template for locked entry count.
// Arguments: count.
var TplJournalSyncLockedCount = assets.TextDesc(
	assets.TextDescKeyWriteJournalSyncLockedCount,
)

// TplJournalSyncUnlockedCount is a format template for unlocked entry count.
// Arguments: count.
var TplJournalSyncUnlockedCount = assets.TextDesc(
	assets.TextDescKeyWriteJournalSyncUnlockedCount,
)
