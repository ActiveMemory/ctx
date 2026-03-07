//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

// prefixError is prepended to all error messages written to stderr.
const prefixError = "Error: "

// tplPathExists is a format template for reporting that a destination path
// already exists. Arguments: original path, resolved destination path.
const tplPathExists = "  %s -> %s (exists)"

// tplExistsWritingAsAlternative is a format template for reporting that a
// file exists and content was written to an alternative filename instead.
// Arguments: original path, alternative path.
const tplExistsWritingAsAlternative = "  ! %s exists, writing as %s"

// tplDryRun is printed when a command runs in dry-run mode.
const tplDryRun = "Dry run — no files will be written."

// tplSource is a format template for reporting a source path.
// Arguments: path.
const tplSource = "  Source: %s"

// tplMirror is a format template for reporting a mirror path.
// Arguments: relative mirror path.
const tplMirror = "  Mirror: %s"

// tplStatusDrift is printed when drift is detected.
const tplStatusDrift = "  Status: drift detected (source is newer)"

// tplStatusNoDrift is printed when no drift is detected.
const tplStatusNoDrift = "  Status: no drift"

// tplArchived is a format template for reporting an archived file.
// Arguments: archive filename.
const tplArchived = "Archived previous mirror to %s"

// tplSynced is a format template for reporting a successful sync.
// Arguments: source label, destination relative path.
const tplSynced = "Synced %s -> %s"

// tplLines is a format template for reporting line counts.
// Arguments: line count.
const tplLines = "  Lines: %d"

// tplLinesPrevious is a format template appended to line counts when a
// previous count is available. Arguments: previous line count.
const tplLinesPrevious = " (was %d)"

// tplNewContent is a format template for reporting new content since last sync.
// Arguments: line count.
const tplNewContent = "  New content: %d lines since last sync"

// tplAddedTo is a format template for confirming an entry was added.
// Arguments: filename.
const tplAddedTo = "✓ Added to %s"

// tplMovingTask is a format template for a task being moved to completed.
// Arguments: truncated task text.
const tplMovingTask = "✓ Moving completed task: %s"

// tplSkippingTask is a format template for a task skipped due to
// incomplete children. Arguments: truncated task text.
const tplSkippingTask = "! Skipping (has incomplete children): %s"

// tplArchivedTasks is a format template for archived tasks summary.
// Arguments: count, archive file path, days threshold.
const tplArchivedTasks = "✓ Archived %d tasks to %s (older than %d days)"

// tplCompletedTask is a format template for a task marked complete.
// Arguments: task text.
const tplCompletedTask = "✓ Completed: %s"

// tplConfigProfileDev is the status output for dev profile.
const tplConfigProfileDev = "active: dev (verbose logging enabled)"

// tplConfigProfileBase is the status output for base profile.
const tplConfigProfileBase = "active: base (defaults)"

// tplConfigProfileNone is the status output when no profile exists.
// Arguments: ctxrc filename.
const tplConfigProfileNone = "active: none (%s does not exist)"

// tplDepsNoProject is printed when no supported project is detected.
const tplDepsNoProject = "No supported project detected."

// tplDepsLookingFor is printed with the list of files checked.
const tplDepsLookingFor = "Looking for: go.mod, package.json, requirements.txt, pyproject.toml, Cargo.toml"

// tplDepsUseType hints at the --type flag.
// Arguments: comma-separated list of builder names.
const tplDepsUseType = "Use --type to force: %s"

// tplDepsNoDeps is printed when no dependencies are found.
const tplDepsNoDeps = "No dependencies found."

// tplSkillsHeader is the heading for the skills list.
const tplSkillsHeader = "Available Skills:"

// tplSkillLine formats a single skill entry.
// Arguments: name, description.
const tplSkillLine = "  /%-22s %s"

// tplHookCopilotSkipped reports that copilot instructions were skipped.
// Arguments: target file path.
const tplHookCopilotSkipped = "  ○ %s (ctx content exists, skipped)"

// tplHookCopilotForceHint tells the user about the --force flag.
const tplHookCopilotForceHint = "  Use --force to overwrite (not yet implemented)."

// tplHookCopilotMerged reports that copilot instructions were merged.
// Arguments: target file path.
const tplHookCopilotMerged = "  ✓ %s (merged)"

// tplHookCopilotCreated reports that copilot instructions were created.
// Arguments: target file path.
const tplHookCopilotCreated = "  ✓ %s"

// tplHookCopilotSessionsDir reports that the sessions directory was created.
// Arguments: sessions directory path.
const tplHookCopilotSessionsDir = "  ✓ %s/"

// tplHookCopilotSummary is the post-write summary for copilot.
const tplHookCopilotSummary = `Copilot Chat (agent mode) will now:
  1. Read .context/ files at session start
  2. Save session summaries to .context/sessions/
  3. Proactively update context during work`

// tplHookUnknownTool reports an unrecognized tool name.
// Arguments: tool name.
const tplHookUnknownTool = "Unknown tool: %s\n"

// tplInitOverwritePrompt prompts the user before overwriting .context/.
// Arguments: context directory path.
const tplInitOverwritePrompt = "%s already exists. Overwrite? [y/N] "

// tplInitAborted is printed when the user declines overwriting.
const tplInitAborted = "Aborted."

// tplInitExistsSkipped reports a file that was skipped because it exists.
// Arguments: filename.
const tplInitExistsSkipped = "  ○ %s (exists, skipped)"

// tplInitFileCreated reports a file that was successfully created.
// Arguments: filename.
const tplInitFileCreated = "  ✓ %s"

// tplInitialized reports successful context initialization.
// Arguments: context directory path.
const tplInitialized = "Context initialized in %s/"

// tplInitWarnNonFatal reports a non-fatal warning during init.
// Arguments: label, error.
const tplInitWarnNonFatal = "  ⚠ %s: %v"

// tplInitScratchpadPlaintext reports a plaintext scratchpad was created.
// Arguments: path.
const tplInitScratchpadPlaintext = "  ✓ %s (plaintext scratchpad)"

// tplInitScratchpadNoKey warns about a missing key for an encrypted scratchpad.
// Arguments: key path.
const tplInitScratchpadNoKey = "  ⚠ Encrypted scratchpad found but no key at %s"

// tplInitScratchpadKeyCreated reports a scratchpad key was generated.
// Arguments: key path.
const tplInitScratchpadKeyCreated = "  ✓ Scratchpad key created at %s"

// tplInitCreatingRootFiles is the heading before project root file creation.
const tplInitCreatingRootFiles = "Creating project root files..."

// tplInitSettingUpPermissions is the heading before permissions setup.
const tplInitSettingUpPermissions = "Setting up Claude Code permissions..."

// tplInitGitignoreUpdated reports .gitignore entries were added.
// Arguments: count of entries added.
const tplInitGitignoreUpdated = "  ✓ .gitignore updated (%d entries added)"

// tplInitGitignoreReview hints how to review the .gitignore changes.
const tplInitGitignoreReview = "  Review with: cat .gitignore"

// tplInitNextSteps is the next-steps guidance block after init completes.
const tplInitNextSteps = `Next steps:
  1. Edit .context/TASKS.md to add your current tasks
  2. Run 'ctx status' to see context summary
  3. Run 'ctx agent' to get AI-ready context packet`

// tplInitPluginInfo is the plugin installation guidance block.
const tplInitPluginInfo = `Claude Code users: install the ctx plugin for hooks & skills:
  /plugin marketplace add ActiveMemory/ctx
  /plugin install ctx@activememory-ctx`

// tplInitPluginNote is the note about local plugin enabling.
const tplInitPluginNote = `Note: local plugin installs are not auto-enabled globally.
Run 'ctx init' again after installing the plugin to enable it,
or manually add to ~/.claude/settings.json:
  {"enabledPlugins": {"ctx@activememory-ctx": true}}`

// tplInitCtxContentExists reports a file skipped because ctx content exists.
// Arguments: path.
const tplInitCtxContentExists = "  ○ %s (ctx content exists, skipped)"

// tplInitMerged reports a file merged during init.
// Arguments: path.
const tplInitMerged = "  ✓ %s (merged)"

// tplInitBackup reports a backup file created.
// Arguments: backup path.
const tplInitBackup = "  ✓ %s (backup)"

// tplInitUpdatedCtxSection reports a file whose ctx section was updated.
// Arguments: path.
const tplInitUpdatedCtxSection = "  ✓ %s (updated ctx section)"

// tplInitUpdatedPlanSection reports a file whose plan section was updated.
// Arguments: path.
const tplInitUpdatedPlanSection = "  ✓ %s (updated plan section)"

// tplInitUpdatedPromptSection reports a file whose prompt section was updated.
// Arguments: path.
const tplInitUpdatedPromptSection = "  ✓ %s (updated prompt section)"

// tplInitFileExistsNoCtx reports a file exists without ctx content.
// Arguments: path.
const tplInitFileExistsNoCtx = "%s exists but has no ctx content."

// tplInitNoChanges reports a settings file with no changes needed.
// Arguments: path.
const tplInitNoChanges = "  ○ %s (no changes needed)"

// tplInitPermsMergedDeduped reports permissions merged and deduped.
// Arguments: path.
const tplInitPermsMergedDeduped = "  ✓ %s (added ctx permissions, removed duplicates)"

// tplInitPermsDeduped reports duplicate permissions removed.
// Arguments: path.
const tplInitPermsDeduped = "  ✓ %s (removed duplicate permissions)"

// tplInitPermsAllowDeny reports allow+deny permissions added.
// Arguments: path.
const tplInitPermsAllowDeny = "  ✓ %s (added ctx allow + deny permissions)"

// tplInitPermsDeny reports deny permissions added.
// Arguments: path.
const tplInitPermsDeny = "  ✓ %s (added ctx deny permissions)"

// tplInitPermsAllow reports ctx permissions added.
// Arguments: path.
const tplInitPermsAllow = "  ✓ %s (added ctx permissions)"

// tplInitMakefileCreated is printed when a new Makefile is created.
const tplInitMakefileCreated = "  ✓ Makefile (created with ctx include)"

// tplInitMakefileIncludes reports Makefile already includes the directive.
// Arguments: filename.
const tplInitMakefileIncludes = "  ○ Makefile (already includes %s)"

// tplInitMakefileAppended reports an include appended to Makefile.
// Arguments: filename.
const tplInitMakefileAppended = "  ✓ Makefile (appended %s include)"

// tplInitPluginSkipped is printed when plugin enablement is skipped.
const tplInitPluginSkipped = "  ○ Plugin enablement skipped (plugin not installed)"

// tplInitPluginAlreadyEnabled is printed when plugin is already enabled.
const tplInitPluginAlreadyEnabled = "  ○ Plugin already enabled globally"

// tplInitPluginEnabled reports plugin enabled globally.
// Arguments: settings path.
const tplInitPluginEnabled = "  ✓ Plugin enabled globally in %s"

// tplInitSkippedDir reports a directory skipped because it exists.
// Arguments: dir.
const tplInitSkippedDir = "  ○ %s/ (exists, skipped)"

// tplInitCreatedDir reports a directory created during init.
// Arguments: dir.
const tplInitCreatedDir = "  ✓ %s/"

// tplInitCreatedWith reports a file created with a qualifier.
// Arguments: path, qualifier.
const tplInitCreatedWith = "  ✓ %s%s"

// tplInitSkippedPlain reports a file skipped without detail.
// Arguments: path.
const tplInitSkippedPlain = "  ○ %s (skipped)"
