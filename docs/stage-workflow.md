<!--
@file docs/stage-workflow.md
@description Defines the repeatable workflow for auditing, implementing, verifying, and learning from each challenge stage.
@depends_on plan.md, collaboration-guide.md, docs/challenge-audit.md, learning/stage-journal.md
@depended_on_by README.md, docs/README.md
-->

# Stage workflow

Every stage moves through the same gates. A stage may be re-scoped at the first two gates; it should not quietly grow during implementation.

## 1. Frame the stage

Write the user-visible capability, the single main concept, prerequisites, explicit non-goals, exact inputs and outputs, invariants, likely failure cases, and cumulative compatibility requirements. Check the proposed size against `docs/challenge-audit.md`.

## 2. Predict and approve

Trace at least one tiny example before code exists. The learner identifies the expected result and, when useful, an intermediate state. Resolve misunderstandings and approve the stage contract before implementation.

## 3. Let AI implement

AI may create the code and tests for the approved contract. It should keep domain logic separable from CLI plumbing where that improves inspectability, preserve earlier behavior, and avoid pulling later-stage concepts forward without a recorded reason.

## 4. Inspect the implementation

Read the important execution path in order: input boundary, representation, core transformation, output boundary, and failure path. Record unfamiliar syntax or library behavior separately from algorithmic ideas so the learner can prioritize what matters.

## 5. Verify engineering behavior

Run the smallest relevant tests plus cumulative regression checks. Capture the command and observed result. Byte-level claims require byte-level comparison; persistence claims require fresh processes; determinism claims require repetition.

Engineering states are: `not_started`, `in_progress`, `verified`, or `blocked`.

## 6. Demonstrate understanding

The learner explains the core mechanism, traces a fresh example, and contributes or evaluates one adversarial case. AI identifies remaining gaps directly and answers follow-up questions.

Learning states are: `not_started`, `reviewing`, `explained`, or `signed_off`.

## 7. Reconcile and continue

Update the stage journal, decisions, metadata edges, file index, and known limitations. Run the metadata validator. Only then select the next stage.

## Definition of stage complete

A stage is fully complete only when engineering status is `verified` and learning status is `signed_off`. Either state may advance without pretending the other has advanced.

