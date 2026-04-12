---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Site
icon: lucide/globe
---

![ctx](../images/ctx-banner.png)

## `ctx site`

Site management commands for the ctx.ist static site.

```bash
ctx site <subcommand>
```

### `ctx site feed`

Generate an Atom 1.0 feed from finalized blog posts in
`docs/blog/`.

```bash
ctx site feed [flags]
```

Scans `docs/blog/` for files matching `YYYY-MM-DD-*.md`,
parses YAML frontmatter, and generates a valid Atom feed.
Only posts with `reviewed_and_finalized: true` are included.
Summaries are extracted from the first paragraph after the
heading.

**Flags**:

| Flag         | Short | Type   | Default           | Description              |
|--------------|-------|--------|-------------------|--------------------------|
| `--out`      | `-o`  | string | `site/feed.xml`   | Output path              |
| `--base-url` |       | string | `https://ctx.ist` | Base URL for entry links |

**Output**:

```
Generated site/feed.xml (21 entries)

Skipped:
  2026-02-25-the-homework-problem.md: not finalized

Warnings:
  2026-02-09-defense-in-depth.md: no summary paragraph found
```

Three buckets: **included** (count), **skipped** (with
reason), **warnings** (included but degraded). `exit 0`
always: warnings inform but do not block.

**Frontmatter requirements**:

| Field                    | Required | Feed mapping                |
|--------------------------|----------|-----------------------------|
| `title`                  | Yes      | `<title>`                   |
| `date`                   | Yes      | `<updated>`                 |
| `reviewed_and_finalized` | Yes      | Draft gate (must be `true`) |
| `author`                 | No       | `<author><name>`            |
| `topics`                 | No       | `<category term="">`        |

**Examples**:

```bash
ctx site feed                                # Generate site/feed.xml
ctx site feed --out /tmp/feed.xml            # Custom output path
ctx site feed --base-url https://example.com # Custom base URL
make site-feed                               # Makefile shortcut
make site                                    # Builds site + feed
```
