# 🤖 Chatly AI Models

## 🔍 Request Categorizer

### OpenAI GPT OSS 20B

- Model ID: openai/gpt-oss-20b:free
- Link: [OpenRouter - GPT OSS 20B](https://openrouter.ai/openai/gpt-oss-20b:free)
- Latency:
    - 🕒 Average: ~400–700 ms
    - ⚡ P95: ~1000 ms
- How It Works:
    - Lightweight open-source model optimized for classification and categorization tasks.
    - Uses instruction-following capabilities to understand task requirements.
    - Fast inference due to smaller parameter count focused on text understanding.
- Best For:
    - Quick categorization of user requests
    - Intent classification
    - Routing decisions
- Notes:
    - Specialized for determining request complexity and routing to appropriate models.
    - Optimized for single-word classification responses.

---

## 💬 Easy Requests (Quick Conversations & Simple Tasks)

### Meta Llama 3.3 8B Instruct

- Model ID: meta-llama/llama-3.3-8b-instruct:free
- Link: [OpenRouter - Llama 3.3 8B Instruct](https://openrouter.ai/meta-llama/llama-3.3-8b-instruct:free)
- Latency:
    - 🕒 Average: ~500–800 ms
    - ⚡ P95: ~1200 ms
- How It Works:
    - Compact 8B parameter model from Meta's Llama 3.3 series with instruction tuning.
    - Uses efficient transformer architecture optimized for fast inference.
    - Trained on diverse conversational data with focus on helpfulness and accuracy.
    - Employs grouped-query attention (GQA) for reduced memory footprint and faster processing.
- Best For:
    - Quick chatbot responses
    - Simple Q&A and information retrieval
    - Basic calculations and conversions
    - General knowledge questions
    - Lightweight conversational assistants
- Notes:
    - Excellent balance of speed and quality for simple tasks.
    - Strong performance on straightforward questions requiring factual answers.

---

## 🧠 Advanced Requests (Complex Reasoning & Analysis)

### OpenAI GPT OSS 20B

- Model ID: openai/gpt-oss-20b:free
- Link: [OpenRouter - GPT OSS 20B](https://openrouter.ai/openai/gpt-oss-20b:free)
- Latency:
    - 🕒 Average: ~800–1200 ms
    - ⚡ P95: ~1800 ms
- How It Works:
    - Open-source 20B parameter model with strong reasoning and analytical capabilities.
    - Uses multi-head attention with optimized inference for balanced speed and quality.
    - Trained on diverse datasets including scientific papers, technical documentation, and reasoning tasks.
    - Employs advanced prompting techniques for improved logical reasoning and nuanced responses.
- Best For:
    - Complex problem-solving and multi-step reasoning
    - Technical explanations and analysis
    - Mathematical and logical tasks
    - Research assistance and detailed explanations
    - Nuanced discussions requiring context awareness
- Notes:
    - Good balance between speed and advanced reasoning capabilities.
    - Handles complex queries while maintaining reasonable response times.

---

## 💻 Coding Requests (Programming & Development)

### Qwen 2.5 Coder 32B Instruct

- Model ID: qwen/qwen-2.5-coder-32b-instruct:free
- Link: [OpenRouter - Qwen 2.5 Coder 32B Instruct](https://openrouter.ai/qwen-2.5-coder-32b-instruct:free)
- Latency:
    - 🕒 Average: ~1200–1600 ms
    - ⚡ P95: ~2200 ms
- How It Works:
    - Specialized 32B parameter code generation model trained on massive programming datasets (GitHub, Stack Overflow, documentation).
    - Uses fill-in-the-middle (FIM) training to understand code context bidirectionally.
    - Supports 80+ programming languages with strong performance on Python, JavaScript, Rust, C++, Go, and Java.
    - Fine-tuned on code completion, debugging, refactoring, explanation, and test generation tasks.
    - Employs instruction-following to understand developer intent and generate contextually appropriate code.
- Best For:
    - Code generation and completion
    - Debugging and error fixing
    - Code explanation and documentation
    - Algorithm implementation
    - Refactoring suggestions
    - Unit test generation
    - Code review and optimization
- Notes:
    - Larger parameter count (32B) provides superior code understanding and generation.
    - Understands code structure, patterns, and best practices across multiple languages.
    - Can generate complete functions, classes, and even small applications.

---

## 🎨 Image Processing Requests (Vision & Multimodal)

### Qwen 2.5 VL 32B Instruct

- Model ID: qwen/qwen2.5-vl-32b-instruct:free
- Link: [OpenRouter - Qwen 2.5 VL 32B Instruct](https://openrouter.ai/qwen2.5-vl-32b-instruct:free)
- Latency:
    - 🕒 Average: ~2500–3500 ms (with image input)
    - ⚡ P95: ~4500 ms
- How It Works:
    - Vision-Language model with 32B parameters combining ViT (Vision Transformer) encoder and language decoder.
    - Processes images through patch embedding, creating visual tokens that are jointly attended with text tokens.
    - Cross-modal attention layers enable reasoning about visual and textual information simultaneously.
    - Instruction-tuned specifically for visual question answering, image analysis, and multimodal tasks.
- Best For:
    - Image captioning and detailed descriptions
    - Visual question answering (VQA)
    - Scene understanding and object detection
    - OCR and text extraction from images
    - Multimodal reasoning combining text and images
- Notes:
    - Strong performance on complex visual reasoning tasks.
    - Supports high-resolution images with detailed analysis.

---

### NVIDIA Nemotron Nano 12B V2 VL

- Model ID: nvidia/nemotron-nano-12b-v2-vl:free
- Link: [OpenRouter - Nemotron Nano 12B V2 VL](https://openrouter.ai/nvidia/nemotron-nano-12b-v2-vl:free)
- Latency:
    - 🕒 Average: ~1800–2400 ms (with image input)
    - ⚡ P95: ~3200 ms
- How It Works:
    - Compact vision-language model optimized by NVIDIA for efficient inference on GPUs.
    - Uses grouped-query attention (GQA) to reduce memory bandwidth requirements.
    - Employs NVIDIA's proprietary quantization and optimization techniques for faster processing.
    - Dual-encoder architecture processes images and text in parallel streams before fusion.
- Best For:
    - Fast image analysis and classification
    - Real-time vision applications
    - Image-to-text generation
    - Visual content moderation
    - Efficient multimodal chatbots
- Notes:
    - Faster than Qwen 2.5 VL while maintaining good accuracy.
    - Optimized for NVIDIA hardware but works well on various platforms.
    - Best choice when speed is critical for image processing tasks.
