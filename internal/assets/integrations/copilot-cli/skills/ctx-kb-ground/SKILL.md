---
name: ctx-kb-ground
description: "External grounding pass for the kb. Reads grounding-sources.md, refreshes the listed sources, and advances source-coverage ledger rows accordingly. Writes a ground closeout for the audit trail. Prompts once if grounding-sources.md is empty; NONE on a line is a per-pass skip."
tools: [bash]
---

# Ground the KB Against External Sources

Refresh `.context/kb/` against the external sources the user has
declared in `.context/ingest/grounding-sources.md`. This is the
*"are we still current?"* pass. It does not mint new evidence by
itself, does not author topic pages, does not modify Confidence
bands. It walks the declared external sources, checks each for
drift against what the kb cites, and advances the source-coverage
ledger to reflect what the refresh found.

If the refresh surfaces material the kb should absorb, this skill
flags it and recommends `/ctx-kb-ingest`. Authoring is ingest's
authority, not this skill's.

Authoritative background reading:
`.context/ingest/KB-RULES.md` §Authority boundary and
§Source-coverage ledger.

## When to Use

- The user says "re-ground the kb", "check upstream",
  "are the docs still current?", or invokes the explicit slash
  form.
- Before a release / handover where source freshness matters.
- After an external vendor has shipped a version bump.
- Periodically (per the user's cadence) as kb hygiene.

## When NOT to Use

- The user has new materials in hand (use `/ctx-kb-ingest`).
- The user is asking a content question (use `/ctx-kb-ask`).
- The user wants a structural audit (use `/ctx-kb-site-review`).

## Authority Boundary (vs Other Skills)

- **`/ctx-kb-ground`**: refreshes external sources listed in
  `grounding-sources.md`; advances source-coverage ledger rows
  for sources where the refresh changed state; writes a ground
  closeout. **May not** mint `EV-###` rows, author prose,
  modify a topic page, or change a Confidence band.
- **`/ctx-kb-ingest`**: handles anything this skill surfaces
  as new material to absorb.
- **`/ctx-kb-ask`**: handles read-only questions about kb
  content.
- **`/ctx-kb-site-review`**: handles structural audit (separate
  surface from source-freshness audit).

## Usage Examples

```text
/ctx-kb-ground
```

No arguments. Sources come from
`.context/ingest/grounding-sources.md`.

## Input Contract

The file `.context/ingest/grounding-sources.md` is the sole
declaration surface. Each non-empty, non-comment line names a
source (URL, in-tree path, MCP resource identifier) the user
wants this skill to track. A line whose value is the literal
`NONE` is a **per-pass skip**: this invocation does nothing and
re-prompts on the next invocation. Lines beginning with `#` are
comments.

There is no CLI argument for sources. To configure what this
skill checks, edit `grounding-sources.md`.

## Pre-Write Gates

Three distinct refusals, each leaves zero residue:

- `.context/` missing → suggest `ctx init` and stop.
- `.context/ingest/` missing → suggest `ctx init --upgrade`
  and stop.
- Kb scope undeclared → refuse with the scope message and stop.

## Refuse-on-Empty

`.context/ingest/grounding-sources.md` may be in three states:

1. **Missing or empty**: file does not exist, or has only
   comments and blank lines. Prompt once:

   > `grounding-sources.md` has no sources. List one source per
   > line (URL, in-tree path, MCP resource). `NONE` on a line
   > is a per-pass skip and re-prompts next invocation.

   Stop. Do not synthesize a list. Do not invent sources from
   the kb's `source-map.md` (that file's authority is ingest;
   grounding's declaration surface is separate by design).

2. **Single line `NONE`**: per-pass skip. Write no closeout;
   return exactly:

   > grounding-sources.md is `NONE` for this pass; skipping.
   > Edit `.context/ingest/grounding-sources.md` to set actual
   > sources, or leave `NONE` to keep skipping.

   The next invocation re-prompts as in (1) above.

3. **One or more sources listed**: proceed to Process.

The empty-and-prompt path is the one exception to the
refuse-on-empty pattern other mode skills enforce; the rationale
is that grounding's declaration lives in a file the user owns
(not in a slash argument), so a one-shot prompt is cheaper than
forcing them to remember the filename.

## Process

1. **Verify pre-write gates.** Refuse cleanly if any gate fails.
   Zero residue on refusal.

2. **Read `.context/ingest/grounding-sources.md`.** Handle the
   three states per §Refuse-on-empty.

3. **For each declared source**, in order of appearance:

   - **Resolve** the source: fetch the URL, stat the in-tree
     path, enumerate the MCP resource.
   - **Cross-reference** against `.context/kb/source-map.md` to
     find the kb's short-name for this source (if any). If
     absent, the source is *new to the kb*; record it as a flag
     to surface in the closeout (do not mint a `source-map.md`
     row; that is ingest's authority).
   - **Check freshness** using the strongest available signal:
     - URL: HTTP Last-Modified header, or ETag, or visible
       version stamp on the page; compare against the
       `source-map.md` row's `dated:` cell (if present).
     - In-tree path: file mtime + git SHA; compare against the
       `evidence-index.md` rows that cite the source by SHA
       (in-repo citations pin to a SHA at extraction time per
       `KB-RULES.md` §Evidence discipline).
     - MCP resource: whatever freshness primitive the resource
       exposes; if none, treat as opaque (record as
       *"freshness opaque"* in the closeout).
   - **Classify the refresh outcome** as one of:
     - **`unchanged`**: source has not drifted since the
       kb's last extraction; no ledger update needed.
     - **`drifted`**: source has changed; the kb's claims
       citing this source may be stale; advance the ledger row
       to a state that reflects the staleness:
       - If the row was `comprehensive`, advance to a typed
         `superseded-pending` annotation in the `Residue` cell
         (do not write a new state name; the state machine in
         `KB-RULES.md` is closed; `Residue` is the
         human-readable annotation surface).
       - If the row was anywhere prior to `comprehensive`,
         leave the state and add a `drifted` note in `Residue`
         + `Next action` set to the explicit
         `/ctx-kb-ingest <slug>` resumption.
     - **`gone`**: source returns 404, file deleted, MCP
       resource removed; flag for the user. The right
       resolution may be `superseded` (with a named successor)
       or `skipped` (out of scope); that judgment is the
       user's, not this skill's.
     - **`freshness opaque`**: no freshness signal available;
       record in `Residue` cell as *"freshness opaque
       (<date checked>)"*; no ledger state change.
   - **Advance the ledger row** only for `drifted` and `gone`
     cases, and only via `Residue` / `Next action` annotation
     (not state change). State transitions out of `comprehensive`
     are ingest's authority.

4. **Write the ground closeout.** Create
   `.context/ingest/closeouts/<TIMESTAMP>-ground-closeout.md`
   with required frontmatter:

   ```yaml
   ---
   sha: <short>
   branch: <name>
   mode: ground
   pass-mode: external-refresh
   life-stage: <bootstrap|maintenance>
   generated-at: <RFC-3339>
   ---
   ```

   Body sections:
   - **Inputs**: declared sources from
     `grounding-sources.md`, count + one bullet each.
   - **Refresh outcomes**: for each source: `unchanged`,
     `drifted`, `gone`, `freshness opaque`, or `new to kb`.
     Cite the kb short-name (or *"new to kb"*) and the
     evidence used to classify (Last-Modified header, version
     stamp, file mtime).
   - **Ledger updates**: every `Residue` / `Next action`
     change applied to a `source-coverage.md` row, with the
     before/after annotation.
   - **Flags**: sources the refresh found `gone`, sources
     classified `new to kb`, sources with conflicting
     freshness signals. Each flag names the source and the
     recommended next pipeline step.
   - **Next pass hint**: explicit invocations to absorb
     drifted / new material (e.g. *"`/ctx-kb-ingest <slug>` to
     refresh `cursor/hooks` against the v1.2 docs"*).

## Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `grounding-sources.md` missing or empty (only comments/blank) | Prompt once with the standard text; stop. No closeout. |
| `grounding-sources.md` single line `NONE` | Skip this pass with the standard skip text; stop. No closeout. |
| `.context/` missing | Refuse; suggest `ctx init`. No residue. |
| `.context/ingest/` missing | Refuse; suggest `ctx init --upgrade`. No residue. |
| Kb scope undeclared | Refuse with the scope message. No residue. |
| Source returns 404 / file deleted / MCP resource removed | Classify `gone`; flag; recommend the user choose `superseded` (with successor) or `skipped` (out of scope). Do not auto-transition. |
| Source unchanged since last extraction | Record `unchanged` in closeout's `Refresh outcomes`; no ledger update. |
| Source drifted since last extraction (URL bumped, file mtime newer than cited SHA) | Record `drifted`; annotate ledger row's `Residue` / `Next action`; recommend `/ctx-kb-ingest <slug>`. Do not modify topic-page prose. |
| Source has no freshness primitive (opaque) | Record `freshness opaque (<date checked>)`; no ledger state change; surface in `Flags` so the user can decide cadence. |
| Source listed in `grounding-sources.md` but not in `source-map.md` (new to kb) | Classify `new to kb`; flag; recommend `/ctx-kb-ingest <source>` to admit. Do not mint a `source-map.md` row from this skill. |
| Source listed in `grounding-sources.md` but URL malformed or path nonexistent | Surface as a per-source error in the closeout's `Flags`; continue with remaining sources; do not abort the pass. |
| Source's `source-map.md` row has `dated:` but `evidence-index.md` rows lack `occurred:` | Flag (temporal-precedence rule needs it); recommend hand-edit. Do not auto-edit. |
| User added a new source to `grounding-sources.md` since the last pass | Treated as a regular declared source; classified per the freshness check; no special path. |
| Mid-pass MCP fetch failure | Record per-source error; continue; do not abort the whole pass. |

## Anti-Patterns

- Minting `EV-###` rows from this skill. Evidence authoring is
  ingest's authority.
- Authoring topic-page prose from this skill. Page authoring is
  ingest's authority.
- Modifying a claim's Confidence band from this skill. Demotion
  is evidence work.
- Auto-transitioning a `comprehensive` ledger row out of
  `comprehensive`. State changes require ingest judgment; this
  skill annotates `Residue` / `Next action` only.
- Synthesising a source list when `grounding-sources.md` is
  empty. The declaration surface is the file the user owns.
- Inventing a freshness signal when none exists. *"Freshness
  opaque"* is the honest classification.
- Skipping the closeout once pre-write gates pass and at least
  one source was processed.

## Output Contract

For pre-write refusals, return only the specified refusal text
and stop. No residue.

For empty / `NONE` cases, return the matching prompt or skip
text and stop. No closeout in those cases.

For passes that processed at least one source, end with this
structured summary:

- **Sources checked**: count + one bullet each, classified
  (`unchanged | drifted | gone | freshness opaque | new to kb`).
- **Ledger updates**: count + one-line categories.
- **Flags**: count + categories.
- **Closeout**: filename on its own line.
- **Next-recommended-action**: explicit invocations to absorb
  drifted / new material (or `none` if every source was
  `unchanged`).

## Quality Checklist

Before reporting completion, verify:

- [ ] Pre-write gates passed (or the matching refusal was
      returned with zero residue).
- [ ] Every declared source from `grounding-sources.md` was
      checked, classified, and recorded in `Refresh outcomes`.
- [ ] No `EV-###` row was minted, no topic-page prose was
      written, no Confidence band was changed, no ledger state
      was transitioned (only `Residue` / `Next action`
      annotated).
- [ ] Every `drifted` / `gone` / `new to kb` source has an
      explicit `Next-recommended-action`.
- [ ] Closeout written with all required frontmatter fields.
