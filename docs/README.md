<!--
@file docs/README.md
@description Explains the purpose, ownership, and timing of each focused project document.
@depends_on plan.md, docs/build-your-own-tokenizer.md, docs/stage-workflow.md, docs/challenge-audit.md, docs/source-orientation.md, docs/tokenizer-foundations.md, docs/research-plan.md, docs/harness-architecture.md, docs/version-control.md
@depended_on_by README.md
-->

# What the docs folder is for

The files in this folder are working agreements and reasoning records. They do not ask the learner to arrive with tokenizer knowledge, and they are not tests whose answers must be supplied all at once.

| Document | What it defines | Who contributes | When it is used |
| --- | --- | --- | --- |
| `build-your-own-tokenizer.md` | The standalone final 12-stage challenge: shared contracts, stage instructions, examples, tests, scope, and source-to-final mapping. | AI drafts from verified implementation evidence; learner reviews clarity, size, and teaching value. | Phase 7 and whenever the challenge is handed to another learner or reviewer. |
| `source-orientation.md` | A high-level reading of what the supplied challenge proposes. | AI drafts from the source; learner asks questions and corrects interpretations. | Phase 1, before detailed tokenizer study or stage verdicts. |
| `tokenizer-foundations.md` | The initial mental model, vocabulary, project context, and learning sequence. | AI teaches and sources; learner builds and revises their own model. | Phase 2, before the stage audit. |
| `research-plan.md` | The questions our foundations work must eventually answer and the standard for sources. | AI curates; learner decides where more depth is needed. | Phase 2 and whenever a knowledge gap appears. |
| `challenge-audit.md` | The rubric and decisions for keeping, splitting, merging, reordering, or deferring proposed stages. | AI gathers evidence and proposes; learner discusses and approves important decisions. | Phase 3, after foundations. |
| `stage-workflow.md` | How one approved stage moves from prediction through code, tests, explanation, and sign-off. | Shared operating agreement. | Phase 5, for every implementation stage. |
| `harness-architecture.md` | The Go starter, separate challenge CLI, manifest format, process boundary, and verification responsibility. | AI implements after learner approval; learner reviews the code and boundary. | Phase 4 and when the tester evolves. |
| `version-control.md` | What belongs on `main` versus `stages`, and how passed stages are committed. | Shared operating agreement. | Throughout the project, especially Phase 5. |

The questions in these documents have different owners. Source-fact questions are answered by reading the supplied material and authoritative references; AI can do the first pass. Design-choice questions are discussed together because they affect what the challenge teaches. Learning-check questions are answered by the learner only after the relevant idea has been taught and demonstrated.

