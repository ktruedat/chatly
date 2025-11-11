use super::{Model, ModelMetrics, ModelType, MetricsWrapper};
use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use std::time::Instant;

/// Request structure pentru OpenRouter API
#[derive(Debug, Serialize)]
struct OpenRouterRequest {
    model: String,
    messages: Vec<Message>,
}

#[derive(Debug, Serialize)]
struct Message {
    role: String,
    content: String,
}

/// Response structure pentru OpenRouter API
#[derive(Debug, Deserialize)]
struct OpenRouterResponse {
    choices: Vec<Choice>,
}

#[derive(Debug, Deserialize)]
struct Choice {
    message: ResponseMessage,
}

#[derive(Debug, Deserialize)]
struct ResponseMessage {
    content: String,
}

/// Model Llama 3.3 8B Instruct (Free)
pub struct LlamaModel {
    name: String,
    model_id: String,
    api_key: String,
    metrics: MetricsWrapper,
    available: bool,
}

impl LlamaModel {
    pub fn new(api_key: String) -> Self {
        Self {
            name: "Llama 3.3 8B Instruct".to_string(),
            model_id: "meta-llama/llama-3.3-8b-instruct:free".to_string(),
            api_key,
            metrics: MetricsWrapper::new(),
            available: true,
        }
    }
    
    async fn call_openrouter(&self, prompt: &str) -> anyhow::Result<String> {
        let client = reqwest::Client::new();
        
        let request = OpenRouterRequest {
            model: self.model_id.clone(),
            messages: vec![Message {
                role: "user".to_string(),
                content: prompt.to_string(),
            }],
        };
        
        let response = client
            .post("https://openrouter.ai/api/v1/chat/completions")
            .header("Authorization", format!("Bearer {}", self.api_key))
            .header("Content-Type", "application/json")
            .json(&request)
            .send()
            .await?;
        
        if !response.status().is_success() {
            let status = response.status();
            let error_text = response.text().await.unwrap_or_default();
            anyhow::bail!("OpenRouter API error: {} - {}", status, error_text);
        }
        
        let result: OpenRouterResponse = response.json().await?;
        
        Ok(result
            .choices
            .first()
            .map(|c| c.message.content.clone())
            .unwrap_or_default())
    }
}

#[async_trait]
impl Model for LlamaModel {
    fn name(&self) -> &str {
        &self.name
    }
    
    fn model_type(&self) -> ModelType {
        ModelType::Text
    }
    
    async fn predict(&self, input: &str) -> anyhow::Result<String> {
        let start = Instant::now();
        
        let result = self.call_openrouter(input).await;
        
        let latency_ms = start.elapsed().as_secs_f64() * 1000.0;
        let success = result.is_ok();
        
        self.update_metrics(latency_ms, success);
        
        result
    }
    
    fn get_metrics(&self) -> ModelMetrics {
        // Blocking read pentru simplitate
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

/// Model Google Gemma 3 12B IT (Free)
pub struct GemmaModel {
    name: String,
    model_id: String,
    api_key: String,
    metrics: MetricsWrapper,
    available: bool,
}

impl GemmaModel {
    pub fn new(api_key: String) -> Self {
        Self {
            name: "Gemma 3 12B IT".to_string(),
            model_id: "google/gemma-3-12b-it:free".to_string(),
            api_key,
            metrics: MetricsWrapper::new(),
            available: true,
        }
    }
    
    async fn call_openrouter(&self, prompt: &str) -> anyhow::Result<String> {
        let client = reqwest::Client::new();
        
        let request = OpenRouterRequest {
            model: self.model_id.clone(),
            messages: vec![Message {
                role: "user".to_string(),
                content: prompt.to_string(),
            }],
        };
        
        let response = client
            .post("https://openrouter.ai/api/v1/chat/completions")
            .header("Authorization", format!("Bearer {}", self.api_key))
            .header("Content-Type", "application/json")
            .json(&request)
            .send()
            .await?;
        
        if !response.status().is_success() {
            let status = response.status();
            let error_text = response.text().await.unwrap_or_default();
            anyhow::bail!("OpenRouter API error: {} - {}", status, error_text);
        }
        
        let result: OpenRouterResponse = response.json().await?;
        
        Ok(result
            .choices
            .first()
            .map(|c| c.message.content.clone())
            .unwrap_or_default())
    }
}

#[async_trait]
impl Model for GemmaModel {
    fn name(&self) -> &str {
        &self.name
    }
    
    fn model_type(&self) -> ModelType {
        ModelType::Text
    }
    
    async fn predict(&self, input: &str) -> anyhow::Result<String> {
        let start = Instant::now();
        
        let result = self.call_openrouter(input).await;
        
        let latency_ms = start.elapsed().as_secs_f64() * 1000.0;
        let success = result.is_ok();
        
        self.update_metrics(latency_ms, success);
        
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
