//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package vscode provides terminal output for VS Code artifact generation.
package vscode

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// InfoCreated reports a VS Code configuration file was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - target: path to the created file
func InfoCreated(cmd *cobra.Command, target string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteVscodeCreated), target))
}

// InfoExistsSkipped reports a VS Code file was skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - target: path to the existing file
func InfoExistsSkipped(cmd *cobra.Command, target string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteVscodeExistsSkipped), target))
}

// InfoRecommendationExists reports the extension recommendation already exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - target: path to the extensions.json file
func InfoRecommendationExists(cmd *cobra.Command, target string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteVscodeRecommendationExists), target))
}

// InfoAddManually reports the file exists but lacks the ctx recommendation.
//
// Parameters:
//   - cmd: Cobra command for output
//   - target: path to the extensions.json file
//   - extensionID: the extension identifier to add
func InfoAddManually(cmd *cobra.Command, target, extensionID string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteVscodeAddManually), target, extensionID))
}

// InfoWarnNonFatal reports a non-fatal error during artifact creation.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: short description of what failed
//   - err: the non-fatal error
func InfoWarnNonFatal(cmd *cobra.Command, name string, err error) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteVscodeWarnNonFatal), name, err))
}
