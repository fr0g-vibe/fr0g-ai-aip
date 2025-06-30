package cli

import (
	"os"
	"testing"
)

func TestExecuteWithConfig_EmptyCommand(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test with empty command
	os.Args = []string{"fr0g-ai-aip", ""}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for empty command")
	}
}

func TestExecuteWithConfig_WhitespaceCommand(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test with whitespace command
	os.Args = []string{"fr0g-ai-aip", "   "}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for whitespace command")
	}
}

func TestExecuteWithConfig_CaseSensitivity(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test case sensitivity
	os.Args = []string{"fr0g-ai-aip", "LIST"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for uppercase command")
	}
}

func TestExecuteWithConfig_SpecialCharacters(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test special characters in command
	os.Args = []string{"fr0g-ai-aip", "list@#$"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err == nil {
		t.Error("Expected error for command with special characters")
	}
}

func TestCreateClient_InvalidConfigurations(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "empty client type",
			config: Config{
				ClientType: "",
			},
			wantErr: true,
		},
		{
			name: "invalid characters in client type",
			config: Config{
				ClientType: "local@#$",
			},
			wantErr: true,
		},
		{
			name: "very long client type",
			config: Config{
				ClientType: "verylongclienttypethatdoesnotexist",
			},
			wantErr: true,
		},
		{
			name: "numeric client type",
			config: Config{
				ClientType: "123",
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := createClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("createClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetConfigFromEnv_InvalidValues(t *testing.T) {
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
	
	// Test with invalid values
	os.Setenv("FR0G_CLIENT_TYPE", "invalid@#$")
	os.Setenv("FR0G_STORAGE_TYPE", "invalid@#$")
	os.Setenv("FR0G_DATA_DIR", "")
	os.Setenv("FR0G_SERVER_URL", "invalid-url")
	
	config := GetConfigFromEnv()
	
	// Should still return a config with the invalid values
	// (validation happens in createClient)
	if config.ClientType != "invalid@#$" {
		t.Errorf("Expected client type 'invalid@#$', got %s", config.ClientType)
	}
	if config.StorageType != "invalid@#$" {
		t.Errorf("Expected storage type 'invalid@#$', got %s", config.StorageType)
	}
}

func TestExecuteWithConfig_LongArguments(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test with very long arguments
	longName := make([]byte, 1000)
	for i := range longName {
		longName[i] = 'a'
	}
	
	os.Args = []string{"fr0g-ai-aip", "create", "-name", string(longName), "-topic", "Test", "-prompt", "Test"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err != nil {
		t.Errorf("Should handle long arguments gracefully, got error: %v", err)
	}
}

func TestExecuteWithConfig_UnicodeArguments(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test with Unicode characters
	os.Args = []string{"fr0g-ai-aip", "create", "-name", "测试专家", "-topic", "测试", "-prompt", "你是一个测试专家"}
	
	config := Config{ClientType: "local", StorageType: "memory"}
	err := ExecuteWithConfig(config)
	if err != nil {
		t.Errorf("Should handle Unicode arguments gracefully, got error: %v", err)
	}
}
