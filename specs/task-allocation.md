# Spec: Task Allocation Across AI Agents (`ctx fleet`)

Route tasks to the best-fit AI agent based on context window size,
token budget, and capabilities. ctx becomes the orchestrator that
knows what work exists and which agent should handle it.

---

## Problem

A developer using Claude Code, Cline, Cursor, and Copilot has four
AI agents — each with different context windows, usage limits, and
strengths. Today each agent works in isolation. There is no way to
route tasks intelligently, track which agent is working on what, or
avoid wasting expensive credits on simple work.

## Solution

**ctx fleet** — a task allocation system built on the shared context
hub. One ctx instance acts as the host, others connect as agents.
The host classifies tasks, allocates them, and publishes assignments
through the hub.

```
                    ┌─────────────────────┐
                    │   ctx hub (:9900)    │
                    │   allocation engine  │
                    └──────┬──────────────┘
                           │ gRPC
            ┌──────────────┼──────────────┐
            │              │              │
    ┌───────┴──────┐ ┌─────┴──────┐ ┌─────┴──────┐
    │ ctx (Claude) │ │ ctx (Cline)│ │ ctx (Cursor)│
    │ 1M context   │ │ 200k ctx   │ │ 200k ctx   │
    │ 5M tokens/d  │ │ 1M tokens/d│ │ 1.5M tok/d │
    └──────────────┘ └────────────┘ └────────────┘
```

**Key principles:**

1. **Hub-based** — reuses the shared context hub
2. **Automatic classification** — ctx estimates task complexity
   via token estimation
3. **Human authority** — plan/dispatch split, manual override
4. **Tiered credit tracking** — real API where available,
   session counting as fallback

---

## Agent Registry

Each agent is registered in `.context/fleet.yaml` with its context
window, daily token budget, tracking mode, and capabilities.

```yaml
agents:
  - id: claude-pro
    tool: claude-code
    context_window: 1000000
    budget:
      daily_tokens: 5000000
      tracking: api             # api | oauth | session
    capabilities: [architecture, multi-file, testing]

  - id: cline-free
    tool: cline
    context_window: 200000
    budget:
      daily_tokens: 1000000
      tracking: session
    capabilities: [single-file, quick-fix, testing]

  - id: cursor-team
    tool: cursor
    context_window: 200000
    budget:
      daily_tokens: 1500000
      tracking: session
    capabilities: [single-file, quick-fix, boilerplate]
```

---

## Task Classification

ctx automatically estimates the token cost of each task:

- Parse task description for file references (explicit paths or
  implicit mentions)
- Count actual tokens of referenced files on disk
- Factor in subtask count and description complexity
- Map estimate to agent context windows

| Token Estimate | Suitable Agents |
|----------------|-----------------|
| < 50k | Any agent |
| 50k–150k | 200k+ context window |
| 150k–500k | 1M context window |
| > 500k | Split task or 1M with budget management |

---

## Allocation Algorithm

Weighted best-fit: heaviest tasks assigned first.

- Filter agents by context window and remaining budget
- Score by: capability match (0.4) + budget remaining (0.3) +
  context window headroom (0.3)
- Assign to highest-scoring agent
- Human override always available via manual assignment

---

## Communication

Uses the shared context hub ([shared-context-hub.md](shared-context-hub.md)).
Two new entry types:

- `assignment` — host dispatches task to an agent
- `assignment-update` — agent reports status back

```
HOST                                AGENT
  │── Publish(assignment) ────────→ │  receives via ctx connect listen
  │                                 │  works on task...
  │←── Publish(assignment-update) ──│  reports completion
```

---

## Credit Tracking

### Tiered Approach

Each provider exposes different levels of usage visibility:

| Provider | Access | Method |
|----------|--------|--------|
| Claude (API) | Full | Usage & Cost Admin API |
| Claude (Pro/Max) | Partial | Internal OAuth endpoint |
| Cline | Depends | Uses underlying provider's API |
| Cursor | None | Dashboard only |
| Copilot | None | Subscription-based |

ctx uses **three tracking modes**, selected per agent in fleet.yaml:

- **`api`** — queries provider usage API for real token consumption
  (Anthropic, OpenRouter)
- **`oauth`** — reads Claude Pro/Max quota from the same endpoint
  Claude Code uses internally (best-effort, may break)
- **`session`** — counts tokens from session history (fallback for
  Cursor, Copilot, anything without an API)

### Error Detection

Across all tiers, ctx also watches session logs for rate limit
signals (HTTP 429, "quota exceeded" messages, throttling gaps).
Throttled agents are excluded from allocation until next reset.

### Budget Reset

Daily auto-reset or manual. Configured in `.ctxrc`:

```yaml
fleet:
  credit_reset: daily
```

---

## CLI Commands

```bash
# Setup
ctx fleet init                          # create fleet.yaml
ctx fleet agents                        # list agents with status
ctx fleet agents add / remove <id>

# Classification & allocation
ctx fleet classify                      # show task token estimates
ctx fleet plan                          # generate allocation (dry-run)
ctx fleet dispatch                      # publish assignments to hub
ctx fleet assign <task> --agent <id>    # manual override

# Monitoring
ctx fleet status                        # show assignment states
ctx fleet credits refresh               # refresh from APIs/sessions
ctx fleet credits reset                 # reset daily budgets

# Agent side
ctx fleet report <task> --status <s>    # report completion/blocked/rejected
```

---

## Integration

- **Shared hub**: primary communication channel for assignments
- **Commit tracing**: adds `agent:<id>` to trace refs
- **Webhooks**: fleet.dispatch/complete/blocked/rejected events
- **`ctx complete`**: auto-updates fleet assignment status
