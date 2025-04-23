package products_customer

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type CatalogRequest struct {
	Id string `json:"id"`
}

type CatalogResponse struct {
	ProductIds []string `json:"product_ids"`
	Message    string   `json:"message"`
}

func Catalog(c *gin.Context) {
	var req CatalogRequest
	var res CatalogResponse
	proCusClient := c.MustGet("proCusClient").(*rpc.ProCusClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("Catalog: invalid request", zap.Error(err))
		return
	}
	logger.Info("Catalog", zap.Any("request", req))
	var err error
	res.ProductIds, res.Message, err = proCusClient.CategoryRequest(c, req.Id)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("Catalog: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"ids": res.ProductIds, "message": res.Message})
	logger.Info("Catalog: response:", zap.String("message", res.Message))
}
