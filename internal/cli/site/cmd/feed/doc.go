//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package feed implements **`ctx site feed`**, the
// hidden subcommand that generates the **RSS / Atom
// feed** for the project's published blog under the
// zensical site directory.
//
// # Public Surface
//
//   - **[Cmd]**: cobra command with `--input`
//     (blog directory; default `docs/blog/`) and
//     `--output` (feed file; default
//     `site/feed.xml`).
//   - **[Run]**: scans the input directory via
//     [internal/cli/site/core/scan], converts each
//     post into a feed item (title, link, summary,
//     pub date), and writes a valid feed XML
//     document.
//
// # Why "Hidden"
//
// `ctx site` lives under [hiddenCmds] in the root
// registration because feed generation is part of
// the publish pipeline (called by `make site` and
// the `_ctx-blog-changelog` skill), not something
// users invoke at the prompt directly.
//
// # Concurrency
//
// Single-process, sequential.
package feed
