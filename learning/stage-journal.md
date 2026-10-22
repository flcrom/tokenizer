<!--
@file learning/stage-journal.md
@description Provides the learner-owned template for predictions, implementation reading, evidence, questions, and stage sign-off.
@depends_on plan.md, collaboration-guide.md
@depended_on_by docs/stage-workflow.md
-->

# Stage journal

Copy the template below for each approved stage. Keep the learner's pre-implementation prediction intact; corrections belong beneath it.

## Stage 1 — Encode ASCII text

### Contract

- Stage name: Encode ASCII text
- New observable behavior: `./your_program.sh encode --text TEXT` prints a JSON array containing one integer ID per input byte, in order.
- Main concept: A base token ID can directly name a byte.
- Explicit non-goals: non-ASCII test coverage, decoding, model loading, merge rules, and the final invalid-input contract.
- Engineering status: `verified`
- Learning status: `signed_off`

### My prediction before implementation

- Tiny inputs and expected outputs, as supplied by the learner:
  - `"Hi!"` → `[72, 105, 33]`
  - `"A A"` → `[65, 32, 65]`
  - `""` → `[]`
- Expected intermediate states: not recorded yet.
- Why: not recorded yet.

### What AI implemented

- Important files and symbols: `app/main.go` contains `main`, `run`, and `encodeBytes`; `app/main_test.go` checks the pure conversion; `challenge-cli/stages/01-encode-ascii.json` defines the black-box contract.
- Execution path in order: `main` passes process arguments to `run`; `run` recognizes `encode`, reads `--text`, calls `encodeBytes`, and sends the resulting slice to Go's JSON encoder; `encodeBytes` converts the string to bytes and copies each byte value into an integer token ID.
- Important invariants: output order and length match the input bytes; no byte is trimmed or skipped; the empty input is represented by a non-nil empty slice so JSON emits `[]`, not `null`; successful output contains exactly one JSON value plus a newline.
- Language or library mechanics I had to learn: `os.Args[1:]` excludes the executable name; Go strings can be converted to `[]byte`; converting a `byte` to `int` preserves its value; `json.Encoder.Encode` appends the required newline.

### Evidence

- Command or check: before implementation, `./challenge test --stage 1 --debug` with a project-local Go cache.
- Observed result: all 6 cases reached the starter and failed with exit code 1.
- Command or check: after implementation, `./challenge test --stage 1 --previous`.
- Observed result: the original 6 black-box cases passed; the accepted suite now contains 7 cases after adding ASCII control whitespace.
- Command or check: root `go test ./...` and `go vet ./...`; challenge module `go test -race ./...` and `go vet ./...`; `sh -n challenge your_program.sh`.
- Observed result: all checks exited 0 when run with the project-local Go build and module caches.
- What this proves: the approved ASCII inputs cross the real process boundary and produce structurally equal JSON arrays with empty stderr and exit code 0; the pure helper and harness regression tests pass; both Go modules pass static analysis.
- What this does not prove: Stage 2 Unicode requirements, decoding, model behavior, the final error contract, or learner understanding.

### My explanation after reading the code

Learner explanation, preserved as given:

> path is basically take the arguments from os, make sure that the arguments are correct using flags and then invoking the encodebytes to convert the string to byte[] to int[]

AI clarification: this correctly identifies the input boundary, argument parsing, and core representation change. The remaining step in the execution path is the output boundary: `json.NewEncoder` serializes the `[]int` and adds the required newline. At this stage the flag parsing provides basic command shape checks; the complete invalid-input contract is deliberately deferred to Stage 16.

### Adversarial case

- Input: `" A\t\nA "`
- Expected behavior: `[32,65,9,10,65,32]`
- Why this case is useful: it catches trimming, word splitting, or code that only handles visible spaces while losing tab or newline bytes. Emoji is not ASCII and belongs to the Stage 2 boundary instead.
- Observed behavior, if run: passed through the black-box harness, producing `[32,65,9,10,65,32]`.

### Questions and sign-off

- Remaining questions: none for Stage 1.
- AI feedback on gaps: the core explanation is correct; include the JSON output boundary when describing the complete execution path.
- Learner acceptance: accepted as satisfactory. The learner's reasoning, preserved as given:

  > the reason it feels satisfactory is becuase not only do i have to write the 4 line code for the encodebyte functoin but i have to establish making the application work from command line, so because of that I think this feels satisfactory enough for the first stage

- Engineering status: `verified`
- Learning status: `signed_off`
- Date signed off: 2026-10-22.

## Stage 2 — Encode and decode UTF-8 base tokens

### Contract

- Stage name: Encode and decode UTF-8 base tokens
- New observable behavior: `encode --text` is verified for arbitrary valid UTF-8 without normalization, and `decode --ids` converts JSON base-token IDs into exact raw bytes without adding a newline.
- Main concept: UTF-8 text is a byte sequence, and base token IDs round-trip those bytes without the tokenizer interpreting characters.
- Explicit non-goals: learned token IDs above 255, merge rules, UTF-8 validation, normalization, and the final invalid-input contract.
- Engineering status: `verified`
- Learning status: `signed_off`

### My prediction before implementation

- Encoding lengths predicted by the learner:
  - `"café"` → 5 IDs
  - `"🙂"` → 4 IDs
  - `"A🙂é"` → 7 IDs
- Initial decoding outputs, preserved as supplied by the learner:
  - `[99,97,102,195,169]` → `[63, 61, 66, c3, a9]`
  - `[65,10,66]` → `[A, , B]`
  - `[]` → `[]`
- Correction discussed before implementation: `63 61 66 C3 A9` is the hexadecimal representation of the output bytes, not a printed array; the terminal renders those valid bytes as `café`. Byte 10 creates a real newline between `A` and `B`. Decoding `[]` writes zero bytes rather than the characters `[]` or a blank line.
- Expected intermediate states: JSON IDs → `[]int` → `[]byte` → raw stdout; no byte-to-Unicode conversion occurs inside the tokenizer.
- Why: each base token ID names one byte directly, so lossless decoding only reconstructs and writes those bytes.

### What AI implemented

- Important files and symbols: `app/main.go` now dispatches through `runEncode` and `runDecode`; `decodeBaseIDs` performs the base-ID-to-byte conversion. `app/main_test.go` covers encoding, decoding, and pure round trips. `challenge-cli/stages/02-utf8-base-round-trip.json` defines five encoding and four byte-exact decoding cases.
- Execution path in order: `main` passes process arguments and stdout to `run`; `run` dispatches `decode`; `runDecode` parses `--ids`, `json.Unmarshal` builds a `[]int`, `decodeBaseIDs` converts each ID to a byte, and `stdout.Write` emits the resulting bytes unchanged.
- Important invariants: one base ID produces exactly one byte in the same position; decoding does not interpret or validate UTF-8; no output newline is added; an empty ID array writes zero bytes; current-stage inputs are valid base IDs from 0 through 255.
- Language or library mechanics I had to learn: separate `flag.FlagSet` values keep command options independent; `json.Unmarshal` parses JSON into an existing Go value; `byte(id)` converts a valid base ID to one byte; `io.Writer.Write` writes raw bytes rather than formatting them as text.

### Evidence

- Command or check: before decoding implementation, `./challenge test --stage 2 --debug`.
- Observed result: all 5 UTF-8 encoding cases passed through the unchanged Stage 1 conversion, while all 4 decoding cases failed because only `encode` was recognized.
- Command or check: after implementation, `./challenge test --stage 2 --previous`.
- Observed result: all 16 cumulative black-box cases passed: 7 from Stage 1 and 9 from revised Stage 2.
- Command or check: root `go test ./...` and `go vet ./...`; challenge module `go test -race ./...` and `go vet ./...`.
- Observed result: all checks exited 0 with the project-local Go caches.
- What this proves: Unicode text is encoded without normalization; base IDs reconstruct exact bytes, including `00 0A FF`; decoding adds no newline; empty input writes zero bytes; existing ASCII behavior remains active; pure UTF-8 fixtures round-trip.
- What this does not prove: learned IDs above 255, model-based decoding, the final invalid-input contract, or learner understanding.

### My explanation after reading the code

Learner explanation, preserved as given:

> stdout.Write takes byte array in input and thats why we are using it. not 100% sure but we need the actual string to be outputted thats is why.

AI clarification accepted during review: `stdout.Write` is used because the tokenizer's output contract is exact bytes, not because Go must construct a string. A terminal may render valid UTF-8 bytes as text, but invalid or invisible bytes must still be emitted unchanged.

### Adversarial case

- Input: `[0,10,255]`
- Expected behavior: write exactly the bytes `00 0A FF`.
- Why this case is useful: it contains a null byte, a newline, and a byte that is not valid UTF-8 on its own. It catches implementations that treat decoded data as Unicode text, stop at a null byte, trim it, or append a formatting newline.
- Observed behavior, if run: the black-box harness compared raw bytes and passed.

### Questions and sign-off

- Remaining questions: none for revised Stage 2.
- AI feedback on gaps: the learner identified the raw-byte output API correctly and accepted the distinction between writing bytes and constructing a string.
- Learner acceptance: accepted by asking to move to the next stage after the code and adversarial-case review.
- Engineering status: `verified`
- Learning status: `signed_off`
- Date signed off: 2026-10-22.

## Stage 3 — Count adjacent pairs

### Contract

- Stage name: Count adjacent pairs
- New observable behavior: `pairs --ids JSON` prints sorted `[leftID,rightID,count]` records for every distinct adjacent pair.
- Main concept: pair frequency is counted over overlapping adjacent windows.
- Explicit non-goals: selecting the most frequent pair, replacing pairs, training, merge IDs, and the final invalid-input contract.
- Engineering status: `verified`
- Learning status: `signed_off`

### My prediction before implementation

- Expected outputs supplied by the learner:
  - `[1,2,1,2,1]` → `[[1,2,2],[2,1,2]]`
  - `[7,7,7,7]` → `[[7,7,3]]`
  - `[42]` → `[]`
- Proposed algorithm, preserved as supplied by the learner:

  > the way i was thinking one should implement this is have 2 pointers, initialize one at the beginning and the other at the second place. make a key value map before this iteration where the key is the pair and the value is the frequency, increament the key for each pair and iterate the first pointer till the second last element
  >
  > then for each key value pair, extract the key,value,pair as left right count and sort it and return the answer

- Expected intermediate states: build a pair-frequency map from `(ids[i], ids[i+1])`, then turn the map entries into records and sort by left ID followed by right ID.
- Why: advancing both adjacent positions by one includes every window, including overlapping pairs.

### What AI implemented

- Important files and symbols: `app/main.go` adds `runPairs`, the shared `parseIDs` helper, the comparable `tokenPair` map key, and `countAdjacentPairs`. `app/main_test.go` unit-tests overlaps, ordering, and short inputs. `challenge-cli/stages/03-count-adjacent-pairs.json` defines six black-box cases.
- Execution path in order: `run` dispatches `pairs`; `runPairs` parses `--ids`; `parseIDs` produces `[]int`; `countAdjacentPairs` scans `(ids[index], ids[index+1])`, increments a map, collects distinct keys, sorts them by left then right ID, and builds `[left,right,count]` records; the JSON encoder writes the records.
- Important invariants: `n` IDs have `n-1` adjacent windows when `n` is at least one; advancing by one preserves overlaps; pair identity includes both ordered IDs; output ordering never depends on Go map iteration; zero pairs produce JSON `[]`, not `null`.
- Language or library mechanics I had to learn: a Go struct containing integers is comparable and can be a map key; `counts[pair]++` supplies zero for an unseen key before incrementing; ranging over a map is unordered; `sort.Slice` applies the explicit left-then-right comparator.

### Evidence

- Command or check: before implementation, `./challenge test --stage 3 --debug`.
- Observed result: all 6 cases failed because `pairs` was not recognized.
- Command or check: after implementation, `./challenge test --stage 3 --previous`.
- Observed result: all 22 cumulative black-box cases passed: 7 from Stage 1, 9 from Stage 2, and 6 from Stage 3.
- Command or check: root `go test ./...` and `go vet ./...`; challenge module `go test -race ./...` and `go vet ./...`.
- Observed result: all checks exited 0 with the project-local Go caches.
- What this proves: adjacent pairs are counted with overlaps, arbitrary token IDs are supported, records are sorted by numeric pair rather than count or map order, short inputs produce `[]`, and earlier behavior remains active.
- What this does not prove: choosing a winning pair, replacing occurrences, training, or learner understanding.

### My explanation after reading the code

The learner's implementation model was recorded before implementation: scan every adjacent `(ids[index], ids[index+1])` pair, increment its map frequency, turn map entries into records, and sort the result.

During code review, the learner asked:

> why exactly are we sorting based on left first and then right after wards and not by frequency ? because when we have to do the merges, we will have to sort this by frequency right and take the most frequent pair and then merge it ?

AI clarification accepted during review: the `pairs` diagnostic uses numeric pair order to provide stable, canonical output. A later training stage performs a separate winner-selection operation using highest frequency first and numeric pair order only as its tie-break. Selecting a winner does not require changing this diagnostic output contract or sorting the entire table by frequency.

### Adversarial case

- Input: `[256,97,256,97]`
- Expected behavior: `[[97,256,1],[256,97,2]]`.
- Why this case is useful: numeric pair order intentionally places the less-frequent pair first. Sorting by frequency would reverse the records, so this case distinguishes diagnostic ordering from future merge selection.
- Observed behavior, if run: the black-box harness passed with numeric left-then-right ordering.

### Questions and sign-off

- Remaining questions: none for Stage 3.
- AI feedback on gaps: the learner's algorithm was correct, and their sorting question identified the precise bug targeted by the adversarial case.
- Learner acceptance: accepted as satisfying enough in terms of implementation before moving to the next stage.
- Engineering status: `verified`
- Learning status: `signed_off`
- Date signed off: 2026-10-22.

## Stage 4 — Train the first merge

### Contract

- Stage name: Train the first merge
- New observable behavior: `train --text TEXT --merges 1` performs one BPE step and prints the canonical version, base vocabulary size, and selected merge-rule array.
- Main concept: one BPE step selects by frequency but replaces left-to-right without overlap.
- Explicit non-goals: tied winners, caller-supplied pairs or new IDs, more than one requested merge, stopping rules, persistence, and the final invalid-input contract.
- Engineering status: `verified`
- Learning status: `signed_off`

### My prediction before implementation

- Learner-approved trace:

  > ababab -> [97,98,97,98,97,98] -> [97,98,3], [98,97,2] -> replace [97,98] with 256 -> [256,256,256]

- Temporary inspectable result originally proposed: `{"pair":[97,98],"count":3,"new_id":256,"ids":[256,256,256]}`. This is now an internal trace rather than the final JSON contract.
- Expected intermediate states: text → byte IDs → pair-frequency map → unique highest-frequency pair → non-overlapping replacement IDs.
- Why: `(97,98)` occurs three times and wins over `(98,97)`, which occurs twice; each winning occurrence is disjoint in this input.

Final contract correction before persistence: the learner chose the document's canonical model JSON. For `train --text 'ababab' --merges 1`, public output is `{"version":1,"base_vocab_size":256,"merges":[[97,98]]}`. The post-replacement `[256,256,256]` remains in internal `trainingResult` tests. Original Stage 5 is removed entirely as a public stage; this revised Stage 4 is an improved implementation of original Stage 6. Original Stage 7 is retained next for tied frequencies.

### What AI implemented

- Important files and symbols: `app/main.go` adds `runTrain`, reusable `pairFrequencies`, `selectUniqueTopPair`, `replacePair`, `firstMergeStep`, and `trainFirstMerge`. `tokenizerModel` defines the public JSON artifact, while the unexported `mergeStep` preserves the detailed algorithm state for unit tests. `challenge-cli/stages/04-train-first-merge.json` defines four black-box cases.
- Execution path in order: `run` dispatches `train`; `runTrain` parses user-supplied `--text` and `--merges`; training encodes bytes, counts pairs, selects the unique highest count, assigns ID `baseVocabSize + 0`, and calls `replacePair`; `trainingResult` retains both the model and transformed IDs, but `runTrain` JSON-encodes only its model.
- Important invariants: base IDs are always 0–255 and `base_vocab_size` therefore remains 256; learned merge index 0 receives ID 256; the ordered merge list is sufficient to derive later learned IDs; counting includes overlapping windows; replacement never reuses a consumed ID; this stage requires one unique winner with count at least two; unmatched IDs retain order.
- Language or library mechanics I had to learn: `flag.Int` parses the user-supplied merge budget; JSON struct tags produce only `version`, `base_vocab_size`, and `merges`; the unexported `trainingResult` can retain corpus-specific IDs for same-package tests without serializing them into the model.

### Evidence

- Command or check: before the temporary implementation, `./challenge test --stage 4 --debug`.
- Observed result: all 4 cases failed because `merge-top` was not recognized.
- Command or check: after revising the contract, `./challenge test --stage 4 --debug`.
- Observed result: all 4 revised cases failed because `train` was not recognized and the usage still advertised `merge-top`.
- Command or check: after clarifying that the public result must retain post-replacement IDs, `./challenge test --stage 4 --debug`.
- Observed result: all 4 cases failed because the `train` result had the model fields but omitted `ids`.
- Command or check: after choosing the canonical persisted format, Stage 4 fixtures removed `ids` and `./challenge test --stage 4 --debug` was rerun.
- Observed result: all 4 cases failed because the serialized result still contained the now-extra `ids` key; separating `tokenizerModel` from `trainingResult` made them pass.
- Command or check: after implementing the revised interface, `./challenge test --stage 4 --previous`.
- Observed result: all 26 cumulative black-box cases passed: 7 from Stage 1, 9 from Stage 2, 6 from Stage 3, and 4 from Stage 4.
- Command or check: root `go test ./...` and `go vet ./...`; challenge module `go test -race ./...` and `go vet ./...`.
- Observed result: all checks exited 0 with the project-local Go caches.
- What this proves: one unique most-frequent pair is selected, assigned ID 256, and replaced left-to-right without overlap; the public command emits only canonical model fields while tests retain transformed IDs; previous behavior remains active.
- What this does not prove: deterministic tie-breaking, repeated merges, safe early stopping, saving and loading the model, encoding with a learned model, or learner understanding.

### My explanation after reading the code

Learner explanation, preserved as given:

> only 2 replacements because as the ids are getting replaced, 2 of them are replaced at once by one new id and so out of the 4 unique pairs, only 2 remain that can be replaced.
> we advance by 2 so that we dont count the ids twice that are already replaced

AI clarification: the four counted windows are four overlapping occurrences of the same pair, not four unique pairs. Replacement consumes two input positions at a time, so five identical IDs contain only two disjoint matches. Advancing by two prevents either consumed position from participating in another replacement.

### Adversarial case

- Input: `train --text 'zzzzz' --merges 1`
- Expected public behavior: `{"version":1,"base_vocab_size":256,"merges":[[122,122]]}`. The focused unit test expects the internal step to count 4 windows and retain transformed IDs `[256,256,122]`.
- Why this case is useful: four overlapping counted windows produce only two non-overlapping replacements. It catches implementations that confuse frequency with replacement count or advance by one after consuming a match.
- Observed behavior, if run: the black-box harness passed with the canonical one-rule model, and the unit test passed with count 4 and replacement IDs `[256,256,122]`.

### Questions and sign-off

- Remaining questions: none for Stage 4.
- AI feedback on gaps: the learner correctly explained consumption and advancing by two; only the distinction between pair occurrences and unique pair values needed correction. After reviewing the revised command and result contract, the learner elected to move to tied-frequency selection.
- Learner acceptance: accepted after reviewing the corrected `train --text TEXT --merges 1` implementation and inspectable result.
- Engineering status: `verified`
- Learning status: `signed_off`
- Date signed off: 2026-10-22.

## Stage 5 — Resolve tied pair frequencies

### Contract

- Stage name: Resolve tied pair frequencies
- New observable behavior: `train --text TEXT --merges 1` succeeds when several pairs share the highest frequency and selects the deterministic course-defined winner.
- Main concept: compare candidates by frequency descending, left ID ascending, then right ID ascending.
- Explicit non-goals: more than one merge, safe early stopping, persistence, and model-based encoding or decoding.
- Engineering status: `verified`
- Learning status: `signed_off`

### My prediction before implementation

Learner proposal, preserved from the stage discussion:

> sorting the id pairs in desc with respect to the freq and then asc with resp to left and then asc with resp to right, that way the top most element will always be the merge we are looking for and it will help us to do multiple merges in the next stages

- Expected intermediate states: pair-frequency map → sorted candidate pairs → first candidate → non-overlapping replacement.
- Why: the first sort key chooses a most-frequent pair, while the next two keys resolve equal frequencies without depending on Go's randomized map iteration order.

### What AI implemented

- Important files and symbols: `selectTopPair` in `app/main.go` copies map keys into a slice and sorts them with the three-level comparator. `TestSelectTopPair` separately covers an ordinary winner, a tied left-ID decision, and a tied right-ID decision. `challenge-cli/stages/05-resolve-tied-pairs.json` adds three process-level cases.
- Execution path in order: Stage 4's training path is unchanged until `firstMergeStep` calls `selectTopPair`; selection sorts all distinct candidates, checks that the first candidate occurs at least twice, and returns it with its frequency; training then assigns ID 256, replaces matches, and emits the existing result shape.
- Important invariants: higher frequency always outranks token IDs; left ID is consulted only when frequencies tie; right ID is consulted only when both frequency and left ID tie; map iteration order cannot change the winner; a pair occurring only once is still not learnable.
- Why this structure helps later: repeated training can recount after every replacement and call the same deterministic selector without changing the comparison rule.

### Evidence

- Command or check: before implementation, `./challenge test --stage 5 --debug`.
- Observed result: both tied cases failed with `most-frequent pair is tied`, while the unique-winner regression passed.
- Command or check: after implementation, `./challenge test --stage 5 --previous`.
- Observed result: all 29 cumulative black-box cases passed: 7 from Stage 1, 9 from Stage 2, 6 from Stage 3, 4 from Stage 4, and 3 from Stage 5.
- What this proves: tied winners are resolved by the specified numeric ordering and the existing unique-winner path still works.
- What this does not prove: repeated merges, safe stopping, persistence, or learner understanding after reading the implementation.

### Adversarial cases

- `bacbac` ties `(98,97)` and `(97,99)` at frequency 2. The smaller left ID selects `(97,99)`, producing `[98,256,98,256]`.
- `abacabac` ties `(97,98)`, `(98,97)`, and `(97,99)` at frequency 2. Left ID 97 eliminates `(98,97)`; right ID 98 then wins over 99, producing `[256,97,99,256,97,99]`.

### Questions and sign-off

- Remaining questions: none for Stage 5.
- AI feedback on gaps: the proposed ordering is correct and naturally reusable. The learner reviewed the comparator and tied examples and confirmed that they match the intended priority.
- Learner acceptance: accepted after reviewing the comparator and both tie-breaking examples.
- Engineering status: `verified`
- Learning status: `signed_off`
- Date signed off: 2026-10-22.

## Stage 6 — Learn successive merges

### Contract

- Stage name: Learn successive merges
- New observable behavior: `train --text TEXT --merges N` accepts a positive `N` and learns that many ordered rules when the input has enough repeated pairs.
- Main concept: every replacement changes adjacency, so training must recount the transformed ID sequence before choosing the next rule.
- Explicit non-goals: zero merges, empty or one-byte inputs, stopping before an oversized budget, persistence, and model-based encoding or decoding.
- Engineering status: `verified`
- Learning status: `signed_off`

### My prediction before implementation

Learner proposal, preserved from the stage discussion:

> counter which counts from 1 to N, for each merge, we iterate the complete ids and merge creating a new list of ids, then we rescan the pairs and make a new sorted pair list and then do this merge thing again until the count.

- Expected intermediate states: current IDs → current pair counts and ordering → selected pair → new IDs → fresh pair counts for the next iteration.
- AI clarification before implementation: the overall approach is correct. The previously sorted candidate list cannot be reused after replacement because it describes adjacencies that may no longer exist.

### What AI implemented

- Important files and symbols: `trainMerges` owns the loop and accumulated rule list; `nextMergeStep` performs one select-and-replace iteration for an existing ID sequence and assigned ID; `trainFirstMerge` remains as a Stage 4-compatible wrapper. `challenge-cli/stages/06-learn-successive-merges.json` compares one-, two-, and three-merge runs.
- Execution path in order: `runTrain` parses a positive `--merges`; `trainMerges` starts from byte IDs; merge index 0 calls `nextMergeStep` with ID 256; the returned IDs become the next iteration's input; subsequent steps use IDs 257, 258, and so on; `trainingResult` retains final transformed IDs while the serialized model contains only rules and fixed metadata.
- Important invariants: the selected pair always comes from the current sequence; counts are rebuilt on every iteration; rule `i` creates token `256 + i`; a rule may reference only base IDs or earlier learned IDs; replacement remains left-to-right and non-overlapping.
- Complexity at this stage: for each requested merge, pair counting and replacement scan the current sequence and selection sorts its distinct pairs. For the challenge's small educational inputs, this direct implementation favors clarity over incremental frequency-maintenance complexity.

### Evidence

- Command or check: before implementation, `./challenge test --stage 6 --debug`.
- Observed result: the one-merge case passed through Stage 5 behavior, while the two- and three-merge cases failed because the CLI accepted only `--merges 1`.
- Command or check: after implementation, `./challenge test --stage 6 --previous`.
- Observed result: all 32 cumulative black-box cases passed: 7 from Stage 1, 9 from Stage 2, 6 from Stage 3, 4 from Stage 4, 3 from Stage 5, and 3 from Stage 6.
- What this proves: successive rules use freshly counted sequences, learned IDs are assigned in order, and later rules can reference earlier learned IDs.
- What this does not prove: safe stopping when the requested budget is too large, zero-merge behavior, persistence, or learner understanding after reading the implementation.

### Trace and adversarial case

For `train --text 'banana banana' --merges 3`:

1. Learn `[97,110]` as ID 256: `[98,256,256,97,32,98,256,256,97]`.
2. Recount and learn `[98,256]` as ID 257: `[257,256,97,32,257,256,97]`.
3. Recount and learn `[256,97]` as ID 258: `[257,258,32,257,258]`.

The one-, two-, and three-merge fixtures form the adversarial set: an implementation that reuses the original frequency table or always assigns ID 256 cannot satisfy all three.

### Questions and sign-off

- Remaining questions: none for Stage 6.
- AI feedback on gaps: the learner's proposed control flow is the standard clear implementation. The learner reviewed the loop, intermediate sequences, learned-ID assignment, and hierarchical third rule and found them satisfactory.
- Learner acceptance: accepted after reviewing the implementation and complete three-merge trace.
- Engineering status: `verified`
- Learning status: `signed_off`
- Date signed off: 2026-10-22.

## Stage 7 — Stop training safely

### Contract

- Stage name: Stop training safely
- New observable behavior: `--merges N` is an upper bound; training succeeds after reaching it or when no pair occurs at least twice. Zero budgets, empty text, and one-byte text are valid.
- Main concept: exhausting useful work is successful termination, not an algorithm error.
- Explicit non-goals: saving models, file input, model-based encoding or decoding, and the final structured error vocabulary.
- Engineering status: `verified`
- Learning status: `signed_off`

### Assessment before implementation

- Learner decision: keep source Stage 9 as written because the boundary is useful and already appears close to the existing behavior.
- Existing pieces: the Stage 6 loop already enforces the upper budget, and `selectTopPair` already identifies an empty candidate set or a best count below 2.
- Missing pieces: `trainMerges` propagated that condition as failure instead of stopping; `runTrain` rejected zero; no black-box cases covered empty input, one-byte input, or oversized budgets.

### What AI implemented

- Important files and symbols: `errNoRepeatedPair` gives the selector's normal exhaustion condition a stable identity; `trainMerges` uses `errors.Is` to break only on that condition; `runTrain` permits explicitly supplied zero while continuing to reject missing or negative budgets. `challenge-cli/stages/07-stop-training-safely.json` covers all source boundaries through canonical model output, while unit tests inspect final IDs.
- Execution path in order: parse an explicitly supplied nonnegative budget; begin with byte IDs and an empty merge list; attempt iterations up to the budget; if selection reports `errNoRepeatedPair`, break without adding a rule; return every rule actually learned and the current IDs.
- Important invariants: `len(merges) <= mergeLimit`; no token is invented from a pair occurring once; zero iterations preserve all input IDs; early stopping preserves rules already learned; ordinary exhaustion writes a successful model rather than stderr.
- Why a sentinel error is used: selection still needs to tell a one-step caller that no merge exists, while the outer trainer needs to recognize that exact condition as a clean stop without swallowing unrelated future errors.

### Evidence

- Command or check: before implementation, `./challenge test --stage 7 --debug`.
- Observed result: all 6 cases failed—five with `no pair occurs at least twice` and the zero-budget case with CLI usage.
- Command or check: after implementation, `./challenge test --stage 7 --previous`.
- Observed result: all 38 cumulative black-box cases passed: 7 from Stage 1, 9 from Stage 2, 6 from Stage 3, 4 from Stage 4, 3 from Stage 5, 3 from Stage 6, and 6 from Stage 7.
- What this proves: both stop conditions work, zero and short inputs preserve their IDs, and oversized budgets return only useful rules.
- What this does not prove: persistence, file I/O, model reuse, the final error contract, or learner understanding after reading the implementation.

### Boundary examples

- `abc --merges 10` learns no rules because every adjacent pair has count 1.
- `aaaaa --merges 10` learns only `[97,97]`, producing `[256,256,97]`; the remaining pairs each occur once.
- `banana banana --merges 99` learns four rules and stops at `[259,32,259]` because both remaining adjacent pairs occur once.
- `abab --merges 0` performs no work and preserves `[97,98,97,98]`.

### Questions and sign-off

- Remaining questions: none for Stage 7.
- AI feedback on gaps: the learner correctly recognized that most mechanics were already in place, reviewed the stopping work, and then identified and corrected the boundary between model state and corpus-specific inspection data.
- Learner acceptance: accepted together with the correction to emit only the canonical three-field model JSON.
- Engineering status: `verified`
- Learning status: `signed_off`
- Date signed off: 2026-10-22.

## Stage 8 — Train from a file and save the model

### Contract

- Stage name: Train from a file and save the model
- New observable behavior: `train` accepts raw corpus bytes through `--input PATH` instead of `--text`, and optional `--output PATH` saves the same canonical model printed to stdout.
- Main concept: training input and learned model output are different file boundaries—the corpus is arbitrary bytes, while the model is canonical JSON.
- Explicit non-goals: loading a model, encoding with learned rules, model validation, recursive directory creation, and the final structured error vocabulary.
- Engineering status: `verified`
- Learning status: `signed_off`

### My prediction before implementation

Learner proposal, preserved from the stage discussion:

> I think we need to just include logic to read from a file and the rest mostly stays the same, we maybe dont have to do the json parsing now since we are reading directly from the doucment

- Correct part: once corpus bytes are obtained, pair counting, deterministic selection, repeated merging, safe stopping, and model construction stay unchanged.
- Clarification: `--input` names the raw training corpus rather than a JSON document, so no JSON parsing is needed for input. JSON serialization is still required to print and optionally save `model.json`; JSON parsing begins when a later stage loads that model.

### What AI implemented

- Important files and symbols: `runTrain` tracks whether `--text`, `--input`, `--output`, and `--merges` were explicitly supplied; `os.ReadFile` loads the complete corpus; `trainBytes` trains directly on `[]byte`; `json.Marshal` produces the canonical model once; `os.WriteFile` persists those bytes. `challenge-cli/stages/08-save-trained-model.json` supplies isolated corpus files and inspects generated models.
- Execution path in order: parse flags → require exactly one corpus source and a nonnegative merge budget → read raw file bytes or convert explicit text to bytes → call `trainBytes` → marshal only `tokenizerModel` → optionally write the encoded model → write the same bytes to stdout.
- Important invariants: text and file input with identical bytes learn identical rules; no newline is trimmed; arbitrary bytes remain unsigned IDs 0–255; stdout and `--output` cannot differ because both receive the same `encodedModel` byte slice; corpus-specific transformed IDs are never persisted.
- Language mechanics: Go strings can contain arbitrary bytes, but a direct `[]byte` trainer makes the file contract explicit. `flag.FlagSet.Visit` distinguishes an explicitly supplied empty `--text ''` from a missing `--text` flag.

### Evidence

- Command or check: before implementation, `./challenge test --stage 8 --debug`.
- Observed result: the existing `--text` case passed; all four file-related cases failed with CLI usage because `--input` and `--output` were unknown.
- Command or check: after implementation, `./challenge test --stage 8 --previous`.
- Observed result: all 43 cumulative black-box cases passed: 7 from Stage 1, 9 from Stage 2, 6 from Stage 3, 4 from Stage 4, 3 from Stage 5, 3 from Stage 6, 6 from Stage 7, and 5 from Stage 8.
- What this proves: file and text corpus sources agree, model persistence matches stdout structurally, newline bytes survive, arbitrary bytes reach training unchanged, and prior stages still pass.
- What this does not prove: model loading, behavior in a fresh encoding process, complete I/O error classification, or learner understanding after reading the implementation.

### Adversarial cases

- `a\na\n` learns `[97,10]`; trimming or line scanning would change or eliminate that rule.
- Raw bytes `FF 00 FF 00` learn `[255,0]`; text-only or UTF-8-decoding input code would corrupt or reject this corpus.
- `banana banana` through `--text` and `--input` must produce the same three-rule canonical model.

### Questions and sign-off

- Remaining questions: none for Stage 8.
- AI feedback on gaps: the proposed reuse of training logic was correct. The learner reviewed the corpus-source validation, raw byte path, and shared stdout/file serialization and found them satisfactory.
- Learner acceptance: accepted after reviewing file input, canonical model output, and the distinction between corpus bytes and model JSON.
- Engineering status: `verified`
- Learning status: `signed_off`
- Date signed off: 2026-10-22.

## Stage template

### Contract

- Stage name:
- New observable behavior:
- Main concept:
- Explicit non-goals:
- Engineering status: `not_started`
- Learning status: `not_started`

### My prediction before implementation

- Tiny input:
- Expected output:
- Expected intermediate states:
- Why:

### What AI implemented

- Important files and symbols:
- Execution path in order:
- Important invariants:
- Language or library mechanics I had to learn:

### Evidence

- Command or check:
- Observed result:
- What this proves:
- What this does not prove:

### My explanation after reading the code

_Explain the mechanism without copying the implementation comments._

### Adversarial case

- Input:
- Expected behavior:
- Why this case is useful:
- Observed behavior, if run:

### Questions and sign-off

- Remaining questions:
- AI feedback on gaps:
- Engineering status:
- Learning status:
- Date signed off:

