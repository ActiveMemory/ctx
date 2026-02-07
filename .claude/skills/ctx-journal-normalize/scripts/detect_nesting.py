#!/usr/bin/env python3
"""Detect files with nested code fences (fences inside open fences).

Files with only flat (non-nested) fence pairs can be auto-verified.
Files with nesting need AI reconstruction.
"""
import re
import glob
import os

FENCE_RE = re.compile(r'^(\s*)(`{3,}|~{3,})(.*)', re.MULTILINE)
VERIFIED_RE = re.compile(r'<!-- fences-verified: \d{4}-\d{2}-\d{2} -->')

def check_nesting(path):
    with open(path) as f:
        content = f.read()

    if VERIFIED_RE.search(content):
        return 'verified', 0

    lines = content.split('\n')
    fence_stack = []  # stack of (backtick_count, line_num)
    nesting_depth = 0
    max_nesting = 0

    # Skip frontmatter
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

        # Check if this closes the current fence
        if fence_stack and fence_char == fence_stack[-1][0] and fence_len >= fence_stack[-1][1] and not info:
            fence_stack.pop()
        else:
            # Opening fence
            fence_stack.append((fence_char, fence_len))
            if len(fence_stack) > max_nesting:
                max_nesting = len(fence_stack)

    return ('nested' if max_nesting > 1 else 'flat'), max_nesting


def main():
    files = sorted(glob.glob('.context/journal/*.md'))
    flat = []
    nested = []
    verified = []

    for f in files:
        name = os.path.basename(f)
        status, depth = check_nesting(f)
        if status == 'verified':
            verified.append(name)
        elif status == 'nested':
            nested.append((name, depth))
        else:
            flat.append(name)

    print(f"=== VERIFIED (already done): {len(verified)} files ===")
    for name in verified:
        print(f"  {name}")

    print(f"\n=== FLAT (auto-verifiable): {len(flat)} files ===")
    for name in flat:
        print(f"  {name}")

    print(f"\n=== NESTED (need AI): {len(nested)} files ===")
    for name, depth in sorted(nested, key=lambda x: -x[1]):
        print(f"  depth={depth}  {name}")

    print(f"\n--- Summary ---")
    print(f"Verified: {len(verified)}")
    print(f"Flat (can auto-verify): {len(flat)}")
    print(f"Nested (need AI): {len(nested)}")


if __name__ == '__main__':
    main()
