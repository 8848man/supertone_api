package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"tts_proxy/internal/domain"
)

type mockTTSAdapter struct {
	SynthesizeFunc func(req *domain.TTSRequest, voiceID string) (*domain.TTSResponse, error)
}

func (m *mockTTSAdapter) Synthesize(req *domain.TTSRequest, voiceID string) (*domain.TTSResponse, error) {
	return m.SynthesizeFunc(req, voiceID)
}

func TestTTSService_Synthesize_Success(t *testing.T) {
	adapter := &mockTTSAdapter{
		SynthesizeFunc: func(req *domain.TTSRequest, voiceID string) (*domain.TTSResponse, error) {
			return &domain.TTSResponse{Audio: []byte("MP3DATA"), Format: "mp3"}, nil
		},
	}
	service := NewTTSService(adapter)

	req := &domain.TTSRequest{
		Text:     "hello",
		Language: "en",
		Style:    "neutral",
		Model:    "sona_speech_1",
		VoiceSettings: map[string]interface{}{
			"pitch_shift":    0,
			"pitch_variance": 1,
			"speed":          1,
		},
	}
	resp, err := service.Synthesize(req, "voice-123")
	assert.NoError(t, err)
	assert.Equal(t, []byte("MP3DATA"), resp.Audio)
	assert.Equal(t, "mp3", resp.Format)
}

func TestTTSService_Synthesize_Error(t *testing.T) {
	adapter := &mockTTSAdapter{
		SynthesizeFunc: func(req *domain.TTSRequest, voiceID string) (*domain.TTSResponse, error) {
			return nil, errors.New("TTS error")
		},
	}
	service := NewTTSService(adapter)

	req := &domain.TTSRequest{
		Text:    "fail",
		Language: "en",
	}
	resp, err := service.Synthesize(req, "voice-123")
	assert.Error(t, err)
	assert.Nil(t, resp)
} 