package auth

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ConfirmEmailRequest struct {
	VerificationCode string `json:"verificationCode"`
}

type ConfirmEmailResponse struct {
	Message string `json:"message"`
}

func ConfirmEmail(c *gin.Context) {
	var req ConfirmEmailRequest
	var res ConfirmEmailResponse
	authClient := c.MustGet("authClient").(*rpc.AuthClient)
	logger := c.MustGet("logger").(*zap.Logger)
	req.VerificationCode = c.Query("token")
	if req.VerificationCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "empty token"})
		logger.Error("ConfirmEmail: invalid request", zap.String("err", "empty token"))
		return
	}
	logger.Info("ConfirmEmail", zap.Any("request", req))
	var err error
	res.Message, err = authClient.ConfirmEmailRequest(c, req.VerificationCode)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("ConfirmEmail: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": res.Message})
	logger.Info("ConfirmEmail: response", zap.String("message", res.Message))
}
