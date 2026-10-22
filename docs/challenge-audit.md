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
| 4 — Count adjacent pairs | `pairs --ids` emits sorted `[left,right,count]` records for every adjacent window | Pair frequency is counted over overlapping adjacent windows | Moderate: sliding windows, compound map keys, frequencies, and deterministic sorting | `keep` as revised Stage 3 | This is an independently observable BPE primitive with several localized failure modes: skipping overlaps, counting values instead of pairs, relying on map order, or sorting by frequency. It is large enough to implement and explain without also introducing replacement semantics. |
| 5 — Replace one pair | `merge --ids --pair --new-id` replaces a caller-supplied pair | Replacement consumes non-overlapping matches | Small and isolated, but the public command is artificial and makes no BPE decision | `remove` as a stage | Do not expose the source's caller-directed `merge` command. Replacement remains a pure, directly unit-tested helper used by training, but does not need its own learner-facing stage. |
| 6 — Train the first merge | Select a unique most-frequent pair, assign ID 256, replace it, and emit a one-rule result | Frequency learns a reusable rule and immediately transforms the training IDs | Several linked operations, but Stage 3 already supplies counting | `keep and revise` as revised Stage 4 | Rename the temporary `merge-top` command to `train --text --merges 1`. Emit the canonical model fields; keep selection and replacement IDs observable through focused internal tests. |
| 7 — Resolve tied frequencies | Extend winner selection when multiple pairs share the highest frequency | Training needs a deterministic tie rule | Focused extension to selection semantics | `keep` as revised Stage 5 | Stage 5 sorts by frequency descending, left ID ascending, then right ID ascending. Focused cases separately prove left-ID priority and right-ID priority while retaining a unique-winner regression. |
| 8 — Learn successive merges | Permit positive `--merges N`, recount after every replacement, and let later rules reference learned IDs | Each merge changes the next frequency table | Substantial but cohesive loop over the established primitives | `keep` as revised Stage 6 | Reuse deterministic selection and replacement in a loop. Compare one-, two-, and three-merge runs on the same text so stale frequency tables and incorrect learned IDs are observable. |
| 9 — Stop training safely | Treat `N` as an upper bound; stop when the best count is below 2; support zero, empty, and one-byte inputs | Normal termination is distinct from failure | Small implementation delta but essential algorithmic contract | `keep` as revised Stage 7 | The loop and minimum-frequency detection already exist, but Stage 7 turns the latter into successful completion and verifies every boundary from the source document. |
| 10 — Save a trained model | Add raw `--input PATH` as an alternative corpus source and optional `--output PATH` persistence | A learned model is portable data; corpus bytes must remain unchanged | Moderate I/O boundary around stable training logic | `keep` as revised Stage 8 | Read the corpus with `os.ReadFile`, train through the byte-oriented core, serialize the canonical model once, and write the same JSON bytes to the file and stdout. |
| 11 — Load a model | Add `encode --model PATH`; parse a saved empty model in a fresh process without changing base-byte output | Persistence becomes useful only when another operation can reconstruct the model | Focused file-read and JSON-boundary stage | `keep` as revised Stage 9 | An empty merge list isolates loading from learned-rule application. Missing and malformed files fail cleanly; whitespace and JSON key order do not matter. Nonempty-rule encoding remains source Stage 12. |
| 12 — Apply a learned pair | Encode new text with a model containing one learned pair | Inference reuses a training rule even when it occurs only once in new input | Focused reuse of established replacement semantics | `merge` with original Stages 13–14 as revised Stage 10 | Empty-model loading already establishes the inference boundary and replacement is already implemented. Keep explicit zero- and one-rule cases, then make the stage substantive with priority and hierarchical composition. |
| 13 — Respect learned merge priority | Apply multiple base-ID rules by saved rank rather than new-input frequency | Training decides priority; inference follows it | Focused loop extension with an adversarial ordering case | `merge` with original Stages 12 and 14 as revised Stage 10 | The combined stage contrasts zero, one, and multiple rules. Reversing two rules must change output. |
| 14 — Encode hierarchical tokens | Permit later rules to contain IDs created by earlier rules | Learned tokens compose into larger learned tokens | Conceptually important but already supported by the established ordered replacement loop | `merge` into revised Stage 10 | Training already creates hierarchical rules and inference already refreshes IDs after each saved rule. Removing the artificial base-ID-only guard plus focused multi-level cases completes encoding without a new algorithm. |

## Revised Stage 4 contract

`train --text TEXT --merges 1` converts text to byte IDs, counts adjacent pairs, chooses the uniquely most-frequent pair, assigns it ID 256, replaces its occurrences left-to-right without overlap, and prints `{"version":1,"base_vocab_size":256,"merges":[[left,right]]}`.

This removes the source's caller-directed Stage 5 entirely and implements its Stage 6 as revised Stage 4. The public result follows the document's canonical model format; the post-replacement IDs remain inspectable in `trainingResult` and focused unit tests rather than becoming model data. `base_vocab_size` is the fixed count of byte tokens 0–255; learned merge index 0 therefore receives ID 256. `--merges` is supplied by the user and represents the maximum number of training iterations, but Stage 4 accepts only the value `1`. Stage 4 fixtures have one unique winner whose overlapping count is at least two; the source's Stage 7 tie-frequency extension is implemented separately as revised Stage 5.

## Revised Stage 5 contract

The command and output shape remain `train --text TEXT --merges 1`. When several candidates share the highest frequency, selection orders all candidates by frequency descending, left token ID ascending, and right token ID ascending, then chooses the first pair. The winning pair must still occur at least twice. This makes selection independent of Go map iteration order and gives later repeated training a reusable deterministic selector. More than one merge and safe early stopping remain out of scope.

## Revised Stage 6 contract

`train --text TEXT --merges N` accepts a positive merge count and attempts exactly that many iterations on fixtures with enough repeated pairs. Each iteration counts the current IDs, selects the deterministic winner, assigns ID `256 + merge index`, performs non-overlapping replacement, appends the ordered rule, and then recounts the transformed sequence. Stdout remains the canonical model; unit tests inspect the transformed IDs separately. Zero merges, empty or one-byte input, and stopping before the requested budget are deliberately deferred to the source's Stage 9.

## Revised Stage 7 contract

`--merges N` is now an upper bound and may be zero. Training stops successfully after `N` rules or as soon as there is no pair with count at least 2. Empty input, one-byte input, and input with only unique adjacent pairs return a canonical model with an empty merge array; their unchanged IDs remain internal. A large budget may therefore produce fewer rules than requested. Missing or negative budgets remain invalid, while the final structured error contract remains deferred.

## Revised Stage 8 contract

`train` requires exactly one corpus source: `--text TEXT` or `--input PATH`. `--input` reads every file byte unchanged, including newlines and bytes that are not valid UTF-8. Training otherwise reuses the existing merge budget and stopping behavior. Stdout always prints the canonical model JSON; optional `--output PATH` writes the same serialized bytes to disk. JSON parsing is not needed yet because this stage creates a model—loading and validating model JSON begins with the source's Stage 11.

## Revised Stage 9 contract

`encode --model PATH --text TEXT` reads and JSON-decodes the canonical model in a fresh process. For this stage, the supplied model must have an empty `merges` array, so encoding still returns the UTF-8 byte IDs unchanged. JSON whitespace and object-key order are irrelevant. An unreadable path or malformed JSON exits nonzero without stdout. A nonempty model is deliberately rejected because applying learned rules is source Stage 12; full semantic validation of model fields remains deferred to the source's final validation stage.

## Revised Stage 10 contract

`encode --model PATH --text TEXT` covers models with zero, one, multiple, or hierarchical learned rules. It starts from UTF-8 byte IDs and visits saved rules in array order. Rule index `i` creates ID `base_vocab_size + i`; every non-overlapping occurrence is replaced before the next rule sees the transformed array. Later rules may therefore consume IDs produced by earlier rules. A saved rule applies even if it occurs once, and new-input frequency never changes priority. Fixtures separately verify zero rules, one rule, ordered and reversed priority, non-overlap, model immutability, and two levels of learned-token composition.

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

