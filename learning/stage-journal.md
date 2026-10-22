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

