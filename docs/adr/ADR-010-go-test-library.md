# ADR-010 — Use testify for Go testing

**Status:** Accepted  
**Date:** 2026-03-18

## Context

The backend needs a testing library decision. Go ships with a built-in testing package but third-party libraries offer additional assertion utilities and test organization features. This decision affects how tests are written, how failures are reported, and how test preconditions are enforced.

## Decision

Use testify (`github.com/stretchr/testify`) for Go tests, specifically the `assert` and `require` packages.

## Options Considered

**Standard library (`testing` package)**

- Pros: No external dependency, always available, idiomatic Go, no version management, lighter weight
- Cons: No assertion helpers, no distinction between fatal and non-fatal assertions, more boilerplate per test case

**testify (`assert` + `require`)**

- Pros: Rich assertion library with clear failure messages, `require` stops test execution on failure (fatal), `assert` continues test execution (non-fatal), table-driven tests are cleaner, failure output includes expected vs actual values automatically, widely adopted in the Go ecosystem
- Cons: External dependency, slightly non-idiomatic Go purists would argue against it, adds to go.mod

## Rationale

testify's `require` and `assert` packages solve a real problem in Go testing. The standard library requires manual `t.Fatalf` and `t.Errorf` calls with hand-formatted messages. testify provides clear, consistent failure output with expected vs actual values automatically. The `require` vs `assert` distinction is meaningful — `require` gates preconditions (if setup fails, stop immediately), `assert` handles independent verifications (continue running to collect all failures). This distinction maps cleanly to how tests should be structured and is worth the external dependency. testify is the most widely adopted Go testing library and is considered standard practice in the Go community.

## Consequences

- Cleaner, more readable test assertions
- Failure messages include expected vs actual values automatically
- `require` stops test execution when preconditions fail which prevents misleading downstream failures
- `assert` continues test execution and collects all failures in one run
- One external dependency added to go.mod
- Test style is consistent with the broader Go ecosystem
