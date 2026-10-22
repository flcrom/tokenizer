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
| `README.md` | Provides the human entry point, current phase, commands, and navigation. | Memory, plan, focused docs, challenge CLI guide, mental-model canvas | — |
| `plan.md` | Defines outcomes, boundaries, phases, and completion criteria. | Collaboration guide | Memory, agents, root README, focused docs, learning records |
| `collaboration-guide.md` | Defines teaching, implementation, explanation, and visual practices. | — | Memory, agents, plan, workflow, research plan, learning records |
| `go.mod` | Declares the standalone learner-application Go module. | — | Program launcher |
| `app/main.go` | Provides the compilable pre-stage tokenizer entry point. | — | Program launcher |
| `your_program.sh` | Compiles and executes the learner program through a stable process boundary. | Application module and entry point | Runner, CLI guide, harness architecture |
| `challenge` | Runs the isolated challenge CLI from the repository root. | Challenge module and entry point | — |
| `challenge-cli/go.mod` | Declares the isolated local challenge Go module. | — | Challenge wrapper, entry point, CLI guide |
| `challenge-cli/cmd/challenge/main.go` | Starts the challenge CLI and maps failures to exit codes. | Challenge module, CLI application | Challenge wrapper, CLI guide |
| `challenge-cli/internal/cli/cli.go` | Implements stage creation, listing, test selection, execution, and reports. | Stage loader and types, runner | CLI entry point, CLI tests, manifest guide |
| `challenge-cli/internal/cli/cli_test.go` | Verifies stage scaffolding, listing, and CLI argument failures. | CLI application | — |
| `challenge-cli/internal/spec/types.go` | Defines and validates stage manifests, fixtures, and expectations. | — | Loader, runner, runner tests, manifest guide |
| `challenge-cli/internal/spec/loader.go` | Strictly loads, validates, sorts, and selects stage manifests. | Stage types | CLI application, loader tests |
| `challenge-cli/internal/spec/loader_test.go` | Verifies strict parsing, ordering, and duplicate rejection. | Stage loader | — |
| `challenge-cli/internal/runner/runner.go` | Executes isolated black-box cases and compares observable results. | Program launcher, stage types | CLI application, runner tests |
| `challenge-cli/internal/runner/runner_test.go` | Verifies process, JSON, bytes, files, exits, and timeout behavior. | Runner, stage types, fake program | — |
| `challenge-cli/testdata/programs/fake_program.sh` | Supplies deterministic behaviors for runner self-tests. | — | Runner tests |
| `challenge-cli/stages/README.md` | Documents generated stage-manifest structure and supported cases. | Stage types, CLI application | CLI guide |
| `challenge-cli/README.md` | Explains local challenge commands, boundaries, and verification. | Challenge module, entry point, manifest guide, launcher | Root README, harness architecture |
| `docs/README.md` | Explains the purpose, ownership, and timing of project documentation. | Plan and all focused docs | Root README |
| `docs/stage-workflow.md` | Defines the repeatable review-build-understand-accept loop. | Plan, collaboration guide, challenge audit, stage journal | Root README, docs index, harness architecture |
| `docs/challenge-audit.md` | Holds the structured evaluation of the handed-over challenge. | Plan, resource index, source orientation | Root README, docs index, stage workflow |
| `docs/source-orientation.md` | Explains the supplied challenge at a high level before stage evaluation. | Resource index | Root README, docs index, challenge audit, tokenizer foundations |
| `docs/tokenizer-foundations.md` | Introduces tokenizer context, concepts, families, project boundaries, and a curated reading path. | Source orientation, research plan, source list | Root README, docs index |
| `docs/research-plan.md` | Defines what foundational tokenizer research must answer and how sources are assessed. | Plan, collaboration guide, resource index | Root README, docs index, tokenizer foundations |
| `docs/harness-architecture.md` | Defines the Go starter, isolated runner, manifests, evidence, and checkpoint lifecycle. | Plan, stage workflow, version control, CLI guide, launcher | Root README, docs index |
| `docs/version-control.md` | Defines branch roles and accepted-stage checkpoint commits. | Plan | Root README, docs index, harness architecture |
| `learning/mental-model.md` | Provides the learner-owned canvas for drawings and evolving explanations. | Plan, collaboration guide | Root README |
| `learning/stage-journal.md` | Provides the learner-owned record for predictions, traces, explanations, and sign-off. | Plan, collaboration guide | Stage workflow |
| `resources/README.md` | Inventories and protects handed-over and collected source material. | Source list | Challenge audit, research plan |
| `resources/sources.md` | Tracks curated tokenizer learning sources, sequence, purpose, and review state. | — | Resource index, tokenizer foundations |
| `scripts/check_metadata.py` | Validates metadata presence, paths, and reciprocal edges. | Memory | Agents |

## Current state

- Phase: 4 complete; ready to review the harness and frame Stage 1.
- Completed: foundations; Go selection; CodeCrafters-style starter; separate local challenge module; stage creation/listing/testing commands; black-box runner self-tests.
- Not completed: learner review of the harness, any tokenizer stage contract, tokenizer behavior, stage acceptance tests, or stage checkpoint commit.
- Protected input: `resources/Build your own Tokenizer.pdf` remains unchanged.
- Git policy: `main` is the living workspace; `stages` is the curated history of the shared scaffold followed by accepted cumulative stage checkpoints.
- Next discussion: learner code review of the harness, followed by the Stage 1 contract and its initially failing black-box cases.

## Durable decisions

### 2026-10-21 — Separate engineering completion from learning sign-off

A passing implementation is necessary but does not prove understanding. Each stage tracks objective engineering evidence separately from the learner's ability to trace, explain, question, and extend the behavior.

### 2026-10-21 — Treat the handed-over challenge as evidence, not authority

The provided document is the starting proposal. Stage boundaries may be kept, split, merged, reordered, or deferred only after a recorded audit based on conceptual load, implementation load, observability, prerequisites, and testability.

### 2026-10-21 — Delay stack and implementation choices

Superseded on 2026-10-22 after the learner completed foundations and approved the Go harness plan. During initial setup, no programming language, package manager, test framework, or code layout was selected.

### 2026-10-21 — Preserve learner-owned thinking artifacts

The learner's drawings, predictions, traces, questions, and explanations are first-class project artifacts. AI may help organize or challenge them but should not silently replace them with polished retrospective prose.

### 2026-10-21 — Keep a dedicated cumulative stage-history branch

Updated on 2026-10-22: `main` is the living workspace, while `stages` begins from the shared scaffold and receives one cumulative checkpoint only after engineering verification, learning sign-off, and explicit learner acceptance. Stage commits use the CodeCrafters-style subject `codecrafters submit [skip ci]` and an annotated stage tag provides the otherwise-missing stage identity.

### 2026-10-22 — Go application and isolated local challenge harness

The learner application is a root Go module with `app/main.go` and the stable `your_program.sh` interface. The local challenge CLI is a separate Go module that creates strict JSON stage manifests and runs the launcher only as a black-box child process in isolated temporary workspaces. It observes arguments, stdin, stdout, stderr, exit status, time, and declared files; it never imports learner code or manages Git.
