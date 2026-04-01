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

// vsMCPServer is a typed MCP server entry.
type vsMCPServer struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// vsMCPFile is the top-level .vscode/mcp.json structure.
type vsMCPFile struct {
	Servers map[string]vsMCPServer `json:"servers"`
}

// writeMCPJSON creates .vscode/mcp.json with the ctx MCP server
// registration. Skips if the file already exists.
func writeMCPJSON(cmd *cobra.Command) error {
	target := filepath.Join(cfgVscode.Dir, cfgVscode.FileMCPJSON)

	if _, statErr := os.Stat(target); statErr == nil {
		writeVscode.InfoExistsSkipped(cmd, target)
		return nil
	}

	file := vsMCPFile{
		Servers: map[string]vsMCPServer{
			"ctx": {
				Command: "ctx",
				Args:    []string{"mcp", "serve"},
			},
		},
	}
	data, _ := json.MarshalIndent(file, "", "  ")
	data = append(data, '\n')

	if writeErr := os.WriteFile(target, data, fs.PermFile); writeErr != nil {
		return writeErr
	}
	writeVscode.InfoCreated(cmd, target)
	return nil
}
