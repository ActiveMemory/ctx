# Learnings Reference

Module-specific, niche, and historical learnings moved from LEARNINGS.md
to keep the main file within token budget. All entries preserved verbatim.

---

## [2026-02-26-100001] ctx init and CLAUDE.md behavior (consolidated)

**Consolidated from**: 4 entries (2026-01-20 to 2026-02-14)

- ctx init is non-destructive: only creates .context/, CLAUDE.md, .claude/, PROMPT.md, and IMPLEMENTATION_PLAN.md. Zero awareness of .cursorrules, .aider.conf.yml, or other tools' configs.
- CLAUDE.md merge insertion is position-aware: findInsertionPoint() finds the first H1, skips trailing blank lines, and inserts there. Never appends to end.
- CLAUDE.md handling is a 3-state machine: no file (create), file without ctx markers (merge/prompt), file with `<!-- ctx:context -->` / `<!-- ctx:end -->` markers (skip or force-replace).
- Always backup before modifying user files: file.bak before modification, marker comments for idempotency, offer merge not overwrite, provide `--merge` escape hatch.

---

## [2026-02-26-100003] Blog and content publishing (consolidated)

**Consolidated from**: 4 entries (2026-01-28 to 2026-02-18)

- Blog posts are living documents: periodic enrichment passes (cross-links, update admonitions, citations) improve reader experience. Schedule enrichment as part of consolidation sessions.
- Blog publishing from ideas/ requires a 7-step checklist: (1) update date/frontmatter, (2) fix relative paths, (3) add cross-links, (4) add "The Arc" section, (5) update blog index, (6) add "See also" in related posts, (7) verify all link targets exist.
- Changelogs document WHAT (for machines/audits); blogs explain WHY (for humans/narrative). When synthesizing session history, output both.
- Cross-repo links to published docs should use ctx.ist URLs, not GitHub blob URLs or relative paths. GitHub won't render zensical admonitions and readers lose navigation context.

---

## [2026-02-26-100004] Worktrees and parallel agents (consolidated)

**Consolidated from**: 5 entries (2026-02-12 to 2026-02-24)

- Git worktrees enable parallel Claude Code sessions without file conflicts. Create as sibling directories (`git worktree add ../ctx-docs -b work/docs`). 3-4 parallel worktrees is the practical limit before merge complexity outweighs productivity gains.
- Parallel agents work cleanly on disjoint file sets: partition by file domain (Go source, docs, specs/context) for zero-conflict parallel execution.
- Worktree agents lack key-dependent features by design: ctx pad fails gracefully (no key), ctx notify silently no-ops. The only real gap is notify — pad and journal are naturally avoided by workflow.
- Multi-agent work requires commit auditing: changes intermingle in the working tree with no commit boundaries. Run `git diff --stat` and group by logical feature before committing.
- rsync between worktrees can clobber permissions and gitignored files. Use `--no-perms` or `--chmod=+x` for scripts, and `--exclude` gitignored paths. The /absorb skill handles these edge cases.

---

## [2026-02-26-100011] Plugin and marketplace architecture (consolidated)

**Consolidated from**: 3 entries (2026-02-16)

- When repo-local .claude/skills/ and a marketplace plugin both define the same skill name, Claude Code lists both: local unprefixed and plugin with `ctx:` namespace prefix. Ensure distributed skills live only in the plugin source.
- Claude Code marketplace plugins source from the repo root where `.claude-plugin/marketplace.json` lives. Edits to skills and hooks under the plugin path take effect on next Claude Code load — no reinstall needed.
- Security docs are most vulnerable to stale paths after architecture migrations. After any file-layout migration, grep security docs for old paths first — stale paths in security guidance give users a false sense of protection.

---

## [2026-02-22-120003] Journal site rendering pipeline (consolidated)

**Consolidated from**: 8 entries (2026-02-20)

- Python-Markdown ends ALL HTML blocks at blank lines — no CommonMark Type 1 exception for `<pre>`. Only fenced code blocks survive blank lines
- `<details>` is Type 6 HTML block — ends at first blank line; collapsible tool output needs CSS/JS approach
- pymdownx.highlight hijacks `<pre><code>` patterns; disable with `use_pygments=false`
- normalizeContent three-layer fix: always run wrapToolOutputs, pick fence depth exceeding inner fences, stripPreWrapper only unescapes when `<pre>` found
- Tool output boundary detection: pre-scan turn numbers, sort+dedup, find min > N, use LAST positional occurrence (last-match-wins)
- Inline code spans with angle brackets: replace backticks with double-quotes and brackets with HTML entities via RegExInlineCodeAngle
- Title sanitization: strip Claude tags, replace angle brackets, strip backticks/hash, truncate to 75 chars (RecallMaxTitleLen)
- AI normalization poor fit for bulk files; code-level pipeline handles rendering at build time

---

## [2026-02-22-120004] Skills evaluation and design (consolidated)

**Consolidated from**: 6 entries (2026-01-23 to 2026-02-15)

- Skills are markdown files in `.claude/commands/` with YAML frontmatter; `$ARGUMENTS` passes args
- Skills that restate or contradict system prompt create tension — check platform defaults first
- Red flags: urgency tags, reasoning overrides, tables labeling hesitation as wrong — discard entirely
- ~80% of external skill files are redundant; apply E/A/R classification, only keep expert-level knowledge delta
- Prefer skills over CLI commands for judgment-heavy workflows; reserve CLI for deterministic operations
- When a skill would edit files controlling agent behavior (permissions, hooks), use a runbook instead

---

## [2026-02-22-120005] Claude Code session and JSONL format (consolidated)

**Consolidated from**: 4 entries (2026-01-28 to 2026-02-04)

- Claude Code, Cursor, Aider each have different JSONL formats; use tool-agnostic Session type with tool-specific parsers
- JSONL files are append-only, never shrink after compaction; file size overreports post-compaction
- Subagent files in `/subagents/` share parent sessionId; skip when scanning (check `isSidechain:true`)
- `slug` field removed in Claude Code v2.1.29+; parse by `sessionId` + valid `type` instead

---

## [2026-02-22-120007] Superseded session-related learnings (consolidated)

**Consolidated from**: 4 entries (2026-01-20)

> **Historical**: `.context/sessions/` was removed in v0.4.0. These learnings are superseded but preserved for context.

- SessionEnd hook fires on all exits including Ctrl+C — hook behavior still accurate for Claude Code, but ctx no longer uses it
- Session filenames used YYYY-MM-DD-HHMMSS-topic.md — journal entries now use `ctx recall export` naming
- Two tiers remain: curated (`.context/*.md`) and full dump (`~/.claude/projects/` + `.context/journal/`); middle `.context/sessions/` tier eliminated
- Auto-load via PreToolUse worked; auto-save via SessionEnd removed because Claude Code retains transcripts natively

---

## [2026-02-27-002832] Zensical site builder quirks (consolidated)

**Consolidated from**: 2 entries (2026-02-21)

- Section-header icons only render when the section has an index.md listed as its first nav entry (navigation.indexes feature). Icons in child page frontmatter don't propagate. Always create `<section>/index.md` with `icon:` frontmatter for new top-level nav sections.
- `zensical serve -a IP:PORT` overrides config file dev_addr. `ctx journal site --serve` hardcodes `zensical serve` with no args — use Make targets with `-a` flag for per-developer overrides instead of config file changes.

---

## [2026-02-27-002831] Journal and recall parsing edge cases (consolidated)

**Consolidated from**: 4 entries (2026-02-03 to 2026-02-24)

- /ctx-journal-normalize is dangerous at scale: on large JSONL files it blows up subagent context windows and produces nondeterministic output. Keep expensive AI skills out of batch pipelines; offer as targeted per-file tools.
- normalizeCodeFences regex treats non-whitespace adjacent to triple-backtick fences as split points, separating language tags. Use plain fences without lang tags in tests (see internal/cli/recall/fmt_test.go).
- Users naturally type inline code fences (`text: ```code`) without proper newline separation. Normalize on export with regex that inserts `\n\n` around fences; apply only to user messages.
- Claude Code injects `<system-reminder>` tags into tool result content in JSONL. When wrapped in code fences, these XML-like tags break markdown rendering. Extract system reminders before wrapping and render as markdown outside the fence.

---

## [2026-02-27-002833] Project identity and structure clarifications (consolidated)

**Consolidated from**: 3 entries (2026-01-20 to 2026-02-06)

- PROMPT.md is a Ralph loop iteration prompt ("what to do next, how to know when done"), not a project briefing. When it drifts into duplicating CLAUDE.md, delete it. Re-introduce only when actively using Ralph loops.
- Only `internal/assets/` (formerly `internal/templates/`) matters for embedded templates — it's where Go embeds files into the binary. A root `templates/` directory is spec baggage. One source of truth: `internal/assets/ ──[ctx init]──> .context/`.
- ctx and Ralph Loop are separate systems: `ctx init` creates `.context/` for context management; Ralph Loop uses PROMPT.md, IMPLEMENTATION_PLAN.md, and specs/ for iterative AI development. ctx does NOT create Ralph Loop infrastructure.

---

## [2026-02-24-015713] url.Parse works for SMB URLs

**Context**: Shell script used sed to extract host/share from smb://host/share. Needed a Go equivalent.

**Lesson**: Go net/url.Parse handles smb:// scheme correctly — u.Host gives hostname, u.Path gives share path with leading slash.

**Application**: Use url.Parse for any custom-scheme URL parsing instead of hand-rolled regex.

---

## [2026-02-22-101959] Loop script templates use positional format verbs — adding a verb requires updating all callers

**Context**: Adding TplLoopNotify to TplLoopScript and TplLoopMaxIter required updating the fmt.Sprintf calls in script.go to pass the new argument at the correct position.

**Lesson**: Go fmt.Sprintf with positional args is fragile — adding a %s to a template shifts all downstream arguments. The comment documenting Args: in the template constant is the only guard.

**Application**: When modifying loop templates, always update the Args: comment and verify script.go passes arguments in the correct order.

---

## [2026-02-20-142442] Default export already preserves enrichment — T2.1 was partially stale

**Context**: Investigated ctx recall export --update and found the default behavior already preserves YAML frontmatter during re-export. The --force flag has a bug where it claims to discard frontmatter but does not.

**Lesson**: Always read the current code before speccing a feature — the need may already be met, and the real work may be a bug fix rather than a new feature.

**Application**: When speccing tasks from the backlog, investigate current state first. Rewrite the task to reflect what is actually needed.

---

## [2026-02-19-214909] GCM authentication makes try-decrypt a reliable format discriminator

**Context**: Needed to auto-detect whether pad merge input files are encrypted or plaintext without relying on file extensions or user flags.

**Lesson**: Authenticated encryption (AES-256-GCM) guarantees that decryption with the wrong key always fails — unlike unauthenticated ciphers that produce silent garbage. This makes 'try decrypt, fall back to plaintext' a safe and simple detection strategy.

**Application**: Use try-decrypt-first as the default pattern for any ctx feature that handles mixed encrypted/plaintext input. No need for format flags or extension-based heuristics.

---

## [2026-01-26-180000] Go json.Marshal Escapes Shell Characters

**Context**: Generated settings.local.json had `2\u003e/dev/null` instead
of `2>/dev/null`.

**Lesson**: Go's json.Marshal escapes `>`, `<`, and `&` as
unicode (`\u003e`, `\u003c`, `\u0026`) by default for HTML safety.
Use `json.Encoder` with `SetEscapeHTML(false)` when generating config files
that contain shell commands.

---
