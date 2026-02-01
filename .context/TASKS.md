- [ ] Add integration tests for CLI commands (drift, sync, decision, learnings, serve, recall) - test actual file operations and command execution #added:2026-02-01-062541

# Tasks — Context CLI

# Tasks

## Phase 0: Cleanup from the previous version

(All tasks archived)

### Phase 1: Parser (DONE)

## Phase 1.a: Cleanup and Release

- [x] T1.2.0.1b Add index markers to DECISIONS.md and LEARNINGS.md templates
      so new projects start with index structure #priority:low #added:2026-01-29
- [x] T1.2.5 feat: implement `--context-dir` global flag to override context directory path
  Documented in cli-reference.md as planned. Should allow `ctx --context-dir /path status`.
  #priority:low #added:2026-01-28
- [x] T1.2.6 feat: implement `--quiet` global flag to suppress non-essential output
  Documented in cli-reference.md as planned.
  #priority:low #added:2026-01-28
- [x] T1.2.7 feat: implement `--no-color` global flag to disable colored output
  Documented in cli-reference.md as planned. Currently `NO_COLOR=1` env var works.
  #priority:low #added:2026-01-28
- [x] Write AST-based test that warns if CLI functions use fmt.Print* instead of
  cmd.Print* #added:2026-01-29-171351 #done:2026-01-31
- [x] feat: `ctx recall export` - export sessions to editable journal files
  - `ctx recall export <session-id>` - export one session
  - `ctx recall export --all` - export all sessions
  - Skip existing files (user may have edited), `--force` to overwrite
  - Output to `.context/journal/YYYY-MM-DD-slug-shortid.md`
    #added:2026-01-28 #done:2026-01-31
- [x] feat: ctx journal - LLM-powered session analysis and synthesis
  - [x] ctx journal site - generate zensical static site #done:2026-01-31
  - [x] ctx serve - convenience wrapper for zensical serve #done:2026-01-31
  - [x] /ctx-journal-enrich - slash command for frontmatter/tags #done:2026-01-31
  - [x] /ctx-journal-summarize - slash command for timeline summaries #done:2026-01-31
- [ ] T1.2.9: upstream CI is broken (again)
- [ ] T1.2.10: Human code review
- [ ] T1.2.11: Human to read all user-facing documentation and update as needed.
- [ ] T1.2.12: cut a release (version number is already bumped)
- [ ] T1.2.13: Compose two blog posts: 1) what has changed after the human-guided
      refactoring, and what we can learn about this.
      2) what has happened since the last release cut.

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

- [x] feat: /ctx-blog slash command - generate blog post draft from recent activity
  - [x] /ctx-blog - from recent activity (sessions, commits, decisions)
  - [x] /ctx-blog-changelog - from commit range with theme
  #added:2026-01-28-072625 #done:2026-01-31

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
  - IDE grep is alternative, this is convenience
    #priority:low