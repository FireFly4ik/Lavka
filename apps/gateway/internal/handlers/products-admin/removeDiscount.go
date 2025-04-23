package products_admin

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type RemoveDiscountRequest struct {
	ProductId string `json:"product_id"`
}

type RemoveDiscountResponse struct {
	Message string `json:"message"`
}

func RemoveDiscount(c *gin.Context) {
	var req RemoveDiscountRequest
	var res RemoveDiscountResponse
	proAdmClient := c.MustGet("proAdmClient").(*rpc.ProAdmClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("RemoveDiscount: invalid request", zap.Error(err))
		return
	}
	logger.Info("RemoveDiscount", zap.Any("request", req))
	var err error
	res.Message, err = proAdmClient.RemoveDiscountRequest(c, req.ProductId)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("RemoveDiscount: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": res.Message})
	logger.Info("RemoveDiscount: response:", zap.String("message", res.Message))
}
