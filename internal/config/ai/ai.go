//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ai

// Output labels and proposal constants.
const (
	ArtifactExtJSON   = ".json"
	ArtifactPrefix    = "proposal-"
	DirAI             = "ai"
	DirProposals      = "proposals"
	EmitSeparator     = ","
	KindProposedPatch = "proposed-patch"
	PromptPrefix      = "Return JSON for requested emit kinds: "
	SchemaMinimal     = `{"type":"object"}`
	StatusProposed    = "proposed"
	TimestampLayout   = "20060102T150405Z"
	WritePingFormat   = "backend: %s\nendpoint: %s\nfirst_model: %s\n"
)
