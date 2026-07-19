# kb-and-vocabulary

## [2026-06-07-180009] KB editorial pipeline (Phase KB) design (consolidated)

**Consolidated from**: 6 entries (2026-05-10 to 2026-05-16)

- Lift the sibling clean-room project's battle-tested editorial pipeline into
  ctx as v1, paired with handover: it is field-tested under production use and
  your-project is already paying the workaround tax (N=1 lived validation); lift
  the whole shape with a non-colliding rename, not hedge-and-defer.
- Mandate git as an architectural precondition: persistent-memory is dishonest
  without an undo layer (git reflog); refuse-on-no-git rather than auto-git-init
  (ctx never modifies the filesystem outside .context/); eliminates commit:none
  dead-code branches. Breaking change in next minor.
- KB ontology is pipeline-only-writer; no /ctx-kb-decide skill: in a KB you
  don't decide, you increase confidence — even NL assertions are
  evidence-capture events, not decision-capture. KB surface stays small (4 mode
  skills + ctx kb note); canonical capture skills unchanged.
- Phase KB ships handover + editorial paired, not split: the closeout/fold
  mechanism is the integration point; shipping paired stresses the fold on day
  one rather than retrofitting it.
- Editorial constitution lives at .context/ingest/KB-RULES.md, not
  CONSTITUTION.md: lifts the sibling project's resolved naming-collision (their
  10-INGEST_RULES.md rename) so ctx CONSTITUTION.md keeps its singular meaning;
  same discipline carries to domain-decisions.md vs DECISIONS.md.
- Phase KB lifts the *current* upstream pipeline shape (pass-mode contract,
  completion circuit breaker, source-coverage state-machine ledger,
  topic-adjacency pre-flight, cold-reader rubric, folder-shaped topics from day
  one, CLI-as-scaffold-authority), superseding the brief's 4-phase model —
  lifting the older shape would re-fight wounds the upstream author already
  healed.

---

## [2026-06-07-180011] Localizable vocabulary and i18n primitives (consolidated)

**Consolidated from**: 5 entries (2026-03-14 to 2026-05-23)

- Session prefixes are parser vocabulary, not i18n text: header-recognition
  patterns move to .ctxrc session_prefixes (default Session:), separating
  content recognition from interface language so users parse multilingual
  session files without code changes.
- Classify rules are user-configurable via .ctxrc (classify_rules overrides
  config/memory defaults) — same pattern as session_prefixes, for
  non-English/specialized domains.
- Spec signal words and the nudge threshold (spec_signal_words,
  spec_nudge_min_len) are .ctxrc-configurable — signal words are language- and
  project-dependent.
- Keep i18n.Fold strict (Unicode case-fold, İ≠i, for identifier
  dedup/parsing/security comparison); add i18n.MatchKey (Fold + NFKD + strip
  combining marks) as a separate diacritic-insensitive primitive for matching
  user input against vocabulary lists. Two explicit-contract primitives beat one
  conflated primitive or an options flag.
- Placeholder overrides use EXTEND, not REPLACE, semantics (diverging from
  SessionPrefixes' REPLACE): the dominant bilingual EN+TR case needs both
  default and added placeholders rejected simultaneously; REPLACE would silently
  lose baseline coverage. Opt-in placeholders_replace:true reserved if REPLACE
  is later wanted.

---

## [2026-03-14-093748] Config-driven freshness check with per-file review URLs

**Status**: Accepted

**Context**: Building a hook to warn when technology-dependent constants go
stale. Initially hardcoded the file list and Anthropic docs URL in the binary,
but this only worked inside the ctx repo and assumed all projects care about
Anthropic docs.

**Decision**: Config-driven freshness check with per-file review URLs

**Rationale**: Making the file list and review URLs configurable via .ctxrc
freshness_files means any project can opt in. Per-file review_url avoids
special-casing by project name — ctx sets Anthropic docs, other projects set
their own vendor links or omit it entirely.

**Consequence**: The hook is a no-op by default (opt-in). ctx's own .ctxrc
carries the tracked files. All nudge text goes through assets/text.yaml for
localization. No project detection logic needed.

---

