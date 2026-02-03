# Tasks — Context CLI

## Phase 1.a: Cleanup and Release

- [ ] Large session exports crash Firefox - Session 15bcbe62 (1710 interactions, 36K lines)
      causes browser warnings even with collapsible `<details>` sections. Investigate:
      1. Split into multiple pages after N interactions (e.g., 200)
      2. Add prev/next navigation links between parts
      3. Consider lazy loading or pagination approaches
      Test case: http://localhost:8000/2026-01-30--15bcbe62/
      #priority:high #added:2026-02-03-070659
- [x] `ctx recall export --all --force` missing sessions bug - Fixed two issues:
      1. CanParse required slug field which newer Claude Code (v2.1.29+) no longer includes
      2. Subagent sessions in /subagents/ dirs share parent sessionId, causing wrong content
      #fixed:2026-02-03
- [ ] T1.2.9: upstream CI is broken (again)
- [x] T1.2.13: Compose two blog posts: 1) what has changed after the human-guided
      refactoring, and what we can learn about this.
      2) what has happened since the last release cut.
- [ ] T1.2.14: There is no documentation about the session site.
- [ ] T1.2.15: Add integration tests for CLI commands (drift, sync, decision,
      learnings, serve, recall) - test actual file operations and command execution
      #added:2026-02-01-062541

## Phase 1.b: Init Improvements

- [ ] T1.3.1: Add `ctx init --ralph` flag — creates Ralph Loop infrastructure
      (PROMPT.md, IMPLEMENTATION_PLAN.md with strict one-task-exit discipline).
      Default `ctx init` remains conversational (no "never ask questions" directive).
      #added:2026-02-03

## Phase 2: Cross-Session Monitoring (`ctx monitor`)

Enable one agent/process to inform another about context health and audits.

### Phase 2.0: Infrastructure

- [x] T2.0.1: Create specs/monitor-architecture.md — overall design
- [x] T2.0.2: Create specs/active-sessions.md — tombstone + heartbeat lifecycle spec
- [x] T2.0.3: Create specs/context-health.md — health metrics, repetition detection
- [x] T2.0.4: Create specs/auditors.md — programmatic vs semantic auditor system
- [x] T2.0.5: Create specs/signals.md — atomic writes, TTL, hysteresis, session-scoped dirs

### Phase 2.1: Active Session Tracking

- [ ] T2.1.1: Create SessionStart hook — writes tombstone to .context/active-sessions/
- [ ] T2.1.2: Create SessionEnd hook — removes tombstone from .context/active-sessions/
- [ ] T2.1.3: Implement internal/monitor/session/tracker.go — GetActiveSessions()
- [ ] T2.1.4: Add orphan detection — sessions with stale/missing transcripts

### Phase 2.2: Context Health Analysis

- [ ] T2.2.1: Implement internal/monitor/health/analyzer.go — parse JSONL, compute metrics
- [ ] T2.2.2: Token estimation — estimate context % from message content lengths
- [ ] T2.2.3: Repetition detection — hash recent content, detect loops
- [ ] T2.2.4: Turn count and activity tracking

### Phase 2.3: Signal System

- [ ] T2.3.1: Implement internal/monitor/signal/writer.go — WriteSignal(), ClearSignal()
- [ ] T2.3.2: Create UserPromptSubmit hook — reads .context/signals/, injects into context
- [ ] T2.3.3: Signal cleanup — auto-remove after delivery

### Phase 2.4: Monitor Command

- [ ] T2.4.1: Create cmd/ctx/monitor.go — Cobra subcommand skeleton
- [ ] T2.4.2: Implement watch loop — poll active sessions, run health checks
- [ ] T2.4.3: Add --interval flag (default 30s)
- [ ] T2.4.4: Add --once flag (run once and exit)

### Phase 2.5: Pluggable Auditors

- [ ] T2.5.1: Define Auditor interface in internal/monitor/auditor/
      - SessionAuditor: CheckSession(session, health) []Alert
      - ProjectAuditor: CheckProject(sessions, projectDir) []Alert
- [ ] T2.5.2: Implement ContextHealthAuditor — warn on high usage, repetition (with hysteresis)
      - At ~70% context: soft nudge "Context getting full, consider wrapping up soon"
      - At ~85% context: clear guidance "Exit this session and start fresh. Run /exit or
        Ctrl+C, then `claude` to continue. ctx remembers across sessions."
      - Detect repetition loops and suggest session restart
- [ ] T2.5.3: Implement AuditReminderAuditor — soft nudge for semantic audits
      - Track last audit times in .context/audit-state.json
      - Remind after configurable thresholds (tasks: 3d, decisions: 7d, specs: 14d)
- [ ] T2.5.4: Add --auditors flag to select which auditors run
- [ ] T2.5.5: Alert identity and dedupe — stable keys, cooldown tracking, state persistence

### Phase 2.6: Semantic Audit Slash Commands

These are LLM-powered audits triggered by user, not continuous monitoring.

- [ ] T2.6.1: Create /ctx-audit-tasks skill — check if incomplete tasks are actually done
- [ ] T2.6.2: Create /ctx-audit-decisions skill — check if code aligns with recorded decisions
- [ ] T2.6.3: Create /ctx-audit-specs skill — check if implementation matches specs/
- [ ] T2.6.4: Update .context/audit-state.json when audits run (for reminder tracking)

## Backlog

Parent command for working with exported sessions (.context/journal/):

Subcommands to explore:
- ctx journal enrich: Add frontmatter/tags (topics, type, outcome, key files)
- ctx journal cluster: Group related sessions, build continuation chains
- ctx journal summarize: Generate timeline summaries, feature narratives
- ctx journal analyze: Find patterns (recurring mistakes, revisited decisions, coupling)
- ctx journal brief <topic>: Generate compressed context packet for a topic
- ctx journal site: Generate static site via zensical (browse, search, timeline)
Additional supporting context:
```text
  Enrichment                                                                                                                                                               
  - Add frontmatter: topics, type (feature/bugfix/exploration), outcome, key files                                                                                         
  - Auto-tag: technologies, libraries, error types                                                                                                                         
  - Extract: decisions made, learnings discovered, tasks completed                                                                                                         
                                                                                                                                                                           
  Organization                                                                                                                                                             
  - Cluster related sessions (same feature across days)                                                                                                                    
  - Build continuation chains ("Part 1 → Part 2 → Part 3")                                                                                                                 
  - Create topic indexes ("All auth-related sessions")                                                                                                                     
                                                                                                                                                                           
  Synthesis                                                                                                                                                                
  - Timeline summaries ("What happened this week")                                                                                                                         
  - Feature narratives ("How we built X" from 5 sessions)                                                                                                                  
  - Decision trails (link decisions to sessions that made them)                                                                                                            
  - FAQ generation from common questions asked                                                                                                                             
                                                                                                                                                                           
  Analysis                                                                                                                                                                 
  - Find recurring mistakes → suggest new learnings                                                                                                                        
  - Detect revisited decisions → smell for bad choices                                                                                                                     
  - Identify files that change together → coupling detection                                                                                                               
  - Time patterns → what takes longer than expected                                                                                                                        
                                                                                                                                                                           
  Context compression                                                                                                                                                      
  - Generate "briefing docs" by topic for future sessions                                                                                                                  
  - "Everything you need to know about the auth system" distilled from 10 sessions                                                                                         
                                                                                                                                                                           
  Static site (zensical)                                                                                                                                                   
  - Browse by date/topic/tag                                                                                                                                               
  - Search across all sessions                                                                                                                                             
  - Related sessions sidebar                                                                                                                                               
  - Timeline visualization                                                                                                                                                 
                                                                                                                                                                           
  Meta/training                                                                                                                                                            
  - Extract good prompt patterns                                                                                                                                           
  - Document what clarifications were needed                                                                                                                               
  - Build project-specific agent guidance 
```

Depends on: ctx recall export (Phase 2)
#priority:low #phase:future #added:2026-01-28-071638


- [ ] feat: ctx enrich - retroactively expand sparse context entries

Finds one-liner learnings/decisions and expands them:
1. Locate sparse entries (missing Context/Lesson/Application)
2. Find originating session via timestamp correlation
3. Read surrounding context from that session
4. Generate full structured entry for human review

Could run as:
- ctx enrich --learnings (expand sparse learnings)
- ctx enrich --decisions (expand sparse decisions)
- ctx enrich --all (both)
- ctx enrich --dry-run (show what would be expanded)

#priority:low #phase:future #added:2026-01-28-073058

- [ ] feat: make config constants configurable via .contextrc

Some hardcoded constants in internal/config/config.go could be user-configurable:
- MaxDecisionsToSummarize (default 3)
- MaxLearningsToSummarize (default 5)
- MaxPreviewLen (default 60)
- WatchAutoSaveInterval (default 5)

Follow the pattern established for token_budget and archive_after_days in internal/rc.
#priority:low #phase:future #added:2026-01-31

- [ ] explore: `ctx recall stats` - analytics/statistics
  - Token usage over time, tool patterns, session durations
  - Explore when we have a clear use case
    #priority:deferred

- [ ] feat: `ctx recall search <query>` - CLI-based search across sessions
  - Simple text search, no server needed
  - Support phrase/literal matching with quotes: `ctx recall search "version history"`
  - Zensical's web search only does term matching, CLI should be more powerful
  - IDE grep is alternative, this is convenience
    #priority:low
