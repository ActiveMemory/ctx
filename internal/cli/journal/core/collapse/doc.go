//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package collapse condenses **large tool-output blocks** in
// journal markdown, the multi-thousand-line shell pastes
// and `ls`/`grep` outputs that bloat an entry without
// adding much signal, into expandable summaries that show
// the first few lines and offer the rest under a
// `<details>` toggle.
//
// The package complements [reduce]: reduce strips bona-fide
// noise (system reminders, orphan fences); collapse
// preserves output but **hides** the bulk so reviewers can
// skim and only expand the tool calls they care about.
//
// # Public Surface
//
//   - **[ToolOutputs](content, opts)**: finds tool-output
//     code blocks larger than a configurable line
//     threshold and replaces them with a `<details>`
//     summary block:
//
//     <details><summary>Tool output (NNN lines)</summary>
//
//     ```
//     ...full output...
//     ```
//
//     </details>
//
//     with the first 5 lines shown above the
//     collapsed block as anchor context. Threshold and
//     preview line count are tunable via [opts].
//
// # Why Not Just Truncate?
//
// Truncating loses information. The journal entry is a
// **record**; the user may need the full output later
// to reconstruct what happened. Collapsing wins on both
// fronts: the rendered page is short and skimmable, the
// raw markdown still contains every byte of the original
// output.
//
// # Concurrency
//
// Pure data transformation. Concurrent callers never
// race.
package collapse
