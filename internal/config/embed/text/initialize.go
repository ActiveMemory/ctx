//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for init setup and abort messages.
const (
	DescKeyWriteInitAborted           = "write.init-aborted"
	DescKeyWriteInitBackup            = "write.init-backup"
	DescKeyWriteInitCreatedDir        = "write.init-created-dir"
	DescKeyWriteInitCreatingRootFiles = "write.init-creating-root-files"
	DescKeyWriteInitCtxContentExists  = "write.init-ctx-content-exists"
	DescKeyWriteInitExistsSkipped     = "write.init-exists-skipped"
)

// DescKeys for init file creation output.
const (
	DescKeyWriteInitFileCreated     = "write.init-file-created"
	DescKeyWriteInitFileExistsNoCtx = "write.init-file-exists-no-ctx"
)

// DescKeys for init gitignore output.
const (
	DescKeyWriteInitGitignoreReview  = "write.init-gitignore-review"
	DescKeyWriteInitGitignoreUpdated = "write.init-gitignore-updated"
)

// DescKeys for init Makefile output.
const (
	DescKeyWriteInitMakefileAppended = "write.init-makefile-appended"
	DescKeyWriteInitMakefileCreated  = "write.init-makefile-created"
	DescKeyWriteInitMakefileIncludes = "write.init-makefile-includes"
)

// DescKeys for init merge and prompt output.
const (
	DescKeyWriteInitMerged          = "write.init-merged"
	DescKeyWriteInitNextStepsBlock  = "write.init-next-steps-block"
	DescKeyWriteInitWorkflowTips    = "write.init-workflow-tips"
	DescKeyWriteInitNoChanges       = "write.init-no-changes"
	DescKeyWriteInitOverwritePrompt = "write.init-overwrite-prompt"
)

// DescKeys for init permission setup output.
const (
	DescKeyWriteInitPermsAllow         = "write.init-perms-allow"
	DescKeyWriteInitPermsAllowDeny     = "write.init-perms-allow-deny"
	DescKeyWriteInitPermsDeduped       = "write.init-perms-deduped"
	DescKeyWriteInitPermsDeny          = "write.init-perms-deny"
	DescKeyWriteInitPermsMergedDeduped = "write.init-perms-merged-deduped"
)

// DescKeys for init plugin enablement output.
const (
	DescKeyWriteInitPluginAlreadyEnabled = "write.init-plugin-already-enabled"
	DescKeyWriteInitPluginEnabled        = "write.init-plugin-enabled"
	DescKeyWriteInitPluginSkipped        = "write.init-plugin-skipped"
)

// DescKeys for init scratchpad setup output.
const (
	DescKeyWriteInitScratchpadKeyCreated = "write.init-scratchpad-key-created"
	DescKeyWriteInitScratchpadNoKey      = "write.init-scratchpad-no-key"
	DescKeyWriteInitScratchpadPlaintext  = "write.init-scratchpad-plaintext"
)

// DescKeys for init skip and directory output.
const (
	DescKeyWriteInitSkippedDir   = "write.init-skipped-dir"
	DescKeyWriteInitSkippedPlain = "write.init-skipped-plain"
)

// DescKeys for init section update output.
const (
	DescKeyWriteInitUpdatedCtxSection = "write.init-updated-ctx-section"
)

// DescKeys for init completion output.
const (
	DescKeyWriteInitGettingStartedSaved  = "write.init-getting-started-saved"
	DescKeyWriteInitSettingUpPermissions = "write.init-setting-up-permissions"
	DescKeyWriteInitWarnNonFatal         = "write.init-warn-non-fatal"
	DescKeyWriteInitialized              = "write.initialized"
)

// Init component labels for InfoWarnNonFatal diagnostic output.
const (
	DescKeyInitLabelEntryTemplates = "init.label-entry-templates"
	DescKeyInitLabelScratchpad     = "init.label-scratchpad"
	DescKeyInitLabelProjectDirs    = "init.label-project-dirs"
	DescKeyInitLabelPermissions    = "init.label-permissions"
	DescKeyInitLabelPluginEnable   = "init.label-plugin-enable"
)

// Init confirmation prompts and mode labels.
const (
	DescKeyInitConfirmClaude = "init.confirm-claude"
)
