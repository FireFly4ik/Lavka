package auth

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type RegisterRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	RedirectURL string `json:"redirect_url"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	var res RegisterResponse
	authClient := c.MustGet("authClient").(*rpc.AuthClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("Register: invalid request", zap.Error(err))
		return
	}
	logger.Info("Register", zap.Any("request", req))
	var err error
	res.Message, err = authClient.RegisterRequest(c, req.Username, req.Email, req.Password, req.RedirectURL)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("Register: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": res.Message})
	logger.Info("Register: response:", zap.String("message", res.Message))
}
