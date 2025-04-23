package auth

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	var res LoginResponse
	authClient := c.MustGet("authClient").(*rpc.AuthClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("Login: invalid request", zap.Error(err))
		return
	}
	logger.Info("Login", zap.Any("request", req))
	var err error
	res.AccessToken, res.RefreshToken, res.Message, err = authClient.LoginRequest(c.Request.Context(), req.Login, req.Password, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("Login: response error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": res.Message, "accessToken": res.AccessToken, "refreshToken": res.RefreshToken})
	logger.Info("Login: response", zap.String("message", res.Message))
}
