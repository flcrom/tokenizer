<!--
@file docs/version-control.md
@description Defines the Git branch roles and cumulative CodeCrafters-style commit policy for verified stages.
@depends_on plan.md
@depended_on_by README.md, docs/README.md, docs/harness-architecture.md
-->

# Version-control strategy

## Branch roles

`main` is the project record. It contains source material, research, learner artifacts, challenge-design decisions, reusable tooling, and approved integrated outcomes.

`stages` is the cumulative implementation timeline. It begins at the shared repository foundation. After a stage passes its approved engineering checks, that stage's code, tests, fixtures, and directly related explanation are committed together. The next stage builds on that commit, matching the cumulative experience of a CodeCrafters challenge.

## Stage commit convention

CodeCrafters-generated commits commonly use this subject:

```text
codecrafters submit [skip ci]
```

We will preserve that subject on `stages`. Because the subject alone does not identify the stage, every verified stage commit will also receive an annotated tag:

```text
stage-01-encode-ascii
stage-02-encode-utf8
```

The tag message records the approved stage name and verification evidence. A stage is not committed or tagged as passed merely because code exists or tests are green. Its engineering status must be `verified`, its learning status must be `signed_off`, and the learner must accept it as satisfactory.

## Relationship between branches

`main` remains the living workspace for design, test authoring, implementation work, and review. `stages` is a curated checkpoint history: the scaffold is shared by both branches, then each accepted stage is represented by one cumulative CodeCrafters-style commit. The local challenge CLI never performs these Git operations.

