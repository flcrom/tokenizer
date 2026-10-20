<!--
@file docs/tokenizer-foundations.md
@description Introduces tokenizer context, core terminology, major families, the project's scope, and a curated learning sequence.
@depends_on docs/source-orientation.md, docs/research-plan.md, resources/sources.md
@depended_on_by README.md, docs/README.md
-->

# Tokenizer foundations

## Start with the real problem

Imagine typing `banana!` into a language model. The model does not directly receive six letters and punctuation. Its neural network performs mathematics on vectors indexed by integers, so something must translate the text into a finite sequence of integer IDs. That translator is the tokenizer.

The complete path is roughly:

```text
human text
  -> tokenizer encode
  -> token IDs
  -> model embeddings and neural-network computation
  -> predicted token IDs
  -> tokenizer decode
  -> output text
```

The tokenizer is adjacent to the LLM, but it is not the LLM. It has its own vocabulary and may have its own training process. Once an LLM is trained against a particular tokenizer, the mapping cannot be swapped casually: token ID `257` means only what that tokenizer's vocabulary says it means, and the model's learned embedding at position `257` is tied to that meaning.

## What, why, where, and when

**What is it?** A tokenizer defines a mapping between text or bytes and token IDs. A token may be a byte, character, word, word fragment, punctuation pattern, or special marker depending on the tokenizer.

**Why is it needed?** Neural networks require finite numeric inputs. Pure word vocabularies become enormous and fail on unseen words; pure character or byte sequences can become very long. Subword tokenizers seek a useful compromise: common sequences become single tokens while uncommon text can still be decomposed into smaller known pieces.

**Where is it used?** It processes training data before or during LLM training, encodes every user prompt at inference time, determines token-based context and billing counts, and decodes the model's predicted IDs into output text.

**When is it trained?** Usually before the language model. Tokenizer training studies a representative text corpus and produces a fixed vocabulary or model. LLM training then consumes IDs produced by that fixed tokenizer. At inference time the tokenizer normally applies its learned rules; it does not retrain on each prompt.

## Six units that must not be confused

- A **character** is the human idea of a written symbol, but even that can be composed in multiple Unicode ways.
- A **Unicode code point** is a numbered element such as `U+00E9` for precomposed `é`.
- A **UTF-8 byte** is one 0–255 value used to encode Unicode text. One code point may use several bytes.
- A **token** is one unit chosen by a tokenizer.
- A **token ID** is the integer name of that token in one specific tokenizer model.
- A **vocabulary** is the complete mapping of available tokens and IDs, often accompanied by rules that determine how text is split.

This challenge starts with all 256 byte values as tokens. It then learns larger tokens by repeatedly joining adjacent token pairs.

## What BPE learns

Byte Pair Encoding begins with small symbols and repeatedly finds a frequent adjacent pair. It creates a new symbol for that pair, replaces non-overlapping occurrences, recounts, and continues. In byte-level BPE, the starting symbols represent UTF-8 bytes.

Training and encoding are different:

- **training** discovers and orders merge rules from a corpus;
- **encoding** applies already learned rules to new input;
- **decoding** expands each token back to bytes and joins the bytes before interpreting valid UTF-8 text.

This distinction explains an initially surprising rule in the source document: a pair that is rare in a new prompt may merge before a more frequent pair because encoding follows the saved training rank, not fresh prompt frequency.

## Major tokenizer families

Word-level tokenizers use whole words, which is intuitive but creates unknown-word and vocabulary-size problems. Character-level tokenizers cover more text but create long sequences. Subword methods sit between them: frequent words or fragments stay compact while rare words split into smaller pieces.

BPE grows a vocabulary through repeated merges. WordPiece uses a related subword idea but chooses vocabulary additions using a different scoring objective. Unigram begins with many candidate pieces and probabilistically removes less useful ones. SentencePiece is a system that can train BPE or Unigram models directly from raw text and represents spaces explicitly. Byte-level BPE adds complete byte coverage to the BPE idea.

Our project studies a deliberately small byte-level BPE core. It is an excellent way to understand merge training, ranks, encoding, decoding, and byte correctness, but it is not the complete preprocessing pipeline of a production GPT tokenizer.

## Recommended learning order

1. Read the Hugging Face tokenizer-algorithms summary for the landscape: word, character, BPE, WordPiece, Unigram, and SentencePiece.
2. Watch Andrej Karpathy's “Let's build the GPT Tokenizer” in sections. It supplies the motivation, Unicode/UTF-8 context, BPE implementation, production complications, and tokenizer-related LLM behavior. Do not try to retain the entire two-hour lecture in one pass.
3. Work through Hugging Face's BPE chapter and predict its merges before revealing the next step.
4. Return to this project's source orientation and trace one tiny training example.
5. Read `minbpe` only after the algorithm makes sense; use it as a comparison implementation, not code to copy blindly.
6. Inspect OpenAI's educational `tiktoken` implementation later, when regex pre-tokenization and merge ranks can be compared with our smaller course model.

The curated links and what each source contains are maintained in `resources/sources.md`.

## Questions to carry forward

- Why do bytes give universal coverage but sometimes produce partial UTF-8 tokens?
- How does the training corpus influence which languages or code patterns receive compact tokens?
- Why must an LLM and its tokenizer agree on the exact token-ID mapping?
- Which production layers—normalization, pre-tokenization, regex splitting, and special tokens—should be taught after the core rather than inside it?

