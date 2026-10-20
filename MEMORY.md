<!--
@file MEMORY.md
@description Routes tokenizer project context and records durable decisions, current state, and metadata rules.
@depends_on plan.md, collaboration-guide.md
@depended_on_by AGENTS.md, README.md, scripts/check_metadata.py
-->

# Tokenizer project memory

## Purpose

This is the compact recovery point for humans and AI collaborators. It records where authoritative context lives, how files relate, which decisions are durable, and what has actually been completed. It is not a transcript or a substitute for source code and tests.

## Context recovery

Start with the file index below. Read metadata headers for likely relevant files, follow only direct edges that can affect the task, and open full files only when they constrain or will be changed by the work.

Use this source-of-truth order when information conflicts:

1. executable code, tests, fixtures, and observed process behavior;
2. approved contracts and focused technical documentation;
3. the untouched source material under `resources/`;
4. this memory and its durable decisions;
5. temporary plans, chat history, and assumptions.

## Metadata standard

Every human-authored file must declare `@file`, `@description`, `@depends_on`, and `@depended_on_by` within its first ten lines. Paths are repository-relative, comma-separated, exact, and reciprocal. Use `none` when a list is empty.

Generated files, dependencies, build output, lockfiles, binaries, media, and formats that cannot safely contain comments are exempt. The CEO-provided PDF is intentionally unchanged and is represented by `resources/README.md`.

## Two-pass workflow

Before changing files, identify the smallest change set, inspect direct dependencies and consumers, state assumptions and success criteria, and choose proportionate verification. After verification, reconcile every changed responsibility and dependency edge, update this index and durable decisions when warranted, then run `python3 scripts/check_metadata.py`.

## File index

| File | Responsibility | Depends on | Depended on by |
| --- | --- | --- | --- |
| `.gitignore` | Excludes generated and machine-local files from version control. | — | — |
| `AGENTS.md` | Gives any AI collaborator the repository operating contract. | Memory, plan, collaboration guide, metadata validator | — |
| `MEMORY.md` | Routes context and records durable project state and decisions. | Plan, collaboration guide | Agents, root README, metadata validator |
| `README.md` | Provides the human entry point, current phase, and navigation. | Memory, plan, docs index, stage workflow, challenge audit, source orientation, tokenizer foundations, research plan, version control, mental-model canvas | — |
| `plan.md` | Defines outcomes, boundaries, phases, and completion criteria. | Collaboration guide | Memory, agents, root README, workflow, audit, research plan, learning records |
| `collaboration-guide.md` | Defines teaching, implementation, explanation, and visual practices. | — | Memory, agents, plan, workflow, research plan, learning records |
| `docs/README.md` | Explains the purpose, ownership, and timing of project documentation. | Plan and all focused docs | Root README |
| `docs/stage-workflow.md` | Defines the repeatable review-build-understand-sign-off loop. | Plan, collaboration guide, challenge audit, stage journal | Root README, docs index |
| `docs/challenge-audit.md` | Holds the structured evaluation of the handed-over challenge. | Plan, resource index, source orientation | Root README, docs index, stage workflow |
| `docs/source-orientation.md` | Explains the supplied challenge at a high level before stage evaluation. | Resource index | Root README, docs index, challenge audit, tokenizer foundations |
| `docs/tokenizer-foundations.md` | Introduces tokenizer context, concepts, families, project boundaries, and a curated reading path. | Source orientation, research plan, source list | Root README, docs index |
| `docs/research-plan.md` | Defines what foundational tokenizer research must answer and how sources are assessed. | Plan, collaboration guide, resource index | Root README, docs index, tokenizer foundations |
| `docs/version-control.md` | Defines branch roles and CodeCrafters-style stage commits. | Plan | Root README, docs index |
| `learning/mental-model.md` | Provides the learner-owned canvas for drawings and evolving explanations. | Plan, collaboration guide | Root README |
| `learning/stage-journal.md` | Provides the learner-owned record for predictions, traces, explanations, and sign-off. | Plan, collaboration guide | Stage workflow |
| `resources/README.md` | Inventories and protects handed-over and collected source material. | Source list | Challenge audit, research plan |
| `resources/sources.md` | Tracks curated tokenizer learning sources, sequence, purpose, and review state. | — | Resource index, tokenizer foundations |
| `scripts/check_metadata.py` | Validates metadata presence, paths, and reciprocal edges. | Memory | Agents |

## Current state

- Phase: 2 — tokenizer foundations, ready to begin interactively.
- Completed: self-contained workspace foundation; high-level orientation to the supplied challenge; initial tokenizer mental model and curated learning sequence.
- Not completed: learner-authored mental model, full foundations discussion, stage-by-stage audit, stage redesign, language selection, implementation, or engineering tests.
- Protected input: `resources/Build your own Tokenizer.pdf` remains unchanged.
- Git policy: `main` holds the project and challenge-design record; `stages` is the cumulative, CodeCrafters-style implementation history.
- Next discussion: work through tokenizer foundations and update the learner's own mental model before auditing stage boundaries.

## Durable decisions

### 2026-10-21 — Separate engineering completion from learning sign-off

A passing implementation is necessary but does not prove understanding. Each stage tracks objective engineering evidence separately from the learner's ability to trace, explain, question, and extend the behavior.

### 2026-10-21 — Treat the handed-over challenge as evidence, not authority

The provided document is the starting proposal. Stage boundaries may be kept, split, merged, reordered, or deferred only after a recorded audit based on conceptual load, implementation load, observability, prerequisites, and testability.

### 2026-10-21 — Delay stack and implementation choices

No programming language, package manager, test framework, or code layout is selected during setup. Those choices follow initial source review and tokenizer foundations so the learner can participate in the decision.

### 2026-10-21 — Preserve learner-owned thinking artifacts

The learner's drawings, predictions, traces, questions, and explanations are first-class project artifacts. AI may help organize or challenge them but should not silently replace them with polished retrospective prose.

### 2026-10-21 — Keep a dedicated cumulative stage-history branch

`main` records the evolving learning project and approved challenge design. `stages` begins from the shared foundation and receives one cumulative commit only after each stage's engineering behavior is verified. Stage commits use the CodeCrafters-style subject `codecrafters submit [skip ci]` and an annotated stage tag provides the otherwise-missing stage identity.
