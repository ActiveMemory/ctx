# CLI command examples

Add cobra Example fields to all CLI commands so every `--help`
output shows practical usage examples.

## Scope

- Expand examples.yaml from 5 add-subtype entries to 136 keys
  covering every command
- Wire `Example: desc.Example(cmd.DescKeyXxx)` into all cmd.go
  files and parent.Cmd()
- Namespace add-subtype examples under `add.` prefix to avoid
  key collision with command desc keys
- Update TestExamplesYAMLLinkage to validate against desc key
  constants instead of entry type constants
- Add drift-prevention comments to all three YAML asset files

## Outcome

- `ctx <any-command> --help` shows usage examples
- examples.yaml, commands.yaml, flags.yaml have update reminders
- Task created for flag-name drift detection test
