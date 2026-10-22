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
| `app/main.go` | Dispatches base and hierarchical learned-rule encoding, model loading, pair diagnostics, deterministic training, and persistence. | — | Program launcher, application tests |
| `app/main_test.go` | Unit-tests base and hierarchical learned-rule encoding, model loading, deterministic training, stopping, and arbitrary corpus bytes. | Application entry point | — |
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
| `examples/stage-12/model.json` | Provides a canonical one-rule model for manually exercising source Stage 12. | — | — |
| `examples/stage-13/priority.json` | Provides a two-rule model demonstrating saved-rank priority. | — | — |
| `examples/stage-13/reverse-priority.json` | Reverses the same two rules to demonstrate that rank changes encoding. | — | — |
| `examples/stage-14/model.json` | Provides the three-rule banana model for hierarchical inference examples. | — | — |
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
| `scripts/check_metadata.py` | Validates metadata presence, paths, and reciprocal edges while excluding local corpus/model runtime directories. | Memory | Agents |

## Current state

- Phase: 5 in progress; revised Stages 1–9 are accepted and checkpointed, and revised Stage 10 complete learned-model encoding is engineering-verified and learner-accepted.
- Completed: accepted and checkpointed Stages 1–9; accepted zero-, one-, priority-, and hierarchical-rule inference; ready-to-run example models; focused unit tests; and 61 cumulative black-box executions.
- Not completed: curated checkpoint for Stage 10, model-based decoding, or full model validation.
- Protected input: `resources/Build your own Tokenizer.pdf` remains unchanged.
- Git policy: `main` is the living workspace; `stages` is the curated history of the shared scaffold followed by accepted cumulative stage checkpoints.
- Next discussion: checkpoint revised Stage 10, then implement source Stage 15 model-based decoding as revised Stage 11.

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

### 2026-10-22 — Keep and accept Stage 1 as the first executable slice

Stage 1 remains independent even though byte-to-ID conversion is only a few lines. Its learning boundary includes establishing the first complete command-line path: argument parsing, byte representation, integer token IDs, JSON output, launcher execution, and black-box verification. Engineering passed seven cases, and the learner explained and accepted the stage.

### 2026-10-22 — Merge original UTF-8 encoding and base decoding stages

Original Stage 2 is conceptually important but provides no new implementation in Go because the accepted Stage 1 already converts the string to its underlying UTF-8 bytes. Official documentation shows similarly direct byte paths in Rust and Ruby, while Python, JavaScript, and Java require only a small standard-library representation change. Original Stages 2 and 3 are therefore merged into revised Stage 2: encode arbitrary valid UTF-8 as base byte IDs and decode base IDs back to exact bytes.

### 2026-10-22 — Accept revised Stage 2 byte round trips

Revised Stage 2 passed 16 cumulative black-box cases plus unit, race, vet, shell, raw-byte, and metadata checks. The learner explained the raw `stdout.Write` boundary and evaluated `[0,10,255]` as an adversarial case that distinguishes exact byte preservation from Unicode-text formatting. The merged stage was accepted before moving to pair counting.

### 2026-10-22 — Keep and accept pair counting as revised Stage 3

Original Stage 4 remains independent as revised Stage 3 because it introduces an observable BPE primitive: overlapping adjacent windows, ordered compound map keys, frequency accumulation, and deterministic numeric reporting. The learner distinguished the diagnostic's left-then-right output order from later frequency-based winner selection and accepted the implementation after 22 cumulative cases passed.

### 2026-10-22 — Use the final-shaped training interface for revised Stage 4

Original Stage 5 is removed as a public challenge stage: caller-directed replacement is an implementation helper, not a useful tokenizer interface. Original Stage 6 becomes revised Stage 4. Its learner-facing command is renamed from the temporary `merge-top` to `train --text TEXT --merges 1`. The public JSON uses the canonical version, base vocabulary size, and ordered merge list, while focused tests retain the post-replacement IDs internally. The merge count is a user-supplied training-iteration budget that this stage restricts to one. Original Stage 7 is retained next to extend winner selection to tied frequencies.

### 2026-10-22 — Order tied pair candidates deterministically in revised Stage 5

Original Stage 7 becomes revised Stage 5. Winner selection sorts candidate pairs by frequency descending, then left ID ascending, then right ID ascending, and chooses the first candidate. This directly expresses the course rule, removes Go map iteration order from the result, and can be reused during every iteration when repeated training is introduced. The stage still requires the winning frequency to be at least two and still permits only one requested merge.

### 2026-10-22 — Recount after every merge in revised Stage 6

Original Stage 8 becomes revised Stage 6. For each positive requested merge, training counts the current IDs, applies the deterministic selector, assigns ID `256 + merge index`, replaces all non-overlapping matches, appends the rule, and repeats using the transformed IDs. The pair ordering from the previous iteration is never reused because every replacement changes adjacency. This stage uses inputs with enough repeated pairs to satisfy the requested count; safe early stopping, zero merges, and short inputs remain original Stage 9's boundary.

### 2026-10-22 — Keep safe stopping as revised Stage 7

Original Stage 9 is retained as revised Stage 7. The existing loop already stops at its merge budget, and selection already detects when no pair occurs twice; the stage changes that detection from an error into successful early completion. It also permits an explicitly supplied zero budget and naturally returns empty merge lists for empty or one-byte input. Missing or negative merge budgets still fail. Final transformed IDs remain internal inspection data rather than persisted model fields.

### 2026-10-22 — Keep persisted model JSON canonical

Training stdout and future `model.json` files contain exactly the document's required fields: `version`, `base_vocab_size`, and ordered `merges`. The corpus-specific transformed IDs are useful for teaching and tests but are not tokenizer state, so `trainingResult` retains them internally while `tokenizerModel` alone is serialized. Previously published stage tags remain immutable historical checkpoints; current cumulative manifests for Stages 4 onward enforce the canonical contract.

### 2026-10-22 — Keep raw-file training and persistence as revised Stage 8

Original Stage 10 becomes revised Stage 8. `train` accepts exactly one corpus source: explicit `--text` or `--input PATH`. File input is read as raw bytes without trimming, line splitting, or UTF-8 interpretation and feeds the existing trainer through `trainBytes`. The canonical model is serialized once with a newline; when `--output PATH` is supplied, those same bytes are written to the file before stdout, preventing the saved and printed models from drifting. Model JSON parsing remains the next loading stage rather than part of persistence.

### 2026-10-22 — Isolate empty-model loading as revised Stage 9

Original Stage 11 becomes revised Stage 9. `encode --model PATH --text TEXT` reads and JSON-decodes a canonical model in a fresh process. Empty models preserve base UTF-8 byte IDs, which isolates persistence and parsing from the separate learned-rule application algorithm in original Stage 12. The implementation rejects nonempty merge lists rather than silently ignoring them. Missing files and malformed JSON fail without stdout; complete semantic validation remains deferred.

### 2026-10-22 — Merge complete learned-model encoding as revised Stage 10

Original Stages 12–14 become revised Stage 10. Its cases explicitly distinguish zero, one, multiple, and hierarchical rules. Encoding starts from UTF-8 byte IDs and applies rules in saved order, assigning rule `i` learned ID `base_vocab_size + i`; every transformed array feeds the next rule, so later rules can consume earlier learned IDs. Training already created such models and the inference loop already had the required mechanics, making original Stage 14 only an artificial guard plus tests rather than an independent implementation stage.
