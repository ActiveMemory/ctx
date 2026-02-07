#!/usr/bin/env python3
"""Journal normalizer: metadata tables + idempotency marker.

Deterministic pass only. Fence reconstruction is left to the AI skill.

Phase 1: Convert **Key**: value metadata blocks to collapsed HTML tables
Phase 2: Add <!-- normalized: YYYY-MM-DD --> marker after frontmatter
"""
import re
import glob
from datetime import date

METADATA_RE = re.compile(r'^\*\*(\w[\w\s/]*?)\*\*:\s*(.*)')
SUMMARY_KEYS = ['Date', 'Duration', 'Turns', 'Model']
NORMALIZED_RE = re.compile(r'<!-- normalized: \d{4}-\d{2}-\d{2} -->')
TODAY = date.today().isoformat()


def fix_metadata(lines):
    """Convert consecutive **Key**: value lines to collapsed HTML tables."""
    out = []
    i = 0
    changed = False

    while i < len(lines):
        m = METADATA_RE.match(lines[i].strip())
        if not m:
            out.append(lines[i])
            i += 1
            continue

        pairs = []
        while i < len(lines):
            m = METADATA_RE.match(lines[i].strip())
            if not m:
                break
            pairs.append((m.group(1), m.group(2)))
            i += 1

        if len(pairs) < 2:
            out.append(f'**{pairs[0][0]}**: {pairs[0][1]}')
            continue

        changed = True
        vals = {k: v for k, v in pairs}
        parts = [vals[k] for k in SUMMARY_KEYS if k in vals and vals[k]]
        summary = ' \u00b7 '.join(parts) if parts else 'Session metadata'

        out.append('<details>')
        out.append(f'<summary>{summary}</summary>')
        out.append('<table>')
        for k, v in pairs:
            out.append(f'<tr><td><strong>{k}</strong></td><td>{v}</td></tr>')
        out.append('</table>')
        out.append('</details>')

    return out, changed


def add_marker(lines):
    """Add <!-- normalized: YYYY-MM-DD --> after frontmatter."""
    marker = f'<!-- normalized: {TODAY} -->'

    if lines and lines[0].strip() == '---':
        for i in range(1, len(lines)):
            if lines[i].strip() == '---':
                return lines[:i+1] + ['', marker] + lines[i+1:], True

    return [marker, ''] + lines, True


def normalize_file(path):
    with open(path) as f:
        content = f.read()

    if NORMALIZED_RE.search(content):
        return False

    lines = content.split('\n')

    lines, _ = fix_metadata(lines)
    lines, _ = add_marker(lines)

    with open(path, 'w') as f:
        f.write('\n'.join(lines))

    return True


def main():
    files = sorted(glob.glob('.context/journal/*.md'))
    processed = 0
    skipped = 0

    for path in files:
        if normalize_file(path):
            processed += 1
            print(f'  normalized: {path}')
        else:
            skipped += 1

    print(f'\n{processed} normalized, {skipped} skipped, {len(files)} total')


if __name__ == '__main__':
    main()
