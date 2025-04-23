package auth

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type RefreshRequest struct {
	Token string `json:"token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Message      string `json:"message"`
}

func Refresh(c *gin.Context) {
	var req RefreshRequest
	var res RefreshResponse
	authClient := c.MustGet("authClient").(*rpc.AuthClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("Refresh: invalid request", zap.Error(err))
		return
	}
	logger.Info("Refresh", zap.Any("request", req))
	var err error
	res.RefreshToken, res.AccessToken, res.Message, err = authClient.RefreshRequest(c, req.Token)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("Refresh: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"accessToken": res.AccessToken, "refreshToken": res.RefreshToken, "message": res.Message})
	logger.Info("Refresh: response:", zap.String("message", res.Message))
}
