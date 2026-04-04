# Implementation Plan: Hooks and Steering

## Overview

Incremental implementation of the Hooks & Steering system for ctx. Each task builds on previous work, starting with foundational types and domain logic, then CLI commands, then MCP extensions, and finally integration wiring. All new packages follow existing `internal/<domain>/` and `internal/cli/<cmd>/` conventions. Go is the implementation language throughout.

## Tasks

- [x] 1. Steering domain package — types and parser
  - [x] 1.1 Create `internal/steering/types.go` with `SteeringFile`, `InclusionMode`, and `SyncReport`
    - Define `InclusionMode` enum (`always`, `auto`, `manual`)
    - Define `SteeringFile` struct with `Name`, `Description`, `Inclusion`, `Tools`, `Priority`, `Body`, `Path`
    - Define `SyncReport` struct with `Written`, `Skipped`, `Errors`
    - Add `doc.go` for the package
    - _Requirements: 1.1, 1.2, 1.5_

  - [x] 1.2 Implement `internal/steering/parse.go` — frontmatter parser and printer
    - `Parse(data []byte, filePath string) (*SteeringFile, error)` — extract YAML frontmatter delimited by `---` and markdown body
    - `Print(sf *SteeringFile) []byte` — serialize back to frontmatter + markdown
    - Apply defaults: `inclusion` → `manual`, `tools` → nil (all), `priority` → 50
    - Return descriptive error on invalid YAML identifying file path and failure
    - Use `gopkg.in/yaml.v3` for YAML parsing
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7_

  - [x] 1.3 Write property test for steering parser round-trip
    - **Property 1: Round-trip consistency** — `Parse(Print(Parse(data))) == Parse(data)` for all valid inputs
    - **Validates: Requirements 1.8, 19.1**

  - [x] 1.4 Implement `internal/steering/filter.go` — inclusion mode filtering
    - `LoadAll(steeringDir string) ([]*SteeringFile, error)` — read all `.md` files and parse
    - `Filter(files, prompt, manualNames, tool string) []*SteeringFile` — apply inclusion rules
    - `always` files included unconditionally; `auto` files included on description substring match; `manual` only when named
    - Sort by ascending priority, then alphabetically by name on tie
    - Filter out files whose `Tools` list excludes the given tool
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

  - [x] 1.5 Write unit tests for steering filter
    - Test each inclusion mode independently
    - Test priority ordering and alphabetical tie-breaking
    - Test tool filtering with explicit tools list and empty (all) tools
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

- [x] 2. Steering sync to tool-native formats
  - [x] 2.1 Implement `internal/steering/sync.go` — tool-native format sync
    - `SyncTool(steeringDir, projectRoot, tool string) (SyncReport, error)` — write steering files to tool-native directory
    - `SyncAll(steeringDir, projectRoot string) (SyncReport, error)` — sync to all supported tools
    - Cursor: `.cursor/rules/<name>.mdc` with Cursor-compatible frontmatter
    - Cline: `.clinerules/<name>.md` plain markdown
    - Kiro: `.kiro/steering/<name>.md` with Kiro frontmatter
    - Skip files whose `tools` list excludes the target tool
    - Skip files whose content hasn't changed (idempotent)
    - Validate output paths resolve within project root boundary
    - _Requirements: 5.3, 5.4, 5.5, 5.8, 5.9, 15.3_

  - [x] 2.2 Write property test for steering sync idempotence
    - **Property 2: Sync idempotence** — running `SyncTool` twice produces identical output files
    - **Validates: Requirements 19.4**

- [x] 3. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 4. Hook domain package — types, discovery, runner, and security
  - [x] 4.1 Create `internal/hook/types.go` with `HookType`, `HookInput`, `HookOutput`, `HookSession`, `HookInfo`, `AggregatedOutput`
    - Define `HookType` constants: `pre-tool-use`, `post-tool-use`, `session-start`, `session-end`, `file-save`, `context-add`
    - `ValidHookTypes()` returns all valid hook type strings
    - Add `doc.go` for the package
    - _Requirements: 6.1, 7.1, 7.2_

  - [x] 4.2 Implement `internal/hook/security.go` — symlink and boundary validation
    - `ValidateHookPath(hooksDir, hookPath string) error` — reject symlinks, validate boundary, check executable bit
    - Reuse patterns from `internal/validate/` (boundary check, `os.Lstat` symlink rejection)
    - _Requirements: 15.1, 15.2, 15.4_

  - [x] 4.3 Implement `internal/hook/discover.go` — hook discovery
    - `Discover(hooksDir string) (map[HookType][]HookInfo, error)` — find all hook scripts grouped by type
    - `FindByName(hooksDir, name string) (*HookInfo, error)` — search all type directories for a hook
    - Skip non-executable scripts with logged warning
    - Skip symlinks (security)
    - Return empty map if hooks directory doesn't exist
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 15.1, 15.2_

  - [x] 4.4 Write unit tests for hook discovery
    - Test discovery with mixed executable/non-executable scripts
    - Test symlink rejection
    - Test missing hooks directory returns empty map
    - Test alphabetical ordering within each hook type
    - _Requirements: 6.1, 6.2, 6.3, 6.4_

  - [x] 4.5 Implement `internal/hook/runner.go` — hook execution engine
    - `RunAll(hooksDir string, hookType HookType, input *HookInput, timeout time.Duration) (*AggregatedOutput, error)`
    - Pass `HookInput` as JSON via stdin, read `HookOutput` as JSON from stdout
    - If hook returns `cancel: true`, halt and return cancellation message
    - Append non-empty `context` fields to aggregated context
    - On non-zero exit: log error, skip hook, continue
    - On invalid JSON stdout: log warning with hook file and parse error, continue
    - Enforce configurable timeout (default 10s); terminate on exceed with logged warning
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8_

  - [x] 4.6 Write unit tests for hook runner
    - Test cancel propagation halts subsequent hooks
    - Test context aggregation from multiple hooks
    - Test non-zero exit code handling (skip and continue)
    - Test invalid JSON output handling (skip and continue)
    - Test timeout enforcement terminates script
    - **Validates: Requirements 7.3, 7.4, 7.5, 7.6, 7.7, 7.8, 19.6**

- [x] 5. Skill domain package
  - [x] 5.1 Create `internal/skill/types.go` with `Skill` struct
    - Define `Skill` struct with `Name`, `Description`, `Body`, `Dir`
    - Add `doc.go` for the package
    - _Requirements: 12.1_

  - [x] 5.2 Implement `internal/skill/load.go` — skill loading
    - `LoadAll(skillsDir string) ([]*Skill, error)` — read all installed skills
    - `Load(skillsDir, name string) (*Skill, error)` — read single skill by name
    - Parse `SKILL.md` frontmatter (`name`, `description`) and markdown body
    - _Requirements: 12.1, 12.4_

  - [x] 5.3 Implement `internal/skill/install.go` and `internal/skill/remove.go`
    - `Install(source, skillsDir string) (*Skill, error)` — copy skill from source, validate `SKILL.md` exists
    - `Remove(skillsDir, name string) error` — delete skill directory
    - Return error if source has no valid `SKILL.md`
    - Return error if skill name not found on remove
    - _Requirements: 12.2, 12.3, 12.5, 12.6_

- [x] 6. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 7. RC package extensions — tool, steering, hooks config
  - [x] 7.1 Add `Tool`, `Steering`, and `Hooks` fields to `CtxRC` in `internal/rc/types.go`
    - Add `Tool string` field with `yaml:"tool"` tag
    - Add `Steering *SteeringRC` and `Hooks *HooksRC` structs
    - `SteeringRC`: `Dir`, `DefaultInclusion`, `DefaultTools`
    - `HooksRC`: `Dir`, `Timeout`, `Enabled`
    - _Requirements: 11.1, 16.1, 16.2_

  - [x] 7.2 Add accessor functions for new RC fields
    - `Tool() string`, `SteeringDir() string`, `HooksDir() string`, `HookTimeout() int`, `HooksEnabled() bool`
    - Apply defaults: steering dir → `.context/steering`, hooks dir → `.context/hooks`, timeout → 10, enabled → true
    - Preserve existing priority hierarchy: CLI flags > env vars > `.ctxrc` > defaults
    - _Requirements: 11.2, 16.3, 16.4_

  - [x] 7.3 Write unit tests for RC tool field resolution
    - Test that `Tool()` returns the configured value
    - Test default values for steering and hooks config
    - Test that `HooksEnabled()` returns false when configured
    - **Validates: Requirements 19.8**

- [x] 8. Add `--tool` persistent flag to root command
  - [x] 8.1 Add `--tool` persistent flag in `internal/bootstrap/cmd.go`
    - Add `ResolveTool(cmd *cobra.Command) (string, error)` helper that reads `--tool` flag, falls back to `rc.Tool()`
    - Return error if neither is set and command requires a tool
    - _Requirements: 11.2, 11.3, 11.6_

- [x] 9. Steering CLI commands
  - [x] 9.1 Create `internal/cli/steering/` package with parent command and `doc.go`
    - `Cmd() *cobra.Command` returning `ctx steering` parent with subcommands
    - Follow existing `internal/cli/<cmd>/` conventions
    - _Requirements: 3.1_

  - [x] 9.2 Implement `ctx steering add <name>` subcommand
    - Create `.context/steering/<name>.md` with default frontmatter and empty body
    - Error if file already exists
    - Error if `.context/` directory does not exist
    - _Requirements: 3.1, 3.2, 3.7_

  - [x] 9.3 Implement `ctx steering list` subcommand
    - Display all steering files with name, inclusion mode, priority, and target tools
    - _Requirements: 3.3_

  - [x] 9.4 Implement `ctx steering preview <prompt>` subcommand
    - Show which steering files would be included for the given prompt text
    - Respect inclusion mode rules using `steering.Filter`
    - _Requirements: 3.4, 19.3_

  - [x] 9.5 Implement `ctx steering init` subcommand
    - Generate foundation files (`product.md`, `tech.md`, `structure.md`, `workflow.md`) in `.context/steering/`
    - Skip existing files and report which were skipped
    - _Requirements: 3.5, 3.6_

  - [x] 9.6 Implement `ctx steering sync` subcommand
    - Without `--tool` or `--all`: read `tool` from `.ctxrc`, sync to that tool's format
    - `--tool <tool>`: sync to specified tool format; error on unsupported tool
    - `--all`: sync to all supported tool formats
    - Error if no tool specified and no `tool` field in `.ctxrc`
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7_

- [x] 10. Hook CLI commands
  - [x] 10.1 Create `internal/cli/hook/` package with parent command and `doc.go`
    - `Cmd() *cobra.Command` returning `ctx hook` parent with subcommands
    - _Requirements: 8.1_

  - [x] 10.2 Implement `ctx hook add <hook-type> <name>` subcommand
    - Create executable script template at `.context/hooks/<hook-type>/<name>.sh`
    - Include shebang, JSON input reading, JSON output structure
    - Error on invalid hook type listing valid types
    - _Requirements: 8.1, 8.2_

  - [x] 10.3 Implement `ctx hook list` subcommand
    - Display all hooks grouped by hook type with name, enabled/disabled status, file path
    - _Requirements: 8.3_

  - [x] 10.4 Implement `ctx hook test <hook-type>` subcommand
    - Accept `--tool` and `--path` flags
    - Construct mock `HookInput`, execute enabled hooks, display `HookOutput` from each
    - _Requirements: 8.4, 19.2_

  - [x] 10.5 Implement `ctx hook enable <name>` and `ctx hook disable <name>` subcommands
    - `enable`: add executable permission bit
    - `disable`: remove executable permission bit
    - Error if hook name not found
    - _Requirements: 8.5, 8.6, 8.7_

- [x] 11. Skill CLI commands
  - [x] 11.1 Create `internal/cli/skill/` package with parent command and `doc.go`
    - `Cmd() *cobra.Command` returning `ctx skill` parent with subcommands
    - _Requirements: 12.1_

  - [x] 11.2 Implement `ctx skill install <source>` subcommand
    - Download or copy skill from source into `.context/skills/<name>/`
    - Error if source has no valid `SKILL.md`
    - _Requirements: 12.2, 12.3_

  - [x] 11.3 Implement `ctx skill list` subcommand
    - Display all installed skills with name and description
    - Return empty list without error if `.context/skills/` doesn't exist
    - _Requirements: 12.4, 14.4_

  - [x] 11.4 Implement `ctx skill remove <name>` subcommand
    - Delete `.context/skills/<name>/` directory
    - Error if skill name not found
    - _Requirements: 12.5, 12.6_

- [x] 12. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 13. Bootstrap registration — wire new CLI commands
  - [x] 13.1 Register `steering`, `hook`, and `skill` commands in `internal/bootstrap/group.go`
    - Add `steering.Cmd` and `hook.Cmd` to `integrations()` group
    - Add `skill.Cmd` to `contextCmds()` group
    - Import new CLI packages
    - _Requirements: 3.1, 8.1, 12.1_

- [x] 14. Extend `ctx init` to create new directories
  - [x] 14.1 Extend `internal/cli/initialize/` to create `.context/steering/`, `.context/hooks/`, and `.context/skills/`
    - Create directories with `0755` permissions
    - Skip existing directories without error
    - _Requirements: 13.1, 13.2, 13.3_

- [x] 15. Extend `ctx agent` for steering and skill integration
  - [x] 15.1 Add steering file inclusion to `ctx agent` context packet assembly
    - After existing tiers, add Tier 6 for steering files (from remaining budget)
    - Include `always` files unconditionally, `auto` files on prompt match
    - Respect token budget — truncate/omit lower-priority steering files when exceeded
    - Produce same output when no `.context/steering/` exists
    - _Requirements: 4.1, 4.2, 4.3, 14.1_

  - [x] 15.2 Add `--skill <name>` flag to `ctx agent`
    - Include named skill's `SKILL.md` content in context packet as Tier 7
    - Error if skill name not found
    - _Requirements: 12.7, 12.8_

- [x] 16. MCP server extensions — new tools
  - [x] 16.1 Implement `SteeringGet` and `Search` handler methods in `internal/mcp/handler/`
    - `SteeringGet(prompt string) (string, error)` — return applicable steering files; if no prompt, return `always` files only
    - `Search(query string) (string, error)` — search across `.context/` files, return excerpts with paths and line numbers
    - _Requirements: 9.1, 9.2, 9.3_

  - [x] 16.2 Implement `SessionStartHooks` and `SessionEndHooks` handler methods
    - `SessionStartHooks() (string, error)` — execute `session-start` hooks, return aggregated context
    - `SessionEndHooks(summary string) (string, error)` — execute `session-end` hooks with summary in parameters
    - Return success with empty context when no hooks exist
    - _Requirements: 10.1, 10.2, 10.3_

  - [x] 16.3 Register new MCP tools in catalog and dispatch
    - Add `ctx_steering_get`, `ctx_search`, `ctx_session_start`, `ctx_session_end` to `internal/mcp/server/catalog/`
    - Add JSON Schema parameter definitions for each tool
    - Add dispatch routes in `internal/mcp/server/dispatch/`
    - Unregistered tool names return JSON-RPC error `-32601`
    - _Requirements: 9.4, 9.5, 10.4_

  - [x] 16.4 Write unit tests for MCP steering and session tools
    - Test `ctx_steering_get` with and without prompt parameter
    - Test `ctx_session_start` with no hooks returns success
    - Test `ctx_session_end` passes summary to hook input
    - **Validates: Requirements 19.5**

- [x] 17. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 18. Drift detection extensions
  - [x] 18.1 Add new issue types and check names to `internal/drift/types.go`
    - `IssueInvalidTool`, `IssueHookNoExec`, `IssueStaleSyncFile`
    - `CheckSteeringTools`, `CheckHookPerms`, `CheckSyncStaleness`, `CheckRCTool`
    - _Requirements: 17.1, 17.2, 17.3, 17.4_

  - [x] 18.2 Implement new drift checks in `internal/drift/detector.go`
    - Check steering files for unsupported tool identifiers
    - Check hook scripts for missing executable permission bit
    - Check synced tool-native files are up to date vs source steering files
    - Check `.ctxrc` `tool` field for unsupported tool identifier
    - _Requirements: 17.1, 17.2, 17.3, 17.4_

  - [x] 18.3 Write unit tests for new drift checks
    - Construct `.context/` directories with known issues, verify correct warnings
    - Test each new check type independently
    - **Validates: Requirements 19.7**

- [x] 19. Backward compatibility verification
  - [x] 19.1 Verify backward compatibility across all extensions
    - Confirm `ctx agent` produces same output when no `.context/steering/` exists
    - Confirm hook runner returns empty results when no `.context/hooks/` exists
    - Confirm commands work without `tool` field when tool is not required
    - Confirm skill list returns empty when no `.context/skills/` exists
    - Confirm existing `CLAUDE.md` and `AGENTS.md` generation is unchanged
    - _Requirements: 14.1, 14.2, 14.3, 14.4, 14.5_

- [x] 20. Final checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties (round-trip, idempotence)
- All new packages include `doc.go` and `testmain_test.go` following existing conventions
- Security validation (symlink rejection, boundary checks) reuses `internal/validate/` patterns
