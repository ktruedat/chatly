// Teste pentru modelele text

#[cfg(test)]
mod tests {
    use super::super::*;
    use crate::models::text_models::{LlamaModel, GemmaModel};
    use std::env;

    fn get_test_api_key() -> String {
        env::var("OPENROUTER_API_KEY")
            .unwrap_or_else(|_| "test-api-key".to_string())
    }

    #[test]
    fn test_llama_model_creation() {
        let model = LlamaModel::new(get_test_api_key());
        assert_eq!(model.name(), "Llama 3.3 8B Instruct");
        assert_eq!(model.model_type(), ModelType::Text);
        assert!(model.is_available());
    }

    #[test]
    fn test_gemma_model_creation() {
        let model = GemmaModel::new(get_test_api_key());
        assert_eq!(model.name(), "Gemma 3 12B IT");
        assert_eq!(model.model_type(), ModelType::Text);
        assert!(model.is_available());
    }

    #[test]
    fn test_initial_metrics() {
        let model = LlamaModel::new(get_test_api_key());
        let metrics = model.get_metrics();
        
        assert_eq!(metrics.total_requests, 0);
        assert_eq!(metrics.successful_requests, 0);
        assert_eq!(metrics.failed_requests, 0);
        assert_eq!(metrics.avg_latency_ms, 0.0);
    }

    #[tokio::test]
    #[ignore] // Ignoră acest test până când API key-ul este setat
    async fn test_llama_predict() {
        let api_key = env::var("OPENROUTER_API_KEY")
            .expect("OPENROUTER_API_KEY must be set for integration tests");
        
        let model = LlamaModel::new(api_key);
        let result = model.predict("Say 'test successful'").await;
        
        assert!(result.is_ok());
        let response = result.unwrap();
        assert!(!response.is_empty());
        
        // Verifică că metricile au fost actualizate
        let metrics = model.get_metrics();
        assert_eq!(metrics.total_requests, 1);
        assert_eq!(metrics.successful_requests, 1);
        assert!(metrics.avg_latency_ms > 0.0);
    }

    #[tokio::test]
    #[ignore] // Ignoră acest test până când API key-ul este setat
    async fn test_gemma_predict() {
        let api_key = env::var("OPENROUTER_API_KEY")
            .expect("OPENROUTER_API_KEY must be set for integration tests");
        
        let model = GemmaModel::new(api_key);
        let result = model.predict("Say 'test successful'").await;
        
        assert!(result.is_ok());
        let response = result.unwrap();
        assert!(!response.is_empty());
    }

    #[tokio::test]
    #[ignore]
    async fn test_metrics_update_on_multiple_requests() {
        let api_key = env::var("OPENROUTER_API_KEY")
            .expect("OPENROUTER_API_KEY must be set for integration tests");
        
        let model = LlamaModel::new(api_key);
        
        // Trimite 3 cereri
        for i in 0..3 {
            let _ = model.predict(&format!("Request {}", i)).await;
        }
        
        // Verifică metricile
        let metrics = model.get_metrics();
        assert_eq!(metrics.total_requests, 3);
        assert!(metrics.avg_latency_ms > 0.0);
        assert!(metrics.p95_latency_ms > 0.0);
    }

    #[test]
    fn test_model_comparison() {
        let llama = LlamaModel::new(get_test_api_key());
        let gemma = GemmaModel::new(get_test_api_key());
        
        // Ambele modele ar trebui să fie de tip Text
        assert_eq!(llama.model_type(), gemma.model_type());
        
        // Dar au nume diferite
        assert_ne!(llama.name(), gemma.name());
    }
}
