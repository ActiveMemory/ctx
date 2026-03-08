//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// tplBootstrapTitle is the heading for bootstrap output.
var tplBootstrapTitle = assets.TextDesc(assets.TextDescKeyWriteBootstrapTitle)

// tplBootstrapSep is the visual separator under the bootstrap heading.
var tplBootstrapSep = assets.TextDesc(assets.TextDescKeyWriteBootstrapSep)

// tplBootstrapDir is a format template for the context directory.
// Arguments: context directory path.
var tplBootstrapDir = assets.TextDesc(assets.TextDescKeyWriteBootstrapDir)

// tplBootstrapFiles is the heading for the file list section.
var tplBootstrapFiles = assets.TextDesc(assets.TextDescKeyWriteBootstrapFiles)

// tplBootstrapRules is the heading for the rules section.
var tplBootstrapRules = assets.TextDesc(assets.TextDescKeyWriteBootstrapRules)

// tplBootstrapNextSteps is the heading for the next steps section.
var tplBootstrapNextSteps = assets.TextDesc(assets.TextDescKeyWriteBootstrapNextSteps)

// tplBootstrapNumbered is a format template for a numbered list item.
// Arguments: index, text.
var tplBootstrapNumbered = assets.TextDesc(assets.TextDescKeyWriteBootstrapNumbered)

// tplBootstrapWarning is a format template for a warning line.
// Arguments: warning text.
var tplBootstrapWarning = assets.TextDesc(assets.TextDescKeyWriteBootstrapWarning)

// prefixError is prepended to all error messages written to stderr.
var prefixError = assets.TextDesc(assets.TextDescKeyWritePrefixError)

// tplPathExists is a format template for reporting that a destination path
// already exists. Arguments: original path, resolved destination path.
var tplPathExists = assets.TextDesc(assets.TextDescKeyWritePathExists)

// tplExistsWritingAsAlternative is a format template for reporting that a
// file exists and content was written to an alternative filename instead.
// Arguments: original path, alternative path.
var tplExistsWritingAsAlternative = assets.TextDesc(assets.TextDescKeyWriteExistsWritingAsAlternative)

// tplDryRun is printed when a command runs in dry-run mode.
var tplDryRun = assets.TextDesc(assets.TextDescKeyWriteDryRun)

// tplSource is a format template for reporting a source path.
// Arguments: path.
var tplSource = assets.TextDesc(assets.TextDescKeyWriteSource)

// tplMirror is a format template for reporting a mirror path.
// Arguments: relative mirror path.
var tplMirror = assets.TextDesc(assets.TextDescKeyWriteMirror)

// tplStatusDrift is printed when drift is detected.
var tplStatusDrift = assets.TextDesc(assets.TextDescKeyWriteStatusDrift)

// tplStatusNoDrift is printed when no drift is detected.
var tplStatusNoDrift = assets.TextDesc(assets.TextDescKeyWriteStatusNoDrift)

// tplArchived is a format template for reporting an archived file.
// Arguments: archive filename.
var tplArchived = assets.TextDesc(assets.TextDescKeyWriteArchived)

// tplSynced is a format template for reporting a successful sync.
// Arguments: source label, destination relative path.
var tplSynced = assets.TextDesc(assets.TextDescKeyWriteSynced)

// tplLines is a format template for reporting line counts.
// Arguments: line count.
var tplLines = assets.TextDesc(assets.TextDescKeyWriteLines)

// tplLinesPrevious is a format template appended to line counts when a
// previous count is available. Arguments: previous line count.
var tplLinesPrevious = assets.TextDesc(assets.TextDescKeyWriteLinesPrevious)

// tplNewContent is a format template for reporting new content since last sync.
// Arguments: line count.
var tplNewContent = assets.TextDesc(assets.TextDescKeyWriteNewContent)

// tplAddedTo is a format template for confirming an entry was added.
// Arguments: filename.
var tplAddedTo = assets.TextDesc(assets.TextDescKeyWriteAddedTo)

// tplMovingTask is a format template for a task being moved to completed.
// Arguments: truncated task text.
var tplMovingTask = assets.TextDesc(assets.TextDescKeyWriteMovingTask)

// tplCompletedTask is a format template for a task marked complete.
// Arguments: task text.
var tplCompletedTask = assets.TextDesc(assets.TextDescKeyWriteCompletedTask)

// tplConfigProfileDev is the status output for dev profile.
var tplConfigProfileDev = assets.TextDesc(assets.TextDescKeyWriteConfigProfileDev)

// tplConfigProfileBase is the status output for base profile.
var tplConfigProfileBase = assets.TextDesc(assets.TextDescKeyWriteConfigProfileBase)

// tplConfigProfileNone is the status output when no profile exists.
// Arguments: ctxrc filename.
var tplConfigProfileNone = assets.TextDesc(assets.TextDescKeyWriteConfigProfileNone)

// tplDepsNoProject is printed when no supported project is detected.
var tplDepsNoProject = assets.TextDesc(assets.TextDescKeyWriteDepsNoProject)

// tplDepsLookingFor is printed with the list of files checked.
var tplDepsLookingFor = assets.TextDesc(assets.TextDescKeyWriteDepsLookingFor)

// tplDepsUseType hints at the --type flag.
// Arguments: comma-separated list of builder names.
var tplDepsUseType = assets.TextDesc(assets.TextDescKeyWriteDepsUseType)

// tplDepsNoDeps is printed when no dependencies are found.
var tplDepsNoDeps = assets.TextDesc(assets.TextDescKeyWriteDepsNoDeps)

// tplSkillsHeader is the heading for the skills list.
var tplSkillsHeader = assets.TextDesc(assets.TextDescKeyWriteSkillsHeader)

// tplSkillLine formats a single skill entry.
// Arguments: name, description.
var tplSkillLine = assets.TextDesc(assets.TextDescKeyWriteSkillLine)

// tplHookCopilotSkipped reports that copilot instructions were skipped.
// Arguments: target file path.
var tplHookCopilotSkipped = assets.TextDesc(assets.TextDescKeyWriteHookCopilotSkipped)

// tplHookCopilotForceHint tells the user about the --force flag.
var tplHookCopilotForceHint = assets.TextDesc(assets.TextDescKeyWriteHookCopilotForceHint)

// tplHookCopilotMerged reports that copilot instructions were merged.
// Arguments: target file path.
var tplHookCopilotMerged = assets.TextDesc(assets.TextDescKeyWriteHookCopilotMerged)

// tplHookCopilotCreated reports that copilot instructions were created.
// Arguments: target file path.
var tplHookCopilotCreated = assets.TextDesc(assets.TextDescKeyWriteHookCopilotCreated)

// tplHookCopilotSessionsDir reports that the sessions directory was created.
// Arguments: sessions directory path.
var tplHookCopilotSessionsDir = assets.TextDesc(assets.TextDescKeyWriteHookCopilotSessionsDir)

// tplHookCopilotSummary is the post-write summary for copilot.
var tplHookCopilotSummary = assets.TextDesc(assets.TextDescKeyWriteHookCopilotSummary)

// tplHookUnknownTool reports an unrecognized tool name.
// Arguments: tool name.
var tplHookUnknownTool = assets.TextDesc(assets.TextDescKeyWriteHookUnknownTool)

// tplInitOverwritePrompt prompts the user before overwriting .context/.
// Arguments: context directory path.
var tplInitOverwritePrompt = assets.TextDesc(assets.TextDescKeyWriteInitOverwritePrompt)

// tplInitAborted is printed when the user declines overwriting.
var tplInitAborted = assets.TextDesc(assets.TextDescKeyWriteInitAborted)

// tplInitExistsSkipped reports a file that was skipped because it exists.
// Arguments: filename.
var tplInitExistsSkipped = assets.TextDesc(assets.TextDescKeyWriteInitExistsSkipped)

// tplInitFileCreated reports a file that was successfully created.
// Arguments: filename.
var tplInitFileCreated = assets.TextDesc(assets.TextDescKeyWriteInitFileCreated)

// tplInitialized reports successful context initialization.
// Arguments: context directory path.
var tplInitialized = assets.TextDesc(assets.TextDescKeyWriteInitialized)

// tplInitWarnNonFatal reports a non-fatal warning during init.
// Arguments: label, error.
var tplInitWarnNonFatal = assets.TextDesc(assets.TextDescKeyWriteInitWarnNonFatal)

// tplInitScratchpadPlaintext reports a plaintext scratchpad was created.
// Arguments: path.
var tplInitScratchpadPlaintext = assets.TextDesc(assets.TextDescKeyWriteInitScratchpadPlaintext)

// tplInitScratchpadNoKey warns about a missing key for an encrypted scratchpad.
// Arguments: key path.
var tplInitScratchpadNoKey = assets.TextDesc(assets.TextDescKeyWriteInitScratchpadNoKey)

// tplInitScratchpadKeyCreated reports a scratchpad key was generated.
// Arguments: key path.
var tplInitScratchpadKeyCreated = assets.TextDesc(assets.TextDescKeyWriteInitScratchpadKeyCreated)

// tplInitCreatingRootFiles is the heading before project root file creation.
var tplInitCreatingRootFiles = assets.TextDesc(assets.TextDescKeyWriteInitCreatingRootFiles)

// tplInitSettingUpPermissions is the heading before permissions setup.
var tplInitSettingUpPermissions = assets.TextDesc(assets.TextDescKeyWriteInitSettingUpPermissions)

// tplInitGitignoreUpdated reports .gitignore entries were added.
// Arguments: count of entries added.
var tplInitGitignoreUpdated = assets.TextDesc(assets.TextDescKeyWriteInitGitignoreUpdated)

// tplInitGitignoreReview hints how to review the .gitignore changes.
var tplInitGitignoreReview = assets.TextDesc(assets.TextDescKeyWriteInitGitignoreReview)

// tplInitNextSteps is the next-steps guidance block after init completes.
var tplInitNextSteps = assets.TextDesc(assets.TextDescKeyWriteInitNextSteps)

// tplInitPluginInfo is the plugin installation guidance block.
var tplInitPluginInfo = assets.TextDesc(assets.TextDescKeyWriteInitPluginInfo)

// tplInitPluginNote is the note about local plugin enabling.
var tplInitPluginNote = assets.TextDesc(assets.TextDescKeyWriteInitPluginNote)

// tplInitCtxContentExists reports a file skipped because ctx content exists.
// Arguments: path.
var tplInitCtxContentExists = assets.TextDesc(
	assets.TextDescKeyWriteInitCtxContentExists,
)

// tplInitMerged reports a file merged during init.
// Arguments: path.
var tplInitMerged = assets.TextDesc(assets.TextDescKeyWriteInitMerged)

// tplInitBackup reports a backup file created.
// Arguments: backup path.
var tplInitBackup = assets.TextDesc(assets.TextDescKeyWriteInitBackup)

// tplInitUpdatedCtxSection reports a file whose ctx section was updated.
// Arguments: path.
var tplInitUpdatedCtxSection = assets.TextDesc(
	assets.TextDescKeyWriteInitUpdatedCtxSection,
)

// tplInitUpdatedPlanSection reports a file whose plan section was updated.
// Arguments: path.
var tplInitUpdatedPlanSection = assets.TextDesc(
	assets.TextDescKeyWriteInitUpdatedPlanSection,
)

// tplInitUpdatedPromptSection reports a file whose prompt section was updated.
// Arguments: path.
var tplInitUpdatedPromptSection = assets.TextDesc(
	assets.TextDescKeyWriteInitUpdatedPromptSection,
)

// tplInitFileExistsNoCtx reports a file exists without ctx content.
// Arguments: path.
var tplInitFileExistsNoCtx = assets.TextDesc(
	assets.TextDescKeyWriteInitFileExistsNoCtx,
)

// tplInitNoChanges reports a settings file with no changes needed.
// Arguments: path.
var tplInitNoChanges = assets.TextDesc(assets.TextDescKeyWriteInitNoChanges)

// tplInitPermsMergedDeduped reports permissions merged and deduped.
// Arguments: path.
var tplInitPermsMergedDeduped = assets.TextDesc(
	assets.TextDescKeyWriteInitPermsMergedDeduped,
)

// tplInitPermsDeduped reports duplicate permissions removed.
// Arguments: path.
var tplInitPermsDeduped = assets.TextDesc(
	assets.TextDescKeyWriteInitPermsDeduped,
)

// tplInitPermsAllowDeny reports allow+deny permissions added.
// Arguments: path.
var tplInitPermsAllowDeny = assets.TextDesc(
	assets.TextDescKeyWriteInitPermsAllowDeny,
)

// tplInitPermsDeny reports deny permissions added.
// Arguments: path.
var tplInitPermsDeny = assets.TextDesc(assets.TextDescKeyWriteInitPermsDeny)

// tplInitPermsAllow reports ctx permissions added.
// Arguments: path.
var tplInitPermsAllow = assets.TextDesc(assets.TextDescKeyWriteInitPermsAllow)

// tplInitMakefileCreated is printed when a new Makefile is created.
var tplInitMakefileCreated = assets.TextDesc(
	assets.TextDescKeyWriteInitMakefileCreated,
)

// tplInitMakefileIncludes reports Makefile already includes the directive.
// Arguments: filename.
var tplInitMakefileIncludes = assets.TextDesc(
	assets.TextDescKeyWriteInitMakefileIncludes,
)

// tplInitMakefileAppended reports an include appended to Makefile.
// Arguments: filename.
var tplInitMakefileAppended = assets.TextDesc(
	assets.TextDescKeyWriteInitMakefileAppended,
)

// tplInitPluginSkipped is printed when plugin enablement is skipped.
var tplInitPluginSkipped = assets.TextDesc(
	assets.TextDescKeyWriteInitPluginSkipped,
)

// tplInitPluginAlreadyEnabled is printed when plugin is already enabled.
var tplInitPluginAlreadyEnabled = assets.TextDesc(
	assets.TextDescKeyWriteInitPluginAlreadyEnabled,
)

// tplInitPluginEnabled reports plugin enabled globally.
// Arguments: settings path.
var tplInitPluginEnabled = assets.TextDesc(
	assets.TextDescKeyWriteInitPluginEnabled,
)

// tplInitSkippedDir reports a directory skipped because it exists.
// Arguments: dir.
var tplInitSkippedDir = assets.TextDesc(
	assets.TextDescKeyWriteInitSkippedDir,
)

// tplInitCreatedDir reports a directory created during init.
// Arguments: dir.
var tplInitCreatedDir = assets.TextDesc(
	assets.TextDescKeyWriteInitCreatedDir,
)

// tplInitCreatedWith reports a file created with a qualifier.
// Arguments: path, qualifier.
var tplInitCreatedWith = assets.TextDesc(
	assets.TextDescKeyWriteInitCreatedWith,
)

// tplInitSkippedPlain reports a file skipped without detail.
// Arguments: path.
var tplInitSkippedPlain = assets.TextDesc(
	assets.TextDescKeyWriteInitSkippedPlain,
)

// tplObsidianGenerated reports successful Obsidian vault generation.
// Arguments: entry count, output directory.
var tplObsidianGenerated = assets.TextDesc(
	assets.TextDescKeyWriteObsidianGenerated,
)

// tplObsidianNextSteps is the post-generation guidance.
// Arguments: output directory.
var tplObsidianNextSteps = assets.TextDesc(
	assets.TextDescKeyWriteObsidianNextSteps,
)

// tplJournalOrphanRemoved reports a removed orphan file.
// Arguments: filename.
var tplJournalOrphanRemoved = assets.TextDesc(
	assets.TextDescKeyWriteJournalOrphanRemoved,
)

// tplJournalSiteGenerated reports successful site generation.
// Arguments: entry count, output directory.
var tplJournalSiteGenerated = assets.TextDesc(
	assets.TextDescKeyWriteJournalSiteGenerated,
)

// tplJournalSiteStarting reports the server is starting.
var tplJournalSiteStarting = assets.TextDesc(
	assets.TextDescKeyWriteJournalSiteStarting,
)

// tplJournalSiteBuilding reports a build is in progress.
var tplJournalSiteBuilding = assets.TextDesc(
	assets.TextDescKeyWriteJournalSiteBuilding,
)

// tplJournalSiteNextSteps shows post-generation guidance.
// Arguments: output directory, zensical binary name.
var tplJournalSiteNextSteps = assets.TextDesc(
	assets.TextDescKeyWriteJournalSiteNextSteps,
)

// tplJournalSiteAlt is the alternative command hint.
var tplJournalSiteAlt = assets.TextDesc(
	assets.TextDescKeyWriteJournalSiteAlt,
)

// tplLoopGenerated reports successful loop script generation.
// Arguments: output file path.
var tplLoopGenerated = assets.TextDesc(
	assets.TextDescKeyWriteLoopGenerated,
)

// tplLoopRunCmd shows how to run the generated script.
// Arguments: output file path.
var tplLoopRunCmd = assets.TextDesc(
	assets.TextDescKeyWriteLoopRunCmd,
)

// tplLoopTool shows the selected tool.
// Arguments: tool name.
var tplLoopTool = assets.TextDesc(assets.TextDescKeyWriteLoopTool)

// tplLoopPrompt shows the prompt file.
// Arguments: prompt file path.
var tplLoopPrompt = assets.TextDesc(assets.TextDescKeyWriteLoopPrompt)

// tplLoopMaxIterations shows the max iterations setting.
// Arguments: count.
var tplLoopMaxIterations = assets.TextDesc(
	assets.TextDescKeyWriteLoopMaxIterations,
)

// tplLoopUnlimited shows unlimited iterations.
var tplLoopUnlimited = assets.TextDesc(assets.TextDescKeyWriteLoopUnlimited)

// tplLoopCompletion shows the completion signal.
// Arguments: signal string.
var tplLoopCompletion = assets.TextDesc(assets.TextDescKeyWriteLoopCompletion)

// tplUnpublishNotFound reports no published block was found.
// Arguments: source filename.
var tplUnpublishNotFound = assets.TextDesc(
	assets.TextDescKeyWriteUnpublishNotFound,
)

// tplUnpublishDone reports the published block was removed.
// Arguments: source filename.
var tplUnpublishDone = assets.TextDesc(assets.TextDescKeyWriteUnpublishDone)

// tplPublishHeader reports publishing has started.
var tplPublishHeader = assets.TextDesc(assets.TextDescKeyWritePublishHeader)

// tplPublishSourceFiles lists the source files used for publishing.
var tplPublishSourceFiles = assets.TextDesc(
	assets.TextDescKeyWritePublishSourceFiles,
)

// tplPublishBudget reports the line budget.
// Arguments: budget.
var tplPublishBudget = assets.TextDesc(assets.TextDescKeyWritePublishBudget)

// tplPublishBlock is the heading for the published block detail.
var tplPublishBlock = assets.TextDesc(assets.TextDescKeyWritePublishBlock)

// tplPublishTasks reports pending tasks count.
// Arguments: count.
var tplPublishTasks = assets.TextDesc(assets.TextDescKeyWritePublishTasks)

// tplPublishDecisions reports recent decisions count.
// Arguments: count.
var tplPublishDecisions = assets.TextDesc(
	assets.TextDescKeyWritePublishDecisions,
)

// tplPublishConventions reports key conventions count.
// Arguments: count.
var tplPublishConventions = assets.TextDesc(
	assets.TextDescKeyWritePublishConventions,
)

// tplPublishLearnings reports recent learnings count.
// Arguments: count.
var tplPublishLearnings = assets.TextDesc(
	assets.TextDescKeyWritePublishLearnings,
)

// tplPublishTotal reports the total line count within budget.
// Arguments: total lines, budget.
var tplPublishTotal = assets.TextDesc(assets.TextDescKeyWritePublishTotal)

// tplPublishDryRun reports a publish dry run.
var tplPublishDryRun = assets.TextDesc(assets.TextDescKeyWritePublishDryRun)

// tplPublishDone reports successful publishing with marker info.
var tplPublishDone = assets.TextDesc(assets.TextDescKeyWritePublishDone)

// tplImportNoEntries reports no entries found in MEMORY.md.
var tplImportNoEntries = assets.TextDesc(assets.TextDescKeyWriteImportNoEntries)

// tplImportScanning reports scanning has started.
// Arguments: source filename.
var tplImportScanning = assets.TextDesc(assets.TextDescKeyWriteImportScanning)

// tplImportFound reports the number of entries found.
// Arguments: count.
var tplImportFound = assets.TextDesc(assets.TextDescKeyWriteImportFound)

// tplImportEntry reports an entry being processed.
// Arguments: truncated title (already quoted).
var tplImportEntry = assets.TextDesc(assets.TextDescKeyWriteImportEntry)

// tplImportClassifiedSkip reports an entry classified as skip.
var tplImportClassifiedSkip = assets.TextDesc(
	assets.TextDescKeyWriteImportClassifiedSkip,
)

// tplImportClassified reports an entry classification.
// Arguments: target file, comma-joined keywords.
var tplImportClassified = assets.TextDesc(assets.TextDescKeyWriteImportClassified)

// tplImportAdded reports an entry added to a target file.
// Arguments: target filename.
var tplImportAdded = assets.TextDesc(assets.TextDescKeyWriteImportAdded)

// tplImportSummaryDryRun is the dry-run summary prefix.
// Arguments: count.
var tplImportSummaryDryRun = assets.TextDesc(
	assets.TextDescKeyWriteImportSummaryDryRun,
)

// tplImportSummary is the import summary prefix.
// Arguments: count.
var tplImportSummary = assets.TextDesc(assets.TextDescKeyWriteImportSummary)

// tplImportSkipped reports skipped entries.
// Arguments: count.
var tplImportSkipped = assets.TextDesc(assets.TextDescKeyWriteImportSkipped)

// tplImportDuplicates reports duplicate entries.
// Arguments: count.
var tplImportDuplicates = assets.TextDesc(assets.TextDescKeyWriteImportDuplicates)

// tplMemoryNoChanges reports no changes since last sync.
var tplMemoryNoChanges = assets.TextDesc(assets.TextDescKeyWriteMemoryNoChanges)

// tplMemoryBridgeHeader is the heading for memory status output.
var tplMemoryBridgeHeader = assets.TextDesc(
	assets.TextDescKeyWriteMemoryBridgeHeader,
)

// tplMemorySourceNotActive reports that auto memory is not active.
var tplMemorySourceNotActive = assets.TextDesc(
	assets.TextDescKeyWriteMemorySourceNotActive,
)

// tplMemorySource is a format template for the source path.
// Arguments: path.
var tplMemorySource = assets.TextDesc(assets.TextDescKeyWriteMemorySource)

// tplMemoryMirror is a format template for the mirror relative path.
// Arguments: relative path.
var tplMemoryMirror = assets.TextDesc(assets.TextDescKeyWriteMemoryMirror)

// tplMemoryLastSync is a format template for the last sync time.
// Arguments: formatted time, human-readable duration.
var tplMemoryLastSync = assets.TextDesc(assets.TextDescKeyWriteMemoryLastSync)

// tplMemoryLastSyncNever reports no sync has occurred.
var tplMemoryLastSyncNever = assets.TextDesc(
	assets.TextDescKeyWriteMemoryLastSyncNever,
)

// tplMemorySourceLines is a format template for MEMORY.md line count.
// Arguments: line count.
var tplMemorySourceLines = assets.TextDesc(assets.TextDescKeyWriteMemorySourceLines)

// tplMemorySourceLinesDrift is a format template for MEMORY.md line count
// when drift is detected. Arguments: line count.
var tplMemorySourceLinesDrift = assets.TextDesc(
	assets.TextDescKeyWriteMemorySourceLinesDrift,
)

// tplMemoryMirrorLines is a format template for mirror line count.
// Arguments: line count.
var tplMemoryMirrorLines = assets.TextDesc(
	assets.TextDescKeyWriteMemoryMirrorLines,
)

// tplMemoryMirrorNotSynced reports the mirror has not been synced.
var tplMemoryMirrorNotSynced = assets.TextDesc(
	assets.TextDescKeyWriteMemoryMirrorNotSynced,
)

// tplMemoryDriftDetected reports drift was detected.
var tplMemoryDriftDetected = assets.TextDesc(
	assets.TextDescKeyWriteMemoryDriftDetected,
)

// tplMemoryDriftNone reports no drift.
var tplMemoryDriftNone = assets.TextDesc(assets.TextDescKeyWriteMemoryDriftNone)

// tplMemoryArchives is a format template for archive snapshot count.
// Arguments: count, archive directory name.
var tplMemoryArchives = assets.TextDesc(assets.TextDescKeyWriteMemoryArchives)

// tplPadEntryAdded is a format template for pad entry confirmation.
// Arguments: entry number.
var tplPadEntryAdded = assets.TextDesc(assets.TextDescKeyWritePadEntryAdded)

// tplPadEntryUpdated is a format template for pad entry update confirmation.
// Arguments: entry number.
var tplPadEntryUpdated = assets.TextDesc(assets.TextDescKeyWritePadEntryUpdated)

// tplPadExportPlan is a format template for a dry-run export line.
// Arguments: label, output path.
var tplPadExportPlan = assets.TextDesc(assets.TextDescKeyWritePadExportPlan)

// tplPadExportDone is a format template for a successfully exported blob.
// Arguments: label.
var tplPadExportDone = assets.TextDesc(assets.TextDescKeyWritePadExportDone)

// tplPadExportWriteFailed is a format template for a failed blob write (stderr).
// Arguments: label, error.
var tplPadExportWriteFailed = assets.TextDesc(
	assets.TextDescKeyWritePadExportWriteFailed,
)

// tplPadExportNone is the message when no blob entries exist to export.
var tplPadExportNone = assets.TextDesc(assets.TextDescKeyWritePadExportNone)

// tplPadExportSummary is a format template for the export summary.
// Arguments: verb ("Exported"/"Would export"), count.
var tplPadExportSummary = assets.TextDesc(assets.TextDescKeyWritePadExportSummary)

// tplPadExportVerbDone is the past-tense verb for export summary.
var tplPadExportVerbDone = assets.TextDesc(assets.TextDescKeyWritePadExportVerbDone)

// tplPadExportVerbDryRun is the dry-run verb for export summary.
var tplPadExportVerbDryRun = assets.TextDesc(
	assets.TextDescKeyWritePadExportVerbDryRun,
)

// tplPadImportNone is the message when no entries were found to import.
var tplPadImportNone = assets.TextDesc(assets.TextDescKeyWritePadImportNone)

// tplPadImportDone is a format template for successful line import.
// Arguments: count.
var tplPadImportDone = assets.TextDesc(assets.TextDescKeyWritePadImportDone)

// tplPadImportBlobAdded is a format template for a successfully imported blob.
// Arguments: filename.
var tplPadImportBlobAdded = assets.TextDesc(
	assets.TextDescKeyWritePadImportBlobAdded,
)

// tplPadImportBlobSkipped is a format template for a skipped blob (stderr).
// Arguments: filename, reason.
var tplPadImportBlobSkipped = assets.TextDesc(
	assets.TextDescKeyWritePadImportBlobSkipped,
)

// tplPadImportBlobTooLarge is a format template for a blob exceeding the size limit (stderr).
// Arguments: filename, max bytes.
var tplPadImportBlobTooLarge = assets.TextDesc(
	assets.TextDescKeyWritePadImportBlobTooLarge,
)

// tplPadImportBlobNone is the message when no files were found to import.
var tplPadImportBlobNone = assets.TextDesc(
	assets.TextDescKeyWritePadImportBlobNone,
)

// tplPadImportBlobSummary is a format template for blob import summary.
// Arguments: added count, skipped count.
var tplPadImportBlobSummary = assets.TextDesc(
	assets.TextDescKeyWritePadImportBlobSummary,
)

// tplPadImportCloseWarning is a format template for file close warning (stderr).
// Arguments: filename, error.
var tplPadImportCloseWarning = assets.TextDesc(
	assets.TextDescKeyWritePadImportCloseWarning,
)

// tplPaused is a format template for the pause confirmation.
// Arguments: session ID.
var tplPaused = assets.TextDesc(assets.TextDescKeyWritePaused)

// tplRestoreNoLocal is the message when golden is restored with no local file.
var tplRestoreNoLocal = assets.TextDesc(assets.TextDescKeyWriteRestoreNoLocal)

// tplRestoreMatch is the message when settings already match golden.
var tplRestoreMatch = assets.TextDesc(assets.TextDescKeyWriteRestoreMatch)

// tplRestoreDroppedHeader is a format template for dropped permissions header.
// Arguments: count.
var tplRestoreDroppedHeader = assets.TextDesc(
	assets.TextDescKeyWriteRestoreDroppedHeader,
)

// tplRestoreRestoredHeader is a format template for restored permissions header.
// Arguments: count.
var tplRestoreRestoredHeader = assets.TextDesc(
	assets.TextDescKeyWriteRestoreRestoredHeader,
)

// tplRestoreDenyDroppedHeader is a format template for dropped deny rules header.
// Arguments: count.
var tplRestoreDenyDroppedHeader = assets.TextDesc(
	assets.TextDescKeyWriteRestoreDenyDroppedHeader,
)

// tplRestoreDenyRestoredHeader is a format template for restored deny rules header.
// Arguments: count.
var tplRestoreDenyRestoredHeader = assets.TextDesc(
	assets.TextDescKeyWriteRestoreDenyRestoredHeader,
)

// tplRestoreRemoved is a format template for a removed permission line.
// Arguments: permission string.
var tplRestoreRemoved = assets.TextDesc(assets.TextDescKeyWriteRestoreRemoved)

// tplRestoreAdded is a format template for an added permission line.
// Arguments: permission string.
var tplRestoreAdded = assets.TextDesc(assets.TextDescKeyWriteRestoreAdded)

// tplRestorePermMatch is the message when only non-permission settings differ.
var tplRestorePermMatch = assets.TextDesc(assets.TextDescKeyWriteRestorePermMatch)

// tplRestoreDone is the message after successful restore.
var tplRestoreDone = assets.TextDesc(assets.TextDescKeyWriteRestoreDone)

// tplSnapshotSaved is a format template for golden image save.
// Arguments: golden file path.
var tplSnapshotSaved = assets.TextDesc(assets.TextDescKeyWriteSnapshotSaved)

// tplSnapshotUpdated is a format template for golden image update.
// Arguments: golden file path.
var tplSnapshotUpdated = assets.TextDesc(assets.TextDescKeyWriteSnapshotUpdated)

// tplResumed is a format template for the resume confirmation.
// Arguments: session ID.
var tplResumed = assets.TextDesc(assets.TextDescKeyWriteResumed)

// tplPadEmpty is the message when the scratchpad has no entries.
var tplPadEmpty = assets.TextDesc(assets.TextDescKeyWritePadEmpty)

// tplPadKeyCreated is a format template for key creation notice (stderr).
// Arguments: key file path.
var tplPadKeyCreated = assets.TextDesc(assets.TextDescKeyWritePadKeyCreated)

// tplPadBlobWritten is a format template for blob file write confirmation.
// Arguments: byte count, output path.
var tplPadBlobWritten = assets.TextDesc(assets.TextDescKeyWritePadBlobWritten)

// tplPadEntryRemoved is a format template for pad entry removal confirmation.
// Arguments: entry number.
var tplPadEntryRemoved = assets.TextDesc(assets.TextDescKeyWritePadEntryRemoved)

// tplPadResolveHeader is a format template for a conflict side header.
// Arguments: side label ("OURS"/"THEIRS").
var tplPadResolveHeader = assets.TextDesc(assets.TextDescKeyWritePadResolveHeader)

// tplPadResolveEntry is a format template for a numbered conflict entry.
// Arguments: 1-based index, display string.
var tplPadResolveEntry = assets.TextDesc(assets.TextDescKeyWritePadResolveEntry)

// tplPadEntryMoved is a format template for pad entry move confirmation.
// Arguments: source position, destination position.
var tplPadEntryMoved = assets.TextDesc(assets.TextDescKeyWritePadEntryMoved)

// tplPadMergeDupe is a format template for a duplicate entry during merge.
// Arguments: display string.
var tplPadMergeDupe = assets.TextDesc(assets.TextDescKeyWritePadMergeDupe)

// tplPadMergeAdded is a format template for a newly added entry during merge.
// Arguments: display string, source file.
var tplPadMergeAdded = assets.TextDesc(assets.TextDescKeyWritePadMergeAdded)

// tplPadMergeBlobConflict is a format template for a blob label conflict warning.
// Arguments: label.
var tplPadMergeBlobConflict = assets.TextDesc(
	assets.TextDescKeyWritePadMergeBlobConflict,
)

// tplPadMergeBinaryWarning is a format template for a binary data warning.
// Arguments: filename.
var tplPadMergeBinaryWarning = assets.TextDesc(
	assets.TextDescKeyWritePadMergeBinaryWarning,
)

// tplPadMergeNone is the message when no entries were found to merge.
var tplPadMergeNone = assets.TextDesc(assets.TextDescKeyWritePadMergeNone)

// tplPadMergeNoneNew is a format template when all entries are duplicates.
// Arguments: dupe count, pluralized "duplicate".
var tplPadMergeNoneNew = assets.TextDesc(assets.TextDescKeyWritePadMergeNoneNew)

// tplPadMergeDryRun is a format template for dry-run merge summary.
// Arguments: added count, pluralized "entry", dupe count, pluralized "duplicate".
var tplPadMergeDryRun = assets.TextDesc(assets.TextDescKeyWritePadMergeDryRun)

// tplPadMergeDone is a format template for successful merge summary.
// Arguments: added count, pluralized "entry", dupe count, pluralized "duplicate".
var tplPadMergeDone = assets.TextDesc(assets.TextDescKeyWritePadMergeDone)

// tplSetupPrompt is the interactive prompt for webhook URL entry.
var tplSetupPrompt = assets.TextDesc(assets.TextDescKeyWriteSetupPrompt)

// tplSetupDone is a format template for successful webhook configuration.
// Arguments: masked URL, encrypted file path.
var tplSetupDone = assets.TextDesc(assets.TextDescKeyWriteSetupDone)

// tplTestNoWebhook is the message when no webhook is configured.
var tplTestNoWebhook = assets.TextDesc(assets.TextDescKeyWriteTestNoWebhook)

// tplTestFiltered is the notice when the test event is filtered.
var tplTestFiltered = assets.TextDesc(assets.TextDescKeyWriteTestFiltered)

// tplTestResult is a format template for webhook test response.
// Arguments: HTTP status code, status text.
var tplTestResult = assets.TextDesc(assets.TextDescKeyWriteTestResult)

// tplTestWorking is the success message after a webhook test.
// Arguments: encrypted file path.
var tplTestWorking = assets.TextDesc(assets.TextDescKeyWriteTestWorking)

// tplPromptCreated is the confirmation after creating a prompt template.
// Arguments: prompt name.
var tplPromptCreated = assets.TextDesc(assets.TextDescKeyWritePromptCreated)

// tplPromptNone is printed when no prompts are found.
var tplPromptNone = assets.TextDesc(assets.TextDescKeyWritePromptNone)

// tplPromptItem is a format template for listing a prompt name.
// Arguments: prompt name.
var tplPromptItem = assets.TextDesc(assets.TextDescKeyWritePromptItem)

// tplPromptRemoved is the confirmation after removing a prompt template.
// Arguments: prompt name.
var tplPromptRemoved = assets.TextDesc(assets.TextDescKeyWritePromptRemoved)

// tplReminderAdded is the confirmation for a newly added reminder.
// Arguments: id, message, suffix (e.g. "  (after 2026-03-10)" or "").
var tplReminderAdded = assets.TextDesc(assets.TextDescKeyWriteReminderAdded)

// tplReminderDismissed is the confirmation for a dismissed reminder.
// Arguments: id, message.
var tplReminderDismissed = assets.TextDesc(assets.TextDescKeyWriteReminderDismissed)

// tplReminderNone is printed when there are no reminders.
var tplReminderNone = assets.TextDesc(assets.TextDescKeyWriteReminderNone)

// tplReminderDismissedAll is the summary after dismissing all reminders.
// Arguments: count.
var tplReminderDismissedAll = assets.TextDesc(
	assets.TextDescKeyWriteReminderDismissedAll,
)

// tplReminderItem is a format template for listing a reminder.
// Arguments: id, message, annotation.
var tplReminderItem = assets.TextDesc(assets.TextDescKeyWriteReminderItem)

// tplReminderNotDue is the annotation for reminders not yet due.
// Arguments: date string.
var tplReminderNotDue = assets.TextDesc(assets.TextDescKeyWriteReminderNotDue)

// tplReminderAfterSuffix formats the date-gate suffix for a reminder.
// Arguments: date string.
var tplReminderAfterSuffix = assets.TextDesc(
	assets.TextDescKeyWriteReminderAfterSuffix,
)

// tplLockUnlockEntry is the confirmation for a single locked/unlocked entry.
// Arguments: filename, verb ("locked" or "unlocked").
var tplLockUnlockEntry = assets.TextDesc(assets.TextDescKeyWriteLockUnlockEntry)

// tplLockUnlockNoChanges is printed when all entries already have the target state.
// Arguments: verb.
var tplLockUnlockNoChanges = assets.TextDesc(
	assets.TextDescKeyWriteLockUnlockNoChanges,
)

// tplLockUnlockSummary is the summary after locking/unlocking entries.
// Arguments: capitalized verb, count.
var tplLockUnlockSummary = assets.TextDesc(assets.TextDescKeyWriteLockUnlockSummary)

// tplBackupResult is a format template for a backup result line.
// Arguments: scope, archive path, formatted size.
var tplBackupResult = assets.TextDesc(assets.TextDescKeyWriteBackupResult)

// tplBackupSMBDest is a format template for the SMB destination suffix.
// Arguments: SMB destination path.
var tplBackupSMBDest = assets.TextDesc(assets.TextDescKeyWriteBackupSMBDest)

// tplStatusTitle is the heading for the status output.
var tplStatusTitle = assets.TextDesc(assets.TextDescKeyWriteStatusTitle)

// tplStatusSeparator is the visual separator under the heading.
var tplStatusSeparator = assets.TextDesc(assets.TextDescKeyWriteStatusSeparator)

// tplStatusDir is a format template for the context directory.
// Arguments: context directory path.
var tplStatusDir = assets.TextDesc(assets.TextDescKeyWriteStatusDir)

// tplStatusFiles is a format template for the total file count.
// Arguments: count.
var tplStatusFiles = assets.TextDesc(assets.TextDescKeyWriteStatusFiles)

// tplStatusTokens is a format template for the token estimate.
// Arguments: formatted token count.
var tplStatusTokens = assets.TextDesc(assets.TextDescKeyWriteStatusTokens)

// tplStatusFilesHeader is the heading for the file list section.
var tplStatusFilesHeader = assets.TextDesc(
	assets.TextDescKeyWriteStatusFilesHeader,
)

// tplStatusFileVerbose is a format template for a verbose file entry.
// Arguments: indicator, name, status, formatted tokens, formatted size.
var tplStatusFileVerbose = assets.TextDesc(
	assets.TextDescKeyWriteStatusFileVerbose,
)

// tplStatusFileCompact is a format template for a compact file entry.
// Arguments: indicator, name, status.
var tplStatusFileCompact = assets.TextDesc(
	assets.TextDescKeyWriteStatusFileCompact,
)

// tplStatusPreviewLine is a format template for a content preview line.
// Arguments: line text.
var tplStatusPreviewLine = assets.TextDesc(
	assets.TextDescKeyWriteStatusPreviewLine,
)

// tplStatusActivityHeader is the heading for the recent activity section.
var tplStatusActivityHeader = assets.TextDesc(
	assets.TextDescKeyWriteStatusActivityHeader,
)

// tplStatusActivityItem is a format template for a recent activity entry.
// Arguments: filename, relative time string.
var tplStatusActivityItem = assets.TextDesc(
	assets.TextDescKeyWriteStatusActivityItem,
)

// tplTimeJustNow is the display string for "just now" relative time.
var tplTimeJustNow = assets.TextDesc(assets.TextDescKeyWriteTimeJustNow)

// tplTimeMinuteAgo is the display string for "1 minute ago".
var tplTimeMinuteAgo = assets.TextDesc(assets.TextDescKeyWriteTimeMinuteAgo)

// tplTimeMinutesAgo is a format template for minutes ago.
// Arguments: count.
var tplTimeMinutesAgo = assets.TextDesc(assets.TextDescKeyWriteTimeMinutesAgo)

// tplTimeHourAgo is the display string for "1 hour ago".
var tplTimeHourAgo = assets.TextDesc(assets.TextDescKeyWriteTimeHourAgo)

// tplTimeHoursAgo is a format template for hours ago.
// Arguments: count.
var tplTimeHoursAgo = assets.TextDesc(assets.TextDescKeyWriteTimeHoursAgo)

// tplTimeDayAgo is the display string for "1 day ago".
var tplTimeDayAgo = assets.TextDesc(assets.TextDescKeyWriteTimeDayAgo)

// tplTimeDaysAgo is a format template for days ago.
// Arguments: count.
var tplTimeDaysAgo = assets.TextDesc(assets.TextDescKeyWriteTimeDaysAgo)

// TplTimeOlderFormat is the Go time layout for dates older than a week.
// Exported because callers must format the fallback date before calling FormatTimeAgo.
// Deprecated: Use config.TimeOlderFormat instead.
const TplTimeOlderFormat = config.TimeOlderFormat

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

// tplJournalSyncNone is printed when no journal entries are found.
var tplJournalSyncNone = assets.TextDesc(assets.TextDescKeyWriteJournalSyncNone)

// tplJournalSyncLocked is a format template for a newly locked entry.
// Arguments: filename.
var tplJournalSyncLocked = assets.TextDesc(
	assets.TextDescKeyWriteJournalSyncLocked,
)

// tplJournalSyncUnlocked is a format template for a newly unlocked entry.
// Arguments: filename.
var tplJournalSyncUnlocked = assets.TextDesc(
	assets.TextDescKeyWriteJournalSyncUnlocked,
)

// tplJournalSyncMatch is printed when state already matches frontmatter.
var tplJournalSyncMatch = assets.TextDesc(
	assets.TextDescKeyWriteJournalSyncMatch,
)

// tplJournalSyncLockedCount is a format template for locked entry count.
// Arguments: count.
var tplJournalSyncLockedCount = assets.TextDesc(
	assets.TextDescKeyWriteJournalSyncLockedCount,
)

// tplJournalSyncUnlockedCount is a format template for unlocked entry count.
// Arguments: count.
var tplJournalSyncUnlockedCount = assets.TextDesc(
	assets.TextDescKeyWriteJournalSyncUnlockedCount,
)
