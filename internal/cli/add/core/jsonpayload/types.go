//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package jsonpayload

// Provenance is the optional commit-trail envelope inside a payload.
//
// Fields:
//   - SessionID: AI session identifier
//   - Branch: Git branch name
//   - Commit: Git commit hash
type Provenance struct {
	SessionID string `json:"session_id"`
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
}

// Payload is the decoded shape of a --json-file argument.
//
// All fields are optional; each add noun consumes only the keys that
// are relevant to it. Extra-but-irrelevant keys (e.g. priority on a
// decision) decode without error and are simply ignored by that noun's
// formatter. A genuinely unknown key is a decode error (see [Load]).
//
// Fields:
//   - Title: Entry content/heading
//   - Body: Extra content appended to Title for tasks (single-line files)
//   - Context: Context field for decisions/learnings
//   - Rationale: Rationale field for decisions
//   - Consequence: Consequence field for decisions
//   - Lesson: Lesson field for learnings
//   - Application: Application field for learnings
//   - Priority: Priority label for tasks
//   - Section: Target section for tasks
//   - Provenance: Optional session/branch/commit envelope
type Payload struct {
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Context     string     `json:"context"`
	Rationale   string     `json:"rationale"`
	Consequence string     `json:"consequence"`
	Lesson      string     `json:"lesson"`
	Application string     `json:"application"`
	Priority    string     `json:"priority"`
	Section     string     `json:"section"`
	Provenance  Provenance `json:"provenance"`
}
