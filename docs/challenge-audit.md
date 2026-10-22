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
| 2 — Encode UTF-8 bytes | Extend `encode --text` to arbitrary valid UTF-8 without normalization | A code point may occupy multiple bytes | Conceptually meaningful; often no implementation change in byte-oriented languages | `merge` with original Stage 3 | In Go, the accepted Stage 1 implementation already passes the proposed UTF-8 cases because `[]byte(text)` exposes the UTF-8 bytes. Rust and Ruby have similarly direct byte APIs. Other popular languages need a small explicit encoding change, but the behavior remains one standard-library conversion. Merge with decoding to form a substantive byte-level round-trip stage. |
| 3 — Decode base token IDs | `decode --ids` writes the exact represented bytes with no added newline | Lossless decoding reconstructs bytes | Moderate: a new command, JSON input, byte reconstruction, and a raw-byte output boundary | `merge` with original Stage 2 | Revised Stage 2 is **Encode and decode UTF-8 base tokens**. It combines the important bytes-versus-characters lesson with observable reversibility and enough new implementation work across languages. |

## Cross-language evidence for the Stage 2 decision

The important comparison is not whether a language *can* encode UTF-8; all of these can. It is whether a natural byte-oriented Stage 1 implementation already handles UTF-8, leaving no implementation work for the original Stage 2.

| Language | Plausible ASCII-only Stage 1 route | Correct UTF-8 byte route | Likely Stage 2 change |
| --- | --- | --- | --- |
| Go | Convert the string to `[]byte`, as this project does | The same `[]byte(text)` conversion | None. Go strings store bytes, and our Stage 1 already uses the correct representation. |
| Rust | Iterate `str::bytes()` or use `as_bytes()` | The same byte iterator or slice | None when Stage 1 is byte-oriented; Rust strings are already UTF-8. |
| Ruby | Call `String#bytes` | The same method when the input string is UTF-8 | Usually none; Ruby strings expose their underlying byte sequence directly. |
| Python | Iterate characters and call `ord`, which works for ASCII | Encode the `str` with UTF-8, then iterate the resulting `bytes` integers | A small but real representation change: `text` becomes `text.encode("utf-8")`. |
| JavaScript / Node.js | Read character or UTF-16 code-unit values, which happen to match ASCII | Create a UTF-8 `Buffer` or use `TextEncoder`, then read its bytes | A small but real API and representation change. Direct string iteration is not UTF-8-byte iteration. |
| Java | Read `char` values, which happen to match ASCII | Call `getBytes(StandardCharsets.UTF_8)`, then convert signed bytes to IDs `0–255` | More visible plumbing because Java `byte` is signed and must be widened as unsigned data. |

The Go result was also checked directly before this decision: `café`, `🙂`, and `A🙂é` produced 5, 4, and 7 IDs respectively without a code change. Precomposed `é` produced `[195,169]`, while decomposed `e` plus U+0301 produced `[101,204,129]`; the implementation does not normalize them.

Official references: [Go strings, bytes, and runes](https://go.dev/blog/strings), [Rust `str::bytes`](https://doc.rust-lang.org/std/primitive.str.html#method.bytes), [Ruby `String`](https://docs.ruby-lang.org/en/master/String.html), [Python text and binary sequence types](https://docs.python.org/3/library/stdtypes.html), [Node.js `Buffer.from`](https://nodejs.org/docs/latest/api/buffer.html#static-method-bufferfromstring-encoding), [Java `String.getBytes`](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/String.html#getBytes(java.nio.charset.Charset)), and the [Java primitive-type ranges](https://docs.oracle.com/javase/specs/jls/se21/html/jls-4.html#jls-4.2.1).

