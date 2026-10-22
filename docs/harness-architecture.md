<!--
@file docs/harness-architecture.md
@description Defines the approved Go starter, black-box challenge CLI, stage manifest, verification, and Git checkpoint architecture.
@depends_on plan.md, docs/stage-workflow.md, docs/version-control.md, challenge-cli/README.md, your_program.sh
@depended_on_by README.md, docs/README.md, docs/build-your-own-tokenizer.md
-->

# Harness architecture

## Boundary

The repository contains two independent Go modules. The root module is the learner's tokenizer application. `challenge-cli/` is the local test-authoring and black-box execution tool. The challenge module never imports the learner module; its only program boundary is the executable root script `./your_program.sh`.

```text
stage manifest
    -> ./challenge test
    -> isolated temporary workspace
    -> ./your_program.sh <arguments>
    -> tokenizer binary
    -> stdout, stderr, exit status, duration, and declared files
```

The CLI does not run Git commands, edit application code, decide learning sign-off, contact CodeCrafters, or advance stages.

## Learner application

The Go application follows the recognizable CodeCrafters layout: `app/main.go` is the entry point and `go.mod` belongs to the repository root. The starter intentionally exits unsuccessfully because no tokenizer stage is implemented yet.

`your_program.sh` compiles the application into ignored local build output and executes it from the caller's original working directory. Preserving the invocation directory lets future tests use relative corpus, model, and output paths inside their isolated workspace.

## Challenge CLI

The root `./challenge` wrapper runs the separate `challenge-cli` module. Its supported commands are:

- `list`: show stage manifests and whether each is a draft or runnable;
- `new-stage`: create a non-overwriting draft manifest;
- `test`: run one stage, optionally all previous stages or one named case;
- `help`: show the local command contract.

The runner supports text or base64 stdin, temporary input files, expected output files, timeouts, repetitions, exit codes, and five output comparison modes: empty, exact text, contained text, exact bytes, and structural JSON.

## Stage definition lifecycle

`new-stage` creates a JSON manifest under `challenge-cli/stages/`. JSON is treated as a generated/tool-owned format and is strictly decoded: unknown fields, duplicate stage numbers or slugs, invalid paths, malformed base64, and invalid expectation combinations are rejected before a test process starts.

For each real stage, we first approve the contract, create and fill its manifest, and confirm the new cases fail against the preceding implementation for the expected reason. A passing harness result establishes engineering evidence only. The learner reviews the code and explicitly accepts the stage before Git records the curated stage checkpoint.

## Git lifecycle

The scaffold is a normal shared foundation commit on `main` and `stages`; it is not a passed tokenizer stage. Later work can happen on `main`, while `stages` receives one cumulative `codecrafters submit [skip ci]` checkpoint plus an annotated stage tag only after the stage is verified and accepted.

## Verification boundary

The harness is verified independently with fake programs that exercise successful text output, stderr capture, structural JSON, raw bytes, nonzero exits, generated files, and timeouts. This does not verify tokenizer behavior; no tokenizer stage exists yet.

