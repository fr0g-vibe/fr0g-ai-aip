package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	HTTP HTTPConfig `yaml:"http"`
	GRPC GRPCConfig `yaml:"grpc"`
	
	// Storage configuration
	Storage StorageConfig `yaml:"storage"`
	
	// Client configuration
	Client ClientConfig `yaml:"client"`
	
	// Security configuration
	Security SecurityConfig `yaml:"security"`
	
	// Logging configuration
	Logging LoggingConfig `yaml:"logging"`
}

type HTTPConfig struct {
	Port            string        `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	EnableTLS       bool          `yaml:"enable_tls"`
	CertFile        string        `yaml:"cert_file"`
	KeyFile         string        `yaml:"key_file"`
}

type GRPCConfig struct {
	Port            string        `yaml:"port"`
	MaxRecvMsgSize  int           `yaml:"max_recv_msg_size"`
	MaxSendMsgSize  int           `yaml:"max_send_msg_size"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout"`
	EnableTLS       bool          `yaml:"enable_tls"`
	CertFile        string        `yaml:"cert_file"`
	KeyFile         string        `yaml:"key_file"`
}

type StorageConfig struct {
	Type    string `yaml:"type"` // memory, file
	DataDir string `yaml:"data_dir"`
}

type ClientConfig struct {
	Type      string `yaml:"type"`       // local, rest, grpc
	ServerURL string `yaml:"server_url"`
	Timeout   time.Duration `yaml:"timeout"`
}

type SecurityConfig struct {
	EnableAuth bool   `yaml:"enable_auth"`
	APIKey     string `yaml:"api_key"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"` // json, text
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	config := &Config{
		HTTP: HTTPConfig{
			Port:            getEnv("FR0G_HTTP_PORT", "8080"),
			ReadTimeout:     getDurationEnv("FR0G_HTTP_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:    getDurationEnv("FR0G_HTTP_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: getDurationEnv("FR0G_HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
			EnableTLS:       getBoolEnv("FR0G_HTTP_ENABLE_TLS", false),
			CertFile:        getEnv("FR0G_HTTP_CERT_FILE", ""),
			KeyFile:         getEnv("FR0G_HTTP_KEY_FILE", ""),
		},
		GRPC: GRPCConfig{
			Port:              getEnv("FR0G_GRPC_PORT", "9090"),
			MaxRecvMsgSize:    getIntEnv("FR0G_GRPC_MAX_RECV_MSG_SIZE", 4*1024*1024), // 4MB
			MaxSendMsgSize:    getIntEnv("FR0G_GRPC_MAX_SEND_MSG_SIZE", 4*1024*1024), // 4MB
			ConnectionTimeout: getDurationEnv("FR0G_GRPC_CONNECTION_TIMEOUT", 5*time.Second),
			EnableTLS:         getBoolEnv("FR0G_GRPC_ENABLE_TLS", false),
			CertFile:          getEnv("FR0G_GRPC_CERT_FILE", ""),
			KeyFile:           getEnv("FR0G_GRPC_KEY_FILE", ""),
		},
		Storage: StorageConfig{
			Type:    getEnv("FR0G_STORAGE_TYPE", "file"),
			DataDir: getEnv("FR0G_DATA_DIR", "./data"),
		},
		Client: ClientConfig{
			Type:      getEnv("FR0G_CLIENT_TYPE", "grpc"),
			ServerURL: getEnv("FR0G_SERVER_URL", "localhost:9090"),
			Timeout:   getDurationEnv("FR0G_CLIENT_TIMEOUT", 30*time.Second),
		},
		Security: SecurityConfig{
			EnableAuth: getBoolEnv("FR0G_ENABLE_AUTH", false),
			APIKey:     getEnv("FR0G_API_KEY", ""),
		},
		Logging: LoggingConfig{
			Level:  getEnv("FR0G_LOG_LEVEL", "info"),
			Format: getEnv("FR0G_LOG_FORMAT", "text"),
		},
	}
	
	// Expand relative paths
	if !filepath.IsAbs(config.Storage.DataDir) {
		if abs, err := filepath.Abs(config.Storage.DataDir); err == nil {
			config.Storage.DataDir = abs
		}
	}
	
	return config
}

// Validate validates the entire configuration
func (c *Config) Validate() error {
	var errors ValidationErrors
	
	// Validate HTTP config
	if httpErrors := c.validateHTTPConfig(); len(httpErrors) > 0 {
		errors = append(errors, httpErrors...)
	}
	
	// Validate gRPC config
	if grpcErrors := c.validateGRPCConfig(); len(grpcErrors) > 0 {
		errors = append(errors, grpcErrors...)
	}
	
	// Validate storage config
	if storageErrors := c.validateStorageConfig(); len(storageErrors) > 0 {
		errors = append(errors, storageErrors...)
	}
	
	// Validate client config
	if clientErrors := c.validateClientConfig(); len(clientErrors) > 0 {
		errors = append(errors, clientErrors...)
	}
	
	// Validate security config
	if securityErrors := c.validateSecurityConfig(); len(securityErrors) > 0 {
		errors = append(errors, securityErrors...)
	}
	
	// Cross-validation
	if crossErrors := c.validateCrossConfig(); len(crossErrors) > 0 {
		errors = append(errors, crossErrors...)
	}
	
	if len(errors) > 0 {
		return errors
	}
	
	return nil
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
