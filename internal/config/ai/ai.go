//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ai

// Output labels and proposal constants.
const (
	ArtifactExtJSON            = ".json"
	ArtifactNonceFormat        = "-%06d"
	ArtifactPrefix             = "proposal-"
	DirAI                      = "ai"
	DirProposals               = "proposals"
	ErrEmitRequired            = "emit is required"
	ErrInvalidArtifact         = "invalid proposal artifact"
	ErrInvalidArtifactResponse = "invalid proposal artifact response"
	EmitSeparator              = ","
	KindProposedPatch          = "proposed-patch"
	PromptPrefix               = "Return JSON for requested emit kinds: "
	SchemaProposalName         = "proposal"
	SchemaMinimal              = `{"type":"object"}`
	StatusProposed             = "proposed"
	TimestampLayout            = "20060102T150405.000000000Z07:00"
	WritePingFormat            = "backend: %s\nendpoint: %s\nfirst_model: %s\n"
)
