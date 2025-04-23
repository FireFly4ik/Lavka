package auth

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ForgotPasswordRequest struct {
	Email       string `json:"email"`
	RedirectURL string `json:"redirect_url"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

func ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	var res ForgotPasswordResponse
	authClient := c.MustGet("authClient").(*rpc.AuthClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("ForgotPassword: invalid request", zap.Error(err))
		return
	}
	logger.Info("ForgotPassword", zap.Any("request", req))
	var err error
	res.Message, err = authClient.ForgotPasswordRequest(c, req.Email, req.RedirectURL)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("ForgotPassword: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	logger.Info("ForgotPassword: response:", zap.String("message", res.Message))
	c.JSON(http.StatusCreated, gin.H{"message": res.Message})
}
