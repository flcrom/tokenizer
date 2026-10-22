<!--
@file docs/challenge-audit.md
@description Provides the rubric and workspace for evaluating and revising the handed-over tokenizer challenge.
@depends_on plan.md, resources/README.md, docs/source-orientation.md
@depended_on_by README.md, docs/README.md, docs/stage-workflow.md
-->

# Challenge audit

## Status

`in_progress`

The source document presents a 16-stage byte-level BPE challenge. The foundations and high-level orientation are complete, so stages are now reviewed immediately before implementation. The table below records decisions cumulatively rather than approving all 16 boundaries in advance.

## Who answers the questions?

These are not prerequisite questions for the learner and not a quiz. AI first answers source-fact questions from the document and cited references. We then discuss the answers in plain language. The learner is asked for a decision only when the answer becomes a course-design choice, such as whether a stage is too large or whether a diagnostic command is worth keeping.

## Audit rubric

For each proposed stage, record:

- the new observable behavior and main concept;
- conceptual prerequisites and implementation prerequisites;
- algorithmic difficulty and ordinary plumbing separately;
- the number of new decisions, edge cases, and failure modes;
- whether an intermediate result is inspectable;
- whether black-box tests can localize likely mistakes;
- what earlier behavior must remain active;
- what the learner should be able to predict, trace, and explain;
- a decision: `keep`, `split`, `merge`, `reorder`, or `defer`;
- the rationale and any replacement stage boundary.

## Size heuristics

A stage is probably too large when it introduces several independent concepts, touches many boundaries at once, needs a long explanation before any behavior is visible, or can fail in ways that the stage tests cannot localize.

A stage is probably too small when it adds only mechanical wiring, has no meaningful independent observation, repeats the same reasoning as an adjacent stage, or cannot support a useful learner prediction or test on its own.

Do not optimize for equal code size or equal duration. A good stage has a coherent learning boundary and useful feedback.

## What the source-level questions mean

| Question | Plain-language meaning | Why we ask it | Who answers, and when |
| --- | --- | --- | --- |
| What finished artifact and workflow does the document propose? | What program exists at the end, and what can a user actually do with it? | It gives us the destination before we inspect the route. | AI drafts from the document in Phase 1; learner only needs to confirm the description makes sense. |
| Which parts are core tokenizer behavior and which are teaching-only diagnostics? | Which commands are needed by the final product, and which exist mainly to reveal an intermediate step? | A useful learning command may not belong to a real tokenizer interface. | AI identifies them in Phase 1; we decide together whether the teaching value justifies them in Phase 3. |
| Which assumptions are course choices rather than properties of BPE generally? | What did the author choose to simplify or make deterministic, even though other BPE implementations can behave differently? | It prevents us from mistaking one educational design for the universal definition of BPE. | AI gathers candidates during Phases 1–2; we validate them from references before the audit. |
| What does the model format preserve, and what does it omit? | What information is saved so another process can encode and decode, and what production metadata is absent? | It tests whether the saved artifact is sufficient for the promised workflow. | AI explains after vocabulary and merge rules are introduced in Phase 2. |
| How does the test strategy establish correctness? | Which observable checks prove byte preservation, deterministic training, persistence, and clean failures? | A challenge needs tests that reveal the first broken concept, not just a wrong final answer. | AI maps the proposed checks in Phase 1; detailed evaluation happens together in Phase 3. |
| Which explicit exclusions matter compared with production tokenizers? | What real-world features are intentionally missing, and do those omissions preserve or distort the learning goal? | It keeps claims honest and identifies possible later extensions. | AI explains the terms in Phase 2; we judge the boundary together in Phase 3. |

## Stage review table

Populate this only after the high-level source orientation.

| Original stage | Capability | Main concept | Load | Decision | Rationale / replacement |
| --- | --- | --- | --- | --- | --- |
| 1 — Encode ASCII text | `encode --text` emits one integer ID per ASCII byte | A token ID can name a byte | Small core conversion plus first end-to-end CLI integration | `keep` — accepted | The conversion is only a few lines, but the stage establishes the command-line input, flag parsing, transformation, JSON output, launcher, and black-box contract together. The learner found that sufficient for a first stage. ASCII keeps the byte/character distinction out of scope until Stage 2. |

