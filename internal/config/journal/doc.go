//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package journal centralizes every constant that
// drives the journal subsystem: session export,
// enrichment pipeline, site generation, recall
// display, and the zensical navigation sidebar.
//
// # Export and Splitting
//
// Sessions exported from Claude Code are split into
// multipart files when they exceed MaxMessagesPerPart
// (200 messages). The MultipartSuffix ("-p") is
// appended to the base slug to form part filenames
// like "session-slug-p2.md".
//
// # Processing Stages
//
// Each journal entry moves through a linear pipeline
// tracked by stage constants:
//
//   - StageExported  -- raw export from Claude Code
//   - StageEnriched  -- metadata added (tags, summary)
//   - StageNormalized -- content normalized for render
//   - StageFencesVerified -- code fences validated
//   - StageLocked    -- read-only, processing complete
//
// Stage state is persisted in File (".state.json")
// inside .context/journal/.
//
// # Site Generation
//
// The site generator uses PopularityThreshold (2) to
// decide which topics earn dedicated pages and
// LineWrapWidth (80) for soft-wrapping output. The
// navigation sidebar caps at MaxRecentSessions (20)
// with titles truncated to MaxNavTitleLen (40 chars).
//
// # Recall Display
//
// The recall show/list commands use PreviewMaxTurns,
// PreviewMaxTextLen, SlugMaxLen, SessionIDShortLen,
// and SessionIDHintLen to control how sessions are
// summarized in terminal output.
//
// # Boilerplate Detection
//
// BoilerplateNoMatch, BoilerplateFilePrefix,
// BoilerplateFileSuffix, and BoilerplateDenied are
// fixed strings emitted by Claude Code. The journal
// parser matches them verbatim to strip tool noise
// from exported conversations.
//
// # Template Variables
//
// VarUnenrichedCount and VarUnimportedCount are the
// keys injected into hook templates so nudge messages
// can report how many entries need processing.
package journal
