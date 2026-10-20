<!--
@file docs/research-plan.md
@description Defines the questions, source standards, and outputs for the tokenizer foundations research phase.
@depends_on plan.md, collaboration-guide.md, resources/README.md
@depended_on_by README.md, docs/README.md, docs/tokenizer-foundations.md
-->

# Tokenizer foundations research plan

## Purpose

The research phase should build enough context to understand what this project is constructing before judging detailed stage boundaries. It should not become an unbounded survey of every tokenizer ever used.

## Questions to answer

1. What problem does tokenization solve, and where does it sit between raw user text and a language model?
2. How do characters, Unicode code points, UTF-8 bytes, tokens, token IDs, and vocabulary entries differ?
3. What happens during tokenizer training, model training, prompt encoding, model inference, and output decoding?
4. How do word-based, character-based, byte-based, BPE, WordPiece, Unigram, and byte-level BPE approaches differ?
5. What roles do normalization, pre-tokenization, merge ranks, special tokens, vocabulary limits, and fallback behavior play?
6. What properties matter: reversibility, determinism, coverage, compression, speed, model compatibility, and handling of multilingual or unusual text?
7. Which details belong to the educational challenge, and which would appear only in a production tokenizer?

## Source standard

Prefer primary sources, official documentation or repositories, seminal papers, and exceptionally clear technical explanations. Use videos when animation or live tracing improves understanding, not merely because the same prose exists in video form. Record why each source is useful, what question it answers, its level, and whether its claims were cross-checked.

Candidate links belong in `resources/sources.md`. Inclusion there does not imply endorsement. The research phase should replace weak or redundant sources rather than accumulating a large unread list.

## Expected outputs

- a short mental model in the learner's own words;
- one end-to-end context diagram from text input to model tokens and back;
- a terminology map with small examples;
- a comparison of major tokenizer families;
- a scoped explanation of the selected byte-level BPE variant;
- open questions that should influence the challenge audit.

