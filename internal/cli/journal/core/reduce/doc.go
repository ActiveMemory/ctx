//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reduce strips noise out of raw AI session JSONL so
// the journal markdown a user reads is the conversation, not
// the wire format. The package is the **noise-removal pass**
// the importer runs before the entry hits disk.
//
// # What Gets Reduced
//
//   - **[StripFences](text)**: removes orphan code-fence
//     markers left by the model when it abandoned a code
//     block mid-response. Without this the renderer
//     enters "code mode" for the rest of the document.
//   - **[StripSystemReminders](text)**: Claude Code
//     injects `<system-reminder>` tags into tool results
//     to nudge the model. The user did not write them
//     and they should not appear in the journal. (See
//     also [internal/parse.StripSystemReminders] which
//     is the shared underlying helper.)
//   - **[CleanToolOutputJSON](text)**: collapses raw
//     JSON tool output into a more readable summary
//     (top-level keys + first values + truncation
//     notice) so a 5,000-line `ls` does not balloon
//     the journal entry. The original is preserved
//     under a `<details>` toggle for archival.
//
// # Idempotency
//
// All three functions are idempotent: running them
// twice in a row on the same input produces the same
// output as running them once. This is what makes
// them safe to run again during re-import.
//
// # Concurrency
//
// All functions are pure data transformations.
// Concurrent callers never race.
package reduce
