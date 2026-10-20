<!--
@file collaboration-guide.md
@description Defines how AI should teach, implement, verify, review understanding, and create visuals in this project.
@depends_on none
@depended_on_by AGENTS.md, MEMORY.md, plan.md, docs/stage-workflow.md, docs/research-plan.md, learning/mental-model.md, learning/stage-journal.md
-->

# Collaboration guide

## Teaching

Begin unfamiliar ideas with a small scenario or byte-level example, then introduce precise vocabulary. Connect every abstraction to the tokenizer being built: its inputs, internal representation, files, commands, tests, or failure behavior.

Do not conceal the path from question to decision. Explain why an option fits, what it costs, and what would make another option preferable. Use short connected paragraphs by default and structured lists when they clarify comparisons or procedures.

## AI-assisted implementation

AI is expected to write much of the code, but not to take silent ownership of requirements or design. Before substantial code changes, establish the stage contract, prerequisites, examples, invariants, failure cases, and verification plan. Keep each implementation small enough to inspect in one learning session.

Before running new behavior, ask the learner for a prediction on one tiny input when practical. After implementation, point to the important code paths and explain them in execution order. Do not flood the learner with every syntax detail; identify which details are essential to correctness and which are ordinary language or library mechanics.

## Understanding and sign-off

Use three kinds of checks without turning every interaction into a quiz:

- tracing: predict intermediate sequences, chosen rules, or output bytes;
- explanation: describe why the behavior is correct and why a tempting alternative is wrong;
- transfer: propose an adversarial or unseen example and reason about it.

Record engineering status and learning status independently. If tests pass but explanation is incomplete, say exactly that. If the learner understands a proposal that has not been implemented, do not call the behavior working.

## Evidence

Match evidence to the claim. Use unit tests for pure helpers, black-box process tests for challenge contracts, raw-byte comparisons for decoding, repeated fresh-process runs for determinism and persistence, and measurements for performance claims. Report commands and observed results that actually ran.

## Drawings and visual artifacts

The learner may think through rough drawings, ASCII sketches, Markdown, or richer diagrams. Preserve rough learner-authored artifacts and date meaningful revisions. AI may label, question, or create a separate clean interpretation, but should not overwrite the original thinking.

For polished diagrams, default to a light theme, black/white/grey foundations, restrained semantic color, small corner radii, strong alignment, and visible direction of data or control flow. State important omissions whenever a diagram simplifies the real system.

