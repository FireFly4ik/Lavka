package products_admin

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ApplyDiscountRequest struct {
	ProductId string  `json:"product_id"`
	Discount  float32 `json:"discount"`
}

type ApplyDiscountResponse struct {
	Message string `json:"message"`
}

func ApplyDiscount(c *gin.Context) {
	var req ApplyDiscountRequest
	var res ApplyDiscountResponse
	proAdmClient := c.MustGet("proAdmClient").(*rpc.ProAdmClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("ApplyDiscount: invalid request", zap.Error(err))
		return
	}
	logger.Info("ApplyDiscount", zap.Any("request", req))
	var err error
	res.Message, err = proAdmClient.ApplyDiscountRequest(c, req.ProductId, req.Discount)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("ApplyDiscount: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": res.Message})
	logger.Info("ApplyDiscount: response:", zap.String("message", res.Message))
}
