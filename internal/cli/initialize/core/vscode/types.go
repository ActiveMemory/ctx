//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package vscode

// vsTask represents a single VS Code task definition for tasks.json.
//
// Fields:
//   - Label: Human-readable name shown in the task picker
//   - Type: Execution type (e.g. "shell")
//   - Command: Shell command to execute
//   - Group: Task group classification (e.g. "none")
//   - Presentation: Terminal display settings
//   - ProblemMatcher: Patterns for parsing task output (empty for ctx)
type vsTask struct {
	Label          string         `json:"label"`
	Type           string         `json:"type"`
	Command        string         `json:"command"`
	Group          string         `json:"group"`
	Presentation   vsPresentation `json:"presentation"`
	ProblemMatcher []string       `json:"problemMatcher"`
}

// vsPresentation controls how the task terminal panel is displayed.
//
// Fields:
//   - Reveal: When to reveal the terminal (e.g. "always")
//   - Panel: Terminal reuse strategy (e.g. "shared")
type vsPresentation struct {
	Reveal string `json:"reveal"`
	Panel  string `json:"panel"`
}

// vsTasksFile is the top-level structure for .vscode/tasks.json.
//
// Fields:
//   - Version: Tasks schema version (e.g. "2.0.0")
//   - Tasks: List of task definitions
type vsTasksFile struct {
	Version string   `json:"version"`
	Tasks   []vsTask `json:"tasks"`
}

// vsMCPServer represents a single MCP server entry in mcp.json.
//
// Fields:
//   - Command: Executable to launch the server
//   - Args: Command-line arguments passed to the server
type vsMCPServer struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// vsMCPFile is the top-level structure for .vscode/mcp.json.
//
// Fields:
//   - Servers: Map of server name to server configuration
type vsMCPFile struct {
	Servers map[string]vsMCPServer `json:"servers"`
}
