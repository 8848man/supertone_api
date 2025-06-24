package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"tts_proxy/internal/domain"
)

type TTSHandler struct {
	TTSService  domain.TTSService
	AuthService domain.AuthService // 현재는 목업/미사용
}

func NewTTSHandler(ttsService domain.TTSService, authService domain.AuthService) *TTSHandler {
	return &TTSHandler{
		TTSService:  ttsService,
		AuthService: authService,
	}
}

// HandleTTS는 /tts/:voiceId POST 요청을 처리합니다.
func (h *TTSHandler) HandleTTS(c *fiber.Ctx) error {
	// URL 경로에서 voiceID 추출
	voiceID := c.Params("voiceId")
	if voiceID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "voice_id is required in URL path"})
	}

	var req domain.TTSRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	resp, err := h.TTSService.Synthesize(&req, voiceID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "audio/wav")
	return c.Status(http.StatusOK).Send(resp.Audio)
} 