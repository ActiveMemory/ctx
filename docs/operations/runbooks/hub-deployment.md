---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Hub Deployment
icon: lucide/server
---

![ctx](../../images/ctx-banner.png)

# Hub Deployment

Linear runbook for setting up a ctx Hub for yourself or a team.
Consolidates pieces currently scattered across hub recipes and
operations docs.

**When to use**: First-time hub setup, or when onboarding a new
team onto an existing hub.

**Prerequisites**: `ctx` binary installed, network connectivity
between hub and clients.

**Companion docs**:

- [Hub overview](../../recipes/hub-overview.md): what the hub
  is and is not
- [Hub operations](../hub.md): data directory, systemd,
  backup, monitoring
- [Hub failure modes](../hub-failure-modes.md): what can go wrong

---

## Step 1: Start the Hub

=== "Quick Start (foreground)"

    ```bash
    ctx hub start
    ```

=== "Production (systemd)"

    See [Hub Operations: Systemd Unit](../hub.md#systemd-unit)
    for the full unit file.

    ```bash
    sudo systemctl enable --now ctx-hub
    ```

The hub creates `admin.token` on first start. Save this token;
it is the only way to register clients.

## Step 2: Generate the Admin Token

On first start, the hub writes `admin.token` to the data
directory (default `~/.ctx/hub-data/`):

```bash
cat ~/.ctx/hub-data/admin.token
```

This token has full admin privileges. Keep it secret.

## Step 3: Register Clients

For each client (person or machine) that will connect:

```bash
# On the hub machine
ctx hub register --name "volkan-laptop" --admin-token <admin-token>
```

This returns a client token. Distribute it securely to the client.

## Step 4: Connect Clients

On each client machine, register the project with the hub. The
`ctx hub *` commands above run on the hub server itself and don't
need a project. The `ctx connection *` commands below are different:
they live inside a project (the encrypted hub config is stored at
`.context/.connect.enc`), so you have to tell `ctx` which project
first.

```bash
# In the project directory on the client machine:
eval "$(ctx activate)"
ctx connection register <hub-address> --token <client-token>
```

Verify the connection:

```bash
ctx connection status
```

If the client doesn't have a project yet, run `ctx init` first, then
`eval "$(ctx activate)"`. See
[Activating a Context Directory](../../recipes/activating-context.md).

## Step 5: Verify Sync

Push a test entry from one client and verify it arrives. Make sure
each client already ran `eval "$(ctx activate)"` from Step 4:
otherwise `ctx add` and `ctx status` fail with
`Error: no context directory specified`.

```bash
# Client A (in its project directory, after activating):
ctx learning add "Hub sync test" --context "Verifying hub setup"

# Client B (in its project directory, after activating):
ctx status   # should show the new learning
```

## Step 6: Configure Backup

Set up regular backups of the hub data directory. See
[Hub Operations: Backup and Restore](../hub.md#backup-and-restore).

Minimum:

```bash
# Add to cron
0 */6 * * * cp ~/.ctx/hub-data/entries.jsonl ~/backups/entries-$(date +\%F).jsonl
```

## Step 7: Configure TLS (When Available)

!!! note "Coming Soon"
    TLS support is planned (H-01/H-02). Until then, run the hub
    on a trusted network or behind a reverse proxy with TLS
    termination.

---

## Team Onboarding Checklist

When adding a new team member to an existing hub:

- [ ] Generate a client token (`ctx hub register --name "<name>"`)
- [ ] Share the token and hub address securely
- [ ] Have them run `ctx connect <hub-address> --token <token>`
- [ ] Verify with `ctx connection status`
- [ ] Point them to the [Hub Getting Started](../../recipes/hub-getting-started.md) recipe

## Troubleshooting

### "Connection Refused"

The hub isn't running or the port is wrong. Check:

```bash
ctx hub status          # on the hub machine
ss -tlnp | grep 9900   # default port
```

### "Authentication Failed"

The client token is wrong or was never registered. Re-register:

```bash
ctx hub register --name "<name>" --admin-token <admin-token>
```

### Entries Not Syncing

Check that the client is listening:

```bash
ctx connection status
```

If connected but not syncing, check the hub logs for sequence
mismatch errors. See
[Hub Failure Modes](../hub-failure-modes.md) for details.
