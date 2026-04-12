//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for steering write output.
const (
	// DescKeyWriteSteeringCreated is the text key for write steering created
	// messages.
	DescKeyWriteSteeringCreated = "write.steering-created"
	// DescKeyWriteSteeringSkipped is the text key for write steering skipped
	// messages.
	DescKeyWriteSteeringSkipped = "write.steering-skipped"
	// DescKeyWriteSteeringInitSummary is the text key for write steering init
	// summary messages.
	DescKeyWriteSteeringInitSummary = "write.steering-init-summary"
	// DescKeyWriteSteeringFileEntry is the text key for write steering file entry
	// messages.
	DescKeyWriteSteeringFileEntry = "write.steering-file-entry"
	// DescKeyWriteSteeringFileCount is the text key for write steering file count
	// messages.
	DescKeyWriteSteeringFileCount = "write.steering-file-count"
	// DescKeyWriteSteeringPreviewHead is the text key for write steering preview
	// head messages.
	DescKeyWriteSteeringPreviewHead = "write.steering-preview-head"
	// DescKeyWriteSteeringPreviewEntry is the text key for write steering preview
	// entry messages.
	DescKeyWriteSteeringPreviewEntry = "write.steering-preview-entry"
	// DescKeyWriteSteeringPreviewCount is the text key for write steering preview
	// count messages.
	DescKeyWriteSteeringPreviewCount = "write.steering-preview-count"
	// DescKeyWriteSteeringSyncWritten is the text key for write steering sync
	// written messages.
	DescKeyWriteSteeringSyncWritten = "write.steering-sync-written"
	// DescKeyWriteSteeringSyncSkipped is the text key for write steering sync
	// skipped messages.
	DescKeyWriteSteeringSyncSkipped = "write.steering-sync-skipped"
	// DescKeyWriteSteeringSyncError is the text key for write steering sync error
	// messages.
	DescKeyWriteSteeringSyncError = "write.steering-sync-error"
	// DescKeyWriteSteeringSyncSummary is the text key for write steering sync
	// summary messages.
	DescKeyWriteSteeringSyncSummary = "write.steering-sync-summary"
	// DescKeyWriteSteeringNoFiles is the message when no steering
	// files exist.
	DescKeyWriteSteeringNoFiles = "write.steering-no-files"
	// DescKeyWriteSteeringNoMatch is the message when no steering
	// files match the prompt.
	DescKeyWriteSteeringNoMatch = "write.steering-no-match"
	// DescKeyWriteSteeringAddModeHint is the one-line hint
	// printed after `ctx steering add` creates a new file,
	// explaining the default inclusion mode and how to
	// switch it.
	DescKeyWriteSteeringAddModeHint = "write.steering-add-mode-hint"
)

// DescKeys for steering foundation file scaffolding text.
// These replace the constants formerly in tpl/tpl_steering.go
// (which had zero format placeholders and didn't belong in tpl/).
const (
	// DescKeyWriteSteeringDescProduct is the text key for the product
	// steering file description.
	DescKeyWriteSteeringDescProduct = "write.steering-desc-product"
	// DescKeyWriteSteeringDescTech is the text key for the tech
	// steering file description.
	DescKeyWriteSteeringDescTech = "write.steering-desc-tech"
	// DescKeyWriteSteeringDescStructure is the text key for the structure
	// steering file description.
	DescKeyWriteSteeringDescStructure = "write.steering-desc-structure"
	// DescKeyWriteSteeringDescWorkflow is the text key for the workflow
	// steering file description.
	DescKeyWriteSteeringDescWorkflow = "write.steering-desc-workflow"
	// DescKeyWriteSteeringGuidance is the text key for the HTML guidance
	// comment prepended to every scaffolded steering body.
	DescKeyWriteSteeringGuidance = "write.steering-guidance"
	// DescKeyWriteSteeringBodyProduct is the text key for the product
	// steering file body template.
	DescKeyWriteSteeringBodyProduct = "write.steering-body-product"
	// DescKeyWriteSteeringBodyTech is the text key for the tech steering
	// file body template.
	DescKeyWriteSteeringBodyTech = "write.steering-body-tech"
	// DescKeyWriteSteeringBodyStructure is the text key for the structure
	// steering file body template.
	DescKeyWriteSteeringBodyStructure = "write.steering-body-structure"
	// DescKeyWriteSteeringBodyWorkflow is the text key for the workflow
	// steering file body template.
	DescKeyWriteSteeringBodyWorkflow = "write.steering-body-workflow"
)
