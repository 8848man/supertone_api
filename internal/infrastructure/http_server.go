package infrastructure

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"tts_proxy/internal/interface/handler"
	"tts_proxy/internal/interface/middleware"
)

type ServerConfig struct {
	Port        string
	TTSEndpoint string
	APIVersion  string
}

type HTTPServer struct {
	App *fiber.App
}

func NewHTTPServer(cfg ServerConfig, ttsHandler *handler.TTSHandler, authMiddleware *middleware.AuthMiddleware) *HTTPServer {
	app := fiber.New()

	// CORS 허용
	app.Use(cors.New())

	// 인증 미들웨어(목업) - 향후 확장
	app.Use(authMiddleware.Handle)

	// API 버전별 라우팅 그룹
	apiGroup := app.Group(fmt.Sprintf("/api/%s", cfg.APIVersion))
	
	// TTS 엔드포인트 - voiceID를 URL 경로 파라미터로 받음
	apiGroup.Post(fmt.Sprintf("%s/:voiceId", cfg.TTSEndpoint), ttsHandler.HandleTTS)

	return &HTTPServer{App: app}
}

func (s *HTTPServer) Start(port string) error {
	return s.App.Listen(":" + port)
} 