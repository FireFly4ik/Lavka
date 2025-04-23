package products_customer

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type SearchProductRequest struct {
	Prefix      string  `json:"prefix"`
	Category    string  `json:"category"`
	MinPrice    float32 `json:"min_price"`
	MaxPrice    float32 `json:"max_price"`
	HasDiscount bool    `json:"has_discount"`
}

type SearchProductResponse struct {
	ProductIds []string `json:"product_ids"`
	Message    string   `json:"message"`
}

func SearchProduct(c *gin.Context) {
	var req SearchProductRequest
	var res SearchProductResponse
	proCusClient := c.MustGet("proCusClient").(*rpc.ProCusClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("SearchProduct: invalid request", zap.Error(err))
		return
	}
	logger.Info("SearchProduct", zap.Any("request", req))
	var err error
	res.ProductIds, res.Message, err = proCusClient.SearchProductRequest(c, req.Prefix, req.Category, req.MinPrice, req.MaxPrice, req.HasDiscount)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("SearchProduct: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"ids": res.ProductIds, "message": res.Message})
	logger.Info("SearchProduct: response:", zap.String("message", res.Message))
}
