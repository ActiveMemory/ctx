# Source Coverage

Per-source completeness ledger maintained by `/ctx-kb-ingest`.
Each row tracks coverage state, not just existence: a source can
be at any of several stages between discovery and comprehensive
understanding. Shape and state-machine transitions live in
`../ingest/schemas/source-coverage.md`.

This file is the canonical answer to *"what is incomplete?"*

| Source         | Topic | State                | EV coverage     | Residue                                                                                                                                                                  | Next action                                                                                       | Updated    |
|----------------|-------|----------------------|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------|------------|
| VLLM-EXAMPLES  | vllm  | topic-page-drafted   | EV-001..EV-005  | Per-category GitHub directories (basic/, generate/, pooling/, speech_to_text/, features/, reasoning/, tool_calling/, applications/, rl/, deployment/, ray_serving/, disaggregated/, observability/) were not fetched. Page-level only. | `/ctx-kb-ingest https://github.com/vllm-project/vllm/tree/main/examples/deployment vllm` (or any single category to deepen the topic) | 2026-05-21 |
