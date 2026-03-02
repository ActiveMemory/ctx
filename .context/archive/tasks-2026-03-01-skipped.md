- [x] Add ctx help command; use-case-oriented cheat sheet for lazy CLI users.
  Should cover: (1) core CLI commands grouped by workflow (getting started,
  tracking decisions, browsing history, AI context), (2) available
  slash-command skills with one-line descriptions, (3) common workflow
  recipes showing how commands and skills combine. One screen,
  no scrolling. Not a skill; a real CLI command. #added:2026-02-06-184257 #done:2026-02-28
- [-] Investigate ctx init overwriting user-generated content in .context/
  files. Commit a9df9dd wiped 18 decisions from DECISIONS.md, replacing with
  empty template. — Already fixed: runInit skips existing files unless --force
  is passed (per-file stat check at run.go:90), and prompts for confirmation
  when essential files exist (run.go:51-64). #added:2026-02-06-182205
- [-] Session pattern analysis skill — rejected. Automated pattern capture from 
  sessions risks training the agent to please rather than push back. Existing 
  mechanisms (learnings, hooks, constitution) already capture process preferences 
  explicitly. See LEARNINGS.md. #added:2026-02-22-212143
- [-] Suppress context checkpoint nudges after wrap-up — replaced by Phase 0.9 
  breakdown below #added:2026-02-24-205402
