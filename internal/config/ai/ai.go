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
	StatusProposed             = "proposed"
	TimestampLayout            = "20060102T150405.000000000Z07:00"
	WritePingFormat            = "backend: %s\nendpoint: %s\nfirst_model: %s\n"
)

// ProposalSchema is the structured response schema for `ctx ai propose`.
const ProposalSchema = `{"type":"object","required":["rows","metadata"],` +
	`"properties":{"rows":{"type":"array","minItems":1,` +
	`"items":{"type":"object","required":["emit","text"],` +
	`"properties":{"emit":{"type":"string"},` +
	`"text":{"type":"string"},` +
	`"start":{"type":"integer"},"end":{"type":"integer"}},` +
	`"additionalProperties":false}},` +
	`"metadata":{"type":"object",` +
	`"required":["backend","model","input","status"],` +
	`"properties":{"backend":{"type":"string"},` +
	`"model":{"type":"string"},"input":{"type":"string"},` +
	`"status":{"type":"string"}},"additionalProperties":false}},` +
	`"additionalProperties":false}`
