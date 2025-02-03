package api

import (
	"github.com/YXRRXY/autoAuth/backend/internal/api/handlers"
	"github.com/YXRRXY/autoAuth/backend/internal/api/middleware"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(h *server.Hertz, authHandler *handlers.AuthHandler) {
	// API 版本分组
	v1 := h.Group("/api/v1")

	// 公开路由
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)

	// 需要认证的路由
	auth := v1.Group("/", middleware.Auth())
	{
		auth.POST("/auth/refresh", authHandler.RefreshToken)
		// 其他需要认证的路由...
	}
}
