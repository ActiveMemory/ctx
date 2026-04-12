---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Multi-machine
icon: lucide/network
---

![ctx](../images/ctx-banner.png)

# `ctx` Hub: Multi-machine

Run the hub on a **LAN host** and connect from project directories
on other workstations. This recipe is the **Story 2 ("small trusted
team")** shape described in the
[`ctx` Hub overview](hub-overview.md) — read that first if
you haven't, especially the trust-model warnings.

This recipe assumes you've already walked through
[Getting Started](hub-getting-started.md) and understand
what flows through the hub (decisions, learnings, conventions,
tasks — **not** journals, scratchpad, or raw context files).

## Topology

```
+------------------+        +------------------+
| workstation A    |        | workstation B    |
|  ~/projects/x    |        |  ~/projects/y    |
|  ctx connection  |        |  ctx connection  |
+---------+--------+        +---------+--------+
          |                           |
          +-----------+   +-----------+
                      v   v
              +-------------------+
              | LAN host "nexus"  |
              | ctx hub start     |
              | --daemon          |
              | :9900             |
              +-------------------+
```

## Step 1 — Start the daemon on the LAN host

On the machine that will hold the hub (call it `nexus`):

```bash
ctx hub start --daemon --port 9900
```

The daemon writes a PID file to `~/.ctx/hub-data/hub.pid`. Stop it
later with:

```bash
ctx hub stop
```

## Step 2 — Firewall and port

Open port `9900/tcp` on `nexus` to the LAN only — **never** expose
the hub to the public internet without a reverse proxy and TLS in
front of it (see [Hub security model](../security/hub.md)).

Typical LAN allowlist rules:

=== "firewalld"

    ```bash
    sudo firewall-cmd --zone=internal \
      --add-port=9900/tcp --permanent
    sudo firewall-cmd --reload
    ```

=== "ufw"

    ```bash
    sudo ufw allow from 192.168.1.0/24 to any port 9900 proto tcp
    ```

=== "nftables"

    ```bash
    sudo nft add rule inet filter input ip saddr 192.168.1.0/24 \
      tcp dport 9900 accept
    ```

## Step 3 — Retrieve the admin token

The daemon prints the admin token to stdout on first run. Running as
a daemon, that output goes to the log instead:

```bash
cat ~/.ctx/hub-data/admin.token
```

Copy the token over a trusted channel (SSH, password manager, or
an encrypted note). **Do not email it or put it in chat.**

## Step 4 — Register projects from each workstation

On workstation `A`:

```bash
cd ~/projects/x
ctx connection register nexus.local:9900 --token ctx_adm_...
ctx connection subscribe decision learning convention
```

On workstation `B`:

```bash
cd ~/projects/y
ctx connection register nexus.local:9900 --token ctx_adm_...
ctx connection subscribe decision learning convention
```

Each registration exchanges the admin token for a **per-project
client token**. Only the client token is persisted in
`.context/.connect.enc`, encrypted with the same AES-256-GCM scheme
ctx uses for notification credentials.

## Step 5 — Verify

From either workstation:

```bash
ctx connection status
```

You should see the ctx Hub address, role (`leader` for single-node),
subscription filters, and the sequence number you're synced to.

## TLS (recommended)

For anything beyond a trusted home LAN, terminate TLS in front of
the hub. The hub speaks gRPC, so the reverse proxy must speak
HTTP/2:

```nginx
server {
    listen 443 ssl http2;
    server_name nexus.example.com;

    ssl_certificate     /etc/letsencrypt/live/nexus.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/nexus.example.com/privkey.pem;

    location / {
        grpc_pass grpc://127.0.0.1:9900;
    }
}
```

Point `ctx connection register` at the public hostname and port 443.

## Handling daemon restarts

The hub is **append-only JSONL** — restarts are safe. Clients keep
their last-seen sequence in `.context/hub/.sync-state.json` and
pick up exactly where they left off on the next `sync` or `listen`
reconnect.

## See also

- [HA cluster recipe](hub-cluster.md) — for redundancy
- [Hub operations](../operations/hub.md) — backup, rotation
- [Hub failure modes](../operations/hub-failure-modes.md)
- [Hub security model](../security/hub.md)
