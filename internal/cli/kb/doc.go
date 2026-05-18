//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package kb wires the `ctx kb` parent command and its
// subcommand tree (Phase KB).
//
// Subcommands:
//
//   - ctx kb topic new <name>: scaffold a folder-shaped topic.
//   - ctx kb note <text>: lightweight capture into findings.md.
//   - ctx kb ingest: invoke the editorial pass (skill-driven).
//   - ctx kb ask: Q&A grounded in the kb (skill-driven).
//   - ctx kb site-review: mechanical audit (skill-driven).
//   - ctx kb ground: re-grounding pass (skill-driven).
//   - ctx kb reindex: refresh the CTX:KB:TOPICS managed block.
//   - ctx kb site build/serve/customize: render via zensical.
//
// See specs/kb-editorial-pipeline.md for the editorial
// pipeline contract.
package kb
