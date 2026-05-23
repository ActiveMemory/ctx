# ctx-remember: Silent Fallback for Absent Companion Tools

`/ctx-remember`'s Companion Tool Check tells the agent to
surface an install hint to the user when a companion tool
(Gemini Search, GitNexus) is absent. This contradicts the
peer-MCP / not-vouched-for-by-ctx position recorded in
`DECISIONS.md` on 2026-05-23. Fix: silently fall back to
built-in capabilities; emit nothing for absent tools.

## Problem

`internal/assets/claude/skills/ctx-remember/SKILL.md`
(lines 142-150, with identical copies in the copilot-cli
and opencode skill variants) instructs the agent:

```
**Check procedure:**

1. Attempt each smoke test silently
2. For tools that respond: note as available (no output needed)
3. For tools that fail or are not connected: append a brief note
   after the readback:
   > "Companion tools: Gemini Search is not connected (web search
   > will fall back to built-in). Install via MCP settings if
   > needed."
4. For GitNexus specifically: if it responds but the current repo
   is not indexed or the index is stale, suggest:
   > "GitNexus index is stale: run `npx gitnexus analyze` to
   > rehydrate."
```

Step 3 surfaces "Install via MCP settings if needed" — a
soft pointer at an install path for a third-party tool ctx
doesn't ship. That's the same coupling pattern today's
DECISIONS.md entry rules out:

> An MCP gateway through ctx would couple ctx to the
> lifecycle of every gatewayed tool. […] That coupling
> is a tax we don't want to pay for a tool we don't ship.

A skill text that tells the user to install a specific
companion tool is a *softer* form of the same coupling:
ctx is implicitly vouching for the install path, the
version compatibility, and (when the install fails) the
support burden. The graceful fallback in steps 1-2 already
handles the functional case; step 3 adds noise that
contradicts the architectural stance.

Step 4 is different: it triggers only when the tool IS
connected (the smoke test passed) and addresses a
working-tool's state. Refreshing a present tool is
operationally part of using it, not vouching for its
install. Step 4 stays.

## Solution

Two identical edits, one per skill copy that contains
the Companion Tool Check section:

1. `internal/assets/claude/skills/ctx-remember/SKILL.md`
2. `internal/assets/integrations/copilot-cli/skills/ctx-remember/SKILL.md`

(The opencode variant at
`internal/assets/integrations/opencode/skills/ctx-remember/SKILL.md`
is a shorter form that doesn't include the Companion Tool
Check section at all — nothing to fix there.)

Each block becomes:

```
**Check procedure:**

1. Attempt each smoke test silently
2. For tools that respond: note as available (no output needed)
3. For tools that fail or are not connected: silently fall
   back to built-in capabilities. Emit no output. ctx does
   not vouch for companion-tool install paths (see
   DECISIONS.md, 2026-05-23 "MCP gateway not worth the
   coupling cost").
4. For GitNexus specifically: if it responds but the current
   repo is not indexed or the index is stale, suggest:
   > "GitNexus index is stale: run `npx gitnexus analyze`
   > to rehydrate."

Present companion status as a one-line note after the readback
only when there's something actionable (stale index). Absent
tools produce no output; the agent uses its built-in
capabilities transparently.
```

The closing paragraph also updates: "If everything is
healthy, say nothing" loses the implicit "say something if
anything is absent" reading.

## Behavioral consequence

Operators who currently see "Companion tools: Gemini
Search is not connected (web search will fall back to
built-in). Install via MCP settings if needed." after
each `/ctx-remember` invocation will stop seeing it. Their
sessions get quieter without losing capability — the
fallback to built-in search was already happening; only
the suggestion to install goes away.

Operators who *want* the hint can re-enable companion
checks broadly via `companion_check: true` in `.ctxrc`
(the flag controls the entire section). The default of
"hint absent → silent" is the new contract.

## Out of Scope

- **Tool-agnostic skill language.** The skill still names
  Gemini Search and GitNexus as the canonical companions.
  Making skill language tool-agnostic (so a user with
  Firecrawl + vLLM gets first-class treatment) is the
  broader "skills hard-code specific tools" smell flagged
  by the user — out of scope here, addressed via a
  separate design pass.
- **Doctor / preflight changes.** `ctx doctor` and similar
  diagnostic surfaces remain free to report companion-tool
  availability; this spec is about the per-invocation
  in-skill nag, not the diagnostic command surface.
- **`companion_check` semantics.** The flag still controls
  whether the entire Companion Tool Check section runs.
  No change to its default or meaning.

## Verification

- `grep -r "Install via MCP" internal/assets/` returns
  zero matches.
- The three SKILL.md files have identical Check Procedure
  blocks, all with the silent-fallback wording.
- `make lint` clean; no Go code changes.
- `go test ./...` clean (no test depends on the install-nag
  text).
