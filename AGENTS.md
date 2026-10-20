<!--
@file AGENTS.md
@description Defines repository-level operating instructions for AI collaborators working on the tokenizer learning challenge.
@depends_on MEMORY.md, plan.md, collaboration-guide.md, scripts/check_metadata.py
@depended_on_by none
-->

# Repository instructions for AI collaborators

Read [MEMORY.md](MEMORY.md) before proposing or making changes. Use its file index to recover the smallest relevant context set instead of reading the whole repository.

Read [plan.md](plan.md) before changing the project sequence, tokenizer scope, implementation language, challenge stages, learning goals, or evaluation method. Read [collaboration-guide.md](collaboration-guide.md) before teaching, implementing code, creating diagrams, or reviewing the learner's explanation.

This project has two separate outcomes:

1. a tokenizer implementation that passes objective checks; and
2. a learner who can explain the important behavior, trace examples, and defend the design.

Never treat the first outcome as proof of the second. AI may write most of the code, but implementation must proceed in inspectable stages. Before a stage is implemented, ask the learner to predict at least one small example. After verification, explain the important code paths, invite an adversarial example, and leave learning sign-off pending until the learner demonstrates understanding.

The CEO-provided challenge is a proposal to evaluate, not an immutable specification. Preserve it unchanged under `resources/`. Record proposed stage changes in `docs/challenge-audit.md` before changing the build plan. Do not begin tokenizer implementation or choose a language unless the current phase and decision log authorize it.

Use the two-pass workflow in `MEMORY.md`: analyze impact before implementation, then reconcile metadata and the file index after verification. Every new human-authored file must have a valid metadata block within its first ten lines. When an edge changes, update both `@depends_on` and the corresponding `@depended_on_by` declaration.

Do not claim behavior works without appropriate evidence. Report only tests, byte comparisons, process exits, traces, or manual checks actually completed. Distinguish verified engineering behavior, learner sign-off, proposals, and open questions.

Before completing any file-changing task, run:

```bash
python3 scripts/check_metadata.py
```

