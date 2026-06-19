# ctx ai backend

## Problem

`ctx` today has no AI backend abstraction. Users who want local-first
AI capabilities (running on vLLM, Ollama, LM Studio, or a
self-hosted OpenAI-compatible server) have no first-class wiring
through ctx — they configure their downstream AI tool (Claude Code,
OpenCode, etc.) by hand. Two user-visible pains follow:

- **Users may not know what they do not know.** Inference is opaque
  to non-ML users; "set `ANTHROPIC_BASE_URL` to your vLLM endpoint"
  is not a recipe non-experts can act on without help.
- **Cloud LLMs break air-gap.** `ctx`'s thesis (Invariant 5,
  local-first / air-gap capable) is currently honoured only for the
  *persistence* layer. AI-shaped capabilities are off-limits to
  classified, isolated, or constrained-environment projects because
  the only paths are cloud APIs.

Block A is the foundation. It provides the abstraction that blocks
B (structured extraction) and C (embedding-backed recall) build on
later, while staying inside the six invariants (Markdown-on-filesystem,
zero runtime deps for core functionality, deterministic assembly,
human authority, local-first/air-gap, no default telemetry).

## Approach

`ctx` grows an **optional, local-first AI backend layer** that talks
to any OpenAI-compatible HTTP endpoint. vLLM is the canonical local
backend (it restores air-gap capability, supports
schema-constrained structured outputs, and ships prefix caching that
rewards stable-prefix prompt structure). The same contract works
against OpenAI, Anthropic, Ollama, LM Studio, and any other
OpenAI-compatible server.

The layer is **strictly additive**:

- Existing commands (`ctx status`, `ctx agent`, ceremonies, hooks)
  keep working with no backend configured. The deterministic core
  is untouched.
- AI commands fail **closed** with a clear "no backend reachable"
  error rather than degrading silently to a non-AI path. There is
  no "use vLLM if available, fall back to deterministic" behaviour.
- The contract floor is OpenAI-compatible HTTP. Anthropic Messages
  is treated as a *strict superset* that some backends also support,
  not as a competing contract.

Block A delivers four things:

1. A **backend registry** with one entrypoint per known backend type
   (`vllm`, `openai`, `anthropic`, `ollama`, `lmstudio`, generic
   `openai-compatible`).
2. A `ctx setup --backend <name>` extension to the existing setup
   family that templates endpoint + auth wiring into `.ctxrc` and
   (where applicable) downstream AI-tool configs.
3. An AI command surface (shape TBD — see Open Questions) that
   exposes one or more verbs against the configured backend. The
   minimum viable verb set is `ping` (reachability) plus *one*
   structured-output consumer from block B chosen during A's
   implementation to validate the pattern end-to-end (see Testing).
4. Reachability/health checks and fail-closed error reporting.

## Behavior

### Happy Path

1. User runs `ctx setup --backend vllm --endpoint http://localhost:8000`.
   This is a backend setup mode on the existing `ctx setup <tool>` family:
   implementation must relax the current exactly-one-positional tool shape so
   `--backend` can run without a tool argument. ctx writes the backend
   definition to `.ctxrc` under the YAML `backends:` mapping.
2. User runs `ctx ai ping` (or equivalent). ctx reads the backend
   definition, performs a reachability check against the endpoint
   (HTTP GET on `/v1/models` for OpenAI-compatible servers), and
   reports the backend name, endpoint, and first model listed.
3. User runs the validation consumer command from block B (e.g.,
   `ctx compact --emit decisions,learnings,tasks <input>` — exact
   shape lives in the B+C supplementary spec). ctx routes the
   request through the configured backend, receives a
   schema-constrained JSON response, and writes a *proposed-patch*
   artifact under the proposal queue (location TBD; see Open
   Questions) — never directly to `.context/*.md`.
4. Existing ceremonies (`/ctx-remember`, `/ctx-wrap-up`, etc.) work
   exactly as before. `ctx status` and `ctx agent` are untouched.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Backend unreachable | Fail closed with a clear error naming the configured endpoint and suggesting `ctx setup --backend <name> --endpoint <url>`. No fallback. |
| Backend reachable but model unavailable | Surface upstream 4xx verbatim; do not retry with a different model. |
| Multiple backends configured (e.g., `vllm` + `openai`) | User must specify `--backend <name>` on the AI command, or set a default via `.ctxrc` `backends.default`. No implicit selection. |
| No backend configured | AI commands print: "no backend configured; run `ctx setup --backend <name>`" and exit non-zero. Non-AI commands are unaffected. |
| API key missing or invalid | Surface upstream auth error verbatim; do not retry. Suggest the env-var or `.ctxrc` key the backend reads. |
| Slow backend | Respect timeout from `.ctxrc` `backends.<name>.timeout` (default TBD). No infinite waits. |
| AI command invoked from inside a deterministic ceremony hook | Fail closed. Coupling deterministic-core hooks to AI availability would violate Invariant 2 ("zero runtime deps for core functionality"). |
| Existing `ANTHROPIC_BASE_URL` already set in user env | Honour it; do not overwrite. `ctx setup --backend` prints a warning if the env vars it would template conflict with what's already set. |
| `.ctxrc` malformed (e.g., missing required key) | Refuse with a clear parse error naming the offending key; do not silently default. |

### Validation Rules

| Field | Rule | Enforced where |
|-------|------|----------------|
| Backend name | Alphanumeric + hyphen; must match a registered backend type | At setup time and at AI-command dispatch |
| Endpoint URL | Must parse as `http://` or `https://`; localhost recommended (not required) for `vllm`-canonical backend | At setup time |
| API key (if any) | Read from env-var (preferred) or `.ctxrc`; **never** allowed in a committed `.ctxrc` if the project's git config marks it as such | Setup-time warning; commit hook (out of scope here, but flagged in Non-Goals) |
| Default backend | If `backends.default` is set, must reference a configured backend | At AI-command dispatch |
| Determinism boundary | `ctx ai` commands must not be invoked by `ctx agent`, `ctx status`, or any hook that fires during deterministic ceremony paths | Unit test guard (see Testing) |

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| No backend configured | `no backend configured; run \`ctx setup --backend <name>\`` | Run setup |
| Backend unreachable | `backend \`<name>\` unreachable at <endpoint>: <underlying error>` | Check endpoint; verify vllm/ollama/etc. is running |
| Model not found | (relay upstream 4xx body verbatim) + `backend \`<name>\` rejected the model selection; check \`/v1/models\` on the endpoint` | Pick a listed model |
| Auth failed | (relay upstream 401/403 verbatim) + `backend \`<name>\` rejected the credential; check <env-var or .ctxrc key>` | Update credential |
| Timeout | `backend \`<name>\` timed out after <duration>; tune \`backends.<name>.timeout\` in .ctxrc` | Increase timeout or use a faster model |
| Multiple backends, none specified | `multiple backends configured; pass \`--backend <name>\` or set \`backends.default\` in .ctxrc` | Pass flag or set default |
| AI command called from deterministic hook | (developer-only) `ctx ai called from deterministic context; this would violate Invariant 2` | Restructure hook |

## Interface

### CLI

**Open question — the brief explicitly punts this to A's spec
author.** Two shapes were enumerated; A's implementation must
commit to one. The decision is expensive to unwind (users will
script against whichever ships).

**Option 1 — new top-level `ctx ai` namespace:**

```
ctx ai ping [--backend <name>]
ctx ai <verb> [--backend <name>] [verb-specific flags]
```

**Option 2 — flags on existing commands:**

```
ctx compact <input> --emit <kinds> --use-ai [--backend <name>]
ctx ingest <input> --extract <kinds> [--backend <name>]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--backend` | string | (resolved from `.ctxrc` `backends.default`) | Selects which configured backend to dispatch through |
| `--endpoint` (setup only) | URL | (per-backend) | Endpoint override at setup time |
| `--api-key-env` (setup only) | string | (per-backend) | Name of the env-var the backend reads for auth |
| (additional flags TBD with interface decision) | | | |

### Skill (if applicable)

No new companion skill ships in Block A. The user-facing surface is covered by
the `ctx setup` docs, the `ctx ai` CLI reference, the vLLM recipe, command
assets/examples, and the agent playbook note listed below. A future `/ctx-ai-*`
skill can be specified after Block A proves the backend contract.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/backend/` | **New package.** Backend registry, contract types (`Backend`, `Request`, `Response`), per-backend implementations (`vllm.go`, `openai.go`, `anthropic.go`, `ollama.go`, `lmstudio.go`, `openaicompat.go`) |
| `internal/cli/ai/` | **New package** (if Option 1). Command surface for the `ai` namespace |
| `internal/cli/setup/cmd/root/` | Extend with `--backend` handling; templating into `.ctxrc` |
| `internal/cli/setup/core/backend/` | **New subpackage.** Setup-time wiring per backend type (env-var templates, downstream-tool config writes) |
| `internal/rc/` | Add `backends:` YAML mapping parsing and validation |
| `.ctxrc` (project-init template) | Add commented-out `backends:` skeleton |
| `internal/assets/context/AGENT_PLAYBOOK.md` | (TBD) note that `ctx ai <verb>` exists and when agents should call it vs. hand-rolling against the AI tool |
| `docs/recipes/` | New recipe `local-inference-with-vllm.md` (or `ai-backend-setup.md`); the user explicitly carved out recipe-restructuring work, so this is *one* recipe added, not a recipe-surface rework |
| `docs/cli/` | New page documenting `ctx ai` (or the chosen flag surface) |

### Key Functions

TBD pending interface decision. Skeleton contract:

```go
// internal/backend/contract.go
type Backend interface {
    Name() string
    Ping(ctx context.Context) error
    Complete(ctx context.Context, req Request) (Response, error)
}

type Registry interface {
    Register(name string, factory func(cfg Config) (Backend, error))
    Resolve(name string) (Backend, error)
    Default() (Backend, error) // honours .ctxrc [backends].default
}
```

### Helpers to Reuse

- `internal/cli/setup/` — existing setup family (Claude Code,
  OpenCode, Cursor, etc.). The `--backend` extension follows the
  same templating-into-config-files pattern.
- `internal/rc/` — `.ctxrc` parsing and validation. Adding a
  `backends:` mapping follows the existing YAML pattern.
- `internal/err/` — typed-string sentinels for backend errors
  (per the recent `entity.Sentinel` convention).
- `internal/assets/commands/text/errors.yaml` — externalised
  error strings.

## Configuration

`.ctxrc` additions use the repo's existing YAML shape:

```yaml
backends:
  default: vllm # optional; required only if more than one backend is configured
  vllm:
    endpoint: http://localhost:8000
    api_key_env: "" # vllm typically runs without auth; empty means none
    timeout: 30s
    default_model: openai/gpt-oss-120b
  openai:
    endpoint: https://api.openai.com
    api_key_env: OPENAI_API_KEY
    timeout: 60s
    default_model: gpt-4o
```

Environment variables: only **read** from env-vars named in
`api_key_env`. ctx never writes credentials to `.ctxrc`.

## Testing

- **Unit:** Backend registry resolution (single, multiple, default,
  missing); `.ctxrc` parsing (well-formed, missing required key,
  malformed table); error sentinels match expected strings.
- **Integration:** Spin up a fake OpenAI-compatible HTTP server (a
  tiny `httptest.Server` per backend test) and drive each backend's
  `Ping` + one `Complete` call against it. Verify fail-closed
  behaviour when the server returns 401/404/500/timeout.
- **Edge cases:** Backend unreachable; multiple backends without
  selector; no backend configured; AI-command-from-deterministic-hook
  guard (compile-time or test-time check that `ctx agent` /
  `ctx status` / canonical hooks do not import `internal/backend`).
- **Validation consumer:** One end-to-end test that exercises the
  chosen B-block validation command (`ctx compact ... --emit
  decisions,learnings,tasks` or equivalent) against the fake
  server, asserts a *proposed-patch* artifact is written to the
  proposal queue, and asserts `.context/*.md` files are
  **unchanged**.

## Non-Goals

This spec explicitly does **not** cover:

- **Cost management.** No usage tracking, no billing integration,
  no cost attribution. Out of manifesto; would dilute the product.
- **Secrets management.** ctx reads from env-vars or `.ctxrc`; it
  does not encrypt, rotate, or vault credentials. Vault/Doppler/etc.
  are separate products.
- **Observability backends.** No metrics emission, no time-series
  storage, no dashboards. The thesis's "no default telemetry"
  (Invariant 6) governs.
- **Multi-backend routing daemon / HTTP proxy.** Killed in the
  brief — structurally cannot solve the failure modes it was
  proposed for, and forces ctx into service-shape.
- **Embeddings, semantic search, vector storage.** Block C.
  Supplementary spec.
- **Structured extraction commands.** Block B (other than the *one*
  validation consumer needed to prove A works end-to-end).
  Supplementary spec.
- **Automatic application of AI-generated content to `.context/*.md`.**
  All AI-produced edits land as *proposed patches* in a review
  queue; ratification is human (or agent) via existing ceremony
  paths.
- **Replacing `/ctx-wrap-up` or any existing ceremony skill.** AI
  commands augment ceremonies via the proposal queue; they do not
  short-circuit them.
- **Touching `ctx agent` or its deterministic assembly.** Sibling,
  not replacement.
- **Recipe-surface restructuring.** Adding *one* recipe for AI
  backend setup is in scope; mapping vLLM's index-style taxonomy
  onto `docs/recipes/` is explicitly out (rejected in the brief).

## Open Questions

These were left open in the brief and must be settled during A's
implementation:

1. **CLI namespace shape:** `ctx ai <verb>` (Option 1) vs. flags
   on existing commands (Option 2). Expensive to unwind. Pick once.
2. **Proposal queue location:** Block A uses `.context/proposals/ai/` as the
   provisional queue. The directory is tracked as substrate metadata; each
   validation run writes one JSON proposed-patch artifact containing backend,
   model, input reference, emit kinds, proposed rows, source spans/citations
   when available, and status metadata. B+C may confirm or relocate this queue
   after Block A ships.
3. **Default extraction model:** A-spec leaves model choice to the
   user; recommended models per task type can be a recipe.
4. **Companion skill:** No new Block A skill; docs/assets/playbook updates cover
   the surface until a dedicated skill is specified.
5. **Validation consumer:** `ctx ai propose <input> --emit ...` is the Block A
   validation-only generic proposer. It proves backend dispatch and proposed
   patch writing without claiming the final B command taxonomy or foreclosing
   later `ctx ai compact` / `ctx ai ingest` commands.

## Task Breakdown

Paste-ready rows for `.context/TASKS.md`. Each is a single block.
Ordering reflects what blocks what — earlier rows must land before
later rows. The `Spec:` reference points at this file.

- [ ] Decide the AI command CLI namespace: `ctx ai <verb>` (new
  top-level) vs. flags on existing commands (`--use-ai`, `--emit`,
  etc.). Foundational; expensive to unwind once shipped. Record
  the call as a `.context/DECISIONS.md` entry naming the chosen
  shape and the rejected alternative with rationale. Blocks every
  other task in this group. Spec: `specs/ctx-ai-backend.md` Open
  Question #1. #priority:medium #added:2026-05-21

- [ ] Implement the backend contract and registry: new
  `internal/backend/` package with `Backend` interface (`Name`,
  `Ping`, `Complete`), `Request`/`Response` types, and a
  `Registry` (`Register`, `Resolve`, `Default`). No per-backend
  implementations yet; this is the abstraction surface that later
  tasks plug into. Unit tests cover single/multiple/default/missing
  backend resolution. Spec: `specs/ctx-ai-backend.md` §Implementation.
  #priority:medium #added:2026-05-21

- [ ] Extend `internal/rc/` to parse and validate the `.ctxrc`
  `backends:` mapping per the spec's Configuration section:
  per-backend `endpoint`, `api_key_env`, `timeout`, `default_model`,
  plus optional `backends.default`. Refuse malformed tables with
  a clear parse error naming the offending key. Add fixtures and
  round-trip tests. Spec: `specs/ctx-ai-backend.md` §Configuration.
  #priority:medium #added:2026-05-21

- [ ] Implement the minimum viable backend set: `vllm` (canonical
  local) and generic `openai-compatible` (the contract floor) in
  `internal/backend/vllm.go` and `internal/backend/openaicompat.go`.
  Both must implement `Ping` (HTTP GET on `/v1/models`) and
  `Complete` (POST `/v1/chat/completions`). Fail closed on
  unreachable / 4xx / 5xx / timeout; never retry with a different
  model. Spec: `specs/ctx-ai-backend.md` §Approach and §Edge Cases.
  #priority:medium #added:2026-05-21

- [ ] Add the named-backend implementations: `openai`, `anthropic`,
  `ollama`, `lmstudio` in `internal/backend/`. Each is a thin
  wrapper over `openaicompat` with backend-specific defaults
  (endpoint, auth header shape, env-var name). Anthropic uses the
  Messages API endpoint where supported but inherits the
  OpenAI-compatible floor for `/v1/chat/completions`. Spec:
  `specs/ctx-ai-backend.md` §Approach. #priority:medium
  #added:2026-05-21

- [ ] Extend the `ctx setup` family with `--backend <name>`:
  templates endpoint + auth wiring into `.ctxrc` and (where
  applicable) downstream AI-tool configs (`ANTHROPIC_BASE_URL`,
  `OPENAI_BASE_URL`). Honours existing env-var values: warn but
  do not overwrite. Lives in new `internal/cli/setup/core/backend/`
  subpackage. Spec: `specs/ctx-ai-backend.md` §Implementation.
  #priority:medium #added:2026-05-21

- [ ] Build the AI command surface per the namespace decision from
  the first task. Minimum verbs: `ping` (reachability + first model
  listed) plus `propose` as the Block A validation-only generic
  proposer. All AI
  commands honour `--backend` flag (falls back to
  `backends.default`), fail closed when no backend configured,
  and surface upstream errors verbatim. Spec:
  `specs/ctx-ai-backend.md` §Interface. #priority:medium
  #added:2026-05-21

- [ ] Add the deterministic-core boundary guard: a unit test (or
  lint check) that fails if `internal/cli/agent/`,
  `internal/cli/status/`, or any deterministic-ceremony hook
  imports `internal/backend/`. This is the structural enforcement
  for Invariant 2 — without it, the additive/optional discipline
  is honour-system only. Spec: `specs/ctx-ai-backend.md` §Validation
  Rules and §Testing. #priority:medium #added:2026-05-21

- [ ] Ship the Block A validation consumer: `ctx ai propose <input>
  --emit decisions,learnings,tasks,open-questions`. This is a
  validation-only generic proposer, not the final B command taxonomy.
  Implements the full pattern end-to-end:
  schema-constrained dispatch through the backend, JSON validation,
  one proposed-patch JSON artifact written to `.context/proposals/ai/`
  with backend, model, input reference, emit kinds, proposed rows,
  source spans/citations when available, and status metadata.
  `.context/*.md` files must remain unchanged. Integration test confirms the round-trip
  against a fake OpenAI-compatible httptest server. Spec:
  `specs/ctx-ai-backend.md` §Testing and Open Question #5.
  #priority:medium #added:2026-05-21

- [ ] Write the documentation deliverables: one new recipe
  (`docs/recipes/local-inference-with-vllm.md` or
  `docs/recipes/ai-backend-setup.md`) covering the
  `ctx setup --backend vllm` flow end-to-end, plus a CLI reference
  page under `docs/cli/` for whichever command surface the
  namespace decision produced. The recipe is *one file*, not a
  recipe-surface rework — that scope was explicitly rejected in
  the brief. Spec: `specs/ctx-ai-backend.md` §Non-Goals.
  #priority:medium #added:2026-05-21

After all ten land, the **B + C supplementary spec
(`specs/ctx-ai-extraction-and-recall.md`) gets re-debated** with
A's surface as ground truth, then promoted to contract specs of
its own before B/C implementation work is broken into tasks.
