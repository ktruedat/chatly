pub mod text_models;
pub mod image_models;

#[cfg(test)]
mod tests;

use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::RwLock;

/// Tipul de model (Text sau Imagine)
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub enum ModelType {
    Text,
    Image,
}

/// Statistici de performanță pentru un model
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ModelMetrics {
    pub total_requests: u64,
    pub successful_requests: u64,
    pub failed_requests: u64,
    pub avg_latency_ms: f64,
    pub p95_latency_ms: f64,
    pub p99_latency_ms: f64,
}

impl Default for ModelMetrics {
    fn default() -> Self {
        Self {
            total_requests: 0,
            successful_requests: 0,
            failed_requests: 0,
            avg_latency_ms: 0.0,
            p95_latency_ms: 0.0,
            p99_latency_ms: 0.0,
        }
    }
}

/// Trait comun pentru toate modelele AI
#[async_trait]
pub trait Model: Send + Sync {
    /// Returnează numele modelului
    fn name(&self) -> &str;
    
    /// Returnează tipul modelului (Text/Image)
    fn model_type(&self) -> ModelType;
    
    /// Execută predicția
    async fn predict(&self, input: &str) -> anyhow::Result<String>;
    
    /// Returnează metricile curente
    fn get_metrics(&self) -> ModelMetrics;
    
    /// Actualizează metricile după o cerere
    fn update_metrics(&self, latency_ms: f64, success: bool);
    
    /// Verifică dacă modelul este disponibil (pentru circuit breaker)
    fn is_available(&self) -> bool;
}

/// Wrapper pentru a gestiona metrici thread-safe
#[derive(Clone)]
pub struct MetricsWrapper {
    pub metrics: Arc<RwLock<ModelMetrics>>,
    pub latency_history: Arc<RwLock<Vec<f64>>>,
}

impl MetricsWrapper {
    pub fn new() -> Self {
        Self {
            metrics: Arc::new(RwLock::new(ModelMetrics::default())),
            latency_history: Arc::new(RwLock::new(Vec::with_capacity(1000))),
        }
    }
    
    pub async fn update(&self, latency_ms: f64, success: bool) {
        let mut metrics = self.metrics.write().await;
        let mut history = self.latency_history.write().await;
        
        metrics.total_requests += 1;
        if success {
            metrics.successful_requests += 1;
        } else {
            metrics.failed_requests += 1;
        }
        
        // Adaugă latența în istoric
        history.push(latency_ms);
        
        // Păstrează doar ultimele 1000 de cereri
        if history.len() > 1000 {
            history.remove(0);
        }
        
        // Calculează media
        if !history.is_empty() {
            metrics.avg_latency_ms = history.iter().sum::<f64>() / history.len() as f64;
            
            // Calculează percentile
            let mut sorted = history.clone();
            sorted.sort_by(|a, b| a.partial_cmp(b).unwrap());
            
            let p95_idx = (sorted.len() as f64 * 0.95) as usize;
            let p99_idx = (sorted.len() as f64 * 0.99) as usize;
            
            metrics.p95_latency_ms = sorted.get(p95_idx).copied().unwrap_or(0.0);
            metrics.p99_latency_ms = sorted.get(p99_idx).copied().unwrap_or(0.0);
        }
    }
    
    pub async fn get_metrics(&self) -> ModelMetrics {
        self.metrics.read().await.clone()
    }
}
