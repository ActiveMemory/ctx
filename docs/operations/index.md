---
title: Operations
icon: lucide/settings
---

![ctx](../images/ctx-banner.png)

Guides for **installing**, **upgrading**, **integrating**, and
**running** `ctx`. Split into three groups by audience.

---

## Day-to-Day

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

## Hub

Operator guides for running a `ctx` Hub, the gRPC server that
fans out structured entries across projects. If you're a client
connecting to a Hub someone else runs, see
[`ctx connect`](../cli/connection.md) and the
[Hub recipes](../recipes/hub-overview.md) instead.

### [Hub Operations](hub.md)

Data directory layout, daemon management, systemd unit,
backup and restore, log rotation, monitoring, and upgrades.

### [Hub Failure Modes](hub-failure-modes.md)

What can go wrong in network, storage, cluster, auth, and
clock layers, and what you should do about each one. Includes
the short-list table oncall engineers will want bookmarked.

---

## Maintainers

Runbooks for people shipping `ctx` itself.

### [Cutting a Release](release.md)

Step-by-step runbook for maintainers: bump version, generate
release notes, run the release script, and verify the result.

---

## Runbooks

Step-by-step procedures you run with your agent. Each runbook
includes a prompt to paste into a Claude Code session and
guidance on triaging the results.

| Runbook | Purpose | When to run |
|---------|---------|-------------|
| [Release checklist](runbooks/release-checklist.md) | Full pre-release sequence | Before every release |
| [Plugin release](runbooks/plugin-release.md) | Plugin-specific release steps | Plugin changes ship |
| [Breaking migration](runbooks/breaking-migration.md) | Guide users across breaking changes | Releases with renames |
| [Hub deployment](runbooks/hub-deployment.md) | Set up a ctx Hub end-to-end | First-time hub setup |
| [New contributor](runbooks/new-contributor.md) | Onboarding: clone to first session | New contributors |
| [Codebase audit](runbooks/codebase-audit.md) | AST audits, magic strings, dead code, doc alignment | Before release, quarterly |
| [Docs semantic audit](runbooks/docs-semantic-audit.md) | Narrative gaps, weak pages, structural problems | Before release, after adding pages |
| [Sanitize permissions](runbooks/sanitize-permissions.md) | Clean `.claude/settings.local.json` of over-broad grants | After heavy permission granting |
| [Architecture exploration](runbooks/architecture-exploration.md) | Systematic architecture docs across repos | New codebase onboarding, reviews |

**Recommended cadence**:

- **Before every release**: release checklist (which includes
  codebase audit + docs semantic audit)
- **Monthly**: sanitize permissions
- **Quarterly**: full sweep of all audit runbooks
