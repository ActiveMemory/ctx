//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgVscode "github.com/ActiveMemory/ctx/internal/config/vscode"
	writeVscode "github.com/ActiveMemory/ctx/internal/write/vscode"
)

// vsTask is a typed VS Code task definition.
type vsTask struct {
	Label          string         `json:"label"`
	Type           string         `json:"type"`
	Command        string         `json:"command"`
	Group          string         `json:"group"`
	Presentation   vsPresentation `json:"presentation"`
	ProblemMatcher []string       `json:"problemMatcher"`
}

// vsPresentation controls how the task terminal is displayed.
type vsPresentation struct {
	Reveal string `json:"reveal"`
	Panel  string `json:"panel"`
}

// vsTasksFile is the top-level .vscode/tasks.json structure.
type vsTasksFile struct {
	Version string   `json:"version"`
	Tasks   []vsTask `json:"tasks"`
}

// writeTasksJSON creates .vscode/tasks.json with ctx command tasks.
// Skips if the file already exists.
func writeTasksJSON(cmd *cobra.Command) error {
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
	data = append(data, '\n')

	if writeErr := os.WriteFile(target, data, fs.PermFile); writeErr != nil {
		return writeErr
	}
	writeVscode.InfoCreated(cmd, target)
	return nil
}
