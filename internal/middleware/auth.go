package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// AuthMiddleware provides API key authentication
func AuthMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health check
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}
			
			// Get API key from header
			providedKey := r.Header.Get("X-API-Key")
			if providedKey == "" {
				// Try Authorization header
				auth := r.Header.Get("Authorization")
				if strings.HasPrefix(auth, "Bearer ") {
					providedKey = strings.TrimPrefix(auth, "Bearer ")
				}
			}
			
			// Validate API key
			if providedKey == "" {
				http.Error(w, "API key required", http.StatusUnauthorized)
				return
			}
			
			if providedKey != apiKey {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}
			
			// Add user context
			ctx := context.WithValue(r.Context(), "authenticated", true)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		w.Header().Set("Access-Control-Max-Age", "86400")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		// Log request
		logRequest(r)
		
		// Process request
		next.ServeHTTP(wrapper, r)
		
		// Log response
		logResponse(r, wrapper.statusCode)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func logRequest(r *http.Request) {
	// In a real implementation, use a proper logger
	// log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
}

func logResponse(r *http.Request, statusCode int) {
	// In a real implementation, use a proper logger
	// log.Printf("Response: %s %s -> %d", r.Method, r.URL.Path, statusCode)
}

// RateLimitMiddleware provides basic rate limiting
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	// Simple in-memory rate limiter
	// In production, use Redis or similar
	clients := make(map[string][]int64)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			now := time.Now().Unix()
			
			// Clean old entries
			if requests, exists := clients[clientIP]; exists {
				var validRequests []int64
				for _, timestamp := range requests {
					if now-timestamp < 60 { // Within last minute
						validRequests = append(validRequests, timestamp)
					}
				}
				clients[clientIP] = validRequests
			}
			
			// Check rate limit
			if len(clients[clientIP]) >= requestsPerMinute {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			
			// Add current request
			clients[clientIP] = append(clients[clientIP], now)
			
			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Use remote address
	return strings.Split(r.RemoteAddr, ":")[0]
}

// CompressionMiddleware adds gzip compression
func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client accepts gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		
		// Set compression header
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		
		// In a real implementation, wrap the response writer with gzip
		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		next.ServeHTTP(w, r)
	})
}

// MetricsMiddleware tracks request metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapper, r)
		
		duration := time.Since(start)
		isError := wrapper.statusCode >= 400
		
		// Record metrics (would integrate with monitoring system)
		recordMetrics(r.URL.Path, duration, isError)
	})
}

func recordMetrics(endpoint string, duration time.Duration, isError bool) {
	// In a real implementation, this would send metrics to a monitoring system
	// like Prometheus, DataDog, etc.
}

// TimeoutMiddleware adds request timeout
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic (in real implementation)
				// log.Printf("Panic recovered: %v", err)
				
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}
