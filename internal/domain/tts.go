package domain

// TTSRequest는 클라이언트가 전달하는 TTS 요청 데이터입니다.
type TTSRequest struct {
	Text          string                 `json:"text"`
	Language      string                 `json:"language"`
	Style         string                 `json:"style"`         // Supertone API용 스타일 (예: "neutral")
	Model         string                 `json:"model"`         // Supertone API용 모델 (예: "sona_speech_1")
	VoiceSettings map[string]interface{} `json:"voice_settings"` // pitch_shift, pitch_variance, speed 등
}

// TTSResponse는 TTS 변환 결과(오디오 바이너리 등)를 나타냅니다.
type TTSResponse struct {
	Audio []byte
	Format string // 예: "wav"
}

// TTSService는 TTS 변환 유즈케이스를 추상화합니다.
type TTSService interface {
	Synthesize(req *TTSRequest, voiceID string) (*TTSResponse, error)
}

// AuthService는 인증/계정 식별을 추상화합니다.
type AuthService interface {
	ValidateToken(token string) (userID string, err error)
} 