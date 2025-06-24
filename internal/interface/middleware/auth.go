package middleware

import (
	"github.com/gofiber/fiber/v2"
	"tts_proxy/internal/domain"
)

type AuthMiddleware struct {
	AuthService domain.AuthService
}

func NewAuthMiddleware(authService domain.AuthService) *AuthMiddleware {
	return &AuthMiddleware{AuthService: authService}
}

// Handle은 인증 미들웨어(목업)로, 실제 인증은 미구현입니다.
func (m *AuthMiddleware) Handle(c *fiber.Ctx) error {
	// 향후: Authorization 헤더에서 토큰 추출 및 검증
	// token := c.Get("Authorization")
	// userID, err := m.AuthService.ValidateToken(token)
	// if err != nil { ... }
	return c.Next()
} 