<!--
@file challenge-cli/README.md
@description Explains the isolated local challenge CLI, its commands, boundaries, and verification responsibilities.
@depends_on challenge-cli/go.mod, challenge-cli/cmd/challenge/main.go, challenge-cli/stages/README.md, your_program.sh
@depended_on_by README.md, docs/harness-architecture.md
-->

# Local challenge CLI

This module is a small local substitute for the part of the CodeCrafters CLI that creates stage definitions and runs black-box tests. It is deliberately not a submission service, Git client, course platform, or tokenizer library.

Run it from the repository root:

```bash
./challenge list
./challenge new-stage --number 1 --slug encode-ascii --title "Encode ASCII text"
./challenge test --stage 1
./challenge test --stage 1 --previous
./challenge test --stage 1 --case "encodes greeting" --debug
```

The test runner starts `./your_program.sh` in a fresh temporary directory for every attempt. It observes only arguments, stdin, stdout, stderr, exit status, duration, and declared files. It never imports the learner application and never runs Git commands.

The CLI's own tests use `testdata/programs/fake_program.sh` to verify the process boundary independently of tokenizer behavior. Stage manifests are described in [stages/README.md](stages/README.md).

