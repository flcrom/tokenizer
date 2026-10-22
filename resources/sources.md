<!--
@file resources/sources.md
@description Tracks candidate tokenizer learning sources, provenance, purpose, and review state.
@depends_on none
@depended_on_by resources/README.md, docs/tokenizer-foundations.md
-->

# Curated sources

The order below moves from mental model to algorithm to implementation. `reviewed_for_path` means the source was inspected enough to understand its teaching role; it does not mean the learner has completed it.

| Order | Source | What it contains and why to use it | Review state |
| --- | --- | --- | --- |
| 1 | [Hugging Face: Tokenization algorithms](https://huggingface.co/docs/transformers/en/tokenizer_summary) | A compact map of word-, character-, and subword-level approaches, followed by BPE, WordPiece, Unigram, and SentencePiece. Read this first for classification, not implementation detail. | `reviewed_for_path` |
| 2 | [Andrej Karpathy: Let's build the GPT Tokenizer](https://www.youtube.com/watch?v=zduSFxRajkE) | A long-form video connecting LLM inputs, Unicode and UTF-8, BPE training, encode/decode, regex pre-tokenization, special tokens, and practical tokenizer quirks. Best used chapter by chapter alongside notes. | `reviewed_for_path` |
| 3 | [Hugging Face LLM Course: Byte-Pair Encoding tokenization](https://huggingface.co/learn/llm-course/en/chapter6/5) | A worked BPE training and tokenization walkthrough with small frequency examples and an implementation. Use it to predict merges by hand. Its introductory example is not itself byte-level BPE, so keep that distinction visible. | `reviewed_for_path` |
| 4 | [Andrej Karpathy: minbpe](https://github.com/karpathy/minbpe) | A minimal educational byte-level BPE codebase with basic, regex-aware, and GPT-4-oriented tokenizers. Read after learning the mechanism; compare boundaries rather than copying it into this challenge. | `reviewed_for_path` |
| 5 | [OpenAI tiktoken educational implementation](https://github.com/openai/tiktoken/blob/main/tiktoken/_educational.py) | A compact reference implementation showing training, merge ranks, encoding, decoding, regex splitting, and visualization. Use later to compare our explicit educational model with a production-adjacent library. | `reviewed_for_path` |
| Interactive bridge | [Tiktokenizer](https://tiktokenizer.vercel.app/) | An interactive inspector for model-specific chat tokenization. It exposes message-role framing, the serialized prompt, token IDs, token count, and whitespace, making it useful for connecting our small BPE model to the extra structure seen by an LLM API. It demonstrates tokenization output rather than teaching BPE training. User-supplied and accessed 2026-10-22. | `reviewed_for_path` |
| Optional bridge | [fast.ai: Let's Build the GPT Tokenizer](https://www.fast.ai/posts/2025-10-16-karpathy-tokenizers) | A searchable text-and-code companion to Karpathy's video. Useful when revisiting a specific explanation without scrubbing through the lecture. | `reviewed_for_path` |
| Background | [GeeksforGeeks: What is Tokenization in NLP?](https://www.geeksforgeeks.org/nlp/what-is-tokenization/) | A short introductory overview retained from the original workspace. It is optional because the sources above provide stronger and more project-specific coverage. | `deprioritized` |

