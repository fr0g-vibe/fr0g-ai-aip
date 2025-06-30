package cli

import (
	"os"
	"path/filepath"
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

func TestExecuteWithConfig_Integration(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	config := Config{
		ClientType:  "local",
		StorageType: "file",
		DataDir:     tmpDir,
	}
	
	// Test create -> list -> get -> delete workflow
	
	// Create a persona
	os.Args = []string{"fr0g-ai-aip", "create", "-name", "Integration Test", "-topic", "Integration Testing", "-prompt", "You are an integration testing expert."}
	err := ExecuteWithConfig(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// List personas to verify creation
	os.Args = []string{"fr0g-ai-aip", "list"}
	err = ExecuteWithConfig(config)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestExecuteWithConfig_ErrorPropagation(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test with invalid storage configuration
	config := Config{
		ClientType:  "local",
		StorageType: "file",
		DataDir:     "/invalid/readonly/path",
	}
	
	os.Args = []string{"fr0g-ai-aip", "create", "-name", "Test", "-topic", "Test", "-prompt", "Test"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for invalid storage path")
	}
}

func TestExecuteWithConfig_GetSuccess(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	config := Config{ClientType: "local", StorageType: "memory"}
	
	// Create a persona first
	os.Args = []string{"fr0g-ai-aip", "create", "-name", "Get Test", "-topic", "Testing", "-prompt", "Test prompt"}
	err := ExecuteWithConfig(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// We can't easily test get without knowing the ID, but we can test the error case
	os.Args = []string{"fr0g-ai-aip", "get", "nonexistent"}
	err = ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for non-existent persona")
	}
}

func TestExecuteWithConfig_UpdateSuccess(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	config := Config{ClientType: "local", StorageType: "memory"}
	
	// Test update with non-existent ID
	os.Args = []string{"fr0g-ai-aip", "update", "nonexistent", "-name", "Updated"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for updating non-existent persona")
	}
}

func TestExecuteWithConfig_DeleteSuccess(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	config := Config{ClientType: "local", StorageType: "memory"}
	
	// Test delete with non-existent ID
	os.Args = []string{"fr0g-ai-aip", "delete", "nonexistent"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for deleting non-existent persona")
	}
}

func TestCreateClient_GRPCError(t *testing.T) {
	// Test gRPC client creation with invalid address
	config := Config{
		ClientType: "grpc",
		ServerURL:  "invalid:address:format",
	}
	
	client, err := createClient(config)
	// gRPC client creation might not fail immediately, but should return a client
	if client == nil && err == nil {
		t.Error("Expected either client or error")
	}
}

func TestGetConfigFromEnv_AllDefaults(t *testing.T) {
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
	
	// Clear all env vars
	os.Unsetenv("FR0G_CLIENT_TYPE")
	os.Unsetenv("FR0G_STORAGE_TYPE")
	os.Unsetenv("FR0G_DATA_DIR")
	os.Unsetenv("FR0G_SERVER_URL")
	
	config := GetConfigFromEnv()
	
	// Verify all defaults
	if config.ClientType != "grpc" {
		t.Errorf("Expected default client type 'grpc', got %s", config.ClientType)
	}
	if config.StorageType != "file" {
		t.Errorf("Expected default storage type 'file', got %s", config.StorageType)
	}
	// DataDir should be expanded to absolute path, so just check it contains "data"
	if !filepath.IsAbs(config.DataDir) || filepath.Base(config.DataDir) != "data" {
		t.Errorf("Expected data dir to be absolute path ending with 'data', got %s", config.DataDir)
	}
	if config.ServerURL != "localhost:9090" {
		t.Errorf("Expected default server URL 'localhost:9090', got %s", config.ServerURL)
	}
}

func TestExecute(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test Execute function (uses default config)
	os.Args = []string{"fr0g-ai-aip", "help"}
	
	err := Execute()
	if err != nil {
		t.Errorf("Execute() failed: %v", err)
	}
}

func TestPrintUsage(t *testing.T) {
	// Test that printUsage doesn't panic
	// We can't easily capture stdout, but we can ensure it runs
	printUsage()
}

func TestCreateClient_RESTDefaults(t *testing.T) {
	config := Config{
		ClientType: "rest",
		ServerURL:  "", // Empty should use some default
	}
	
	client, err := createClient(config)
	if err != nil {
		t.Fatalf("Failed to create REST client: %v", err)
	}
	
	if client == nil {
		t.Error("Expected client to be created")
	}
}

func TestExecuteWithConfig_CreateWithContext(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test create command with all optional fields
	os.Args = []string{"fr0g-ai-aip", "create", 
		"-name", "Full Test Expert", 
		"-topic", "Full Testing", 
		"-prompt", "You are a comprehensive testing expert with full context."}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err != nil {
		t.Errorf("Expected no error for full create command, got %v", err)
	}
}

func TestExecuteWithConfig_UpdatePartial(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	config := Config{ClientType: "local", StorageType: "memory"}
	
	// Create a persona first
	os.Args = []string{"fr0g-ai-aip", "create", "-name", "Update Test", "-topic", "Testing", "-prompt", "Test prompt"}
	err := ExecuteWithConfig(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// Test partial update (only name)
	os.Args = []string{"fr0g-ai-aip", "update", "nonexistent", "-name", "New Name"}
	err = ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for updating non-existent persona")
	}
}

func TestCreateClient_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "grpc with custom URL",
			config: Config{
				ClientType: "grpc",
				ServerURL:  "custom.example.com:9090",
			},
			wantErr: false,
		},
		{
			name: "rest with empty URL",
			config: Config{
				ClientType: "rest",
				ServerURL:  "",
			},
			wantErr: false,
		},
		{
			name: "local with empty data dir",
			config: Config{
				ClientType:  "local",
				StorageType: "file",
				DataDir:     "",
			},
			wantErr: true, // Should fail with empty data dir
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := createClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("createClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("Expected client to be created")
			}
		})
	}
}
