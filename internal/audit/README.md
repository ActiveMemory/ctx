# internal/audit/

AST-based codebase invariant tests. Read this before changing any file
in this directory.

## The codebase is clean. Keep it that way.

Every test in this package passed with **zero violations** when last
reviewed. There are no grandfathered exceptions, no "we'll fix it
later" carve-outs, and no blanket package exemptions.

## If a test fails after your change

**Fix your code, not the test.** The test is enforcing a project
convention that existed before your change. If your code introduces
a new violation, the correct response is to make your code conform.

## If you genuinely need an exception

Exceptions are not added inline. The process is:

1. **Create a dedicated PR** that contains _only_ the exception
   change (the allowlist entry, the grandfathered bump, etc.).
2. **Justify every entry** in the PR description: what it is, why
   it cannot conform, and when (if ever) the exception will be
   removed.
3. **Never widen an allowlist to make a drive-by test pass.** That
   is the single most common way these checks degrade.

## If you are an AI agent

**STOP. Read this file before touching any `_test.go` in this
directory.**

- Do NOT add packages to `exemptStringPackages` or similar
  allowlists.
- Do NOT bump `grandfathered*` counters.
- Do NOT add entries to `linuxOnlyExports` or any other bypass map.
- Do NOT weaken a check to make your code pass.

If a test fails, the fix belongs in the code under test, not here.
If you believe an exception is truly warranted, surface it to the
user and let them decide; do not silently widen a check.
