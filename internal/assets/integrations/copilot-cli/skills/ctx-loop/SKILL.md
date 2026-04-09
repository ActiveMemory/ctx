---
name: ctx-loop
description: "Generate autonomous iteration loop scripts for headless AI tool runs with configurable completion signals."
tools: [bash, read, write]
---

Generate shell scripts for autonomous AI iteration loops.

## When to Use

- Setting up CI-driven AI workflows
- When a task needs autonomous iteration with checks
- For batch processing with verification gates

## When NOT to Use

- Interactive work (just do it in the session)
- Simple single-run tasks
- When safety checks aren't defined

## Process

### 1. Define the loop

Gather:
- **Command**: what to run each iteration
- **Completion signal**: how to detect "done" (exit code, output pattern, file exists)
- **Max iterations**: safety limit (default: 10)
- **Checkpoint command**: what to run between iterations

### 2. Generate the script

```bash
#!/bin/bash
set -euo pipefail
MAX_ITER=${1:-10}
for i in $(seq 1 "$MAX_ITER"); do
  echo "=== Iteration $i/$MAX_ITER ==="
  # Run the task
  {command}
  # Check completion
  if {completion_check}; then
    echo "✅ Complete after $i iterations"
    exit 0
  fi
  # Checkpoint
  {checkpoint}
done
echo "❌ Max iterations reached"
exit 1
```

### 3. Write and verify

Write to the requested location. Verify with `bash -n`.

## Quality Checklist

- [ ] Max iterations has a sane default
- [ ] Completion signal is well-defined
- [ ] Script has `set -euo pipefail`
- [ ] Script passes `bash -n` syntax check
