---
title: ctx ai
icon: lucide/brain-circuit
---

![ctx](../images/ctx-banner.png)

## `ctx ai`

Dispatch through a configured AI backend (optional, local-first AI
capabilities).

```bash
ctx ai <verb> [flags]
```

`ctx ai` is the namespace for optional AI capabilities that talk to
a configured backend (vLLM, OpenAI, Anthropic, Ollama, LM Studio,
or any OpenAI-compatible HTTP endpoint). The backend is read from
`.ctxrc`; configure one with [`ctx setup --backend <name>`](setup.md#ctx-setup).

All `ctx ai *` commands **fail closed** when no backend is
configured: there is no silent degradation to a non-AI path. The
deterministic ctx core (`ctx agent`, `ctx status`, ceremony hooks)
is untouched whether or not a backend is set up.

!!! note "Strictly additive"
    `ctx ai` is a Block A foundation. It does not change any
    existing command's behavior, and `ctx agent` / `ctx status` /
    deterministic ceremony hooks may not import the backend layer
    (enforced by a boundary-guard unit test). Removing the
    `backends:` section from `.ctxrc` returns the project to a
    fully deterministic state.

### `ctx ai ping`

Check reachability of the configured AI backend.

```bash
ctx ai ping [--backend <name>]
```

Issues a `GET /v1/models` against the configured backend and
reports the backend name and the first model the server listed.
On the vLLM backend, refused connections during the cold-start
window are retried so the command returns once weights finish
loading (default window: 90 seconds).

**Flags**:

| Flag        | Description                                                                                |
|-------------|--------------------------------------------------------------------------------------------|
| `--backend` | Pick which configured backend to ping (defaults to the `default_backend:` key in `.ctxrc`) |

**Exit codes**:

| Code | Meaning                                                                  |
|------|--------------------------------------------------------------------------|
| 0    | Backend reachable and listed at least one model                          |
| 1    | No backends configured; unknown name; unreachable; lists no models       |

**Examples**:

```bash
ctx ai ping
# backend "vllm" reachable; first model: "meta-llama/Llama-3.1-8B-Instruct"

ctx ai ping --backend openai
# backend "openai" reachable; first model: "gpt-4o-mini"
```

### `ctx ai extract`

Extract decisions, learnings, tasks, and open questions from
stdin into a proposal file.

```bash
ctx ai extract [--backend <name>] < input.md
```

Reads free text from stdin, dispatches a JSON-mode chat completion
through the configured backend, and writes the response as a
proposal under `.context/proposals/<TS>-extract.md`. The canonical
`.context/*.md` files are **never written directly**: every
AI-produced edit lands as a proposed patch in the proposal queue
for human (or agent) ratification.

The proposal file is markdown wrapping the model's raw JSON
response inside a fenced code block. Use the proposal queue as a
draft store: after review, apply selected candidates with the
appropriate `ctx <noun> add` command (`ctx decision add`,
`ctx learning add`, `ctx task add`, etc.).

**Flags**:

| Flag        | Description                                                                                |
|-------------|--------------------------------------------------------------------------------------------|
| `--backend` | Pick which configured backend to dispatch through (defaults to `default_backend:`)         |

**Exit codes**:

| Code | Meaning                                                            |
|------|--------------------------------------------------------------------|
| 0    | Proposal written                                                   |
| 1    | Empty input; no backends configured; unreachable; upstream failure |

**Examples**:

```bash
cat session-notes.md | ctx ai extract
# proposal written to .context/proposals/20260523T220000Z-extract.md

ctx journal source --limit 5 | ctx ai extract --backend vllm
```

!!! note "Proposals are gitignored by default"
    Fresh `ctx init` projects ignore `.context/proposals/*` with a
    `.gitkeep` carve-out, mirroring the `handovers/` shape. The
    queue is a per-session draft store, not a tracked artifact.

## Configuration

`ctx ai` reads its backend roster from `.ctxrc`. A minimal
configuration:

```yaml
# .ctxrc
default_backend: vllm
backends:
  - name: vllm
    endpoint: http://localhost:8000
  - name: openai
    api_key_env: OPENAI_API_KEY
```

**Per-backend fields**:

| Field           | Type     | Description                                                                                   |
|-----------------|----------|-----------------------------------------------------------------------------------------------|
| `name`          | string   | Backend type label. Must match a registered backend (see Supported backends below).           |
| `endpoint`      | string   | Base URL. Must be `http://` or `https://`. Per-vendor defaults apply when omitted.            |
| `api_key_env`   | string   | Env-var name that holds the bearer token. Empty means no `Authorization` header is sent.      |
| `timeout`       | duration | Per-request deadline (e.g. `30s`, `2m`). Defaults to 30 seconds.                              |
| `default_model` | string   | Model ID used when a request does not specify one. Empty requires explicit model selection.   |

**Top-level keys**:

| Key               | Description                                                                                 |
|-------------------|---------------------------------------------------------------------------------------------|
| `default_backend` | Which backend `ctx ai *` commands use when `--backend` is not passed. Required if 2+ entries. |

**Supported backends**:

| Name                | Default endpoint              | Default `api_key_env` | Notes                                                                                |
|---------------------|-------------------------------|-----------------------|--------------------------------------------------------------------------------------|
| `vllm`              | `http://localhost:8000`       | *(none)*              | Canonical local backend. `Ping` retries refused connections during the cold-start window. |
| `openai-compatible` | *(required: no default)*      | *(none)*              | Contract floor; configure when the vendor isn't one of the named wrappers.           |
| `openai`            | `https://api.openai.com`      | `OPENAI_API_KEY`      | OpenAI-compatible `/v1/chat/completions`.                                            |
| `anthropic`         | `https://api.anthropic.com`   | `ANTHROPIC_API_KEY`   | Uses Anthropic's OpenAI-compatible endpoint (not the native Messages API).           |
| `ollama`            | `http://localhost:11434`      | *(none)*              | Local Ollama daemon; unauthenticated by default.                                     |
| `lmstudio`          | `http://localhost:1234`       | *(none)*              | Local LM Studio server; unauthenticated by default.                                  |

User values in `.ctxrc` always win — per-vendor defaults only
fill the gaps.

## See also

* [Local Inference with vLLM](../recipes/local-inference-with-vllm.md): end-to-end recipe
* [`ctx setup --backend`](setup.md#ctx-setup): write a backend entry to `.ctxrc`
