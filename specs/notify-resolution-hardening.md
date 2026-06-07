---
title: Notify Resolution Hardening
status: proposed
date: 2026-06-02
owner: jose
scope: |-
  behavioral — key resolution (crypto.ResolveKeyPath), notify
  delivery (internal/notify), warn config, docs. Net deletion +
  one visibility fix.
supersedes-section-of:
  - specs/released/v0.8.0/global-encryption-key.md (the
    project-local `.context/.ctx.key` auto-detection tier only;
    the single-global-key decision and the `key_path` override
    survive in stronger form)
related:
  - specs/cwd-anchored-context.md
  - specs/require-git.md
brief-source: |-
  Debated in session ee7d6d68 (2026-06-02) while picking up
  TASKS.md item P0.8.5 ("enable webhook notifications in
  worktrees"). The task premise turned out to be partly wrong;
  the conversation reframed it. Key turns, all empirically
  verified against a built binary + isolated repo + fake webhook
  sink: (1) notify ALREADY works in a worktree with the default
  global key — the failure only reproduces with a project-local
  key, which is itself a deprecated footgun. (2) The fire path
  silently swallows a configured-but-undeliverable webhook —
  that is the real, always-present defect. (3) ctx cannot and
  must not distinguish a worktree from N side-by-side terminals
  in the same dir; whether config (incl. notify.events) reaches
  a worktree is already controlled by whether the user git-tracks
  `.ctxrc`. So "make worktrees fire" needs no code — it is an
  emergent property the user already controls. Convergence:
  remove the project-local key tier, surface the silent failure,
  document the `.ctxrc`-tracking control. Smaller than the task
  implied and aligned with the global-key + cwd-anchored
  simplification lineage.
---

# Spec: Notify Resolution Hardening

## Problem

`TASKS.md` P0.8.5 claims `ctx hook notify` "silently fails because
`.context.key` is gitignored and absent in worktrees." Investigation
found the premise is mostly wrong, but it surfaced two genuine defects
and one design question that this spec settles.

What is actually true (all verified empirically — built binary, isolated
git repo, isolated `$HOME`, local HTTP sink; the real committed webhook
was never fired):

1. **The default config already works in worktrees.** The key file is
   `.ctx.key` (not `.context.key`), and the default is a single global
   key at `~/.ctx/.ctx.key` — shared across every worktree and
   side-by-side terminal on a machine. `.context/.notify.enc` is
   committed, so it rides into a same-branch worktree. A fire from a
   worktree delivered HTTP 200.

2. **A project-local key silently breaks it.** When a project has a
   distinct `<contextDir>/.ctx.key` (the `os.Stat`-gated tier 2 of
   `crypto.ResolveKeyPath`), that file is gitignored and therefore
   absent in a fresh worktree checkout. Resolution silently falls back
   to the global key, decryption fails with
   `cipher: message authentication failed`, and the **fire path drops
   the notification with no stderr and exit 0**. The `test` subcommand
   surfaces the error; the fire path (what autonomous agents call) does
   not.

3. **The project already deprecated project-local keys.** The v0.8.0
   `global-encryption-key.md` spec collapsed per-project keys into a
   single global key precisely because they were "over-engineered for
   the common case" and (per `DECISIONS.md` 2026-03-01) "a security
   antipattern [key next to ciphertext] and broke in worktrees." The
   implicit `.context/.ctx.key` auto-detect tier is the vestige that
   still produces the worktree divergence.

4. **Whether config reaches a worktree is already a user decision.** ctx
   is CWD-anchored (`specs/cwd-anchored-context.md`): it reads
   `$PWD/.ctxrc` or applies defaults, and it cannot distinguish a git
   worktree from N terminals `cd`'d into the same project. So the right
   axis is not "is this a worktree" but "did the user git-track
   `.ctxrc`": a committed `.ctxrc` propagates `notify.events` to every
   checkout (verified); a gitignored/profile `.ctxrc` does not, and the
   worktree falls to defaults. `load()` has **no `.ctxrc.base`
   fallback**, so a tracked profile template does not auto-activate in a
   fresh checkout.

## Decision

Three changes. Two are code (one is net deletion), one is documentation.

### Change 1 — Remove the implicit project-local key tier

`crypto.ResolveKeyPath` drops the `os.Stat(<contextDir>/.ctx.key)`
auto-detection tier. Resolution becomes:

1. **Explicit override** — `.ctxrc key_path` (with `~/` expansion).
2. **Global default** — `~/.ctx/.ctx.key`.
3. **Degenerate fallback** — `<contextDir>/.ctx.key` **only** when the
   home directory is unavailable (so `GlobalKeyPath()` returned `""`).
   This is not the footgun: it is a last resort when no global location
   can be computed, not "silently prefer a stray file next to the
   ciphertext."

Rationale: a project-local key is the single thing that makes a worktree
behave differently from a side-by-side terminal; removing it makes the
two indistinguishable (the desired model) and deletes a documented
security antipattern. Genuine per-project isolation remains available via
the explicit `key_path` override pointing anywhere the user controls.

### Change 2 — Surface configured-but-undeliverable webhooks

`internal/notify` distinguishes "not configured" (legitimately silent)
from "configured but could not deliver" (must be visible), on the fire
path as well as the test path.

`LoadWebhook`: the webhook is "configured" exactly when
`.context/.notify.enc` exists, so a **missing `.notify.enc` is the one
silent "not configured" signal** — `LoadWebhook` stats it first and
returns `("", nil)` on not-exist. `os.Stat` returns an unwrapped
`*fs.PathError`, so the not-exist check is reliable regardless of how
`crypto.LoadKey` wraps its error (its wrap goes through the text
registry, on which neither `os.IsNotExist` nor `errors.Is` is dependable
in an uninitialized test binary — see LEARNINGS 2026-06-02). Once
`.notify.enc` exists, every remaining failure propagates: a missing,
unreadable, or invalid key (`crypto.LoadKey`), an unreadable
`.notify.enc` (`io.SafeReadUserFile`), or a `crypto.Decrypt` failure
(wrong key — the worktree/cross-machine case). This deliberately surfaces
the absent-key-with-present-enc case rather than mistaking it for "no
webhook" (it satisfies the Edge case below and the test-path visibility
in Problem item 2).

`Send`: keep the legitimate silences (event not in `notify.events`;
`url == "" && err == nil`); on a non-nil `LoadWebhook` error (configured
but undeliverable — the worktree/wrong-key case) emit a non-fatal
`logWarn.Warn` and return nil; on `PostJSON` failure emit a non-fatal
`logWarn.Warn` and return nil (delivery failed, but no longer invisible);
on marshal failure warn likewise. `Send` stays non-fatal throughout
(fire-and-forget), but never silent about a real failure.

### Change 3 — Document the `.ctxrc`-tracking control (no code)

We do **not** add machinery to make worktrees fire, and we do not
special-case worktrees anywhere. Instead we document the existing,
correct behavior:

- Whether `notify.events` (and all config) reaches a worktree / parallel
  checkout is controlled by whether the user git-tracks `.ctxrc`.
  Committing `.ctxrc` → fires everywhere; gitignored/profile `.ctxrc` →
  worktrees use defaults.
- Committing `.ctxrc` is **secret-safe**: it holds `notify.events`,
  `key_path`, rotation days, etc. — never the webhook secret, which
  lives encrypted in `.context/.notify.enc`.
- `.ctxrc.base`/`.ctxrc.dev` are profile *templates*, not fallbacks: a
  fresh checkout with a gitignored active `.ctxrc` uses hardcoded
  defaults, ignoring any tracked base.

## Non-goals

- **No worktree detection.** ctx must not branch behavior on "am I in a
  worktree"; it cannot tell a worktree from side-by-side terminals, and
  the CWD-anchored model forbids it.
- **No config propagation / key-copy into worktrees** (rejected approach
  B). It is unenforceable (agent-driven), redundant under a global key,
  and would be ctx deciding intent it cannot observe.
- **No `git rev-parse --git-common-dir` key fallback** (rejected approach
  A). Removing the project-local tier eliminates the divergence outright,
  so there is nothing to redirect.
- **No change to `Send`'s fire-and-forget contract.** Failures become
  visible (warn), not fatal.
- **Cross-machine / fresh-clone decrypt** (committed `.notify.enc` +
  per-user global key on another machine) is out of scope; Change 2
  makes that failure visible, which is sufficient here.

## Detailed changes

| File | Change |
|------|--------|
| `internal/crypto/keypath.go` | `ResolveKeyPath`: delete the tier-2 `os.Stat(local)` block; keep override → global → home-unavailable fallback. Rewrite the doc comment to the 3-step order. |
| `internal/crypto/keypath_test.go` | Remove/replace `TestResolveKeyPath_ProjectLocalBeforeGlobal`; add a test asserting a present `<ctxDir>/.ctx.key` is **ignored** in favor of global; keep override + global + home-unavailable cases. |
| `internal/rc/rc.go` | `KeyPath` doc comment: update the resolution order to override → global (drop the stale project-local-as-tier wording). |
| `internal/notify/notify.go` | `LoadWebhook`: stat `.notify.enc` first (the silent not-configured signal); once present, surface key/read/decrypt errors. `Send`: warn (non-fatal) on configured-but-undeliverable, POST failure, marshal failure. Update both doc comments. |
| `internal/config/warn/warn.go` | Add `NotifyWebhookLoad`, `NotifyWebhookPost`, `NotifyWebhookMarshal` format keys. |
| `internal/log/warn/warn.go` | Add an exported `SetSink(io.Writer) func()` test seam so callers can assert/silence warnings (the doc comment already anticipated this). |
| `internal/notify/notify_test.go` | Cover: invalid-size key → propagated; enc-present-but-key-absent → propagated; corrupted enc → propagated; `Send` warns exactly once on undeliverable and on POST failure (capturing the sink), stays silent + warning-free when not configured. |
| `docs/recipes/parallel-worktrees.md` | Add the `.ctxrc`-tracking control + secret-safety + `.ctxrc.base`-not-a-fallback notes. |
| `internal/cli/notify/doc.go` | Cross-reference the worktree/`.ctxrc`-tracking behavior. |

## Migration

There is no auto-migration and currently no legacy-key warning in the
tree (the old `MigrateKeyFile` warning was already removed). After
Change 1:

- Projects on the **default global key**: no effect.
- Projects with a **project-local `.context/.ctx.key`**: key ops resolve
  to the global key. Their existing `.notify.enc` / pad encrypted with
  the local key will fail to decrypt — and Change 2 makes that failure
  **visible** (a warning on the fire path; a surfaced error on pad/test
  paths) instead of silent. Documented remedy: back up the project-local
  key, then either re-key to the global key (re-run `ctx hook notify
  setup`, re-encrypt the pad) or set an explicit `key_path` override
  pointing at the backed-up key.

Optional follow-up (not required by this spec): a `ctx doctor` check that
detects a stranded `<contextDir>/.ctx.key` and prints the remedy
proactively rather than relying on the surfaced decrypt failure.

## Test plan

- Unit: `ResolveKeyPath` ignores a present project-local key (global
  wins); override still wins; home-unavailable still falls back to local.
- Unit: `LoadWebhook` returns `("", nil)` when `.notify.enc` is missing;
  propagates an invalid-size-key error, an enc-present-but-key-absent
  error, and a decrypt (corrupted-enc) error.
- Unit: `Send` is silent and warning-free when not configured (event
  subscribed, no `.notify.enc`) and when the event is filtered; warns
  exactly once on a configured-but-undeliverable webhook (captured via
  `logWarn.SetSink`); warns exactly once on a POST transport failure;
  stays non-fatal (returns nil) in all cases.
- Full `make lint` + `make test` green; `internal/audit` suite green
  (mixed-visibility, magic-strings, doc-comment audits).

## Edge cases

- **Home dir unavailable** (`os.UserHomeDir` fails): `GlobalKeyPath()`
  returns `""`; resolution uses the degenerate `<contextDir>/.ctx.key`
  fallback. Preserved.
- **`key_path` override (or global key) absent while `.notify.enc`
  exists**: because `LoadWebhook` stats `.notify.enc` first and only then
  loads the key, a configured webhook whose key is missing surfaces an
  error (`Send` warns; the test path reports it) rather than silently
  reading as "not configured". This is the visible-failure behavior the
  override Edge case requires.
- **No `.context/` at `$PWD`**: `rc.ContextDir()` /  `rc.KeyPath()`
  return `errCtx.ErrNoCtxHere`; `LoadWebhook` propagates it; `Send`
  warns (configured-but-undeliverable) — acceptable, since `Send`'s
  callers run only after `RequireContextDir`.
