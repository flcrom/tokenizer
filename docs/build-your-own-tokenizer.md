<!--
@file docs/build-your-own-tokenizer.md
@description Defines the final standalone 12-stage byte-level BPE tokenizer challenge built and verified in this repository.
@depends_on plan.md, docs/challenge-audit.md, docs/harness-architecture.md, docs/version-control.md, resources/README.md, resources/sources.md
@depended_on_by README.md, docs/README.md
-->

# Build your own tokenizer

## A 12-stage byte-level BPE challenge

Build a small tokenizer that learns merge rules from a corpus, saves them, and uses them to turn new text into token IDs and back into the original bytes.

The challenge is cumulative. Each stage adds one observable behavior to the same command-line program. The examples are deliberately small enough to trace by hand, while the final program contains the complete train-encode-decode core of a byte-level Byte Pair Encoding tokenizer.

This is an AI-assisted learning challenge. AI may write substantial amounts of code, but the learner should predict small examples, inspect the important implementation paths, understand the tests, and be able to explain why the resulting behavior is correct.

## What you will build

The finished program has three everyday operations:

1. `train` learns an ordered list of byte-pair merge rules from a corpus;
2. `encode` applies a saved model to new text and emits integer token IDs;
3. `decode` expands token IDs and writes the original bytes.

The course also keeps `pairs` as a diagnostic command. It exposes the frequency-counting primitive before that primitive becomes part of training.

Every possible byte already has a base token:

- byte value `0` has token ID `0`;
- byte value `1` has token ID `1`;
- ...
- byte value `255` has token ID `255`.

Learned tokens start at ID `256`. This means that unfamiliar words, Unicode text, and even arbitrary byte sequences remain representable without an unknown-byte token.

The model and its token IDs belong only to this challenge. They are not automatically compatible with GPT, Claude, Llama, or another existing model.

## The finished workflow

Train three merges from a small corpus:

```console
$ ./your_program.sh train --text 'banana banana' --merges 3
{"version":1,"base_vocab_size":256,"merges":[[97,110],[98,256],[256,97]]}
```

Save the same model:

```console
$ ./your_program.sh train --text 'banana banana' --merges 3 --output model.json
{"version":1,"base_vocab_size":256,"merges":[[97,110],[98,256],[256,97]]}
```

The same canonical JSON is written to `model.json` and stdout.

Encode new text in a fresh process:

```console
$ ./your_program.sh encode --model model.json --text 'banana!'
[257,258,33]
```

Decode those IDs:

```console
$ ./your_program.sh decode --model model.json --ids '[257,258,33]'
banana!
```

In this model:

- token `256` expands to `an`;
- token `257` expands to `ban`;
- token `258` expands to `ana`.

The model is therefore both learned data and a recipe for reconstructing every learned token.

## Program and command contract

Implement the tokenizer in any programming language. The only required execution boundary is:

```text
./your_program.sh <command> [arguments]
```

`your_program.sh` is the language adapter. It may compile or invoke the learner's implementation, but callers do not need to know which language is behind it. A starter repository for a particular language only needs to provide this launcher and an initial program entry point.

The local challenge runner is independent of the learner implementation. It does not import application code and does not manage Git. It invokes `./your_program.sh` as a black-box child process in a fresh temporary workspace and observes only:

- arguments and declared input files;
- stdout and stderr;
- exit status and timeout;
- declared output files.

Useful commands are:

```console
$ ./challenge list
$ ./challenge test --stage 4
$ ./challenge test --stage 4 --previous
$ ./challenge test --stage 4 --case 'trains one alternating-pair rule' --debug
```

Run the current stage with `--previous` before considering it complete. Every stage must preserve all earlier behavior.

## Shared behavior across the challenge

### Bytes and text

Convert `--text` to its UTF-8 bytes without normalization, trimming, case conversion, word splitting, or inserted separators. The complete input is one byte sequence, including spaces and newlines.

A visible character is not necessarily one byte. For example:

```text
é  -> [195,169]
🙂 -> [240,159,153,130]
```

Precomposed `é` and decomposed `e` followed by U+0301 are different byte sequences and must remain different.

### Pair counting

Count every adjacent window, including overlapping windows. The three IDs `[97,97,97]` contain two pair windows:

```text
positions 0-1 -> [97,97]
positions 1-2 -> [97,97]
```

Its pair frequency is therefore two.

### Pair replacement

Replacement scans from left to right. When the chosen pair matches, append the new token ID and advance by two. Otherwise copy the current ID and advance by one.

Replacement does not overlap. Replacing `[97,97]` in `[97,97,97]` therefore produces `[256,97]`, not two copies of token `256`.

### Training

For each requested merge:

1. count adjacent pairs in the current IDs;
2. select the highest-frequency pair;
3. break a frequency tie by smaller left ID, then smaller right ID;
4. stop successfully if the winning count is below two;
5. assign ID `256 + merge_index`;
6. replace every non-overlapping occurrence;
7. append the pair to the model and recount the transformed sequence.

`--merges N` is a maximum number of new rules, not a target vocabulary size. Training may stop early.

### Encoding with a model

Start with UTF-8 byte IDs. Apply each saved merge rule in model order, replacing every non-overlapping match in the current sequence. The output from one rule becomes the input to the next rule.

Do not recount frequencies or learn new pairs while encoding. Training decided the priority; inference follows the saved order. A learned rule may apply even if it occurs only once in the new input.

### Decoding

Each base token expands to its one byte. Each learned token expands to the concatenation of its two operands' already-known byte expansions. Concatenate all requested token expansions and write the resulting bytes exactly once, without adding a newline.

Do not convert each individual token to UTF-8 text. A learned token may contain only part of a multibyte character even though the complete byte stream is valid UTF-8.

### Output

Successful `encode`, `pairs`, and `train` commands print one JSON value followed by a newline. JSON object whitespace and key order are not significant; array order is significant.

Successful `decode` writes raw bytes with no added newline.

After the final stage, every failure exits with status `1`, leaves stdout empty, and writes one JSON error object to stderr.

## The model format

A saved model is a JSON object with exactly three required fields:

```json
{
  "version": 1,
  "base_vocab_size": 256,
  "merges": [
    [97, 110],
    [98, 256],
    [256, 97]
  ]
}
```

The rules are ordered. Merge entry `i` creates token ID:

```text
base_vocab_size + i
```

With a base vocabulary size of `256`, entries zero, one, and two create IDs `256`, `257`, and `258`.

Each merge must contain two nonnegative integer IDs. Both operands must already exist when the rule is defined, so each must be smaller than the rule's new ID. A rule may refer to base tokens or earlier learned tokens, but never itself or a future token. The same pair may not be defined twice.

Unknown extra object fields may be ignored. Transformed corpus IDs, pair frequencies, and textual token spellings are not model state and are not stored.

## The 12 cumulative stages

Complete these stages in order. Inputs may be treated as valid until Stage 12 introduces the complete error contract.

### Stage 1: Encode ASCII text

**New capability:** Add `encode --text` for ASCII input. With no model, emit one integer ID for every byte in order.

**Main concept:** A token ID can name a byte.

```console
$ ./your_program.sh encode --text 'Hi!'
[72,105,33]
```

#### Requirements

- Preserve the byte order.
- Preserve spaces, tabs, newlines, repeated characters, and punctuation.
- Empty text produces `[]`.
- Do not add IDs that were not present in the input.

#### Objective checks

- Encode one byte, several bytes, spaces, all ASCII whitespace, repeated characters, punctuation, and empty text.
- Compare the parsed JSON array rather than visual formatting.

#### Implementation guidance

Use the language's byte or UTF-8 encoding facilities and convert each resulting byte value to an integer before JSON encoding the array. Do not iterate characters or Unicode code points and assume their numeric values are bytes. Some languages expose a string's encoded bytes directly; others require an explicit UTF-8 encoding operation.

This is intentionally the first complete executable slice: argument parsing, conversion, JSON output, the launcher, and black-box testing all become real even though the core conversion is small.

### Stage 2: Encode and decode UTF-8 base tokens

**New capability:** Extend encoding to arbitrary valid UTF-8 and add `decode --ids` for base token IDs.

**Main concept:** Characters and bytes are different units, and decoding reconstructs bytes rather than character objects.

```console
$ ./your_program.sh encode --text 'café'
[99,97,102,195,169]

$ ./your_program.sh decode --ids '[99,97,102,195,169]'
café
```

#### Requirements

- Keep every encoded ID between `0` and `255`.
- Do not normalize text.
- Parse `--ids` as a JSON array.
- Convert base IDs to bytes and write them directly to stdout.
- Add no newline during decoding.

#### Objective checks

- Encode accented text, emoji, and mixed ASCII/Unicode text.
- Distinguish precomposed `é` (`[195,169]`) from decomposed `e` plus U+0301 (`[101,204,129]`).
- Decode an empty array to zero bytes.
- Decode `[0,10,255]` and compare exact bytes `00 0A FF` rather than terminal rendering.
- Preserve an embedded newline exactly.

#### Why the original stages were merged

The original proposal separated UTF-8 encoding and base decoding. In many languages, extending an ASCII byte path to UTF-8 is a small standard-library representation change; in byte-oriented string models it may require no algorithmic change. Combining that extension with decoding creates a meaningful, language-independent byte-level round trip while preserving both concepts in focused tests.

### Stage 3: Count adjacent pairs

**New capability:** Add the diagnostic command `pairs --ids`. Emit `[left,right,count]` records for every distinct adjacent pair.

**Main concept:** Pair frequency is counted over overlapping adjacent windows.

```console
$ ./your_program.sh pairs --ids '[97,98,97,98]'
[[97,98,2],[98,97,1]]
```

#### Requirements

- Count windows at positions `(0,1)`, `(1,2)`, and so on.
- Treat a pair as an ordered compound key.
- Sort output numerically by left ID and then right ID, not by frequency.
- Return `[]` for zero or one input ID.

#### Objective checks

- `[97,97,97]` produces `[[97,97,2]]`.
- `[1,2,1,2]` distinguishes the two ordered pairs.
- Tied counts remain ordered numerically.
- Learned IDs such as `256` are counted exactly like base IDs.

This diagnostic order is for inspectable output. The training winner introduced next uses frequency first.

### Stage 4: Train the first merge

**New capability:** Add `train --text TEXT --merges 1`. Select the most frequent pair, assign ID `256`, replace it internally, and emit a one-rule model.

**Main concept:** Training turns repeated adjacency into a reusable rule.

```console
$ ./your_program.sh train --text 'ababab' --merges 1
{"version":1,"base_vocab_size":256,"merges":[[97,98]]}
```

The initial IDs are:

```text
[97,98,97,98,97,98]
```

The pair counts include:

```text
[97,98] -> 3
[98,97] -> 2
```

Replacing the winner internally produces:

```text
[256,256,256]
```

#### Requirements

- This stage accepts a user-supplied merge budget of exactly one.
- Select a fixture with one unique winner whose count is at least two.
- Store the winning pair, not its frequency, in `merges`.
- Keep transformed IDs inspectable in focused tests but out of the saved model.

#### Objective checks

- Learn an alternating pair.
- Learn a repeated same-ID pair.
- Learn the same pair across separated occurrences.
- Confirm replacement consumes matches without overlap.

The original proposal exposed caller-directed pair replacement as its own command. This challenge keeps replacement as a focused helper rather than an artificial public stage.

### Stage 5: Resolve tied pair frequencies

**New capability:** Make winner selection deterministic when several pairs share the highest frequency.

**Main concept:** Reproducible training needs an explicit tie rule.

```console
$ ./your_program.sh train --text 'babab' --merges 1
{"version":1,"base_vocab_size":256,"merges":[[97,98]]}
```

Both `[98,97]` and `[97,98]` occur twice. `[97,98]` wins because its left ID is smaller, even though `[98,97]` appeared first.

#### Requirements

Order candidate pairs by:

1. frequency descending;
2. left token ID ascending;
3. right token ID ascending.

Choose the first candidate after applying that ordering. Do not depend on hash-map or dictionary iteration order, or on whichever pair appeared first in the text.

#### Objective checks

- Prove left-ID priority with tied counts.
- Prove right-ID priority when the left IDs match.
- Retain a unique-winner regression case.
- Repeated fresh-process runs must choose the same model.

### Stage 6: Learn successive merges

**New capability:** Permit positive `--merges N`, recount after every replacement, and allow later rules to reference learned IDs.

**Main concept:** Every merge changes the next frequency table.

```console
$ ./your_program.sh train --text 'banana banana' --merges 3
{"version":1,"base_vocab_size":256,"merges":[[97,110],[98,256],[256,97]]}
```

The transformed sequences are:

```text
initial: [98,97,110,97,110,97,32,98,97,110,97,110,97]
rule 0:  [98,256,256,97,32,98,256,256,97]
rule 1:  [257,256,97,32,257,256,97]
rule 2:  [257,258,32,257,258]
```

#### Requirements

- Run at most `N` iterations.
- Recount pairs from the current IDs on every iteration.
- Assign learned IDs consecutively from `256`.
- Append each selected pair to the model in learning order.
- Allow a new rule to contain an earlier learned ID.

#### Objective checks

- Compare one-, two-, and three-merge runs on the same corpus.
- Ensure stale initial frequencies produce the wrong answer.
- Verify the learned IDs are exactly `256`, `257`, and `258`.

### Stage 7: Stop training safely

**New capability:** Treat the merge count as a budget and finish successfully when no pair repeats.

**Main concept:** Normal algorithm termination is not an error.

```console
$ ./your_program.sh train --text 'abc' --merges 10
{"version":1,"base_vocab_size":256,"merges":[]}
```

#### Requirements

- Stop after `N` rules or when the best frequency is below two.
- Permit an explicitly supplied zero budget.
- Empty and one-byte input produce an empty merge list.
- Never invent a rule merely to reach the requested count.

#### Objective checks

- `--merges 0` performs no work.
- Empty text and one-byte text remain valid inputs.
- `aaaaa` with a large budget learns only `[97,97]`.
- `banana banana` with budget `99` stops after four useful rules:

```json
[[97,110],[98,256],[256,97],[257,258]]
```

### Stage 8: Train from a file and save the model

**New capability:** Add raw file input and optional model persistence.

**Main concept:** A learned model is portable data that outlives the training process.

```console
$ ./your_program.sh train --input corpus.txt --merges 3 --output model.json
{"version":1,"base_vocab_size":256,"merges":[[97,110],[98,256],[256,97]]}
```

#### Requirements

- Accept exactly one corpus source: `--text TEXT` or `--input PATH`.
- Read file contents as raw bytes.
- Do not trim a trailing newline or interpret the file as lines.
- Serialize the canonical model once.
- When `--output` is supplied, write exactly the same bytes to the file and stdout.

#### Objective checks

- Equivalent text and file bytes learn the same model.
- The requested model file exists after the process exits.
- Parsed file and stdout models are equal.
- Preserve newline bytes and arbitrary non-UTF-8 corpus bytes.

### Stage 9: Load an empty model

**New capability:** Add `--model PATH` to `encode` and reconstruct the model in a fresh process. Begin with an empty merge list.

**Main concept:** Serialization is useful only when another process can load it.

Given `empty.json`:

```json
{"version":1,"base_vocab_size":256,"merges":[]}
```

the command remains equivalent to base-byte encoding:

```console
$ ./your_program.sh encode --model empty.json --text 'Hi!'
[72,105,33]
```

#### Requirements

- Read and JSON-decode the model from disk.
- Do not depend on in-memory state from a previous training process.
- Accept ordinary JSON whitespace and reordered object keys.
- Preserve no-model encoding behavior.
- At this stage, a missing path or malformed JSON must fail without stdout; the final error codes arrive in Stage 12.

#### Objective checks

- Load an empty model in a fresh process.
- Compare it with no-model encoding.
- Prove the model path is actually read by testing missing and malformed files.

### Stage 10: Encode with learned merge rules

**New capability:** Apply zero, one, multiple, prioritized, and hierarchical learned rules to new text.

**Main concept:** Encoding replays the ordered transformations learned during training.

Using the banana model:

```console
$ ./your_program.sh encode --model model.json --text 'banana!'
[257,258,33]
```

The trace is:

```text
bytes:              [98,97,110,97,110,97,33]
apply [97,110]:     [98,256,256,97,33]
apply [98,256]:     [257,256,97,33]
apply [256,97]:     [257,258,33]
```

#### Requirements

- Empty models leave byte IDs unchanged.
- Apply each rule in saved order to the current ID sequence.
- Replace all non-overlapping occurrences of that rule.
- A rule may apply only once in the new input and must still be used.
- Do not rank pairs using frequencies in the new input.
- Later rules may consume learned IDs created by earlier rules.
- Encoding must not mutate the model or add rules.

#### Priority example

For a model whose merges are `[[97,98],[98,99]]`, encoding `abc bc bc` applies `ab` before the more frequent `bc`:

```text
[256,99,32,257,32,257]
```

Reversing the stored rule order changes the result:

```text
[97,256,32,256,32,256]
```

#### Objective checks

- Cover zero rules, one rule, several base-ID rules, and several hierarchical rules.
- Verify non-overlapping replacement on repeated input.
- Reverse a two-rule model and require a different encoding.
- Encode `banana!`, `banana banana`, and inputs where one learned token enables another.
- Repeat a priority-sensitive case in fresh processes to prove determinism.

#### Why three original stages were merged

Once encoding applies one rule to the current IDs, multiple rules require only iteration. Refreshing the ID sequence naturally enables hierarchical tokens, and training already emits such models. Keeping zero-, one-, priority-, and hierarchy-focused tests preserves the concepts without manufacturing separate implementation stages.

### Stage 11: Decode learned tokens

**New capability:** Add `--model PATH` to `decode` and expand learned token IDs into exact bytes.

**Main concept:** A learned token's meaning is a recursively composed byte sequence.

```console
$ ./your_program.sh decode --model model.json --ids '[257,258,33]'
banana!
```

#### Requirements

- Build base expansions first: ID `n` expands to byte `n` for `0 <= n < 256`.
- Build learned expansions in merge order.
- For each merge, concatenate the left expansion and right expansion.
- Permit base and learned IDs in the same decode request.
- Concatenate requested expansions before writing bytes.
- Add no newline and perform no per-token UTF-8 conversion.

#### Important byte-boundary example

With model rule `[[169,195]]`, learned token `256` expands to bytes `[169,195]`. That token alone is not valid UTF-8. It can still participate in a complete byte stream such as:

```text
[195] + [169,195] + [169]
   C3       A9 C3       A9
```

Together those bytes represent `éé`. Treating each token as independent text would corrupt the round trip.

#### Objective checks

- Decode hierarchical banana tokens.
- Mix base and learned IDs.
- Decode models trained on multibyte text.
- Preserve unseen multilingual base bytes.
- Verify raw byte equality when a learned-token expansion is not independently valid UTF-8.

### Stage 12: Reject invalid inputs cleanly

**New capability:** Establish one stable failure contract across every command and validate model graphs before trusting them.

**Main concept:** Command input, filesystem operations, model definitions, and vocabulary membership are different failure domains.

```console
$ ./your_program.sh decode --model model.json --ids '[999]'
```

Observable result:

```text
stdout: empty
stderr: {"error":"unknown_token"}
exit status: 1
```

Every failure uses exactly one of four codes.

#### `invalid_input`

Use this when the invocation itself is malformed:

- unknown commands or flags;
- missing required flags;
- malformed, null, non-integer, or negative `--ids` values;
- non-integer or negative merge budgets;
- conflicting or missing corpus sources.

#### `io_error`

Use this when an operating-system read or write fails:

- unreadable model path;
- unreadable corpus path;
- unwritable output path;
- failed stdout write.

#### `invalid_model`

Use this when a readable model does not satisfy the contract:

- malformed model JSON;
- missing or wrongly typed required fields;
- unsupported version or base vocabulary size;
- a merge that does not contain exactly two integers;
- negative operands;
- self-references or future references;
- duplicate pair rules.

Unknown extra object fields remain allowed.

#### `unknown_token`

Use this when decode receives a nonnegative ID that is not defined by the selected vocabulary.

With three model rules, the valid IDs are `0` through `258`. ID `999` is therefore unknown. Without a model, every ID above `255` is unknown.

A negative ID is `invalid_input`, not `unknown_token`, because it is malformed command data rather than a possible vocabulary ID.

#### Model-graph invariant

Rule `i` creates ID `256 + i`. Both operands must be smaller than that new ID:

- rule `0` may use IDs `0..255`;
- rule `1` may also use ID `256`;
- rule `0` may not use ID `256`, because that is a self-reference;
- no rule may refer to an ID created in the future.

#### Objective checks

- Cover all four public error codes.
- Require status `1`, empty stdout, and structural JSON on stderr.
- Reject malformed command values before producing success output.
- Reject every invalid-model family above.
- Confirm unknown extra fields do not break a valid model.
- Keep every valid case from Stages 1-11 passing.

## Complete worked example

Consider the corpus `banana banana` with a budget of three merges.

### 1. Convert the corpus to bytes

```text
b  a  n  a  n  a     b  a  n  a  n  a
98 97 110 97 110 97 32 98 97 110 97 110 97
```

### 2. Learn rule zero

The pair `[97,110]` (`an`) has the highest frequency and becomes token `256`:

```text
[98,256,256,97,32,98,256,256,97]
```

### 3. Recount and learn rule one

The current pair `[98,256]` becomes token `257`:

```text
[257,256,97,32,257,256,97]
```

Token `257` expands to bytes for `ban` because it combines byte `98` (`b`) with token `256` (`an`).

### 4. Recount and learn rule two

The current pair `[256,97]` becomes token `258`:

```text
[257,258,32,257,258]
```

Token `258` expands to `ana`.

### 5. Save the model

```json
{
  "version": 1,
  "base_vocab_size": 256,
  "merges": [[97,110],[98,256],[256,97]]
}
```

### 6. Encode unseen text

Encoding `banana!` starts from bytes, applies those three saved rules in order, and produces:

```text
[257,258,33]
```

No frequencies are counted during inference.

### 7. Decode the IDs

The decoder builds these expansions:

```text
256 -> bytes(97) + bytes(110) -> an
257 -> bytes(98) + token(256)  -> ban
258 -> token(256) + bytes(97)  -> ana
```

Concatenating expansions for `[257,258,33]` produces the exact bytes for `banana!`.

## How the challenge is tested

The local suite contains 90 declared black-box cases and 92 process executions; one priority-sensitive case runs three times. The final cumulative command is:

```console
$ ./challenge test --stage 12 --previous
```

The suite:

- invokes `./your_program.sh` in an isolated temporary directory;
- parses successful JSON structurally;
- compares decode output as exact bytes;
- checks stderr and exit status;
- creates temporary corpus and model fixtures;
- checks persisted output files;
- repeats determinism-sensitive cases;
- keeps every earlier stage active.

Unit tests remain useful for pure helpers such as pair selection, replacement, model validation, and learned-token expansion. Black-box tests define the learner-facing contract.

Before considering the complete implementation verified, run the cumulative black-box suite and the unit/static checks appropriate to the chosen language:

```console
$ ./challenge test --stage 12 --previous
$ <your language's unit-test command>
$ <your language's formatter, linter, or static-analysis command>
```

Passing tests prove the covered behavior. They do not by themselves prove that the learner understands the algorithm.

## Learning with AI

AI can scaffold commands, implement helpers, diagnose failures, and propose test cases. The learner retains ownership of requirements and understanding.

For each stage:

1. read the contract and predict one small result;
2. ask AI to implement within that boundary;
3. run the new stage and all previous stages;
4. inspect the important execution path;
5. explain why the result is correct;
6. consider one adversarial or boundary example;
7. accept the stage only when both engineering and understanding are satisfactory.

Three questions anchor the core mechanisms:

1. Why do three identical bytes contain two adjacent pair windows but only one non-overlapping replacement?
2. Why can a pair that is rare in new text merge before a more common pair during encoding?
3. Why must learned-token byte expansions be concatenated before interpreting the result as UTF-8?

## Scope and deliberate limitations

This challenge includes:

- a byte vocabulary covering IDs `0..255`;
- deterministic byte-level BPE training;
- a compact ordered pair-rule model;
- fresh-process persistence;
- model-based encoding and decoding;
- raw-byte round trips;
- structured validation errors.

It deliberately excludes:

- Unicode normalization;
- regex or word-level pre-tokenization;
- special tokens and chat templates;
- padding and truncation;
- compatibility with an existing LLM's vocabulary;
- streaming input;
- an optimized heap-based or incremental trainer;
- production-scale memory and throughput guarantees;
- neural-network training, GPUs, or model inference.

The whole corpus is currently held and repeatedly rescanned in memory. This makes the algorithm easy to inspect but is not the approach expected from a production tokenizer trainer.

## From the supplied 16 stages to these 12 stages

The handed-over document was treated as a proposal to test rather than an immutable specification. Implementation evidence produced this final mapping:

| Supplied stage | Final stage | Decision |
| --- | --- | --- |
| 1: Encode ASCII | 1 | Kept as the first end-to-end executable slice. |
| 2: Encode UTF-8 | 2 | Merged with base decoding because UTF-8 conversion is usually a small representation step and the combined stage creates a meaningful byte round trip in any language. |
| 3: Decode base IDs | 2 | Merged to create a substantive byte round-trip stage. |
| 4: Count pairs | 3 | Kept as an independently observable BPE primitive. |
| 5: Replace one pair | - | Removed as a public stage; retained as an internal, unit-tested helper. |
| 6: Train first merge | 4 | Revised to use the final-shaped `train --merges 1` interface. |
| 7: Break ties | 5 | Kept as a focused deterministic-selection extension. |
| 8: Successive merges | 6 | Kept; requires recounting after every transformed sequence. |
| 9: Stop safely | 7 | Kept as the training termination boundary. |
| 10: Save model | 8 | Kept with raw file input and canonical persistence. |
| 11: Load model | 9 | Kept with an empty model to isolate loading from rule application. |
| 12: Apply one pair | 10 | Merged with priority and hierarchy while retaining separate test groups. |
| 13: Respect priority | 10 | Merged because the ordered rule loop is already required. |
| 14: Hierarchical tokens | 10 | Merged because refreshed IDs naturally allow later learned-token rules. |
| 15: Decode learned tokens | 11 | Kept as a meaningful inverse-model operation. |
| 16: Reject invalid inputs | 12 | Kept whole after implementation and measured-diff review. |

The final sequence removes one artificial public command and avoids stages whose behavior was already implemented by an earlier natural design, while preserving each important concept in examples and tests.

## Suggested learning path

The following is a flexible sequence, not a measured completion forecast:

- Sessions 1-2: Stages 1-3, byte round trips and pair counting.
- Sessions 3-4: Stages 4-7, training, ties, repeated merging, and stopping.
- Session 5: Stages 8-9, persistence across processes.
- Sessions 6-7: Stages 10-11, complete model inference and learned-token decoding.
- Sessions 8-10 as needed: Stage 12, cumulative regression work, integrated validation, and explaining the complete algorithm.

Faster learners can combine sessions. A difficult stage can take an entire session. The useful boundary is understanding and observable progress, not a fixed calendar duration.

## Further resources

- [Hugging Face: Tokenization algorithms](https://huggingface.co/docs/transformers/en/tokenizer_summary) gives a compact comparison of BPE, WordPiece, Unigram, and related approaches.
- [Andrej Karpathy: Let's build the GPT Tokenizer](https://www.youtube.com/watch?v=zduSFxRajkE) connects Unicode, UTF-8, BPE, special tokens, and production tokenizer behavior.
- [Hugging Face LLM Course: Byte-Pair Encoding tokenization](https://huggingface.co/learn/llm-course/en/chapter6/5) provides a worked BPE learning example.
- [Karpathy's minbpe](https://github.com/karpathy/minbpe) is a compact educational byte-level BPE implementation.
- [OpenAI's tiktoken educational implementation](https://github.com/openai/tiktoken/blob/main/tiktoken/_educational.py) is useful for comparing this explicit model with a production-adjacent implementation.
- [Tiktokenizer](https://tiktokenizer.vercel.app/) interactively shows model-specific tokens, IDs, whitespace, and chat-message framing.

These resources help connect the small inspectable challenge to real tokenizer systems. Their exact algorithms and file formats are not automatically part of this course contract.
