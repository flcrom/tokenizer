<!--
@file README.md
@description Introduces the tokenizer learning project, its present boundary, navigation, and immediate next step.
@depends_on MEMORY.md, plan.md, docs/README.md, docs/stage-workflow.md, docs/challenge-audit.md, docs/source-orientation.md, docs/tokenizer-foundations.md, docs/research-plan.md, docs/version-control.md, learning/mental-model.md
@depended_on_by none
-->

# Build and understand a tokenizer

This workspace has two linked goals: build a real tokenizer with substantial AI assistance, and turn that implementation journey into a rough CodeCrafters-style challenge whose stages remain small, observable, testable, and teachable.

The learner does not need to type every line. The non-negotiable requirement is conscious ownership: they should understand what each stage adds, trace the important algorithms, inspect the code that AI writes, interpret the evidence, and explain the relevant trade-offs and failure modes.

## Current boundary

The repository has completed **Phase 0: workspace foundation** and the first **Phase 1: source orientation** pass. No implementation language or tokenizer code has been selected. The CEO-provided document is preserved as source material and has been described at a high level, but its individual stage boundaries have not been approved.

The next step is **Phase 2: tokenizer foundations**. We will build a usable mental model of tokenization, bytes, token IDs, vocabularies, training, encoding, decoding, and major tokenizer families before auditing all 16 stages.

## Where to start

- [plan.md](plan.md) defines the outcomes, phases, and explicit non-goals.
- [docs/README.md](docs/README.md) explains what each document is for, who contributes to it, and when it is used.
- [docs/source-orientation.md](docs/source-orientation.md) gives the high-level reading of the supplied challenge.
- [docs/tokenizer-foundations.md](docs/tokenizer-foundations.md) begins the plain-language tokenizer lesson and gives a curated learning sequence.
- [docs/stage-workflow.md](docs/stage-workflow.md) defines how every later challenge stage moves from prediction to engineering evidence to learning sign-off.
- [docs/challenge-audit.md](docs/challenge-audit.md) is the future stage-by-stage review surface.
- [docs/research-plan.md](docs/research-plan.md) defines the questions the tokenizer foundations phase must answer.
- [docs/version-control.md](docs/version-control.md) defines the `main` and `stages` branch roles and stage-commit convention.
- [learning/mental-model.md](learning/mental-model.md) is the learner's place for drawings and evolving models.
- [MEMORY.md](MEMORY.md) is the compact recovery point for another AI model.

Before completing a file-changing task, run:

```bash
python3 scripts/check_metadata.py
```

