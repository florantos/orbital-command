# ADR-002 — Use Lefthook as a pre-commit hook tool

**Status:** Accepted  
**Date:** 2026-03-17

## Context

The project needs a way to ensure that ensure the codebase is complying with linting, formatting, and testing standards.

## Decision

Use Lefthook for pre-commit hooks for both backend and frontend.

## Options Considered

**Husky**

- Pros: Large ecosystem, simple setup, familar to most developers
- Cons: Javascript/Node specific (doesn't work outside of JS/Node projects), slower than Lefthook, one file per hook

**Lefthook**

- Pros: Runs hooks in parallel, faster than Husky, single config file, language agnostic
- Cons: smaller ecosystem, less widely known

## Rationale

Our project needs a way to enforce code quality in the backend and the frontend. We use Go and Typescript making lefthook the obvious choice. The speed gains and configuration simplicity is a plus.

## Consequences

- Lefthook is less known so there is a small learning curve
- lefthook is younger and has a smaller ecosystem
- Hooks run in parallel so pre-commit checks are faster
- No Node depency required to run hooks
