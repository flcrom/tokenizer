<!--
@file plan.md
@description Defines the tokenizer project's learning outcomes, phased path, scope boundaries, and completion standards.
@depends_on collaboration-guide.md
@depended_on_by AGENTS.md, MEMORY.md, README.md, docs/README.md, docs/stage-workflow.md, docs/challenge-audit.md, docs/research-plan.md, docs/harness-architecture.md, docs/version-control.md, learning/mental-model.md, learning/stage-journal.md
-->

# Tokenizer project plan

## Outcomes

This project should produce both a working tokenizer and an evidence-backed draft challenge. It should also develop the learner's ability to direct AI, inspect AI-generated implementation, reason from examples, test claims, and communicate the finished system honestly.

At the end, the learner should be able to explain where tokenizers sit in an AI system, distinguish major tokenizer families, describe the chosen algorithm and model format, trace training and inference on small inputs, navigate the implementation, interpret the test strategy, and identify what the educational implementation omits relative to production tokenizers.

## Working principles

- Begin with a mental model and concrete examples before implementation details.
- Treat the handed-over challenge as a proposal to test, not a specification to copy.
- Prefer stages with one main concept and one observable capability.
- Let AI implement within an understood, approved contract.
- Keep tests and code-inspection prompts close to the stage they explain.
- Track engineering evidence and learner understanding separately.
- Use small deterministic examples before realistic corpora or performance work.
- Add technologies only when the project creates a concrete need for them.

## Phases

### Phase 0 — Workspace foundation

Create self-contained instructions, context recovery, source inventory, audit rubric, research questions, and learner-owned note templates. Do not select the implementation stack or write tokenizer code.

### Phase 1 — High-level source orientation

Read the handed-over challenge for its intended product, algorithm family, workflow, stage sequence, testing model, and explicit exclusions. Capture questions and uncertainties without conducting the full stage audit yet.

### Phase 2 — Tokenizer foundations

Use high-quality primary and pedagogical sources to understand why tokenization exists, where it sits in training and inference, bytes versus characters, vocabulary and token IDs, major algorithm families, special tokens, normalization and pre-tokenization, reversibility, and production trade-offs. Update the learner's mental model during this phase.

### Phase 3 — Challenge audit and redesign

Evaluate every proposed stage against the audit rubric. Keep, split, merge, reorder, or defer stages with a written rationale. Define the revised stage contract and cumulative black-box acceptance tests before implementation begins.

### Phase 4 — Stack and harness decision

Completed on 2026-10-22. Go is the implementation language. The learner application follows the CodeCrafters-style `app/main.go` and `your_program.sh` interface. A separate Go module under `challenge-cli/` authors declarative stage manifests and runs black-box process tests; it does not import the application or manage Git.

### Phase 5 — Incremental build

Implement the approved stages using the loop in `docs/stage-workflow.md`. Each stage preserves earlier behavior, produces objective evidence, and remains pending for learning sign-off until the learner can explain it.

### Phase 6 — Integrated validation

Run the complete black-box suite in fresh processes, verify byte-exact round trips and deterministic training, test invalid inputs, inspect performance on a modest corpus, and document known limitations without overstating production readiness.

### Phase 7 — Challenge retrospective

Compare the built journey with the original proposal. Record which stages changed, where learners are likely to struggle, which tests best localize mistakes, what AI contributed, and what a production or CodeCrafters-ready version would still require.

## Current non-goals

Until later phases explicitly add them, this project does not promise compatibility with an existing LLM tokenizer, production-scale training, distributed processing, GPU acceleration, a web interface, model training, or a polished public CodeCrafters course definition.

## Completion standard

The project is complete only when the approved implementation is verified, the revised challenge is documented, integrated behavior is demonstrated, important limitations are explicit, and the learner can explain the mechanisms recorded in the learning journal. A green test suite alone is not completion.

