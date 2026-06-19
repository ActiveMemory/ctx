---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: AI Backends
icon: lucide/brain-circuit
---

![ctx](../images/ctx-banner.png)

## `ctx ai`

Optional AI backend commands. These commands are additive: deterministic
commands such as `ctx status`, `ctx agent`, and hooks do not require a backend.

Configure a backend first:

```bash
ctx setup --backend vllm --endpoint http://localhost:8000 --write
```

### `ctx ai ping`

Check the selected backend by reading `/v1/models`.

```bash
ctx ai ping [--backend <name>]
```

Output includes the resolved backend name, endpoint, and first model listed.

### `ctx ai propose`

Send an input file to the selected backend and write a reviewable proposal
artifact. It never edits `.context/*.md` directly.

```bash
ctx ai propose <input> --emit decisions,learnings,tasks,open-questions \
  [--backend <name>]
```

Artifacts are written under `.context/proposals/ai/` as JSON proposed-patch
records with backend, model, input, emit kinds, status, and the decoded response.

### Backend selection

`ctx ai` uses this order:

1. `--backend <name>`
2. `.ctxrc` `backends.default`
3. the only configured backend, if exactly one exists

If multiple backends exist and none is selected, the command fails closed.
If no backend is configured or the backend is unreachable, the command fails
closed; non-AI commands are unaffected.
