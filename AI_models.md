# 🤖 Chatly AI Models

## 📝 Text Models (OpenRouter Free)

### 1. Llama 3.3 8B Instruct

- Model ID: meta-llama/llama-3.3-8b-instruct:free
- Link: [OpenRouter - Llama 3.3 8B Instruct](https://openrouter.ai/meta-llama/llama-3.3-8b-instruct:free)
- Latency:
    - 🕒 Average: ~850–1200 ms
    - ⚡ P95: ~1700 ms
- How It Works:
    - Transformer-based architecture trained by Meta on diverse multilingual datasets.
    - Fine-tuned for following instructions, summarization, and conversational flow.
    - Works by encoding input text into tokens, processing them through self-attention layers, and decoding coherent responses.
- Best For:
    - Conversational chatbots
    - General-purpose text generation
    - Code or logic explanations
- Notes:
    - Very stable and lightweight, performs well under low latency conditions.

---

### 2. Google Gemma 3 12B IT

- Model ID: google/gemma-3-12b-it:free
- Link: [OpenRouter - Gemma 3 12B IT](https://openrouter.ai/google/gemma-3-12b-it:free)
- Latency:
    - 🕒 Average: ~1100–1500 ms
    - ⚡ P95: ~2100 ms
- How It Works:
    - Uses Google’s Gemma transformer stack with deep layer normalization and dense attention heads.
    - Instruction-tuned (IT) — optimized for task reasoning, question answering, and contextual understanding.
    - Employs reinforcement learning from human feedback (RLHF) for improved instruction alignment.
- Best For:
    - Complex reasoning and multi-turn dialogue
    - Technical explanations or summarization
    - Creative and analytical text tasks
- Notes:
    - Slower than Llama 3.3, but produces more accurate and nuanced results for longer contexts.

---

## 🎨 Image Models

### 1. Qwen 2.5 VL 32B

- Model ID: qwen/qwen-2.5-vl-32b
- Link: [OpenRouter - Qwen 2.5 VL 32B](https://openrouter.ai/qwen/qwen-2.5-vl-32b)
- Latency:
    - 🕒 Average: ~2500–3200 ms (image input)
    - ⚡ P95: ~4200 ms
- How It Works:
    - Vision-Language model combining a large image encoder (ViT-based) with a 32B-parameter text decoder.
    - Processes visual embeddings and text jointly to perform reasoning about image content.
    - Can describe, classify, or answer questions about images using multimodal attention.
- Best For:
    - Image captioning and analysis
    - Visual question answering (VQA)
    - Multimodal reasoning with both text and images
- Notes:
    - Excellent balance between image understanding and textual reasoning.
    - Latency higher due to heavy visual embedding computation.

---

### 2. Llama 4 Scout

- Model ID: meta-llama/llama-4-scout
- Link: [OpenRouter - Llama 4 Scout](https://openrouter.ai/meta-llama/llama-4-scout)
- Latency:
    - 🕒 Average: ~2300–2800 ms
    - ⚡ P95: ~3800 ms
- How It Works:
    - A multimodal extension of Meta’s Llama 4 series.
    - Uses a dual-stream transformer architecture for parallel text and vision processing.
    - Can interpret visual content, generate image captions, and reason about objects and scenes.
    - Supports visual grounding (e.g., understanding spatial relationships in images).
- Best For:
    - Image description and contextual analysis
    - Scene reasoning and multimodal dialogue
    - Combining text and image tasks in a single pipeline
- Notes:
    - Faster than Qwen 2.5 VL for simpler images.
    - Produces very human-like explanations of visual content.
