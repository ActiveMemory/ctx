//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

import "time"

// CloseoutFrontmatter holds the six required frontmatter fields
// of a closeout file. PassMode is empty for non-ingest modes
// (ask, site-review, ground, note); LifeStage is empty when the
// pass did not perform a topic-page synthesis.
type CloseoutFrontmatter struct {
	SHA         string    `yaml:"sha"`
	Branch      string    `yaml:"branch"`
	Mode        string    `yaml:"mode"`
	PassMode    string    `yaml:"pass-mode,omitempty"`
	LifeStage   string    `yaml:"life-stage,omitempty"`
	GeneratedAt time.Time `yaml:"generated-at"`
}

// CloseoutFile pairs a closeout's on-disk path with its parsed
// frontmatter and the raw body bytes (everything after the
// closing `---`).
type CloseoutFile struct {
	Path        string
	Frontmatter CloseoutFrontmatter
	Body        string
}
