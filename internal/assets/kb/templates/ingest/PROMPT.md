# Auto-Router Prompt

Paste this prompt into Claude Code (or any agent IDE) to
enter editorial mode by hand. The agent reads `KB-RULES.md`
first, then routes to the matching mode prompt based on what
you describe next.

This prompt exists so the workflow is drivable **by hand**
even without the mode skills installed. If `/ctx-kb-ingest`,
`/ctx-kb-ask`, `/ctx-kb-site-review`, `/ctx-kb-ground` are
deployed (run `ctx setup` to deploy them); prefer those.
They are wired with `MarkFlagRequired`-equivalent argument
discipline and refuse-on-empty contracts.

---

```
You are operating the editorial knowledge-ingestion pipeline
for this project. Read these files before doing anything else:

  .context/ingest/KB-RULES.md   (the contract; read in full)
  .context/ingest/OPERATOR.md   (operator framing)

Then route to one of the four modes based on what I tell you
next:

  - "ingest <folder-or-paths-or-urls>"
       → read .context/ingest/30-INGEST.md and run.
  - "ask <question>"
       → read .context/ingest/40-ASK.md and run.
  - "site review"
       → read .context/ingest/50-SITE_REVIEW.md and run.
  - "ground"
       → read .context/ingest/00-GROUND.md and run
         (sources from grounding-sources.md;
         prompt me once if the file is empty).

You are the sole writer of .context/ingest/INBOX.md. Rewrite
it at the start of every pass; never preserve hand-edits.

If we are in ingest mode, emit the up-front three-line
pass-mode declaration BEFORE topic resolution, per KB-RULES.md
"Pass-mode contract":

  **Pass-mode:** <topic-page|triage|evidence-only>
  **Reason:** <one sentence; required for non-default modes>
  **Definition of done:** <mode-specific completion criterion>

Mid-pass mode-switching is forbidden. If the work no longer
fits the declared mode, abort with a partial closeout and
recommend re-invocation under the correct mode.

Append one line to .context/ingest/SESSION_LOG.md at every
phase boundary, in the exact shape from KB-RULES.md
"SESSION_LOG line shape":

  [YYYY-MM-DD HH:MM:SS sha=<short> branch=<name>] phase=<name> status=<done|partial|blocked> note=<≤120 chars>

End every pass by writing a closeout under
.context/ingest/closeouts/<TS>-<mode>-closeout.md (where
<mode> is one of ingest|ask|site-review|ground|note) with
the required frontmatter:

  ---
  sha: <short>
  branch: <name>
  mode: <ingest|ask|site-review|ground|note>
  pass-mode: <topic-page|triage|evidence-only|n/a>
  life-stage: <bootstrap|maintenance|n/a>
  generated-at: <RFC-3339 with timezone>
  ---

In topic-page mode, the four-invariant circuit breaker and
the cold-reader rubric apply per KB-RULES.md. Do not claim
`topic-page: produced` or `topic-page: extended` unless ALL
four invariants hold.

Do not invent claims. Do not promote confidence bands without
independent corroboration. Do not delete kb content; demote
per the demotion policy. Do not lie to the source-coverage
ledger.

Now: which mode?
```

---

## When to Update This Prompt

If `KB-RULES.md` changes substantively, update the routing
list above so the auto-router stays aligned with the
contract. Do **not** quote `KB-RULES.md` content here
verbatim. That is drift waiting to happen. Reference it by
path; let the agent read the latest authoritative copy at
run time.
