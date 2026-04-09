---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Serve
icon: lucide/server
---

## `ctx serve`

Serve a journal site locally, or start the shared context hub.

### Static site (default)

```bash
ctx serve                    # Serve .context/journal/site/
ctx serve ./my-site          # Serve a specific directory
```

### Shared context hub

Start a gRPC hub server for cross-project knowledge sharing.

```bash
ctx serve --shared                          # Start on default port 9900
ctx serve --shared --port 8080              # Custom port
ctx serve --shared --data-dir /path/to/data # Custom data directory
```

On first run, generates an admin token and prints it to stdout.
Save this token — it's needed for `ctx connect register`.

**Default data directory:** `~/.ctx/hub-data/`

### Daemon mode

Run the hub as a background process:

```bash
ctx serve --shared --daemon          # Start in background
ctx serve --stop                     # Stop the running daemon
```

The daemon writes a PID file to `<data-dir>/hub.pid`.

### Cluster mode

For high availability, run multiple hubs with Raft leader election:

```bash
ctx serve --shared --port 9900 --peers host2:9901,host3:9901
```

Raft is used only for leader election. Data replication uses
sequence-based gRPC sync (append-only, no conflicts).

### Validation

The hub validates all published entries:
- **Type** must be `decision`, `learning`, `convention`, or `task`
- **ID** and **Origin** are required (non-empty)
- **Content** max 1MB (text-only — decisions, learnings, conventions)
- **Duplicate registration** is rejected (one token per project)

### Flags

| Flag | Description |
|------|-------------|
| `--shared` | Start the shared context hub |
| `--port` | Hub listen port (default 9900) |
| `--data-dir` | Hub data directory (default ~/.ctx/hub-data/) |
| `--daemon` | Run hub in the background |
| `--stop` | Stop a running hub daemon |
| `--peers` | Comma-separated peer addresses for cluster mode |
