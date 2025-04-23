package products_customer

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ProductRequest struct {
	ProductID string `json:"productID"`
}

type ProductResponse struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Image       *string  `json:"image"`
	Price       float32  `json:"price"`
	Discount    float32  `json:"discount"`
	Category    []string `json:"categories"`
	Message     string   `json:"message"`
}

func Product(c *gin.Context) {
	var req ProductRequest
	var res ProductResponse
	proCusClient := c.MustGet("proCusClient").(*rpc.ProCusClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("Product: invalid request", zap.Error(err))
		return
	}
	logger.Info("Product", zap.Any("request", req))
	var err error
	res.Name, res.Description, res.Image, res.Price, res.Discount, res.Category, res.Message, err = proCusClient.ProductRequest(c, req.ProductID)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("Product: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"name": res.Name, "description": res.Description, "image": res.Image, "price": res.Price, "discount": res.Discount, "category": res.Category, "message": res.Message})
	logger.Info("Product: response:", zap.String("message", res.Message))
}
