use super::{Model, ModelMetrics, ModelType, MetricsWrapper};
use async_trait::async_trait;
use std::time::Instant;

// Modele de imagini/vision gratuite folosind OpenRouter

pub struct QwenVisionModel {
    name: String,
    model_id: String,
    api_key: String,
    metrics: MetricsWrapper,
    available: bool,
}

impl QwenVisionModel {
    pub fn new(api_key: String) -> Self {
        Self {
            name: "Qwen 2.5 VL 32B".to_string(),
            model_id: "qwen/qwen2.5-vl-32b-instruct:free".to_string(),
            api_key,
            metrics: MetricsWrapper::new(),
            available: true,
        }
    }
}

#[async_trait]
impl Model for QwenVisionModel {
    fn name(&self) -> &str {
        &self.name
    }
    
    fn model_type(&self) -> ModelType {
        ModelType::Image
    }
    
    async fn predict(&self, input: &str) -> anyhow::Result<String> {
        let start = Instant::now();
        
        // TODO: Implementare apel OpenRouter pentru analiza imagini/vision
        let result = Ok(format!("Vision analysis for: {}", input));
        
        let latency_ms = start.elapsed().as_secs_f64() * 1000.0;
        self.update_metrics(latency_ms, result.is_ok());
        
        result
    }
    
    fn get_metrics(&self) -> ModelMetrics {
        let metrics_clone = self.metrics.metrics.clone();
        tokio::task::block_in_place(|| {
            tokio::runtime::Handle::current().block_on(async {
                metrics_clone.read().await.clone()
            })
        })
    }
    
    fn update_metrics(&self, latency_ms: f64, success: bool) {
        let metrics = self.metrics.clone();
        tokio::spawn(async move {
            metrics.update(latency_ms, success).await;
        });
    }
    
    fn is_available(&self) -> bool {
        self.available
    }
}

pub struct LlamaScoutModel {
    name: String,
    model_id: String,
    api_key: String,
    metrics: MetricsWrapper,
    available: bool,
}

impl LlamaScoutModel {
    pub fn new(api_key: String) -> Self {
        Self {
            name: "Llama 4 Scout".to_string(),
            model_id: "meta-llama/llama-4-scout:free".to_string(),
            api_key,
            metrics: MetricsWrapper::new(),
            available: true,
        }
    }
}

#[async_trait]
impl Model for LlamaScoutModel {
    fn name(&self) -> &str {
        &self.name
    }
    
    fn model_type(&self) -> ModelType {
        ModelType::Image
    }
    
    async fn predict(&self, input: &str) -> anyhow::Result<String> {
        let start = Instant::now();
        
        // TODO: Implementare apel OpenRouter pentru analiza imagini/vision
        let result = Ok(format!("Vision analysis for: {}", input));
        
        let latency_ms = start.elapsed().as_secs_f64() * 1000.0;
        self.update_metrics(latency_ms, result.is_ok());
        
        result
    }
    
    fn get_metrics(&self) -> ModelMetrics {
        let metrics_clone = self.metrics.metrics.clone();
        tokio::task::block_in_place(|| {
            tokio::runtime::Handle::current().block_on(async {
                metrics_clone.read().await.clone()
            })
        })
    }
    
    fn update_metrics(&self, latency_ms: f64, success: bool) {
        let metrics = self.metrics.clone();
        tokio::spawn(async move {
            metrics.update(latency_ms, success).await;
        });
    }
    
    fn is_available(&self) -> bool {
        self.available
    }
}
