//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// VS Code artifact output formatting keys.
const (
	// DescKeyWriteVscodeCreated reports a VS Code config file was created.
	DescKeyWriteVscodeCreated = "write.vscode-created"
	// DescKeyWriteVscodeExistsSkipped reports a file was skipped (exists).
	DescKeyWriteVscodeExistsSkipped = "write.vscode-exists-skipped"
	// DescKeyWriteVscodeRecommendationExists reports the extension
	// recommendation already exists.
	// DescKeyWriteVscodeRecommendationExists is the text key for write vscode
	// recommendation exists messages.
	DescKeyWriteVscodeRecommendationExists = "write.vscode-recommendation-exists"
	// DescKeyWriteVscodeAddManually reports the file exists but lacks
	// the ctx recommendation.
	// DescKeyWriteVscodeAddManually is the text key for write vscode add manually
	// messages.
	DescKeyWriteVscodeAddManually = "write.vscode-add-manually"
	// DescKeyWriteVscodeWarnNonFatal reports a non-fatal error during
	// artifact creation.
	// DescKeyWriteVscodeWarnNonFatal is the text key for write vscode warn non
	// fatal messages.
	DescKeyWriteVscodeWarnNonFatal = "write.vscode-warn-non-fatal"
)
