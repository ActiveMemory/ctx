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
	writeVscode "github.com/ActiveMemory/ctx/internal/write/vscode"
)

// createTasksJSON creates .vscode/tasks.json with ctx command tasks.
//
// Skips if the file already exists to preserve user customizations.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if writing the file fails
func createTasksJSON(cmd *cobra.Command) error {
	target := filepath.Join(cfgVscode.Dir, cfgVscode.FileTasksJSON)

	if _, statErr := os.Stat(target); statErr == nil {
		writeVscode.InfoExistsSkipped(cmd, target)
		return nil
	}

	tasks := make([]vsTask, len(cfgVscode.Tasks))
	for i, t := range cfgVscode.Tasks {
		tasks[i] = vsTask{
			Label:   t.Label,
			Type:    cfgVscode.TypeShell,
			Command: t.Command,
			Group:   cfgVscode.GroupNone,
			Presentation: vsPresentation{
				Reveal: cfgVscode.RevealAlways,
				Panel:  cfgVscode.PanelShared,
			},
			ProblemMatcher: []string{},
		}
	}

	file := vsTasksFile{
		Version: cfgVscode.TasksVersion,
		Tasks:   tasks,
	}
	data, _ := json.MarshalIndent(file, "", "  ")
	data = append(data, token.NewlineLF...)

	if writeErr := os.WriteFile(target, data, fs.PermFile); writeErr != nil {
		return writeErr
	}
	writeVscode.InfoCreated(cmd, target)
	return nil
}
