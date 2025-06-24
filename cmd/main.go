package main

import (
	"log"

	"tts_proxy/internal/infrastructure"
	"tts_proxy/internal/interface/handler"
	"tts_proxy/internal/interface/middleware"
	"tts_proxy/internal/usecase"
	"tts_proxy/pkg/config"
)

type mockAuthService struct{}
func (m *mockAuthService) ValidateToken(token string) (string, error) { return "", nil }

func main() {
	cfg := config.LoadConfig()
	ttsConfig := config.LoadTTSConfig()

	ttsAdapter := infrastructure.NewTTSProxyAdapter(infrastructure.TTSProxyConfig{
		APIURL: ttsConfig.APIURL,
		APIKey: ttsConfig.APIKey,
	})
	ttsService := usecase.NewTTSService(ttsAdapter)
	authService := &mockAuthService{} // 실제 구현시 대체
	ttsHandler := handler.NewTTSHandler(ttsService, authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	server := infrastructure.NewHTTPServer(infrastructure.ServerConfig{
		Port:        cfg.Port,
		TTSEndpoint: cfg.TTSEndpoint,
		APIVersion:  cfg.APIVersion,
	}, ttsHandler, authMiddleware)
	
	log.Printf("[INFO] Server starting on :%s", cfg.Port)
	log.Printf("[INFO] Using TTS Provider: %s", ttsConfig.Provider)
	log.Printf("[INFO] API Endpoint: /api/%s%s", cfg.APIVersion, cfg.TTSEndpoint)
	if err := server.Start(cfg.Port); err != nil {
		log.Fatalf("[FATAL] Server error: %v", err)
	}
} 