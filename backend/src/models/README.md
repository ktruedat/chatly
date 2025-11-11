# 🤖 Modele AI - Chatly

## 📝 Modele Text (OpenRouter Free)

### 1. **Llama 3.3 8B Instruct**
- **ID Model:** `meta-llama/llama-3.3-8b-instruct:free`
- **Link:** https://openrouter.ai/meta-llama/llama-3.3-8b-instruct:free
- **Caracteristici:**
  - Model open-source de la Meta
  - 8 billion parametri
  - Optimizat pentru instrucțiuni
  - Gratuit pe OpenRouter
  - Bun pentru sarcini generale de conversație

### 2. **Google Gemma 3 12B IT**
- **ID Model:** `google/gemma-3-12b-it:free`
- **Link:** https://openrouter.ai/google/gemma-3-12b-it:free
- **Caracteristici:**
  - Model de la Google
  - 12 billion parametri
  - Instruction-tuned
  - Gratuit pe OpenRouter
  - Mai performant pentru sarcini complexe

---

## 🎨 Modele Imagini (Placeholder)

### 1. **DALL-E 3** (TO DO)
- Model pentru generare imagini de la OpenAI
- Calitate înaltă
- Urmează să fie configurat

### 2. **Stable Diffusion XL** (TO DO)
- Model open-source pentru generare imagini
- Flexible și customizabil
- Urmează să fie configurat

---

## 🔧 Configurare

### Setare API Key

1. **Obține API Key de la OpenRouter:**
   - Accesează: https://openrouter.ai/keys
   - Creează un API key nou
   - Copiază key-ul

2. **Setează variabila de mediu:**

**Windows PowerShell:**
```powershell
$env:OPENROUTER_API_KEY="sk-or-v1-..."
```

**Linux/Mac:**
```bash
export OPENROUTER_API_KEY="sk-or-v1-..."
```

3. **Sau creează fișier `.env`:**
```
OPENROUTER_API_KEY=sk-or-v1-...
```

---

## 📊 Structura Codului

```
models/
├── mod.rs              # Definește trait-ul Model și tipuri comune
├── text_models.rs      # Implementări pentru Llama și Gemma
└── image_models.rs     # Implementări pentru modele de imagini (placeholder)
```

### Trait `Model`

Toate modelele implementează acest trait comun:

```rust
#[async_trait]
pub trait Model: Send + Sync {
    fn name(&self) -> &str;
    fn model_type(&self) -> ModelType;
    async fn predict(&self, input: &str) -> anyhow::Result<String>;
    fn get_metrics(&self) -> ModelMetrics;
    fn update_metrics(&self, latency_ms: f64, success: bool);
    fn is_available(&self) -> bool;
}
```

---

## 🚀 Utilizare

### Exemplu de apel:

```rust
use models::text_models::LlamaModel;
use models::Model;

#[tokio::main]
async fn main() {
    let model = LlamaModel::new("your-api-key".to_string());
    
    let response = model.predict("Explain quantum computing").await;
    
    match response {
        Ok(text) => println!("Response: {}", text),
        Err(e) => eprintln!("Error: {}", e),
    }
    
    // Verifică metricile
    let metrics = model.get_metrics();
    println!("Latență medie: {:.2}ms", metrics.avg_latency_ms);
    println!("P95: {:.2}ms", metrics.p95_latency_ms);
}
```

---

## 📈 Metrici Colectate

Pentru fiecare model, sistemul colectează:

- ✅ **Număr total de cereri**
- ✅ **Cereri reușite vs eșuate**
- ✅ **Latență medie (avg)**
- ✅ **Latență P95** - 95% din cereri sunt mai rapide
- ✅ **Latență P99** - 99% din cereri sunt mai rapide

Aceste metrici sunt folosite de **Router Adaptiv** pentru a alege automat modelul optim.

---

## 🔄 Rutare Adaptivă

Sistemul alege automat modelul în funcție de:

1. **Tip de date:**
   - Text → alege din Llama/Gemma
   - Imagine → alege din DALL-E/Stable Diffusion

2. **Performanță:**
   - Preferă modelul cu **P95 latency** mai mic
   - Evită modelele cu **rata mare de erori**

3. **Disponibilitate:**
   - Circuit breaker detectează modele defecte
   - Redirecționează automat către model funcțional

---

## 🛡️ Circuit Breaker

Protecție automată împotriva eșecurilor:

- **Closed** - normal, trimite cereri
- **Open** - prea multe erori, oprește cererile
- **Half-Open** - testează recuperarea

---

## 📝 TO DO pentru echipă

- [ ] Configurare API key în `.env`
- [ ] Testare modele text cu cereri reale
- [ ] Implementare modele pentru imagini
- [ ] Implementare router adaptiv
- [ ] Implementare circuit breaker
- [ ] Testare load cu trafic simulat
- [ ] Documentare endpoint-uri HTTP
