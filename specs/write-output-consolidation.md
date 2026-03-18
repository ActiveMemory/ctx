# Spec: Write Package Output Consolidation

Consolidate multi-line imperative `cmd.Println` sequences in `internal/write/`
into pre-computed single-print patterns where the output shape is clearer and
conditionals are separated from I/O.

## Problem

Functions in `internal/write/` mix control flow with output. A typical
pattern looks like:

```go
cmd.Println(fmt.Sprintf(config.TplLoopGenerated, outputFile))
cmd.Println()
cmd.Println(heading)
cmd.Println(fmt.Sprintf(config.TplLoopRunCmd, outputFile))
cmd.Println()
cmd.Println(fmt.Sprintf(config.TplLoopTool, tool))
cmd.Println(fmt.Sprintf(config.TplLoopPrompt, promptFile))
if maxIterations > 0 {
    cmd.Println(fmt.Sprintf(config.TplLoopMaxIterations, maxIterations))
} else {
    cmd.Println(config.TplLoopUnlimited)
}
cmd.Println(fmt.Sprintf(config.TplLoopCompletion, completionMsg))
```

Problems with this pattern:

1. **Output shape is invisible** — you can't see the final block without
   mentally executing the function line by line
2. **Conditionals interleaved with I/O** — business logic and presentation
   are tangled
3. **Many small Tpl\* constants** — each line gets its own YAML key and
   config variable, fragmenting a single output block across 3 layers
   (YAML → config → write function)
4. **Testing is harder** — verifying output requires mocking `cmd` and
   asserting individual lines rather than checking one string

## Analysis Summary

Full audit of `internal/write/` (26 Go files, ~160 functions, 337
`cmd.Println`/`cmd.Print` calls):

| Category | Count | Description |
|----------|------:|-------------|
| Trivial | 27 | 1 Println, no conditionals — already optimal |
| Simple | 42 | 2-3 Printlns, no conditionals — marginal gain |
| Multi-line | 16 | 4+ Printlns, no conditionals — good candidates |
| Conditional | 22 | if/else interleaved with output — needs pre-computation |
| Complex | 18 | Loops, nested conditionals — leave as-is |

**Target population:** ~38 functions (Multi-line + Conditional).

## Solution

### Pattern: Pre-compute, then print

Separate conditional logic from output. Build the full output string first,
then emit it in one call.

**Before:**

```go
func InfoLoopGenerated(
    cmd *cobra.Command,
    outputFile, heading, tool, promptFile string,
    maxIterations int,
    completionMsg string,
) {
    cmd.Println(fmt.Sprintf(config.TplLoopGenerated, outputFile))
    cmd.Println()
    cmd.Println(heading)
    cmd.Println(fmt.Sprintf(config.TplLoopRunCmd, outputFile))
    cmd.Println()
    cmd.Println(fmt.Sprintf(config.TplLoopTool, tool))
    cmd.Println(fmt.Sprintf(config.TplLoopPrompt, promptFile))
    if maxIterations > 0 {
        cmd.Println(fmt.Sprintf(config.TplLoopMaxIterations, maxIterations))
    } else {
        cmd.Println(config.TplLoopUnlimited)
    }
    cmd.Println(fmt.Sprintf(config.TplLoopCompletion, completionMsg))
}
```

**After:**

```go
func InfoLoopGenerated(
    cmd *cobra.Command,
    outputFile, heading, tool, promptFile string,
    maxIterations int,
    completionMsg string,
) {
    iterLine := config.TplLoopUnlimited
    if maxIterations > 0 {
        iterLine = fmt.Sprintf(config.TplLoopMaxIterations, maxIterations)
    }
    cmd.Println(fmt.Sprintf(config.TplLoopBlock,
        outputFile, heading, outputFile, tool, promptFile, iterLine, completionMsg,
    ))
}
```

Where `TplLoopBlock` is a single multiline YAML entry:

```yaml
write.loop-block:
  short: |
    ✓ Generated %s

    %s
      ./%s

    Tool: %s
    Prompt: %s
    %s
    Completion signal: %s
```

### When to consolidate

A function is a candidate when **all** of these hold:

1. It has 4+ `cmd.Println` calls
2. The output represents a single logical block (not independent messages)
3. Conditionals can be pre-computed into a string before the print call
4. No loops over dynamic-length collections

### When NOT to consolidate

- **Trivial/Simple functions** (1-3 Printlns) — the overhead of a block
  template exceeds the clarity gain
- **Loop-based functions** (`BootstrapText`, `StatusActivity`,
  `LoadAssembled`) — iteration over dynamic arrays requires imperative code
- **Functions with many independent conditionals** (`PublishPlan` with 4
  optional sections) — pre-computing 4 optional blocks into one format
  string is less readable than the current approach

### Template constants migration

For each consolidated function:

1. **Add** one `TplXxxBlock` multiline YAML entry in the appropriate
   `text/*.yaml` file
2. **Add** one `TextDescKeyWriteXxxBlock` constant in `embed.go`
3. **Add** one `config.TplXxxBlock` variable in `config/config.go`
4. **Remove** the individual `TplXxx*` line constants that the block replaces
5. **Remove** corresponding `TextDescKey*` constants and YAML entries

### Impact on existing Tpl\* constants

Each consolidated function eliminates 4-8 individual constants. For ~20
realistic candidates, that's 80-160 fewer constants across the YAML → config
→ write chain.

## Scope

### In scope

- Functions in `internal/write/` with 4+ Printlns and pre-computable
  conditionals (~20 realistic candidates from the 38 target population)
- Corresponding YAML, embed.go, and config/config.go constant cleanup
- Existing tests updated to match new output format

### Out of scope

- Migrating to `text/template` — overkill for simple format strings with
  at most one conditional
- Changing the write package's public API signatures — callers stay the same
- Consolidating trivial/simple functions (1-3 Printlns)
- Refactoring loop-based complex functions

## Candidate Functions

### Tier 1: Multi-line, no conditionals (straightforward)

| Function | File | Printlns | Constants eliminated |
|----------|------|:--------:|:--------------------:|
| `InfoInitNextSteps` | info.go | 5 | 3 |
| `InfoObsidianGenerated` | info.go | 4 | 2 |
| `InfoJournalSiteGenerated` | info.go | 6 | 4 |
| `InfoDepsNoProject` | info.go | 3 | 3 |
| `ArchiveDryRun` | task.go | 5 | ~4 |
| `ImportScanHeader` | import.go | 3 | ~2 |

### Tier 2: Conditional, pre-computable (moderate)

| Function | File | Printlns | Conditionals | Pre-computation |
|----------|------|:--------:|:------------:|-----------------|
| `InfoLoopGenerated` | info.go | 7 | 1 (maxIterations) | iterLine string |
| `SyncResult` | sync.go | 8 | 3 (mirror, source, new) | 3 pre-computed lines |
| `CtxSyncHeader` | sync/ctxsync.go | 5 | 1 (dryRun) | dryRunLine string |
| `CtxSyncAction` | sync/ctxsync.go | 3 | 1 (suggestion) | suggestionLine string |
| `PruneSummary` | prune.go | 4 | 1 (dryRun) | suffix string |
| `SessionMetadata` | recall.go | 9+ | 2 (branch, model) | 2 optional lines |
| `TestResult` | notify.go | 2 | 1 (statusCode) | suffix string |
| `SyncDryRun` | sync.go | 4 | 1 (hasDrift) | driftLine string |

### Tier 3: Skip (complex or low-value)

- `BootstrapText` — loops over dynamic arrays
- `StatusActivity` — loop iteration
- `LoadAssembled` — complex token budgeting loop
- `PublishPlan` — 4 independent optional sections
- `ImportSummary` — 4 conditional counts + dryRun variant
- `RestoreDiff` — helper function calls within output
- All trivial/simple functions

## Sequencing

This work fits naturally as an extension of Phase WC (Write Consolidation).
Tier 1 (no conditionals) should be done first as a proof of the pattern
before tackling Tier 2.

## Non-Goals

- Full templating engine (`text/template`, `html/template`)
- Changing function signatures or caller code
- Consolidating output across different functions
- i18n beyond what the existing YAML asset system provides

## Risks

- **Multiline YAML format strings** are harder to diff in PRs — mitigated
  by keeping each block under ~10 lines
- **Format verb count must match args** — `fmt.Sprintf` catches this at
  runtime, not compile time. Mitigated by existing test patterns that
  exercise each function.
