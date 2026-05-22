# Evidence Index

Append-only ledger of atomic claims minted by `/ctx-kb-ingest`.
Shape lives in `../ingest/schemas/evidence-index.md`. Every
`EV-###` cited from a topic page, glossary entry, timeline event,
contradiction, or domain-decision MUST resolve to a row here.

| id     | claim                                                                                                                                                          | source         | locator              | sha | confidence | tags                       | occurred | extracted  |
|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------|----------------------|-----|------------|----------------------------|----------|------------|
| EV-001 | vLLM organizes its public examples landing page into thirteen named categories: basic, generate, pooling, speech_to_text, features, reasoning, tool_calling, applications, rl, deployment, ray_serving, disaggregated, observability. | VLLM-EXAMPLES  | category table       |     | high       | vllm, taxonomy             |          | 2026-05-21 |
| EV-002 | The vLLM examples landing page contains no inline code; each of its thirteen categories links to a directory under github.com/vllm-project/vllm/tree/main/examples/, and the example files live next to source. | VLLM-EXAMPLES  | category table links |     | high       | vllm, docs-structure       |          | 2026-05-21 |
| EV-003 | vLLM's example taxonomy operates on two dimensions simultaneously without labeling them: capability (generate, pooling, reasoning, tool_calling, speech_to_text) and deployment context (basic, deployment, ray_serving, disaggregated, observability), with crosscut entries (features, applications, rl). | VLLM-EXAMPLES  | category table       |     | medium     | vllm, taxonomy, analysis   |          | 2026-05-21 |
| EV-004 | The vLLM examples landing page contains no numbered workflows, no procedural recipes, and no "if you want X follow steps 1..N" framing; categories are presented as a flat bulleted index.                       | VLLM-EXAMPLES  | full page            |     | high       | vllm, docs-structure, contrast |       | 2026-05-21 |
| EV-005 | The vLLM examples landing page introduces its taxonomy with the verbatim line "vLLM's examples are organized into the following categories", presenting taxonomy-as-taxonomy with no entry-point guidance.       | VLLM-EXAMPLES  | intro paragraph      |     | high       | vllm, docs-structure       |          | 2026-05-21 |
