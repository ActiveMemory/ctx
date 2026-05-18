//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package handover writes per-session handover artifacts under
// `.context/handovers/`. The handover is the session-to-session
// glue: a former-agent-to-next-agent note created by
// `/ctx-wrap-up` at session end and read by `/ctx-remember` at
// session start. It is universal to every ctx project and does
// not depend on the editorial pipeline.
//
// A handover carries a past-tense summary plus a future-tense
// "first action for the next session", with optional highlights
// and open questions. Files are timestamped (`<TS>-<slug>.md`)
// so multiple concurrent agent runs never overwrite each other.
//
// # The ceremony
//
// The "skies will fall" ceremony every ctx session runs is:
//
//  1. `/ctx-remember` at session start reads the latest handover.
//  2. The session does its work.
//  3. `/ctx-wrap-up` at session end delegates to this writer as
//     its final step.
//
// Skipping step 3 means the next session's `/ctx-remember` has
// no handover to read and recall degenerates to probabilistic
// reconstruction from canonical files plus journal.
//
// # Optional closeout folding (Phase KB only)
//
// When `.context/kb/` exists, the project also runs the editorial
// pipeline, which writes closeouts under
// `.context/ingest/closeouts/` per pass. As an implementation
// detail, [Write] additionally folds postdated closeouts into the
// handover's `## Folded Closeouts` section and archives them
// under `.context/archive/closeouts/`. A `--no-fold` write skips
// this fold (mid-session checkpoint). When `.context/kb/` is
// absent there are no closeouts to fold; the handover is still
// written normally. The fold is orthogonal to the handover
// mechanism; the handover itself is not a KB feature.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/write/closeout]
//     supplies the optional fold-source files.
//   - [github.com/ActiveMemory/ctx/internal/gitmeta] stamps
//     `sha` / `branch` provenance.
//   - [github.com/ActiveMemory/ctx/internal/cli/handover/core/path]
//     resolves `.context/handovers/`.
package handover
