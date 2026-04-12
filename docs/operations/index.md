---
title: Operations
icon: lucide/settings
---

![ctx](../images/ctx-banner.png)

Guides for **installing**, **upgrading**, **integrating**, and
**running** `ctx`. Split into three groups by audience.

---

## Hub

Operator guides for running a `ctx` Hub — the gRPC server that
fans out structured entries across projects. If you're a client
connecting to a Hub someone else runs, see
[`ctx connect`](../cli/connection.md) and the
[Hub recipes](../recipes/hub-overview.md) instead.

### [Hub Operations](hub.md)

Data directory layout, daemon management, systemd unit,
backup and restore, log rotation, monitoring, and upgrades.

### [Hub Failure Modes](hub-failure-modes.md)

What can go wrong in network, storage, cluster, auth, and
clock layers — and what you should do about each one. Includes
the short-list table oncall engineers will want bookmarked.

---

## Operating `ctx`

Everyday operation guides for anyone running `ctx` in a
project or adopting it in a team.

### [Integration](migration.md)

Adopt `ctx` in an existing project: initialize context files,
migrate from other tools, and onboard team members.

### [Upgrade](upgrading.md)

Upgrade between versions with step-by-step migration notes
and breaking-change guidance.

### [AI Tools](integrations.md)

Configure `ctx` with Claude Code, Cursor, Aider, Copilot,
Windsurf, and other AI coding tools.

### [Autonomous Loops](autonomous-loop.md)

Run an unattended AI agent that works through tasks overnight,
with `ctx` providing persistent memory between iterations.

---

## Maintainers

Runbooks for people shipping `ctx` itself.

### [Cutting a Release](release.md)

Step-by-step runbook for maintainers: bump version, generate
release notes, run the release script, and verify the result.

---

## Runbooks

Step-by-step procedures in `hack/runbooks/`. Run these
regularly — they catch problems that linters and tests cannot.

| Runbook | Purpose | When to run |
|---------|---------|-------------|
| [Codebase audit](../../hack/runbooks/codebase-audit.md) | AST audits, magic strings, dead code, doc alignment | Before release, quarterly |
| [Docs semantic audit](../../hack/runbooks/docs-semantic-audit.md) | Narrative gaps, weak pages, structural problems | Before release, after adding pages |
| [Sanitize permissions](../../hack/runbooks/sanitize-permissions.md) | Clean `.claude/settings.local.json` of over-broad grants | After heavy permission granting |

**Recommended cadence**:

- **Before every release**: codebase audit + docs semantic audit
- **Monthly**: sanitize permissions
- **Quarterly**: full sweep of all three

The `_ctx-release` skill runs the codebase audit automatically
as part of its pre-release checks. The other two are manual —
add them to your release checklist.
