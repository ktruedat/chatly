mod models;
mod metrics;
mod router;

use models::{Model, ModelType};
use models::text_models::{LlamaModel, GemmaModel};
use models::image_models::{QwenVisionModel, LlamaScoutModel};
use router::{AdaptiveRouter, adaptive_router::RoutingStrategy};
use std::sync::Arc;

#[tokio::main]
async fn main() {
    // Încarcă variabilele de mediu din fișierul .env
    dotenv::dotenv().ok();
    
    println!("🚀 Chatly AI Server - Starting...\n");
    
    // Citește API key din variabilă de mediu
    let api_key = std::env::var("OPENROUTER_API_KEY")
        .expect("❌ ERROR: OPENROUTER_API_KEY nu este setat! Setează-l cu: $env:OPENROUTER_API_KEY=\"sk-or-v1-...\"");
    
    println!("📝 Inițializare modele text...");
    let text_models: Vec<Arc<dyn Model>> = vec![
        Arc::new(LlamaModel::new(api_key.clone())),
        Arc::new(GemmaModel::new(api_key.clone())),
    ];
    
    println!("🎨 Inițializare modele imagini...");
    let image_models: Vec<Arc<dyn Model>> = vec![
        Arc::new(QwenVisionModel::new(api_key.clone())),
        Arc::new(LlamaScoutModel::new(api_key.clone())),
    ];
    
    // Afișare modele disponibile
    println!("\n✅ Modele text disponibile:");
    for model in &text_models {
        println!("  - {}", model.name());
    }
    
    println!("\n✅ Modele imagini disponibile:");
    for model in &image_models {
        println!("  - {}", model.name());
    }
    
    // Inițializare router adaptiv
    println!("\n� Inițializare Router Adaptiv...");
    let router = AdaptiveRouter::new(
        text_models.clone(),
        image_models.clone(),
        RoutingStrategy::LowestLatency,
    );
    
    println!("✅ Router adaptiv configurat cu strategia: LowestLatency");
    println!("   📌 Circuit Breaker activat pentru protecție la eșecuri");
    
    // Test predicții
    println!("\n🧪 Testare router adaptiv...");
    
    println!("\n� Test predicție text:");
    match router.predict("Salut! Cum te cheamă?", ModelType::Text).await {
        Ok(response) => println!("   ✅ Răspuns: {}", response),
        Err(e) => println!("   ❌ Eroare: {}", e),
    }
    
    println!("\n🎨 Test predicție vision:");
    match router.predict("Descrie această imagine", ModelType::Image).await {
        Ok(response) => println!("   ✅ Răspuns: {}", response),
        Err(e) => println!("   ❌ Eroare: {}", e),
    }
    
    // Afișare metrici
    println!("\n📊 Metrici modele:");
    let all_metrics = router.get_all_metrics().await;
    for (name, model_type, metrics) in all_metrics {
        println!("\n  {} ({:?}):", name, model_type);
        println!("    Total cereri: {}", metrics.total_requests);
        println!("    Cereri reușite: {}", metrics.successful_requests);
        println!("    Cereri eșuate: {}", metrics.failed_requests);
        println!("    Latență medie: {:.2}ms", metrics.avg_latency_ms);
    }
    
    println!("\n✅ Server pregătit pentru cereri!");
    println!("💡 Următorii pași:");
    println!("  1. ✅ Router adaptiv implementat");
    println!("  2. ✅ Circuit breaker implementat");
    println!("  3. ⏳ Implementare endpoint HTTP pentru /predict");
    println!("  4. ⏳ Implementare endpoint HTTP pentru /metrics");
}
