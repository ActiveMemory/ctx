#!/usr/bin/env python3
"""Fix consolidation annotations on closing fence lines.

`consolidateToolRuns` appends (×N) to closing fence lines like:
    ``` (×2)

This makes the line look like an opening fence (with info string) instead
of a closing fence, causing cascade nesting errors.

Fix: split to two lines:
    ```
    (×2)
"""
import re
import glob
import os

PATTERN = re.compile(r'^(```+)\s+(\(×\d+\))$')


def fix_file(path):
    with open(path) as f:
        lines = f.readlines()

    changed = False
    out = []
    for line in lines:
        m = PATTERN.match(line.rstrip('\n'))
        if m:
            out.append(m.group(1) + '\n')
            out.append(m.group(2) + '\n')
            changed = True
        else:
            out.append(line)

    if changed:
        with open(path, 'w') as f:
            f.writelines(out)

    return changed


def main():
    files = sorted(glob.glob('.context/journal/*.md'))
    fixed = 0

    for path in files:
        if fix_file(path):
            fixed += 1
            print(f'  fixed: {os.path.basename(path)}')

    print(f'\n{fixed} files fixed')


if __name__ == '__main__':
    main()
