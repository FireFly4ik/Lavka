package auth

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type InfoRequest struct {
	Token string `json:"token"`
}

type InfoResponse struct {
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Phone    *string `json:"phone"`
	Avatar   *string `json:"avatar"`
	Message  string  `json:"message"`
}

func ClientInfo(c *gin.Context) {
	var req InfoRequest
	var res InfoResponse
	authClient := c.MustGet("authClient").(*rpc.AuthClient)
	logger := c.MustGet("logger").(*zap.Logger)
	tokenString := c.GetHeader("Authorization")
	tokenString = tokenString[7:]
	req.Token = tokenString
	logger.Info("Client Info", zap.Any("request", req))
	var err error
	res.Email, res.Username, res.Phone, res.Avatar, res.Message, err = authClient.ClientInfoRequest(c, req.Token)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("Client Info: response error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"email": res.Email, "username": res.Username, "phone": &res.Phone, "avatar": &res.Avatar, "message": res.Message})
	logger.Info("Client Info: response", zap.String("message", res.Message))
}
