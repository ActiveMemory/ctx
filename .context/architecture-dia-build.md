# Build, Dependencies, and File Layout (Diagrams)

Parent: [ARCHITECTURE.md](ARCHITECTURE.md)

## External Dependencies

```
  go.mod (3 direct dependencies):
  ┌──────────────────────────────────────────────────────┐
  │  github.com/fatih/color     — terminal coloring      │
  │  github.com/spf13/cobra     — CLI framework          │
  │  gopkg.in/yaml.v3           — YAML parsing           │
  └──────────────────────────────────────────────────────┘

  External tools (optional, not Go dependencies):
  ┌──────────────────────────────────────────────────────┐
  │  zensical    — static site generation (journal, docs)│
  │  gpg         — commit signing                        │
  └──────────────────────────────────────────────────────┘
```

## Build & Release Pipeline

```
  Developer workstation             GitHub Actions
  ┌──────────────────────┐         ┌──────────────────────┐
  │ make build           │         │ ci.yml               │
  │   CGO_ENABLED=0      │         │   go build           │
  │   -ldflags version   │         │   go test -v ./...   │
  │                      │         │   go vet             │
  │ make audit           │         │   golangci-lint      │
  │   gofmt              │         │   DCO check (PRs)    │
  │   go vet             │         │                      │
  │   golangci-lint      │         │ release.yml          │
  │   lint-drift.sh      │         │   (on v* tag)        │
  │   lint-docs.sh       │         │   test + build-all   │
  │   go test ./...      │         │   6 platform binaries│
  │                      │         │   GitHub release     │
  │ make smoke           │         └──────────────────────┘
  │   Integration tests  │
  │                      │
  │ hack/release.sh      │
  │   VERSION bump       │
  │   release notes      │
  │   plugin version sync│
  │   test + smoke       │
  │   build-all          │
  │   signed git tag     │
  └──────────────────────┘

  Build targets: darwin/amd64, darwin/arm64,
                 linux/amd64, linux/arm64,
                 windows/amd64, windows/arm64
```

## File Layout

```
ctx/
├── cmd/ctx/                     # Main entry point (main.go)
├── internal/
│   ├── bootstrap/               # CLI initialization, command registration
│   ├── claude/                  # Claude Code hooks, skills, settings types
│   ├── cli/                     # 23 command packages
│   │   ├── add/                 #   ctx add
│   │   ├── agent/               #   ctx agent
│   │   ├── compact/             #   ctx compact
│   │   ├── complete/            #   ctx complete
│   │   ├── decision/            #   ctx decision
│   │   ├── drift/               #   ctx drift
│   │   ├── hook/                #   ctx hook
│   │   ├── initialize/          #   ctx init
│   │   ├── journal/             #   ctx journal
│   │   ├── learnings/           #   ctx learnings
│   │   ├── load/                #   ctx load
│   │   ├── loop/                #   ctx loop
│   │   ├── notify/              #   ctx notify
│   │   ├── pad/                 #   ctx pad
│   │   ├── permissions/         #   ctx permissions
│   │   ├── recall/              #   ctx recall
│   │   ├── remind/              #   ctx remind
│   │   ├── serve/               #   ctx serve
│   │   ├── status/              #   ctx status
│   │   ├── sync/                #   ctx sync
│   │   ├── system/              #   ctx system
│   │   ├── task/                #   ctx task
│   │   └── watch/               #   ctx watch
│   ├── config/                  # Constants, regex, file names, read order
│   ├── context/                 # Context loading, token estimation
│   ├── crypto/                  # AES-256-GCM encryption, key management
│   ├── drift/                   # Drift detection engine (7 checks)
│   ├── index/                   # Index table generation for DECISIONS/LEARNINGS
│   ├── journal/
│   │   └── state/               # Journal processing pipeline state
│   ├── notify/                  # Webhook notifications, encrypted URL storage
│   ├── rc/                      # Runtime config (.ctxrc, env, CLI flags)
│   ├── recall/
│   │   └── parser/              # Session transcript parsing (JSONL, Markdown)
│   ├── sysinfo/                 # OS resource metrics (memory, disk, load)
│   ├── task/                    # Task checkbox parsing
│   ├── assets/                  # Embedded templates (go:embed)
│   │   ├── claude/
│   │   │   ├── .claude-plugin/  #   Plugin manifest (plugin.json)
│   │   │   ├── hooks/           #   Hook definitions (hooks.json)
│   │   │   └── skills/          #   30 skill templates (*/SKILL.md)
│   │   ├── entry-templates/     #   Decision/learning entry templates
│   │   ├── ralph/               #   Ralph autonomous loop PROMPT.md
│   │   └── tools/               #   Helper scripts (cleanup, watch)
│   └── validation/              # Input sanitization, path boundary checks
├── docs/                        # Documentation site source
├── site/                        # Generated static site (zensical)
├── hack/                        # Build/release scripts, runbooks
├── editors/vscode/              # VS Code extension (@ctx chat participant)
├── specs/                       # Feature specifications
├── .context/                    # This project's own context directory
│   ├── CONSTITUTION.md          #   Inviolable rules
│   ├── TASKS.md                 #   Current work items
│   ├── CONVENTIONS.md           #   Code patterns and standards
│   ├── ARCHITECTURE.md          #   This file
│   ├── DECISIONS.md             #   Architectural decisions
│   ├── LEARNINGS.md             #   Gotchas and tips
│   ├── GLOSSARY.md              #   Domain terms
│   ├── DETAILED_DESIGN.md       #   Deep per-module reference
│   ├── AGENT_PLAYBOOK.md        #   Meta instructions for AI agents
│   ├── journal/                 #   Exported session transcripts
│   ├── sessions/                #   Session snapshots
│   └── archive/                 #   Archived tasks
├── .claude/                     # Claude Code integration
│   ├── settings.local.json      #   Hooks and permissions
│   └── skills/                  #   Live skill definitions (30 skills)
├── .claude-plugin/              # Plugin marketplace manifest
├── Makefile                     # Build, test, lint, release targets
├── VERSION                      # Single source of truth (0.7.0)
└── go.mod                       # Go 1.25.6, 3 direct dependencies
```
