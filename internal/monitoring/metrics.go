package monitoring

import (
	"sync"
	"time"
)

// Metrics holds application metrics
type Metrics struct {
	mu sync.RWMutex
	
	// Request metrics
	RequestCount    map[string]int64  // endpoint -> count
	RequestDuration map[string][]time.Duration // endpoint -> durations
	ErrorCount      map[string]int64  // endpoint -> error count
	
	// Business metrics
	PersonaCount   int64
	IdentityCount  int64
	CommunityCount int64
	
	// System metrics
	StartTime      time.Time
	LastRequest    time.Time
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		RequestCount:    make(map[string]int64),
		RequestDuration: make(map[string][]time.Duration),
		ErrorCount:      make(map[string]int64),
		StartTime:       time.Now(),
	}
}

// RecordRequest records a request metric
func (m *Metrics) RecordRequest(endpoint string, duration time.Duration, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.RequestCount[endpoint]++
	m.RequestDuration[endpoint] = append(m.RequestDuration[endpoint], duration)
	m.LastRequest = time.Now()
	
	if isError {
		m.ErrorCount[endpoint]++
	}
	
	// Keep only last 1000 durations per endpoint
	if len(m.RequestDuration[endpoint]) > 1000 {
		m.RequestDuration[endpoint] = m.RequestDuration[endpoint][len(m.RequestDuration[endpoint])-1000:]
	}
}

// UpdateBusinessMetrics updates business-related metrics
func (m *Metrics) UpdateBusinessMetrics(personas, identities, communities int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.PersonaCount = personas
	m.IdentityCount = identities
	m.CommunityCount = communities
}

// GetSummary returns a metrics summary
func (m *Metrics) GetSummary() MetricsSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	summary := MetricsSummary{
		Uptime:         time.Since(m.StartTime),
		PersonaCount:   m.PersonaCount,
		IdentityCount:  m.IdentityCount,
		CommunityCount: m.CommunityCount,
		Endpoints:      make(map[string]EndpointMetrics),
	}
	
	for endpoint, count := range m.RequestCount {
		durations := m.RequestDuration[endpoint]
		errors := m.ErrorCount[endpoint]
		
		var avgDuration time.Duration
		if len(durations) > 0 {
			total := time.Duration(0)
			for _, d := range durations {
				total += d
			}
			avgDuration = total / time.Duration(len(durations))
		}
		
		errorRate := float64(0)
		if count > 0 {
			errorRate = float64(errors) / float64(count)
		}
		
		summary.Endpoints[endpoint] = EndpointMetrics{
			RequestCount:    count,
			ErrorCount:      errors,
			ErrorRate:       errorRate,
			AverageDuration: avgDuration,
		}
	}
	
	return summary
}

// MetricsSummary represents a summary of metrics
type MetricsSummary struct {
	Uptime         time.Duration                `json:"uptime"`
	PersonaCount   int64                        `json:"persona_count"`
	IdentityCount  int64                        `json:"identity_count"`
	CommunityCount int64                        `json:"community_count"`
	Endpoints      map[string]EndpointMetrics   `json:"endpoints"`
}

// EndpointMetrics represents metrics for a specific endpoint
type EndpointMetrics struct {
	RequestCount    int64         `json:"request_count"`
	ErrorCount      int64         `json:"error_count"`
	ErrorRate       float64       `json:"error_rate"`
	AverageDuration time.Duration `json:"average_duration"`
}

// Global metrics instance
var globalMetrics = NewMetrics()

// GetGlobalMetrics returns the global metrics instance
func GetGlobalMetrics() *Metrics {
	return globalMetrics
}

// HealthCheck represents system health status
type HealthCheck struct {
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Uptime      time.Duration     `json:"uptime"`
	Version     string            `json:"version"`
	Checks      map[string]string `json:"checks"`
	Metrics     MetricsSummary    `json:"metrics"`
}

// GetHealthCheck returns current health status
func GetHealthCheck(version string) HealthCheck {
	metrics := globalMetrics.GetSummary()
	
	checks := make(map[string]string)
	checks["storage"] = "ok"
	checks["memory"] = "ok"
	
	status := "healthy"
	for _, check := range checks {
		if check != "ok" {
			status = "degraded"
			break
		}
	}
	
	return HealthCheck{
		Status:    status,
		Timestamp: time.Now(),
		Uptime:    metrics.Uptime,
		Version:   version,
		Checks:    checks,
		Metrics:   metrics,
	}
}

// PerformanceProfiler tracks performance bottlenecks
type PerformanceProfiler struct {
	mu sync.RWMutex
	operations map[string][]OperationMetric
}

type OperationMetric struct {
	Duration  time.Duration
	Timestamp time.Time
	Success   bool
	Details   map[string]interface{}
}

// NewPerformanceProfiler creates a new performance profiler
func NewPerformanceProfiler() *PerformanceProfiler {
	return &PerformanceProfiler{
		operations: make(map[string][]OperationMetric),
	}
}

// RecordOperation records an operation for profiling
func (p *PerformanceProfiler) RecordOperation(name string, duration time.Duration, success bool, details map[string]interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	metric := OperationMetric{
		Duration:  duration,
		Timestamp: time.Now(),
		Success:   success,
		Details:   details,
	}
	
	p.operations[name] = append(p.operations[name], metric)
	
	// Keep only last 100 operations per type
	if len(p.operations[name]) > 100 {
		p.operations[name] = p.operations[name][len(p.operations[name])-100:]
	}
}

// GetSlowOperations returns operations that took longer than threshold
func (p *PerformanceProfiler) GetSlowOperations(threshold time.Duration) map[string][]OperationMetric {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	slow := make(map[string][]OperationMetric)
	
	for name, operations := range p.operations {
		for _, op := range operations {
			if op.Duration > threshold {
				slow[name] = append(slow[name], op)
			}
		}
	}
	
	return slow
}

// Global profiler instance
var globalProfiler = NewPerformanceProfiler()

// GetGlobalProfiler returns the global profiler instance
func GetGlobalProfiler() *PerformanceProfiler {
	return globalProfiler
}
