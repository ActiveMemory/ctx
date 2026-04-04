//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for setup wizard write output.
const (
	// DescKeyWriteSetupDone is the text key for write setup done messages.
	DescKeyWriteSetupDone = "write.setup-done"
	// DescKeyWriteSetupPrompt is the text key for write setup prompt messages.
	DescKeyWriteSetupPrompt = "write.setup-prompt"
	// DescKeyWriteSetupDeployComplete is the text key for write setup deploy
	// complete messages.
	DescKeyWriteSetupDeployComplete = "write.setup-deploy-complete"
	// DescKeyWriteSetupDeployMCP is the text key for write setup deploy mcp
	// messages.
	DescKeyWriteSetupDeployMCP = "write.setup-deploy-mcp"
	// DescKeyWriteSetupDeploySteering is the text key for write setup deploy
	// steering messages.
	DescKeyWriteSetupDeploySteering = "write.setup-deploy-steering"
	// DescKeyWriteSetupDeployExists is the text key for write setup deploy exists
	// messages.
	DescKeyWriteSetupDeployExists = "write.setup-deploy-exists"
	// DescKeyWriteSetupDeployCreated is the text key for write setup deploy
	// created messages.
	DescKeyWriteSetupDeployCreated = "write.setup-deploy-created"
	// DescKeyWriteSetupDeploySynced is the text key for write setup deploy synced
	// messages.
	DescKeyWriteSetupDeploySynced = "write.setup-deploy-synced"
	// DescKeyWriteSetupDeploySkipSteer is the text key for write setup deploy
	// skip steer messages.
	DescKeyWriteSetupDeploySkipSteer = "write.setup-deploy-skip-steer"
)

// DescKeys for setup integration instruction output.
const (
	// DescKeyWriteSetupCursorHead is the Cursor section header.
	DescKeyWriteSetupCursorHead = "write.setup-cursor-head"
	// DescKeyWriteSetupCursorRun is the Cursor run command hint.
	DescKeyWriteSetupCursorRun = "write.setup-cursor-run"
	// DescKeyWriteSetupCursorMCP is the Cursor MCP config path.
	DescKeyWriteSetupCursorMCP = "write.setup-cursor-mcp"
	// DescKeyWriteSetupCursorSync is the Cursor steering sync path.
	DescKeyWriteSetupCursorSync = "write.setup-cursor-sync"
	// DescKeyWriteSetupKiroHead is the Kiro section header.
	DescKeyWriteSetupKiroHead = "write.setup-kiro-head"
	// DescKeyWriteSetupKiroRun is the Kiro run command hint.
	DescKeyWriteSetupKiroRun = "write.setup-kiro-run"
	// DescKeyWriteSetupKiroMCP is the Kiro MCP config path.
	DescKeyWriteSetupKiroMCP = "write.setup-kiro-mcp"
	// DescKeyWriteSetupKiroSync is the Kiro steering sync path.
	DescKeyWriteSetupKiroSync = "write.setup-kiro-sync"
	// DescKeyWriteSetupClineHead is the Cline section header.
	DescKeyWriteSetupClineHead = "write.setup-cline-head"
	// DescKeyWriteSetupClineRun is the Cline run command hint.
	DescKeyWriteSetupClineRun = "write.setup-cline-run"
	// DescKeyWriteSetupClineMCP is the Cline MCP config path.
	DescKeyWriteSetupClineMCP = "write.setup-cline-mcp"
	// DescKeyWriteSetupClineSync is the Cline steering sync path.
	DescKeyWriteSetupClineSync = "write.setup-cline-sync"
	// DescKeyWriteSetupNoSteeringToSync is the message when no
	// steering files are available for sync.
	DescKeyWriteSetupNoSteeringToSync = "write.setup-no-steering-to-sync"
)
