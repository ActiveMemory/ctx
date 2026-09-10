---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Hub Operations
icon: lucide/settings
---

![ctx](../images/ctx-banner.png)

# `ctx` Hub: Operations

Running the `ctx` `ctx` Hub in production. This page is
for **operators**: people running a hub for themselves or a
team, not people writing to a hub someone else is running.

If you have not read it yet, start with the
[`ctx` Hub overview](../recipes/hub-overview.md). It
explains what the hub is, the two user stories it supports
(personal cross-project brain vs small trusted team), and what
it does **not** do. A client-side tour is in
[Getting Started](../recipes/hub-getting-started.md).

!!! info "Operator Cheat Sheet"
    - The hub fans out four entry types only: `decision`,
      `learning`, `convention`, `task`. Journals, scratchpad,
      and other local state are out of scope.
    - Identity is per-**project**, not per-user. Attribution is
      limited to `Origin`, which is self-asserted by the
      publishing client.
    - The data model is an **append-only JSONL log** plus two
      small JSON sidecar files. Nothing is rewritten in place.

## Data Directory Layout

The hub stores everything under a single data directory
(default `~/.ctx/hub-data/`, override with `--data-dir`).

```
<data-dir>/
  admin.token        # Initial admin token (chmod 600)
  clients.json       # Registered client tokens and project names
  meta.json          # Sequence counter, version, cluster metadata
  entries.jsonl      # Append-only log (single source of truth)
  hub.pid            # Daemon PID file (daemon mode only)
  raft/              # Raft state (cluster mode only)
    log.db
    stable.db
    snapshots/
```

**Invariants:**

* `entries.jsonl` is **append-only**. Every line is a valid JSON
  object. Corrupt lines are fatal at startup: fix or truncate
  before restart.
* `meta.json` is authoritative for the next sequence number. On
  restart, the hub reads the last valid line of `entries.jsonl` and
  refuses to start if the sequences disagree.
* `clients.json` holds hashed client tokens; losing it invalidates
  all client registrations.

## Starting and Stopping

=== "Foreground"

    ```bash
    ctx hub start                    # Ctrl-C to stop
    ctx hub start --port 8080        # Custom port
    ctx hub start --data-dir /srv/ctx-hub
    ```

=== "Daemon"

    ```bash
    ctx hub start --daemon           # Fork to background
    ctx hub stop                      # Graceful shutdown
    ```

`--stop` sends SIGTERM to the PID in `hub.pid`, waits for
in-flight RPCs to drain, then exits. If the daemon is wedged,
remove `hub.pid` and send `SIGKILL` manually. `entries.jsonl` is
crash-safe, so you will not lose accepted writes.

## Systemd Unit

For production single-node deployments, run the hub as a systemd
service instead of `--daemon`:

```ini
# /etc/systemd/system/ctx-hub.service
[Unit]
Description=ctx `ctx` Hub
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=ctx
Group=ctx
ExecStart=/usr/local/bin/ctx hub start --port 9900 \
    --data-dir /var/lib/ctx-hub
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/ctx-hub
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable --now ctx-hub
sudo journalctl -u ctx-hub -f
```

## Backup and Restore

Because `entries.jsonl` is append-only, backups are trivial:

```bash
# Hot backup, safe while the hub is running.
cp <data-dir>/entries.jsonl backups/entries-$(date +%F).jsonl
cp <data-dir>/meta.json      backups/meta-$(date +%F).json
cp <data-dir>/clients.json   backups/clients-$(date +%F).json
```

For a consistent snapshot across all three files, stop the hub,
copy, then start again, or use a filesystem-level snapshot (LVM,
ZFS, Btrfs).

**Restore:**

```bash
ctx hub stop                           # Stop the hub
cp backups/entries-2026-04-10.jsonl <data-dir>/entries.jsonl
cp backups/meta-2026-04-10.json      <data-dir>/meta.json
cp backups/clients-2026-04-10.json   <data-dir>/clients.json
ctx hub start --daemon
```

Clients that pushed sequences **above** the restored watermark
will re-publish, because the hub now reports a lower sequence
than what clients have on disk. Nothing is lost, but the store is
append-only and does not deduplicate by entry ID: those entries
come back with new sequence numbers, so the shared feed shows
them twice. Prune the duplicates offline if they matter.

## Log Rotation

`entries.jsonl` grows unbounded. For long-lived hubs, rotate it
offline:

```bash
ctx hub stop
mv <data-dir>/entries.jsonl <data-dir>/entries-$(date +%F).jsonl.old
# Replay the last N days into a fresh entries.jsonl if you want a
# trimmed active log, or leave the old file in place as history.
ctx hub start --daemon
```

Do **not** truncate `entries.jsonl` while the hub is running.
The hub holds an open file handle; an in-place truncation confuses
the sequence counter and loses writes.

## Monitoring

Liveness probe:

```bash
ctx hub status
```

The command exits non-zero when the RPC fails — hub unreachable,
token rejected — so a wrapper can treat that as the liveness
signal. On a reachable node it prints what the node knows about
itself, and the grading is yours to do from those lines:

```
Role: Leader
Leader: 10.0.0.5:9901
Entries: 1248  Peers: 2
```

For cluster deployments, watch for:

- **Role flaps**: the `Role:` line changing more than once per
  hour suggests network instability or disk contention.
- **A leader nobody can name**: `Leader: unknown (election in
  progress)` is normal for a second or two after a node starts
  or a leader dies. Persisting past that, on a node that is
  itself reachable, means it cannot see a quorum.
- **A peer count that disagrees between nodes**: each node
  reports its own committed Raft configuration, so two nodes
  disagreeing on `Peers:` means they never agreed on membership.
- **`entries.jsonl` growth rate**: sudden spikes often indicate a
  misbehaving `ctx connection listen` reconnect loop.

## Upgrading

The JSONL format is versioned in `meta.json`. `ctx` refuses to start
against a newer store version than it understands; older store
versions are upgraded in place at first start after an upgrade.

**Always back up `<data-dir>/` before upgrading.**

## See Also

- [`ctx` Hub failure modes](hub-failure-modes.md)
- [`ctx` Hub security model](../security/hub.md)
- [`ctx serve` reference](../cli/serve.md)
- [`ctx hub` reference](../cli/hub.md)
