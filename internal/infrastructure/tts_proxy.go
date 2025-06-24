package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"tts_proxy/internal/domain"
)

type TTSProxyConfig struct {
	APIURL string
	APIKey string
}

type TTSProxyAdapter struct {
	config TTSProxyConfig
	client *http.Client
}

func NewTTSProxyAdapter(config TTSProxyConfig) *TTSProxyAdapter {
	return &TTSProxyAdapter{
		config: config,
		client: &http.Client{},
	}
}

// Synthesize는 외부 TTS API에 요청을 전달하고 오디오를 반환합니다.
func (a *TTSProxyAdapter) Synthesize(req *domain.TTSRequest, voiceID string) (*domain.TTSResponse, error) {
	// Supertone API 스펙에 맞는 URL 구성: BASEURL/v1/text-to-speech/{voiceId}?output_format=wav
	if voiceID == "" {
		return nil, errors.New("voice_id is required")
	}
	
	apiURL := fmt.Sprintf("%s/v1/text-to-speech/%s?output_format=wav", a.config.APIURL, voiceID)
	
	// Supertone API 요청 본문 구성 (실제 API 명세에 맞춤)
	payload := map[string]interface{}{
		"text":     req.Text,
		"language": req.Language,
	}
	
	// style이 비어있지 않으면 추가
	if req.Style != "" {
		payload["style"] = req.Style
	}
	
	// model이 비어있지 않으면 추가
	if req.Model != "" {
		payload["model"] = req.Model
	}
	
	// voice_settings가 nil이 아니고 비어있지 않으면 추가
	if req.VoiceSettings != nil && len(req.VoiceSettings) > 0 {
		payload["voice_settings"] = req.VoiceSettings
	}
	
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// 디버깅을 위한 로그 출력
	log.Printf("[DEBUG] API URL: %s", apiURL)
	log.Printf("[DEBUG] Request Body: %s", string(requestBody))

	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-sup-api-key", a.config.APIKey)

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 에러 응답 본문도 읽어서 로그에 출력
		errorBody, _ := io.ReadAll(resp.Body)
		log.Printf("[ERROR] API Error Response: %s", string(errorBody))
		return nil, errors.New("TTS API error: " + resp.Status)
	}

	// WAV 바이너리 데이터 읽기
	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &domain.TTSResponse{
		Audio:  audio,
		Format: "wav",
	}, nil
} 