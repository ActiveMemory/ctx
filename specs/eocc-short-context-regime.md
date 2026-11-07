# EoCC Short-Context Regime and Handover Packet Injection

## Problem

`arXiv:2611.05306` (*You don't need long context*) demonstrates two
load-bearing facts about agentic software engineering workloads:

1. **Long contexts degrade decode efficiency without improving accuracy.**
   Serving frontier and open-weight coding models at 256K tokens costs
   3.6x to 7x more per generated token and doubles inter-token latency
   compared to 16K-48K contexts. With End-of-Conversation Compaction
   (EoCC) triggered at 18K-36K tokens (inside a 24K-48K total limit),
   Claude Opus 4.5 and GPT 5.6 Sol match or beat full-history accuracy on
   SWE-bench Verified (+4.8 and +4.4 percentage points) and Terminal-Bench
   (+5.6 percentage points).
2. **Fresh contexts require immediate semantic state injection.** In EoCC,
   when a conversation crosses the compaction threshold, the agent writes
   a structured execution brief and starts a fresh conversation seeded
   with the system prompt, task, and that brief.

In `ctx` today, three gaps prevent an end-to-end EoCC workflow:

1. `check-context-size` (`internal/cli/system/cmd/checkcontextsize/run.go`)
   only triggers percentage-based nudges at 60% and 90% of the model's
   context window (`EffectiveContextWindow`). On a 200K model, 60% is
   120K tokens; on a 1M model, 60% is 600K tokens. Both are far past the
   18K-36K EoCC sweet spot.
2. `ctx agent` (`internal/cli/agent/core/budget/assemble.go`) compiles
   `CONSTITUTION.md`, `TASKS.md`, `CONVENTIONS.md`, `DECISIONS.md`, and
   `LEARNINGS.md`, but never includes the latest handover from
   `.context/handovers/`. When a user starts a fresh session (or when a
   tool compacts context), the auto-injected `ctx agent` packet lacks the
   immediate execution state (`Summary`, `Established facts`,
   `Applied changes and verification`, `Dead ends`, `Next session`).
3. OpenCode's `experimental.session.compacting` handler
   (`internal/assets/integrations/opencode/plugin/index.ts`) runs
   `ctx system bootstrap` (which outputs only directory metadata) instead
   of injecting the `ctx agent` context packet.

## Approach

1. **Configurable EoCC Token Threshold in `.ctxrc`:**
   Add `eocc_threshold_tokens` (`int`, default `0` = disabled) to
   `rc.CtxRC`. When set (for example, `24000` or `36000`),
   `check-context-size` evaluates whether the current session token count
   has reached `rc.EoCCThresholdTokens()`.
2. **EoCC Compaction Nudge in `check-context-size`:**
   When `tokens >= eoccThreshold` (and the session is not in wrap-up
   suppression), `check-context-size` emits a high-priority EoCC
   compaction nudge (`hook.VariantEoCC`) instructing the agent to stop
   calling tools, persist the forward-looking execution brief via
   `ctx handover write --no-fold`, and advise resetting the conversation
   (`/clear` or a fresh session).
3. **Latest Handover Inclusion in `ctx agent` (Tier 1.5):**
   Extend `AssemblePacket` in `internal/cli/agent/core/budget/assemble.go`
   to load the newest handover file from `.context/handovers/` (if present)
   as **Tier 1.5** immediately after `Tier 1 (Constitution)`. Allocate up
   to 15% of the total token budget (clamped to remaining budget) to the
   latest handover's execution-brief sections (`Summary`,
   `Established facts`, `Applied changes and verification`, `Dead ends`,
   `Next session`), rendered under a `## Latest Handover` section in the
   Markdown packet and a `"handover"` field in the JSON packet. Provide a
   `--no-handover` flag on `ctx agent` to opt out when a caller wants only
   canonical files.
4. **OpenCode Compaction Parity:**
   Update `internal/assets/integrations/opencode/plugin/index.ts` so
   `experimental.session.compacting` runs `ctx agent --budget 4000`
   (which now includes the latest handover) and appends the packet to
   `output.context`.

## Behavior

### Happy Path

1. A project sets `eocc_threshold_tokens: 30000` in `.ctxrc`.
2. During an interactive or autonomous session, `check-context-size` reads
   live token usage from the session log (`coreSession.ReadTokenInfo`).
3. Once `tokens >= 30000`, `check-context-size` emits the EoCC compaction
   nudge box and logs an `eocc` event.
4. The agent runs `ctx handover write "<title>" --no-fold ...` with the
   five execution-brief fields and runs `ctx system mark-wrapped-up` (or
   the user resets the session).
5. On the next session start (or after compaction in Pi / OpenCode /
   Claude Code), `ctx agent` runs automatically, loads the newly written
   handover into Tier 1.5 of the context packet, and presents the exact
   execution state to the fresh model instance.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `eocc_threshold_tokens` is unset, `0`, or negative in `.ctxrc` | EoCC threshold trigger is disabled; `check-context-size` uses only the existing 60%/90% window triggers and `billing_token_warn`. |
| `eocc_threshold_tokens` is set below `4000` (e.g., `500`) | `rc.EoCCThresholdTokens()` clamps values in `(0, 4000)` to `4000` (minimum viable context budget) so a typo cannot trigger compaction on turn 1. |
| `.context/handovers/` is absent or empty when `ctx agent` runs | Tier 1.5 is skipped cleanly with zero tokens consumed; the rest of the packet assembles as before. |
| Latest handover exceeds 15% of the `ctx agent` token budget | The handover body is truncated on a line boundary to fit `min(int(float64(budget)*0.15), remaining)` with a pointer to the handover file path in `.context/handovers/`. |
| User ran `/ctx-wrap-up` recently (`check.WrappedUpRecently()`) | EoCC nudge is suppressed alongside checkpoint and window nudges so wrap-up ceremony commands are not interrupted. |
| Both `eocc_threshold_tokens` and the 60% checkpoint trigger are crossed on the same turn | The EoCC nudge takes precedence over the generic 60% checkpoint nudge (and marks the checkpoint guard) to avoid emitting two overlapping compaction boxes on one turn. |

### Validation Rules

- `.ctxrc` schema (`internal/assets/schema/ctxrc.schema.json`) and
  `internal/rc/validate.go` accept `eocc_threshold_tokens` as a
  non-negative integer.
- `rc.EoCCThresholdTokens()` returns `0` when unset or `<= 0`, and clamps
  positive values below `DefaultMinEoCCThreshold` (`4000`) to `4000`.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| Unreadable handover file in `.context/handovers/` during `ctx agent` | Non-fatal warning logged; Tier 1.5 skipped so `ctx agent` still succeeds. | Fix or remove the corrupted file in `.context/handovers/`. |
| Hook relay failure in `EmitEoCCWarning` | Honours log-first principle: returns error without printing the nudge box or touching the guard file, retrying on the next prompt. | Inspect `.context/logs/check-context-size.log`. |

## Interface

### CLI

```bash
ctx agent [--budget 8000] [--format md|json] [--no-handover]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--no-handover` | `bool` | `false` | Exclude the latest handover from the assembled context packet. |

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/config/runtime/runtime.go` | Add `DefaultEoCCThresholdTokens = 0`, `DefaultMinEoCCThresholdTokens = 4000`, and `DefaultHandoverBudgetPct = 15`. |
| `internal/rc/types.go` | Add `EoCCThresholdTokens int` (`yaml:"eocc_threshold_tokens"`) to `CtxRC`. |
| `internal/rc/default.go` | Alias `DefaultEoCCThresholdTokens` and `DefaultMinEoCCThresholdTokens` from `internal/config/runtime`. |
| `internal/rc/rc.go` | Add `EoCCThresholdTokens() int` accessor with clamping. |
| `internal/assets/schema/ctxrc.schema.json` | Add `eocc_threshold_tokens` property definition. |
| `internal/rc/schema_test.go` & `internal/rc/rc_test.go` | Add schema parity and default/override/clamp tests for `eocc_threshold_tokens`. |
| `internal/config/hook/variant.go` | Add `VariantEoCC = "eocc"`. |
| `internal/config/event/size.go` | Add `EoCC = "eocc"` event constant. |
| `internal/config/stats/context.go` | Add `ContextSizeEoCCNudgedPrefix` constant. |
| `internal/config/embed/text/check_context_size.go` | Add `DescKeyWriteCheckContextSizeEoCC*` constants. |
| `internal/assets/commands/text/hooks.yaml` & `write.yaml` | Add EoCC fallback, log, relay, and box title text entries. |
| `internal/assets/hooks/messages/check-context-size/eocc.txt` | Add embedded EoCC hook message template. |
| `internal/assets/hooks/messages/registry.yaml` | Register `eocc` template variant under `check-context-size`. |
| `internal/cli/system/core/nudge/context_size.go` | Add `EmitEoCCWarning` builder and relay helper. |
| `internal/cli/system/cmd/checkcontextsize/run.go` | Evaluate `rc.EoCCThresholdTokens()` against `tokens` and emit the EoCC nudge when crossed. |
| `internal/config/flag/flag.go` | Add `NoHandover = "no-handover"` flag constant. |
| `internal/config/embed/flag/agent.go` | Add `DescKeyAgentNoHandover = "agent.no-handover"`. |
| `internal/assets/commands/flags.yaml` | Add `agent.no-handover` flag description. |
| `internal/write/handover/latest.go` | Add `LatestFile(handoversDir string) (File, bool, error)` to return the parsed latest handover file and body. |
| `internal/cli/agent/cmd/root/cmd.go` | Bind `--no-handover` bool flag and pass to `Run`. |
| `internal/cli/agent/cmd/root/run.go` | Load the latest handover body when `!noHandover` and pass into `OutputAgentMarkdown` / `OutputAgentJSON`. |
| `internal/cli/agent/core/budget/types.go` | Add `Handover string` (`json:"handover,omitempty"`) to `packet` and `AssembledPacket`. |
| `internal/cli/agent/core/budget/out.go` | Forward `handoverBody` to `AssemblePacket`. |
| `internal/cli/agent/core/budget/assemble.go` | Add Tier 1.5 handover budget allocation and fitting. |
| `internal/cli/agent/core/budget/render.go` | Render `## Latest Handover` section in Markdown output when `pkt.Handover != ""`. |
| `internal/assets/integrations/opencode/plugin/index.ts` | Replace `ctx system bootstrap` in `experimental.session.compacting` with `ctx agent --budget 4000`. |
| `docs/home/configuration.md` & `docs/cli/agent.md` | Document `eocc_threshold_tokens` and `--no-handover`. |

### Key Functions

```go
// EoCCThresholdTokens returns the configured EoCC compaction token
// threshold, or 0 when disabled. Positive values below 4000 clamp to 4000.
func EoCCThresholdTokens() int

// EmitEoCCWarning builds and relays the EoCC compaction nudge box.
func EmitEoCCWarning(
	logFile, sessionID string,
	count, tokens, threshold int,
) (string, error)
```

### Helpers to Reuse

- `handover.LatestFile` (`internal/write/handover/latest.go`) to locate and
  parse the newest handover file in `.context/handovers/`.
- `ctxToken.EstimateString` (`internal/context/token`) for Tier 1.5 budget
  accounting.
- `message.Load` and `message.NudgeBox`
  (`internal/cli/system/core/message`) for the EoCC hook template.

## Configuration

```yaml
# .ctxrc
eocc_threshold_tokens: 30000
```

## Testing

- **Unit (`internal/rc/rc_test.go`):** Verify `EoCCThresholdTokens()`
  returns `0` by default, clamps `1000` to `4000`, and returns `30000`
  when configured.
- **Unit (`internal/cli/agent/core/budget/budget_test.go`):** Verify
  `AssemblePacket` includes the newest handover in Tier 1.5, respects the
  15% cap, omits it when `--no-handover` is set, and behaves identically
  when `.context/handovers/` is empty.
- **Unit (`internal/cli/system/core/nudge/trigger_test.go`):** Verify EoCC
  threshold firing, wrap-up suppression, and precedence over the 60%
  checkpoint nudge.

## Non-Goals

- Automatically terminating or clearing an interactive host conversation
  from inside a hook subprocess (host CLIs such as Claude Code and Codex
  do not expose a hook return code that clears conversation history
  in-place).
- Removing the existing 60% checkpoint, 90% window warning, or
  `billing_token_warn` triggers.
