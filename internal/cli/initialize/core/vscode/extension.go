//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package vscode

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgVscode "github.com/ActiveMemory/ctx/internal/config/vscode"
	"github.com/ActiveMemory/ctx/internal/io"
	writeVscode "github.com/ActiveMemory/ctx/internal/write/vscode"
)

// createExtensionsJSON creates .vscode/extensions.json with the ctx
// extension recommendation.
//
// If the file exists and already contains the recommendation, it is
// skipped. If the file exists without the recommendation, the user
// is prompted to add it manually.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if reading or writing the file fails
func createExtensionsJSON(cmd *cobra.Command) error {
	target := filepath.Join(cfgVscode.Dir, cfgVscode.FileExtensionsJSON)

	if _, statErr := os.Stat(target); statErr == nil {
		data, readErr := io.SafeReadUserFile(target)
		if readErr != nil {
			return readErr
		}
		var existing map[string][]string
		if json.Unmarshal(data, &existing) == nil {
			for _, r := range existing[cfgVscode.KeyRecommendations] {
				if r == cfgVscode.ExtensionID {
					writeVscode.InfoRecommendationExists(cmd, target)
					return nil
				}
			}
		}
		writeVscode.InfoAddManually(cmd, target, cfgVscode.ExtensionID)
		return nil
	}

	content := map[string][]string{
		cfgVscode.KeyRecommendations: {cfgVscode.ExtensionID},
	}
	data, _ := json.MarshalIndent(content, "", "  ")
	data = append(data, token.NewlineLF...)

	if writeErr := os.WriteFile(target, data, fs.PermFile); writeErr != nil {
		return writeErr
	}
	writeVscode.InfoCreated(cmd, target)
	return nil
}
