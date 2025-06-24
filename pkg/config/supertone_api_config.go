package config

// TTSProvider는 TTS API 제공자 타입을 정의합니다.
type TTSProvider string

const (
	SupertoneProvider TTSProvider = "supertone"
	// 향후 다른 제공자 추가 가능
	// AzureProvider    TTSProvider = "azure"
	// GoogleProvider   TTSProvider = "google"
)

// TTSAPIConfig는 TTS API 설정을 추상화합니다.
type TTSAPIConfig struct {
	Provider TTSProvider
	APIURL   string
	APIKey   string
	// 향후 확장 가능한 설정들
	Timeout  int // 초 단위
	Retries  int
}

// SupertoneConfig는 Supertone API 전용 설정을 반환합니다.
func SupertoneConfig() *TTSAPIConfig {
	// 보안 파일에서 API 키와 URL을 읽기
	apiKey, err := GetSupertoneAPIKey()
	if err != nil {
		// 폴백: 환경 변수 사용
		apiKey = getEnvOrDefault("SUPERTONE_API_KEY", "")
	}
	
	apiURL, err := GetSupertoneAPIURL()
	if err != nil {
		// 폴백: 환경 변수 사용
		apiURL = getEnvOrDefault("SUPERTONE_API_URL", "https://supertoneapi.com")
	}

	return &TTSAPIConfig{
		Provider: SupertoneProvider,
		APIURL:   apiURL,
		APIKey:   apiKey,
		Timeout:  30, // 30초
		Retries:  3,
	}
}

// LoadTTSConfig는 환경 변수에 따라 적절한 TTS 설정을 로드합니다.
func LoadTTSConfig() *TTSAPIConfig {
	provider := TTSProvider(getEnvOrDefault("TTS_PROVIDER", "supertone"))
	
	switch provider {
	case SupertoneProvider:
		return SupertoneConfig()
	default:
		// 기본값으로 Supertone 사용
		return SupertoneConfig()
	}
} 