//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package vscode

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

// vsMCPServer is a typed MCP server entry.
type vsMCPServer struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// vsMCPFile is the top-level .vscode/mcp.json structure.
type vsMCPFile struct {
	Servers map[string]vsMCPServer `json:"servers"`
}
