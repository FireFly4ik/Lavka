package auth

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ResetPasswordRequest struct {
	VerificationCode string `json:"verificationCode"`
	Password         string `json:"password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	var res ResetPasswordResponse
	authClient := c.MustGet("authClient").(*rpc.AuthClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("ResetPassword: invalid request", zap.Error(err))
		return
	}
	logger.Info("ResetPassword: request", zap.Any("request", req))
	var err error
	res.Message, err = authClient.ResetPasswordRequest(c, req.VerificationCode, req.Password)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("ResetPassword: response error", zap.Error(err))
		return
	}
	logger.Info("ResetPassword: response", zap.String("message", res.Message))
}
