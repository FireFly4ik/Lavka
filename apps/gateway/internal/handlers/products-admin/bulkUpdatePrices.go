package products_admin

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type BulkUpdateProductsRequest struct {
	Products []UpdateProductRequest `json:"products"`
}

type BulkUpdateProductsResponse struct {
	Message string `json:"message"`
}

func BulkUpdateProducts(c *gin.Context) {
	var req BulkUpdateProductsRequest
	var res BulkUpdateProductsResponse
	proAdmClient := c.MustGet("proAdmClient").(*rpc.ProAdmClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("BulkUpdateProducts: invalid request", zap.Error(err))
		return
	}
	logger.Info("BulkUpdateProducts", zap.Any("request", req))
	for _, product := range req.Products {
		_, err := proAdmClient.UpdateProductRequest(c, product.ProductId, product.Name, product.Description, product.Image, product.Price, product.Discount, product.Category)
		if err != nil {
			if res.Message == "" {
				c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
			}
			logger.Error("BulkUpdateProducts: response error:", zap.String("message", res.Message), zap.Error(err))
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": res.Message})
	logger.Info("BulkUpdateProducts: response:", zap.String("message", res.Message))
}
