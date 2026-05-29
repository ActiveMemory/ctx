//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import "time"

// Header is the parsed shape of an audit report's YAML
// frontmatter.
//
// Fields:
//   - Kind: report kind (often matches the file basename)
//   - Status: "findings" or "clean"
//   - CommitRange: git ref range the audit covered
//   - GeneratedAt: UTC timestamp when the audit ran
//   - Generator: skill or tool name that produced the
//     report
//   - Digest: opaque content digest for staleness detection
type Header struct {
	Kind        string    `yaml:"kind"`
	Status      string    `yaml:"status"`
	CommitRange string    `yaml:"commit-range"`
	GeneratedAt time.Time `yaml:"generated-at"`
	Generator   string    `yaml:"generator"`
	Digest      string    `yaml:"digest"`
}
