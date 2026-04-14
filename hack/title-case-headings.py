#!/usr/bin/env python3
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
"""Title-case headings (H1-H6) and admonition titles in Markdown files.

Style: AP-leaning. Lowercase ALL articles, prepositions (any length), and
coordinating conjunctions when they appear MID-phrase. First word, last word,
and the first word after a colon are always capitalized. Subordinating
conjunctions (when, while, after, before, etc.) are capitalized.

Protected verbatim:
  - Backticked code spans `...`
  - Markdown link URLs `(...)` immediately following `]`
  - Markdown reference-style link labels `][label]`
  - All-uppercase tokens of length >= 2 (acronyms)
  - Mixed-case tokens like macOS, GitHub, JavaScript, JSONL
  - Single uppercase letter labels (A, B, ... in 'Appendix A:')
  - Brand 'ctx' always lowercase, including possessive 'ctx's'
  - Version-number tokens (v0, v0.8.0)

Skipped contexts:
  - YAML frontmatter at file head
  - Fenced code blocks ```...```

Usage:
  hack/title-case-headings.py <path>            # dry-run, prints diffs
  hack/title-case-headings.py --apply <path>    # write changes in place

<path> may be a single Markdown file or a directory (recursively scanned for
*.md). Exits non-zero if any changes are needed (in dry-run mode), so it's
safe to wire into CI.
"""
import re
import sys
import pathlib

ARTICLES = {'a', 'an', 'the'}
PREPOSITIONS = {
    # AP-strict; ambiguous words (after/before/since/until/past/near/down/up/off)
    # excluded so they cap when they're conjunctions or adj/adv.
    'about','above','across','against','along','among','around','as','at',
    'behind','below','beneath','beside','between','beyond','by','despite',
    'during','except','for','from','in','inside','into','like','of','on',
    'onto','out','outside','over','per','plus','regarding','than','through',
    'throughout','till','to','toward','under','underneath','unto','upon',
    'versus','via','vs','with','within','without',
}
COORD_CONJ = {'and','or','but','nor','so','yet','for'}

LOWER_MID = ARTICLES | PREPOSITIONS | COORD_CONJ
BRAND_LOWER = {'ctx'}

WORD_RE = re.compile(r"[A-Za-z][A-Za-z0-9'/]*")

def title_case_word(word, force_cap=False):
    if not word:
        return word
    # Hyphenated: each segment treated as own word; segments after the first
    # are always capitalized (Chicago hyphen rule).
    if '-' in word:
        segs = word.split('-')
        new = [title_case_word(segs[0], force_cap=force_cap)]
        for s in segs[1:]:
            new.append(title_case_word(s, force_cap=True))
        return '-'.join(new)
    lw = word.lower()
    # Brand and brand-with-suffix (ctx, ctx's, ctxs)
    for brand in BRAND_LOWER:
        if lw == brand:
            return brand
        if lw.startswith(brand) and len(lw) > len(brand):
            tail = lw[len(brand):]
            # Allow possessive or short plural-like suffix
            if tail in ("'s", 's', "'") or (tail.startswith("'") and tail[1:].isalpha() and len(tail) <= 3):
                return brand + word[len(brand):]
    # Single uppercase letter label (e.g., 'A' in 'Appendix A')
    if len(word) == 1 and word.isupper():
        return word
    # Acronym already all-caps
    if len(word) >= 2 and word.isupper():
        return word
    # Version-number token: lowercase 'v' followed by digits (v0, v1, v0.8 etc.)
    if re.match(r'^v\d', word):
        return word
    # Mixed-case (interior caps): preserve (macOS, GitHub, JavaScript)
    if any(c.isupper() for c in word[1:]):
        return word
    # Mid-phrase function word
    if not force_cap and lw in LOWER_MID:
        return lw
    # Default: cap first letter
    return word[0].upper() + word[1:]

PROTECT_RE = re.compile(
    r'(`[^`]*`)'                  # backtick code span
    r'|(\][ ]*\([^)]*\))'         # markdown inline link URL incl. ']'
    r'|(\][ ]*\[[^\]]+\])'        # markdown reference-style link label
    # Brand tagline — italic lowercase, with or without quotes/punctuation.
    r'|(\*do you remember\??\*)'
    # CLI long-flag tokens (--keep-frontmatter, --keep-frontmatter=false)
    r'|(--[a-z][a-z0-9_-]*(?:=\S+)?)'
    # Slash-prefixed commands (/ctx-remember, /ctx-decision-add)
    r'|(/[a-z][a-z0-9_-]*)'
)

def split_protected(text):
    pieces = []
    last = 0
    for m in PROTECT_RE.finditer(text):
        if m.start() > last:
            pieces.append(('plain', text[last:m.start()]))
        pieces.append(('protected', m.group(0)))
        last = m.end()
    if last < len(text):
        pieces.append(('plain', text[last:]))
    return pieces

def title_case_text(text):
    pieces = split_protected(text)
    total = 0
    for kind, t in pieces:
        if kind == 'protected':
            total += 1
        else:
            total += len(WORD_RE.findall(t))
    if total == 0:
        return text
    overall = 0
    after_colon = False
    out = []
    for kind, t in pieces:
        if kind == 'protected':
            overall += 1
            out.append(t)
            continue
        last = 0
        buf = []
        for m in WORD_RE.finditer(t):
            literal = t[last:m.start()]
            if re.search(r':\s*$', literal):
                after_colon = True
            buf.append(literal)
            overall += 1
            is_first = (overall == 1)
            is_last = (overall == total)
            # Contraction tail: when the preceding literal ends with an
            # apostrophe (e.g. after a backtick brand span like `ctx`'s),
            # the word is a contraction suffix (s, ll, ve, t, re, d, m).
            # Preserve as lowercase, never capitalize.
            is_contraction_tail = literal.endswith("'") and len(m.group(0)) <= 3
            # Filename extension: when the preceding literal ends with '.'
            # and the word is a short lowercase token (md, sh, py, json,
            # yaml, txt, ...), preserve as lowercase.
            word_text = m.group(0)
            is_ext_tail = (
                literal.endswith('.')
                and len(word_text) <= 5
                and word_text.islower()
            )
            force = (is_first or is_last or after_colon) and not is_contraction_tail and not is_ext_tail
            if (is_contraction_tail or is_ext_tail) and not (is_first or is_last):
                buf.append(word_text.lower())
            else:
                buf.append(title_case_word(word_text, force_cap=force))
            after_colon = False
            last = m.end()
        trailing = t[last:]
        if re.search(r':\s*$', trailing):
            after_colon = True
        buf.append(trailing)
        out.append(''.join(buf))
    return ''.join(out)

def process_file(path):
    src = path.read_text(encoding='utf-8')
    lines = src.split('\n')
    in_fence = False
    in_fm = False
    out_lines = []
    changes = []
    for i, line in enumerate(lines):
        if i == 0 and line.strip() == '---':
            in_fm = True
            out_lines.append(line); continue
        if in_fm and line.strip() == '---':
            in_fm = False
            out_lines.append(line); continue
        if in_fm:
            out_lines.append(line); continue
        if line.startswith('```'):
            in_fence = not in_fence
            out_lines.append(line); continue
        if in_fence:
            out_lines.append(line); continue
        m = re.match(r'^(#{1,6})\s+(.*)$', line)
        if m:
            hashes, text = m.groups()
            new_text = title_case_text(text.rstrip())
            new_line = f"{hashes} {new_text}"
            if new_line != line:
                changes.append((i + 1, line, new_line))
            out_lines.append(new_line); continue
        am = re.match(r'^(\s*)([!?]{3})(\s+\w[\w-]*\s+)"([^"]+)"\s*$', line)
        if am:
            indent, marker, mid, title = am.groups()
            new_title = title_case_text(title)
            new_line = f'{indent}{marker}{mid}"{new_title}"'
            if new_line != line:
                changes.append((i + 1, line, new_line))
            out_lines.append(new_line); continue
        out_lines.append(line)
    return '\n'.join(out_lines), changes

def iter_md_files(target: pathlib.Path):
    if target.is_file():
        if target.suffix.lower() == '.md':
            yield target
        return
    if target.is_dir():
        yield from sorted(target.rglob('*.md'))

def main():
    args = sys.argv[1:]
    if '-h' in args or '--help' in args or not args:
        print(__doc__)
        sys.exit(0 if args else 2)
    apply = False
    if '--apply' in args:
        apply = True
        args.remove('--apply')
    if len(args) != 1:
        print("error: expected exactly one <path> argument", file=sys.stderr)
        sys.exit(2)
    target = pathlib.Path(args[0])
    if not target.exists():
        print(f"error: path not found: {target}", file=sys.stderr)
        sys.exit(2)
    total = 0
    files = 0
    for md in iter_md_files(target):
        new_src, changes = process_file(md)
        if changes:
            files += 1
            total += len(changes)
            for ln, old, new in changes:
                print(f"{md}:{ln}")
                print(f"  - {old}")
                print(f"  + {new}")
            if apply:
                md.write_text(new_src, encoding='utf-8')
    mode = 'APPLIED' if apply else 'dry-run'
    print(f"\n=== {total} changes across {files} files ({mode}) ===")
    sys.exit(0 if (apply or total == 0) else 1)

if __name__ == '__main__':
    main()
