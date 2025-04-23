package products_admin

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UpdateProductRequest struct {
	ProductId   string   `json:"product_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Price       float32  `json:"price"`
	Discount    float32  `json:"discount"`
	Category    []string `json:"category"`
}

type UpdateProductResponse struct {
	Message string `json:"message"`
}

func UpdateProduct(c *gin.Context) {
	var req UpdateProductRequest
	var res UpdateProductResponse
	proAdmClient := c.MustGet("proAdmClient").(*rpc.ProAdmClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("UpdateProduct: invalid request", zap.Error(err))
		return
	}
	logger.Info("UpdateProduct", zap.Any("request", req))
	var err error
	res.Message, err = proAdmClient.UpdateProductRequest(c, req.ProductId, req.Name, req.Description, req.Image, req.Price, req.Discount, req.Category)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("UpdateProduct: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": res.Message})
	logger.Info("UpdateProduct: response:", zap.String("message", res.Message))
}
