# Model Context Window: Claude 5 (Fable/Mythos/Sonnet 5) Detection

**Status:** Implemented

## Problem

`ModelContextWindow` (internal/cli/system/core/session/session_token.go)
recognizes 1M-window models by the explicit `[1m]` suffix or the `opus`
family substring. The Claude 5 Mythos-class models (`claude-fable-5`,
`claude-mythos-5`) match neither, so detection falls through to
`rc.DefaultContextWindow` (200k). Every consumer of
`EffectiveContextWindow` — the check-context-size hook, heartbeat,
nudges, provenance — then computes percentages against 200k on a 1M
session: at ~208k real usage the hook warned "104% full" while the
actual window was 79% free.

## Fix

Add `ModelFable` ("fable"), `ModelMythos` ("mythos"), and
`ModelSonnet5` ("sonnet-5") family substrings to
`internal/config/claude` and include them in the always-1M branch
of `ModelContextWindow`, alongside `ModelOpus`. Substring matching
over the folded model ID mirrors the existing Opus detection, so
dated variants (`claude-fable-5-YYYYMMDD`) are covered. Sonnet 5's
1M window is standard per Anthropic's current model catalog
(unlike earlier Sonnets, whose 1M is an explicit `[1m]` opt-in
already handled by `ModelSuffix1M`). "sonnet-5" does not
substring-match `claude-sonnet-4-5`, so legacy Sonnets are
unaffected.

## Verification

- `TestModelContextWindow` gains cases: `claude-fable-5`,
  `claude-mythos-5`, and `claude-sonnet-5` →
  `claude.ContextWindow1M`.
- `go test ./internal/cli/system/core/session/...` passes.

## Non-Goals

- Re-mapping `claude-sonnet-4-6` (also 1M at the API level): ctx
  deliberately models Claude Code's per-session gating for that
  generation — 200k default, 1M via the `[1m]` suffix or settings
  opt-in — and the existing test pins that behavior. Revisit only
  with session-level evidence.
- A remote model-capability registry; the substring table is the
  intended low-maintenance mechanism.
