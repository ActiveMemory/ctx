# Auditor System

## Overview

Auditors detect issues and generate alerts. Two categories:

| Type | Method | Frequency | Cost | Examples |
|------|--------|-----------|------|----------|
| **Programmatic** | Code, heuristics | Continuous | Cheap | Context health, repetition, file mtime |
| **Semantic** | LLM judgment | On-demand | Expensive | Task drift, decision alignment, spec compliance |

`ctx monitor` runs programmatic auditors continuously.
Semantic audits are slash commands the user triggers.

## Programmatic Auditors

### Interface

```go
// Per-session auditors (context health, repetition)
type SessionAuditor interface {
    Name() string
    CheckSession(session Session, health *ContextHealth) []Alert
}

// Project-level auditors (file staleness, audit reminders)
type ProjectAuditor interface {
    Name() string
    CheckProject(sessions []Session, projectDir string) []Alert
}
```

### Built-in Auditors

| Auditor | Type | Checks |
|---------|------|--------|
| `ContextHealthAuditor` | Session | Token usage, repetition |
| `AuditReminderAuditor` | Project | Time since last semantic audit |

### Monitor Execution

```go
func (m *Monitor) runAuditors(sessions []Session) []Alert {
    var alerts []Alert

    // Per-session auditors
    for _, s := range sessions {
        if s.State != StateActive {
            continue
        }
        health := m.analyzeHealth(s.TranscriptPath)
        for _, auditor := range m.sessionAuditors {
            alerts = append(alerts, auditor.CheckSession(s, health)...)
        }
    }

    // Project auditors
    for _, auditor := range m.projectAuditors {
        alerts = append(alerts, auditor.CheckProject(sessions, m.projectDir)...)
    }

    return alerts
}
```

## Alert Identity

Every alert has a stable identity for dedupe and cooldowns:

```go
type Alert struct {
    // Identity
    Auditor     string    // "context-health", "audit-reminder"
    Code        string    // "CTX_HEALTH_85", "AUDIT_TASKS_DUE"
    SessionID   string    // Empty for project-level alerts

    // Content
    Severity    string    // "info", "warning", "critical"
    Summary     string    // What's wrong (fact, not judgment)
    Action      string    // What to do (concrete, specific)
    TTL         int       // Seconds until expiration

    // Metadata
    GeneratedAt time.Time
}

func (a Alert) Key() string {
    return fmt.Sprintf("%s:%s:%s", a.Auditor, a.SessionID, a.Code)
}
```

### Dedupe and Cooldowns

Central management via alert key:

```go
type AlertTracker struct {
    LastEmitted  map[string]time.Time
    ActiveAlerts map[string]bool
    Cooldowns    map[string]time.Duration
}

func (t *AlertTracker) ShouldEmit(alert Alert) bool {
    key := alert.Key()
    cooldown := t.Cooldowns[alert.Auditor]

    if last, ok := t.LastEmitted[key]; ok {
        if time.Since(last) < cooldown {
            return false
        }
    }

    t.LastEmitted[key] = time.Now()
    return true
}
```

### Cooldown Defaults

| Auditor | Cooldown |
|---------|----------|
| context-health | 15 min |
| audit-reminder | 4 hours |

## Human-Centered Framing

Alerts are for humans. They must be:
- **Short**: One line summary
- **Factual**: State what is, not judgment
- **Actionable**: Concrete next step
- **Non-nagging**: Users ignore repeated warnings

### Alert Format

```
{Summary}
→ {Action}
```

### Examples

**Bad:**
```
⚠️ Context is getting dangerously high at 86%
```

**Good:**
```
Context: 86% full, repetition in last 8 turns.
→ Summarize plan in 5 bullets, then /clear and reload DECISIONS.md
```

### Tone Guidelines

| Do | Don't |
|----|-------|
| "Context: 86% full" | "Context dangerously high" |
| "Last 8 turns" | "Recently" |
| "Summarize in 5 bullets" | "Consider summarizing" |
| One action | Menu of options |

## Semantic Audits (Slash Commands)

These require LLM judgment — too expensive and slow for continuous monitoring.
User triggers them explicitly.

### Available Commands

| Command | Purpose |
|---------|---------|
| `/ctx-audit-tasks` | Are incomplete tasks actually done? |
| `/ctx-audit-decisions` | Does code align with recorded decisions? |
| `/ctx-audit-specs` | Has implementation drifted from specs/? |

### Soft Nudges

Monitor doesn't run semantic audits automatically, but can remind:

```go
type AuditReminderAuditor struct{}

func (a *AuditReminderAuditor) CheckProject(sessions []Session, dir string) []Alert {
    state := loadAuditState(dir)
    var alerts []Alert

    for auditType, lastRun := range state.LastAudit {
        threshold := auditThresholds[auditType]
        if time.Since(lastRun) > threshold {
            alerts = append(alerts, Alert{
                Auditor:  "audit-reminder",
                Code:     "AUDIT_" + strings.ToUpper(auditType) + "_DUE",
                Severity: "info",
                Summary:  fmt.Sprintf("Last %s audit: %s ago", auditType, humanDuration(time.Since(lastRun))),
                Action:   fmt.Sprintf("Consider running /ctx-audit-%s", auditType),
            })
        }
    }

    return alerts
}
```

### Audit State Tracking

```json
// .context/audit-state.json
{
  "lastAudit": {
    "tasks": "2026-01-31T10:00:00Z",
    "decisions": "2026-01-29T14:00:00Z",
    "specs": "2026-01-25T09:00:00Z"
  }
}
```

Slash commands update this when they run.

### Nudge Thresholds

| Audit Type | Remind After |
|------------|--------------|
| tasks | 3 days |
| decisions | 7 days |
| specs | 14 days |

### Why Not Auto-Run Semantic Audits?

| Concern | Issue |
|---------|-------|
| Cost | API calls without explicit user action |
| Consent | "I didn't ask for this" |
| Staleness | Results from 3am stale by 9am |
| Interactivity | User can't follow up |

User opt-in via cron if they want automation:

```bash
# User's choice, not default
0 9 * * * cd /project && claude --print "/ctx-audit-tasks" > .context/audits/$(date +%F).md
```

## Adding New Auditors

### Programmatic Auditor

```go
type MyAuditor struct{}

func (a *MyAuditor) Name() string { return "my-auditor" }

func (a *MyAuditor) CheckSession(s Session, h *ContextHealth) []Alert {
    // Your logic here
    return nil
}

// Register
monitor.AddSessionAuditor(&MyAuditor{})
```

### Semantic Audit (Slash Command)

Create `.claude/commands/ctx-audit-foo.md`:

```markdown
---
description: Check if foo is aligned with bar
allowed-tools: [Read, Glob, Grep]
---

Read .context/DECISIONS.md and check if the codebase follows the recorded
decisions. Report any discrepancies.

After completing, update .context/audit-state.json with:
  "foo": "{current timestamp}"
```

## Configuration

`.ctxrc`:

```yaml
monitor:
  interval: 30s
  inject_min_severity: warning

  auditors:
    context-health:
      enabled: true
      cooldown: 15m
    audit-reminder:
      enabled: true
      cooldown: 4h

  thresholds:
    context_warn: 85
    context_critical: 95
    repetition: 0.5
```
