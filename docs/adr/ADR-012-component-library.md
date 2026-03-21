# ADR-012 — Use Radix UI for headless component primitives

**Status:** Accepted  
**Date:** 2026-03-21

## Context

The frontend needs interactive components e.g modals, dialogs, dropdowns, select menus, that meet accessibility standards and work correctly across browsers and assistive technologies. Building these from scratch is deceptively complex: accessible modals require focus trapping, escape key handling, scroll locking, correct aria attributes, and proper focus restoration on close. The choice affects how much accessibility work is done manually versus handled by a library, and how much control is retained over visual styling.

## Decision

Use Radix UI primitives for complex interactive components that have significant accessibility requirements. Style all Radix components using CSS Modules and design tokens consistent with ADR-011.

## Options Considered

**Build from scratch**

- Pros: Zero dependencies, complete control over implementation, no library API to learn
- Cons: Accessible modals, dialogs, and dropdowns are genuinely hard to implement correctly: focus trapping, aria attributes, keyboard navigation, and scroll locking are all easy to get wrong, every component requires significant effort to meet WCAG standards, bugs are likely and hard to catch without dedicated accessibility testing

**Radix UI**

- Pros: Zero visual opinions, all styling done with CSS Modules and design tokens, accessibility built in (focus management, keyboard navigation, aria attributes, scroll locking), battle tested across thousands of production applications, install only the primitives you need, no conflict with ADR-011 CSS approach
- Cons: Additional dependency, requires learning Radix's composition API, some components have specific structural requirements

**shadcn/ui**

- Pros: Pre-styled components ready to use, built on Radix so accessibility is included, you own the copied source code
- Cons: Built on Tailwind which.directly conflicts with ADR-011 which chose CSS Modules and design tokens, the default aesthetic is clean and minimal which fights against the design system, overriding styles would require more work than starting from scratch with Radix

## Rationale

Radix UI is the correct fit for this project. It provides accessibility primitives with no styling opinions, which means all visual styling is done with CSS Modules and design tokens exactly as established in ADR-011. The sci-fi aesthetic is fully preserved — Radix contributes behavior, not appearance. Building accessible interactive components from scratch is not a good use of time when a battle-tested solution exists. shadcn/ui was eliminated because its Tailwind dependency directly conflicts with the styling architecture already established.

## Consequences

- Accessible interactive components without manual aria and focus management
- No visual opinions imposed
- CSS Modules and design tokens remain the single styling approach
- Additional `@radix-ui/*` dependencies added per primitive as needed
- Radix's composition API must be learned for each primitive used
- Components built on Radix will have correct keyboard navigation and screen reader support from the start
