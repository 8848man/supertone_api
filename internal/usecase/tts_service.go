package usecase

import "tts_proxy/internal/domain"

// TTSAdapter는 외부 TTS API 호출을 추상화합니다.
type TTSAdapter interface {
	Synthesize(req *domain.TTSRequest, voiceID string) (*domain.TTSResponse, error)
}

// ttsService는 TTSService의 실제 구현체입니다.
type ttsService struct {
	adapter TTSAdapter
}

// NewTTSService는 TTSService 구현체를 생성합니다.
func NewTTSService(adapter TTSAdapter) domain.TTSService {
	return &ttsService{adapter: adapter}
}

// Synthesize는 외부 TTSAdapter를 통해 TTS 변환을 수행합니다.
func (s *ttsService) Synthesize(req *domain.TTSRequest, voiceID string) (*domain.TTSResponse, error) {
	return s.adapter.Synthesize(req, voiceID)
} 