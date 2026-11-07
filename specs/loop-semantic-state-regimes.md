# Loop Semantic State Regimes (EoCC and SAC)

## Problem

`ctx loop` (`internal/cli/loop/core/script/script.go`) generates a bash
script (`loop.sh`) that runs an AI CLI tool in a `while true` loop until a
completion signal (`SYSTEM_CONVERGED`, `DONE`, etc.) appears in the output.
Because each loop iteration invokes a fresh process (for example,
`claude --print "$(cat PROMPT.md)"`), `ctx loop` already resets the active
context window between iterations.

However, compared to the End-of-Conversation Compaction (EoCC) and
Single-Action Compaction (SAC) harnesses in `arXiv:2611.05306` (*You don't
need long context*, Sections 2.4.1, 3, and 5), `ctx loop` has three
structural gaps:

1. **No semantic state carry-over ($S_t \to S_{t+1}$):** Each iteration
   receives only the static `PROMPT.md` contents. It does not automatically
   inject the assembled `ctx agent` packet or the latest forward-looking
   handover written by the preceding iteration.
2. **No compaction contract at iteration end:** The loop prompt does not
   instruct the agent to write an updated execution brief
   (`ctx handover write --no-fold`) before exiting the iteration, forcing
   the next iteration to re-discover state from raw files and `git log`.
3. **No turn/action cap per iteration:** In `arXiv:2611.05306`, EoCC bounds
   each stretch by a context/turn limit, and SAC restricts each conversation
   to a single environment action followed by a semantic state update,
   reducing median peak context from ~30K tokens (200K max) to 7-9K tokens
   while maintaining SWE-bench and Terminal-Bench accuracy.

## Approach

Extend `ctx loop` with explicit execution regimes (`--regime` / `-r`) and
per-iteration turn bounds (`--max-turns`):

1. **`--regime <mode>` (`-r`):**
   - `standard` (default): Preserves current `ctx loop` output byte-for-byte
     for backward compatibility with existing golden tests and scripts.
   - `eocc`: Configures the loop for End-of-Conversation Compaction. Each
     iteration receives a composite prompt combining `PROMPT.md`,
     `ctx agent --budget 4000` (which includes the latest handover state),
     and the EoCC handover requirement: before finishing the iteration (or
     when reaching the turn cap), the agent must persist its updated
     forward-looking execution brief via `ctx handover write --no-fold`.
   - `sac`: Configures the loop for Single-Action Compaction (SSSM cycle).
     Each iteration instructs the agent to: (1) read the injected semantic
     state, (2) execute at most one environment action (tool call / edit /
     test command), (3) observe the outcome, and (4) immediately persist
     the updated semantic state via `ctx handover write --no-fold` and exit
     the iteration (emitting the completion signal only when the task's
     completion criteria are fully verified).
2. **`--max-turns <n>`:**
   - Optional integer flag specifying the maximum turns per loop iteration.
   - When `--regime sac` is selected and `--max-turns` is omitted, defaults
     to `2` (1 action turn + 1 state-update turn) for tools that support
     CLI turn caps (such as `claude --print --max-turns 2`), in addition to
     the prompt-level single-action rule.

## Behavior

### Happy Path

1. The user runs:
   ```bash
   ctx loop --prompt PROMPT.md --tool claude --regime sac --max-iterations 25
   ```
2. `ctx loop` validates the regime and generates `loop.sh`.
3. On each iteration of `loop.sh`:
   - The script captures `STATE_PACKET=$(ctx agent --budget 4000 2>/dev/null || true)`.
   - The script constructs the iteration input from `PROMPT.md`,
     `$STATE_PACKET`, and the regime's state-update instructions.
   - The script invokes `claude --print --max-turns 2 "$ITER_PROMPT"`.
   - The agent performs one action, updates the handover via
     `ctx handover write --no-fold`, and exits.
   - If the output contains `SYSTEM_CONVERGED`, the loop notifies via
     `ctx hook notify --event loop` and exits 0; otherwise, the next
     iteration starts with a fresh 7-9K token context containing the newly
     written handover state.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `--regime` is omitted or set to `"standard"` | Generated script is byte-identical to existing `ctx loop` output (all existing golden tests pass unchanged). |
| Unknown `--regime` value (e.g., `--regime foo`) | Command fails with a validation error listing valid regimes (`standard`, `eocc`, `sac`). Zero files written. |
| `--tool aider` or `--tool generic` used with `--regime sac` and `--max-turns` | Because `aider` and `generic` do not share Claude's `--max-turns` CLI flag, the single-action constraint is enforced through the injected regime instructions in the prompt envelope without appending an unsupported CLI flag. |
| `ctx agent` returns non-zero inside the loop script (e.g., temporary lock) | The script's `|| true` fallback keeps the loop running with `PROMPT.md` and the regime instructions. |
| Negative `--max-turns` passed | Rejected at CLI validation with an error before generating the script. |

### Validation Rules

- `--regime` must be one of `cfgLoop.RegimeStandard` (`"standard"`),
  `cfgLoop.RegimeEoCC` (`"eocc"`), or `cfgLoop.RegimeSAC` (`"sac"`).
- `--max-turns` must be `>= 0`.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| Invalid `--regime` value | `Error: invalid loop regime "foo": must be one of standard, eocc, sac` | Pass `--regime standard`, `--regime eocc`, or `--regime sac`. |
| Negative `--max-turns` | `Error: --max-turns must be non-negative` | Pass `0` (default) or a positive integer. |

## Interface

### CLI

```bash
ctx loop [--prompt PROMPT.md] [--tool claude|aider|generic] \
  [--regime standard|eocc|sac] [--max-turns <n>] \
  [--max-iterations <n>] [--completion <signal>] [--output loop.sh]
```

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--regime` | `-r` | `string` | `"standard"` | Context regime: `standard`, `eocc`, or `sac`. |
| `--max-turns` | | `int` | `0` (`2` for `sac` on `claude`) | Maximum agent turns per loop iteration (`0` for tool default). |

### Skill

```text
/ctx-loop
```

Update `internal/assets/claude/skills/ctx-loop/SKILL.md` (and synced tool
skill copies) to document `--regime eocc`, `--regime sac`, and
`--max-turns`, explaining when to choose `eocc` (multi-step stretches with
compaction between stretches) versus `sac` (strict single-action SSSM
steps with minimal 7-9K peak context).

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/config/loop/loop.go` | Add constants `RegimeStandard = "standard"`, `RegimeEoCC = "eocc"`, `RegimeSAC = "sac"`, and `DefaultSACMaxTurns = 2`. |
| `internal/config/flag/flag.go` | Add `Regime = "regime"`, `ShortRegime = "r"`, and `MaxTurns = "max-turns"`. |
| `internal/config/embed/flag/loop.go` | Add description keys `DescKeyLoopRegime` and `DescKeyLoopMaxTurns`. |
| `internal/assets/commands/flags.yaml` | Add flag descriptions for `loop.regime` and `loop.max-turns`. |
| `internal/config/embed/text/err_config.go` | Add `DescKeyErrConfigInvalidRegime` and `DescKeyErrConfigNegativeMaxTurns`. |
| `internal/assets/commands/text/errors.yaml` | Add error templates for invalid loop regime and negative `--max-turns`. |
| `internal/err/config/config.go` | Add error constructors `InvalidRegime` and `NegativeMaxTurns`. |
| `internal/assets/tpl/types.go` | Extend `tpl.LoopData` with `Regime`, `IsEoCC`, `IsSAC`, `MaxTurns`, and `SupportsMaxTurns`. |
| `internal/cli/loop/core/script/types.go` | Define `Opts` struct for loop script generation. |
| `internal/cli/loop/cmd/root/cmd.go` | Bind `--regime` (`-r`) and `--max-turns` flags and pass to `Run`. |
| `internal/cli/loop/cmd/root/run.go` | Validate `regime` and `maxTurns` and forward `script.Opts` to `script.Generate`. |
| `internal/cli/loop/core/script/script.go` | Update `Generate(opts Opts)` to populate `tpl.LoopData` with regime state-packet assembly and `--max-turns` CLI wiring. |
| `internal/assets/tpl/templates/loop-script.sh.tmpl` | Conditionally render the `ctx agent --budget 4000` state-packet injection and EoCC/SAC execution-brief instructions when `Regime` is `eocc` or `sac`, keeping `standard` output byte-identical. |
| `internal/cli/loop/core/script/script_test.go` | Keep existing golden tests for `standard` and add golden/unit tests for `eocc` and `sac` regimes. |
| `internal/assets/claude/skills/ctx-loop/SKILL.md` | Document `--regime` and `--max-turns` (and sync to `internal/assets/codex/skills/`, `internal/assets/integrations/copilot-cli/skills/`, and `internal/assets/integrations/opencode/skills/`). |
| `docs/cli/loop.md` | Document EoCC and SAC loop modes with examples. |

### Key Functions

```go
// Opts holds the configuration for generating a Ralph loop script.
type Opts struct {
	PromptFile    string
	Tool          string
	Regime        string
	MaxTurns      int
	MaxIterations int
	CompletionMsg string
}

func Generate(opts Opts) (string, error)
```

### Helpers to Reuse

- `tpl.Render` (`internal/assets/tpl/render.go`)
- `desc.Text` (`internal/assets/read/desc`)

## Configuration

No new `.ctxrc` keys required; loop parameters are explicit CLI flags baked
into the generated `loop.sh` script.

## Testing

- **Golden tests (`internal/cli/loop/core/script/script_test.go`):**
  - Existing `claude-max.golden`, `claude-nomax.golden`,
    `aider-nomax.golden`, and `generic-max.golden` remain unchanged under
    `--regime standard`.
  - New golden tests for `claude-eocc` and `claude-sac` verify `ctx agent`
    packet capture, `ctx handover write --no-fold` instructions, and
    `--max-turns` flag emission.
- **Validation tests:** Verify invalid `--regime` and negative
  `--max-turns` return expected errors.

## Non-Goals

- Turning `ctx loop` into an in-process Go subprocess supervisor; `ctx loop`
  remains a transparent, inspectable shell script generator.
- Forcing SAC single-action resets during interactive human chat sessions.
