//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package kb supplies directory and filename constants for the
// ctx knowledge-base editorial pipeline (Phase KB).
//
// Names live here, not in [github.com/ActiveMemory/ctx/internal/cli/kb],
// to honor the project-wide rule that magic strings and values
// belong in internal/config/.
//
// # Layout
//
// Under the project's .context/ directory:
//
//	.context/
//	├── kb/
//	│   ├── index.md                  (KBIndex)
//	│   ├── evidence-index.md         (EvidenceIndex)
//	│   ├── glossary.md               (Glossary)
//	│   ├── contradictions.md         (Contradictions)
//	│   ├── outstanding-questions.md  (OutstandingQuestions)
//	│   ├── domain-decisions.md       (DomainDecisions)
//	│   ├── timeline.md               (Timeline)
//	│   ├── source-map.md             (SourceMap)
//	│   ├── source-coverage.md        (SourceCoverage)
//	│   ├── relationship-map.md       (RelationshipMap)
//	│   └── topics/<slug>/index.md    (Topics dir + TopicIndex filename)
//	├── ingest/
//	│   ├── KB-RULES.md               (Rules)
//	│   ├── 00-GROUND.md              (ModeGround)
//	│   ├── 30-INGEST.md              (ModeIngest)
//	│   ├── 40-ASK.md                 (ModeAsk)
//	│   ├── 50-SITE_REVIEW.md         (ModeSiteReview)
//	│   ├── INBOX.md                  (Inbox)
//	│   ├── SESSION_LOG.md            (SessionLog)
//	│   ├── grounding-sources.md      (GroundingSources)
//	│   ├── OPERATOR.md               (Operator)
//	│   ├── PROMPT.md                 (Prompt)
//	│   ├── closeouts/                (CloseoutsSubdir)
//	│   ├── schemas/                  (SchemasSubdir)
//	│   └── findings.md               (Findings; lazy-init via ctx kb note)
//	├── handovers/                    (HandoversSubdir)
//	├── archive/closeouts/            (ArchiveCloseoutsSubdir)
//	└── site/kb/                      (SiteKBSubdir; gitignored)
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/core/path]
//     composes these constants with [rc.ContextDir] to produce
//     full filesystem paths.
//   - [github.com/ActiveMemory/ctx/internal/write/kb] writes the
//     evidence-index, glossary, source-coverage, and other
//     per-artifact files.
package kb
