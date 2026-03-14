---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Version History
icon: lucide/history
---

![ctx](../images/ctx-banner.png)

## Version History

Documentation snapshots for each release. 

Tap the corresponding **view docs** to view the docs as they were at that release.

## Releases

| Version | Release Date | Documentation                                                     |
|---------|--------------|-------------------------------------------------------------------|
| v0.6.0  | 2026-02-16   | [view docs](https://github.com/ActiveMemory/ctx/tree/v0.6.0/docs) |
| v0.3.0  | 2026-02-07   | [view docs](https://github.com/ActiveMemory/ctx/tree/v0.3.0/docs) |
| v0.2.0  | 2026-02-01   | [view docs](https://github.com/ActiveMemory/ctx/tree/v0.2.0/docs) |
| v0.1.2  | 2026-01-27   | [view docs](https://github.com/ActiveMemory/ctx/tree/v0.1.2/docs) |
| v0.1.1  | 2026-01-26   | [view docs](https://github.com/ActiveMemory/ctx/tree/v0.1.1/docs) |
| v0.1.0  | 2026-01-25   | [view docs](https://github.com/ActiveMemory/ctx/tree/v0.1.0/docs) |

### `v0.6.0`: The Integration Release

Plugin architecture: hooks and skills converted from shell scripts to Go
subcommands, shipped as a Claude Code marketplace plugin. Multi-tool hook
generation for Cursor, Aider, Copilot, and Windsurf. Webhook notifications
with encrypted URL storage.

### `v0.3.0`: The Discipline Release

Journal static site generation via zensical. 49-skill audit and fix pass
(positive framing, phantom reference removal, scope tightening).
Context consolidation skill. golangci-lint v2 migration.

### `v0.2.0`: The Archaeology Release

Session journal system: `ctx recall export` converts Claude Code JSONL
transcripts to browsable Markdown. Constants refactor with semantic
prefixes (`Dir*`, `File*`, `Filename*`). CRLF handling for Windows
compatibility.

### `v0.1.2`

Default Claude Code permissions deployed on `ctx init`. Prompting guide
published as a standalone documentation page.

### `v0.1.1`

Bug fixes: hook schema key format corrected, JSON unicode escaping
fixed in context file output.

### `v0.1.0`: Initial Release

CLI with 15 subcommands, 6 context file types (CONSTITUTION, TASKS,
CONVENTIONS, ARCHITECTURE, DECISIONS, LEARNINGS), Makefile build system,
and Claude Code hook integration.

## Latest

The [main documentation](../index.md) always reflects the latest development version.

For the most recent stable release, see
[v0.6.0](https://github.com/ActiveMemory/ctx/tree/v0.6.0/docs).

## Changelog

For detailed changes between versions, see the 
[GitHub Releases](https://github.com/ActiveMemory/ctx/releases) page.
