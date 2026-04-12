---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Serve
icon: lucide/server
---

![ctx](../images/ctx-banner.png)

## `ctx serve`

Serve a static site locally via [zensical](https://pypi.org/project/zensical/).

With no argument, serves the journal site at
`.context/journal-site`. With a directory argument, serves
that directory if it contains a `zensical.toml`.

```bash
ctx serve                             # Serve .context/journal-site
ctx serve ./my-site                   # Serve a specific directory
ctx serve ./docs                      # Serve any zensical site
```

!!! info "This command does NOT start a hub"
    `ctx serve` is purely for static-site serving. To run a
    `ctx` Hub for cross-project knowledge sharing, use
    [`ctx hub start`](hub.md). That command lives in its
    own group because the hub is a gRPC server, not a
    static site.

**Requires zensical to be installed**:

```bash
pipx install zensical
```

### Arguments

| Argument     | Description                                      |
|--------------|--------------------------------------------------|
| `[directory]` | Directory containing a `zensical.toml` to serve |

When omitted, serves `.context/journal-site` by default — the
directory produced by `ctx journal site`.

**Examples**:

```bash
ctx serve                         # Default: serve .context/journal-site
ctx serve ./my-site               # Serve a specific directory
ctx serve ./docs                  # Serve any zensical site
```

### See also

- [`ctx journal`](journal.md) — generate the journal site
  that `ctx serve` displays.
- [`ctx hub start`](hub.md) — for running a `ctx` Hub
  server, not a static site.
- [Browsing and enriching past sessions](../recipes/session-archaeology.md)
  — the recipe that combines `ctx journal` and `ctx serve`.
