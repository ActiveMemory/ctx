---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Setup
icon: lucide/toy-brick
---

![ctx](../images/ctx-banner.png)

## `ctx setup`

Generate AI tool integration configuration.

```bash
ctx setup <tool> [flags]
```

**Flags**:

| Flag      | Short | Description                                                                 |
|-----------|-------|-----------------------------------------------------------------------------|
| `--write` | `-w`  | Write the generated config to disk (e.g. `.github/copilot-instructions.md`) |

**Supported tools**:

| Tool          | Description                                  |
|---------------|----------------------------------------------|
| `claude-code` | Redirects to plugin install instructions     |
| `cursor`      | Cursor IDE                                   |
| `kiro`        | Kiro IDE                                     |
| `cline`       | Cline (VS Code extension)                    |
| `aider`       | Aider CLI                                    |
| `copilot`     | GitHub Copilot                               |
| `opencode`    | OpenCode (terminal-first AI coding agent)    |
| `windsurf`    | Windsurf IDE                                 |

!!! note "Claude Code Uses the Plugin System"
    Claude Code integration is now provided via the `ctx` plugin.
    Running `ctx setup claude-code` prints plugin install instructions.

**Examples**:

```bash
# Print hook instructions to stdout
ctx setup cursor
ctx setup aider

# Generate and write .github/copilot-instructions.md
ctx setup copilot --write

# Generate MCP config and sync steering files
ctx setup kiro --write
ctx setup cursor --write
ctx setup cline --write

# Generate OpenCode plugin, skills, AGENTS.md, and global MCP config
ctx setup opencode --write
```

## `ctx setup --backend`

Write an AI backend entry into `.ctxrc` instead of templating an
editor tool. The `--backend` mode is mutually exclusive with the
positional tool argument: backend writing and editor-tool
templating are separate operations.

```bash
ctx setup --backend <name> [--endpoint <url>] [--api-key-env <var>]
```

The writer round-trips `.ctxrc` through a yaml.Node tree so
unrelated top-level keys and comments are preserved. Re-running
the same `--backend` name updates the existing entry in place;
new names append.

**Flags**:

| Flag             | Description                                                                                  |
|------------------|----------------------------------------------------------------------------------------------|
| `--backend`      | Backend type label (`vllm`, `openai`, `anthropic`, `ollama`, `lmstudio`, `openai-compatible`) |
| `--endpoint`     | Override the backend's default endpoint URL                                                  |
| `--api-key-env`  | Override the env-var name the backend reads for its bearer token                             |

Per-vendor defaults (endpoint + env-var name) auto-apply for the
six known backends; user values in `.ctxrc` always win.

**Examples**:

```bash
# Wire a local vLLM at the default endpoint
ctx setup --backend vllm

# Wire vLLM on a non-default port
ctx setup --backend vllm --endpoint http://gpu-host:8080

# Wire OpenAI with the canonical env-var name
ctx setup --backend openai
# (api_key_env defaults to OPENAI_API_KEY)

# Wire a generic OpenAI-compatible endpoint with a custom key var
ctx setup --backend openai-compatible \
  --endpoint https://router.example.com \
  --api-key-env ROUTER_API_KEY
```

See [`ctx ai`](ai.md#ctx-ai) for the verbs that dispatch through a
configured backend, and the
[Local Inference with vLLM](../recipes/local-inference-with-vllm.md)
recipe for the end-to-end flow.
