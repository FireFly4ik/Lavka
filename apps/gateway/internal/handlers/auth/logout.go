package auth

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type LogoutRequest struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

func Logout(c *gin.Context) {
	var req LogoutRequest
	var res LogoutResponse
	authClient := c.MustGet("authClient").(*rpc.AuthClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("Logout: invalid request", zap.Error(err))
		return
	}
	logger.Info("Logout", zap.Any("request", req))
	var err error
	res.Message, err = authClient.LogoutRequest(c, req.Token)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("Logout: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	logger.Info("Logout: response:", zap.String("message", res.Message))
}
