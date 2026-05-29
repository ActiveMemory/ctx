//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package store

import "time"

// Report is one parsed audit report from
// `.context/audit/<id>.md`.
//
// Fields:
//   - ID: report basename without extension (e.g. "surface")
//   - Path: absolute path of the report file
//   - Kind: kind label from frontmatter (often == ID)
//   - Status: "findings" or "clean"
//   - CommitRange: git ref range the audit covered
//   - GeneratedAt: UTC timestamp from the frontmatter
//   - Generator: skill or tool that produced the report
//   - Digest: opaque content digest for staleness detection
//   - Body: report body verbatim (everything after the
//     frontmatter), suitable for direct verbatim relay
type Report struct {
	ID          string
	Path        string
	Kind        string
	Status      string
	CommitRange string
	GeneratedAt time.Time
	Generator   string
	Digest      string
	Body        string
}

// DismissalLedger is the persisted dismissal state, keyed
// by report id and bound to the digest that was dismissed
// (so a fresh audit overwrite cancels the prior dismissal).
type DismissalLedger struct {
	Entries map[string]DismissedAt `json:"entries"`
}

// DismissedAt records when and against what digest a
// report id was dismissed.
//
// Fields:
//   - Digest: report digest at the moment of dismissal
//   - At: UTC timestamp of the dismissal
type DismissedAt struct {
	Digest string    `json:"digest"`
	At     time.Time `json:"at"`
}
