package main

import (
	"testing"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/config"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
)

func TestAppValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		httpPort    string
		grpcPort    string
		expectError bool
	}{
		{
			name:        "different ports should pass",
			httpPort:    "8080",
			grpcPort:    "9090",
			expectError: false,
		},
		{
			name:        "same ports should fail",
			httpPort:    "8080",
			grpcPort:    "8080",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				config: &config.Config{
					HTTP: config.HTTPConfig{Port: tt.httpPort},
					GRPC: config.GRPCConfig{Port: tt.grpcPort},
				},
			}

			err := app.ValidateConfig()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateConfig() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestCreateStorage(t *testing.T) {
	tests := []struct {
		name        string
		storageType string
		dataDir     string
		expectError bool
	}{
		{
			name:        "memory storage should work",
			storageType: "memory",
			dataDir:     "",
			expectError: false,
		},
		{
			name:        "file storage without datadir should fail",
			storageType: "file",
			dataDir:     "",
			expectError: true,
		},
		{
			name:        "invalid storage type should fail",
			storageType: "invalid",
			dataDir:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.StorageConfig{
				Type:    tt.storageType,
				DataDir: tt.dataDir,
			}

			_, err := createStorage(cfg)
			if (err != nil) != tt.expectError {
				t.Errorf("createStorage() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestAppCreateServers(t *testing.T) {
	// Create a minimal valid app
	app := &App{
		config: &config.Config{
			HTTP: config.HTTPConfig{
				Port:         "8080",
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
			},
			GRPC: config.GRPCConfig{
				Port:           "9090",
				MaxRecvMsgSize: 1024 * 1024,
				MaxSendMsgSize: 1024 * 1024,
			},
			Security: config.SecurityConfig{
				EnableAuth: false,
			},
		},
	}

	// Create storage and service
	store, err := createStorage(config.StorageConfig{Type: "memory"})
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	app.service = persona.NewService(store)

	// Test server creation
	httpServer, grpcServer, err := app.CreateServers()
	if err != nil {
		t.Errorf("CreateServers() error = %v", err)
	}

	if httpServer == nil {
		t.Error("HTTP server should not be nil")
	}

	if grpcServer == nil {
		t.Error("gRPC server should not be nil")
	}
}

func TestAppValidateConfigIntegration(t *testing.T) {
	// Test with a complete config structure
	cfg := &config.Config{
		HTTP: config.HTTPConfig{
			Port:         "8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		GRPC: config.GRPCConfig{
			Port:           "9090",
			MaxRecvMsgSize: 1024 * 1024,
			MaxSendMsgSize: 1024 * 1024,
		},
		Storage: config.StorageConfig{
			Type: "memory",
		},
		Security: config.SecurityConfig{
			EnableAuth: false,
		},
	}

	app := &App{config: cfg}
	err := app.ValidateConfig()
	if err != nil {
		t.Errorf("ValidateConfig() with valid config should not error: %v", err)
	}
}
