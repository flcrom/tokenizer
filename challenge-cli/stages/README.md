<!--
@file challenge-cli/stages/README.md
@description Documents the declarative stage manifest format created and executed by the local challenge CLI.
@depends_on challenge-cli/internal/spec/types.go, challenge-cli/internal/cli/cli.go
@depended_on_by challenge-cli/README.md
-->

# Stage manifests

Each approved stage receives one `NN-slug.json` manifest. JSON is intentionally used as a tool-owned format: `./challenge new-stage` can create it without editing Go source, and the generic runner can strictly validate it before running the learner's program.

No tokenizer stage is defined during the harness checkpoint. A new draft begins with:

```bash
./challenge new-stage --number 1 --slug encode-ascii --title "Encode ASCII text"
```

A case can supply arguments, stdin as text or base64, temporary input files, a timeout, and repetitions. It can require an exit code; empty, exact-text, exact-byte, substring, or structural-JSON stdout and stderr; and generated file contents. Paths must remain within the runner's temporary case workspace.

Example shape:

```json
{
  "version": 1,
  "stage": 1,
  "slug": "encode-ascii",
  "title": "Encode ASCII text",
  "cases": [
    {
      "name": "encodes greeting",
      "args": ["encode", "--text", "Hi!"],
      "expect": {
        "stdout": {"mode": "json", "json": [72, 105, 33]},
        "stderr": {"mode": "empty"},
        "exit_code": 0
      }
    }
  ]
}
```

