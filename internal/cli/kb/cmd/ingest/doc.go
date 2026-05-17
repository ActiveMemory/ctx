//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ingest implements `ctx kb ingest`.
//
// The editorial pass itself is performed by the /ctx-kb-ingest
// skill (the agent reads .context/ingest/30-INGEST.md and follows
// the pass-mode + circuit-breaker contract). This CLI surface
// surfaces the workflow to non-Claude users (printing the
// canonical invocation), and refuses cleanly on empty input so
// the contract is the same in both surfaces.
//
// # Refusal contract
//
//   - No source argument:
//     [github.com/ActiveMemory/ctx/internal/err/kb/cli.ErrIngestNoSources].
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb/cli] supplies
//     the hint strings.
package ingest
