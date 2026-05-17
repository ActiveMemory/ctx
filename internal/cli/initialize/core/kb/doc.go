//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package kb extracts the embedded KB editorial-pipeline
// templates from [github.com/ActiveMemory/ctx/internal/assets]
// into a freshly-initialized .context/ directory.
//
// Layout written (relative to contextDir):
//
//	kb/
//	├── index.md                  (from kb/templates/kb/index.md)
//	└── topics/.gitkeep
//	ingest/
//	├── KB-RULES.md
//	├── 00-GROUND.md
//	├── 30-INGEST.md
//	├── 40-ASK.md
//	├── 50-SITE_REVIEW.md
//	├── OPERATOR.md
//	├── PROMPT.md
//	├── closeouts/.gitkeep
//	└── schemas/
//	    └── *.md (10 files)
//	handovers/.gitkeep
//
// Existing files are left untouched (init never overwrites
// curated content).
package kb
