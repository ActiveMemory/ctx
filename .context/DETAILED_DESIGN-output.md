# Detailed Design: Output Layer

Modules: write/*, err/*, assets/read/*

## internal/write/*

**Purpose**: 46 output formatting packages — one per command/feature.
Centralizes all user-facing terminal output.

**Pattern**: Every write function takes `*cobra.Command` for output
routing. Nil commands are safe no-ops.

```go
func Added(cmd *cobra.Command, filename string) {
    if cmd == nil { return }
    cmd.Println(desc.Text(text.DescKeyAddSuccess), filename)
}
```

**Key packages**:

| Package | Purpose | Key Functions |
|---------|---------|---------------|
| write/add | Add command output | Added(), SpecNudge() |
| write/agent | Agent packet rendering | Packet() |
| write/status | Context health display | Header(), FileItem(), Activity() |
| write/journal | Journal operations | InfoOrphanRemoved(), InfoSiteGenerated() |
| write/pad | Scratchpad operations | InfoPathConversionExists() |
| write/notify | Webhook setup/test | SetupPrompt(), SetupDone(), TestResult() |
| write/bootstrap | Init output | CommunityFooter(), Dir(), Text(), JSON() |
| write/doctor | Health reporting | Report() (text), JSON() |
| write/compact | Task compaction | InfoMovingTask(), SectionsRemoved() |
| write/config | Config output | ProfileStatus(), Schema(), SwitchConfirm() |
| write/drift | Drift reports | formatted violation output |
| write/trace | Trace display | formatted trace history |
| write/change | Change output | List() |
| write/complete | Task completion | Completed() |
| write/archive | Task archiving | Skipping(), NoCompleted(), Success() |
| write/line | Line output | generic line formatting |
| write/message | Message templates | TemplateVars(), OverrideCreated() |
| write/err | Error display | formatted error output |

**Output modes**:
- Terminal text: most packages via cmd.Println()
- JSON: bootstrap, doctor, resource have JSON() variants
- Pre-rendered markdown: agent/Packet()

**Relationship to CLI**: One-directional. CLI cmd/root/run.go
imports write packages at the handler level. Core business logic
never imports write packages.

**Danger zones**:
1. All text from config/embed/text — missing DescKey causes empty
   output (not a crash), which is hard to detect visually.
2. write/ packages reference desc.Text() at call time, not at
   init — YAML lookup failures appear at runtime.
3. No unit tests for most write packages — output correctness
   verified only through integration tests.

**Extension points**:
- Add new write package per new command
- Add JSON() variant for structured output
- Text content is in YAML — update without code changes

**Dependencies**: spf13/cobra, assets/read/desc, config/embed/text,
format

---

## internal/err/*

**Purpose**: 35 error constructor packages — one per functional
domain. Centralizes error messages with YAML text lookup.

**Pattern**: Every error constructor returns `error`. Messages
looked up from embedded YAML. No sentinel errors (var ErrX);
all errors come from functions.

```go
func FileNotFound(path string) error {
    return fmt.Errorf(
        desc.Text(text.DescKeyErrFsFileNotFound), path,
    )
}
```

**One custom error type**: `err/context/NotFoundError` allows
`errors.As()` matching.

**Key packages**:

| Package | Constructors | Domain |
|---------|-------------|--------|
| err/fs | 27 | All filesystem operations (read, write, stat, mkdir) |
| err/add | 6 | Entry add operations |
| err/cli | 4 | CLI validation (flags, args, selections) |
| err/context | 4 | Context directory validation |
| err/journal | 11 | Journal operations |
| err/task | 6 | Task file operations |
| err/pad | 14 | Scratchpad validation |
| err/mcp | 2 | MCP tool validation |
| err/trace | 6 | Git/hook operations |
| err/crypto | 3 | Encryption operations |
| err/notify | 3 | Webhook operations |

**Error wrapping convention**:
```go
// Pattern: template with path + cause
fmt.Errorf(desc.Text(text.DescKeyErrFsFileWrite), path, cause)
// Output: "failed to write /path/to/file: underlying error"
```

**Danger zones**:
1. Error messages are in YAML — missing keys produce empty error
   messages. Caught by TestDescKeyYAMLLinkage audit.
2. fmt.Errorf with %w is used inconsistently — some errors wrap
   causes, others don't. errors.Is/As behavior varies.
3. No error codes — programmatic error handling relies on string
   matching or the single custom type.

**Extension points**:
- Add new err/ package per new domain
- Error text is in YAML — update messages without code changes

**Dependencies**: assets/read/desc, config/embed/text, fmt, errors

---

## internal/assets/read/*

**Purpose**: 14 typed accessor packages that provide domain-specific
interfaces to the embedded filesystem. Avoids "god-object" assets
package.

**Design rule**: Dependencies flow one way: read/* imports assets.FS;
assets/ never imports read/*. Domain packages expose clean
interfaces.

**Key packages**:

| Package | Purpose | Key Functions |
|---------|---------|---------------|
| read/desc | Text description YAML lookup | Text(key), Flag(name), Command(key), Example() |
| read/skill | Skill definitions | List(), Content(name) |
| read/entry | Entry templates | List(), ForName(name) |
| read/claude | Claude Code integration | Md() |
| read/hook | Hook assets | TraceScript(), Message(), MessageRegistry() |
| read/agent | Copilot integration | CopilotInstructions(), AgentsMd() |
| read/catalog | Template listing | List() |
| read/journal | Journal styling | ExtraCSS() |
| read/project | Directory README templates | Readme(dir) |
| read/schema | Config schema | Schema() |
| read/template | Generic template reader | Template(name) |
| read/lookup | Pre-loaded YAML maps | TextDesc(), StopWords(), Init() |
| read/makefile | Makefile templates | (presumed) |
| read/philosophy | Philosophy docs | (presumed) |

**Data flow**: `lookup.Init()` (called at startup) pre-loads all
YAML maps. Subsequent `desc.Text()` calls are O(1) map lookups.

**Danger zones**:
1. lookup.Init() must be called before any desc accessor —
   calling desc.Text() before Init() returns empty strings
   silently (no panic, no warning).
2. read/desc is imported by 80+ packages — changes to its API
   have massive blast radius.
3. YAML file structure changes (key renames, file splits) require
   synchronized code + YAML updates.

**Extension points**:
- Add new read/* package for new embedded asset category
- YAML-backed text content is extensible without code changes

**Dependencies**: assets (embed.FS), config/embed/* (key constants)
