---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Local Inference with vLLM
icon: lucide/server
---

![ctx](../images/ctx-banner.png)

Use a local vLLM OpenAI-compatible server for optional `ctx ai` commands.
This does not change deterministic ctx behavior: `ctx status`, `ctx agent`,
and hooks keep working when no backend is configured or reachable.

## 1. Start vLLM

Run vLLM with its OpenAI-compatible API enabled. The exact model is your choice;
ctx only needs `/v1/models` and `/v1/chat/completions`.

## 2. Configure ctx

From the project root:

```bash
ctx setup --backend vllm --endpoint http://localhost:8000 --write
```

This writes a `.ctxrc` `backends:` block similar to:

```yaml
backends:
  default: vllm
  vllm:
    endpoint: http://localhost:8000
```

Add a model or timeout if needed:

```yaml
backends:
  default: vllm
  vllm:
    endpoint: http://localhost:8000
    default_model: Qwen/Qwen2.5-Coder-7B-Instruct
    timeout: 30s
```

## 3. Verify reachability

```bash
ctx ai ping
```

Expected output names the backend, endpoint, and first model listed by vLLM.

## 4. Write a reviewable proposal

```bash
ctx ai propose notes.md --emit decisions,learnings,tasks,open-questions
```

ctx writes one JSON proposed-patch artifact under `.context/proposals/ai/`.
Review it before changing canonical `.context/*.md` files.

## Failure behavior

AI commands fail closed when no backend is configured, multiple backends are
ambiguous, the backend is unreachable, or the upstream response is invalid.
They do not fall back to a different model and do not affect non-AI commands.
