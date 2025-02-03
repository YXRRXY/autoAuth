package handlers

import (
	"context"

	"github.com/YXRRXY/autoAuth/backend/internal/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 处理注册请求
func (h *AuthHandler) Register(c context.Context, ctx *app.RequestContext) {
	var req service.RegisterRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": "无效的请求参数: " + err.Error(),
		})
		return
	}

	if err := h.authService.Register(c, &req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusCreated, map[string]interface{}{
		"message": "注册成功",
	})
}

// Login 处理登录请求
func (h *AuthHandler) Login(c context.Context, ctx *app.RequestContext) {
	var req LoginRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": "无效的请求参数",
		})
		return
	}

	token, user, err := h.authService.Login(c, req.Username, req.Password)
	if err != nil {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(c context.Context, ctx *app.RequestContext) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": "未登录",
		})
		return
	}

	// 使用 service 层获取用户信息
	user, err := h.authService.GetUserByID(c, userID.(uint))
	if err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"error": "用户不存在",
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"user": user,
	})
}

// RefreshToken 处理令牌刷新请求
func (h *AuthHandler) RefreshToken(c context.Context, ctx *app.RequestContext) {
	refreshToken := string(ctx.GetHeader("Refresh-Token"))
	if refreshToken == "" {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": "未提供刷新令牌",
		})
		return
	}

	// 使用刷新令牌生成新的令牌对
	tokenPair, err := h.authService.RefreshToken(c, refreshToken)
	if err != nil {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": "刷新令牌无效或已过期",
		})
		return
	}

	ctx.JSON(consts.StatusOK, tokenPair)
}
