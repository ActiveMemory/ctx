# text-markdown-serialization

## [2026-07-17-081010] An insert helper that returns content[:i]+x instead of content[:i]+x+content[i:] silently drops the tail

**Context**: insert.AfterHeader (the fallback beforeFirstEntry takes for a knowledge file with no ## [ entries) returned content[:insertPoint]+entry, truncating the file at the insertion point and discarding everything after it. Masked in practice because a ctx-init'd file has nothing after its comment block, so the dropped tail was empty; it bites the moment any non-entry section sits below the preamble of an as-yet-entry-less file. Same family as the LEARNINGS clobber bug index.Validate guards: silent memory loss, git-only recovery.

**Lesson**: A string-splice 'insert' must reattach the tail. content[:i]+x is append-with-truncation, NOT insert. The sibling Task() splice did it right (content[:i]+x+sep+content[i:]); AfterHeader was the odd one out and nobody noticed because the package had ZERO tests. EOF-anchored fixtures mask the bug precisely because the tail is empty there.

**Application**: Audit any result built from content[:i] for a matching content[i:] on the other side of the inserted text; a lone content[:i] is the smell. Test insert/splice helpers with a NON-empty tail, not just an EOF-anchored fixture. When touching a data-mutating helper with no tests, adding the test is part of the fix, not optional.

---

## [2026-07-17-081010] Uninitialized desc.Text() returns empty, and strings.Index(s, "") == 0 makes anchor-based inserts silently match at offset 0

**Context**: Writing a layout proof for the insert package, both cases passed on the first run — for the wrong reason. A *_test.go that never calls lookup.Init() (via TestMain) gets "" back from every desc.Text() call, because the embedded asset lookup is uninitialized in that test binary. Anchor logic such as beforeFirstEntry does strings.Index(content, desc.Text(headingKey)); with a "" needle that returns 0, so 'insert before the anchor' prepends to the whole file (above the H1) and any assertion of the form 'entry appears before X' passes trivially.

**Lesson**: strings.Index(s, "") == 0: an empty needle 'matches' at the start of any string. So a text-driven test whose helper resolves labels through desc.Text MUST initialize the asset lookup, or it exercises a degenerate offset-0 code path that production never takes — and passes for the wrong reason. This nearly caused a false 'measurement gate validated' report when the gate was in fact broken.

**Application**: Any package test that calls desc.Text (directly, or transitively through the code under test) needs a TestMain calling lookup.Init() (see internal/cli/system/core/session/testmain_test.go for the pattern). When a text/anchor-driven test passes suspiciously easily, dump desc.Text() of the keys involved and assert they are non-empty before trusting the assertions.

---

## [2026-07-04-152957] Typed JSON round-trips silently drop user-owned keys

**Context**: The init permissions merge read settings.local.json into a typed struct and re-marshaled it, so every key ctx does not model (env, a user's statusLine) vanished on re-init

**Lesson**: A typed read-modify-write of a file users also edit by hand is silent data loss; unknown fields do not survive encoding/json round-trips

**Application**: Mutate shared JSON via raw-map surgery (map[string]json.RawMessage), rewriting only the keys ctx owns; see internal/cli/initialize/core/merge/settings.go

---

## [2026-05-23-001000] Unicode block separation makes diacritic-stripping surgical — no per-script handling needed for Arabic/Indic/Hebrew/CJK

**Context**: While building `i18n.MatchKey` (commit 978582f5) for
diacritic-insensitive placeholder matching, the natural reflex was "this is
going to need per-script special cases — CJK doesn't have case, Arabic has
shadda/fatha that are meaning-changing, Bengali vowel signs are
script-essential, Hebrew niqqud distinguishes words." I sized the work assuming
we'd need a script-aware policy, possibly with a locale config or an opt-in flag
for "strip all combining marks" vs "strip only Latin-style decoration".
Empirical test across Turkish/German/French/Spanish/Catalan/Czech/Vietnamese
(should collapse) and Arabic/Bengali/Devanagari/Hindi/Hebrew/Chinese/Korean
(should preserve) showed the entire policy fits in one numeric range:
U+0300..U+036F.

**Lesson**: Unicode pre-separated combining marks by intent at the codepoint
level. The "Combining Diacritical Marks" block (U+0300–U+036F) holds
Latin/general decorative marks: acute, grave, diaeresis, tilde, cedilla, caron,
the Turkish combining dot, the Vietnamese horn, etc. Script-essential marks live
in separate blocks per script: Arabic in U+0610–U+06ED, Bengali in
U+0980–U+09FF, Devanagari in U+0900–U+097F, Hebrew niqqud in
U+0591–U+05C7, and so on. The block boundaries are not coincidental — they
encode the same distinction a reasonable design would want to make. So a narrow
byte-range strip is exactly the right primitive: it expresses "remove
decoration, keep structural marks" in one comparison, without needing to know
anything about the input's script.

**Application**: When designing comparison/normalization primitives for
international input, check the Unicode block boundaries before reaching for
per-script special cases or a config field. Often the standardization committee
already drew the line you want, and an arithmetic range check (`r >= 0x0300 && r
<= 0x036F`) does the work. Verify empirically across the scripts you care about
— but expect the answer to be cleaner than your initial sizing. The general
rule: when Unicode has put related characters in their own block, treat that
block as a meaningful unit of policy. (For ctx, this is now
`cfgI18n.CombiningMarksLatinStart`/`End` and the `MatchKey` implementation in
`internal/i18n/matchkey.go`.)

---

## [2026-05-11-231025] Naive Markdown line-sweep corrupts multi-line code spans and YAML lists

**Context**: Performed a programmatic typographic sweep across docs/*.md to wrap
bare 'ctx' tokens in backticks (commit 61aab858). 81 source files, 236 lines
changed. First pass corrupted two indented JSON snippets in MkDocs admonitions
because the fence regex anchored to start-of-line and missed admonition-indented
fences. After fixing the fence regex, two more corruptions surfaced (multi-line
inline-code spans where the opening backtick is on line N and the closing on
line N+1: the line-at-a-time transformer treated each line independently,
leading to misjudged span boundaries on the second line). After post-sweep
validation, a YAML parse error on docs/blog/2026-02-03-the-attention-budget.md
surfaced one more breakage: a 'topics:' list-item like '- ctx primitives' got
wrapped to '- `ctx` primitives', which is invalid YAML (a value starting with
backtick is not a valid unquoted scalar). Total: 2 multi-line span corruptions +
1 YAML breakage, all detected only by post-sweep validation (make site + grep
audit), not by the dry-run.

**Lesson**: A naive line-at-a-time regex sweep across Markdown documents must
respect a wider 'skip' set than the obvious cases. The full safe-skip list is:
(1) triple-backtick fenced code blocks, BOTH root-level and indented inside
MkDocs admonitions or list items (fence regex must allow leading whitespace,
e.g. '^\\s*```'); (2) inline backtick spans on the same line; (3) multi-line
inline-code spans crossing line boundaries (line-at-a-time logic cannot detect
both ends, so either track fence-like 'odd-count' state across lines or treat
any unclosed-on-line backtick as 'protect rest of line'); (4) the ENTIRE YAML
frontmatter block (delimited by '---' at top and next '---'), not just specific
keys like title/description/icon, because list-item values under
topics/tags/keywords are also YAML and break on a leading backtick; (5) image
alt-text '![alt]' (alt-text does not render in monotype); (6) link-reference
definitions '[name]: url "title"'; (7) project copyright header comment blocks.
Dry-run output never catches YAML or multi-line span breakage; validation MUST
include a parser-level check (make site for YAML, post-grep for '``name`'
double-backtick patterns near the wrapped token).

**Application**: When designing any future programmatic sweep across docs/
(typography passes, internationalization, brand renames, em-dash replacement,
link-text rewrites): (1) implement the full skip set above, not a subset; (2)
for fence detection use '^\\s*```', not '^```'; (3) for the frontmatter, detect
the entire block between '---' delimiters, not specific keys; (4) for multi-line
inline-code, choose between cross-line backtick-pair tracking (complex but
correct) or the simpler 'unclosed backtick protects rest of line' heuristic
(corrupts ~1 per 100 files but recoverable manually); (5) ALWAYS validate
post-sweep with 'make site' (zensical surfaces YAML errors) and a grep for
'``\\w' double-backtick patterns near the wrapped token; (6) commit only after
both validations are clean. For one-shot sweeps the script can be ad-hoc, but
record the validation gate as part of the commit message so the next contributor
knows what to check.

---

## [2026-04-14-010105] Brand-name handling in title-case engines must cover possessives

**Context**: First pass of hack/title-case-headings.py produced 'Ctx's' from
'ctx's' because the brand check matched the bare token only.

**Lesson**: A brand allowlist needs to recognize <brand>, <brand>'s, <brand>s,
and short apostrophe-suffixed variants. Single-word matching misses contractions
and possessives.

**Application**: When adding a new always-lowercase brand to
hack/title-case-headings.py, extend the suffix-aware loop in title_case_word,
not just the BRAND_LOWER set.

---

## [2026-04-03-133244] desc.Text() is the single highest-connectivity symbol in the codebase

**Context**: GitNexus enrichment during architecture analysis revealed
desc.Text() (internal/assets/read/desc/desc.go:75) has 30+ direct callers
spanning every architectural layer (MCP handler, format, index, tidy, trace,
memory, sysinfo, io) and participates in 53 execution flows.

**Lesson**: TestDescKeyYAMLLinkage is the most critical guard in the codebase
— it protects the symbol with the widest blast radius. If YAML text loading
breaks, the entire CLI and MCP server output blank strings silently (no crash,
no warning).

**Application**: Treat desc.Text() as a frozen API — add new functions rather
than modifying the existing signature. Any change to config/embed/text or
assets/read/desc should be followed by running the linkage audit. Monitor this
symbol during major refactors.

---

## [2026-03-06-141504] Stats sort uses string comparison on RFC3339 timestamps with mixed timezones

**Context**: ctx system stats showed only old sessions, hiding the current one

**Lesson**: RFC3339 string comparison breaks when entries mix UTC (Z) and offset
(-08:00) formats — 13:00-08:00 sorts before 18:00Z lexicographically despite
being later in absolute time

**Application**: Always parse to time.Time before comparing RFC3339 timestamps;
never rely on lexicographic sort

---

## [2026-03-01-095709] TASKS.md template checkbox syntax inside HTML comments is parsed by RegExTaskMultiline

**Context**: Template had example checkboxes (- [x], - [ ]) in HTML comments
that the line-based regex matched as real tasks, causing
TestArchiveCommand_NoCompletedTasks to fail

**Lesson**: RegExTaskMultiline is line-based and has no awareness of HTML
comment blocks — checkbox-like patterns inside comments get counted as real
tasks

**Application**: Use backtick-quoted or indented references instead of actual
checkbox syntax in template comments. When adding examples to TASKS.md
templates, avoid patterns that match regExTaskPattern

---

