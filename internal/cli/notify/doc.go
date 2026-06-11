//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package notify implements the **`ctx hook notify`**
// command surface (webhook send, setup, and test)
// that wraps the in-process [internal/notify] engine for
// CLI use.
//
// The command lives under `ctx hook` rather than at root
// because notifications belong to the **hook subsystem**
// (delivered when hooks fire); see
// `internal/cli/hook/hook.go` for the parent registration.
//
// # Subcommands
//
//   - **`ctx hook notify [message]`**: fire-and-forget
//     send. Required: `--event <name>`. Optional:
//     `--session-id`, `--hook`, `--variant`. Honors
//     the `notify.events` filter in `.ctxrc`; silent
//     no-op when the event is not whitelisted.
//   - **`ctx hook notify setup`**: interactive prompt
//     to capture and encrypt the webhook URL. See
//     [internal/cli/notify/cmd/setup].
//   - **`ctx hook notify test`**: sends a test event,
//     **bypassing** the event filter so users can
//     verify connectivity without subscribing the test
//     event first. See [internal/cli/notify/cmd/test].
//
// # Worktrees and parallel checkouts
//
// `ctx hook notify` does not special-case git worktrees — it
// cannot distinguish one from several terminals open in the same
// project. Whether the `notify.events` filter reaches a worktree
// depends on whether `.ctxrc` is git-tracked (committed → fires
// everywhere; gitignored → worktrees fall back to defaults). The
// key is the single global `~/.ctx/.ctx.key`, shared by all
// checkouts. See `docs/recipes/parallel-worktrees.md`.
//
// # Concurrency
//
// Stateless. The CLI command spawns one HTTP request
// and exits.
package notify
