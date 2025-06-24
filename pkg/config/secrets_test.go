package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSecrets(t *testing.T) {
	// 테스트용 임시 디렉토리 생성
	tempDir, err := os.MkdirTemp("", "test_secrets")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 현재 작업 디렉토리 저장
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// 임시 디렉토리로 변경
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// 테스트용 secrets 디렉토리 생성
	secretsDir := filepath.Join(tempDir, "config", "secrets")
	if err := os.MkdirAll(secretsDir, 0755); err != nil {
		t.Fatalf("Failed to create secrets directory: %v", err)
	}

	// 테스트용 secrets 파일 생성
	testSecrets := `{
		"supertone": {
			"api_key": "test_api_key_123",
			"api_url": "https://test-supertoneapi.com"
		},
		"other_provider": {
			"api_key": "test_other_key_456",
			"api_url": "https://test-other-api.com"
		}
	}`

	secretsFile := filepath.Join(secretsDir, "api_keys.json")
	if err := os.WriteFile(secretsFile, []byte(testSecrets), 0644); err != nil {
		t.Fatalf("Failed to write test secrets file: %v", err)
	}

	// LoadSecrets 테스트
	secrets, err := LoadSecrets()
	if err != nil {
		t.Fatalf("LoadSecrets failed: %v", err)
	}

	// 결과 검증
	if secrets.Supertone.APIKey != "test_api_key_123" {
		t.Errorf("Expected Supertone API key 'test_api_key_123', got '%s'", secrets.Supertone.APIKey)
	}

	if secrets.Supertone.APIURL != "https://test-supertoneapi.com" {
		t.Errorf("Expected Supertone API URL 'https://test-supertoneapi.com', got '%s'", secrets.Supertone.APIURL)
	}

	if secrets.OtherProvider.APIKey != "test_other_key_456" {
		t.Errorf("Expected Other Provider API key 'test_other_key_456', got '%s'", secrets.OtherProvider.APIKey)
	}
}

func TestGetProviderConfig(t *testing.T) {
	// 테스트용 임시 디렉토리 생성
	tempDir, err := os.MkdirTemp("", "test_provider_config")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 현재 작업 디렉토리 저장
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// 임시 디렉토리로 변경
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// 테스트용 secrets 디렉토리 생성
	secretsDir := filepath.Join(tempDir, "config", "secrets")
	if err := os.MkdirAll(secretsDir, 0755); err != nil {
		t.Fatalf("Failed to create secrets directory: %v", err)
	}

	// 테스트용 secrets 파일 생성
	testSecrets := `{
		"supertone": {
			"api_key": "test_api_key_123",
			"api_url": "https://test-supertoneapi.com"
		},
		"other_provider": {
			"api_key": "test_other_key_456",
			"api_url": "https://test-other-api.com"
		}
	}`

	secretsFile := filepath.Join(secretsDir, "api_keys.json")
	if err := os.WriteFile(secretsFile, []byte(testSecrets), 0644); err != nil {
		t.Fatalf("Failed to write test secrets file: %v", err)
	}

	// GetProviderConfig 테스트
	apiKey, apiURL, err := GetProviderConfig("supertone")
	if err != nil {
		t.Fatalf("GetProviderConfig failed for supertone: %v", err)
	}

	if apiKey != "test_api_key_123" {
		t.Errorf("Expected Supertone API key 'test_api_key_123', got '%s'", apiKey)
	}

	if apiURL != "https://test-supertoneapi.com" {
		t.Errorf("Expected Supertone API URL 'https://test-supertoneapi.com', got '%s'", apiURL)
	}

	// 알 수 없는 제공자 테스트
	_, _, err = GetProviderConfig("unknown_provider")
	if err == nil {
		t.Error("Expected error for unknown provider, but got none")
	}
} 