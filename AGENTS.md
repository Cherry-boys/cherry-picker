## Startup rule

Before doing any substantial work in this repository:

1. Run `pnpm docs:list`.
2. Read the docs that match the current task.
3. Only then inspect code and make changes.

## Expectation

When starting a task, mention briefly which docs were checked or why the docs listing step could not run.

## Build, Test, and Development Commands
- Use `pnpm` (Corepack-enabled) for repository scripts.
- Run all unit + integration tests across projects: `pnpm run test:all`.
- Run only integration tests across projects: `pnpm run test:integration:all`.
- Run only unit tests across projects: `pnpm run test:unit:all`.
- Run type-check across Nx projects: `pnpm run typecheck:all`.
- Run docs index helper before starting substantial work: `pnpm run docs:list`.
