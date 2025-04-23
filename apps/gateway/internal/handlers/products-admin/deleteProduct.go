package products_admin

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type DeleteProductRequest struct {
	ProductId string `json:"product_id"`
}

type DeleteProductResponse struct {
	Message string `json:"message"`
}

func DeleteProduct(c *gin.Context) {
	var req DeleteProductRequest
	var res DeleteProductResponse
	proAdmClient := c.MustGet("proAdmClient").(*rpc.ProAdmClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("DeleteProduct: invalid request", zap.Error(err))
		return
	}
	logger.Info("DeleteProduct", zap.Any("request", req))
	var err error
	res.Message, err = proAdmClient.DeleteProductRequest(c, req.ProductId)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("DeleteProduct: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": res.Message})
	logger.Info("DeleteProduct: response:", zap.String("message", res.Message))
}
