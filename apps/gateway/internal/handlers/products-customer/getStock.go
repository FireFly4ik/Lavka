package products_customer

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type GetStockRequest struct {
	ProductId string `json:"productId"`
	MarketId  string `json:"marketId"`
}

type GetStockResponse struct {
	Stock   int32  `json:"stock"`
	Message string `json:"message"`
}

func GetStock(c *gin.Context) {
	var req GetStockRequest
	var res GetStockResponse
	proCusClient := c.MustGet("proCusClient").(*rpc.ProCusClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("GetStock: invalid request", zap.Error(err))
		return
	}
	logger.Info("GetStock", zap.Any("request", req))
	var err error
	res.Stock, res.Message, err = proCusClient.GetStockRequest(c, req.ProductId, req.MarketId)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("GetStock: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"stock": res.Stock, "message": res.Message})
	logger.Info("GetStock: response:", zap.String("message", res.Message))
}
