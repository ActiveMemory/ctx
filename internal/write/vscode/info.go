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
)

// InfoCreated reports a VS Code configuration file was created.
func InfoCreated(cmd *cobra.Command, target string) {
	cmd.Println(fmt.Sprintf("  ✓ %s", target))
}

// InfoExistsSkipped reports a VS Code file was skipped because it exists.
func InfoExistsSkipped(cmd *cobra.Command, target string) {
	cmd.Println(fmt.Sprintf("  ○ %s (exists, skipped)", target))
}

// InfoRecommendationExists reports the extension recommendation already exists.
func InfoRecommendationExists(cmd *cobra.Command, target string) {
	cmd.Println(fmt.Sprintf("  ○ %s (recommendation exists)", target))
}

// InfoAddManually reports the file exists but lacks the ctx recommendation.
func InfoAddManually(cmd *cobra.Command, target, extensionID string) {
	cmd.Println(fmt.Sprintf("  ○ %s (exists, add %s manually)", target, extensionID))
}

// InfoWarnNonFatal reports a non-fatal error during artifact creation.
func InfoWarnNonFatal(cmd *cobra.Command, name string, err error) {
	cmd.Println(fmt.Sprintf("  ⚠ %s: %v", name, err))
}
