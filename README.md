<!--
@file README.md
@description Introduces the tokenizer learning project, its present boundary, navigation, and immediate next step.
@depends_on MEMORY.md, plan.md, docs/README.md, docs/stage-workflow.md, docs/challenge-audit.md, docs/source-orientation.md, docs/tokenizer-foundations.md, docs/research-plan.md, docs/harness-architecture.md, docs/version-control.md, challenge-cli/README.md, learning/mental-model.md
@depended_on_by none
-->

# Build and understand a tokenizer

This workspace has two linked goals: build a real tokenizer with substantial AI assistance, and turn that implementation journey into a rough CodeCrafters-style challenge whose stages remain small, observable, testable, and teachable.

The learner does not need to type every line. The non-negotiable requirement is conscious ownership: they should understand what each stage adds, trace the important algorithms, inspect the code that AI writes, interpret the evidence, and explain the relevant trade-offs and failure modes.

## Current boundary

The repository has completed the workspace, source-orientation, foundations, and stack/harness decisions. Go is the selected language, and the local black-box challenge harness is implemented. Stages 1–3 are accepted and checkpointed. Redesigned Stage 4 exposes `train --text TEXT --merges 1`, learns one merge rule, and emits both the model fields and post-merge IDs; it is engineering-verified and learner-accepted.

The next step is revised Stage 5, which retains the source's Stage 7 tie-breaking extension and orders candidates by frequency descending, left ID ascending, then right ID ascending.

## Where to start

- [plan.md](plan.md) defines the outcomes, phases, and explicit non-goals.
- [docs/README.md](docs/README.md) explains what each document is for, who contributes to it, and when it is used.
- [docs/source-orientation.md](docs/source-orientation.md) gives the high-level reading of the supplied challenge.
- [docs/tokenizer-foundations.md](docs/tokenizer-foundations.md) begins the plain-language tokenizer lesson and gives a curated learning sequence.
- [docs/stage-workflow.md](docs/stage-workflow.md) defines how every later challenge stage moves from prediction to engineering evidence to learning sign-off.
- [docs/challenge-audit.md](docs/challenge-audit.md) is the future stage-by-stage review surface.
- [docs/research-plan.md](docs/research-plan.md) defines the questions the tokenizer foundations phase must answer.
- [docs/harness-architecture.md](docs/harness-architecture.md) defines the Go starter, isolated challenge CLI, manifest contract, and verification boundary.
- [docs/version-control.md](docs/version-control.md) defines the `main` and `stages` branch roles and stage-commit convention.
- [challenge-cli/README.md](challenge-cli/README.md) documents the local tester commands and limitations.
- [learning/mental-model.md](learning/mental-model.md) is the learner's place for drawings and evolving models.
- [MEMORY.md](MEMORY.md) is the compact recovery point for another AI model.

Before completing a file-changing task, run:

```bash
python3 scripts/check_metadata.py
```

The starter and local challenge interface are:

```bash
./your_program.sh
./challenge list
./challenge help
```

