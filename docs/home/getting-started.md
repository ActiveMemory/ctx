---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Getting Started
icon: lucide/rocket
---

![ctx](images/ctx-banner.png)

## Prerequisites

`ctx` does not require `git`, but using version control with your `.context/`
directory is strongly recommended:

AI sessions occasionally modify or overwrite context files inadvertently.
With `git`, the AI can check history and restore lost content:
Without it, the data is gone.

Also, several `ctx` features (*journal changelog, blog generation*) also use
`git` history directly.

## Installation

A full ctx installation has two parts:

1. **The `ctx` binary**: the CLI tool itself
2. **The Claude Code plugin**: hooks and skills that make Claude Code
   context-aware

You need **both**. The binary alone gives you the CLI but Claude Code
won't autoload context, nudge you to persist decisions, or provide
the `/ctx-*` skills.

Pick one of the options below. Each is a complete path from zero to
a working setup.

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

1. Launch `claude`
2. Type `/plugin` and press Enter
3. Select **Marketplaces** → **Add Marketplace**
4. Enter the path to the **root of your clone**,
   e.g. `~/WORKSPACE/ctx`
   (this is where `.claude-plugin/marketplace.json` lives: It points
   Claude Code to the actual plugin in `internal/assets/claude`)
5. Back in `/plugin`, select **Install** and choose `ctx`

This points Claude Code at the plugin source on disk. Changes you make
to hooks or skills take effect immediately: No reinstall is needed.

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
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.6.0/ctx-0.6.0-linux-amd64
    chmod +x ctx-0.6.0-linux-amd64
    sudo mv ctx-0.6.0-linux-amd64 /usr/local/bin/ctx
    ```

=== "Linux (ARM64)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.6.0/ctx-0.6.0-linux-arm64
    chmod +x ctx-0.6.0-linux-arm64
    sudo mv ctx-0.6.0-linux-arm64 /usr/local/bin/ctx
    ```

=== "macOS (Apple Silicon)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.6.0/ctx-0.6.0-darwin-arm64
    chmod +x ctx-0.6.0-darwin-arm64
    sudo mv ctx-0.6.0-darwin-arm64 /usr/local/bin/ctx
    ```

=== "macOS (Intel)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.6.0/ctx-0.6.0-darwin-amd64
    chmod +x ctx-0.6.0-darwin-amd64
    sudo mv ctx-0.6.0-darwin-amd64 /usr/local/bin/ctx
    ```

=== "Windows"

    Download `ctx-0.6.0-windows-amd64.exe` from the releases page and add it to your `PATH`.

**Install the Claude Code plugin** from the marketplace:

1. Launch `claude`
2. Type `/plugin` and press Enter
3. Select **Marketplaces** → **Add Marketplace**
4. Enter `ActiveMemory/ctx`
5. Back in `/plugin`, select **Install** and choose `ctx`

**Verify:**

```bash
ctx --version       # binary is in PATH
claude /plugin list # plugin is installed
```

#### Verifying Checksums

Each binary has a corresponding `.sha256` checksum file. To verify your download:

```bash
# Download the checksum file
curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.6.0/ctx-0.6.0-linux-amd64.sha256

# Verify the binary
sha256sum -c ctx-0.6.0-linux-amd64.sha256
```

On macOS, use `shasum -a 256 -c` instead of `sha256sum -c`.

----

??? note "Plugin Details"
    After installation (either option) you get:

    * **Context autoloading**: `ctx agent` runs on every tool use (with cooldown)
    * **Persistence nudges**: reminders to capture learnings and decisions
    * **Post-commit hooks**: nudge context capture after `git commit`
    * **Context size monitoring**: alerts as sessions grow large
    * **25+ skills**: `/ctx-status`, `/ctx-add-task`, `/ctx-recall`, and more

    See [Integrations](integrations.md#claude-code-full-integration) for the
    full hook and skill reference.

## Quick Start

### 1. Initialize Context

```bash
cd your-project
ctx init
```

This creates a `.context/` directory with template files and a
`.scratchpad.key` for the [encrypted scratchpad](scratchpad.md).
For Claude Code, install the [ctx plugin](integrations.md#claude-code-full-integration)
for automatic hooks and skills.

### 2. Check Status

```bash
ctx status
```

Shows context summary: files present, token estimate, and recent activity.

### 3. Start Using with AI

With Claude Code (and the ctx plugin installed), context loads automatically
via hooks.

With **VS Code Copilot Chat**, install the
[ctx extension](integrations.md#vs-code-chat-extension-ctx) and use
`@ctx /status`, `@ctx /agent`, and other slash commands directly in chat.
Run `ctx hook copilot --write` to generate `.github/copilot-instructions.md`
for automatic context loading.

For other tools, paste the output of:

```bash
ctx agent --budget 8000
```

### 4. Verify It Works

Ask your AI: **"Do you remember?"**

It should cite specific context: current tasks, recent decisions,
or previous session topics.

----

**Next Up**:

* [Your First Session →](first-session.md) — a step-by-step walkthrough from `ctx init` to verified recall
* [Common Workflows →](common-workflows.md) — day-to-day commands for tracking context, checking health, and browsing history
