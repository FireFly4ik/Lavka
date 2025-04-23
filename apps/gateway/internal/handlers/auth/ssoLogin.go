package auth

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type SSOLoginRequest struct {
	SSOToken string `json:"ssoToken"`
}

func SSOLogin(c *gin.Context) {
	var req SSOLoginRequest
	var res LoginResponse
	authClient := c.MustGet("authClient").(*rpc.AuthClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("SSOLogin: invalid request", zap.Error(err))
		return
	}
	logger.Info("SSOLogin: request", zap.Any("request", req))
	var err error
	res.AccessToken, res.RefreshToken, res.Message, err = authClient.SSOLoginRequest(c, req.SSOToken)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("SSOLogin: response error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": res.Message, "accessToken": res.AccessToken, "refreshToken": res.RefreshToken})
	logger.Info("Login: response", zap.String("message", res.Message))
}
