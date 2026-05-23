---
title: "Local Inference with vLLM"
icon: lucide/cpu
---

![ctx](../images/ctx-banner.png)

## The Problem

You want `ctx` skills to use a model you control: a local vLLM
instance on your workstation or a shared LAN box. Sending session
notes to a SaaS provider is off the table for compliance, cost,
or data-residency reasons.

`ctx ai` lets you wire any OpenAI-compatible HTTP endpoint into
`ctx` as an optional, fail-closed AI layer. The deterministic ctx
core (`ctx agent`, `ctx status`, ceremony hooks) keeps working
whether the backend is up or not.

This recipe walks the end-to-end vLLM flow: start the server,
write it into `.ctxrc` with `ctx setup --backend vllm`, verify
reachability with `ctx ai ping`, and drive a structured extraction
with `ctx ai extract`. Everything that touches `.context/` lands
in the proposal queue for review — no canonical file is
overwritten by the model.

## TL;DR

```bash
# Start vLLM (any OpenAI-compatible endpoint works; this is the canonical local one)
vllm serve meta-llama/Llama-3.1-8B-Instruct --port 8000

# Wire it into .ctxrc
ctx setup --backend vllm --endpoint http://localhost:8000

# Verify the backend is reachable
ctx ai ping
# backend "vllm" reachable; first model: "meta-llama/Llama-3.1-8B-Instruct"

# Run an extraction; output lands in the proposal queue, not in .context/*.md
ctx journal source --limit 5 | ctx ai extract
# proposal written to .context/proposals/20260523T220000Z-extract.md
```

The proposal file is markdown wrapping the model's JSON response.
Review it, then apply selected entries with the usual
`ctx <noun> add` commands.

## Commands and Skills Used

| Command/Skill                  | Role in this workflow                                                                                            |
|--------------------------------|------------------------------------------------------------------------------------------------------------------|
| `ctx setup --backend <name>`   | Write a backend entry into `.ctxrc` (idempotent; updates an existing entry in place, appends a new one otherwise) |
| `ctx ai ping`                  | `GET /v1/models` against the configured backend; reports the first model listed                                  |
| `ctx ai extract`               | Read stdin, dispatch a JSON-mode completion, write the response as a proposal file                               |
| `ctx journal source`           | Stream recent session text to pipe into `ctx ai extract`                                                         |
| `ctx decision/learning/task add` | Apply candidates from a reviewed proposal back into the canonical files                                          |

## The Workflow

### Step 1: Start vLLM

vLLM is the canonical local backend. Any model that fits in your
GPU works; this recipe uses Llama 3.1 8B Instruct as a small,
widely-available example.

```bash
vllm serve meta-llama/Llama-3.1-8B-Instruct --port 8000
```

The first launch downloads weights, which can take several
minutes. `ctx ai ping` knows about this cold-start window — see
Step 3.

!!! tip "Other OpenAI-compatible servers work too"
    `ctx` also ships per-vendor wrappers for `openai`,
    `anthropic`, `ollama`, and `lmstudio`, plus a generic
    `openai-compatible` backend for anything else that speaks the
    OpenAI HTTP surface. The flow below is identical; only the
    `--backend` value and the default endpoint change. See
    [`ctx ai` reference](../cli/ai.md#configuration) for the full
    backend roster.

### Step 2: Wire the backend into `.ctxrc`

```bash
ctx setup --backend vllm --endpoint http://localhost:8000
```

This appends (or updates in place) a `backends:` entry in
`.ctxrc`, preserving comments and other top-level keys. The
default endpoint for `vllm` is `http://localhost:8000`, so the
`--endpoint` flag is only necessary when the server listens
elsewhere.

The resulting `.ctxrc` snippet:

```yaml
default_backend: vllm
backends:
  - name: vllm
    endpoint: http://localhost:8000
```

**Flags**:

| Flag             | Description                                                                                  |
|------------------|----------------------------------------------------------------------------------------------|
| `--backend`      | Backend type label (`vllm`, `openai`, `anthropic`, `ollama`, `lmstudio`, `openai-compatible`) |
| `--endpoint`     | Override the backend's default endpoint URL                                                  |
| `--api-key-env`  | Override the env-var name the backend reads for its bearer token                             |

The `--backend` mode is mutually exclusive with the positional
tool argument (`ctx setup cursor`, `ctx setup aider`, etc.):
backend writing and editor-tool templating are separate
operations.

### Step 3: Verify reachability

```bash
ctx ai ping
# backend "vllm" reachable; first model: "meta-llama/Llama-3.1-8B-Instruct"
```

`ctx ai ping` issues a `GET /v1/models` against the configured
endpoint, asserts HTTP 200, parses the `data[].id` array, and
prints the first model.

On the `vllm` backend, refused connections during the **cold-start
window** (default 90 seconds) are retried at 500ms intervals.
This is necessary because the vLLM listener isn't bound while
weights are loading, so the OS returns `ECONNREFUSED` rather than
HTTP 503. By the time `ctx ai ping` returns, the server is ready
to accept chat requests.

When more than one backend is configured, either pass `--backend`
explicitly or set `default_backend:` in `.ctxrc`:

```bash
ctx ai ping --backend openai
```

If `default_backend:` is unset and multiple backends are
configured, `ctx ai ping` fails with an ambiguous-default error.

### Step 4: Drive a structured extraction

`ctx ai extract` reads free text from stdin and asks the model to
emit a JSON structure of decisions, learnings, tasks, and open
questions. The response lands as a proposal file:

```bash
ctx journal source --limit 5 | ctx ai extract
# proposal written to .context/proposals/20260523T220000Z-extract.md
```

The proposal file is markdown wrapping the raw JSON response
inside a fenced code block. Each candidate keeps its category,
title, and body so a reviewer can sort and apply selectively.

!!! warning "`.context/*.md` is never written directly"
    `ctx ai extract` writes only to `.context/proposals/`. The
    canonical files (`DECISIONS.md`, `LEARNINGS.md`, `TASKS.md`,
    `CONVENTIONS.md`) are read but never modified — a dedicated
    regression test asserts this. Every AI-produced edit goes
    through the proposal queue.

### Step 5: Review and apply

Open the proposal, drop entries you don't want, then re-enter the
surviving candidates with the usual structured commands:

```bash
# Inspect what the model proposed
$EDITOR .context/proposals/20260523T220000Z-extract.md

# Apply selected candidates
ctx decision add "..."
ctx learning add "..."
ctx task add "..."
```

The proposal file stays in `.context/proposals/` as an audit
trail. It is gitignored by default — fresh `ctx init` projects
get `.context/proposals/*` + `!.context/proposals/.gitkeep` in
their `.gitignore` template.

## Tips

* **Pick a model that fits.** vLLM hard-fails when the model
  doesn't fit in GPU memory. Start with a 7B or 8B instruct model
  for a single workstation; reach for 70B-class models only when
  you have the hardware to back it.
* **The cold-start retry is scoped to `Ping`.** Once a `ctx ai
  ping` returns successfully, subsequent `ctx ai extract` calls
  reach a live listener. There's no need for separate retry logic
  on the extract path.
* **Multiple backends are first-class.** Configure both `vllm` and
  `openai` in `.ctxrc`, set `default_backend: vllm`, and override
  per-call with `ctx ai extract --backend openai` when the local
  model isn't strong enough for a specific input.
* **The proposal queue is a draft store, not a tracked artifact.**
  Treat each file as ephemeral: write, review, apply, discard.
  The gitignore carve-out keeps the queue out of version control
  by default.
* **Removing the backend removes the AI layer.** Deleting the
  `backends:` section of `.ctxrc` returns the project to
  deterministic-only mode. The ctx core is unaffected — no
  `ctx agent` or hook run depends on `ctx ai`.
* **`ctx ai` honours timeouts from `.ctxrc`.** Per-backend
  `timeout: 2m` overrides the 30s default; useful when running a
  large local model that takes a while to first-token.

## Troubleshooting

| Symptom                                                            | Likely cause / fix                                                                                                  |
|--------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------|
| `no backend configured; run `ctx setup --backend <name>``          | `.ctxrc` has no `backends:` entries; run `ctx setup --backend vllm`                                                 |
| `backend "vllm" unreachable at http://localhost:8000`              | vLLM isn't running, or it's bound to a different port; check `vllm serve` output                                    |
| `backend "vllm" reachable but lists no models`                     | vLLM is up but didn't load any model; check launch args (a `--model` was likely missed)                             |
| `multiple backends configured; pass `--backend <name>` or set `default_backend:`` | Set `default_backend:` in `.ctxrc` or pass `--backend` per call                                                     |
| Cold-start retry never converges                                   | Window default is 90 seconds; if your weights take longer to load, increase per-backend `timeout:` to match         |

## See also

* [`ctx ai` CLI reference](../cli/ai.md): full flag and behavior reference for `ping` and `extract`
* [`ctx setup` CLI reference](../cli/setup.md): full flag reference for editor-tool templating and `--backend` mode
* [Setup across AI Tools](multi-tool-setup.md): editor-tool integration (Cursor, Aider, Copilot, etc.); the companion to this recipe
