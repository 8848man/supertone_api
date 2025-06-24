package infrastructure

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"tts_proxy/internal/domain"
)

type mockRoundTripper struct {
	RoundTripFunc func(req *http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req), nil
}

func TestTTSProxyAdapter_Synthesize_Success(t *testing.T) {
	mockRT := &mockRoundTripper{
		RoundTripFunc: func(req *http.Request) *http.Response {
			// URL이 올바른 형태인지 확인
			expectedURL := "https://supertoneapi.com/v1/text-to-speech/test-voice-123"
			assert.Equal(t, expectedURL, req.URL.String())
			
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader("WAVDATA")),
			}
		},
	}
	adapter := &TTSProxyAdapter{
		config: TTSProxyConfig{APIURL: "https://supertoneapi.com", APIKey: "key"},
		client: &http.Client{Transport: mockRT},
	}

	req := &domain.TTSRequest{
		Text:     "hi",
		Language: "en",
		Style:    "neutral",
		Model:    "sona_speech_1",
		VoiceSettings: map[string]interface{}{
			"pitch_shift":    0,
			"pitch_variance": 1,
			"speed":          1,
		},
	}
	resp, err := adapter.Synthesize(req, "test-voice-123")
	assert.NoError(t, err)
	assert.Equal(t, []byte("MP3DATA"), resp.Audio)
	assert.Equal(t, "mp3", resp.Format)
}

func TestTTSProxyAdapter_Synthesize_MissingVoiceID(t *testing.T) {
	adapter := &TTSProxyAdapter{
		config: TTSProxyConfig{APIURL: "https://supertoneapi.com", APIKey: "key"},
		client: &http.Client{},
	}

	req := &domain.TTSRequest{Text: "hi", Language: "en"}
	resp, err := adapter.Synthesize(req, "") // 빈 voiceID
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "voice_id is required")
	assert.Nil(t, resp)
}

func TestTTSProxyAdapter_Synthesize_Error(t *testing.T) {
	mockRT := &mockRoundTripper{
		RoundTripFunc: func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: ioutil.NopCloser(strings.NewReader("error")),
			}
		},
	}
	adapter := &TTSProxyAdapter{
		config: TTSProxyConfig{APIURL: "https://supertoneapi.com", APIKey: "key"},
		client: &http.Client{Transport: mockRT},
	}

	req := &domain.TTSRequest{
		Text:    "fail",
		Language: "en",
	}
	resp, err := adapter.Synthesize(req, "test-voice-123")
	assert.Error(t, err)
	assert.Nil(t, resp)
} 