//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package message loads, renders, and lets users override
// the **hook message templates** every nudge in ctx prints
// — the small bits of text like "Context checkpoint
// reached" or "Consider archiving completed tasks" that
// hooks emit through the VERBATIM relay channel.
//
// The package backs `ctx hook message list / show / edit /
// reset` plus the in-process consumers in every
// `cmd/check_*` hook.
//
// # Public Surface
//
//   - **[Load](key, vars)** — resolves a message key to
//     its rendered text. Looks up the user's per-project
//     override at [OverridePath](key) first; falls back
//     to the embedded default from
//     [internal/assets/hooks/messages] when no override
//     exists. Sprintf-substitutes [vars] into the
//     template.
//   - **[BoxLines](text)** — wraps text in the
//     box-drawing border ctx uses to make hook nudges
//     visually distinct from normal output.
//   - **[NudgeBox](text)** — boxed-and-prefixed
//     convenience for the standard nudge banner.
//   - **[FormatTemplateVars](vars)** — exposes the
//     normalized key/value pairs hooks pass to [Load].
//   - **[OverridePath](key)** — returns the per-project
//     override file path for `key`. The CLI uses this
//     for the `edit` subcommand; the resolver uses it
//     to detect override presence.
//
// # Override Workflow
//
// Users edit messages by running `ctx hook message edit
// <key>`. The CLI:
//
//  1. Computes [OverridePath](key).
//  2. Materializes the embedded default at that path
//     so the user has something concrete to edit.
//  3. Spawns `$EDITOR` on the file.
//  4. Subsequent loads pick up the override
//     automatically.
//
// `ctx hook message reset <key>` deletes the override
// so the embedded default takes over again.
//
// # Concurrency
//
// File reads are scoped per call; no module-level
// caches. Concurrent loads are safe.
//
// # Related Packages
//
//   - [internal/cli/message]              — the
//     `ctx hook message *` CLI surface.
//   - [internal/assets/hooks/messages]    — the
//     embedded default templates.
//   - [internal/cli/system/cmd/check_*]   — every hook
//     consumes a message through this package.
//   - [internal/cli/system/core/nudge]    — the
//     emitter that renders [NudgeBox] output.
package message
