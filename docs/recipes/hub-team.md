---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Team knowledge bus
icon: lucide/users
---

![ctx](../images/ctx-banner.png)

# Team knowledge bus

This recipe shows **how a small trusted team uses a `ctx`
Hub as a shared knowledge bus** — the "Story 2" shape
from the [Hub overview](hub-overview.md). You're not
building a wiki, you're not replacing your issue tracker,
and you're not running a multi-tenant service. You're
connecting 3-10 developers who trust each other so that
lessons, decisions, and conventions flow between them
without ceremony.

**Prerequisites**:

- A running `ctx` Hub on a LAN host or internal server
  everyone on the team can reach. See
  [Multi-machine setup](hub-multi-machine.md) for the
  deployment guide.
- Each team member has `ctx` installed and has
  `ctx connection register`-ed their working projects with
  the hub.

## Trust model — read this first

The hub assumes **everyone holding a client token is
friendly**. There's no per-user attribution you can rely
on, no read ACL beyond subscription filters, and `Origin`
is self-asserted by the publishing client. Treat the hub
like a team wiki: useful because everyone can write to
it, **not** because it can prove *who* wrote what.

If your team is:

- ✅ 3-10 engineers, all known to each other, all
  trusted with production access
- ✅ On a single internal network or behind a VPN
- ✅ Comfortable with "the hub assumes friendly
  participants"

…this recipe fits. If your team is:

- ❌ Larger than ~15, with turnover
- ❌ Includes contractors, untrusted agents, or
  compromised-workstation concerns
- ❌ Needs audit trails that prove who published what
- ❌ Requires per-team-member isolation

…you're in "Story 3" territory, which the hub **does
not** support today. Use a wiki or a dedicated knowledge
platform instead.

## The team's three verbs

Everyone on the team does three things, same as in the
[personal recipe](hub-personal.md), but with different
social expectations:

1. **Record** — when you learn something that would save
   a teammate time, capture it with `ctx add --share`.
2. **Subscribe** — every engineer's project directories
   subscribe to the types the team cares about.
3. **Load** — agents pick up shared entries automatically
   via the auto-sync hook and the `--include-hub` flag
   in the PreToolUse hook pipeline.

The operational shape is identical to solo use. What's
different is the *culture* around publishing: when do
you `--share`, and what belongs on the hub vs. in your
local `.context/`.

## What goes on the hub (team rules of thumb)

**Share it if it's true for more than one person.** The
central question: "would the next teammate who hits this
problem save time if they already knew this?" If yes,
`--share`. If no, record it locally and move on.

**Decisions**:

- ✅ Cross-service decisions (database choice, auth
  model, deployment pattern, monitoring stack).
- ✅ Policy decisions that apply to all services
  (naming, API versioning, error-message format).
- ❌ Internal implementation decisions inside a single
  service ("chose a map over a slice here because lookups
  dominate").
- ❌ One-off tactical calls for a specific PR.

**Learnings**:

- ✅ Gotchas, surprising behavior, flaky infrastructure
  quirks — anything you'd tell a teammate over coffee
  with "watch out for X".
- ✅ Lessons from incidents — right after the postmortem
  is the highest-value time to share.
- ❌ Internal debugging notes that only make sense with
  context from your current branch.

**Conventions**:

- ✅ Repo layout, commit message format, pre-commit
  hooks, review expectations.
- ✅ Language-level style decisions that apply across
  services.
- ❌ Per-service idioms ("in `billing/` we prefer…").

**Tasks**: almost always project-local. Don't subscribe
to `task` unless the team has a specific reason (e.g., a
cross-cutting migration you want visible everywhere).

## A realistic week

**Monday — 3 AM incident, shared learning**

On-call engineer Alice gets paged: the payment service
starts returning 500s after a dependency update. After
an hour she finds the culprit — a breaking change in a
transitive gRPC dep that only manifests under high
concurrency. Postmortem on Tuesday, but right now she
records the learning:

```bash
ctx add learning --share \
  --context "Payment service 3 AM incident, 2026-04-03" \
  --lesson  "grpc-go v1.62+ changes DialContext behavior under high concurrency: connections from a single channel can deadlock if the server emits GOAWAY mid-stream. Symptom: 500 errors cluster in 30s bursts, no error in grpc client logs." \
  --application "Any service on grpc-go. Pin to v1.61 or patch with keepalive: https://github.com/grpc/grpc-go/issues/..." 
```

By Tuesday morning, every other engineer's agent
context packet contains this learning. When Bob starts
work on the `ledger` service (which also uses grpc-go),
his Claude Code session already knows about the gotcha
without Bob having to read the incident channel.

**Wednesday — cross-service decision**

The team agrees on a new pattern for API versioning —
header-based instead of URL-based. Platform lead Carol
records the decision:

```bash
ctx add decision --share \
  --context "Need consistent API versioning across all 6 services. Current URL-based /v1/ isn't working for gradual rollouts." \
  --rationale "Header-based versioning lets us route by header at the edge, which makes canary rollouts trivial. URL-based versioning forces clients to update their paths." \
  --consequence "All new endpoints use X-API-Version header. Existing /v1/ endpoints stay. Deprecation schedule in q3." \
  "Use header-based API versioning for new endpoints"
```

Every engineer's next session knows about this decision
automatically. When Dave starts adding endpoints to the
`inventory` service on Thursday, Claude already prompts
him for the header pattern instead of defaulting to
`/v1/`.

**Friday — convention drift caught at review**

Dave notices that his PR auto-formatted some error
messages to end with periods. He recalls the team
convention is "no trailing period" but can't remember
where it was documented. He runs `ctx connection status`,
sees the hub is healthy, greps his local
`.context/hub/conventions.md`, and finds:

```markdown
## [2026-03-12] Error message format
Lowercase start, no trailing period, single sentence.
```

He fixes the PR. No lookup on the wiki, no question in
chat, no context-switch penalty.

## Workflow tips for teams

**Designate a "champion" for decisions.** The team lead
or platform engineer should be the person who explicitly
`--share`s cross-cutting decisions. Other team members
share learnings freely but should ask "should this be a
decision?" in review before `--share`ing a decision. This
keeps the decision stream signal-rich.

**Publish postmortem learnings immediately, not after
the meeting.** The postmortem itself is a document; the
*actionable rules* that come out of it belong on the
hub, and they should land within an hour of the
incident. "Share fast, edit later" is the rule.

**Delete noisy entries, don't tolerate them.** The hub
is append-only, but the `.context/hub/` mirror on each
client is just markdown. If a shared learning turns out
to be wrong or obsolete, remove it from local mirrors
and stop the hub daemon to truncate `entries.jsonl`
(see [Hub operations](../operations/hub.md)). Noisy
shared feeds lose trust fast.

**Don't subscribe every project to every type.** For
backend engineers, subscribing to `decision + learning +
convention` is usually right. For platform or DevOps
projects, adding `task` makes sense. For a prototype or
experiment project, subscribing only to `convention`
might be enough.

**Run a single hub, not one per team.** If two teams
need to share knowledge, they should share a hub.
Splitting hubs by team creates silos — which is often
exactly the thing you were trying to solve.

## Operational concerns

The team recipe assumes someone owns the hub host. That
person (or a small group) is responsible for:

- **Uptime**: the hub is infrastructure; treat it like
  any other internal service you run. See
  [Hub operations](../operations/hub.md).
- **Backups**: `entries.jsonl` is the source of truth.
  Snapshot it to the same backup tier as your other
  internal data.
- **Upgrades**: cadence the team agrees on. Major
  upgrades may require everyone to re-register, so do
  them at natural breaks.
- **Failures**: see
  [Hub failure modes](../operations/hub-failure-modes.md)
  for the standard oncall playbook.

**Optional but recommended**: run a 3-node Raft cluster
so the hub survives individual node failures. See
[HA cluster](hub-cluster.md). For teams under 10 people,
a single-node hub with daily backups is usually fine.

## Token management

Every team member has a client token stored in their
`.context/.connect.enc`. Rules of thumb:

- **One token per engineer per project.** Not one token
  per team; not one shared token. Each engineer
  registers each of their working projects separately.
- **Token compromise = revoke immediately.** When an
  engineer leaves, their tokens should be removed from
  `clients.json` on the hub. This is a manual operation
  today; see [Hub security](../security/hub.md) for the
  revocation steps.
- **No checked-in tokens.** `.context/.connect.enc` is
  encrypted with the local machine key, but don't push
  it to shared repos — it's per-workstation.

## What this recipe is *not*

**Not a wiki replacement.** The hub is for structured
entries, not prose. Put your architecture overviews,
onboarding docs, and design discussions in a real wiki.

**Not an audit log.** `Origin` on the hub is
self-asserted. If compliance requires provenance, the
hub is the wrong tool.

**Not a ticket system.** Task sharing works, but
mature teams already have Jira/Linear/Github Issues.
Don't try to replace those with hub tasks — use the
hub for lightweight cross-project todos that your
existing tracker doesn't capture well.

**Not a production service for end users.** This is
internal team infrastructure. Do not expose the hub to
customers, partners, or the open internet.

## See also

- [Hub overview](hub-overview.md) — when to use the
  hub and when not to.
- [Personal cross-project brain](hub-personal.md) —
  the single-developer companion recipe.
- [Multi-machine setup](hub-multi-machine.md) —
  standing up the hub on a LAN host.
- [HA cluster](hub-cluster.md) — optional redundancy
  for larger teams.
- [Hub operations](../operations/hub.md) — backup,
  rotation, monitoring.
- [Hub security](../security/hub.md) — threat model
  and hardening checklist.
