# Code Review

Review this code change focusing on:

- **Correctness**: Does the logic do what it claims? Are there off-by-one errors, nil dereferences, or race conditions?
- **Edge cases**: What happens with empty input, max values, concurrent access, or partial failures?
- **Naming clarity**: Do function, variable, and type names communicate intent without needing comments?
- **Test coverage gaps**: What behavior is untested? What inputs would exercise uncovered paths?
- **Convention adherence**: Does this follow the project patterns documented in `.context/CONVENTIONS.md`?

Flag but don't fix style issues. Focus your review on substance over formatting.
