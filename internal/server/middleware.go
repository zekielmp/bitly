package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zekielmp/Bitly/internal/models"
	"github.com/zekielmp/Bitly/internal/utils"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Authoization Beare JWT
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(ctx, "Authorization header required", ctx.Err())
			ctx.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 || tokenString[0] != "bearer" {
			utils.UnauthorizedResponse(ctx, "Invalid authorization bearer format", ctx.Err())
			ctx.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString[1], s.config.Jwt.SecretKey)
		if err != nil {
			utils.UnauthorizedResponse(ctx, "Invalid token", ctx.Err())
			ctx.Abort()
			return

		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)
		ctx.Set("user_role", claims.Role)

		ctx.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("user_role")
		if !exists {
			utils.ForbiddenResponse(ctx, "Forbidden", ctx.Err())
			ctx.Abort()
			return
		}
		if role != string(models.UserRoleAdmin) {
			utils.ForbiddenResponse(ctx, "Forbidden", ctx.Err())
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
