---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Hub Security Model
icon: lucide/shield
---

![ctx](../images/ctx-banner.png)

# `ctx` Hub: Security model

What the hub defends against, what it **does not** defend against,
and the concrete mechanisms in play.

## Threat model

The hub is designed for **trusted cross-project knowledge sharing**
within a team or homelab. It assumes:

- The hub host is trusted. Anyone with root on that box can read
  every entry ever published.
- Network is semi-trusted. Hub traffic is gRPC over TCP; TLS is
  **strongly recommended** but not mandatory.
- Client machines are trusted enough to hold a per-project client
  token. Losing a client token is roughly equivalent to losing an
  API key — scoped damage, not total compromise.
- Entry content is **not** secret. Decisions, learnings, and
  conventions may be indexed by AI agents, rendered in docs,
  shared across projects. Do not push credentials or PII into
  the hub.

The hub is **not** a secure messaging system, a secrets store, or
a compliance-grade audit log. If your threat model needs those,
use a dedicated tool and keep the hub for knowledge sharing.

## Mechanisms

### Bearer tokens

All RPCs except `Register` require a bearer token in gRPC
metadata. Two kinds of tokens exist:

| Kind         | Format        | Scope                           | Lifetime            |
|--------------|---------------|---------------------------------|---------------------|
| Admin token  | `ctx_adm_...` | Register new projects           | Manual rotate       |
| Client token | `ctx_cli_...` | Publish, Sync, Listen, Status   | Project lifetime    |

Tokens are compared in **constant time** (`crypto/subtle`) to
prevent timing oracles, and looked up via an `O(1)` hash map so
the comparison cost does not depend on the total number of
registered clients.

### Client-side encryption at rest

`.context/.connect.enc` stores the client token and hub address,
encrypted with **AES-256-GCM** using the same scheme the
notification subsystem uses. The key is derived from ctx's local
keyring (see `internal/crypto`).

An attacker with read access to the project directory cannot
learn the client token without also breaking ctx's local
keyring.

### Hub-side token storage

!!! warning "Tokens are stored in plaintext on the hub host"
    `<data-dir>/clients.json` currently stores client tokens
    **verbatim**, not hashed. Anyone with read access to the
    hub's data directory sees every registered client's token
    and can impersonate any project that has ever registered.

    Mitigations **today**:

    - Run the hub as an unprivileged user and lock the data
      directory with `chmod 700 <data-dir>`.
    - Use the systemd unit in
      [Operations](../operations/hub.md#systemd-unit),
      which enables `ProtectSystem=strict`,
      `NoNewPrivileges=true`, and a dedicated user.
    - Never expose `<data-dir>` over NFS, SMB, or shared
      filesystems.
    - Treat `<data-dir>` the same way you'd treat
      `/etc/shadow` — back it up encrypted, never check it
      into version control.

    Hashing `clients.json` and moving to keyring-backed storage
    is tracked as a follow-up in the PR #60 task group. Until
    that lands, assume a hub host compromise equals total hub
    compromise.

### Input validation

Every published entry is validated before it touches the log:

- **Type** must be one of: `decision`, `learning`, `convention`,
  `task`. Unknown types are rejected.
- **ID** and **Origin** are required and non-empty.
- **Content size** is capped at **1 MB**. Reasonable for text,
  hostile for attempts to fill the disk.
- **Duplicate project registration** is rejected — a client that
  replays an old `Register` call gets an error, not a second
  token.

### No script execution

The hub never interprets entry content. There is no expression
language, no template evaluation, no markdown rendering at
ingest. Content is stored as bytes and fanned out to clients
verbatim.

### Audit trail

`entries.jsonl` is append-only. Every accepted publish is
recorded with the publishing project's origin tag and sequence
number. Nothing is ever deleted by the hub; retention is managed
manually by the operator (see
[log rotation](../operations/hub.md#log-rotation)).

## What the hub does **not** defend against

- **Untrusted entry senders.** A client with a valid token can
  publish anything (within the 1 MB cap). There is no content
  validation beyond shape.
- **Denial of service from a registered client.** A misbehaving
  client can publish until disk is full. Monitor
  `entries.jsonl` growth.
- **Network eavesdropping without TLS.** Plain gRPC leaks entry
  content and tokens. Use a TLS-terminating reverse proxy
  (see [Multi-machine recipe](../recipes/hub-multi-machine.md#tls-recommended)).
- **Host compromise.** Root on the hub host = access to every
  entry and every token. Harden the host.
- **Accidental secret upload.** The hub will happily fan out a
  decision containing an API key. Sanitize content before
  publishing.

## Operational hardening checklist

- [ ] Run the hub as an **unprivileged user** with
      `NoNewPrivileges=true` and `ProtectSystem=strict` (see
      the systemd unit in
      [Operations](../operations/hub.md#systemd-unit)).
- [ ] Terminate **TLS** in front of the hub for anything beyond
      a trusted LAN.
- [ ] Restrict the listen port with firewall rules to the
      client subnet only.
- [ ] Back up `<data-dir>/admin.token` to a secrets manager; do
      not leave it in shell history.
- [ ] Rotate the admin token when a team member with access
      leaves. Client tokens keep working across rotations.
- [ ] Monitor `entries.jsonl` growth; alert on sudden spikes.
- [ ] Run NTP on all clients to prevent entry-timestamp skew.
- [ ] Do not publish from machines you do not trust.

## Responsible disclosure

Security issues in the hub follow the same process as the rest
of ctx — see [Reporting](reporting.md).

## See also

- [`ctx` Hub Operations](../operations/hub.md)
- [`ctx` Hub failure modes](../operations/hub-failure-modes.md)
- [HA cluster recipe](../recipes/hub-cluster.md)
