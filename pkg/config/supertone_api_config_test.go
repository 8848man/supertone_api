package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestSecrets(t *testing.T) func() {
	// 테스트용 임시 디렉토리 생성
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// 현재 작업 디렉토리 저장
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// 임시 디렉토리로 변경
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// 테스트용 secrets 디렉토리 생성
	secretsDir := filepath.Join(tempDir, "config", "secrets")
	if err := os.MkdirAll(secretsDir, 0755); err != nil {
		t.Fatalf("Failed to create secrets directory: %v", err)
	}

	// 테스트용 secrets 파일 생성
	testSecrets := `{
		"supertone": {
			"api_key": "your_supertone_api_key_here",
			"api_url": "https://supertoneapi.com"
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

	// 클린업 함수 반환
	return func() {
		os.Chdir(originalWd)
		os.RemoveAll(tempDir)
	}
}

func TestSupertoneConfig_DefaultValues(t *testing.T) {
	cleanup := setupTestSecrets(t)
	defer cleanup()

	// 환경 변수 초기화
	os.Unsetenv("SUPERTONE_API_URL")
	os.Unsetenv("SUPERTONE_API_KEY")
	
	config := SupertoneConfig()
	
	assert.Equal(t, SupertoneProvider, config.Provider)
	assert.Equal(t, "https://supertoneapi.com", config.APIURL)
	assert.Equal(t, "your_supertone_api_key_here", config.APIKey)
	assert.Equal(t, 30, config.Timeout)
	assert.Equal(t, 3, config.Retries)
}

func TestSupertoneConfig_CustomValues(t *testing.T) {
	cleanup := setupTestSecrets(t)
	defer cleanup()

	// 환경 변수 설정 (보안 파일보다 우선순위가 낮음)
	os.Setenv("SUPERTONE_API_URL", "https://custom-supertone.com")
	os.Setenv("SUPERTONE_API_KEY", "custom-api-key")
	
	config := SupertoneConfig()
	
	assert.Equal(t, SupertoneProvider, config.Provider)
	// 보안 파일의 값이 우선됨
	assert.Equal(t, "https://supertoneapi.com", config.APIURL)
	assert.Equal(t, "your_supertone_api_key_here", config.APIKey)
	
	// 환경 변수 정리
	os.Unsetenv("SUPERTONE_API_URL")
	os.Unsetenv("SUPERTONE_API_KEY")
}

func TestLoadTTSConfig_DefaultProvider(t *testing.T) {
	cleanup := setupTestSecrets(t)
	defer cleanup()

	// 환경 변수 초기화
	os.Unsetenv("TTS_PROVIDER")
	
	config := LoadTTSConfig()
	
	assert.Equal(t, SupertoneProvider, config.Provider)
	assert.Equal(t, "https://supertoneapi.com", config.APIURL)
	assert.Equal(t, "your_supertone_api_key_here", config.APIKey)
}

func TestLoadTTSConfig_SupertoneProvider(t *testing.T) {
	cleanup := setupTestSecrets(t)
	defer cleanup()

	// 환경 변수 설정
	os.Setenv("TTS_PROVIDER", "supertone")
	
	config := LoadTTSConfig()
	
	assert.Equal(t, SupertoneProvider, config.Provider)
	assert.Equal(t, "https://supertoneapi.com", config.APIURL)
	assert.Equal(t, "your_supertone_api_key_here", config.APIKey)
	
	// 환경 변수 정리
	os.Unsetenv("TTS_PROVIDER")
} 