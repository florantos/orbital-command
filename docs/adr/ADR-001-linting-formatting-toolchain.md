# ADR-001 — Use ESLint and Prettier for linting and formatting

**Status:** Accepted  
**Date:** 2026-03-17

## Context

The frontend needs a linter for catching bugs, performance issues, and enforcing coding standard and best practices. It also needs a formatter to maintain code style consistency across the frontend.

## Decision

Use ESLint and Prettier together for linting and code formatting.

## Options Considered

**Biome**

- Pros: A Single tool that handles both linting and formatting, significatly faster than ESLint + Prettier, handles conflicts between linter and formatter automatically, import sorting built in, includes most ESLint rules/plugins, growing ecosystem
- Cons: Less battle tested in large production codebases, younger and smaller ecosystem, missing React Compiler lint rules with no committed roadmap to add them

**ESLint & Prettier**

- Pros: Massive ecosystem, battle tested across millions of projects, complete React Compiler lint rules support, familiar for every developer
- Cons: Two tools to configure and maintain, need to handle conflicts between ESLint and Prettier manually, slower than Biome, more configuration overhead needed upfront

## Rationale

Our project uses React with the React Compiler, so not having full React Compiler linting rules disqualifies Biome. The speed gains as well as the quicker setup and lack of configuration overhead for Biome is attractive. Biome is something we can re-evaluate later should it add full support for React Compiler. ESLint and Prettier is also familiar to most if not all developers.

## Consequences

- ESLint and Prettier must be configured to not conflict
- The initial setup overhead for ESLint and Prettier is real but mitigated mostly since it comes built in with a baseline configuration with Vite
- Eslint is substantially slower than Biome and it will become noticable as the project grows.
