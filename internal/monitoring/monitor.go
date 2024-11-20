package monitoring

import (
    "sync"
    "time"
)

type ResourceMetric struct {
    Name      string
    Type      string
    Status    string
    Timestamp time.Time
    Metrics   map[string]float64
}

type Monitor struct {
    metrics map[string]*ResourceMetric
    mu      sync.RWMutex
    logger  *Logger
}

func NewMonitor(logger *Logger) *Monitor {
    return &Monitor{
        metrics: make(map[string]*ResourceMetric),
        logger:  logger,
    }
}

func (m *Monitor) RecordMetric(resourceName string, metricType string, value float64) {
    m.mu.Lock()
    defer m.mu.Unlock()

    metric, exists := m.metrics[resourceName]
    if !exists {
        metric = &ResourceMetric{
            Name:      resourceName,
            Timestamp: time.Now(),
            Metrics:   make(map[string]float64),
        }
        m.metrics[resourceName] = metric
    }

    metric.Metrics[metricType] = value
    metric.Timestamp = time.Now()
}

func (m *Monitor) GetMetrics(resourceName string) *ResourceMetric {
    m.mu.RLock()
    defer m.mu.RUnlock()

    return m.metrics[resourceName]
}

func (m *Monitor) StartHealthCheck(interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()

        for range ticker.C {
            m.performHealthCheck()
        }
    }()
}

func (m *Monitor) performHealthCheck() {
    m.mu.RLock()