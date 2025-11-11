use crate::models::{Model, ModelType};
use std::sync::Arc;
use std::time::{Duration, Instant};
use tokio::sync::RwLock;

/// Strategii de rutare pentru selectarea modelului
#[derive(Debug, Clone, PartialEq)]
pub enum RoutingStrategy {
    /// Alege modelul cu cea mai mică latență medie
    LowestLatency,
    /// Alege modelul cu cea mai mare rată de succes
    HighestSuccessRate,
    /// Round-robin între modele disponibile
    RoundRobin,
    /// Weighted round-robin bazat pe performanță
    WeightedRoundRobin,
}

/// Circuit breaker pentru protecție la eșecuri
#[derive(Debug, Clone)]
struct CircuitBreaker {
    failure_threshold: u32,
    success_threshold: u32,
    timeout: Duration,
    failures: u32,
    successes: u32,
    state: CircuitState,
    last_failure_time: Option<Instant>,
}

#[derive(Debug, Clone, PartialEq)]
enum CircuitState {
    Closed,   // Funcționează normal
    Open,     // Blocat din cauza eșecurilor
    HalfOpen, // În testare după timeout
}

impl CircuitBreaker {
    fn new() -> Self {
        Self {
            failure_threshold: 5,
            success_threshold: 2,
            timeout: Duration::from_secs(60),
            failures: 0,
            successes: 0,
            state: CircuitState::Closed,
            last_failure_time: None,
        }
    }

    fn record_success(&mut self) {
        match self.state {
            CircuitState::Closed => {
                self.failures = 0;
            }
            CircuitState::HalfOpen => {
                self.successes += 1;
                if self.successes >= self.success_threshold {
                    self.state = CircuitState::Closed;
                    self.failures = 0;
                    self.successes = 0;
                }
            }
            CircuitState::Open => {}
        }
    }

    fn record_failure(&mut self) {
        self.last_failure_time = Some(Instant::now());
        
        match self.state {
            CircuitState::Closed => {
                self.failures += 1;
                if self.failures >= self.failure_threshold {
                    self.state = CircuitState::Open;
                }
            }
            CircuitState::HalfOpen => {
                self.state = CircuitState::Open;
                self.successes = 0;
            }
            CircuitState::Open => {}
        }
    }

    fn is_available(&mut self) -> bool {
        match self.state {
            CircuitState::Closed => true,
            CircuitState::Open => {
                if let Some(last_failure) = self.last_failure_time {
                    if last_failure.elapsed() >= self.timeout {
                        self.state = CircuitState::HalfOpen;
                        self.successes = 0;
                        true
                    } else {
                        false
                    }
                } else {
                    false
                }
            }
            CircuitState::HalfOpen => true,
        }
    }
}

/// Router adaptiv care selectează cel mai bun model bazat pe metrici
pub struct AdaptiveRouter {
    text_models: Vec<Arc<dyn Model>>,
    image_models: Vec<Arc<dyn Model>>,
    strategy: RoutingStrategy,
    circuit_breakers: Arc<RwLock<Vec<CircuitBreaker>>>,
    round_robin_index: Arc<RwLock<usize>>,
}

impl AdaptiveRouter {
    /// Creează un router nou cu modelele specificate
    pub fn new(
        text_models: Vec<Arc<dyn Model>>,
        image_models: Vec<Arc<dyn Model>>,
        strategy: RoutingStrategy,
    ) -> Self {
        let total_models = text_models.len() + image_models.len();
        let circuit_breakers = vec![CircuitBreaker::new(); total_models];

        Self {
            text_models,
            image_models,
            strategy,
            circuit_breakers: Arc::new(RwLock::new(circuit_breakers)),
            round_robin_index: Arc::new(RwLock::new(0)),
        }
    }

    /// Selectează cel mai bun model pentru un tip specific
    pub async fn select_model(&self, model_type: ModelType) -> Option<Arc<dyn Model>> {
        let models = match model_type {
            ModelType::Text => &self.text_models,
            ModelType::Image => &self.image_models,
        };

        if models.is_empty() {
            return None;
        }

        match self.strategy {
            RoutingStrategy::LowestLatency => self.select_by_lowest_latency(models).await,
            RoutingStrategy::HighestSuccessRate => self.select_by_success_rate(models).await,
            RoutingStrategy::RoundRobin => self.select_by_round_robin(models).await,
            RoutingStrategy::WeightedRoundRobin => self.select_by_weighted_round_robin(models).await,
        }
    }

    /// Selectează modelul cu cea mai mică latență
    async fn select_by_lowest_latency(&self, models: &[Arc<dyn Model>]) -> Option<Arc<dyn Model>> {
        let mut best_model: Option<Arc<dyn Model>> = None;
        let mut lowest_latency = f64::MAX;

        let circuit_breakers = self.circuit_breakers.read().await;

        for (idx, model) in models.iter().enumerate() {
            if !model.is_available() {
                continue;
            }

            // Verifică circuit breaker
            let cb_idx = self.get_circuit_breaker_index(model).await;
            if cb_idx < circuit_breakers.len() {
                let mut cb = circuit_breakers[cb_idx].clone();
                if !cb.is_available() {
                    continue;
                }
            }

            let metrics = model.get_metrics();
            
            // Prioritate pentru modele cu cereri reușite
            if metrics.successful_requests > 0 {
                let avg_latency = metrics.avg_latency_ms;
                if avg_latency < lowest_latency {
                    lowest_latency = avg_latency;
                    best_model = Some(Arc::clone(model));
                }
            } else if best_model.is_none() {
                // Dacă nu avem încă un model selectat, folosește primul disponibil
                best_model = Some(Arc::clone(model));
            }
        }

        best_model
    }

    /// Selectează modelul cu cea mai mare rată de succes
    async fn select_by_success_rate(&self, models: &[Arc<dyn Model>]) -> Option<Arc<dyn Model>> {
        let mut best_model: Option<Arc<dyn Model>> = None;
        let mut highest_rate = 0.0;

        let circuit_breakers = self.circuit_breakers.read().await;

        for (idx, model) in models.iter().enumerate() {
            if !model.is_available() {
                continue;
            }

            // Verifică circuit breaker
            let cb_idx = self.get_circuit_breaker_index(model).await;
            if cb_idx < circuit_breakers.len() {
                let mut cb = circuit_breakers[cb_idx].clone();
                if !cb.is_available() {
                    continue;
                }
            }

            let metrics = model.get_metrics();
            
            if metrics.total_requests > 0 {
                let success_rate = metrics.successful_requests as f64 / metrics.total_requests as f64;
                if success_rate > highest_rate {
                    highest_rate = success_rate;
                    best_model = Some(Arc::clone(model));
                }
            } else if best_model.is_none() {
                best_model = Some(Arc::clone(model));
            }
        }

        best_model
    }

    /// Selectează modelul folosind round-robin
    async fn select_by_round_robin(&self, models: &[Arc<dyn Model>]) -> Option<Arc<dyn Model>> {
        let circuit_breakers = self.circuit_breakers.read().await;
        let mut index = self.round_robin_index.write().await;

        let start_index = *index;
        loop {
            let model = &models[*index % models.len()];
            *index = (*index + 1) % models.len();

            if !model.is_available() {
                if *index == start_index {
                    return None; // Am încercat toate modelele
                }
                continue;
            }

            // Verifică circuit breaker
            let cb_idx = self.get_circuit_breaker_index(model).await;
            if cb_idx < circuit_breakers.len() {
                let mut cb = circuit_breakers[cb_idx].clone();
                if !cb.is_available() {
                    if *index == start_index {
                        return None;
                    }
                    continue;
                }
            }

            return Some(Arc::clone(model));
        }
    }

    /// Selectează modelul folosind weighted round-robin (bazat pe performanță)
    async fn select_by_weighted_round_robin(&self, models: &[Arc<dyn Model>]) -> Option<Arc<dyn Model>> {
        // Calculează weight-uri bazate pe latență și rată de succes
        let mut weights: Vec<f64> = Vec::new();
        let circuit_breakers = self.circuit_breakers.read().await;

        for (idx, model) in models.iter().enumerate() {
            if !model.is_available() {
                weights.push(0.0);
                continue;
            }

            // Verifică circuit breaker
            let cb_idx = self.get_circuit_breaker_index(model).await;
            if cb_idx < circuit_breakers.len() {
                let mut cb = circuit_breakers[cb_idx].clone();
                if !cb.is_available() {
                    weights.push(0.0);
                    continue;
                }
            }

            let metrics = model.get_metrics();
            
            if metrics.total_requests == 0 {
                weights.push(1.0); // Weight implicit pentru modele noi
            } else {
                let success_rate = metrics.successful_requests as f64 / metrics.total_requests as f64;
                let latency_score = if metrics.avg_latency_ms > 0.0 {
                    1000.0 / metrics.avg_latency_ms // Scor mai mare pentru latență mai mică
                } else {
                    1.0
                };
                
                // Combinație de rată de succes și performanță
                let weight = success_rate * 0.7 + (latency_score.min(1.0)) * 0.3;
                weights.push(weight);
            }
        }

        // Selectează bazat pe weight-uri
        let total_weight: f64 = weights.iter().sum();
        if total_weight == 0.0 {
            return None;
        }

        // Pentru simplitate, selectăm modelul cu cel mai mare weight
        let max_weight_idx = weights
            .iter()
            .enumerate()
            .max_by(|(_, a), (_, b)| a.partial_cmp(b).unwrap())
            .map(|(idx, _)| idx)?;

        Some(Arc::clone(&models[max_weight_idx]))
    }

    /// Execută o predicție folosind routerul adaptiv
    pub async fn predict(&self, input: &str, model_type: ModelType) -> anyhow::Result<String> {
        let model = self
            .select_model(model_type)
            .await
            .ok_or_else(|| anyhow::anyhow!("Nu există modele disponibile pentru tipul specificat"))?;

        let start = Instant::now();
        let result = model.predict(input).await;
        let latency_ms = start.elapsed().as_secs_f64() * 1000.0;

        // Actualizează circuit breaker
        let cb_idx = self.get_circuit_breaker_index(&model).await;
        let mut circuit_breakers = self.circuit_breakers.write().await;
        
        if cb_idx < circuit_breakers.len() {
            if result.is_ok() {
                circuit_breakers[cb_idx].record_success();
            } else {
                circuit_breakers[cb_idx].record_failure();
            }
        }

        result
    }

    /// Obține indexul circuit breaker-ului pentru un model
    async fn get_circuit_breaker_index(&self, model: &Arc<dyn Model>) -> usize {
        let model_name = model.name();
        
        // Caută în text models
        for (idx, m) in self.text_models.iter().enumerate() {
            if m.name() == model_name {
                return idx;
            }
        }
        
        // Caută în image models
        let text_offset = self.text_models.len();
        for (idx, m) in self.image_models.iter().enumerate() {
            if m.name() == model_name {
                return text_offset + idx;
            }
        }
        
        0 // Fallback
    }

    /// Obține statistici despre toate modelele
    pub async fn get_all_metrics(&self) -> Vec<(String, ModelType, crate::models::ModelMetrics)> {
        let mut all_metrics = Vec::new();

        for model in &self.text_models {
            all_metrics.push((
                model.name().to_string(),
                ModelType::Text,
                model.get_metrics(),
            ));
        }

        for model in &self.image_models {
            all_metrics.push((
                model.name().to_string(),
                ModelType::Image,
                model.get_metrics(),
            ));
        }

        all_metrics
    }

    /// Schimbă strategia de rutare
    pub fn set_strategy(&mut self, strategy: RoutingStrategy) {
        self.strategy = strategy;
    }

    /// Obține strategia curentă
    pub fn get_strategy(&self) -> &RoutingStrategy {
        &self.strategy
    }
}
