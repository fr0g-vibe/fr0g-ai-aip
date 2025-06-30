package cli

import (
	"os"
	"testing"
)

func TestGetConfigFromEnv(t *testing.T) {
	// Save original env vars
	originalClientType := os.Getenv("FR0G_CLIENT_TYPE")
	originalStorageType := os.Getenv("FR0G_STORAGE_TYPE")
	originalDataDir := os.Getenv("FR0G_DATA_DIR")
	originalServerURL := os.Getenv("FR0G_SERVER_URL")
	
	// Clean up after test
	defer func() {
		os.Setenv("FR0G_CLIENT_TYPE", originalClientType)
		os.Setenv("FR0G_STORAGE_TYPE", originalStorageType)
		os.Setenv("FR0G_DATA_DIR", originalDataDir)
		os.Setenv("FR0G_SERVER_URL", originalServerURL)
	}()
	
	// Test default config
	os.Unsetenv("FR0G_CLIENT_TYPE")
	os.Unsetenv("FR0G_STORAGE_TYPE")
	os.Unsetenv("FR0G_DATA_DIR")
	os.Unsetenv("FR0G_SERVER_URL")
	
	config := GetConfigFromEnv()
	if config.ClientType != "grpc" {
		t.Errorf("Expected default client type 'grpc', got %s", config.ClientType)
	}
	if config.StorageType != "file" {
		t.Errorf("Expected default storage type 'file', got %s", config.StorageType)
	}
	
	// Test custom env vars
	os.Setenv("FR0G_CLIENT_TYPE", "rest")
	os.Setenv("FR0G_STORAGE_TYPE", "memory")
	os.Setenv("FR0G_DATA_DIR", "/tmp/test")
	os.Setenv("FR0G_SERVER_URL", "http://example.com")
	
	config = GetConfigFromEnv()
	if config.ClientType != "rest" {
		t.Errorf("Expected client type 'rest', got %s", config.ClientType)
	}
	if config.StorageType != "memory" {
		t.Errorf("Expected storage type 'memory', got %s", config.StorageType)
	}
	if config.DataDir != "/tmp/test" {
		t.Errorf("Expected data dir '/tmp/test', got %s", config.DataDir)
	}
	if config.ServerURL != "http://example.com" {
		t.Errorf("Expected server URL 'http://example.com', got %s", config.ServerURL)
	}
}

func TestCreateClient(t *testing.T) {
	tests := []struct {
		name       string
		config     Config
		expectErr  bool
	}{
		{
			name: "local memory client",
			config: Config{
				ClientType:  "local",
				StorageType: "memory",
			},
			expectErr: false,
		},
		{
			name: "local file client",
			config: Config{
				ClientType:  "local",
				StorageType: "file",
				DataDir:     "/tmp/test-personas",
			},
			expectErr: false,
		},
		{
			name: "rest client",
			config: Config{
				ClientType: "rest",
				ServerURL:  "http://localhost:8080",
			},
			expectErr: false,
		},
		{
			name: "grpc client",
			config: Config{
				ClientType: "grpc",
				ServerURL:  "localhost:9090",
			},
			expectErr: false,
		},
		{
			name: "unknown client type",
			config: Config{
				ClientType: "unknown",
			},
			expectErr: true,
		},
		{
			name: "unknown storage type",
			config: Config{
				ClientType:  "local",
				StorageType: "unknown",
			},
			expectErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := createClient(tt.config)
			if (err != nil) != tt.expectErr {
				t.Errorf("createClient() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && client == nil {
				t.Error("Expected client to be created")
			}
		})
	}
}

func TestExecuteWithConfig_Help(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test help command
	os.Args = []string{"fr0g-ai-aip", "help"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err != nil {
		t.Errorf("Expected no error for help command, got %v", err)
	}
}

func TestExecuteWithConfig_NoArgs(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test with no arguments (should show usage)
	os.Args = []string{"fr0g-ai-aip"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err != nil {
		t.Errorf("Expected no error for no args, got %v", err)
	}
}

func TestExecuteWithConfig_UnknownCommand(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test unknown command
	os.Args = []string{"fr0g-ai-aip", "unknown"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for unknown command")
	}
}

func TestExecuteWithConfig_CreateMissingArgs(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test create command with missing arguments
	os.Args = []string{"fr0g-ai-aip", "create"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for create command with missing args")
	}
}

func TestExecuteWithConfig_GetMissingID(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test get command with missing ID
	os.Args = []string{"fr0g-ai-aip", "get"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for get command with missing ID")
	}
}

func TestExecuteWithConfig_UpdateMissingID(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test update command with missing ID
	os.Args = []string{"fr0g-ai-aip", "update"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for update command with missing ID")
	}
}

func TestExecuteWithConfig_DeleteMissingID(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test delete command with missing ID
	os.Args = []string{"fr0g-ai-aip", "delete"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for delete command with missing ID")
	}
}

func TestServeCommand(t *testing.T) {
	err := serveCommand()
	if err == nil {
		t.Error("Expected error for serve command (not implemented)")
	}
}

func TestExecuteWithConfig_CreateSuccess(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test create command with valid arguments
	os.Args = []string{"fr0g-ai-aip", "create", "-name", "Test Expert", "-topic", "Testing", "-prompt", "You are a test expert"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err != nil {
		t.Errorf("Expected no error for valid create command, got %v", err)
	}
}

func TestExecuteWithConfig_ListSuccess(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test list command
	os.Args = []string{"fr0g-ai-aip", "list"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err != nil {
		t.Errorf("Expected no error for list command, got %v", err)
	}
}

func TestCreateClient_FileStorageError(t *testing.T) {
	config := Config{
		ClientType:  "local",
		StorageType: "file",
		DataDir:     "/invalid/path/that/cannot/be/created",
	}
	
	_, err := createClient(config)
	if err == nil {
		t.Error("Expected error for invalid file storage path")
	}
}

func TestGetConfigFromEnv_PathExpansion(t *testing.T) {
	// Save original env vars
	originalDataDir := os.Getenv("FR0G_DATA_DIR")
	defer func() {
		os.Setenv("FR0G_DATA_DIR", originalDataDir)
	}()
	
	// Test relative path expansion
	os.Setenv("FR0G_DATA_DIR", "./test-data")
	
	config := GetConfigFromEnv()
	if config.DataDir == "./test-data" {
		t.Error("Expected relative path to be expanded to absolute path")
	}
}
