package middleware

import (
	"context"
	"strings"

	"github.com/YXRRXY/autoAuth/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Auth 认证中间件
func Auth() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		auth := string(ctx.GetHeader("Authorization"))
		if auth == "" {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "未提供认证令牌",
			})
			ctx.Abort()
			return
		}

		// 检查并提取 Bearer 令牌
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "无效的认证格式",
			})
			ctx.Abort()
			return
		}

		// 验证访问令牌
		claims, err := utils.GetJWTHandler().ValidateAccessToken(parts[1])
		if err != nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "无效的认证令牌: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		// 将用户信息存储在上下文中
		ctx.Set("userID", claims.UserID)
		ctx.Set("username", claims.Username)

		ctx.Next(c)
	}
}
