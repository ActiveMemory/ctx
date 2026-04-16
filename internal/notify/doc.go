//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package notify implements **fire-and-forget webhook
// notifications**: ctx posts a small JSON payload to a
// user-configured URL when something interesting happens
// (loop completion, hook nudge, version mismatch, key-rotation
// reminder, etc.) and never blocks the caller waiting for the
// response.
//
// The package is what backs `ctx hook notify`,
// `ctx hook notify setup`, and `ctx hook notify test` on the
// CLI side, plus the in-process callers like the autonomous
// loop runner.
//
// # End-to-End Flow
//
//  1. **Setup** ([SaveWebhook]) encrypts a webhook URL with
//     AES-256-GCM ([internal/crypto]) and writes
//     `.context/.notify.enc`. The same per-machine key
//     protects the scratchpad; a fresh key is generated and
//     saved on first use if none exists.
//  2. **Send** ([Send]) loads + decrypts the URL via
//     [LoadWebhook], gates on the configured event filter
//     via [EventAllowed], builds an [entity.NotifyPayload],
//     and ships it to [PostJSON].
//  3. **PostJSON** does the actual HTTP: short timeout,
//     `Content-Type: application/json`, single attempt, no
//     retry. The intent is "best-effort signal", not "guaranteed
//     delivery".
//
// All three functions return cleanly when nothing is
// configured: `("", nil)` from [LoadWebhook] when either
// the key or the encrypted URL file is missing, and a
// silent noop from [Send].
//
// # Event Filter
//
// `notify.events` in `.ctxrc` is **opt-in**: empty list
// means **no events fire** (not "all events"). Recognized
// events: `loop`, `nudge`, `relay`, `heartbeat`. The filter
// is enforced by [EventAllowed].
//
// **`ctx hook notify test` bypasses the filter** as a
// special case so users can verify connectivity without
// having to subscribe their target event first; the test
// path warns when an unfiltered event would normally have
// been dropped.
//
// # Template References
//
// Some emitters attach a [entity.TemplateRef] (hook name +
// variant) to the payload so downstream relays can render a
// canonical message. [template_ref.go] holds the helpers
// that resolve a [TemplateRef] to its rendered string at
// the receiving end (used by integrations that re-emit
// via Slack/Discord/ntfy.sh).
//
// # Encryption Key
//
// The encryption key is shared by both `ctx pad` and
// `ctx hook notify`. Rotating it (every
// `key_rotation_days`, default 90) requires re-running
// `ctx pad init` *and* `ctx hook notify setup`. The
// rotation nudge fires from
// `internal/cli/system/cmd/check_version`.
//
// # Concurrency
//
// All exported functions are safe to call concurrently;
// they hold no module-level state. The HTTP client is the
// stdlib default, connection-pooled and goroutine-safe.
package notify
