#!/usr/bin/env python3
"""Auto-add fences-verified marker to flat (non-nested) fence files."""
import re
import glob
import os
from datetime import date

FENCE_RE = re.compile(r'^\s*(`{3,}|~{3,})(.*)', re.MULTILINE)
VERIFIED_RE = re.compile(r'<!-- fences-verified: \d{4}-\d{2}-\d{2} -->')
NORMALIZED_RE = re.compile(r'<!-- normalized: \d{4}-\d{2}-\d{2} -->')
TODAY = date.today().isoformat()


def has_nesting(path):
    with open(path) as f:
        content = f.read()

    if VERIFIED_RE.search(content):
        return None  # already verified

    lines = content.split('\n')
    fence_stack = []

    start = 0
    if lines and lines[0].strip() == '---':
        for i in range(1, len(lines)):
            if lines[i].strip() == '---':
                start = i + 1
                break

    for i in range(start, len(lines)):
        line = lines[i]
        m = re.match(r'^\s*(`{3,}|~{3,})(.*)', line)
        if not m:
            continue

        fence_char = m.group(1)[0]
        fence_len = len(m.group(1))
        info = m.group(2).strip()

        if fence_stack and fence_char == fence_stack[-1][0] and fence_len >= fence_stack[-1][1] and not info:
            fence_stack.pop()
        else:
            fence_stack.append((fence_char, fence_len))
            if len(fence_stack) > 1:
                return True  # nested

    return False  # flat


def add_verified_marker(path):
    with open(path) as f:
        content = f.read()

    marker = f'<!-- fences-verified: {TODAY} -->'

    # Insert after normalized marker if present
    nm = NORMALIZED_RE.search(content)
    if nm:
        pos = nm.end()
        content = content[:pos] + '\n' + marker + content[pos:]
    else:
        # Insert after frontmatter
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
        result = has_nesting(path)
        if result is None:
            skipped += 1
            continue
        if result is False:
            add_verified_marker(path)
            verified += 1
            print(f'  verified: {os.path.basename(path)}')

    print(f'\n{verified} auto-verified, {skipped} already done')


if __name__ == '__main__':
    main()
