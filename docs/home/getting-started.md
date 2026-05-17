---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Getting Started
icon: lucide/rocket
---

![ctx](../images/ctx-banner.png)

## Prerequisites

`ctx` does not require `git`, but using version control with your `.context/`
directory is strongly recommended:

AI sessions occasionally modify or overwrite context files inadvertently.
With `git`, the AI can check history and restore lost content:
Without it, the data is gone.

Also, several `ctx` features (*journal changelog, blog generation*) also use
`git` history directly.

## Installation

Every setup starts with **the `ctx` binary**: the CLI tool itself.

If you use **Claude Code**, you also install the **`ctx` plugin**, which
adds hooks (context autoloading, persistence nudges) and 25+ `/ctx-*`
skills. For other AI tools, `ctx` integrates via generated instruction
files or manual context pasting: see
[Integrations](../operations/integrations.md) for tool-specific setup.

Pick one of the options below to install the binary. Claude Code users
should also follow the plugin steps included in each option.

### Option 1: Build from Source (*Recommended*)

Requires [Go](https://go.dev/) (*version defined in 
[`go.mod`](https://github.com/ActiveMemory/ctx/blob/main/go.mod)*) and
[Claude Code](https://docs.anthropic.com/en/docs/claude-code/overview).

```bash
git clone https://github.com/ActiveMemory/ctx.git
cd ctx
make build
sudo make install
```

**Install the Claude Code plugin** from your local clone:

1. Launch `claude`;
2. Type `/plugin` and press Enter;
3. Select **Marketplaces** → **Add Marketplace**
4. Enter the path to the **root of your clone**,
   e.g. `~/WORKSPACE/ctx`
   (*this is where `.claude-plugin/marketplace.json` lives: It points
   Claude Code to the actual plugin in `internal/assets/claude`*)
5. Back in `/plugin`, select **Install** and choose `ctx`

This points Claude Code at the plugin source on disk. Changes you make
to hooks or skills take effect immediately: No reinstall is needed.

!!! warning "Local Installs Need Manual Enablement"
    Unlike marketplace installs, local plugin installs are **not**
    auto-enabled globally. The plugin will only work in projects that
    explicitly enable it. Run `ctx init` in each project (it auto-enables
    the plugin), or add the entry to `~/.claude/settings.json` manually:

    ```json
    { "enabledPlugins": { "ctx@activememory-ctx": true } }
    ```

**Verify:**

```bash
ctx --version       # binary is in PATH
claude /plugin list # plugin is installed
```

!!! tip "Use the Source, Luke"
    Building from source gives you the latest features and bug fixes.

    Since `ctx` is predominantly a developer tool, this is the
    **recommended approach**: 

    You get the freshest code, can inspect what
    you are installing, and the plugin stays in sync with the binary.

### Option 2: Binary Download + Marketplace

Pre-built binaries are available from the
[releases page](https://github.com/ActiveMemory/ctx/releases).

=== "Linux (x86_64)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.8.1/ctx-0.8.1-linux-amd64
    chmod +x ctx-0.8.1-linux-amd64
    sudo mv ctx-0.8.1-linux-amd64 /usr/local/bin/ctx
    ```

=== "Linux (ARM64)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.8.1/ctx-0.8.1-linux-arm64
    chmod +x ctx-0.8.1-linux-arm64
    sudo mv ctx-0.8.1-linux-arm64 /usr/local/bin/ctx
    ```

=== "macOS (Apple Silicon)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.8.1/ctx-0.8.1-darwin-arm64
    chmod +x ctx-0.8.1-darwin-arm64
    sudo mv ctx-0.8.1-darwin-arm64 /usr/local/bin/ctx
    ```

=== "macOS (Intel)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.8.1/ctx-0.8.1-darwin-amd64
    chmod +x ctx-0.8.1-darwin-amd64
    sudo mv ctx-0.8.1-darwin-amd64 /usr/local/bin/ctx
    ```

=== "Windows"

    Download `ctx-0.8.1-windows-amd64.exe` from the releases page and add it to your `PATH`.

**Claude Code users**: install the plugin from the marketplace:

1. Launch `claude`;
2. Type `/plugin` and press Enter;
3. Select **Marketplaces** → **Add Marketplace**;
4. Enter `ActiveMemory/ctx`;
5. Back in `/plugin`, select **Install** and choose `ctx`.

**Other tool users**: see [Integrations](../operations/integrations.md) for
tool-specific setup (Cursor, Copilot, Aider, Windsurf, etc.).

!!! note "Verify the Plugin Is Enabled"
    After installing, confirm the plugin is enabled globally. Check
    `~/.claude/settings.json` for an `enabledPlugins` entry. If missing,
    run `ctx init` in your project (it auto-enables the plugin), or add
    it manually:

    ```json
    { "enabledPlugins": { "ctx@activememory-ctx": true } }
    ```

**Verify:**

```bash
ctx --version       # binary is in PATH
claude /plugin list # plugin is installed (Claude Code only)
```

#### Verifying Checksums

Each binary has a corresponding `.sha256` checksum file. To verify your download:

```bash
# Download the checksum file
curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.8.1/ctx-0.8.1-linux-amd64.sha256

# Verify the binary
sha256sum -c ctx-0.8.1-linux-amd64.sha256
```

On macOS, use `shasum -a 256 -c` instead of `sha256sum -c`.

----

??? note "Plugin Details"
    After installation (*either option*) you get:

    * **Context autoloading**: `ctx agent` runs on every tool use (*with cooldown*)
    * **Persistence nudges**: reminders to capture learnings and decisions
    * **Post-commit hooks**: nudge context capture after `git commit`
    * **Context size monitoring**: alerts as sessions grow large
    * **Project skills**: `/ctx-status`, `/ctx-task-add`, `/ctx-history`, and more

    See [Integrations](../operations/integrations.md#claude-code-full-integration) for the
    full hook and skill reference.

## Quick Start

### 1. Initialize Context

```bash
cd your-project
ctx init
```

This creates a `.context/` directory with template files and an
encryption key at `~/.ctx/` for the
[encrypted scratchpad](../reference/scratchpad.md).
For Claude Code, install the [`ctx` plugin](../operations/integrations.md#claude-code-full-integration)
for automatic hooks and skills.

`ctx init` also scaffolds four **foundation steering files** in
`.context/steering/`: `product.md`, `tech.md`, `structure.md`,
`workflow.md`. **They are placeholders until you customize
them** (see the next step); skipping that step has consequences,
so it is broken out as its own numbered beat rather than
buried here.

### 2. Customize Your Steering Files

Steering files are **behavioral rules prepended to every AI
prompt**: the layer that tells your AI *how to act* on this
specific project. They are distinct from decisions (*what* was
chosen) and conventions (*how* the codebase is written); see
[`ctx` for Steering Files](../recipes/steering.md) for the full
model.

`ctx init` scaffolded four foundation files; open each and
fill it in:

| File            | What to fill in                                        |
|-----------------|--------------------------------------------------------|
| `product.md`    | What the project is, who uses it, what's out of scope  |
| `tech.md`       | Languages, frameworks, runtime, hard constraints       |
| `structure.md`  | Directory layout, where new files go, naming rules     |
| `workflow.md`   | Branch strategy, commit conventions, pre-commit checks |

Each scaffolded file ships with a **tombstone marker** line
(`<!-- remove this after you edit the steering file !-->`).
**As long as the marker is present, the file is silently
skipped** on every load path: the agent context packet, MCP
`ctx_steering_get`, and native-tool sync (Cursor / Cline /
Kiro). The skip is deliberate: injecting unfilled placeholders
into AI prompts is worse than no steering at all, because the
AI tries to follow "Describe the product..." as if it were a
rule.

**Replace each file's body with real content, then delete the
tombstone line.** When the line is gone, the file becomes
active on the next AI tool call.

Don't want steering at all? Pass `--no-steering-init` to
`ctx init` to skip the scaffold entirely. Existing edits are
never clobbered by re-running `ctx init`.

Inclusion modes (`always` / `auto` / `manual`), priority, and
tool scoping are covered in
[Writing Steering Files](../recipes/steering.md) and
[`ctx steering`](../cli/steering.md).

### 3. Activate the Project

Tell `ctx` which `.context/` directory the rest of these commands
should use:

```bash
eval "$(ctx activate)"
```

You only need to run this once per terminal. If you skip it, the
next steps fail with `Error: no context directory specified`. Direnv
users can wire it into `.envrc` and forget about it. For more
options (multiple `.context/` directories, scripts, CI), see
[Activating a Context Directory](../recipes/activating-context.md).

### 4. Check Status

```bash
ctx status
```

Shows context summary: files present, token estimate, and recent activity.

### 5. Start Using with AI

With Claude Code (*and the `ctx` plugin installed*), context loads automatically
via hooks.

With **VS Code Copilot Chat**, install the
[`ctx` extension](../operations/integrations.md#vs-code-chat-extension-ctx) and use
`@ctx /status`, `@ctx /agent`, and other slash commands directly in chat.
Run `ctx setup copilot --write` to generate `.github/copilot-instructions.md`
for automatic context loading.

For other tools, paste the output of:

```bash
ctx agent --budget 8000
```

### 5B. Set Up for Your AI Tool

If you use an MCP-compatible tool, generate the integration config
with `ctx setup`:

=== "Kiro"

    ```bash
    ctx setup kiro --write
    # Creates .kiro/settings/mcp.json and syncs steering files
    ```

=== "Cursor"

    ```bash
    ctx setup cursor --write
    # Creates .cursor/mcp.json and syncs steering files
    ```

=== "Cline"

    ```bash
    ctx setup cline --write
    # Creates .vscode/mcp.json and syncs steering files
    ```

This registers the `ctx` MCP server and syncs any
[steering files](../cli/steering.md) into the tool's
native format. Re-run after adding or changing steering files.

### 6. Verify It Works

Ask your AI: **"Do you remember?"**

It should cite specific context: current tasks, recent decisions,
or previous session topics.

### 7. Set Up Companion Tools (Highly Recommended)

`ctx` works on its own, but two companion MCP servers unlock significantly
better agent behavior. The investment is small and the benefits compound
over sessions:

* **[Gemini Search](https://github.com/nicobailon/gemini-code-search-mcp)**
: grounded web search with citations. Skills like `/ctx-code-review`
  and `/ctx-explain` use it for up-to-date documentation lookups instead
  of relying on training data.

* **[GitNexus](https://github.com/nicobailon/gitnexus-mcp)**: code
  knowledge graph with symbol resolution, blast radius analysis, and
  domain clustering. Skills like `/ctx-refactor` and `/ctx-code-review`
  use it for impact analysis and dependency awareness.

```bash
# Index your project for GitNexus (run once, then after major changes)
npx gitnexus analyze
```

Both are optional MCP servers: if they are not connected, skills degrade
gracefully to built-in capabilities. See
[Companion Tools](../recipes/multi-tool-setup.md#companion-tools-highly-recommended)
for setup details and verification.

----

**Next Up**:

* [Your First Session →](first-session.md): a step-by-step walkthrough
  from `ctx init` to verified recall
* [Common Workflows →](common-workflows.md): day-to-day commands for
  tracking context, checking health, and browsing history
