package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupertoneConfig_DefaultValues(t *testing.T) {
	// 환경 변수 초기화
	os.Unsetenv("SUPERTONE_API_URL")
	os.Unsetenv("SUPERTONE_API_KEY")
	
	config := SupertoneConfig()
	
	assert.Equal(t, SupertoneProvider, config.Provider)
	assert.Equal(t, "https://supertoneapi.com", config.APIURL)
	assert.Equal(t, "5ccd5ef313ccb9aa15795df6a1c03fd8", config.APIKey)
	assert.Equal(t, 30, config.Timeout)
	assert.Equal(t, 3, config.Retries)
}

func TestSupertoneConfig_CustomValues(t *testing.T) {
	// 환경 변수 설정
	os.Setenv("SUPERTONE_API_URL", "https://custom-supertone.com")
	os.Setenv("SUPERTONE_API_KEY", "custom-api-key")
	
	config := SupertoneConfig()
	
	assert.Equal(t, SupertoneProvider, config.Provider)
	assert.Equal(t, "https://custom-supertone.com", config.APIURL)
	assert.Equal(t, "custom-api-key", config.APIKey)
	
	// 환경 변수 정리
	os.Unsetenv("SUPERTONE_API_URL")
	os.Unsetenv("SUPERTONE_API_KEY")
}

func TestLoadTTSConfig_DefaultProvider(t *testing.T) {
	// 환경 변수 초기화
	os.Unsetenv("TTS_PROVIDER")
	
	config := LoadTTSConfig()
	
	assert.Equal(t, SupertoneProvider, config.Provider)
	assert.Equal(t, "https://supertoneapi.com", config.APIURL)
	assert.Equal(t, "5ccd5ef313ccb9aa15795df6a1c03fd8", config.APIKey)
}

func TestLoadTTSConfig_SupertoneProvider(t *testing.T) {
	// 환경 변수 설정
	os.Setenv("TTS_PROVIDER", "supertone")
	
	config := LoadTTSConfig()
	
	assert.Equal(t, SupertoneProvider, config.Provider)
	assert.Equal(t, "https://supertoneapi.com", config.APIURL)
	assert.Equal(t, "5ccd5ef313ccb9aa15795df6a1c03fd8", config.APIKey)
	
	// 환경 변수 정리
	os.Unsetenv("TTS_PROVIDER")
} 