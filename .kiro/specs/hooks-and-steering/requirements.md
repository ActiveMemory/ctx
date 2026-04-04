# Requirements Document

## Introduction

This document specifies the Hooks & Steering system for ctx — a universal AI-context layer CLI tool (`github.com/ActiveMemory/ctx`). The system extends ctx from a persistence-only tool into a behavioral guidance and lifecycle automation platform that works across all major AI coding assistants (Claude Code, Cursor, Cline, Kiro, Codex).

The feature spans four phases: a unified Steering Layer, an event-driven Hooks System, enhanced MCP Server tools, and a Skills System. Each phase is independently useful and additive — nothing breaks existing usage. Cross-tool configuration is handled by extending the existing `.ctxrc` profile system (`.ctxrc.kiro`, `.ctxrc.claude`, etc.) with new `tool`, `steering`, and `hooks` sections — no separate profiles directory is needed.

All new artifacts are plain files (markdown, YAML, shell scripts) stored under `.context/`, making them git-versionable, human-readable, and tool-agnostic. ctx remains the single source of truth; other tools sync from ctx. ctx does not auto-detect which AI tool the user is using — the user specifies the target tool once via the `tool` field in `.ctxrc`, and all commands read from that. A `--tool` CLI flag is available as a one-off override.

## Glossary

- **Steering_File**: A markdown file with YAML frontmatter stored in `.context/steering/` that defines persistent behavioral rules injected into AI tool prompts.
- **Inclusion_Mode**: One of three modes (`always`, `auto`, `manual`) that determines when a Steering_File is injected into an AI prompt.
- **Hook**: An executable script in `.context/hooks/<hook-type>/` that fires at a specific lifecycle event, receives JSON via stdin, and returns JSON via stdout.
- **Hook_Type**: A lifecycle event category: `pre-tool-use`, `post-tool-use`, `session-start`, `session-end`, `file-save`, or `context-add`.
- **Hook_Runner**: The internal component that discovers, validates, and executes Hook scripts, passing the Hook_Input and reading the Hook_Output.
- **Hook_Input**: A JSON object sent to a Hook via stdin containing `hookType`, `tool`, `parameters`, `session`, `timestamp`, and `ctxVersion` fields.
- **Hook_Output**: A JSON object returned by a Hook via stdout containing `cancel` (boolean), `context` (optional string), and `message` (optional string) fields.
- **MCP_Server**: The Model Context Protocol server (`ctx mcp serve`) that exposes ctx operations as JSON-RPC 2.0 tools over stdin/stdout.
- **MCP_Tool**: A callable operation exposed by the MCP_Server (e.g., `ctx_status`, `ctx_steering_get`).
- **Profile**: A named `.ctxrc.<name>` configuration file (e.g., `.ctxrc.kiro`, `.ctxrc.claude`) that defines tool-specific settings using the existing ctx config switch mechanism.
- **Skill**: A reusable instruction bundle stored in `.context/skills/` containing a `SKILL.md` manifest with YAML frontmatter and markdown instructions.
- **Context_Packet**: The assembled markdown output produced by `ctx agent` containing prioritized context files fitted to a token budget.
- **Steering_Sync**: The process of converting Steering_Files into tool-native formats (e.g., `.cursor/rules/*.mdc`, `.clinerules/*.md`, `.kiro/steering/*.md`).
- **Foundation_Files**: Auto-generated Steering_Files (`product.md`, `tech.md`, `structure.md`, `workflow.md`) created by `ctx steering init` from codebase analysis.
- **Tool_Identifier**: A string identifying a supported AI tool: `claude`, `cursor`, `cline`, `kiro`, or `codex`.
- **Frontmatter_Parser**: The component that parses YAML frontmatter delimited by `---` from the top of markdown files.

## Requirements

### Requirement 1: Steering File Storage and Format

**User Story:** As a developer, I want to store behavioral rules as markdown files with structured frontmatter in `.context/steering/`, so that AI tools receive consistent guidance that is human-readable and version-controlled.

#### Acceptance Criteria

1. THE Steering_File SHALL contain YAML frontmatter delimited by `---` lines followed by markdown body content.
2. THE Frontmatter_Parser SHALL extract `name` (string), `description` (string), `inclusion` (one of `always`, `auto`, `manual`), `tools` (list of Tool_Identifiers), and `priority` (integer) fields from Steering_File frontmatter.
3. WHEN a Steering_File omits the `inclusion` field, THE Frontmatter_Parser SHALL default the Inclusion_Mode to `manual`.
4. WHEN a Steering_File omits the `tools` field, THE Frontmatter_Parser SHALL apply the Steering_File to all supported tools.
5. WHEN a Steering_File omits the `priority` field, THE Frontmatter_Parser SHALL default the priority to `50`.
6. WHEN a Steering_File contains invalid YAML frontmatter, THE Frontmatter_Parser SHALL return a descriptive error identifying the file path and the parsing failure.
7. THE Frontmatter_Parser SHALL format Steering_File objects back into valid frontmatter-plus-markdown files.
8. FOR ALL valid Steering_File objects, parsing then printing then parsing SHALL produce an equivalent object (round-trip property).

### Requirement 2: Steering Inclusion Modes

**User Story:** As a developer, I want steering files to be injected into AI prompts based on configurable inclusion modes, so that relevant rules appear automatically while irrelevant ones stay out of the way.

#### Acceptance Criteria

1. WHILE the Inclusion_Mode of a Steering_File is `always`, THE Steering_Layer SHALL include that Steering_File in every Context_Packet regardless of prompt content.
2. WHILE the Inclusion_Mode of a Steering_File is `auto`, THE Steering_Layer SHALL include that Steering_File in a Context_Packet only when the prompt description matches the Steering_File `description` field.
3. WHILE the Inclusion_Mode of a Steering_File is `manual`, THE Steering_Layer SHALL include that Steering_File in a Context_Packet only when the user explicitly references the Steering_File by name.
4. THE Steering_Layer SHALL inject Steering_Files ordered by ascending `priority` value (lower priority number injected first).
5. WHEN two Steering_Files share the same `priority` value, THE Steering_Layer SHALL order them alphabetically by `name`.

### Requirement 3: Steering CLI Commands

**User Story:** As a developer, I want CLI commands to create, list, preview, and initialize steering files, so that I can manage behavioral rules without manually editing files.

#### Acceptance Criteria

1. WHEN the user runs `ctx steering add <name>`, THE CLI SHALL create a new Steering_File at `.context/steering/<name>.md` with default frontmatter and an empty markdown body.
2. WHEN the user runs `ctx steering add` with a name that already exists, THE CLI SHALL return an error stating the file already exists.
3. WHEN the user runs `ctx steering list`, THE CLI SHALL display all Steering_Files with their name, Inclusion_Mode, priority, and target tools.
4. WHEN the user runs `ctx steering preview <prompt>`, THE CLI SHALL display the list of Steering_Files that would be included for the given prompt text, respecting Inclusion_Mode rules.
5. WHEN the user runs `ctx steering init`, THE CLI SHALL generate Foundation_Files (`product.md`, `tech.md`, `structure.md`, `workflow.md`) in `.context/steering/` by analyzing the current codebase.
6. WHEN Foundation_Files already exist and the user runs `ctx steering init`, THE CLI SHALL skip existing files and report which files were skipped.
7. WHEN the `.context/` directory does not exist, THE CLI SHALL return an error instructing the user to run `ctx init` first.

### Requirement 4: Steering Integration with Context Packet

**User Story:** As a developer, I want `ctx agent` to include relevant steering files in its output, so that AI tools automatically receive behavioral guidance alongside project context.

#### Acceptance Criteria

1. WHEN `ctx agent` assembles a Context_Packet, THE Agent_Command SHALL include applicable Steering_Files after the existing priority-ordered context files.
2. THE Agent_Command SHALL respect the token budget when including Steering_Files, truncating or omitting lower-priority Steering_Files when the budget is exceeded.
3. WHEN no Steering_Files exist in `.context/steering/`, THE Agent_Command SHALL produce the same Context_Packet as the current implementation without errors.


### Requirement 5: Steering Sync to Tool-Native Formats

**User Story:** As a developer, I want to sync steering files to tool-native formats for Cursor, Cline, and Kiro, so that ctx remains the single source of truth while each tool receives rules in its expected format.

#### Acceptance Criteria

1. WHEN the user runs `ctx steering sync` without the `--tool` flag and without the `--all` flag, THE CLI SHALL read the `tool` field from the active `.ctxrc` and sync to that tool's native format.
2. WHEN the user runs `ctx steering sync` without the `--tool` flag and without the `--all` flag and no `tool` field is set in `.ctxrc`, THE CLI SHALL return an error instructing the user to specify a tool with `--tool <tool>`, use `--all`, or set the `tool` field in `.ctxrc`.
3. WHEN the user runs `ctx steering sync --tool cursor`, THE CLI SHALL write each applicable Steering_File as a `.cursor/rules/<name>.mdc` file with Cursor-compatible frontmatter.
4. WHEN the user runs `ctx steering sync --tool cline`, THE CLI SHALL write each applicable Steering_File as a `.clinerules/<name>.md` file.
5. WHEN the user runs `ctx steering sync --tool kiro`, THE CLI SHALL write each applicable Steering_File as a `.kiro/steering/<name>.md` file with Kiro-compatible frontmatter.
6. WHEN the user runs `ctx steering sync --tool` with an unsupported Tool_Identifier, THE CLI SHALL return an error listing the supported Tool_Identifiers.
7. WHEN the user runs `ctx steering sync --all`, THE CLI SHALL sync Steering_Files to all supported tool formats.
8. WHEN a synced tool-native file already exists and the source Steering_File has not changed, THE CLI SHALL skip the file and not overwrite the existing content.
9. WHEN a Steering_File specifies a `tools` list that excludes a given Tool_Identifier, THE Steering_Sync SHALL skip that Steering_File for the excluded tool.

### Requirement 6: Hook Storage and Discovery

**User Story:** As a developer, I want hooks stored as executable scripts in `.context/hooks/<hook-type>/`, so that lifecycle automation is file-based, git-versionable, and language-agnostic.

#### Acceptance Criteria

1. THE Hook_Runner SHALL discover executable scripts in `.context/hooks/<hook-type>/` directories where `<hook-type>` is one of: `pre-tool-use`, `post-tool-use`, `session-start`, `session-end`, `file-save`, `context-add`.
2. WHEN a script in a Hook directory lacks the executable permission bit, THE Hook_Runner SHALL skip that script and log a warning identifying the file path.
3. THE Hook_Runner SHALL execute discovered hooks in alphabetical order by filename within each Hook_Type directory.
4. WHEN the `.context/hooks/` directory does not exist, THE Hook_Runner SHALL return an empty hook list without error.

### Requirement 7: Hook Input/Output Contract

**User Story:** As a developer, I want hooks to receive structured JSON input and return structured JSON output, so that hook scripts can make informed decisions and communicate results back to the system.

#### Acceptance Criteria

1. THE Hook_Runner SHALL pass a Hook_Input JSON object to each hook script via stdin containing: `hookType` (string), `tool` (string), `parameters` (object), `session` (object with `id` and `model` fields), `timestamp` (ISO 8601 string), and `ctxVersion` (string).
2. THE Hook_Runner SHALL read a Hook_Output JSON object from each hook script via stdout containing: `cancel` (boolean), `context` (optional string), and `message` (optional string).
3. WHEN a hook script returns `cancel: true` in the Hook_Output, THE Hook_Runner SHALL halt execution of subsequent hooks for that event and return the cancellation message.
4. WHEN a hook script returns a non-empty `context` field in the Hook_Output, THE Hook_Runner SHALL append that text to the AI conversation context.
5. IF a hook script exits with a non-zero exit code, THEN THE Hook_Runner SHALL log the error, skip that hook, and continue executing remaining hooks for the event.
6. IF a hook script produces invalid JSON on stdout, THEN THE Hook_Runner SHALL log a warning identifying the hook file and the parse error, and continue executing remaining hooks.
7. THE Hook_Runner SHALL enforce a configurable timeout (default 10 seconds) per hook execution.
8. IF a hook script exceeds the timeout, THEN THE Hook_Runner SHALL terminate the script process and log a timeout warning.

### Requirement 8: Hook CLI Commands

**User Story:** As a developer, I want CLI commands to create, list, test, enable, and disable hooks, so that I can manage lifecycle automation without manually editing files.

#### Acceptance Criteria

1. WHEN the user runs `ctx hook add <hook-type> <name>`, THE CLI SHALL create an executable script template at `.context/hooks/<hook-type>/<name>.sh` with the correct shebang, JSON input reading, and JSON output structure.
2. WHEN the user runs `ctx hook add` with an invalid Hook_Type, THE CLI SHALL return an error listing the valid Hook_Types.
3. WHEN the user runs `ctx hook list`, THE CLI SHALL display all hooks grouped by Hook_Type, showing name, enabled/disabled status, and file path.
4. WHEN the user runs `ctx hook test <hook-type> --tool <tool> --path <path>`, THE CLI SHALL construct a mock Hook_Input, execute all enabled hooks for that Hook_Type, and display the Hook_Output from each.
5. WHEN the user runs `ctx hook disable <name>`, THE CLI SHALL remove the executable permission bit from the hook script.
6. WHEN the user runs `ctx hook enable <name>`, THE CLI SHALL add the executable permission bit to the hook script.
7. WHEN the user runs `ctx hook disable` or `ctx hook enable` with a name that does not match any hook file, THE CLI SHALL return an error stating the hook was not found.

### Requirement 9: MCP Server — New Steering and Search Tools

**User Story:** As a developer, I want the MCP server to expose steering retrieval and context search as callable tools, so that AI tools can dynamically request relevant guidance and search context mid-session.

#### Acceptance Criteria

1. WHEN an MCP client calls the `ctx_steering_get` MCP_Tool with an optional `prompt` parameter, THE MCP_Server SHALL return the list of applicable Steering_Files for that prompt, respecting Inclusion_Mode rules.
2. WHEN an MCP client calls the `ctx_steering_get` MCP_Tool without a `prompt` parameter, THE MCP_Server SHALL return all Steering_Files with Inclusion_Mode `always`.
3. WHEN an MCP client calls the `ctx_search` MCP_Tool with a `query` parameter, THE MCP_Server SHALL search across all `.context/` files and return matching excerpts with file paths and line numbers.
4. THE MCP_Server SHALL register `ctx_steering_get` and `ctx_search` in the MCP tool catalog with JSON Schema parameter definitions.
5. WHEN the MCP_Server receives a request for an unregistered tool name, THE MCP_Server SHALL return a JSON-RPC error response with error code `-32601` (method not found).

### Requirement 10: MCP Server — Session Lifecycle Tools

**User Story:** As a developer, I want the MCP server to expose session start and session end tools, so that AI tools can signal lifecycle events and trigger hooks automatically.

#### Acceptance Criteria

1. WHEN an MCP client calls the `ctx_session_start` MCP_Tool, THE MCP_Server SHALL execute all enabled `session-start` hooks and return the aggregated context from hook outputs.
2. WHEN an MCP client calls the `ctx_session_end` MCP_Tool with an optional `summary` parameter, THE MCP_Server SHALL execute all enabled `session-end` hooks, passing the summary in the Hook_Input parameters.
3. WHEN an MCP client calls the `ctx_session_start` MCP_Tool and no `session-start` hooks exist, THE MCP_Server SHALL return a success response with empty context.
4. THE MCP_Server SHALL register `ctx_session_start` and `ctx_session_end` in the MCP tool catalog with JSON Schema parameter definitions.


### Requirement 11: Tool Configuration via .ctxrc Profiles

**User Story:** As a developer, I want to set my active AI tool once in `.ctxrc` and have all ctx commands respect it, so that I don't have to repeat `--tool` on every command.

#### Acceptance Criteria

1. THE RC_Package SHALL support a `tool` field in `.ctxrc` containing a single Tool_Identifier (e.g., `kiro`, `claude`, `cursor`, `cline`, `codex`).
2. WHEN the `tool` field is set in `.ctxrc`, ALL commands that accept a `--tool` flag SHALL use the `.ctxrc` value as the default.
3. WHEN the user provides a `--tool` CLI flag, THE CLI SHALL use the flag value and ignore the `.ctxrc` `tool` field.
4. WHEN the user runs `ctx config switch kiro` and a `.ctxrc.kiro` file exists, THE CLI SHALL copy `.ctxrc.kiro` to `.ctxrc`, activating the kiro tool configuration including the `tool` field.
5. THE user SHALL be able to create tool-specific `.ctxrc` profiles (e.g., `.ctxrc.kiro`, `.ctxrc.claude`) with different `tool`, `steering`, `hooks`, and `token_budget` settings.
6. WHEN the `tool` field is not set in `.ctxrc` and no `--tool` flag is provided, commands that require a Tool_Identifier SHALL return an error instructing the user to set the `tool` field or provide `--tool`.

### Requirement 12: Skills System

**User Story:** As a developer, I want to install, list, and remove reusable instruction bundles, so that I can share and reuse AI guidance across projects.

#### Acceptance Criteria

1. THE Skill SHALL be a directory in `.context/skills/<skill-name>/` containing a `SKILL.md` file with YAML frontmatter (`name`, `description`) and markdown instruction body.
2. WHEN the user runs `ctx skill install <source>`, THE CLI SHALL download or copy the Skill from the source path into `.context/skills/<skill-name>/`.
3. WHEN the user runs `ctx skill install` with a source that does not contain a valid `SKILL.md`, THE CLI SHALL return an error stating the source is not a valid skill.
4. WHEN the user runs `ctx skill list`, THE CLI SHALL display all installed skills with their name and description from the `SKILL.md` frontmatter.
5. WHEN the user runs `ctx skill remove <name>`, THE CLI SHALL delete the `.context/skills/<name>/` directory.
6. WHEN the user runs `ctx skill remove` with a name that does not match any installed skill, THE CLI SHALL return an error stating the skill was not found.
7. WHEN the user runs `ctx agent --skill <name>`, THE Agent_Command SHALL include the named Skill's `SKILL.md` content in the Context_Packet.
8. WHEN the user runs `ctx agent --skill <name>` with a name that does not match any installed skill, THE Agent_Command SHALL return an error stating the skill was not found.

### Requirement 13: Directory Initialization

**User Story:** As a developer, I want `ctx init` to create the new directories for steering, hooks, profiles, and skills, so that the project is ready for the full hooks-and-steering system from the start.

#### Acceptance Criteria

1. WHEN the user runs `ctx init`, THE Initialize_Command SHALL create `.context/steering/`, `.context/hooks/`, and `.context/skills/` directories alongside existing `.context/` subdirectories.
2. WHEN the directories already exist and the user runs `ctx init`, THE Initialize_Command SHALL skip existing directories without error.
3. THE Initialize_Command SHALL set directory permissions to `0755` for all newly created directories.

### Requirement 14: Backward Compatibility

**User Story:** As a developer, I want the hooks-and-steering system to be fully additive, so that existing ctx workflows, CLAUDE.md generation, and AGENTS.md generation continue to work without modification.

#### Acceptance Criteria

1. WHEN no `.context/steering/` directory exists, THE Agent_Command SHALL produce the same Context_Packet as the current implementation.
2. WHEN no `.context/hooks/` directory exists, THE Hook_Runner SHALL return empty results without error.
3. WHEN no `tool` field is set in `.ctxrc`, commands that do not require a Tool_Identifier SHALL continue to function with default behavior.
4. WHEN no `.context/skills/` directory exists, THE CLI SHALL report an empty skill list without error.
5. THE existing `CLAUDE.md` and `AGENTS.md` generation commands SHALL continue to function without modification.

### Requirement 15: Security Constraints

**User Story:** As a developer, I want the hooks-and-steering system to follow ctx's existing security model, so that no new attack vectors are introduced.

#### Acceptance Criteria

1. THE Hook_Runner SHALL reject hook scripts that are symlinks (consistent with ctx's symlink rejection defense layer).
2. THE Hook_Runner SHALL validate that all hook script paths resolve within the `.context/hooks/` directory boundary (consistent with ctx's boundary validation).
3. THE Steering_Sync SHALL validate that all output paths resolve within the project root directory boundary.
4. THE Hook_Runner SHALL execute hook scripts with the same user permissions as the ctx process, without privilege escalation.
5. IF a hook script attempts to write outside the project root, THEN THE Hook_Runner SHALL block the write and log a security warning.

### Requirement 16: Configuration Integration

**User Story:** As a developer, I want hooks-and-steering settings to integrate with the existing `.ctxrc` configuration system, so that all ctx configuration remains in one place.

#### Acceptance Criteria

1. THE RC_Package SHALL support a `steering` section in `.ctxrc` with fields: `dir` (path override, default `.context/steering`), `default_inclusion` (default Inclusion_Mode), and `default_tools` (default Tool_Identifier list).
2. THE RC_Package SHALL support a `hooks` section in `.ctxrc` with fields: `dir` (path override, default `.context/hooks`), `timeout` (integer seconds, default 10), and `enabled` (boolean, default true).
3. WHEN the `hooks.enabled` field in `.ctxrc` is set to `false`, THE Hook_Runner SHALL skip all hook execution.
4. THE RC_Package SHALL resolve hooks-and-steering configuration using the existing priority hierarchy: CLI flags > environment variables > `.ctxrc` > hardcoded defaults.

### Requirement 17: Drift Detection for Steering and Hooks

**User Story:** As a developer, I want `ctx drift` to detect issues with steering files and hooks, so that I am warned about stale or misconfigured behavioral guidance.

#### Acceptance Criteria

1. WHEN `ctx drift` runs and a Steering_File references a Tool_Identifier not in the supported list, THE Drift_Detector SHALL report a warning for that file.
2. WHEN `ctx drift` runs and a hook script in `.context/hooks/` lacks the executable permission bit, THE Drift_Detector SHALL report a warning for that file.
3. WHEN `ctx drift` runs and synced tool-native files are out of date compared to their source Steering_Files, THE Drift_Detector SHALL report a warning listing the stale files.
4. WHEN `ctx drift` runs and the `tool` field in `.ctxrc` contains an unsupported Tool_Identifier, THE Drift_Detector SHALL report a warning.

### Requirement 18: Use Cases

**User Story:** As a developer, I want documented use cases that demonstrate how the hooks-and-steering system solves real workflow problems, so that I understand the practical value of each phase.

#### Acceptance Criteria

1. THE Documentation SHALL describe a use case where a `session-start` hook automatically injects the full Context_Packet into an AI session, replacing manual CLAUDE.md editing.
2. THE Documentation SHALL describe a use case where a `pre-tool-use` hook blocks AI writes to a frozen legacy directory and returns a decision reference.
3. THE Documentation SHALL describe a use case where a `post-tool-use` hook runs a linter after AI file writes and injects lint results into the conversation.
4. THE Documentation SHALL describe a use case where `ctx steering sync --all` propagates a single set of API design rules to Cursor, Cline, and Kiro simultaneously.
5. THE Documentation SHALL describe a use case where `ctx config switch kiro` activates a `.ctxrc.kiro` profile with kiro-specific tool, budget, and steering settings, and `ctx config switch claude` switches back to Claude Code settings.
6. THE Documentation SHALL describe a use case where a Skill bundle for React patterns is installed from a remote source and activated via `ctx agent --skill react-patterns`.

### Requirement 19: Test Mechanisms

**User Story:** As a developer, I want each component to have clear test mechanisms, so that correctness can be verified through automated testing.

#### Acceptance Criteria

1. THE Frontmatter_Parser SHALL be testable via round-trip property: for all valid Steering_File inputs, `parse(print(parse(input))) == parse(input)`.
2. THE Hook_Runner SHALL be testable via the `ctx hook test` command, which constructs mock Hook_Input and verifies Hook_Output structure.
3. THE Steering_Layer inclusion logic SHALL be testable via the `ctx steering preview` command, which shows which files match a given prompt without side effects.
4. THE Steering_Sync SHALL be testable via idempotence property: running `ctx steering sync --tool <tool>` twice in succession SHALL produce identical output files.
5. THE MCP_Server new tools SHALL be testable via JSON-RPC requests sent over stdin, verifying response structure matches the MCP protocol specification.
6. THE Hook_Runner timeout enforcement SHALL be testable by providing a hook script that sleeps beyond the configured timeout and verifying the script is terminated.
7. THE Drift_Detector new checks SHALL be testable by constructing `.context/` directories with known issues and verifying the correct warnings are reported.
8. THE `.ctxrc` tool field SHALL be testable by verifying that commands read the `tool` value and apply it as the default Tool_Identifier.
