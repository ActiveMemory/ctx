//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for sync operations write output.
const (
	// DescKeyWriteSynced is the text key for write synced messages.
	DescKeyWriteSynced = "write.synced"
	// DescKeyWriteSyncAction is the text key for write sync action messages.
	DescKeyWriteSyncAction = "write.sync-action"
	// DescKeyWriteSyncDryRun is the text key for write sync dry run messages.
	DescKeyWriteSyncDryRun = "write.sync-dry-run"
	// DescKeyWriteSyncDryRunSummary is the text key for write sync dry run
	// summary messages.
	DescKeyWriteSyncDryRunSummary = "write.sync-dry-run-summary"
	// DescKeyWriteSyncHeader is the text key for write sync header messages.
	DescKeyWriteSyncHeader = "write.sync-header"
	// DescKeyWriteSyncInSync is the text key for write sync in sync messages.
	DescKeyWriteSyncInSync = "write.sync-in-sync"
	// DescKeyWriteSyncSeparator is the text key for write sync separator messages.
	DescKeyWriteSyncSeparator = "write.sync-separator"
	// DescKeyWriteSyncSuggestion is the text key for write sync suggestion
	// messages.
	DescKeyWriteSyncSuggestion = "write.sync-suggestion"
	// DescKeyWriteSyncSummary is the text key for write sync summary messages.
	DescKeyWriteSyncSummary = "write.sync-summary"
)

// DescKeys for sync topic names.
const (
	// DescKeySyncTopicEslint is the text key for sync topic eslint messages.
	DescKeySyncTopicEslint = "sync.topic.eslint"
	// DescKeySyncTopicPrettier is the text key for sync topic prettier messages.
	DescKeySyncTopicPrettier = "sync.topic.prettier"
	// DescKeySyncTopicTSConfig is the text key for sync topic ts config messages.
	DescKeySyncTopicTSConfig = "sync.topic.tsconfig"
	// DescKeySyncTopicEditorConfig is the text key for sync topic editor config
	// messages.
	DescKeySyncTopicEditorConfig = "sync.topic.editorconfig"
	// DescKeySyncTopicMakefile is the text key for sync topic makefile messages.
	DescKeySyncTopicMakefile = "sync.topic.makefile"
	// DescKeySyncTopicDockerfile is the text key for sync topic dockerfile
	// messages.
	DescKeySyncTopicDockerfile = "sync.topic.dockerfile"
)

// DescKeys for sync rule descriptions.
const (
	// DescKeySyncDepsDescription is the text key for sync deps description
	// messages.
	DescKeySyncDepsDescription = "sync.deps.description"
	// DescKeySyncDepsSuggestion is the text key for sync deps suggestion messages.
	DescKeySyncDepsSuggestion = "sync.deps.suggestion"
	// DescKeySyncConfigDescription is the text key for sync config description
	// messages.
	DescKeySyncConfigDescription = "sync.config.description"
	// DescKeySyncConfigSuggestion is the text key for sync config suggestion
	// messages.
	DescKeySyncConfigSuggestion = "sync.config.suggestion"
	// DescKeySyncDirDescription is the text key for sync dir description messages.
	DescKeySyncDirDescription = "sync.dir.description"
	// DescKeySyncDirSuggestion is the text key for sync dir suggestion messages.
	DescKeySyncDirSuggestion = "sync.dir.suggestion"
)
