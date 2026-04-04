//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

// MCP server identity constants.
const (
	// ResourceURIPrefix is the URI scheme prefix for MCP context resources.
	ResourceURIPrefix = "ctx://context/"
	// JSONRPCVersion is the JSON-RPC protocol version string.
	JSONRPCVersion = "2.0"
	// Name is the server name reported during initialization.
	Name = "ctx"
	// Command is the binary name used to launch the MCP server.
	Command = "ctx"
	// SubcommandServe is the serve subcommand under mcp.
	SubcommandServe = "serve"
)

// PollIntervalSec is the default interval in seconds for
// resource change polling.
const PollIntervalSec = 5

// Args returns the CLI arguments to launch the ctx MCP server.
func Args() []string {
	return []string{"mcp", SubcommandServe}
}
