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
ctx setup --backend <name> [flags]
```

**Flags**:

| Flag            | Short | Description                                                                 |
|-----------------|-------|-----------------------------------------------------------------------------|
| `--write`       | `-w`  | Write the generated config to disk (e.g. `.github/copilot-instructions.md`) |
| `--backend`     |       | Configure a `ctx ai` backend in `.ctxrc`                                    |
| `--endpoint`    |       | Backend endpoint URL                                                        |
| `--api-key-env` |       | Environment variable that contains the backend API key                      |
| `--model`       |       | Default model for backend requests                                          |
| `--timeout`     |       | Backend request timeout duration                                            |

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

**Supported backends**:

| Backend               | Description                                      |
|-----------------------|--------------------------------------------------|
| `vllm`                | Local vLLM OpenAI-compatible server              |
| `openai-compatible`   | Generic OpenAI-compatible HTTP backend           |
| `openai`              | OpenAI-compatible backend with OpenAI defaults   |
| `anthropic`           | Anthropic backend using the compatibility floor  |
| `ollama`              | Local Ollama OpenAI-compatible endpoint defaults |
| `lmstudio`            | Local LM Studio endpoint defaults                |

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

# Configure a local vLLM backend for ctx ai commands
ctx setup --backend vllm --endpoint http://localhost:8000 --write

# Preview an OpenAI-compatible backend config without writing
ctx setup --backend openai-compatible \
  --endpoint https://llm.example.com \
  --api-key-env OPENAI_API_KEY \
  --model gpt-4.1-mini
```
