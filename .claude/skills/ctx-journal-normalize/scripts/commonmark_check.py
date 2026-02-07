#!/usr/bin/env python3
"""Check fence correctness using CommonMark rules.

CommonMark fence rules:
- Opening fence: 3+ backticks/tildes, optional info string
- Closing fence: same char type, >= count, NO info string
- Inside a fence, EVERYTHING is literal until the closing fence
- A line with fewer backticks than the opening is just text

This detects files where fences are correctly matched vs broken.
"""
import re
import glob
import os

VERIFIED_RE = re.compile(r'<!-- fences-verified: \d{4}-\d{2}-\d{2} -->')
FENCE_RE = re.compile(r'^(\s*)(`{3,}|~{3,})(.*?)$')


def check_commonmark(path):
    with open(path) as f:
        content = f.read()

    if VERIFIED_RE.search(content):
        return 'verified', []

    lines = content.split('\n')

    # Skip frontmatter
    start = 0
    if lines and lines[0].strip() == '---':
        for i in range(1, len(lines)):
            if lines[i].strip() == '---':
                start = i + 1
                break

    problems = []
    in_fence = False
    fence_char = None
    fence_count = 0
    fence_start = 0

    for i in range(start, len(lines)):
        line = lines[i]
        m = FENCE_RE.match(line)
        if not m:
            continue

        char = m.group(2)[0]
        count = len(m.group(2))
        info = m.group(3).strip()

        if not in_fence:
            # Opening fence
            in_fence = True
            fence_char = char
            fence_count = count
            fence_start = i + 1
        else:
            # Inside a fence — check for closing
            if char == fence_char and count >= fence_count and not info:
                # Valid close
                in_fence = False
            # If char != fence_char, or count < fence_count, or has info:
            # just literal text inside the fence — fine

    if in_fence:
        problems.append(f'unclosed fence at line {fence_start} ({fence_char}*{fence_count})')

    return ('ok' if not problems else 'broken'), problems


def main():
    files = sorted(glob.glob('.context/journal/*.md'))
    verified = []
    ok = []
    broken = []

    for f in files:
        name = os.path.basename(f)
        status, problems = check_commonmark(f)
        if status == 'verified':
            verified.append(name)
        elif status == 'ok':
            ok.append(name)
        else:
            broken.append((name, problems))

    print(f'=== VERIFIED: {len(verified)} ===')
    print(f'=== CORRECT (can auto-verify): {len(ok)} ===')
    for name in ok:
        print(f'  {name}')
    print(f'\n=== BROKEN (need fixing): {len(broken)} ===')
    for name, problems in broken:
        print(f'  {name}')
        for p in problems:
            print(f'    {p}')
    print(f'\n--- Summary ---')
    print(f'Verified: {len(verified)}')
    print(f'Correct: {len(ok)}')
    print(f'Broken: {len(broken)}')


if __name__ == '__main__':
    main()
