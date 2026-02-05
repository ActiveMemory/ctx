![ctx](assets/ctx-banner.png)

# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in `ctx`, please report it responsibly.

**Do NOT open a public issue for security vulnerabilities.**

Instead, please use one of these methods:

### Email

Send details to **security@ctx.ist**

### GitHub Private Reporting

1. Go to the [Security tab](https://github.com/ActiveMemory/ctx/security) of
   this repository
2. Click "Report a vulnerability"
3. Provide a detailed description of the issue

### What to Include

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (*if any*)

### Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial assessment**: Within 7 days
- **Resolution target**: Within 30 days (*depending on severity*)

## Trust Model

`ctx` operates within a single trust boundary: **the local filesystem**.

The person who authors `.context/` files is the same person who runs the
agent that reads them. There is no remote input, no shared state, and no
server component.

This means:

* **ctx does not sanitize context files for prompt injection.** This is a
  deliberate design choice, not an oversight. The files are authored by the
  developer who owns the machine: Sanitizing their own instructions back
  to them would be counterproductive.
* **If you place adversarial instructions in your own `.context/` files,
  your agent will follow them.** This is expected behavior. You control the
  context; the agent trusts it.
* **Shared repositories should review `.context/` files in code review**,
  the same way you would review any committed configuration. A malicious
  contributor could add harmful instructions to `CONSTITUTION.md` or
  `TASKS.md`: Treat these files with the same scrutiny as CI/CD config
  or Makefiles.

## Security Considerations

`ctx` is designed with security in mind:

* **No secrets in context**: The constitution explicitly forbids storing
  secrets, tokens, API keys, or credentials in `.context/` files
* **Local only**: `ctx` runs entirely locally with no external network calls
* **No code execution**: `ctx` reads and writes Markdown files only; it does
  not execute arbitrary code
* **Git-tracked**: All context files are meant to be committed, so they should
  never contain sensitive data

### Best Practices

1. **Review before committing**: Always review `.context/` files before 
   committing
2. **Use .gitignore**: If you must store sensitive notes locally, add them 
   to `.gitignore`
3. **Drift detection**: Run `ctx drift` to check for potential secrets in 
   your project

## Attribution

We appreciate responsible disclosure and will acknowledge security researchers
who report valid vulnerabilities (*unless they prefer to remain anonymous*).

