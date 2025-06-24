package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"tts_proxy/internal/domain"
)

type mockTTSService struct {
	SynthesizeFunc func(req *domain.TTSRequest, voiceID string) (*domain.TTSResponse, error)
}

func (m *mockTTSService) Synthesize(req *domain.TTSRequest, voiceID string) (*domain.TTSResponse, error) {
	return m.SynthesizeFunc(req, voiceID)
}

type mockAuthService struct{}
func (m *mockAuthService) ValidateToken(token string) (string, error) { return "", nil }

func TestHandleTTS_Success(t *testing.T) {
	app := fiber.New()
	mockService := &mockTTSService{
		SynthesizeFunc: func(req *domain.TTSRequest, voiceID string) (*domain.TTSResponse, error) {
			return &domain.TTSResponse{Audio: []byte("MP3DATA"), Format: "mp3"}, nil
		},
	}
	handler := NewTTSHandler(mockService, &mockAuthService{})
	app.Post("/tts/:voiceId", handler.HandleTTS)

	body, _ := json.Marshal(domain.TTSRequest{
		Text:     "hi",
		Language: "en",
		Style:    "neutral",
		Model:    "sona_speech_1",
		VoiceSettings: map[string]interface{}{
			"pitch_shift":    0,
			"pitch_variance": 1,
			"speed":          1,
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/tts/voice-123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	assert.Equal(t, "MP3DATA", buf.String())
	assert.Equal(t, "audio/mpeg", resp.Header.Get("Content-Type"))
}

func TestHandleTTS_BadRequest(t *testing.T) {
	app := fiber.New()
	handler := NewTTSHandler(&mockTTSService{}, &mockAuthService{})
	app.Post("/tts/:voiceId", handler.HandleTTS)

	req := httptest.NewRequest(http.MethodPost, "/tts/voice-123", bytes.NewReader([]byte("notjson")))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHandleTTS_MissingVoiceID(t *testing.T) {
	app := fiber.New()
	handler := NewTTSHandler(&mockTTSService{}, &mockAuthService{})
	app.Post("/tts/:voiceId", handler.HandleTTS)

	body, _ := json.Marshal(domain.TTSRequest{
		Text:     "hi",
		Language: "en",
	})
	req := httptest.NewRequest(http.MethodPost, "/tts/", bytes.NewReader(body)) // voiceId 없음
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode) // Fiber는 경로가 없으면 404 반환
} 