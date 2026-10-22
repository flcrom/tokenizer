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

