package products_customer

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type DiscountResponse struct {
	ProductsID []string `json:"products_id"`
	Message    string   `json:"message"`
}

func Discount(c *gin.Context) {
	var res DiscountResponse
	proCusClient := c.MustGet("proCusClient").(*rpc.ProCusClient)
	logger := c.MustGet("logger").(*zap.Logger)
	var err error
	res.ProductsID, res.Message, err = proCusClient.DiscountRequest(c)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("Discount: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"ids": res.ProductsID, "message": res.Message})
	logger.Info("Discount: response:", zap.String("message", res.Message))
}
