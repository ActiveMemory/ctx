//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package run

import backendPkg "github.com/ActiveMemory/ctx/internal/backend"

// PingResult is printable backend reachability information.
//
// Fields:
//   - Backend: resolved backend name
//   - Endpoint: configured backend endpoint
//   - FirstModel: first model returned by backend model listing
type PingResult struct {
	Backend    string
	Endpoint   string
	FirstModel string
}

// ProposalArtifact is the validation-only proposed-patch artifact.
//
// Fields:
//   - Kind: artifact kind
//   - Backend: backend that generated the response
//   - Model: model reported by the backend
//   - Input: input file path
//   - Emit: requested proposal kinds
//   - Status: artifact status
//   - Response: decoded backend JSON response
type ProposalArtifact struct {
	Kind     string           `json:"kind"`
	Backend  string           `json:"backend"`
	Model    string           `json:"model"`
	Input    string           `json:"input"`
	Emit     []string         `json:"emit"`
	Status   string           `json:"status"`
	Response ProposalResponse `json:"response"`
}

// ProposalResponse is the structured proposal payload.
type ProposalResponse struct {
	Rows     []ProposalRow `json:"rows"`
	Metadata ProposalMeta  `json:"metadata"`
}

// ProposalRow is one proposed artifact row.
type ProposalRow struct {
	Emit  string `json:"emit"`
	Text  string `json:"text"`
	Start int    `json:"start,omitempty"`
	End   int    `json:"end,omitempty"`
}

// ProposalMeta contains reviewable response metadata.
type ProposalMeta struct {
	Backend string `json:"backend"`
	Model   string `json:"model"`
	Input   string `json:"input"`
	Status  string `json:"status"`
}

// resolvedBackend carries the selected backend and metadata.
type resolvedBackend struct {
	name    string
	backend backendPkg.Backend
}
