<!--
@file docs/source-orientation.md
@description Gives a high-level, non-evaluative orientation to the supplied 16-stage tokenizer challenge.
@depends_on resources/README.md
@depended_on_by README.md, docs/README.md, docs/challenge-audit.md, docs/tokenizer-foundations.md
-->

# High-level orientation to the supplied challenge

## The destination

The document proposes a command-line program that learns a small byte-level Byte Pair Encoding model from text, saves that learned model, uses it to turn new text into integer token IDs, and turns those IDs back into the original bytes.

That gives the finished system three everyday operations:

1. **train**: inspect a corpus and learn an ordered list of merge rules;
2. **encode**: apply the saved rules to new text and produce token IDs;
3. **decode**: expand token IDs back into exactly the original bytes.

The IDs belong only to the model produced by this project. They are not automatically GPT, Claude, Llama, or any other existing model's token IDs.

## The simplest useful mental model

An LLM consumes numbers, not strings. A tokenizer is a reversible translator between text and a sequence of integer IDs. This challenge starts with one token for every possible byte value, then learns that frequently adjacent token pairs can be represented by new token IDs.

For example, training might learn that byte IDs for `a` and `n` often occur together and assign their pair a new ID, `256`. Later it may combine that learned token with another token. The ordered merge rules are therefore both the learned model and the recipe needed to reconstruct what every learned token means.

## Why the challenge begins with bytes

UTF-8 turns text into bytes. There are exactly 256 possible byte values, so a base vocabulary containing IDs `0` through `255` can represent any UTF-8 input without needing an “unknown character” token. One visible character may require multiple bytes, which is why the document separates ASCII, Unicode, and raw-byte decoding concepts.

## How the 16 stages are grouped

- Stages 1–3 establish byte encoding and lossless decoding.
- Stages 4–5 expose pair counting and pair replacement as diagnostic commands.
- Stages 6–9 learn one or more deterministic merge rules and stop safely.
- Stages 10–12 save a model, load it in another process, and apply a learned pair.
- Stages 13–15 handle rule priority, hierarchical learned tokens, and decoding learned IDs.
- Stage 16 adds a consistent validation and error contract.

This grouping describes the proposal; it does not yet say that every boundary is correctly sized.

## Core behavior versus teaching scaffolding

`train`, `encode`, and `decode` form the final product workflow. The proposed `pairs` and `merge` commands expose two internal algorithmic steps so a learner and the black-box test suite can inspect them independently. They are pedagogical diagnostics rather than commands a typical end user would expect from a production tokenizer.

## Deliberate simplifications

The document defines a small course-specific model format, one whole-input byte sequence, deterministic tie-breaking, a minimum pair frequency of two during training, and an explicit ordered pair-rule representation. It excludes normalization, pre-tokenization, regex splitting, special tokens, compatibility with existing model vocabularies, padding, truncation, streaming, and optimized large-corpus training.

These boundaries make the algorithm inspectable. Later research must distinguish which rules define BPE generally and which belong only to this course.

## What you are expected to understand now

At this point, it is enough to understand the destination and the major phases. You are not expected to trace all 16 stages, judge their sizes, understand UTF-8 edge cases, or approve the proposed model format. Those become meaningful after the foundations phase.

