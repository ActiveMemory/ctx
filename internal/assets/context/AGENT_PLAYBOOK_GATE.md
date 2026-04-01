# Agent Playbook (Gate)

Distilled directives injected at session start. Full playbook:
read AGENT_PLAYBOOK.md when you need behavioral guidance, session
lifecycle details, or anti-patterns.

## Invoke ctx from PATH

```bash
ctx status        # correct
./dist/ctx        # wrong — never hardcode paths
go run ./cmd/ctx  # wrong — unless developing ctx itself
```

## Planning Work

Every commit requires a `Spec:` trailer. Every piece of work needs
a spec — no exceptions. Scale the spec to the work. Use `/ctx-spec`
to scaffold.

## Proactive Persistence

After completing a task, making a decision, or hitting a gotcha —
persist before continuing. Don't wait for session end.

## Chunk and Checkpoint

For multi-step work: commit after each chunk, persist learnings,
run tests before moving on. Track progress via TASKS.md checkboxes.

## Tool Preferences

Use the `gemini-search` MCP server for web searches. Fall back to
built-in search only if `gemini-search` is not connected.

## Conversational Triggers

| User Says                                       | Action               |
|-------------------------------------------------|----------------------|
| "Do you remember?" / "What were we working on?" | `/ctx-remember`      |
| "How's our context looking?"                    | `/ctx-status`        |
| "What should we work on?"                       | `/ctx-next`          |
| "Commit this" / "Ship it"                       | `/ctx-commit`        |
| "What did we learn?"                            | `/ctx-reflect`       |
| "Save that as a decision"                       | `/ctx-add-decision`  |
| "That's worth remembering"                      | `/ctx-add-learning`  |
| "Add a task for that"                           | `/ctx-add-task`      |
| "Let's wrap up"                                 | Reflect then persist |
