#!/usr/bin/env python3
"""Auto-verify files that pass CommonMark fence validation."""
import re
import glob
import os
from datetime import date

VERIFIED_RE = re.compile(r'<!-- fences-verified: \d{4}-\d{2}-\d{2} -->')
NORMALIZED_RE = re.compile(r'<!-- normalized: \d{4}-\d{2}-\d{2} -->')
FENCE_RE = re.compile(r'^(\s*)(`{3,}|~{3,})(.*?)$')
TODAY = date.today().isoformat()


def is_commonmark_correct(path):
    with open(path) as f:
        content = f.read()

    if VERIFIED_RE.search(content):
        return None  # already verified

    lines = content.split('\n')

    start = 0
    if lines and lines[0].strip() == '---':
        for i in range(1, len(lines)):
            if lines[i].strip() == '---':
                start = i + 1
                break

    in_fence = False
    fence_char = None
    fence_count = 0

    for i in range(start, len(lines)):
        line = lines[i]
        m = FENCE_RE.match(line)
        if not m:
            continue

        char = m.group(2)[0]
        count = len(m.group(2))
        info = m.group(3).strip()

        if not in_fence:
            in_fence = True
            fence_char = char
            fence_count = count
        else:
            if char == fence_char and count >= fence_count and not info:
                in_fence = False

    return not in_fence  # True if all fences closed


def add_verified_marker(path):
    with open(path) as f:
        content = f.read()

    marker = f'<!-- fences-verified: {TODAY} -->'
    nm = NORMALIZED_RE.search(content)
    if nm:
        pos = nm.end()
        content = content[:pos] + '\n' + marker + content[pos:]
    else:
        lines = content.split('\n')
        if lines and lines[0].strip() == '---':
            for i in range(1, len(lines)):
                if lines[i].strip() == '---':
                    lines.insert(i + 1, '')
                    lines.insert(i + 2, marker)
                    content = '\n'.join(lines)
                    break
        else:
            content = marker + '\n' + content

    with open(path, 'w') as f:
        f.write(content)


def main():
    files = sorted(glob.glob('.context/journal/*.md'))
    verified = 0
    skipped = 0

    for path in files:
        result = is_commonmark_correct(path)
        if result is None:
            skipped += 1
            continue
        if result:
            add_verified_marker(path)
            verified += 1
            print(f'  verified: {os.path.basename(path)}')

    print(f'\n{verified} auto-verified, {skipped} already done')


if __name__ == '__main__':
    main()
