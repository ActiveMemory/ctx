//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for init setup and abort messages.
const (
	// DescKeyWriteInitAborted is the text key for write init aborted messages.
	DescKeyWriteInitAborted = "write.init-aborted"
	// DescKeyWriteInitBackup is the text key for write init backup messages.
	DescKeyWriteInitBackup = "write.init-backup"
	// DescKeyWriteInitCreatedDir is the text key for write init created dir
	// messages.
	DescKeyWriteInitCreatedDir = "write.init-created-dir"
	// DescKeyWriteInitCreatingRootFiles is the text key for write init creating
	// root files messages.
	DescKeyWriteInitCreatingRootFiles = "write.init-creating-root-files"
	// DescKeyWriteInitCtxContentExists is the text key for write init ctx content
	// exists messages.
	DescKeyWriteInitCtxContentExists = "write.init-ctx-content-exists"
	// DescKeyWriteInitExistsSkipped is the text key for write init exists skipped
	// messages.
	DescKeyWriteInitExistsSkipped = "write.init-exists-skipped"
)

// DescKeys for init file creation output.
const (
	// DescKeyWriteInitFileCreated is the text key for write init file created
	// messages.
	DescKeyWriteInitFileCreated = "write.init-file-created"
	// DescKeyWriteInitFileExistsNoCtx is the text key for write init file exists
	// no ctx messages.
	DescKeyWriteInitFileExistsNoCtx = "write.init-file-exists-no-ctx"
)

// DescKeys for init gitignore output.
const (
	// DescKeyWriteInitGitignoreReview is the text key for write init gitignore
	// review messages.
	DescKeyWriteInitGitignoreReview = "write.init-gitignore-review"
	// DescKeyWriteInitGitignoreUpdated is the text key for write init gitignore
	// updated messages.
	DescKeyWriteInitGitignoreUpdated = "write.init-gitignore-updated"
)

// DescKeys for init Makefile output.
const (
	// DescKeyWriteInitMakefileAppended is the text key for write init makefile
	// appended messages.
	DescKeyWriteInitMakefileAppended = "write.init-makefile-appended"
	// DescKeyWriteInitMakefileCreated is the text key for write init makefile
	// created messages.
	DescKeyWriteInitMakefileCreated = "write.init-makefile-created"
	// DescKeyWriteInitMakefileIncludes is the text key for write init makefile
	// includes messages.
	DescKeyWriteInitMakefileIncludes = "write.init-makefile-includes"
)

// DescKeys for init merge and prompt output.
const (
	// DescKeyWriteInitMerged is the text key for write init merged messages.
	DescKeyWriteInitMerged = "write.init-merged"
	// DescKeyWriteInitNextStepsBlock is the text key for write init next steps
	// block messages.
	DescKeyWriteInitNextStepsBlock = "write.init-next-steps-block"
	// DescKeyWriteInitWorkflowTips is the text key for write init workflow tips
	// messages.
	DescKeyWriteInitWorkflowTips = "write.init-workflow-tips"
	// DescKeyWriteInitNoChanges is the text key for write init no changes
	// messages.
	DescKeyWriteInitNoChanges = "write.init-no-changes"
	// DescKeyWriteInitOverwritePrompt is the text key for write init overwrite
	// prompt messages.
	DescKeyWriteInitOverwritePrompt = "write.init-overwrite-prompt"
)

// DescKeys for init permission setup output.
const (
	// DescKeyWriteInitPermsAllow is the text key for write init perms allow
	// messages.
	DescKeyWriteInitPermsAllow = "write.init-perms-allow"
	// DescKeyWriteInitPermsAllowDeny is the text key for write init perms allow
	// deny messages.
	DescKeyWriteInitPermsAllowDeny = "write.init-perms-allow-deny"
	// DescKeyWriteInitPermsDeduped is the text key for write init perms deduped
	// messages.
	DescKeyWriteInitPermsDeduped = "write.init-perms-deduped"
	// DescKeyWriteInitPermsDeny is the text key for write init perms deny
	// messages.
	DescKeyWriteInitPermsDeny = "write.init-perms-deny"
	// DescKeyWriteInitPermsMergedDeduped is the text key for write init perms
	// merged deduped messages.
	DescKeyWriteInitPermsMergedDeduped = "write.init-perms-merged-deduped"
)

// DescKeys for init plugin enablement output.
const (
	// DescKeyWriteInitPluginAlreadyEnabled is the text key for write init plugin
	// already enabled messages.
	DescKeyWriteInitPluginAlreadyEnabled = "write.init-plugin-already-enabled"
	// DescKeyWriteInitPluginEnabled is the text key for write init plugin enabled
	// messages.
	DescKeyWriteInitPluginEnabled = "write.init-plugin-enabled"
	// DescKeyWriteInitPluginSkipped is the text key for write init plugin skipped
	// messages.
	DescKeyWriteInitPluginSkipped = "write.init-plugin-skipped"
	// DescKeyWriteInitPluginLocalAlreadyEnabled is the text key for write init
	// plugin already enabled locally messages.
	DescKeyWriteInitPluginLocalAlreadyEnabled = "write.init-plugin-local-already-enabled"
	// DescKeyWriteInitPluginLocalEnabled is the text key for write init plugin
	// enabled locally messages.
	DescKeyWriteInitPluginLocalEnabled = "write.init-plugin-local-enabled"
)

// DescKeys for init scratchpad setup output.
const (
	// DescKeyWriteInitScratchpadKeyCreated is the text key for write init
	// scratchpad key created messages.
	DescKeyWriteInitScratchpadKeyCreated = "write.init-scratchpad-key-created"
	// DescKeyWriteInitScratchpadNoKey is the text key for write init scratchpad
	// no key messages.
	DescKeyWriteInitScratchpadNoKey = "write.init-scratchpad-no-key"
	// DescKeyWriteInitScratchpadPlaintext is the text key for write init
	// scratchpad plaintext messages.
	DescKeyWriteInitScratchpadPlaintext = "write.init-scratchpad-plaintext"
)

// DescKeys for init skip and directory output.
const (
	// DescKeyWriteInitSkippedDir is the text key for write init skipped dir
	// messages.
	DescKeyWriteInitSkippedDir = "write.init-skipped-dir"
	// DescKeyWriteInitSkippedPlain is the text key for write init skipped plain
	// messages.
	DescKeyWriteInitSkippedPlain = "write.init-skipped-plain"
)

// DescKeys for init section update output.
const (
	// DescKeyWriteInitUpdatedCtxSection is the text key for write init updated
	// ctx section messages.
	DescKeyWriteInitUpdatedCtxSection = "write.init-updated-ctx-section"
)

// DescKeys for init completion output.
const (
	// DescKeyWriteInitGettingStartedSaved is the text key for write init getting
	// started saved messages.
	DescKeyWriteInitGettingStartedSaved = "write.init-getting-started-saved"
	// DescKeyWriteInitSettingUpPermissions is the text key for write init setting
	// up permissions messages.
	DescKeyWriteInitSettingUpPermissions = "write.init-setting-up-permissions"
	// DescKeyWriteInitWarnNonFatal is the text key for write init warn non fatal
	// messages.
	DescKeyWriteInitWarnNonFatal = "write.init-warn-non-fatal"
	// DescKeyWriteInitialized is the text key for write initialized messages.
	DescKeyWriteInitialized = "write.initialized"
)

// Init component labels for InfoWarnNonFatal diagnostic output.
const (
	// DescKeyInitLabelEntryTemplates is the text key for init label entry
	// templates messages.
	DescKeyInitLabelEntryTemplates = "init.label-entry-templates"
	// DescKeyInitLabelScratchpad is the text key for init label scratchpad
	// messages.
	DescKeyInitLabelScratchpad = "init.label-scratchpad"
	// DescKeyInitLabelProjectDirs is the text key for init label project dirs
	// messages.
	DescKeyInitLabelProjectDirs = "init.label-project-dirs"
	// DescKeyInitLabelPermissions is the text key for init label permissions
	// messages.
	DescKeyInitLabelPermissions = "init.label-permissions"
	// DescKeyInitLabelPluginEnable is the text key for init label plugin enable
	// messages.
	DescKeyInitLabelPluginEnable = "init.label-plugin-enable"
)

// Init confirmation prompts and mode labels.
const (
	// DescKeyInitConfirmClaude is the text key for init confirm claude messages.
	DescKeyInitConfirmClaude = "init.confirm-claude"
)
