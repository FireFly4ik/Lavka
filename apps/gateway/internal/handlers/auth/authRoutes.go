package auth

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/middleware"
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterAuthRoutes(r *gin.Engine, log *zap.Logger, authClient *rpc.AuthClient, jwtKey []byte) {
	authGroup := r.Group("auth")
	authGroup.Use(func(c *gin.Context) {
		c.Set("authClient", authClient)
		c.Set("logger", log)
		c.Next()
	})
	{
		authGroup.POST("/login", Login)
		authGroup.POST("/ssologin", SSOLogin)
		authGroup.POST("/register", Register)
		authGroup.PUT("/refresh", Refresh)
		authGroup.PUT("/resetpassword", ResetPassword)
		authGroup.GET("/confirmemail", ConfirmEmail)
		authGroup.GET("/forgotpassword", ForgotPassword)
		authGroup.DELETE("/logout", Logout)
		authGroup.Group("/", middleware.AuthMiddleware(jwtKey)).GET("/me", ClientInfo)
	}
	log.Info("auth routes were registered")
}
