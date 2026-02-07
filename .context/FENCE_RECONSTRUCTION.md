# Fence Reconstruction Tracker

Reconstruct proper backtick fence nesting in journal source files so the site
renders with code highlighting instead of stripped plain text.

## Architecture

- `ctx journal site` strips ALL fence markers from site copies (via `stripFences`)
- Files with `<!-- fences-verified: YYYY-MM-DD -->` skip stripping (fences trusted)
- Two markers: `normalized` (metadata done) and `fences-verified` (fences correct)

## Process per file

1. Read the source file from `.context/journal/`
2. Reconstruct fences: innermost code gets 3 backticks, each nesting level adds 1
3. Add `<!-- fences-verified: YYYY-MM-DD -->` marker
4. Write back
5. Check the box below

## Scripts

In `.claude/skills/ctx-journal-normalize/scripts/`:
- `normalize.py` — metadata tables + normalized marker
- `detect_nesting.py` — stack-based nesting detector (useful but overly conservative)
- `auto_verify_flat.py` — auto-add fences-verified to flat files
- `fix_consolidation_fences.py` — split `` ``` (×N) `` to separate lines
- `commonmark_check.py` — **correct** CommonMark fence validator
- `auto_verify_commonmark.py` — auto-verify files passing CommonMark validation

## Learnings

1. **Stack-based nesting detection is misleading.** A stack model counts
   ```` wrapping ``` as "depth 2", but CommonMark renders this correctly
   (inner ``` is literal text inside the 4-backtick fence). The stack detector
   flagged 83 files as "nested (need AI)" — CommonMark validation found only
   8 actually broken.

2. **`consolidateToolRuns` bug**: appends `(×N)` to closing fence lines,
   making them look like opening fences. E.g. `` ``` (×2) `` parses as an
   opening fence with info string "(×2)". Fix: split to separate lines.
   112 occurrences across 45 files, resolved deterministically.

3. **CommonMark fence rules are simpler than they seem:**
   - Opening: 3+ backticks/tildes + optional info string
   - Closing: same char, >= count, NO info string
   - Inside a fence, everything is literal until close
   - ```` properly wraps ``` content

4. **Deterministic > AI**: Each scripted fix reduces the AI workload
   dramatically. Started with 98 files "needing AI", ended with 8 after
   three deterministic passes.

## Progress

- [x] 15 files with no fences (nothing to do)
- [x] 129 files verified (scripts + 1 manual fix)
- [x] 8 files with broken fences — fixed by removing stray fence markers

**All 137 journal files now have `<!-- fences-verified: 2026-02-06 -->` markers.**

## Broken files — all fixed

- [x] `2026-02-03--3d875788.md` — removed stray ``` at L3916
- [x] `2026-01-22--f99b05b4.md` — removed stray ```` at L3750
- [x] `2026-01-23--06e093ff.md` — removed stray ```` at L903
- [x] `2026-01-23--76fe2ab9-p11.md` — removed stray ```` at L3509
- [x] `2026-01-23--9105e489.md` — removed stray ```` at L6409
- [x] `2026-01-23-lazy-fluttering-eich-9105e489.md` — removed stray ``` at L5321
- [x] `2026-01-27--9830db47.md` — removed stray ```` at L4463
- [x] `2026-01-27-gleaming-wobbling-sutherland-37944bf2-p4.md` — removed stray ``` at L3627

## Pattern for the 8 broken files

All 8 had stray fence markers — either ```` after `</details>` tags (5 files)
or ``` at block boundaries (3 files). These were leftovers from the
`consolidateToolRuns` processing that weren't caught by the deterministic
`fix_consolidation_fences.py` pass. Simple removal fixed all of them — no
fence reconstruction or nesting adjustment was needed.
